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

package edgednssvr

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var (
	trustedIPs = []string{
		"127.0.0.0/8", // IPv4 loopback
		// RFC1918:
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16"}

	ipBlocks []*net.IPNet
)

func init() {

	for _, cidr := range trustedIPs {
		_, ipBlock, err := net.ParseCIDR(cidr)

		if err != nil {
			panic(err)
		}
		ipBlocks = append(ipBlocks, ipBlock)
	}
}

// LoggingHandler logs the DNS message.
func LoggingHandler(next dns.Handler) dns.Handler {
	return dns.HandlerFunc(func(w dns.ResponseWriter, q *dns.Msg) {
		log.Debugf("Lookup query for %s from client %s", q.Question[0].Name, w.RemoteAddr().String())
		next.ServeDNS(w, q)
	})
}

// UnsupportedHandler checks the OP code and IP address
// if it is unsupported request it replies with an error.
func UnsupportedHandler(next dns.Handler) dns.Handler {
	return dns.HandlerFunc(func(w dns.ResponseWriter, q *dns.Msg) {
		if q.Opcode != dns.OpcodeQuery {
			log.Noticef("Received unsupported DNS Opcode %s", dns.OpcodeToString[q.Opcode])

			m := &dns.Msg{}
			m.SetRcode(q, dns.RcodeRefused)

			if err := w.WriteMsg(m); err != nil {
				log.Errf("Failed to reply to client: %s", err)
			}
		}

		if !validIP(w.RemoteAddr().String()) {
			// Close the connection of an invalid source IP
			log.Errf("Invalid source ip")
			_ = w.Close()
			return
		}

		next.ServeDNS(w, q)
	})
}

// LimitRequestRate for client connection.
func LimitRequestRate(allowed func(string) bool, next dns.Handler) dns.Handler {
	return dns.HandlerFunc(func(w dns.ResponseWriter, q *dns.Msg) {

		var id string
		// Check for Non-Existent Domain request to prevent
		// multi-client DoS attacks on invalid domains.
		if q.Rcode == dns.RcodeNameError {
			id = "NXDOMAIN"
		} else {
			id = w.RemoteAddr().String()
		}

		// Check to see if the client has exceeded the rate limit
		if !allowed(id) {
			_ = w.Close()
			return
		}

		next.ServeDNS(w, q)
	})
}

// AuthorityHandler performs an authoritative lookup.
func (r *Responder) AuthorityHandler(next dns.Handler) dns.Handler {
	return dns.HandlerFunc(func(w dns.ResponseWriter, q *dns.Msg) {
		// If the query is not type A, pass the lookup to the external forwarding servers
		if q.Question[0].Qtype != dns.TypeA {
			log.Debugf("Unsupported query type %s for %s", q.Question[0].Qtype, q.Question[0].Name)
			next.ServeDNS(w, q)
			return
		}

		rrs, err := r.storage.GetRRSet(q.Question[0].Name, q.Question[0].Qtype)
		if err != nil {
			if strings.Contains(err.Error(), "no authoritative records found") {
				log.Debugf("No authoritative records found for %s", q.Question[0].Name)
			} else {
				log.Errf("Error getting authoritative records for %s from storage: %s", q.Question[0].Name, err.Error())
			}
			next.ServeDNS(w, q)
			return
		}

		// Only randomize answers if number of answers is greater than 1
		if len(*rrs) > 1 {
			shuffle(*rrs)
		}

		m := &dns.Msg{}
		m.SetReply(q)
		m.Authoritative = true
		m.Answer = *rrs

		if err := w.WriteMsg(m); err != nil {
			log.Errf("Failed to reply to client: %s", err)
		}
	})
}

// IterativeForwardingHandler queries each name server, in order, and returns the first answer it receives.
func (r *Responder) IterativeForwardingHandler(next dns.Handler) dns.Handler {
	return dns.HandlerFunc(func(w dns.ResponseWriter, q *dns.Msg) {
		forwardRequestIterative := func() (*dns.Msg, error) {
			errs := make(map[string]error)

			for _, nsrvr := range r.Forwarders() {
				// If the request takes longer than 5 seconds, cancel and proceed to the next name server.
				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				defer cancel()

				m, err := r.exchanger.Exchange(ctx, q, nsrvr)
				if err != nil {
					log.Debug(err)
					errs[nsrvr] = err
					continue
				}

				return m, nil
			}

			return nil, fmt.Errorf("unable to resolve: %s: forwarder errors: %v", q.Question[0].Name, errs)
		}

		metricVal := 0
		defer func() {
			r.metric("lookup", q.Question[0].Name, metricVal)
		}()

		m, err := forwardRequestIterative()
		if err != nil {
			log.Errf("Failed to find answer: %s", err)
			next.ServeDNS(w, q)
			return
		}

		if err := w.WriteMsg(m); err != nil {
			log.Errf("Failed to reply to client: %s", err)
		}

		metricVal = 1
	})
}

// errorHandler replies with a SERVFAIL error.
func errorHandler() dns.Handler {
	return dns.HandlerFunc(func(w dns.ResponseWriter, q *dns.Msg) {
		m := &dns.Msg{}

		m.SetReply(q)
		m.SetRcode(q, dns.RcodeServerFailure)

		if err := w.WriteMsg(m); err != nil {
			log.Errf("Failed to reply to client: %s", err)
		}
	})
}

// Shuffle the order of byte arrays, allowing DNS answers to be randomized.
func shuffle(rrs []dns.RR) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	for n := len(rrs); n > 1; n-- {
		randIndex := r.Intn(n)
		rrs[n-1], rrs[randIndex] = rrs[randIndex], rrs[n-1]
	}
}

// validIP check if the source address is a trusted IP.
func validIP(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	IP := net.ParseIP(host)

	for _, ipBlock := range ipBlocks {

		if ipBlock.Contains(IP) {
			return true
		}
	}
	return false
}

//adding more trusted ips
func addtrustedIPs(trustedips []string) {
	for _, ip := range trustedips {
		if ip != "" {
			_, ipBlock, err := net.ParseCIDR(ip)

			if err != nil {
				log.Errf("Failed to add trusted ip %s : %s", ip, err)
			} else {
				ipBlocks = append(ipBlocks, ipBlock)
				log.Infof("Added trusted ip: %s", ip)
			}
		}

	}
}
