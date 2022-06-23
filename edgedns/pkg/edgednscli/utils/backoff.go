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
package utils

import (
	"math"
	"time"
)

// Calculator takes the number of times an action has been attempted and calculates the duration to wait before
// trying again.
type Calculator func(numAttempts uint) time.Duration

// Exponential returns a Calculator closure that calculates (2^numAttempts * startInterval) as its duration.
// The duration will be kept at or below the provided ceiling.
func Exponential(startInterval, ceiling time.Duration) Calculator {
	return func(numAttempts uint) time.Duration {
		if numAttempts == math.MaxUint64 {
			return ceiling
		}

		backOff := math.Exp2(float64(numAttempts)) * startInterval.Seconds()

		if backOff > ceiling.Seconds() {
			return ceiling
		}

		return time.Duration(backOff) * time.Second
	}
}

// Network provides a convenience exponential back off calculator that starts at one second and goes up
// to one minute.
func Network() Calculator {
	return Exponential(time.Second, time.Minute)
}

// Fibonacci uses fibonacci numbers to back off by 1-second intervals up to the (n-1)th fibonacci number, then restarts.
// When numAttemps == 0, a 0-second delay is returned.
func Fibonacci(n uint) Calculator {
	return func(numAttempts uint) time.Duration {
		return time.Duration(fib((numAttempts)%n)) * time.Second
	}
}

var fibs = make(map[uint]uint)

func fib(n uint) uint {
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 1
	default:
		if fibs[n] == 0 {
			// use recursion to memoize
			fibs[n] = fib(n-2) + fib(n-1)
		}
	}

	return fibs[n]
}
