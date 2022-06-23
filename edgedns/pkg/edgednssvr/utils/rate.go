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
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// Limit implements the Token Bucket algorithm for many buckets identified by a
// unique string key.
//
// In this variation the only bucket accessor is without replacement.
func Limit(nPerSecond float64, burstSize int) (take func(string) bool, stop func()) {
	_, take, stop = LimitWithOptionalReplacement(nPerSecond, burstSize)
	return
}

// LimitWithOptionalReplacement implements the Token Bucket algorithm for many
// buckets identified by a unique string key.
//
// This variation provides two accessors to the bucket: with and without token
// replacement.
func LimitWithOptionalReplacement(nPerSecond float64, burstSize int) (takeWithReplacement,
	take func(string) bool, stop func()) {
	ctx, cancel := context.WithCancel(context.Background())

	var mu sync.RWMutex

	buckets := make(map[string]func(bool) bool)

	checkForToken := func(withReplacement bool) func(string) bool {
		return func(key string) bool {
			// Load or store new bucket
			mu.RLock()
			take, ok := buckets[key]
			mu.RUnlock()

			if !ok {
				mu.Lock()
				take, ok = buckets[key]

				if !ok {
					take = newBucket(ctx, nPerSecond, burstSize)
					buckets[key] = take
				}
				mu.Unlock()
			}
			// Check token availability
			return take(withReplacement)
		}
	}

	return checkForToken(true), checkForToken(false), cancel
}

// Construct a self-refilling bucket.
func newBucket(ctx context.Context, nPerSecond float64, burstSize int) (take func(withReplacement bool) bool) { // nolint: gocognit,lll
	// Start with a filled bucket
	bkt := make(chan struct{}, burstSize)

	for i := 0; i < burstSize; i++ {
		bkt <- struct{}{}
	}

	// Start a routine to refill the bucket
	go func() {
		period := time.Duration(math.MaxInt64) // infinity
		if nPerSecond > 0 {
			period = time.Duration(float64(time.Second) / nPerSecond)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(period):
				select {
				case bkt <- struct{}{}:
				default:
				}
			}
		}
	}()

	// Return whether a token is immediately available
	return func(withReplacement bool) bool {
		if withReplacement {
			return len(bkt) > 0
		}

		select {
		case <-ctx.Done():
			fmt.Println("case 1")
			return false
		case <-bkt:
			return true
		default:
			return false
		}
	}
}
