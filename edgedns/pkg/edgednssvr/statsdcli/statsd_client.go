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

package statsdcli

import (
	"fmt"
	"strconv"

	"github.com/alexcesaro/statsd"
	"github.com/smart-edge-open/edge-services/common/log"
)

// StatsdConfig handles all monitoring configuration
type StatsdConfig struct {
	Address string
	Port    int
}

// Statsd actual client
type Statsd struct {
	cfg    StatsdConfig
	client *statsd.Client
}

// StatsdClient handles all monitoring data
type StatsdClient interface {
	Connect() error
	Gauge(target, action string, value int)
	Close()
}

// Connect used to connect to the actual statsd server
func (r *Statsd) Connect() error {

	client, err := statsd.New(statsd.Address(r.cfg.Address + ":" + strconv.Itoa(r.cfg.Port)))
	if err != nil {
		log.Errf("Could not connect to the stats service: %s", err)
		return err
	}
	r.client = client

	return nil
}

// Gauge used to store metrics
func (r *Statsd) Gauge(target, action string, value int) {
	r.client.Gauge(fmt.Sprintf("status.edgedns.%s.%s", target, action), value)
}

// Close used to close the connection towards the statsd server
func (r *Statsd) Close() {
	r.client.Close()
}

// NewClient used to initialize the client based on configuration
// either to actual statsd or mock
func NewClient(cfg StatsdConfig) StatsdClient {
	if cfg.Address != "" && cfg.Port > 0 {
		return &Statsd{
			cfg:    cfg,
			client: &statsd.Client{},
		}
	}
	return &StatsdMock{
		name: "Stub Statsd Client",
	}
}

// StatsdMock is the mock client
type StatsdMock struct {
	name string
}

// Connect mock method of StatsdMock
func (m *StatsdMock) Connect() error {
	fmt.Println("Stub Statsd Client Connect called")
	return nil
}

// Close mock method of StatsdMock
func (m *StatsdMock) Close() {
	fmt.Println("Stub Statsd Client Close called")
}

// Gauge mock method of StatsdMock
func (m *StatsdMock) Gauge(target, action string, value int) {
	fmt.Println("Stub Statsd Client Gauge called: target:", target, " action:", action, " value:", value)
}
