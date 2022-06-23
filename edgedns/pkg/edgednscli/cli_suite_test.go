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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednscli"
)

var (
	testTmpFolder string
	cliPKI        edgednscli.PKIPaths
	fakeSvr       *ControlServer
)

const serverTestAddress = "localhost:14204"

func TestCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cli Suite")
}

var _ = BeforeSuite(func() {

	var err error

	testTmpFolder, err = ioutil.TempDir("/tmp", "dns_test")
	Expect(err).ShouldNot(HaveOccurred())

	Expect(prepareTestCredentials(testTmpFolder)).ToNot(HaveOccurred())

	cliPKI = edgednscli.PKIPaths{
		CrtPath:            filepath.Join(testTmpFolder, "c_cert.pem"),
		KeyPath:            filepath.Join(testTmpFolder, "c_key.pem"),
		CAPath:             filepath.Join(testTmpFolder, "cacerts.pem"),
		ServerNameOverride: "",
	}

	pki := &ControlServerPKI{
		Crt: filepath.Join(testTmpFolder, "cert.pem"),
		Key: filepath.Join(testTmpFolder, "key.pem"),
		CA:  filepath.Join(testTmpFolder, "cacerts.pem"),
	}

	fakeSvr = &ControlServer{
		Address: serverTestAddress,
		PKI:     pki,
	}

	Expect(fakeSvr.StartServer()).ToNot(HaveOccurred())

	time.Sleep(1 * time.Second)
})

var _ = AfterSuite(func() {
	Expect(fakeSvr.GracefulStop()).ToNot(HaveOccurred())

	err := os.RemoveAll(testTmpFolder)
	Expect(err).ShouldNot(HaveOccurred())
})
