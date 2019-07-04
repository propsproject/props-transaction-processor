// Code generated by protoc-gen-go. DO NOT EDIT.
// source: balance.proto

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

type BalanceType int32

const (
	BalanceType_USER   BalanceType = 0
	BalanceType_WALLET BalanceType = 1
)

var BalanceType_name = map[int32]string{
	0: "USER",
	1: "WALLET",
}

var BalanceType_value = map[string]int32{
	"USER":   0,
	"WALLET": 1,
}

func (x BalanceType) String() string {
	return proto.EnumName(BalanceType_name, int32(x))
}

func (BalanceType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ee25a00b628521b1, []int{0}
}

type UpdateType int32

const (
	UpdateType_PENDING_PROPS_BALANCE UpdateType = 0
	UpdateType_PROPS_BALANCE         UpdateType = 1
)

var UpdateType_name = map[int32]string{
	0: "PENDING_PROPS_BALANCE",
	1: "PROPS_BALANCE",
}

var UpdateType_value = map[string]int32{
	"PENDING_PROPS_BALANCE": 0,
	"PROPS_BALANCE":         1,
}

func (x UpdateType) String() string {
	return proto.EnumName(UpdateType_name, int32(x))
}

func (UpdateType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ee25a00b628521b1, []int{1}
}

type BalanceDetails struct {
	Pending              string     `protobuf:"bytes,1,opt,name=pending,proto3" json:"pending,omitempty"`
	TotalPending         string     `protobuf:"bytes,2,opt,name=total_pending,json=totalPending,proto3" json:"total_pending,omitempty"`
	Transferable         string     `protobuf:"bytes,3,opt,name=transferable,proto3" json:"transferable,omitempty"`
	Bonded               string     `protobuf:"bytes,4,opt,name=bonded,proto3" json:"bonded,omitempty"`
	Delegated            string     `protobuf:"bytes,5,opt,name=delegated,proto3" json:"delegated,omitempty"`
	DelegatedTo          string     `protobuf:"bytes,6,opt,name=delegatedTo,proto3" json:"delegatedTo,omitempty"`
	Timestamp            int64      `protobuf:"varint,7,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	LastEthBlockId       int64      `protobuf:"varint,8,opt,name=last_eth_block_id,json=lastEthBlockId,proto3" json:"last_eth_block_id,omitempty"`
	LastUpdateType       UpdateType `protobuf:"varint,9,opt,name=last_update_type,json=lastUpdateType,proto3,enum=pending_props_pb.UpdateType" json:"last_update_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *BalanceDetails) Reset()         { *m = BalanceDetails{} }
func (m *BalanceDetails) String() string { return proto.CompactTextString(m) }
func (*BalanceDetails) ProtoMessage()    {}
func (*BalanceDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee25a00b628521b1, []int{0}
}

func (m *BalanceDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BalanceDetails.Unmarshal(m, b)
}
func (m *BalanceDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BalanceDetails.Marshal(b, m, deterministic)
}
func (m *BalanceDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BalanceDetails.Merge(m, src)
}
func (m *BalanceDetails) XXX_Size() int {
	return xxx_messageInfo_BalanceDetails.Size(m)
}
func (m *BalanceDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_BalanceDetails.DiscardUnknown(m)
}

var xxx_messageInfo_BalanceDetails proto.InternalMessageInfo

func (m *BalanceDetails) GetPending() string {
	if m != nil {
		return m.Pending
	}
	return ""
}

func (m *BalanceDetails) GetTotalPending() string {
	if m != nil {
		return m.TotalPending
	}
	return ""
}

func (m *BalanceDetails) GetTransferable() string {
	if m != nil {
		return m.Transferable
	}
	return ""
}

func (m *BalanceDetails) GetBonded() string {
	if m != nil {
		return m.Bonded
	}
	return ""
}

func (m *BalanceDetails) GetDelegated() string {
	if m != nil {
		return m.Delegated
	}
	return ""
}

func (m *BalanceDetails) GetDelegatedTo() string {
	if m != nil {
		return m.DelegatedTo
	}
	return ""
}

func (m *BalanceDetails) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *BalanceDetails) GetLastEthBlockId() int64 {
	if m != nil {
		return m.LastEthBlockId
	}
	return 0
}

func (m *BalanceDetails) GetLastUpdateType() UpdateType {
	if m != nil {
		return m.LastUpdateType
	}
	return UpdateType_PENDING_PROPS_BALANCE
}

