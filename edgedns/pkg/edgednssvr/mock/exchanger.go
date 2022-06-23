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

package mock

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/miekg/dns"
)

const (
	// ForwarderBad is a DNS forwarder that always fails.
	ForwarderBad = "1.2.3.4"
	// ForwarderGood is a DNS forwarder that always succeeds.
	ForwarderGood = "5.6.7.8"
)

// Exchanger implements edgedns.Exchanger.
type Exchanger struct {
	// If err is:
	// - True: Exchange returns an error
	// - False: Exchange returns a reply message with one A record
	err bool
	mux sync.RWMutex

	// ForwarderCh tracks the forwarders addresses used by the Exchange method.
	ForwarderCh chan string
}

// Exchange returns a message or an error.
func (ex *Exchanger) Exchange(_ context.Context, q *dns.Msg, addr string) (*dns.Msg, error) {
	ex.mux.RLock()
	defer ex.mux.RUnlock()

	if ex.err {
		return nil, fmt.Errorf("ERROR ERROR ERROR")
	}

	if addr == ForwarderBad {
		return nil, fmt.Errorf("ERROR BAD FORWARDER")
	}

	// Record the forwarder address
	ex.ForwarderCh <- addr

	reply := &dns.Msg{}
	reply.SetReply(q)
	reply.Answer = []dns.RR{
		&dns.A{
			Hdr: dns.RR_Header{
				Name:   q.Question[0].Name,
				Rrtype: q.Question[0].Qtype,
				Class:  dns.ClassINET,
				Ttl:    10,
			},
			A: net.IP{},
		},
	}

	return reply, nil
}

// SetErr sets ex.err.
func (ex *Exchanger) SetErr(e bool) {
	ex.mux.Lock()
	ex.err = e
	ex.mux.Unlock()
}
