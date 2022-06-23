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

package edgednscli_test

import (
	cli "github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednscli"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednscli/pb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI test", func() {

	When("DNS CLI Add is called", func() {
		Context("With correct add parameters", func() {
			It("Should pass", func() {

				cliCfg := cli.AppFlags{
					Address: serverTestAddress,
					PKI:     &cliPKI,
				}
				client, err := cli.PbDNSClient(&cliCfg)
				Expect(err).ShouldNot(HaveOccurred())

				rec := "baz.bar.foo.com:1.1.1.1,1.1.1.2,1.1.1.3,1.1.1.4"

				err = cli.AddRecord(rec, pb.RType_A, client)

				Expect(err).ShouldNot(HaveOccurred())

			})
		})
		Context("Correct add, wrong record_type field", func() {
			It("Should fail", func() {

				cliCfg := cli.AppFlags{
					Address: serverTestAddress,
					PKI:     &cliPKI,
				}
				client, err := cli.PbDNSClient(&cliCfg)
				Expect(err).ShouldNot(HaveOccurred())

				rec := "baz.bar.foo.com:1.1.1.1,1.1.1.2,1.1.1.3,1.1.1.4"

				err = cli.AddRecord(rec, pb.RType_None, client)

				Expect(err).Should(HaveOccurred())

			})
		})
		Context("Wrong dnsserver address", func() {
			It("Should fail", func() {

				cliCfg := cli.AppFlags{
					Address: ":1",
					PKI:     &cliPKI,
				}
				client, err := cli.PbDNSClient(&cliCfg)
				Expect(err).Should(HaveOccurred())
				Expect(client).To(BeNil())
			})
		})
	})

	When("DNS CLI Del is called", func() {
		Context("With correct del parameters", func() {
			It("Should pass", func() {
				cliCfg := cli.AppFlags{
					Address: serverTestAddress,
					PKI:     &cliPKI,
				}
				client, err := cli.PbDNSClient(&cliCfg)
				Expect(err).ShouldNot(HaveOccurred())

				fqdn := "baz.bar.foo.com"

				err = cli.DelRecord(fqdn, pb.RType_A, client)
				Expect(err).ShouldNot(HaveOccurred())

			})
		})
		Context("Correct del, wrong record_type field", func() {
			It("Should fail", func() {
				cliCfg := cli.AppFlags{
					Address: serverTestAddress,
					PKI:     &cliPKI,
				}
				client, err := cli.PbDNSClient(&cliCfg)
				Expect(err).ShouldNot(HaveOccurred())

				fqdn := "baz.bar.foo.com"

				err = cli.DelRecord(fqdn, pb.RType_None, client)
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("Wrong address", func() {
			It("Should fail", func() {
				cliCfg := cli.AppFlags{
					Address: ":1",
					PKI:     &cliPKI,
				}
				client, err := cli.PbDNSClient(&cliCfg)
				Expect(err).Should(HaveOccurred())
				Expect(client).To(BeNil())
			})
		})
	})
	When("DNS CLI Get All is called", func() {
		Context("With correct parameters", func() {
			It("Should pass", func() {

				cliCfg := cli.AppFlags{
					Address: serverTestAddress,
					PKI:     &cliPKI,
				}
				client, err := cli.PbDNSClient(&cliCfg)
				Expect(err).ShouldNot(HaveOccurred())
				recs, err := cli.FilterRecords(client, func(fqdn string) bool {
					return true
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(recs).ToNot(BeNil())
				cli.PrintRecs(recs)
			})
		})
	})
	When("DNS CLI Set Forwarders is called", func() {
		Context("With correct parameters", func() {
			It("Should pass", func() {

				cliCfg := cli.AppFlags{
					Address: serverTestAddress,
					PKI:     &cliPKI,
				}
				client, err := cli.PbDNSClient(&cliCfg)
				Expect(err).ShouldNot(HaveOccurred())
				err = cli.SetForwarders("8.8.8.8", client)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
