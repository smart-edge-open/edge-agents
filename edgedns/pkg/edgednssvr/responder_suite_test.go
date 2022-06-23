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
package edgednssvr_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	edgedns "github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/grpc"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/mock"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr/storage"
)

var (
	dnsServer       *edgedns.Responder
	idleConnsClosed chan struct{}
	exchanger       *mock.Exchanger
)

func TestDns(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Edge DNS Integration Suite")
}

const eport = 60420

var _ = BeforeSuite(func() {
	var err error
	pn := config.GinkgoConfig.ParallelNode
	port := eport + pn
	addr4 := "127.0.0.1"
	sock := fmt.Sprintf("dns_%d.sock", pn)
	db := fmt.Sprintf("dns_%d.db", pn)

	cfg := edgedns.Config{
		Addr4: addr4,
		Port:  port,
	}

	stg := &storage.BoltDB{
		Filename: db,
	}

	ctl := &grpc.ControlServer{
		Sock: sock,
	}

	exchanger = &mock.Exchanger{
		ForwarderCh: make(chan string, 2),
	}

	dnsServer = edgedns.NewResponder(cfg, stg, ctl, exchanger)

	idleConnsClosed = make(chan struct{})

	go func() {
		err = dnsServer.Start()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		signal.Notify(dnsServer.Sig, syscall.SIGINT, syscall.SIGTERM)

		// Receive OS signals and listener errors from Start()
		s := <-dnsServer.Sig
		switch s {
		case syscall.SIGCHLD:
			fmt.Println("Child listener/service unexpectedly died")
		default:
			fmt.Printf("Signal (%v) received\n", s)
		}
		dnsServer.Stop()
		_ = os.Remove(stg.Filename)
		close(idleConnsClosed)
	}()

	// Wait for listeners
	time.Sleep(1 * time.Second)
})

var _ = AfterSuite(func() {
	// Signal Shutdown
	select {
	case dnsServer.Sig <- syscall.SIGINT:
		fmt.Println("Stopping test server")
	default:
		fmt.Println("Shutdown receiver already executed.")
	}

	// Wait for Shutdown to complete
	<-idleConnsClosed
})
