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
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/statsdcli"
)

// StatsReporter will ping the stats service at the specified interval
type StatsReporter struct {
	interval time.Duration
}

// NewStatsReporter returns a new StatsReporter with all internals initialized
func NewStatsReporter(interval time.Duration) *StatsReporter {
	return &StatsReporter{
		interval: interval,
	}
}

// Start will send heartbeats to the stats service at the specified interval.
func (sr *StatsReporter) Start(ctx context.Context) {
	for {
		select {
		case <-time.After(sr.interval):
			sr.heartbeat()
		case <-ctx.Done():
			return
		}
	}
}

func (*StatsReporter) heartbeat() {
	stcfg := statsdcli.StatsdConfig{
		Address: statsdip,
		Port:    statsdport,
	}
	cli := statsdcli.NewClient(stcfg)
	if err := cli.Connect(); err != nil {
		fmt.Println("statsd connection failed", err.Error())
		return
	}
	cli.Gauge("status.edgedns..up", "", 1)
	cli.Close()
}