type Balance struct {
	UserId               string          `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ApplicationId        string          `protobuf:"bytes,2,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	BalanceDetails       *BalanceDetails `protobuf:"bytes,3,opt,name=balance_details,json=balanceDetails,proto3" json:"balance_details,omitempty"`
	PreCutoffDetails     *BalanceDetails `protobuf:"bytes,4,opt,name=pre_cutoff_details,json=preCutoffDetails,proto3" json:"pre_cutoff_details,omitempty"`
	Type                 BalanceType     `protobuf:"varint,5,opt,name=type,proto3,enum=pending_props_pb.BalanceType" json:"type,omitempty"`
	LinkedWallet         string          `protobuf:"bytes,6,opt,name=linked_wallet,json=linkedWallet,proto3" json:"linked_wallet,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Balance) Reset()         { *m = Balance{} }
func (m *Balance) String() string { return proto.CompactTextString(m) }
func (*Balance) ProtoMessage()    {}
func (*Balance) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee25a00b628521b1, []int{1}
}

func (m *Balance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Balance.Unmarshal(m, b)
}
func (m *Balance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Balance.Marshal(b, m, deterministic)
}
func (m *Balance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Balance.Merge(m, src)
}
func (m *Balance) XXX_Size() int {
	return xxx_messageInfo_Balance.Size(m)
}
func (m *Balance) XXX_DiscardUnknown() {
	xxx_messageInfo_Balance.DiscardUnknown(m)
}

var xxx_messageInfo_Balance proto.InternalMessageInfo

func (m *Balance) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *Balance) GetApplicationId() string {
	if m != nil {
		return m.ApplicationId
	}
	return ""
}

func (m *Balance) GetBalanceDetails() *BalanceDetails {
	if m != nil {
		return m.BalanceDetails
	}
	return nil
}

func (m *Balance) GetPreCutoffDetails() *BalanceDetails {
	if m != nil {
		return m.PreCutoffDetails
	}
	return nil
}

func (m *Balance) GetType() BalanceType {
	if m != nil {
		return m.Type
	}
	return BalanceType_USER
}

func (m *Balance) GetLinkedWallet() string {
	if m != nil {
		return m.LinkedWallet
	}
	return ""
}

func init() {
	proto.RegisterEnum("pending_props_pb.BalanceType", BalanceType_name, BalanceType_value)
	proto.RegisterEnum("pending_props_pb.UpdateType", UpdateType_name, UpdateType_value)
	proto.RegisterType((*BalanceDetails)(nil), "pending_props_pb.BalanceDetails")
	proto.RegisterType((*Balance)(nil), "pending_props_pb.Balance")
}

func init() { proto.RegisterFile("balance.proto", fileDescriptor_ee25a00b628521b1) }

var fileDescriptor_ee25a00b628521b1 = []byte{
	// 451 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0x51, 0x6b, 0xdb, 0x30,
	0x14, 0x85, 0xeb, 0x34, 0x75, 0x9a, 0x9b, 0xc4, 0x73, 0x05, 0xdb, 0x34, 0xe8, 0xc0, 0xa4, 0x0c,
	0xb2, 0x3e, 0x04, 0xd6, 0xbd, 0xed, 0x2d, 0x69, 0xbd, 0x61, 0x08, 0x59, 0x70, 0x53, 0xfa, 0x28,
	0xe4, 0xe8, 0xa6, 0x35, 0x55, 0x2d, 0x61, 0x2b, 0x8c, 0xfe, 0xb7, 0xfd, 0x94, 0xfd, 0x98, 0x61,
	0xd9, 0x49, 0x9c, 0x8d, 0xc1, 0x1e, 0xef, 0x77, 0xce, 0xbd, 0x36, 0xe7, 0xd8, 0x30, 0x48, 0xb8,
	0xe4, 0xd9, 0x0a, 0xc7, 0x3a, 0x57, 0x46, 0x11, 0x5f, 0x63, 0x26, 0xd2, 0xec, 0x81, 0xe9, 0x5c,
	0xe9, 0x82, 0xe9, 0x64, 0xf8, 0xab, 0x05, 0xde, 0xb4, 0xf2, 0xdc, 0xa0, 0xe1, 0xa9, 0x2c, 0x08,
	0x85, 0x4e, 0x6d, 0xa3, 0x4e, 0xe0, 0x8c, 0xba, 0xf1, 0x76, 0x24, 0x17, 0x30, 0x30, 0xca, 0x70,
	0xc9, 0xb6, 0x7a, 0xcb, 0xea, 0x7d, 0x0b, 0x17, 0xb5, 0x69, 0x08, 0x7d, 0x93, 0xf3, 0xac, 0x58,
	0x63, 0xce, 0x13, 0x89, 0xf4, 0xb8, 0xf6, 0x34, 0x18, 0x79, 0x03, 0x6e, 0xa2, 0x32, 0x81, 0x82,
	0xb6, 0xad, 0x5a, 0x4f, 0xe4, 0x1c, 0xba, 0x02, 0x25, 0x3e, 0x70, 0x83, 0x82, 0x9e, 0x58, 0x69,
	0x0f, 0x48, 0x00, 0xbd, 0xdd, 0xb0, 0x54, 0xd4, 0xb5, 0x7a, 0x13, 0x95, 0xfb, 0x26, 0x7d, 0xc6,
	0xc2, 0xf0, 0x67, 0x4d, 0x3b, 0x81, 0x33, 0x3a, 0x8e, 0xf7, 0x80, 0x7c, 0x84, 0x33, 0xc9, 0x0b,
	0xc3, 0xd0, 0x3c, 0xb2, 0x44, 0xaa, 0xd5, 0x13, 0x4b, 0x05, 0x3d, 0xb5, 0x2e, 0xaf, 0x14, 0x42,
	0xf3, 0x38, 0x2d, 0x71, 0x24, 0xc8, 0x57, 0xf0, 0xad, 0x75, 0xa3, 0x05, 0x37, 0xc8, 0xcc, 0x8b,
	0x46, 0xda, 0x0d, 0x9c, 0x91, 0x77, 0x75, 0x3e, 0xfe, 0x33, 0xc3, 0xf1, 0x9d, 0x35, 0x2d, 0x5f,
	0x34, 0x56, 0x77, 0xf6, 0xf3, 0xf0, 0x67, 0x0b, 0x3a, 0x75, 0xbc, 0xe4, 0x2d, 0x74, 0x36, 0x05,
	0xe6, 0xe5, 0x43, 0xab, 0x5c, 0xdd, 0x72, 0x8c, 0x04, 0xf9, 0x00, 0x1e, 0xd7, 0x5a, 0xa6, 0x2b,
	0x6e, 0x52, 0x95, 0x95, 0x7a, 0x95, 0xeb, 0xa0, 0x41, 0x23, 0x41, 0x22, 0x78, 0x55, 0xb7, 0xc9,
	0x44, 0x55, 0x95, 0xcd, 0xb6, 0x77, 0x15, 0xfc, 0xfd, 0x4a, 0x87, 0x95, 0xc6, 0x5e, 0x72, 0x58,
	0xf1, 0x1c, 0x88, 0xce, 0x91, 0xad, 0x36, 0x46, 0xad, 0xd7, 0xbb, 0x6b, 0xed, 0xff, 0xbc, 0xe6,
	0xeb, 0x1c, 0xaf, 0xed, 0xea, 0xf6, 0xde, 0x27, 0x68, 0xdb, 0x88, 0x4e, 0x6c, 0x44, 0xef, 0xff,
	0x79, 0xc1, 0x66, 0x64, 0xad, 0xe5, 0xb7, 0x24, 0xd3, 0xec, 0x09, 0x05, 0xfb, 0xc1, 0xa5, 0x44,
	0x53, 0xd7, 0xd9, 0xaf, 0xe0, 0xbd, 0x65, 0x97, 0x17, 0xd0, 0x6b, 0x6c, 0x92, 0x53, 0x68, 0xdf,
	0xdd, 0x86, 0xb1, 0x7f, 0x44, 0x00, 0xdc, 0xfb, 0xc9, 0x6c, 0x16, 0x2e, 0x7d, 0xe7, 0xf2, 0x0b,
	0xc0, 0x3e, 0x71, 0xf2, 0x0e, 0x5e, 0x2f, 0xc2, 0xf9, 0x4d, 0x34, 0xff, 0xc6, 0x16, 0xf1, 0xf7,
	0xc5, 0x2d, 0x9b, 0x4e, 0x66, 0x93, 0xf9, 0x75, 0xe8, 0x1f, 0x91, 0x33, 0x18, 0x1c, 0x22, 0x27,
	0x71, 0xed, 0x7f, 0xf1, 0xf9, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x90, 0x6a, 0x87, 0xee, 0x28,
	0x03, 0x00, 0x00,
}
