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

package mock

import (
	"github.com/miekg/dns"
	"github.com/stretchr/testify/mock"
)

// StorageMock mocks edgedns.Storage.
type StorageMock struct {
	mock.Mock
}

// Start mock method
func (m *StorageMock) Start() error {
	args := m.Called()
	return args.Error(0)
}

// Stop mock method
func (m *StorageMock) Stop() error {
	args := m.Called()
	return args.Error(0)
}

// SetHostRRSet mock method
func (m *StorageMock) SetHostRRSet(rrtype uint16, fqdn []byte, addrs [][]byte) error {
	args := m.Called(rrtype, fqdn, addrs)
	return args.Error(0)
}

// GetRRSet mock method
func (m *StorageMock) GetRRSet(name string, rrtype uint16) (*[]dns.RR, error) {
	args := m.Called(name, rrtype)
	return args.Get(0).(*[]dns.RR), args.Error(1)
}

// GetAllRRSets mock method
func (m *StorageMock) GetAllRRSets() (map[string][][]byte, error) {
	args := m.Called()
	return args.Get(0).(map[string][][]byte), args.Error(1)
}

// DelRRSet mock method
func (m *StorageMock) DelRRSet(rrtype uint16, fqdn []byte) error {
	args := m.Called(rrtype, fqdn)
	return args.Error(0)
}

// GetForwarders mock method
func (m *StorageMock) GetForwarders() ([][]byte, error) {
	args := m.Called()
	return args.Get(0).([][]byte), args.Error(1)
}

// SetForwarders mock method
func (m *StorageMock) SetForwarders(addrs [][]byte) error {
	args := m.Called(addrs)
	return args.Error(0)
}
