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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.13.0
// source: resolver.proto

package pb

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// DNS Resource Record (https://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml#dns-parameters-4)
type RType int32

const (
	RType_None       RType = 0
	RType_A          RType = 1
	RType_NS         RType = 2
	RType_MD         RType = 3
	RType_MF         RType = 4
	RType_CNAME      RType = 5
	RType_SOA        RType = 6
	RType_MB         RType = 7
	RType_MG         RType = 8
	RType_MR         RType = 9
	RType_NULL       RType = 10
	RType_PTR        RType = 12
	RType_HINFO      RType = 13
	RType_MINFO      RType = 14
	RType_MX         RType = 15
	RType_TXT        RType = 16
	RType_RP         RType = 17
	RType_AFSDB      RType = 18
	RType_X25        RType = 19
	RType_ISDN       RType = 20
	RType_RT         RType = 21
	RType_NSAPPTR    RType = 23
	RType_SIG        RType = 24
	RType_KEY        RType = 25
	RType_PX         RType = 26
	RType_GPOS       RType = 27
	RType_AAAA       RType = 28
	RType_LOC        RType = 29
	RType_NXT        RType = 30
	RType_EID        RType = 31
	RType_NIMLOC     RType = 32
	RType_SRV        RType = 33
	RType_ATMA       RType = 34
	RType_NAPTR      RType = 35
	RType_KX         RType = 36
	RType_CERT       RType = 37
	RType_DNAME      RType = 39
	RType_OPT        RType = 41 // EDNS
	RType_DS         RType = 43
	RType_SSHFP      RType = 44
	RType_RRSIG      RType = 46
	RType_NSEC       RType = 47
	RType_DNSKEY     RType = 48
	RType_DHCID      RType = 49
	RType_NSEC3      RType = 50
	RType_NSEC3PARAM RType = 51
	RType_TLSA       RType = 52
	RType_SMIMEA     RType = 53
	RType_HIP        RType = 55
	RType_NINFO      RType = 56
	RType_RKEY       RType = 57
	RType_TALINK     RType = 58
	RType_CDS        RType = 59
	RType_CDNSKEY    RType = 60
	RType_OPENPGPKEY RType = 61
	RType_SPF        RType = 99
	RType_UINFO      RType = 100
	RType_UID        RType = 101
	RType_GID        RType = 102
	RType_UNSPEC     RType = 103
	RType_NID        RType = 104
	RType_L32        RType = 105
	RType_L64        RType = 106
	RType_LP         RType = 107
	RType_EUI48      RType = 108
	RType_EUI64      RType = 109
	RType_URI        RType = 256
	RType_CAA        RType = 257
	RType_AVC        RType = 258
	RType_TKEY       RType = 249
	RType_TSIG       RType = 250
	// valid Question.Q only
	RType_IXFR     RType = 251
	RType_AXFR     RType = 252
	RType_MAILB    RType = 253
	RType_MAILA    RType = 254
	RType_ANY      RType = 255
	RType_TA       RType = 32768
	RType_DLV      RType = 32769
	RType_Reserved RType = 65535
)

