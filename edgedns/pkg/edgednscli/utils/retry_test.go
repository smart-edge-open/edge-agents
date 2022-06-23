// INTEL CONFIDENTIAL
//
// Copyright 2021-2021 Intel Corporation.
//
// This software and the related documents are Intel copyrighted materials, and your use of
// them is governed by the express license under which they were provided to you ("License").
// Unless the License provides otherwise, you may not use, modify, copy, publish, distribute,
// disclose or transmit this software or the related documents without Intel's prior written permission.
//
// This software and the related documents are provided as is, with no express or implied warranties,
// other than those that are expressly stated in the License.
package utils_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednscli/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Retry", func() {
	var (
		attemptCount int
	)

	BeforeEach(func() {
		attemptCount = 0
	})

	Describe("Repeatedly", func() {
		Context("Success", func() {
			It("Should retry repeatedly without backing off and succeed", func() {
				count := 0
				Expect(
					utils.Repeatedly(
						context.Background(),
						func(ctx context.Context) error {
							count++
							if count == 5 {
								return nil
							}
							return errors.New("Repeatedly error")
						}),
				).To(BeNil())

				Expect(count).To(Equal(5))
			})
		})

		Context("Failure", func() {
			It("Should retry repeatedly without backing off until it is canceled", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				var (
					count = 0
					done  = make(chan bool, 1)
					done2 = make(chan bool, 1)
					err   error
				)
				go func() {
					err = utils.Repeatedly(
						ctx,
						func(ctx context.Context) error {
							count++
							if count == 5 {
								done <- true
							}
							return errors.New("Repeatedly error")
						})
					done2 <- true
				}()

				<-done
				Expect(count).To(BeNumerically(">=", 5))
				cancel()
				<-done2
				Expect(err).To(Equal(context.Canceled))
			})
		})
	})

	Describe("WithBackOff", func() {
		Context("Success", func() {
			It("Should retry repeatedly with backoff and succeed", func() {
				var (
					count = 0
					start = time.Now()
				)
				Expect(
					utils.WithBackOff(
						context.Background(),
						func(numAttempts uint) time.Duration {
							switch numAttempts {
							case 0:
								return 0 * time.Millisecond
							case 1:
								return 10 * time.Millisecond
							case 2:
								return 10 * time.Millisecond
							case 3:
								return 20 * time.Millisecond
							case 4:
								return 30 * time.Millisecond
							default:
								Fail("Too many retry attempts!")
							}
							panic("shouldn't be here")
						},
						func(ctx context.Context) error {
							count++
							if count == 5 {
								return nil
							}
							return errors.New("WithBackOff error")
						}),
				).To(BeNil())

				Expect(count).To(Equal(5))
				Expect(time.Since(start)).To(BeNumerically("<", 100*time.Millisecond))
			})
		})

		Context("Failure", func() {
			It("Should retry repeatedly with backoff until it is canceled", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				var (
					count = 0
					done  = make(chan bool, 1)
					done2 = make(chan bool, 1)
					err   error
					start = time.Now()
				)
				go func() {
					err = utils.WithBackOff(
						ctx,
						func(numAttempts uint) time.Duration {
							switch numAttempts {
							case 0:
								return 0 * time.Millisecond
							case 1:
								return 10 * time.Millisecond
							case 2:
								return 10 * time.Millisecond
							case 3:
								return 20 * time.Millisecond
							case 4:
								return 30 * time.Millisecond
							default:
								return 50 * time.Millisecond
							}
						},
						func(ctx context.Context) error {
							count++
							if count == 5 {
								done <- true
							}
							return errors.New("WithBackOff error")
						})
					done2 <- true
				}()

				<-done
				Expect(count).To(Equal(5))
				cancel()
				<-done2
				Expect(err).To(Equal(context.Canceled))
				Expect(time.Since(start)).To(BeNumerically("<", 100*time.Millisecond))
			})
		})
	})

	Describe("EverySecond", func() {
		Context("Success", func() {
			It("Should retry every second and succeed", func() {
				var (
					count = 0
					start = time.Now()
				)
				Expect(
					utils.EverySecond(
						context.Background(),
						func(ctx context.Context) error {
							count++
							if count == 2 {
								return nil
							}
							return errors.New("EverySecond error")
						}),
				).To(BeNil())

				Expect(count).To(Equal(2))
				Expect(time.Since(start)).To(BeNumerically(">=", 1*time.Second))
			})
		})

		Context("Failure", func() {
			It("Should retry every second until it is canceled", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				var (
					count = 0
					done  = make(chan bool, 1)
					done2 = make(chan bool, 1)
					err   error
					start = time.Now()
				)
				go func() {
					err = utils.EverySecond(
						ctx,
						func(ctx context.Context) error {
							count++
							if count == 2 {
								done <- true
							}
							return errors.New("WithBackOff error")
						})
					done2 <- true
				}()

				<-done
				Expect(count).To(Equal(2))
				cancel()
				<-done2
				Expect(err).To(Equal(context.Canceled))
				Expect(time.Since(start)).To(BeNumerically(">=", 1*time.Second))
			})
		})
	})

	Describe("UpToNTimes", func() {
		Context("Success", func() {
			It("Should retry up to N times and succeed", func() {
				count := 0

				Expect(
					utils.UpToNTimes(
						context.Background(),
						nil,
						5,
						func(ctx context.Context) error {
							count++
							if count == 5 {

								return nil
							}
							return errors.New("UpToNTimes error")
						}),
				).To(BeNil())

				Expect(count).To(Equal(5))
			})
		})

		Context("Failure", func() {
			It("Should retry up to N times and return ErrMaxed", func() {
				count := 0

				Expect(
					utils.UpToNTimes(
						context.Background(),
						nil,
						5,
						func(ctx context.Context) error {
							count++
							return errors.New("UpToNTimes error")
						}),
				).To(Equal(&utils.ErrMaxed{5}))

				Expect(count).To(Equal(5))
			})
		})
	})

	Context("Using a binary exponential backoff with ceiling", func() {
		It("Should produce expected durations", func() {
			boc := utils.Exponential(time.Second, 128*time.Second)

			var boTests = []struct {
				attempt  uint          // input
				expected time.Duration // expected result
			}{
				{0, 1 * time.Second},
				{1, 2 * time.Second},
				{2, 4 * time.Second},
				{3, 8 * time.Second},
				{4, 16 * time.Second},
				{5, 32 * time.Second},
				{6, 64 * time.Second},
				{7, 128 * time.Second},
				{8, 128 * time.Second}, // successive attempts will be limited to the ceiling
				{9, 128 * time.Second}, // successive attempts will be limited to the ceiling
			}

			for _, td := range boTests {
				Expect(boc(td.attempt)).To(Equal(td.expected))
			}
		})

		It("Should handle overflow durations", func() {
			boc := utils.Exponential(time.Second, 128*time.Second)
			Expect(boc(1069395)).To(Equal(128 * time.Second))
		})
	})

	Context("With finite retries", func() {
		boc := utils.Exponential(time.Microsecond, time.Microsecond)
		count := func(tempCtx context.Context) error {
			attemptCount = attemptCount + 1
			return fmt.Errorf("Error for testing retry policy")
		}
		succeed := func(tempCtx context.Context) error {
			attemptCount = attemptCount + 1
			return nil
		}

		It("Should be able to complete on first try", func() {
			err := utils.UpToNTimes(context.Background(), boc, 2, succeed)
			Expect(err).NotTo(HaveOccurred()) // The function should complete successfully
			Expect(attemptCount).To(Equal(1)) // The function should only run once
		})

		It("Should reach maximum attempts", func() {
			err := utils.UpToNTimes(context.Background(), boc, 5, count)
			Expect(err).To(HaveOccurred())
			Expect(attemptCount).To(Equal(5))
		})
	})

	Context("With infinite retries", func() {
		infinite := func(tempCtx context.Context) error {
			attemptCount = attemptCount + 1
			if attemptCount > 50 {
				return nil
			}
			return fmt.Errorf("Error for testing retry policy")
		}

		It("Should run 'infinite' retries", func() {
			boc := utils.Exponential(time.Microsecond, 100*time.Microsecond)
			err := utils.UpToNTimes(context.Background(), boc, utils.InfiniteRetries, infinite)
			Expect(err).NotTo(HaveOccurred())
			Expect(attemptCount).To(Equal(51))
		})
	})

	Context("With a cancellable function", func() {
		infinite := func(tempCtx context.Context) error {
			attemptCount = attemptCount + 1
			return fmt.Errorf("Error for testing retry policy")
		}

		It("Should be able to stop correctly", func() {
			boc := utils.Exponential(time.Second, time.Minute)
			started := make(chan struct{})
			finished := make(chan struct{})

			var err error
			ctx, cancel := context.WithCancel(context.Background())

			go func() { // Start the retry routine
				defer GinkgoRecover()

				close(started)
				err = utils.UpToNTimes(ctx, boc, utils.InfiniteRetries, infinite)
				close(finished)
			}()

			Eventually(started).Should(BeClosed(), "Test function failed to start")
			time.Sleep(100 * time.Millisecond)
			cancel()
			Eventually(finished, time.Second).Should(BeClosed(), "Function failed to stop in time after context cancel")

			Expect(err).To(MatchError(ctx.Err()))
			Expect(attemptCount).NotTo(Equal(0))
		})
	})
})

