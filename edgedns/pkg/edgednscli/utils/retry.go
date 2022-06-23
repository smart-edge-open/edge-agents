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
	"time"

	logger "github.com/smart-edge-open/edge-services/common/log"
)

// InfiniteRetries has a value of math.MaxUint32
const InfiniteRetries = math.MaxUint32

// Repeatedly executes the given action until it succeeds or is canceled. It retries immediately without
// backing off.
func Repeatedly(ctx context.Context, action func(context.Context) error) error {
	return WithBackOff(
		ctx,
		nil,
		action)
}

// WithBackOff executes the given action until it succeeds or is canceled. It retries according to the
// backoff calculator provided.
func WithBackOff(ctx context.Context, backOffCalc Calculator, action func(context.Context) error) error {
	return UpToNTimes(
		ctx,
		backOffCalc,
		InfiniteRetries,
		action)
}

// EverySecond executes the given action at one-second intervals until it succeeds or is canceled.
func EverySecond(ctx context.Context, action func(context.Context) error) error {
	return UpToNTimes(
		ctx,
		func(_ uint) time.Duration {
			return time.Second
		},
		InfiniteRetries,
		action)
}

// UpToNTimes executes a function up to maxTries times using the provided back off calculator. If the backOffCalc
// is nil, the function is retried immediately. If the function completes successfully or the context is
// canceled, NTimes returns nil. Otherwise, if maxTries is exceeded, it returns an ErrMaxed error.
func UpToNTimes(ctx context.Context, backOffCalc Calculator, maxTries uint, action func(context.Context) error) error {
	var (
		i   uint // explicitly declare as uint to avoid needing multiple int() casts in for loop
		err error
	)

	// Main loop. i starts at 1 and goes up to maxTries.
	for i = 1; i <= maxTries; i++ {
		if maxTries == InfiniteRetries {
			if backOffCalc != nil {
				logger.Infof("Retry: Attempt %d (%s since last attempt)", i, backOffCalc(i-1))
			} else {
				logger.Infof("Retry: Attempt %d (no backoff)", i)
			}
		} else {
			if backOffCalc != nil {
				logger.Infof("Retry: Attempt %d of %d (%s since last attempt)", i, maxTries, backOffCalc(i-1))
			} else {
				logger.Infof("Retry: Attempt %d of %d (no backoff)", i, maxTries)
			}
		}

		// Check to see if context is canceled
		select {
		case <-ctx.Done():
			logger.Warning("Retry: Received cancel signal; command run loop exiting.")
			return ctx.Err()
		default:
			// Just keep moving on
		}

		// Wrap the function to recover from panic errors
		func() {
			defer func() {
				if panicErr := recover(); panicErr != nil {
					err = fmt.Errorf("Retry: Action panicked and run loop recovered with error: %s", panicErr)
				}
			}()

			err = action(ctx)
		}()

		// Successful completion of task
		if err == nil {
			return nil
		}

		// An error occurred in the task, log it and wait appropriate amount of time
		logger.Infof("Retry: Task iteration %d did not complete successfully: %s", i, err)
		if backOffCalc != nil {
			logger.Infof("Waiting %s until retrying", backOffCalc(i))
			// Wait the given back off duration before continuing or until context is canceled
			select {
			case <-ctx.Done():
				logger.Warning("Retry: Received cancel signal; command run loop exiting.")
				return ctx.Err()
			case <-time.After(backOffCalc(i)):
				continue
			}
		} else {
			logger.Info("Retrying immediately")
		}
	}

	// We have exceeded the max attempts
	logger.Errf("Retry: Task errored too many times - aborting.")
	return &ErrMaxed{maxTries}
}