// Enum value maps for RType.
var (
	RType_name = map[int32]string{
		0:     "None",
		1:     "A",
		2:     "NS",
		3:     "MD",
		4:     "MF",
		5:     "CNAME",
		6:     "SOA",
		7:     "MB",
		8:     "MG",
		9:     "MR",
		10:    "NULL",
		12:    "PTR",
		13:    "HINFO",
		14:    "MINFO",
		15:    "MX",
		16:    "TXT",
		17:    "RP",
		18:    "AFSDB",
		19:    "X25",
		20:    "ISDN",
		21:    "RT",
		23:    "NSAPPTR",
		24:    "SIG",
		25:    "KEY",
		26:    "PX",
		27:    "GPOS",
		28:    "AAAA",
		29:    "LOC",
		30:    "NXT",
		31:    "EID",
		32:    "NIMLOC",
		33:    "SRV",
		34:    "ATMA",
		35:    "NAPTR",
		36:    "KX",
		37:    "CERT",
		39:    "DNAME",
		41:    "OPT",
		43:    "DS",
		44:    "SSHFP",
		46:    "RRSIG",
		47:    "NSEC",
		48:    "DNSKEY",
		49:    "DHCID",
		50:    "NSEC3",
		51:    "NSEC3PARAM",
		52:    "TLSA",
		53:    "SMIMEA",
		55:    "HIP",
		56:    "NINFO",
		57:    "RKEY",
		58:    "TALINK",
		59:    "CDS",
		60:    "CDNSKEY",
		61:    "OPENPGPKEY",
		99:    "SPF",
		100:   "UINFO",
		101:   "UID",
		102:   "GID",
		103:   "UNSPEC",
		104:   "NID",
		105:   "L32",
		106:   "L64",
		107:   "LP",
		108:   "EUI48",
		109:   "EUI64",
		256:   "URI",
		257:   "CAA",
		258:   "AVC",
		249:   "TKEY",
		250:   "TSIG",
		251:   "IXFR",
		252:   "AXFR",
		253:   "MAILB",
		254:   "MAILA",
		255:   "ANY",
		32768: "TA",
		32769: "DLV",
		65535: "Reserved",
	}
	RType_value = map[string]int32{
		"None":       0,
		"A":          1,
		"NS":         2,
		"MD":         3,
		"MF":         4,
		"CNAME":      5,
		"SOA":        6,
		"MB":         7,
		"MG":         8,
		"MR":         9,
		"NULL":       10,
		"PTR":        12,
		"HINFO":      13,
		"MINFO":      14,
		"MX":         15,
		"TXT":        16,
		"RP":         17,
		"AFSDB":      18,
		"X25":        19,
		"ISDN":       20,
		"RT":         21,
		"NSAPPTR":    23,
		"SIG":        24,
		"KEY":        25,
		"PX":         26,
		"GPOS":       27,
		"AAAA":       28,
		"LOC":        29,
		"NXT":        30,
		"EID":        31,
		"NIMLOC":     32,
		"SRV":        33,
		"ATMA":       34,
		"NAPTR":      35,
		"KX":         36,
		"CERT":       37,
		"DNAME":      39,
		"OPT":        41,
		"DS":         43,
		"SSHFP":      44,
		"RRSIG":      46,
		"NSEC":       47,
		"DNSKEY":     48,
		"DHCID":      49,
		"NSEC3":      50,
		"NSEC3PARAM": 51,
		"TLSA":       52,
		"SMIMEA":     53,
		"HIP":        55,
		"NINFO":      56,
		"RKEY":       57,
		"TALINK":     58,
		"CDS":        59,
		"CDNSKEY":    60,
		"OPENPGPKEY": 61,
		"SPF":        99,
		"UINFO":      100,
		"UID":        101,
		"GID":        102,
		"UNSPEC":     103,
		"NID":        104,
		"L32":        105,
		"L64":        106,
		"LP":         107,
		"EUI48":      108,
		"EUI64":      109,
		"URI":        256,
		"CAA":        257,
		"AVC":        258,
		"TKEY":       249,
		"TSIG":       250,
		"IXFR":       251,
		"AXFR":       252,
		"MAILB":      253,
		"MAILA":      254,
		"ANY":        255,
		"TA":         32768,
		"DLV":        32769,
		"Reserved":   65535,
	}
)

func (x RType) Enum() *RType {
	p := new(RType)
	*p = x
	return p
}

func (x RType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RType) Descriptor() protoreflect.EnumDescriptor {
	return file_resolver_proto_enumTypes[0].Descriptor()
}

func (RType) Type() protoreflect.EnumType {
	return &file_resolver_proto_enumTypes[0]
}

