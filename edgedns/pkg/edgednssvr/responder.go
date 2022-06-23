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

package edgednssvr

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/miekg/dns"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/statsdcli"
	rate "github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/utils"
	logger "github.com/smart-edge-open/edge-services/common/log"
)

var log = logger.DefaultLogger.WithField("edgedns", nil)

// Storage is a backend persistence for all records
type Storage interface {
	Start() error
	Stop() error

	// SetHostRRSet Creates or updates all resource records for a given FQDN
	// 				and resource record type
	//
	// rrtype 		Resource Record Type (A or AAAA)
	// fqdn			Fully Qualified Domain Name
	// addrs		One or more IP addresses for the FQDN
	SetHostRRSet(rrtype uint16, fqdn []byte, addrs [][]byte) error

	// GetRRSet returns all resources records for an FQDN and resource type
	GetRRSet(name string, rrtype uint16) (*[]dns.RR, error)

	// GetAllRRSets returns all resource records for all FQDNs. RR type is assumed A.
	GetAllRRSets() (map[string][][]byte, error)

	// DelRRSet removes a RR set for a given FQDN and resource type
	DelRRSet(rrtype uint16, fqdn []byte) error

	// GetForwarders gets the responder's forwarder configuration.
	GetForwarders() ([][]byte, error)

	// SetForwarders set the responder's forwarder configuration.
	SetForwarders(addrs [][]byte) error
}

// ControlServer provides an API to administer the runtime state
// of the Responder records
type ControlServer interface {
	Start(stg Storage, r *Responder) error
	GracefulStop() error
}

// Exchanger resolves DNS queries.
type Exchanger interface {
	Exchange(ctx context.Context, q *dns.Msg, addr string) (*dns.Msg, error)
}

// Config contains all runtime configuration parameters
type Config struct {
	Addr4      string
	Port       int
	forwarders []string
	StatsdCfg  statsdcli.StatsdConfig
	TrustedIPs []string
}

// Responder handles all DNS queries
type Responder struct {
	Sig          chan os.Signal // Shutdown signals
	cfg          Config
	mux          sync.RWMutex // Guards configuration
	server4      *dns.Server
	storage      Storage
	control      ControlServer
	exchanger    Exchanger
	statsdclient statsdcli.StatsdClient
}

// NewResponder returns a new DNS Responder (Server)
func NewResponder(cfg Config, stg Storage, ctl ControlServer, ex Exchanger) *Responder {
	return &Responder{
		Sig:       make(chan os.Signal),
		cfg:       cfg,
		storage:   stg,
		control:   ctl,
		exchanger: ex,
	}
}

// Start will register and start all services
func (r *Responder) Start() error {
	log.Infof("Starting Edge DNS Server")

	// Start DB backend
	if err := r.storage.Start(); err != nil {
		return fmt.Errorf("unable to start DB: %s", err)
	}

	// Load configuration
	if err := r.loadConfig(); err != nil {
		return fmt.Errorf("unable to load configuration: %v", err)
	}

	// Start gRPC API
	if err := r.control.Start(r.storage, r); err != nil {
		return err
	}
	r.statsdclient = statsdcli.NewClient(r.cfg.StatsdCfg)

	// Define a per-request rate limiting function
	allowed, stopLimit := rate.Limit(50.0, 100)

	// HandleFunc uses DefaultMsgAcceptFunc,
	// which checks the request and will reject if:
	//
	// * isn't a request (don't respond in that case).
	// * opcode isn't OpcodeQuery or OpcodeNotify
	// * Zero bit isn't zero
	// * has more than 1 question in the question section
	// * has more than 1 RR in the Answer section
	// * has more than 0 RRs in the Authority section
	// * has more than 2 RRs in the Additional section

	//adding more trusted ips
	addtrustedIPs(r.cfg.TrustedIPs)
	dns.Handle(
		".",
		UnsupportedHandler(
			LimitRequestRate(allowed,
				LoggingHandler(
					r.AuthorityHandler(
						r.IterativeForwardingHandler(
							errorHandler(),
						),
					),
				),
			),
		),
	)

	go func() {
		defer stopLimit()
		// Start DNS Listeners
		r.startListeners()
	}()

	return nil
}

