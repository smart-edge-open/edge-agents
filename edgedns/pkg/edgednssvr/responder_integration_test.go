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
package edgednssvr_test

import (
	"fmt"

	"github.com/miekg/dns"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/mock"
	client "github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/test"
)

// Send a DNS query to the test server
func query(d string, t uint16) (msg *dns.Msg, err error) {
	ns := fmt.Sprintf("127.0.0.1:%d", eport+config.GinkgoConfig.ParallelNode)
	dnsClient := new(dns.Client)
	q := &dns.Msg{}
	q.SetQuestion(d, t)

	msg, _, err = dnsClient.Exchange(q, ns)
	return msg, err
}

// Extract IP addresses as string values from a DNS response
func parseAnswers(m *dns.Msg) ([]string, error) {
	q := m.Question[0]
	var addrs []string
	switch q.Qtype {
	case dns.TypeA:
		for _, i := range m.Answer {
			ans, ok := i.(*dns.A)
			if !ok {
				return nil, fmt.Errorf("IPv4 Answer is not an A record")
			}
			addrs = append(addrs, ans.A.String())
		}
	default:
		return nil, fmt.Errorf("Unknown type: %s", q.String())
	}

	return addrs, nil
}

var _ = Describe("Responder", func() {

	var apiClient *client.ControlClient
	var msg *dns.Msg
	var err error

	BeforeEach(func() {
		sock := fmt.Sprintf("dns_%d.sock", config.GinkgoConfig.ParallelNode)
		apiClient = client.NewControlClient(&sock)

		Expect(dnsServer.SetForwarders([]string{mock.ForwarderBad, mock.ForwarderGood})).To(Succeed())
	})

	It("Sets authoritative A records", func(done Done) {
		Expect(apiClient.Connect()).To(Succeed())
		defer apiClient.Close()

		addrsIn := []string{"1.2.6.7", "3.4.5.6", "7.8.4.1"}

		Expect(apiClient.SetA("baz.foo.com", addrsIn)).To(Succeed())

		msg, err = query("baz.foo.com.", dns.TypeA)
		Expect(err).NotTo(HaveOccurred())

		var addrsOut []string
		addrsOut, err = parseAnswers(msg)
		Expect(err).NotTo(HaveOccurred())
		Expect(addrsOut).Should(HaveLen(3))

		Expect(addrsOut).Should(ConsistOf(addrsIn))
		Expect(exchanger.ForwarderCh).ShouldNot(Receive())

		close(done)
	})

	// TODO figure out \x encoding and deal with any existing records added
	// by other tests.
	It("Gets all authoritative A records", func(done Done) {
		Expect(apiClient.Connect()).To(Succeed())
		defer apiClient.Close()

		addrsIn := []string{"1.2.6.7", "3.4.5.6", "7.8.4.1"}

		Expect(apiClient.SetA("baz.foo.com", addrsIn)).To(Succeed())

		records, err := apiClient.GetAll() // nolint: govet
		Expect(err).NotTo(HaveOccurred())
		Expect(records).To(Equal(map[string][]string{
			"baz.foo.com.": {"1.2.6.7", "3.4.5.6", "7.8.4.1"},
		}))

		close(done)
	})

	It("Deletes authoritative A records", func(done Done) {
		Expect(apiClient.Connect()).To(Succeed())
		defer apiClient.Close()

		Expect(apiClient.SetA("baz.bar.foo.com",
			[]string{"42.24.42"})).To(Succeed())

		msg, err = query("baz.bar.foo.com.", dns.TypeA)
		Expect(err).NotTo(HaveOccurred())
		Expect(msg.Rcode).Should(Equal(dns.RcodeSuccess))
		Expect(exchanger.ForwarderCh).ShouldNot(Receive())

		Expect(apiClient.DeleteA("baz.bar.foo.com")).To(Succeed())

		exchanger.SetErr(true)

		msg, err = query("baz.bar.foo.com.", dns.TypeA)
		Expect(err).NotTo(HaveOccurred())
		Expect(msg.Rcode).Should(Equal(dns.RcodeServerFailure))
		Expect(exchanger.ForwarderCh).ShouldNot(Receive())

		close(done)
	})

	It("Randomizes query results", func(done Done) {
		Expect(apiClient.Connect()).To(Succeed())
		defer apiClient.Close()

		addrsIn := []string{"1.42.6.7", "3.42.5.6", "7.8.42.1"}

		Expect(apiClient.SetA("rnd.foo.com", addrsIn)).To(Succeed())

		var rcnt int
		for j := 1; j < 6; j++ {
			msg, err = query("rnd.foo.com.", dns.TypeA)
			Expect(err).NotTo(HaveOccurred())

			var addrsOut []string
			addrsOut, err = parseAnswers(msg)
			Expect(err).NotTo(HaveOccurred())

			Expect(addrsOut).Should(HaveLen(3))
			Expect(exchanger.ForwarderCh).ShouldNot(Receive())

			for i, v := range addrsOut {
				if v != addrsIn[i] {
					rcnt++
				}
			}
		}
		fmt.Printf("Queries randomized %d times\n", rcnt)
		Expect(rcnt).Should(BeNumerically(">", 2))

		close(done)
	})

	It("Returns SERVFAIL for unanswerable queries", func(done Done) {
		exchanger.SetErr(true)

		msg, err = query("oblivion.dev.null.", dns.TypeA)
		Expect(err).NotTo(HaveOccurred())
		Expect(msg.Rcode).Should(Equal(dns.RcodeServerFailure))
		Expect(exchanger.ForwarderCh).ShouldNot(Receive())

		close(done)
	})

	It("Does not allow multiple questions in a query", func(done Done) {
		ns := fmt.Sprintf("127.0.0.1:%d",
			eport+config.GinkgoConfig.ParallelNode)
		dnsClient := new(dns.Client)

		m := &dns.Msg{}
		m.Question = []dns.Question{
			{
				Name:   "a.b.c.",
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET},
			{
				Name:   "e.f.g.",
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET},
		}

		resp, _, err := dnsClient.Exchange(m, ns)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Rcode).Should(Equal(dns.RcodeFormatError))
		Expect(exchanger.ForwarderCh).ShouldNot(Receive())

		close(done)
	})

	It("Only allows queries", func(done Done) {
		ns := fmt.Sprintf("127.0.0.1:%d",
			eport+config.GinkgoConfig.ParallelNode)
		dnsClient := new(dns.Client)

		m := &dns.Msg{}
		m.SetNotify("example.org.")
		soa, _ := dns.NewRR("example.org. IN SOA sns.dns.icann.org." +
			"noc.dns.icann.org. 2018112827 7200 3600 1209600 3600")
		m.Answer = []dns.RR{soa}

		resp, _, err := dnsClient.Exchange(m, ns)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Rcode).Should(Equal(dns.RcodeRefused))
		Expect(exchanger.ForwarderCh).ShouldNot(Receive())

		close(done)
	})

	It("Delegates non-authoritative queries", func(done Done) {
		exchanger.SetErr(false)

		msg, err := query("google.com.", dns.TypeA)
		Expect(err).NotTo(HaveOccurred())
		Expect(msg.Answer).NotTo(BeEmpty())
		Expect(<-exchanger.ForwarderCh).To(Equal(mock.ForwarderGood))

		close(done)
	})

	It("Iterates over forwarders", func(done Done) {
		exchanger.SetErr(false)

		msg, err := query("google.com.", dns.TypeA)
		Expect(err).NotTo(HaveOccurred())
		Expect(msg.Answer).NotTo(BeEmpty())

		// The first lookup using the first forwarder (ForwarderBad) should fail. The lookup with the second forwarder
		// (ForwarderGood) should always succed.
		Expect(<-exchanger.ForwarderCh).To(Equal(mock.ForwarderGood))

		close(done)
	})
})