func (x RType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RType.Descriptor instead.
func (RType) EnumDescriptor() ([]byte, []int) {
	return file_resolver_proto_rawDescGZIP(), []int{0}
}

type HostRecordSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RecordType RType    `protobuf:"varint,1,opt,name=record_type,json=recordType,proto3,enum=pb.RType" json:"record_type,omitempty"`
	Fqdn       string   `protobuf:"bytes,2,opt,name=fqdn,proto3" json:"fqdn,omitempty"`
	Addresses  [][]byte `protobuf:"bytes,3,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *HostRecordSet) Reset() {
	*x = HostRecordSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resolver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HostRecordSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HostRecordSet) ProtoMessage() {}

func (x *HostRecordSet) ProtoReflect() protoreflect.Message {
	mi := &file_resolver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HostRecordSet.ProtoReflect.Descriptor instead.
func (*HostRecordSet) Descriptor() ([]byte, []int) {
	return file_resolver_proto_rawDescGZIP(), []int{0}
}

func (x *HostRecordSet) GetRecordType() RType {
	if x != nil {
		return x.RecordType
	}
	return RType_None
}

func (x *HostRecordSet) GetFqdn() string {
	if x != nil {
		return x.Fqdn
	}
	return ""
}

func (x *HostRecordSet) GetAddresses() [][]byte {
	if x != nil {
		return x.Addresses
	}
	return nil
}

type HostRecordSets struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RecordSets []*HostRecordSet `protobuf:"bytes,1,rep,name=record_sets,json=recordSets,proto3" json:"record_sets,omitempty"`
}

func (x *HostRecordSets) Reset() {
	*x = HostRecordSets{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resolver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HostRecordSets) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HostRecordSets) ProtoMessage() {}

func (x *HostRecordSets) ProtoReflect() protoreflect.Message {
	mi := &file_resolver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HostRecordSets.ProtoReflect.Descriptor instead.
func (*HostRecordSets) Descriptor() ([]byte, []int) {
	return file_resolver_proto_rawDescGZIP(), []int{1}
}

func (x *HostRecordSets) GetRecordSets() []*HostRecordSet {
	if x != nil {
		return x.RecordSets
	}
	return nil
}

// RecordSet represents all values associated with an FQDN and type
//
// Example: An A record for foo.example.org may have one or more addresses,
// and the corresponding RecordSet includes all addresses.
type RecordSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RecordType RType  `protobuf:"varint,1,opt,name=record_type,json=recordType,proto3,enum=pb.RType" json:"record_type,omitempty"`
	Fqdn       string `protobuf:"bytes,2,opt,name=fqdn,proto3" json:"fqdn,omitempty"`
}

func (x *RecordSet) Reset() {
	*x = RecordSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resolver_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordSet) ProtoMessage() {}

func (x *RecordSet) ProtoReflect() protoreflect.Message {
	mi := &file_resolver_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordSet.ProtoReflect.Descriptor instead.
func (*RecordSet) Descriptor() ([]byte, []int) {
	return file_resolver_proto_rawDescGZIP(), []int{2}
}

func (x *RecordSet) GetRecordType() RType {
	if x != nil {
		return x.RecordType
	}
	return RType_None
}

func (x *RecordSet) GetFqdn() string {
	if x != nil {
		return x.Fqdn
	}
	return ""
}

// Forwarders is a set of DNS forwarder addresses the Control service will use to resolve non-authoritative queries.
type Forwarders struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addresses [][]byte `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *Forwarders) Reset() {
	*x = Forwarders{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resolver_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Forwarders) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Forwarders) ProtoMessage() {}