// loadConfig reads and sets the configuration from storage.
func (r *Responder) loadConfig() error {
	fwdrs, err := r.storage.GetForwarders()
	if err != nil {
		return fmt.Errorf("unable to get forwarders from storage: %v", err)
	}

	var addrs []string
	for _, fwdr := range fwdrs {
		ip := net.IP(fwdr)
		if ip.To4() == nil && ip.To16() == nil {
			return fmt.Errorf("unable to parse forwarder IP %s from storage: %s", fwdr, err)
		}
		addrs = append(addrs, ip.String())
	}

	r.mux.Lock()
	r.cfg.forwarders = addrs
	r.mux.Unlock()

	return nil
}

func (r *Responder) startListeners() {
	if len(r.cfg.Addr4) > 0 {
		log.Infof("Starting IPv4 DNS Listener at %s:%d",
			r.cfg.Addr4, r.cfg.Port)
		r.server4 = &dns.Server{Addr: r.cfg.Addr4 + ":" +
			strconv.Itoa(r.cfg.Port), Net: "udp"}

		if err := r.server4.ListenAndServe(); err != nil {
			log.Errf("IPv4 listener error: %s", err)
			r.Sig <- syscall.SIGCHLD
		}

	}

	if len(r.cfg.Addr4) == 0 {
		log.Infoln("Starting DNS Listener on all addresses")
		r.server4 = &dns.Server{Addr: ":" +
			strconv.Itoa(r.cfg.Port), Net: "udp"}
		if err := r.server4.ListenAndServe(); err != nil {
			log.Errf("Any-address listener error: %s", err)
			r.Sig <- syscall.SIGCHLD
		}

	}
}

// Stop all listeners
func (r *Responder) Stop() {
	log.Debugln("Edge DNS Server shutdown started")

	if r.server4 != nil {
		log.Debugln("Stopping IPv4 Responder")
		if err := r.server4.Shutdown(); err != nil {
			log.Errf("IPv4 listener shutdown error: %s", err)
		}
	}

	log.Debugln("Stopping API")
	if err := r.control.GracefulStop(); err != nil {
		log.Errf("Control Server Shutdown error: %s", err)
	}

	log.Debugln("Stopping DB")
	if err := r.storage.Stop(); err != nil {
		log.Errf("DB Shutdown error: %s", err)
	}

	log.Infoln("Edge DNS Server stopped")
}

// SetForwarders sets the forwarders. Providing an empty slice will clear the fowarders.
func (r *Responder) SetForwarders(fwdrs []string) error {
	metricVal := 0
	defer func() {
		r.metric("forwarders", strings.Join(fwdrs, ","), metricVal)
	}()

	var bytes [][]byte
	for _, fwdr := range fwdrs {
		ip := net.ParseIP(fwdr)
		if ip == nil {
			return fmt.Errorf("unable to parse IP %s", fwdr)
		}
		bytes = append(bytes, ip)
	}

	r.mux.Lock()
	r.cfg.forwarders = fwdrs
	r.mux.Unlock()

	if err := r.storage.SetForwarders(bytes); err != nil {
		return err
	}

	metricVal = 1

	return nil
}

func (r *Responder) metric(target, action string, value int) {
	//Connect to statsd server
	if err := r.statsdclient.Connect(); err != nil {
		return
	}
	r.statsdclient.Gauge(target, action, value)
	r.statsdclient.Close()
}

// Forwarders returns the forwarders.
func (r *Responder) Forwarders() []string {
	r.mux.RLock()
	fwdrs := make([]string, len(r.cfg.forwarders))
	_ = copy(fwdrs, r.cfg.forwarders)
	r.mux.RUnlock()

	return fwdrs
}