var _ = Describe("Backoff", func() {
	Describe("Network", func() {
		It("Should back off exponentially up to a minute", func() {
			fmt.Println("STEP: creating the calculator")
			calc := utils.Network()
			Expect(calc).ToNot(BeNil())

			fmt.Println("STEP: Verifying that attempt  0 returns a  1-second delay")
			Expect(calc(0)).To(Equal(1 * time.Second))

			fmt.Println("STEP: Verifying that attempt  1 returns a  2-second delay")
			Expect(calc(1)).To(Equal(2 * time.Second))

			fmt.Println("STEP: Verifying that attempt  2 returns a  4-second delay")
			Expect(calc(2)).To(Equal(4 * time.Second))

			fmt.Println("STEP: Verifying that attempt  3 returns a  8-second delay")
			Expect(calc(3)).To(Equal(8 * time.Second))

			fmt.Println("STEP: Verifying that attempt  4 returns a 16-second delay")
			Expect(calc(4)).To(Equal(16 * time.Second))

			fmt.Println("STEP: Verifying that attempt  5 returns a 32-second delay")
			Expect(calc(5)).To(Equal(32 * time.Second))

			fmt.Println("STEP: Verifying that attempts 6+ return a 60-second delay")
			Expect(calc(6)).To(Equal(60 * time.Second))
			Expect(calc(7)).To(Equal(60 * time.Second))
			Expect(calc(10)).To(Equal(60 * time.Second))
			Expect(calc(54)).To(Equal(60 * time.Second))
			Expect(calc(1000)).To(Equal(60 * time.Second))
			Expect(calc(math.MaxUint64 - 1)).To(Equal(60 * time.Second))
			Expect(calc(math.MaxUint64)).To(Equal(60 * time.Second))
		})
	})

	Describe("Fibonacci", func() {
		It("Should back off using fibonacci up to n = 10, then restart", func() {
			fmt.Println("STEP: creating the calculator")
			calc := utils.Fibonacci(11)
			Expect(calc).ToNot(BeNil())

			fmt.Println("STEP: Verifying that attempt  0 returns a  0-second delay")
			Expect(calc(0)).To(Equal(0 * time.Second))

			fmt.Println("STEP: Verifying that attempt  1 returns a  1-second delay")
			Expect(calc(1)).To(Equal(1 * time.Second))

			fmt.Println("STEP: Verifying that attempt  3 returns a  2-second delay")
			Expect(calc(3)).To(Equal(2 * time.Second))

			fmt.Println("STEP: Verifying that attempt  4 returns a  3-second delay")
			Expect(calc(4)).To(Equal(3 * time.Second))

			fmt.Println("STEP: Verifying that attempt  5 returns a  5-second delay")
			Expect(calc(5)).To(Equal(5 * time.Second))

			fmt.Println("STEP: Verifying that attempt  6 returns an 8-second delay")
			Expect(calc(6)).To(Equal(8 * time.Second))

			fmt.Println("STEP: Verifying that attempt  7 returns a 13-second delay")
			Expect(calc(7)).To(Equal(13 * time.Second))

			fmt.Println("STEP: Verifying that attempt  8 returns a 21-second delay")
			Expect(calc(8)).To(Equal(21 * time.Second))

			fmt.Println("STEP: Verifying that attempt  9 returns a 34-second delay")
			Expect(calc(9)).To(Equal(34 * time.Second))

			fmt.Println("STEP: Verifying that attempt 10 returns a 55-second delay")
			Expect(calc(10)).To(Equal(55 * time.Second))

			fmt.Println("STEP: Verifying that attempt 11 returns a  0-second delay")
			Expect(calc(11)).To(Equal(0 * time.Second))

			fmt.Println("STEP: Verifying that attempt 12 returns a  1-second delay")
			Expect(calc(12)).To(Equal(1 * time.Second))

			fmt.Println("STEP: Verifying boundary values")
			Expect(calc(math.MaxUint64 - 1)).To(Equal(2 * time.Second))
			Expect(calc(math.MaxUint64)).To(Equal(3 * time.Second))
		})
	})
})
