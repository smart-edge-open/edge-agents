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
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	edgedns "github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/dns"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/grpc"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/statsdcli"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/storage"

	d "github.com/miekg/dns"
	logger "github.com/smart-edge-open/edge-services/common/log"
)

var (
	log = logger.DefaultLogger.WithField("main", nil)

	logLvl     string
	syslogAddr string
	hbInterval int
	v4         string
	port       int
	sock       string
	db         string
	addr       string
	pkiCrtPath string
	pkiKeyPath string
	pkiCAPath  string
	statsdip   string
	statsdport int
	trustedips string
)

func main() {
	os.Exit(mainWithExitCode())
}
func mainWithExitCode() int {

	flag.StringVar(&logLvl, "log", "info", "Log level.\nSupported values: debug, info,"+
		" notice, warning, error, critical, alert, emergency")
	flag.StringVar(&syslogAddr, "syslog", "", "Syslog address")
	flag.StringVar(&v4, "4", "0.0.0.0", "IPv4 listener and outbound source address")
	flag.IntVar(&port, "port", 5053, "listener UDP port")
	flag.StringVar(&sock, "sock", "/run/edgedns.sock", "API socket path used by default. "+
		"This parameter is not used if 'address' is defined.")
	flag.StringVar(&addr, "address", "", "API IP address. If defined, socket parameter is not used.")
	flag.StringVar(&db, "db", "/var/lib/edgedns/rrsets.db", "Database file path")
	flag.IntVar(&hbInterval, "hb", 60, "Heartbeat interval in s")
	flag.StringVar(&pkiCrtPath, "cert", "certs/cert.pem", "PKI Cert Path")
	flag.StringVar(&pkiKeyPath, "key", "certs/key.pem", "PKI Key Path")
	flag.StringVar(&pkiCAPath, "ca", "certs/root.pem", "PKI CA Path")
	flag.StringVar(&statsdip, "statsdip", "", "IP address of external statsd server")
	flag.IntVar(&statsdport, "statsdport", 0, "Port of external statsd server")
	flag.StringVar(&trustedips, "trustedips", "", "Trusted ip range for client")
	flag.Parse()

	lvl, err := logger.ParseLevel(logLvl)
	if err != nil {
		log.Errf("Failed to parse log level: %s", err.Error())
		return 1
	}
	logger.SetLevel(lvl)

	err = logger.ConnectSyslog(syslogAddr)
	if err != nil {
		if syslogAddr != "" {
			log.Errf("Syslog(%s) connection failed: %s", syslogAddr, err.Error())
			return 1
		}
		log.Warningf("Fail to connect to local syslog")

	}

	sockPath := path.Dir(sock)
	if _, err = os.Stat(sockPath); os.IsNotExist(err) {
		err = os.MkdirAll(sockPath, 0750)
		if err != nil {
			log.Err(err)
			return 1
		}
	}

	dbPath := path.Dir(db)
	if _, err = os.Stat(dbPath); os.IsNotExist(err) {
		err = os.MkdirAll(dbPath, 0750)
		if err != nil {
			log.Err(err)
			return 1
		}
	}

	cfg := edgedns.Config{
		Addr4:      v4,
		Port:       port,
		TrustedIPs: strings.Split(trustedips, ","),
	}

	st := statsdcli.StatsdConfig{
		Address: statsdip,
		Port:    statsdport,
	}
	cfg.StatsdCfg = st

	stg := &storage.BoltDB{
		Filename: db,
	}

	pki := &grpc.ControlServerPKI{
		Crt: pkiCrtPath,
		Key: pkiKeyPath,
		Ca:  pkiCAPath,
	}

	ctl := &grpc.ControlServer{
		Sock:    sock,
		Address: addr,
		PKI:     pki,
	}
	ex := &dns.NetExchanger{
		DNS: &d.Client{
			Dialer: &net.Dialer{
				Timeout: 5 * time.Second,
				LocalAddr: &net.UDPAddr{
					IP: net.ParseIP(v4), // Use the listener's IP address as the source for outbound requests
				},
			},
		},
	}
	srv := edgedns.NewResponder(cfg, stg, ctl, ex)

	if err = srv.Start(); err != nil {
		log.Err(err)
		return 1
	}
	defer srv.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go NewStatsReporter(time.Duration(hbInterval) * time.Second).Start(ctx)

	// Receive OS signals and listener errors from Start()
	signal.Notify(srv.Sig, syscall.SIGINT, syscall.SIGTERM)
	sig := <-srv.Sig
	switch sig {
	case syscall.SIGCHLD:
		log.Err("Child listener/service unexpectedly died")
		return 1
	default:
		log.Infof("Signal (%v) received, shutting down", sig)
		return 0
	}
}