func (x *Forwarders) ProtoReflect() protoreflect.Message {
	mi := &file_resolver_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Forwarders.ProtoReflect.Descriptor instead.
func (*Forwarders) Descriptor() ([]byte, []int) {
	return file_resolver_proto_rawDescGZIP(), []int{3}
}

func (x *Forwarders) GetAddresses() [][]byte {
	if x != nil {
		return x.Addresses
	}
	return nil
}

var File_resolver_proto protoreflect.FileDescriptor

var file_resolver_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x6d, 0x0a, 0x0d, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53,
	0x65, 0x74, 0x12, 0x2a, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x66, 0x71, 0x64, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x71,
	0x64, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73,
	0x22, 0x44, 0x0a, 0x0e, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x65,
	0x74, 0x73, 0x12, 0x32, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x65, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x48, 0x6f, 0x73,
	0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x74, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x53, 0x65, 0x74, 0x73, 0x22, 0x4b, 0x0a, 0x09, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x53, 0x65, 0x74, 0x12, 0x2a, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x71, 0x64, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66,
	0x71, 0x64, 0x6e, 0x22, 0x2a, 0x0a, 0x0a, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72,
	0x73, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0c, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x2a,
	0xa6, 0x06, 0x0a, 0x05, 0x52, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x6f, 0x6e,
	0x65, 0x10, 0x00, 0x12, 0x05, 0x0a, 0x01, 0x41, 0x10, 0x01, 0x12, 0x06, 0x0a, 0x02, 0x4e, 0x53,
	0x10, 0x02, 0x12, 0x06, 0x0a, 0x02, 0x4d, 0x44, 0x10, 0x03, 0x12, 0x06, 0x0a, 0x02, 0x4d, 0x46,
	0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x4e, 0x41, 0x4d, 0x45, 0x10, 0x05, 0x12, 0x07, 0x0a,
	0x03, 0x53, 0x4f, 0x41, 0x10, 0x06, 0x12, 0x06, 0x0a, 0x02, 0x4d, 0x42, 0x10, 0x07, 0x12, 0x06,
	0x0a, 0x02, 0x4d, 0x47, 0x10, 0x08, 0x12, 0x06, 0x0a, 0x02, 0x4d, 0x52, 0x10, 0x09, 0x12, 0x08,
	0x0a, 0x04, 0x4e, 0x55, 0x4c, 0x4c, 0x10, 0x0a, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x54, 0x52, 0x10,
	0x0c, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x0d, 0x12, 0x09, 0x0a, 0x05,
	0x4d, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x0e, 0x12, 0x06, 0x0a, 0x02, 0x4d, 0x58, 0x10, 0x0f, 0x12,
	0x07, 0x0a, 0x03, 0x54, 0x58, 0x54, 0x10, 0x10, 0x12, 0x06, 0x0a, 0x02, 0x52, 0x50, 0x10, 0x11,
	0x12, 0x09, 0x0a, 0x05, 0x41, 0x46, 0x53, 0x44, 0x42, 0x10, 0x12, 0x12, 0x07, 0x0a, 0x03, 0x58,
	0x32, 0x35, 0x10, 0x13, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x53, 0x44, 0x4e, 0x10, 0x14, 0x12, 0x06,
	0x0a, 0x02, 0x52, 0x54, 0x10, 0x15, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x53, 0x41, 0x50, 0x50, 0x54,
	0x52, 0x10, 0x17, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x49, 0x47, 0x10, 0x18, 0x12, 0x07, 0x0a, 0x03,
	0x4b, 0x45, 0x59, 0x10, 0x19, 0x12, 0x06, 0x0a, 0x02, 0x50, 0x58, 0x10, 0x1a, 0x12, 0x08, 0x0a,
	0x04, 0x47, 0x50, 0x4f, 0x53, 0x10, 0x1b, 0x12, 0x08, 0x0a, 0x04, 0x41, 0x41, 0x41, 0x41, 0x10,
	0x1c, 0x12, 0x07, 0x0a, 0x03, 0x4c, 0x4f, 0x43, 0x10, 0x1d, 0x12, 0x07, 0x0a, 0x03, 0x4e, 0x58,
	0x54, 0x10, 0x1e, 0x12, 0x07, 0x0a, 0x03, 0x45, 0x49, 0x44, 0x10, 0x1f, 0x12, 0x0a, 0x0a, 0x06,
	0x4e, 0x49, 0x4d, 0x4c, 0x4f, 0x43, 0x10, 0x20, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x52, 0x56, 0x10,
	0x21, 0x12, 0x08, 0x0a, 0x04, 0x41, 0x54, 0x4d, 0x41, 0x10, 0x22, 0x12, 0x09, 0x0a, 0x05, 0x4e,
	0x41, 0x50, 0x54, 0x52, 0x10, 0x23, 0x12, 0x06, 0x0a, 0x02, 0x4b, 0x58, 0x10, 0x24, 0x12, 0x08,
	0x0a, 0x04, 0x43, 0x45, 0x52, 0x54, 0x10, 0x25, 0x12, 0x09, 0x0a, 0x05, 0x44, 0x4e, 0x41, 0x4d,
	0x45, 0x10, 0x27, 0x12, 0x07, 0x0a, 0x03, 0x4f, 0x50, 0x54, 0x10, 0x29, 0x12, 0x06, 0x0a, 0x02,
	0x44, 0x53, 0x10, 0x2b, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x53, 0x48, 0x46, 0x50, 0x10, 0x2c, 0x12,
	0x09, 0x0a, 0x05, 0x52, 0x52, 0x53, 0x49, 0x47, 0x10, 0x2e, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x53,
	0x45, 0x43, 0x10, 0x2f, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x4e, 0x53, 0x4b, 0x45, 0x59, 0x10, 0x30,
	0x12, 0x09, 0x0a, 0x05, 0x44, 0x48, 0x43, 0x49, 0x44, 0x10, 0x31, 0x12, 0x09, 0x0a, 0x05, 0x4e,
	0x53, 0x45, 0x43, 0x33, 0x10, 0x32, 0x12, 0x0e, 0x0a, 0x0a, 0x4e, 0x53, 0x45, 0x43, 0x33, 0x50,
	0x41, 0x52, 0x41, 0x4d, 0x10, 0x33, 0x12, 0x08, 0x0a, 0x04, 0x54, 0x4c, 0x53, 0x41, 0x10, 0x34,
	0x12, 0x0a, 0x0a, 0x06, 0x53, 0x4d, 0x49, 0x4d, 0x45, 0x41, 0x10, 0x35, 0x12, 0x07, 0x0a, 0x03,
	0x48, 0x49, 0x50, 0x10, 0x37, 0x12, 0x09, 0x0a, 0x05, 0x4e, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x38,
	0x12, 0x08, 0x0a, 0x04, 0x52, 0x4b, 0x45, 0x59, 0x10, 0x39, 0x12, 0x0a, 0x0a, 0x06, 0x54, 0x41,
	0x4c, 0x49, 0x4e, 0x4b, 0x10, 0x3a, 0x12, 0x07, 0x0a, 0x03, 0x43, 0x44, 0x53, 0x10, 0x3b, 0x12,
	0x0b, 0x0a, 0x07, 0x43, 0x44, 0x4e, 0x53, 0x4b, 0x45, 0x59, 0x10, 0x3c, 0x12, 0x0e, 0x0a, 0x0a,
	0x4f, 0x50, 0x45, 0x4e, 0x50, 0x47, 0x50, 0x4b, 0x45, 0x59, 0x10, 0x3d, 0x12, 0x07, 0x0a, 0x03,
	0x53, 0x50, 0x46, 0x10, 0x63, 0x12, 0x09, 0x0a, 0x05, 0x55, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x64,
	0x12, 0x07, 0x0a, 0x03, 0x55, 0x49, 0x44, 0x10, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x47, 0x49, 0x44,
	0x10, 0x66, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x10, 0x67, 0x12, 0x07,
	0x0a, 0x03, 0x4e, 0x49, 0x44, 0x10, 0x68, 0x12, 0x07, 0x0a, 0x03, 0x4c, 0x33, 0x32, 0x10, 0x69,
	0x12, 0x07, 0x0a, 0x03, 0x4c, 0x36, 0x34, 0x10, 0x6a, 0x12, 0x06, 0x0a, 0x02, 0x4c, 0x50, 0x10,
	0x6b, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x55, 0x49, 0x34, 0x38, 0x10, 0x6c, 0x12, 0x09, 0x0a, 0x05,
	0x45, 0x55, 0x49, 0x36, 0x34, 0x10, 0x6d, 0x12, 0x08, 0x0a, 0x03, 0x55, 0x52, 0x49, 0x10, 0x80,
	0x02, 0x12, 0x08, 0x0a, 0x03, 0x43, 0x41, 0x41, 0x10, 0x81, 0x02, 0x12, 0x08, 0x0a, 0x03, 0x41,
	0x56, 0x43, 0x10, 0x82, 0x02, 0x12, 0x09, 0x0a, 0x04, 0x54, 0x4b, 0x45, 0x59, 0x10, 0xf9, 0x01,
	0x12, 0x09, 0x0a, 0x04, 0x54, 0x53, 0x49, 0x47, 0x10, 0xfa, 0x01, 0x12, 0x09, 0x0a, 0x04, 0x49,
	0x58, 0x46, 0x52, 0x10, 0xfb, 0x01, 0x12, 0x09, 0x0a, 0x04, 0x41, 0x58, 0x46, 0x52, 0x10, 0xfc,
	0x01, 0x12, 0x0a, 0x0a, 0x05, 0x4d, 0x41, 0x49, 0x4c, 0x42, 0x10, 0xfd, 0x01, 0x12, 0x0a, 0x0a,
	0x05, 0x4d, 0x41, 0x49, 0x4c, 0x41, 0x10, 0xfe, 0x01, 0x12, 0x08, 0x0a, 0x03, 0x41, 0x4e, 0x59,
	0x10, 0xff, 0x01, 0x12, 0x08, 0x0a, 0x02, 0x54, 0x41, 0x10, 0x80, 0x80, 0x02, 0x12, 0x09, 0x0a,
	0x03, 0x44, 0x4c, 0x56, 0x10, 0x81, 0x80, 0x02, 0x12, 0x0e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x64, 0x10, 0xff, 0xff, 0x03, 0x32, 0xc1, 0x02, 0x0a, 0x07, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x12, 0x43, 0x0a, 0x14, 0x53, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x74, 0x61, 0x74, 0x69, 0x76, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x11, 0x2e, 0x70,
	0x62, 0x2e, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x41, 0x6c, 0x6c, 0x48, 0x6f, 0x73, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x53, 0x65, 0x74, 0x73, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x61, 0x74, 0x69, 0x76, 0x65, 0x12, 0x0d, 0x2e,
	0x70, 0x62, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x72,
	0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65, 0x72, 0x73, 0x22,
	0x00, 0x12, 0x39, 0x0a, 0x0d, 0x53, 0x65, 0x74, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65,
	0x72, 0x73, 0x12, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x65,
	0x72, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resolver_proto_rawDescOnce sync.Once
	file_resolver_proto_rawDescData = file_resolver_proto_rawDesc
)

