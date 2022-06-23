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
package storage_test

import (
	"fmt"
	"os"

	"github.com/miekg/dns"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/storage"
)

var _ = Describe("BoltDB Storage", func() {

	var stg *storage.BoltDB
	BeforeEach(func() {
		f := fmt.Sprintf("unit_%d.db", config.GinkgoConfig.ParallelNode)
		stg = &storage.BoltDB{
			Filename: f,
		}
	})

	AfterEach(func() {
		Expect(stg.Stop()).To(Succeed())
		Expect(os.Remove(stg.Filename)).To(Succeed())
	})

	It("Handles unsupported types", func() {
		Expect(stg.Start()).To(Succeed())
		err := stg.DelRRSet(dns.TypeAVC, []byte("foo.example.com"))
		Expect(err).NotTo(BeNil())
	})

	It("Gets all record sets", func() {
		Expect(stg.Start()).To(Succeed())
		Expect(stg.SetHostRRSet(dns.TypeA, []byte("foobar.com"), [][]byte{[]byte("1.2.3.4")})).To(Succeed())
		Expect(stg.SetHostRRSet(dns.TypeA, []byte("bazork.com"), [][]byte{[]byte("2.3.4.5"),
			[]byte("3.3.3.3")})).To(Succeed())

		rs, err := stg.GetAllRRSets()
		Expect(err).To(BeNil())

		Expect(rs).To(Equal(map[string][][]byte{
			"foobar.com.": {[]byte("1.2.3.4")},
			"bazork.com.": {[]byte("2.3.4.5"), []byte("3.3.3.3")},
		}))
	})
})
