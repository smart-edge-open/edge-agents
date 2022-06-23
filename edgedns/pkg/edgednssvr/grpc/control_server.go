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

package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"

	edgedns "github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr"

	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/pb"

	"github.com/golang/protobuf/ptypes/empty"
	logger "github.com/smart-edge-open/edge-services/common/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var log = logger.DefaultLogger.WithField("grpc", nil)

// ControlServerPKI defines PKI paths to enable encrypted GRPC server
type ControlServerPKI struct {
	Crt string
	Key string
	Ca  string
}

// ControlServer implements the ControlServer API
type ControlServer struct {
	Sock      string
	Address   string
	PKI       *ControlServerPKI
	server    *grpc.Server
	storage   edgedns.Storage
	responder *edgedns.Responder
}

var _ pb.ControlServer = &ControlServer{}

func readPKI(crtPath, keyPath,
	caPath string) (*credentials.TransportCredentials, error) {

	srvCert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed load server key pair: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(filepath.Clean(caPath))
	if err != nil {
		return nil, fmt.Errorf("Failed read ca certificates: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Errf("Failed appends CA certs from %s", caPath)
		return nil, fmt.Errorf("Failed appends CA certs from %s", caPath)
	}

	creds := credentials.NewTLS(&tls.Config{ // nolint: gosec
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{srvCert},
		ClientCAs:    certPool,
	})

	return &creds, nil
}

func (cs *ControlServer) startIPServer(stg edgedns.Storage) error {
	log.Infof("Starting IP API at %s", cs.Address)
	tc, err := readPKI(filepath.Clean(cs.PKI.Crt),
		filepath.Clean(cs.PKI.Key),
		filepath.Clean(cs.PKI.Ca))
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
			log.Errf("API listener exited unexpectedly: %s", err)
		}
	}()
	return nil
}

func (cs *ControlServer) startSocketServer(stg edgedns.Storage) error {
	log.Infof("Starting socket API at %s", cs.Sock)
	lis, err := net.Listen("unix", cs.Sock)
	if err != nil {
		return fmt.Errorf("Failed to start API listener: %v", err)
	}

	cs.server = grpc.NewServer()
	pb.RegisterControlServer(cs.server, cs)
	go func() {
		if err := cs.server.Serve(lis); err != nil {
			log.Errf("API listener exited unexpectedly: %s", err)
		}
	}()
	return nil
}

// New creates a new ControlServer with storage and a responder. Used for injecting
// a mock storage in tests.
func New(stg edgedns.Storage) *ControlServer {
	return &ControlServer{
		storage: stg,
	}
}

// Start listens on a Unix domain socket only if address is empty.
// If IP address is provided socket file path is ignored and
// server starts to listen on IP address.
func (cs *ControlServer) Start(stg edgedns.Storage, rs *edgedns.Responder) error {
	cs.storage = stg
	cs.responder = rs

	if cs.Address != "" {
		return cs.startIPServer(stg)
	}
	return cs.startSocketServer(stg)
}

// GracefulStop shuts down connetions and removes the Unix domain socket
func (cs *ControlServer) GracefulStop() error {
	cs.server.GracefulStop()
	return nil
}

// SetAuthoritativeHost sets an Authoritative or Forwarder address
// for a given domain
func (cs *ControlServer) SetAuthoritativeHost(ctx context.Context,
	rr *pb.HostRecordSet) (*empty.Empty, error) {

	log.Infof("[API] SetAuthoritativeHost: %s (%d)",
		rr.Fqdn, len(rr.Addresses))
	if rr.RecordType != pb.RType_A {
		return &empty.Empty{}, status.Error(codes.Unimplemented,
			"only A records are supported")
	}

	if rr.Fqdn == "" {
		return &empty.Empty{}, status.Error(codes.InvalidArgument,
			"Fqdn cannot be empty")
	}

	// Lowercase fqdn for lookup purposes. When querying, the fdqn
	// should be lowercased for a proper match.
	fqdn := strings.ToLower(rr.GetFqdn())

	err := cs.storage.SetHostRRSet(uint16(rr.RecordType),
		[]byte(fqdn),
		rr.Addresses)
	if err != nil {
		log.Errf("Failed to set authoritative record: %s", err)
		return &empty.Empty{}, status.Error(codes.Internal,
			"unknown internal DB error occurred")
	}
	return &empty.Empty{}, nil
}

// GetAllHosts returns all records for all domains
func (cs *ControlServer) GetAllHosts(ctx context.Context,
	_ *empty.Empty) (*pb.HostRecordSets, error) {

	log.Infof("[API] GetAllHosts")

	rrs, err := cs.storage.GetAllRRSets()
	if err != nil {
		log.Errf("Failed to get all authoritative records: %s", err)
		return &pb.HostRecordSets{}, status.Error(codes.Internal,
			"unknown internal DB error occurred")
	}

	rs := &pb.HostRecordSets{}
	for fqdn, rr := range rrs {
		rs.RecordSets = append(rs.RecordSets, &pb.HostRecordSet{
			RecordType: pb.RType_A,
			Fqdn:       fqdn,
			Addresses:  rr,
		})
	}

	return rs, nil
}

// DeleteAuthoritative deletes the Resource Record
// for a given Query type and domain
func (cs *ControlServer) DeleteAuthoritative(ctx context.Context,
	rr *pb.RecordSet) (*empty.Empty, error) {

	log.Infof("[API] DeleteAuthoritative: %s", rr.Fqdn)
	if rr.RecordType == pb.RType_None {
		return &empty.Empty{}, status.Error(codes.InvalidArgument,
			"you must specify a record type")
	}
	if err := cs.storage.DelRRSet(uint16(rr.RecordType),
		[]byte(rr.Fqdn)); err != nil {

		log.Errf("Failed to delete authoritative record: %s", err)
		return &empty.Empty{}, status.Error(codes.Internal,
			"unknown internal DB error occurred")
	}
	return &empty.Empty{}, nil
}

// GetForwarders returns the responder's DNS forwarders.
func (cs *ControlServer) GetForwarders(_ context.Context, _ *empty.Empty) (*pb.Forwarders, error) {
	resp := &pb.Forwarders{}

	// TODO: Allocate fixed slice and return
	for _, fwdr := range cs.responder.Forwarders() {
		resp.Addresses = append(resp.Addresses, net.ParseIP(fwdr))
	}

	return resp, nil
}

// SetForwarders sets the responder's DNS forwarders.
func (cs *ControlServer) SetForwarders(_ context.Context, req *pb.Forwarders) (*empty.Empty, error) {
	var fwdrs []string
	for _, addr := range req.Addresses {
		fwdrs = append(fwdrs, net.IP(addr).String())
	}

	if err := cs.responder.SetForwarders(fwdrs); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to set responder forwarders: %s", err.Error())
	}

	log.Infof("Set responder forwarders to %v", fwdrs)

	return &empty.Empty{}, nil
}