func file_resolver_proto_rawDescGZIP() []byte {
	file_resolver_proto_rawDescOnce.Do(func() {
		file_resolver_proto_rawDescData = protoimpl.X.CompressGZIP(file_resolver_proto_rawDescData)
	})
	return file_resolver_proto_rawDescData
}

var file_resolver_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resolver_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_resolver_proto_goTypes = []interface{}{
	(RType)(0),             // 0: pb.RType
	(*HostRecordSet)(nil),  // 1: pb.HostRecordSet
	(*HostRecordSets)(nil), // 2: pb.HostRecordSets
	(*RecordSet)(nil),      // 3: pb.RecordSet
	(*Forwarders)(nil),     // 4: pb.Forwarders
	(*empty.Empty)(nil),    // 5: google.protobuf.Empty
}
var file_resolver_proto_depIdxs = []int32{
	0, // 0: pb.HostRecordSet.record_type:type_name -> pb.RType
	1, // 1: pb.HostRecordSets.record_sets:type_name -> pb.HostRecordSet
	0, // 2: pb.RecordSet.record_type:type_name -> pb.RType
	1, // 3: pb.Control.SetAuthoritativeHost:input_type -> pb.HostRecordSet
	5, // 4: pb.Control.GetAllHosts:input_type -> google.protobuf.Empty
	3, // 5: pb.Control.DeleteAuthoritative:input_type -> pb.RecordSet
	5, // 6: pb.Control.GetForwarders:input_type -> google.protobuf.Empty
	4, // 7: pb.Control.SetForwarders:input_type -> pb.Forwarders
	5, // 8: pb.Control.SetAuthoritativeHost:output_type -> google.protobuf.Empty
	2, // 9: pb.Control.GetAllHosts:output_type -> pb.HostRecordSets
	5, // 10: pb.Control.DeleteAuthoritative:output_type -> google.protobuf.Empty
	4, // 11: pb.Control.GetForwarders:output_type -> pb.Forwarders
	5, // 12: pb.Control.SetForwarders:output_type -> google.protobuf.Empty
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_resolver_proto_init() }
func file_resolver_proto_init() {
	if File_resolver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resolver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HostRecordSet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_resolver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HostRecordSets); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_resolver_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordSet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_resolver_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Forwarders); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resolver_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_resolver_proto_goTypes,
		DependencyIndexes: file_resolver_proto_depIdxs,
		EnumInfos:         file_resolver_proto_enumTypes,
		MessageInfos:      file_resolver_proto_msgTypes,
	}.Build()
	File_resolver_proto = out.File
	file_resolver_proto_rawDesc = nil
	file_resolver_proto_goTypes = nil
	file_resolver_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ControlClient is the client API for Control service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ControlClient interface {
	SetAuthoritativeHost(ctx context.Context, in *HostRecordSet, opts ...grpc.CallOption) (*empty.Empty, error)
	GetAllHosts(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*HostRecordSets, error)
	DeleteAuthoritative(ctx context.Context, in *RecordSet, opts ...grpc.CallOption) (*empty.Empty, error)
	GetForwarders(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Forwarders, error)
	SetForwarders(ctx context.Context, in *Forwarders, opts ...grpc.CallOption) (*empty.Empty, error)
}

