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

package test

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/pb"
	"google.golang.org/grpc"
)

type rcpClient func(ctx context.Context) error

// ControlClient represents an API client
type ControlClient struct {
	sock *string
	cc   *grpc.ClientConn
	rc   pb.ControlClient
}

// NewControlClient returns a new Client
func NewControlClient(sock *string) *ControlClient {
	if f, err := os.Stat(*sock); err != nil || f.Mode()&os.ModeSocket == 0 {
		fmt.Printf("Invalid API socket: %s %v\n", *sock, err)
		os.Exit(66)
	}

	return &ControlClient{
		sock: sock,
	}
}

// Connect returns a new connected Client
func (c *ControlClient) Connect() error {
	var err error

	c.cc, err = grpc.Dial("unix://"+*c.sock, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Unable to connect: %v", err)
	}
	fmt.Printf("Connected to API socket at %s\n", *c.sock)
	c.rc = pb.NewControlClient(c.cc)
	return nil
}

// Close will terminate the connection to the API
func (c *ControlClient) Close() {
	err := c.cc.Close()
	if err != nil {
		fmt.Printf("Client disconnect error: %s", err)
	} else {
		fmt.Println("Client Disconnected")
	}
}

func (c *ControlClient) callRPC(f rcpClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := f(ctx)
	if err != nil {
		return fmt.Errorf("Unable to perform RPC: %v", err)
	}
	return nil
}

func newHostRecord(rtype pb.RType,
	fqdn string, addrs []string) *pb.HostRecordSet {

	var addrBytes [][]byte
	if len(addrs) > 0 {
		for _, addr := range addrs {
			addrBytes = append(addrBytes, []byte(net.ParseIP(addr)))
		}
	}
	return &pb.HostRecordSet{
		RecordType: rtype,
		Fqdn:       fqdn,
		Addresses:  addrBytes,
	}
}

// SetA sets an A record for a FQDN
func (c *ControlClient) SetA(fqdn string, addrs []string) error {
	fmt.Printf("Setting %d IPv4 address(es) for %s\n", len(addrs), fqdn)
	return c.callRPC(func(ctx context.Context) error {
		_, err := pb.NewControlClient(c.cc).SetAuthoritativeHost(ctx,
			newHostRecord(pb.RType_A, fqdn, addrs))
		return err
	})
}

// GetAll gets all records
func (c *ControlClient) GetAll() (map[string][]string, error) {
	fmt.Println("Getting all records")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	recSets, err := pb.NewControlClient(c.cc).GetAllHosts(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	allRecs := make(map[string][]string)
	for _, recSet := range recSets.RecordSets {
		var addrs []string
		for i, addr := range recSet.Addresses {
			addrs = append(addrs, net.IP(addr).String())
			if i < len(recSet.Addresses)-1 {
				fmt.Printf(", ")
			}
		}
		allRecs[recSet.Fqdn] = addrs
	}

	return allRecs, nil
}

// DeleteA deletes an authoritative entry for given FQDN
func (c *ControlClient) DeleteA(fqdn string) error {
	fmt.Printf("Deleting IPv4 address(es) for %s\n", fqdn)
	return c.callRPC(func(ctx context.Context) error {
		_, err := pb.NewControlClient(c.cc).DeleteAuthoritative(ctx,
			&pb.RecordSet{
				RecordType: pb.RType_A,
				Fqdn:       fqdn,
			})
		return err
	})
}
