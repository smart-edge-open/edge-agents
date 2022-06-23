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
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednscli"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednscli/pb"
)

const name = "edgedns-client"

var (
	code               int
	fwdrs              string
	addr               string
	pkiCrtPath         string
	pkiKeyPath         string
	pkiCAPath          string
	serverNameOverride string

	addFS       *flag.FlagSet
	delFS       *flag.FlagSet
	getAllFS    *flag.FlagSet
	getPlatFS   *flag.FlagSet
	getClientFS *flag.FlagSet
)

func main() {

	defer func() {
		os.Exit(code)
	}()

	flag.StringVar(&edgednscli.SvcOpts.DNSSocketPath, "socket-path", "/run/edgedns.sock", "Unix domain socket path")
	flag.StringVar(&addr, "address", ":4204", "EdgeDNS API address")
	flag.StringVar(&pkiCrtPath, "cert", "certs/cert.pem", "PKI Cert Path")
	flag.StringVar(&pkiKeyPath, "key", "certs/key.pem", "PKI Key Path")
	flag.StringVar(&pkiCAPath, "ca", "certs/root.pem", "PKI CA Path")
	flag.StringVar(&serverNameOverride, "name", "", "PKI Server Name to override while grpc connection")
	flag.StringVar(&fwdrs, "set-forwarders", "", "set DNS forwarders")

	flag.Parse()

	addFS = flag.NewFlagSet("add", flag.ContinueOnError)
	addFS.StringVar(&edgednscli.RecordOpts.Rec, "A", "", "domain:ip1,ip2,ip3")

	delFS = flag.NewFlagSet("del", flag.ContinueOnError)
	delFS.StringVar(&edgednscli.RecordOpts.Rec, "A", "", "domain")

	getAllFS = flag.NewFlagSet("get-all-records", flag.ContinueOnError)
	getPlatFS = flag.NewFlagSet("get-platform-records", flag.ContinueOnError)
	getClientFS = flag.NewFlagSet("get-client-records", flag.ContinueOnError)

	pki := &edgednscli.PKIPaths{
		CrtPath:            pkiCrtPath,
		KeyPath:            pkiKeyPath,
		CAPath:             pkiCAPath,
		ServerNameOverride: serverNameOverride,
	}
	cfg := edgednscli.AppFlags{
		Sock:    edgednscli.SvcOpts.DNSSocketPath,
		Address: addr,
		PKI:     pki}

	cli, err := edgednscli.PbDNSClient(&cfg)
	if err != nil {
		fmt.Printf("Unable to create Edge DNS client: %v\n", err)
		code = 1
		return
	}
	// If the set-forwarders flag is set without a value, clear the forwarders in Edge DNS.
	if strings.Contains(strings.Join(os.Args, " "), "set-forwarders") {
		if err := edgednscli.SetForwarders(fwdrs, cli); err != nil {
			fmt.Printf("Unable to set DNS forwarders: %v\n", err)
			code = 1
			return
		}
		return
	}
	options(flag.Args(), cli)
}

//options list the main functions of EDNS
func options(arg []string, cli pb.ControlClient) {
	switch arg[0] {
	case "add":
		if err := addFS.Parse(arg[1:]); err == nil {
			if err := edgednscli.AddRecord(edgednscli.RecordOpts.Rec, pb.RType_A, cli); err != nil {
				fmt.Println(err)
				code = 1
			}
		}
	case "del":
		if err := delFS.Parse(arg[1:]); err == nil {
			if err := edgednscli.DelRecord(edgednscli.RecordOpts.Rec, pb.RType_A, cli); err != nil {
				fmt.Println(err)
				code = 1
			}
		}
	case "help":
		usage()
	case "get-all-records":
		filter(cli, true)
	default:
		extraoptions(arg[0], cli)
	}
}

//extraoptions lists some extra functions of EDNS
func extraoptions(arg string, cli pb.ControlClient) {
	switch arg {
	case "get-platform-records":
		filter(cli, false)
	case "get-client-records":
		filter(cli, false)
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unknown command: %s\n", arg)
		usage()
		code = 1
	}
}
func filter(cli pb.ControlClient, getall bool) {
	recs, err := edgednscli.FilterRecords(cli, func(fqdn string) bool {
		if getall {
			return true
		}
		return !strings.HasSuffix(fqdn, ".mec.")
	})
	if err != nil {
		fmt.Println(err)
		code = 1
		return
	}
	edgednscli.PrintRecs(recs)
}
func usage() {
	fmt.Printf("usage %s [GLOBAL_OPTION] CMD [OPTION]\n", name)
	fmt.Println("global options:")
	flag.PrintDefaults()
	fmt.Println(addFS.Name())
	addFS.PrintDefaults()
	fmt.Println(delFS.Name())
	delFS.PrintDefaults()
	fmt.Println(getAllFS.Name())
	getAllFS.PrintDefaults()
	fmt.Println(getPlatFS.Name())
	getPlatFS.PrintDefaults()
	fmt.Println(getClientFS.Name())
	getClientFS.PrintDefaults()
}
