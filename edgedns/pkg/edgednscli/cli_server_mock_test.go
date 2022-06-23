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
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednscli/pb"
	"google.golang.org/grpc"
)

// ControlServerPKI defines PKI paths to enable encrypted GRPC server
type ControlServerPKI struct {
	Crt string
	Key string
	CA  string
}

// ControlServer implements the ControlServer API
type ControlServer struct {
	Address string
	PKI     *ControlServerPKI
	server  *grpc.Server
}

func (cs *ControlServer) StartServer() error {
	fmt.Println("Starting IP API at: ", cs.Address)

	tc, err := readTestPKICredentials(filepath.Clean(cs.PKI.Crt),
		filepath.Clean(cs.PKI.Key),
		filepath.Clean(cs.PKI.CA))
	if err != nil {
		return fmt.Errorf("failed to read pki: %v", err)
	}

	lis, err := net.Listen("tcp", cs.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	cs.server = grpc.NewServer(grpc.Creds(*tc))
	pb.RegisterControlServer(cs.server, cs)
	go func() {
		if err := cs.server.Serve(lis); err != nil {
			fmt.Printf("API listener exited unexpectedly: %s", err)
		}
	}()
	return nil
}

// GracefulStop shuts down connetions and removes the Unix domain socket
func (cs *ControlServer) GracefulStop() error {
	cs.server.GracefulStop()
	return nil
}

// SetAuthoritativeHost is a mock representation of regular server part of
// 'SetAuthoritativeHost' API function. It sets fileds of a internal struct
// 'setRequestl' which can be used to examine the correctness of cli messages
// inside of UT.
func (cs *ControlServer) SetAuthoritativeHost(ctx context.Context,
	rr *pb.HostRecordSet) (*empty.Empty, error) {

	fmt.Printf("[Test Server] SetAuthoritativeHost: %s %s %v",
		rr.RecordType, rr.Fqdn, rr.Addresses)

	return &empty.Empty{}, nil
}

// DeleteAuthoritative is a mock representation of regular server part of
// 'DeleteAuthoritative' API function. It sets fileds of a internal struct
// 'delRequest' which can be used to examine the correctness of cli messages
// inside of UT.
func (cs *ControlServer) DeleteAuthoritative(ctx context.Context,
	rr *pb.RecordSet) (*empty.Empty, error) {

	fmt.Printf("[Test Server] DeleteAuthoritative: [%s %s]",
		rr.RecordType, rr.Fqdn)

	return &empty.Empty{}, nil
}

// GetAllHosts is a mock representation of regular server part of
// 'GetAllHosts' API function. It sets fileds of a internal struct
// 'getAllReq' which can be used to examine the correctness of cli messages
// inside of UT.
func (cs *ControlServer) GetAllHosts(context.Context, *empty.Empty) (*pb.HostRecordSets, error) {

	addr := make([][]byte, 1)
	addr[0] = []byte("1.1.1.1")

	r := &pb.HostRecordSet{
		RecordType: pb.RType_A,
		Fqdn:       "test.com",
		Addresses:  addr,
	}
	rs := make([]*pb.HostRecordSet, 1)
	rs[0] = r
	rrs := &pb.HostRecordSets{
		RecordSets: rs,
	}
	fmt.Printf("[Test Server] GetAllHosts: %s %s %v",
		r.RecordType, r.Fqdn, r.Addresses)

	return rrs, nil
}

func (cs *ControlServer) GetForwarders(context.Context, *empty.Empty) (*pb.Forwarders, error) {
	addr := make([][]byte, 1)
	addr[0] = []byte("8.8.8.8")
	fwdrs := &pb.Forwarders{
		Addresses: addr,
	}
	fmt.Printf("[Test Server] GetForwarders: %v",
		fwdrs.Addresses)
	return fwdrs, nil
}
func (cs *ControlServer) SetForwarders(context.Context, *pb.Forwarders) (*empty.Empty, error) {
	addr := make([][]byte, 1)
	addr[0] = []byte("8.8.8.8")
	fwdrs := &pb.Forwarders{
		Addresses: addr,
	}

	fmt.Printf("[Test Server] SetForwarders: %v",
		fwdrs.Addresses)
	return &empty.Empty{}, nil
}
