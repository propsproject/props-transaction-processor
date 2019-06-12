// Code generated by protoc-gen-go. DO NOT EDIT.
// source: users.proto

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

type ApplicationUser struct {
	ApplicationId        string   `protobuf:"bytes,1,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	UserId               string   `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Signature            string   `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	Timestamp            int64    `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApplicationUser) Reset()         { *m = ApplicationUser{} }
func (m *ApplicationUser) String() string { return proto.CompactTextString(m) }
func (*ApplicationUser) ProtoMessage()    {}
func (*ApplicationUser) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{0}
}

func (m *ApplicationUser) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApplicationUser.Unmarshal(m, b)
}
func (m *ApplicationUser) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApplicationUser.Marshal(b, m, deterministic)
}
func (m *ApplicationUser) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplicationUser.Merge(m, src)
}
func (m *ApplicationUser) XXX_Size() int {
	return xxx_messageInfo_ApplicationUser.Size(m)
}
func (m *ApplicationUser) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplicationUser.DiscardUnknown(m)
}

var xxx_messageInfo_ApplicationUser proto.InternalMessageInfo

func (m *ApplicationUser) GetApplicationId() string {
	if m != nil {
		return m.ApplicationId
	}
	return ""
}

func (m *ApplicationUser) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *ApplicationUser) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func (m *ApplicationUser) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type WalletToUser struct {
	Address              string             `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Users                []*ApplicationUser `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *WalletToUser) Reset()         { *m = WalletToUser{} }
func (m *WalletToUser) String() string { return proto.CompactTextString(m) }
func (*WalletToUser) ProtoMessage()    {}
func (*WalletToUser) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{1}
}

func (m *WalletToUser) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WalletToUser.Unmarshal(m, b)
}
func (m *WalletToUser) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WalletToUser.Marshal(b, m, deterministic)
}
func (m *WalletToUser) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WalletToUser.Merge(m, src)
}
func (m *WalletToUser) XXX_Size() int {
	return xxx_messageInfo_WalletToUser.Size(m)
}
func (m *WalletToUser) XXX_DiscardUnknown() {
	xxx_messageInfo_WalletToUser.DiscardUnknown(m)
}

var xxx_messageInfo_WalletToUser proto.InternalMessageInfo

func (m *WalletToUser) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *WalletToUser) GetUsers() []*ApplicationUser {
	if m != nil {
		return m.Users
	}
	return nil
}

func init() {
	proto.RegisterType((*ApplicationUser)(nil), "pending_props_pb.ApplicationUser")
	proto.RegisterType((*WalletToUser)(nil), "pending_props_pb.WalletToUser")
}

func init() { proto.RegisterFile("users.proto", fileDescriptor_030765f334c86cea) }

var fileDescriptor_030765f334c86cea = []byte{
	// 206 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x3d, 0x4b, 0xc6, 0x30,
	0x14, 0x85, 0xc9, 0x5b, 0x6d, 0xe9, 0xad, 0x5f, 0x64, 0x31, 0x83, 0x43, 0x2d, 0x08, 0x9d, 0x3a,
	0xe8, 0xe0, 0xec, 0xd8, 0xb5, 0x28, 0x8e, 0x25, 0x35, 0xa1, 0x04, 0xda, 0xe4, 0x92, 0x9b, 0xfe,
	0x0c, 0xff, 0xb3, 0x24, 0x2a, 0x7d, 0xe9, 0x78, 0x9e, 0x27, 0x84, 0x73, 0x2e, 0x54, 0x1b, 0x69,
	0x4f, 0x1d, 0x7a, 0x17, 0x1c, 0xbf, 0x43, 0x6d, 0x95, 0xb1, 0xf3, 0x88, 0xde, 0x21, 0x8d, 0x38,
	0x35, 0xdf, 0x0c, 0x6e, 0xdf, 0x10, 0x17, 0xf3, 0x25, 0x83, 0x71, 0xf6, 0x83, 0xb4, 0xe7, 0x4f,
	0x70, 0x23, 0x77, 0x34, 0x1a, 0x25, 0x58, 0xcd, 0xda, 0x72, 0xb8, 0x3e, 0xa3, 0xbd, 0xe2, 0xf7,
	0x50, 0xc4, 0xbf, 0xa3, 0x3f, 0x25, 0x9f, 0xc7, 0xd8, 0x2b, 0xfe, 0x00, 0x25, 0x99, 0xd9, 0xca,
	0xb0, 0x79, 0x2d, 0xb2, 0xa4, 0x76, 0x10, 0x6d, 0x30, 0xab, 0xa6, 0x20, 0x57, 0x14, 0x17, 0x35,
	0x6b, 0xb3, 0x61, 0x07, 0x8d, 0x84, 0xab, 0x4f, 0xb9, 0x2c, 0x3a, 0xbc, 0xbb, 0xd4, 0x45, 0x40,
	0x21, 0x95, 0xf2, 0x9a, 0xe8, 0xaf, 0xc4, 0x7f, 0xe4, 0xaf, 0x70, 0x99, 0xa6, 0x89, 0x53, 0x9d,
	0xb5, 0xd5, 0xf3, 0x63, 0x77, 0xdc, 0xd6, 0x1d, 0x76, 0x0d, 0xbf, 0xef, 0xa7, 0x3c, 0xdd, 0xe2,
	0xe5, 0x27, 0x00, 0x00, 0xff, 0xff, 0xc2, 0xef, 0x12, 0xd1, 0x1a, 0x01, 0x00, 0x00,
}