type controlClient struct {
	cc grpc.ClientConnInterface
}

func NewControlClient(cc grpc.ClientConnInterface) ControlClient {
	return &controlClient{cc}
}

func (c *controlClient) SetAuthoritativeHost(ctx context.Context, in *HostRecordSet, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/pb.Control/SetAuthoritativeHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) GetAllHosts(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*HostRecordSets, error) {
	out := new(HostRecordSets)
	err := c.cc.Invoke(ctx, "/pb.Control/GetAllHosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) DeleteAuthoritative(ctx context.Context, in *RecordSet, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/pb.Control/DeleteAuthoritative", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) GetForwarders(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Forwarders, error) {
	out := new(Forwarders)
	err := c.cc.Invoke(ctx, "/pb.Control/GetForwarders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) SetForwarders(ctx context.Context, in *Forwarders, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/pb.Control/SetForwarders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ControlServer is the server API for Control service.
type ControlServer interface {
	SetAuthoritativeHost(context.Context, *HostRecordSet) (*empty.Empty, error)
	GetAllHosts(context.Context, *empty.Empty) (*HostRecordSets, error)
	DeleteAuthoritative(context.Context, *RecordSet) (*empty.Empty, error)
	GetForwarders(context.Context, *empty.Empty) (*Forwarders, error)
	SetForwarders(context.Context, *Forwarders) (*empty.Empty, error)
}

// UnimplementedControlServer can be embedded to have forward compatible implementations.
type UnimplementedControlServer struct {
}

func (*UnimplementedControlServer) SetAuthoritativeHost(context.Context, *HostRecordSet) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAuthoritativeHost not implemented")
}
func (*UnimplementedControlServer) GetAllHosts(context.Context, *empty.Empty) (*HostRecordSets, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllHosts not implemented")
}
func (*UnimplementedControlServer) DeleteAuthoritative(context.Context, *RecordSet) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAuthoritative not implemented")
}
func (*UnimplementedControlServer) GetForwarders(context.Context, *empty.Empty) (*Forwarders, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetForwarders not implemented")
}
func (*UnimplementedControlServer) SetForwarders(context.Context, *Forwarders) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetForwarders not implemented")
}

