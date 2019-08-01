// Code generated by protoc-gen-go. DO NOT EDIT.
// source: reward_entities.proto

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

type RewardEntityType int32

const (
	RewardEntityType_VALIDATOR   RewardEntityType = 0
	RewardEntityType_APPLICATION RewardEntityType = 1
)

var RewardEntityType_name = map[int32]string{
	0: "VALIDATOR",
	1: "APPLICATION",
}

var RewardEntityType_value = map[string]int32{
	"VALIDATOR":   0,
	"APPLICATION": 1,
}

func (x RewardEntityType) String() string {
	return proto.EnumName(RewardEntityType_name, int32(x))
}

func (RewardEntityType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ca39a3517e7f2e26, []int{0}
}

type RewardEntityState int32

const (
	RewardEntityState_ACTIVE   RewardEntityState = 0
	RewardEntityState_INACTIVE RewardEntityState = 1
)

var RewardEntityState_name = map[int32]string{
	0: "ACTIVE",
	1: "INACTIVE",
}

var RewardEntityState_value = map[string]int32{
	"ACTIVE":   0,
	"INACTIVE": 1,
}

func (x RewardEntityState) String() string {
	return proto.EnumName(RewardEntityState_name, int32(x))
}

func (RewardEntityState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ca39a3517e7f2e26, []int{1}
}

type RewardEntity struct {
	Type                 RewardEntityType  `protobuf:"varint,1,opt,name=type,proto3,enum=pending_props_pb.RewardEntityType" json:"type,omitempty"`
	Name                 string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Address              string            `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	RewardsAddress       string            `protobuf:"bytes,4,opt,name=rewardsAddress,proto3" json:"rewardsAddress,omitempty"`
	SidechainAddress     string            `protobuf:"bytes,5,opt,name=sidechainAddress,proto3" json:"sidechainAddress,omitempty"`
	Status               RewardEntityState `protobuf:"varint,6,opt,name=status,proto3,enum=pending_props_pb.RewardEntityState" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RewardEntity) Reset()         { *m = RewardEntity{} }
func (m *RewardEntity) String() string { return proto.CompactTextString(m) }
func (*RewardEntity) ProtoMessage()    {}
func (*RewardEntity) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca39a3517e7f2e26, []int{0}
}

func (m *RewardEntity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RewardEntity.Unmarshal(m, b)
}
func (m *RewardEntity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RewardEntity.Marshal(b, m, deterministic)
}
func (m *RewardEntity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RewardEntity.Merge(m, src)
}
func (m *RewardEntity) XXX_Size() int {
	return xxx_messageInfo_RewardEntity.Size(m)
}
func (m *RewardEntity) XXX_DiscardUnknown() {
	xxx_messageInfo_RewardEntity.DiscardUnknown(m)
}

var xxx_messageInfo_RewardEntity proto.InternalMessageInfo

func (m *RewardEntity) GetType() RewardEntityType {
	if m != nil {
		return m.Type
	}
	return RewardEntityType_VALIDATOR
}

func (m *RewardEntity) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RewardEntity) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *RewardEntity) GetRewardsAddress() string {
	if m != nil {
		return m.RewardsAddress
	}
	return ""
}

func (m *RewardEntity) GetSidechainAddress() string {
	if m != nil {
		return m.SidechainAddress
	}
	return ""
}

func (m *RewardEntity) GetStatus() RewardEntityState {
	if m != nil {
		return m.Status
	}
	return RewardEntityState_ACTIVE
}

func init() {
	proto.RegisterEnum("pending_props_pb.RewardEntityType", RewardEntityType_name, RewardEntityType_value)
	proto.RegisterEnum("pending_props_pb.RewardEntityState", RewardEntityState_name, RewardEntityState_value)
	proto.RegisterType((*RewardEntity)(nil), "pending_props_pb.RewardEntity")
}

func init() { proto.RegisterFile("reward_entities.proto", fileDescriptor_ca39a3517e7f2e26) }

var fileDescriptor_ca39a3517e7f2e26 = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0xd1, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0x07, 0xf0, 0x75, 0xd6, 0xea, 0x9e, 0x73, 0xc6, 0x07, 0x42, 0x8f, 0x63, 0x82, 0x8c, 0x82,
	0x3d, 0x4c, 0xf0, 0xe2, 0x29, 0xcc, 0x1d, 0x0a, 0x63, 0x1b, 0xb5, 0xec, 0x5a, 0x32, 0xfb, 0xd0,
	0x1c, 0x4c, 0x43, 0x12, 0x91, 0x7e, 0x00, 0xbf, 0xb7, 0x18, 0x3b, 0x90, 0x0e, 0x76, 0xcb, 0xfb,
	0xf3, 0x7b, 0xe1, 0x1f, 0x02, 0x37, 0x86, 0xbe, 0x84, 0xa9, 0x4a, 0x52, 0x4e, 0x3a, 0x49, 0x36,
	0xd5, 0xa6, 0x76, 0x35, 0x32, 0x4d, 0xaa, 0x92, 0xea, 0xad, 0xd4, 0xa6, 0xd6, 0xb6, 0xd4, 0xbb,
	0xc9, 0x77, 0x1f, 0x86, 0xb9, 0xb7, 0x8b, 0x5f, 0xda, 0xe0, 0x23, 0x84, 0xae, 0xd1, 0x14, 0x07,
	0xe3, 0x60, 0x3a, 0x9a, 0x4d, 0xd2, 0xee, 0x46, 0xfa, 0x5f, 0x17, 0x8d, 0xa6, 0xdc, 0x7b, 0x44,
	0x08, 0x95, 0xf8, 0xa0, 0xb8, 0x3f, 0x0e, 0xa6, 0x83, 0xdc, 0x9f, 0x31, 0x86, 0x33, 0x51, 0x55,
	0x86, 0xac, 0x8d, 0x4f, 0x7c, 0xbc, 0x1f, 0xf1, 0x0e, 0x46, 0x7f, 0x0d, 0x2d, 0x6f, 0x41, 0xe8,
	0x41, 0x27, 0xc5, 0x04, 0x98, 0x95, 0x15, 0xbd, 0xbe, 0x0b, 0xa9, 0xf6, 0xf2, 0xd4, 0xcb, 0x83,
	0x1c, 0x9f, 0x20, 0xb2, 0x4e, 0xb8, 0x4f, 0x1b, 0x47, 0xbe, 0xfb, 0xed, 0xf1, 0xee, 0x2f, 0x4e,
	0x38, 0xca, 0xdb, 0x95, 0x64, 0x06, 0xac, 0xfb, 0x30, 0xbc, 0x84, 0xc1, 0x96, 0x2f, 0xb3, 0x67,
	0x5e, 0xac, 0x73, 0xd6, 0xc3, 0x2b, 0xb8, 0xe0, 0x9b, 0xcd, 0x32, 0x9b, 0xf3, 0x22, 0x5b, 0xaf,
	0x58, 0x90, 0xdc, 0xc3, 0xf5, 0xc1, 0x85, 0x08, 0x10, 0xf1, 0x79, 0x91, 0x6d, 0x17, 0xac, 0x87,
	0x43, 0x38, 0xcf, 0x56, 0xed, 0x14, 0xec, 0x22, 0xff, 0x07, 0x0f, 0x3f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x24, 0x6a, 0x63, 0x7a, 0x9c, 0x01, 0x00, 0x00,
}