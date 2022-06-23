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

package dns

import (
	"context"
	"fmt"

	"github.com/miekg/dns"
	logger "github.com/smart-edge-open/edge-services/common/log"
)

var log = logger.DefaultLogger.WithField("exchanger", nil)

// NetExchanger resolves DNS queries using the network to reach external name servers.
type NetExchanger struct {
	DNS *dns.Client
}

// Exchange synchonously resolves a DNS query.
func (ex *NetExchanger) Exchange(_ context.Context, q *dns.Msg, addr string) (*dns.Msg, error) {
	m, rtt, err := ex.DNS.Exchange(q, addr+":53")
	if err != nil {
		return nil, fmt.Errorf("unable to resolve %s: %s: %w", addr, q.Question[0].Name, err)
	}
	log.Debugf("Lookup %s from upstream %s query took %v", q.Question[0].Name, addr, rtt)

	return m, nil
}
