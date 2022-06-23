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
	"fmt"
	"sync"
	"time"

	rate "github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Limit to allowed burst and per second rate", func() {
	It("events are limited to the burst size when received at max processing rate", func() {
		keys := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
		burstSize := 10
		takeFn, stopFn := rate.Limit(2.0, burstSize)
		defer stopFn()
		counter := make(map[string]int)
		eventsPerKey := 1_000
		By(fmt.Sprintf("creating %d keys with %d elements each", len(keys), eventsPerKey))
		for _, key := range events(eventsPerKey, keys) {
			if takeFn(key) {
				counter[key]++
			}
		}
		Expect(counter).Should(HaveLen(len(keys)))
		for _, key := range keys {
			Expect(counter).Should(HaveKeyWithValue(key, burstSize))
		}
	})

	It("receives events equal to the burst size plus the average rate times time in some time period", func() {
		keys := []string{"1", "2", "3", "4", "5"}
		burstSize := 10
		takeFn, stopFn := rate.Limit(2.0, burstSize)
		defer stopFn()
		counter := make(map[string]int)
		eventsPerKey := 1_000
		By(fmt.Sprintf("creating %d keys with %d elements each", len(keys), eventsPerKey))
		for i, key := range events(eventsPerKey, keys) {
			if takeFn(key) {
				counter[key]++
			}
			if i == eventsPerKey {
				// pause for 2 seconds
				// 10 events should be added per initial burst
				// and then 2 events per second for a total of 14
				Eventually(time.After(2200*time.Millisecond), 3*time.Second).Should(Receive())
			}
		}
		Expect(counter).Should(HaveLen(len(keys)))
		for _, key := range keys {
			Expect(counter).Should(HaveKeyWithValue(key, burstSize+4))
		}
	})

	It("uses concurrent event creators and events are limited to the burst size for each bucket", func() {
		keys := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
		burstSize := 10
		takeFn, stopFn := rate.Limit(2.0, burstSize)
		defer stopFn()
		var mu sync.Mutex
		counter := make(map[string]int)
		eventsPerKey := 1_000
		routines := 2
		By(fmt.Sprintf("creating %d keys with %d elements each in %d goroutines", len(keys), eventsPerKey, routines))
		var wg sync.WaitGroup
		for i := 0; i < routines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, key := range events(eventsPerKey, keys) {
					if takeFn(key) {
						mu.Lock()
						counter[key]++
						mu.Unlock()
					}
				}
			}()
		}
		wg.Wait()

		Expect(counter).Should(HaveLen(len(keys)))
		for _, key := range keys {
			Expect(counter).Should(HaveKeyWithValue(key, burstSize))
		}
	})

	It("event removal function operates correctly", func() {
		const (
			event        = "event"
			nPerSecond   = 0.08325
			burstSize    = 5
			refillPeriod = 13 * time.Second
		)

		// nPerSecond of 0.08325 means limiter refills once per 12s which is 5 per minute
		// burstSize of 5 means limiter is initialized with a capacity of 5
		loginAllowed, loginFailed, cleanup := rate.LimitWithOptionalReplacement(nPerSecond, burstSize)
		defer cleanup()

		Expect(loginAllowed(event)).To(BeTrue())

		By("Removing initial event capacity from limiter")
		for n := 1; n <= burstSize; n++ {
			Expect(loginFailed(event)).To(BeTrue())
		}

		By("Testing that an event is rate limited")
		Expect(loginAllowed(event)).To(BeFalse())
		Expect(loginFailed(event)).To(BeFalse())

		By("Waiting for event capacity to be refilled")
		<-time.After(refillPeriod)

		By("Testing that an event is refilled")
		Expect(loginAllowed(event)).To(BeTrue())
		Expect(loginFailed(event)).To(BeTrue())

		// loginFailed() removed the current event capacity
		By("Testing that an event is rate limited")
		Expect(loginAllowed(event)).To(BeFalse())
		Expect(loginFailed(event)).To(BeFalse())
	})
})

func events(n int, keys []string) []string {
	out := make([]string, n*len(keys))
	j := 0

	for i := 0; i < n; i++ {
		for _, key := range keys {
			// since some tests pause after n events are taken from the output slice,
			// we should add keys in this pattern if n=4: abc,abc,abc,abc and not aaaa,bbbb,cccc
			out[j] = key
			j++
		}
	}

	return out
}