func RegisterControlServer(s *grpc.Server, srv ControlServer) {
	s.RegisterService(&_Control_serviceDesc, srv)
}

func _Control_SetAuthoritativeHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostRecordSet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).SetAuthoritativeHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Control/SetAuthoritativeHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).SetAuthoritativeHost(ctx, req.(*HostRecordSet))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_GetAllHosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).GetAllHosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Control/GetAllHosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).GetAllHosts(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_DeleteAuthoritative_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordSet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).DeleteAuthoritative(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Control/DeleteAuthoritative",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).DeleteAuthoritative(ctx, req.(*RecordSet))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_GetForwarders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).GetForwarders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Control/GetForwarders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).GetForwarders(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_SetForwarders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Forwarders)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).SetForwarders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Control/SetForwarders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).SetForwarders(ctx, req.(*Forwarders))
	}
	return interceptor(ctx, in, info, handler)
}

var _Control_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Control",
	HandlerType: (*ControlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetAuthoritativeHost",
			Handler:    _Control_SetAuthoritativeHost_Handler,
		},
		{
			MethodName: "GetAllHosts",
			Handler:    _Control_GetAllHosts_Handler,
		},
		{
			MethodName: "DeleteAuthoritative",
			Handler:    _Control_DeleteAuthoritative_Handler,
		},
		{
			MethodName: "GetForwarders",
			Handler:    _Control_GetForwarders_Handler,
		},
		{
			MethodName: "SetForwarders",
			Handler:    _Control_SetForwarders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "resolver.proto",
}
