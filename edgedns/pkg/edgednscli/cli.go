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

package edgednscli

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednscli/pb"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednscli/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// PKIPaths defines paths to files needed to create encrypted grpc connection
type PKIPaths struct {
	CrtPath            string
	KeyPath            string
	CAPath             string
	ServerNameOverride string
}

// AppFlags defines config flags set during sturtup of app
type AppFlags struct {
	Sock    string
	Address string
	PKI     *PKIPaths
}

var SvcOpts struct {
	DNSSocketPath string
}

var RecordOpts struct {
	Rec string
}

const grpcDialTimeoutSec = 1

func readPKI(cfg *AppFlags) (*credentials.TransportCredentials, error) {

	ca, err := ioutil.ReadFile(cfg.PKI.CAPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Printf("Append certs failed from %s", cfg.PKI.CAPath)
		return nil, fmt.Errorf("append certs failed from %s",
			cfg.PKI.CAPath)
	}

	srvCert, err := tls.LoadX509KeyPair(cfg.PKI.CrtPath, cfg.PKI.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load server key pair: %v", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		ServerName:   cfg.PKI.ServerNameOverride,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{srvCert},
		RootCAs:      certPool,
	})

	return &creds, nil
}

func SetForwarders(fwdrs string, cli pb.ControlClient) error {
	var forwarders [][]byte
	if len(fwdrs) != 0 {
		addrs := strings.Split(fwdrs, " ")

		for _, addr := range addrs {
			// Skip setting the DNS service to itself
			if addr == "192.168.216.2" {
				continue
			}

			ip := net.ParseIP(addr)
			if ip == nil {
				return fmt.Errorf("%s is not a valid IP address", addr)
			}

			forwarders = append(forwarders, ip)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return utils.WithBackOff(
		ctx,
		utils.Network(),
		func(ctx context.Context) error {
			if _, err := cli.SetForwarders(ctx, &pb.Forwarders{Addresses: forwarders}); err != nil {
				return err
			}
			return nil
		},
	)
}

func AddRecord(rec string, rt pb.RType, c pb.ControlClient) error {

	recs, err := validateAdd(rec, rt)
	if err != nil {
		return fmt.Errorf("could not validate add record parameters: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := c.SetAuthoritativeHost(ctx, recs); err != nil {
		return fmt.Errorf("could not add record set for %s: %v", recs.Fqdn, err)
	}
	fmt.Printf("added (%d) addresses for record for %s\n", len(recs.Addresses), recs.Fqdn)
	return nil
}

func DelRecord(fqdn string, rt pb.RType, c pb.ControlClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	recSet, err := validateDel(fqdn, rt)
	if err != nil {
		return fmt.Errorf("could not validate del record parameters: %v", err)
	}
	if _, err := c.DeleteAuthoritative(ctx, recSet); err != nil {
		return fmt.Errorf("could not delete record %s: %v", recSet.Fqdn, err)
	}
	fmt.Printf("deleted record for %s\n", recSet.Fqdn)
	return nil
}

func FilterRecords(c pb.ControlClient, display func(fqdn string) bool) (map[string][][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hostRecSets, err := c.GetAllHosts(ctx, &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not get all records: %v", err)
	}

	filtered := make(map[string][][]byte)
	for _, hostRecSet := range hostRecSets.RecordSets {
		if display(hostRecSet.Fqdn) {
			filtered[hostRecSet.Fqdn] = hostRecSet.Addresses
		}
	}
	return filtered, err
}

func PrintRecs(recs map[string][][]byte) {
	for fqdn, addrs := range recs {
		fmt.Printf("%s: ", fqdn)
		for i, addr := range addrs {
			fmt.Printf("%s", net.IP(addr))
			if i < len(addrs)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Println()
	}
}

//func PbDNSClient(sockPath string) (pb.ControlClient, error) {
func PbDNSClient(cfg *AppFlags) (pb.ControlClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), grpcDialTimeoutSec*time.Second)
	defer cancel()
	var dnsConn *grpc.ClientConn
	var err error
	if cfg.Address == "" {
		dialer := func(ctx context.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", addr)
		}

		dnsConn, err = grpc.DialContext(ctx, cfg.Sock,
			grpc.WithBlock(),
			grpc.WithInsecure(),
			grpc.WithContextDialer(dialer),
		)

		if err != nil {
			return nil, errors.Wrapf(err, "problem dialing service at: %s", cfg.Sock)
		}
		fmt.Printf("gRPC connection established to edgedns service at %s\n", cfg.Sock)
	} else {
		tc, err := readPKI(cfg)
		if err != nil {
			return nil, fmt.Errorf("PKI failure: %v", err)
		}

		fmt.Printf("Connecting to EdgeDNS server(%s)", cfg.Address)

		dnsConn, err = grpc.DialContext(ctx, cfg.Address,
			grpc.WithTransportCredentials(*tc), grpc.WithBlock())
		if err != nil {
			return nil, fmt.Errorf("failed to grpc dial: %v", err)
		}
		fmt.Printf("gRPC connection established to edgedns service at %s\n", cfg.Address)
	}

	return pb.NewControlClient(dnsConn), nil
}

func validateAdd(rec string, rt pb.RType) (*pb.HostRecordSet, error) {
	t := strings.Split(rec, ":")
	if len(t) != 2 {
		return nil, errors.New("record must be in domain:ip format")
	}
	fmt.Printf("add record for %s\n", RecordOpts.Rec)

	if rt != pb.RType_A {
		return nil, errors.New("RecordType needs to be set to \"A\" as default")
	}
	addrs := strings.Split(t[1], ",")
	parsedaddrs, err := parseAddresses(addrs)
	if err != nil {
		return nil, fmt.Errorf("DNS address translation failure: %v", err)
	}
	recs := &pb.HostRecordSet{
		RecordType: rt,
		Fqdn:       t[0],
		Addresses:  parsedaddrs,
	}
	return recs, nil
}
func validateDel(fqdn string, rt pb.RType) (*pb.RecordSet, error) {

	if rt != pb.RType_A {
		return nil, errors.New("RecordType needs to be set to \"A\" as default")
	}

	recSet := &pb.RecordSet{
		RecordType: rt,
		Fqdn:       fqdn,
	}

	return recSet, nil
}
func parseAddresses(addresses []string) ([][]byte, error) {
	var outputByteSlice [][]byte
	for _, ipString := range addresses {
		ip := net.ParseIP(ipString)
		if ip == nil {
			return outputByteSlice,
				fmt.Errorf("wrong IP address provided: %s", ipString)
		}
		outputByteSlice = append(outputByteSlice, ip)
	}
	return outputByteSlice, nil
}
