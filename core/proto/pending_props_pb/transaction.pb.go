// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transaction.proto

package pending_props_pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Transaction struct {
	Type                 Method   `protobuf:"varint,1,opt,name=type,proto3,enum=pending_props_pb.Method" json:"type,omitempty"`
	UserId               string   `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ApplicationId        string   `protobuf:"bytes,3,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	Timestamp            int64    `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Amount               string   `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount,omitempty"`
	Description          string   `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	TxHash               string   `protobuf:"bytes,7,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	Wallet               string   `protobuf:"bytes,8,opt,name=wallet,proto3" json:"wallet,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0}
}

func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (m *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(m, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetType() Method {
	if m != nil {
		return m.Type
	}
	return Method_ISSUE
}

func (m *Transaction) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *Transaction) GetApplicationId() string {
	if m != nil {
		return m.ApplicationId
	}
	return ""
}

func (m *Transaction) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Transaction) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func (m *Transaction) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Transaction) GetTxHash() string {
	if m != nil {
		return m.TxHash
	}
	return ""
}

func (m *Transaction) GetWallet() string {
	if m != nil {
		return m.Wallet
	}
	return ""
}

func init() {
	proto.RegisterType((*Transaction)(nil), "pending_props_pb.Transaction")
}

func init() { proto.RegisterFile("transaction.proto", fileDescriptor_2cc4e03d2c28c490) }

var fileDescriptor_2cc4e03d2c28c490 = []byte{
	// 238 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x31, 0x6b, 0xc3, 0x30,
	0x10, 0x85, 0x71, 0x92, 0x3a, 0x8d, 0x42, 0x42, 0xab, 0xa1, 0x3d, 0x4a, 0x07, 0x53, 0x28, 0x78,
	0x28, 0x1e, 0xda, 0x3f, 0xd1, 0x0c, 0x5d, 0x4c, 0x77, 0x73, 0xb1, 0x44, 0x2d, 0xb0, 0xa5, 0x43,
	0xba, 0xd0, 0x64, 0xef, 0x0f, 0x2f, 0x92, 0x86, 0x84, 0x8c, 0xef, 0x7b, 0xf7, 0x3e, 0x84, 0xc4,
	0x3d, 0x7b, 0xb4, 0x01, 0x7b, 0x36, 0xce, 0x36, 0xe4, 0x1d, 0x3b, 0x79, 0x47, 0xda, 0x2a, 0x63,
	0x7f, 0x3a, 0xf2, 0x8e, 0x42, 0x47, 0xfb, 0xa7, 0x0d, 0xe1, 0x69, 0x74, 0xa8, 0xf2, 0xc1, 0xcb,
	0xdf, 0x4c, 0xac, 0xbf, 0xcf, 0x33, 0xf9, 0x26, 0x16, 0x7c, 0x22, 0x0d, 0x45, 0x55, 0xd4, 0xdb,
	0x77, 0x68, 0xae, 0xf7, 0xcd, 0x97, 0xe6, 0xc1, 0xa9, 0x36, 0x5d, 0xc9, 0x47, 0xb1, 0x3c, 0x04,
	0xed, 0x3b, 0xa3, 0x60, 0x56, 0x15, 0xf5, 0xaa, 0x2d, 0x63, 0xdc, 0x29, 0xf9, 0x2a, 0xb6, 0x48,
	0x34, 0x9a, 0x1e, 0xa3, 0x35, 0xf6, 0xf3, 0xd4, 0x6f, 0x2e, 0xe8, 0x4e, 0xc9, 0x67, 0xb1, 0x62,
	0x33, 0xe9, 0xc0, 0x38, 0x11, 0x2c, 0xaa, 0xa2, 0x9e, 0xb7, 0x67, 0x20, 0x1f, 0x44, 0x89, 0x93,
	0x3b, 0x58, 0x86, 0x9b, 0x2c, 0xcf, 0x49, 0x56, 0x62, 0xad, 0x74, 0xe8, 0xbd, 0xa1, 0xa8, 0x81,
	0x32, 0x95, 0x97, 0x28, 0xbe, 0x8b, 0x8f, 0xdd, 0x80, 0x61, 0x80, 0x65, 0x9e, 0xf2, 0xf1, 0x13,
	0xc3, 0x10, 0x95, 0xbf, 0x38, 0x8e, 0x9a, 0xe1, 0x36, 0xf3, 0x9c, 0xf6, 0x65, 0xfa, 0x8d, 0x8f,
	0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0x69, 0x6e, 0x03, 0x59, 0x43, 0x01, 0x00, 0x00,
}
