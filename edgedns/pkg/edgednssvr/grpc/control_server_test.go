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
package grpc_test

import (
	"context"

	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/grpc"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/mock"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/pb"

	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("gRPC server", func() {
	var (
		ctx context.Context
		cs  *grpc.ControlServer
		ms  mock.StorageMock
	)

	BeforeEach(func() {
		ctx = context.Background()
		cs = grpc.New(&ms)
	})

	It("Sets an authoritative host", func() {
		ms.On("SetHostRRSet",
			uint16(pb.RType_A),
			[]byte("foobar.com"),
			[][]byte{[]byte("1.2.3.4")},
		).Return(nil)

		_, err := cs.SetAuthoritativeHost(ctx, &pb.HostRecordSet{
			RecordType: pb.RType_A,
			Fqdn:       "foobar.com",
			Addresses:  [][]byte{[]byte("1.2.3.4")},
		})
		Expect(err).To(BeNil())
	})

	It("Gets all hosts", func() {
		ms.On("GetAllRRSets").Return(map[string][][]byte{
			"foobar.com": {[]byte("1.2.3.4")},
		}, nil)

		rs, err := cs.GetAllHosts(ctx, &empty.Empty{})
		Expect(err).To(BeNil())
		Expect(rs).To(Equal(&pb.HostRecordSets{
			RecordSets: []*pb.HostRecordSet{
				{
					RecordType: pb.RType_A,
					Fqdn:       "foobar.com",
					Addresses:  [][]byte{[]byte("1.2.3.4")},
				},
			},
		}))
	})

	It("Deletes an authoritative host", func() {
		ms.On("DelRRSet",
			uint16(pb.RType_A),
			[]byte("foobar.com"),
		).Return(nil)

		_, err := cs.DeleteAuthoritative(ctx, &pb.RecordSet{
			RecordType: pb.RType_A,
			Fqdn:       "foobar.com",
		})
		Expect(err).To(BeNil())
	})
})
