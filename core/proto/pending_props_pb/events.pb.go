// Code generated by protoc-gen-go. DO NOT EDIT.
// source: events.proto

package pending_props_pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EventType int32

const (
	EventType_EarningIssued       EventType = 0
	EventType_EarningRevoked      EventType = 1
	EventType_EarningSettled      EventType = 2
	EventType_BalanceUpdated      EventType = 3
	EventType_LastEthBlockUpdated EventType = 4
	EventType_WalletLinked        EventType = 5
	EventType_WalletUnlinked      EventType = 6
	EventType_TransactionAdded    EventType = 7
)

var EventType_name = map[int32]string{
	0: "EarningIssued",
	1: "EarningRevoked",
	2: "EarningSettled",
	3: "BalanceUpdated",
	4: "LastEthBlockUpdated",
	5: "WalletLinked",
	6: "WalletUnlinked",
	7: "TransactionAdded",
}
var EventType_value = map[string]int32{
	"EarningIssued":       0,
	"EarningRevoked":      1,
	"EarningSettled":      2,
	"BalanceUpdated":      3,
	"LastEthBlockUpdated": 4,
	"WalletLinked":        5,
	"WalletUnlinked":      6,
	"TransactionAdded":    7,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}
func (EventType) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

type TransactionEvent struct {
	Transaction  *Transaction `protobuf:"bytes,1,opt,name=transaction" json:"transaction,omitempty"`
	Type         Method       `protobuf:"varint,2,opt,name=type,enum=pending_props_pb.Method" json:"type,omitempty"`
	StateAddress string       `protobuf:"bytes,3,opt,name=stateAddress" json:"stateAddress,omitempty"`
	Message      string       `protobuf:"bytes,4,opt,name=message" json:"message,omitempty"`
	Description  string       `protobuf:"bytes,5,opt,name=description" json:"description,omitempty"`
}

func (m *TransactionEvent) Reset()                    { *m = TransactionEvent{} }
func (m *TransactionEvent) String() string            { return proto.CompactTextString(m) }
func (*TransactionEvent) ProtoMessage()               {}
func (*TransactionEvent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *TransactionEvent) GetTransaction() *Transaction {
	if m != nil {
		return m.Transaction
	}
	return nil
}

func (m *TransactionEvent) GetType() Method {
	if m != nil {
		return m.Type
	}
	return Method_ISSUE
}

func (m *TransactionEvent) GetStateAddress() string {
	if m != nil {
		return m.StateAddress
	}
	return ""
}

func (m *TransactionEvent) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *TransactionEvent) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type BalanceEvent struct {
	Balance     *Balance `protobuf:"bytes,1,opt,name=balance" json:"balance,omitempty"`
	Message     string   `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Description string   `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
}

func (m *BalanceEvent) Reset()                    { *m = BalanceEvent{} }
func (m *BalanceEvent) String() string            { return proto.CompactTextString(m) }
func (*BalanceEvent) ProtoMessage()               {}
func (*BalanceEvent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *BalanceEvent) GetBalance() *Balance {
	if m != nil {
		return m.Balance
	}
	return nil
}

func (m *BalanceEvent) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *BalanceEvent) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type LastEthBlockEvent struct {
	BlockId   int64  `protobuf:"varint,1,opt,name=blockId" json:"blockId,omitempty"`
	Message   string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Timestamp int64  `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *LastEthBlockEvent) Reset()                    { *m = LastEthBlockEvent{} }
func (m *LastEthBlockEvent) String() string            { return proto.CompactTextString(m) }
func (*LastEthBlockEvent) ProtoMessage()               {}
func (*LastEthBlockEvent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *LastEthBlockEvent) GetBlockId() int64 {
	if m != nil {
		return m.BlockId
	}
	return 0
}

func (m *LastEthBlockEvent) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *LastEthBlockEvent) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type WalletLinkedEvent struct {
	User          *ApplicationUser `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
	WalletToUsers *WalletToUser    `protobuf:"bytes,2,opt,name=walletToUsers" json:"walletToUsers,omitempty"`
	Message       string           `protobuf:"bytes,3,opt,name=message" json:"message,omitempty"`
}

func (m *WalletLinkedEvent) Reset()                    { *m = WalletLinkedEvent{} }
func (m *WalletLinkedEvent) String() string            { return proto.CompactTextString(m) }
func (*WalletLinkedEvent) ProtoMessage()               {}
func (*WalletLinkedEvent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *WalletLinkedEvent) GetUser() *ApplicationUser {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *WalletLinkedEvent) GetWalletToUsers() *WalletToUser {
	if m != nil {
		return m.WalletToUsers
	}
	return nil
}

func (m *WalletLinkedEvent) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type WalletUnlinkedEvent struct {
	User          *ApplicationUser `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
	WalletToUsers *WalletToUser    `protobuf:"bytes,2,opt,name=walletToUsers" json:"walletToUsers,omitempty"`
	Message       string           `protobuf:"bytes,3,opt,name=message" json:"message,omitempty"`
}

func (m *WalletUnlinkedEvent) Reset()                    { *m = WalletUnlinkedEvent{} }
func (m *WalletUnlinkedEvent) String() string            { return proto.CompactTextString(m) }
func (*WalletUnlinkedEvent) ProtoMessage()               {}
func (*WalletUnlinkedEvent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *WalletUnlinkedEvent) GetUser() *ApplicationUser {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *WalletUnlinkedEvent) GetWalletToUsers() *WalletToUser {
	if m != nil {
		return m.WalletToUsers
	}
	return nil
}

func (m *WalletUnlinkedEvent) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*TransactionEvent)(nil), "pending_props_pb.TransactionEvent")
	proto.RegisterType((*BalanceEvent)(nil), "pending_props_pb.BalanceEvent")
	proto.RegisterType((*LastEthBlockEvent)(nil), "pending_props_pb.LastEthBlockEvent")
	proto.RegisterType((*WalletLinkedEvent)(nil), "pending_props_pb.WalletLinkedEvent")
	proto.RegisterType((*WalletUnlinkedEvent)(nil), "pending_props_pb.WalletUnlinkedEvent")
	proto.RegisterEnum("pending_props_pb.EventType", EventType_name, EventType_value)
}

func init() { proto.RegisterFile("events.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 465 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x53, 0xc1, 0x6a, 0xdb, 0x40,
	0x10, 0xed, 0x5a, 0x4e, 0x8c, 0xc7, 0x76, 0x58, 0x6f, 0x0a, 0x55, 0x4d, 0x5b, 0x5c, 0x9d, 0x4c,
	0x29, 0x3e, 0x38, 0xf4, 0x5c, 0x1c, 0xea, 0x43, 0x20, 0xbd, 0x6c, 0x6d, 0x72, 0x0c, 0x6b, 0xed,
	0xe0, 0x08, 0xcb, 0xbb, 0x8b, 0x76, 0x93, 0x92, 0x53, 0xff, 0xa7, 0x14, 0xfa, 0x4d, 0xfd, 0x93,
	0xa2, 0x95, 0x4c, 0x56, 0x15, 0xe4, 0xdc, 0xe3, 0xbc, 0x79, 0x4f, 0xf3, 0xde, 0x68, 0x16, 0x86,
	0xf8, 0x80, 0xca, 0xd9, 0xb9, 0x29, 0xb4, 0xd3, 0x8c, 0x1a, 0x54, 0x32, 0x53, 0xbb, 0x5b, 0x53,
	0x68, 0x63, 0x6f, 0xcd, 0x76, 0x32, 0x32, 0xe2, 0x31, 0xd7, 0x42, 0x56, 0x84, 0xc9, 0x68, 0x2b,
	0x72, 0xa1, 0x52, 0xac, 0xcb, 0xc1, 0xbd, 0xc5, 0xa2, 0x16, 0x4f, 0xc6, 0xae, 0x10, 0xca, 0x8a,
	0xd4, 0x65, 0x5a, 0x55, 0x50, 0xf2, 0x87, 0x00, 0x5d, 0x3f, 0xa1, 0xab, 0x72, 0x16, 0xfb, 0x0c,
	0x83, 0x80, 0x19, 0x93, 0x29, 0x99, 0x0d, 0x16, 0x6f, 0xe7, 0xff, 0x8e, 0x9e, 0x07, 0x42, 0x1e,
	0x2a, 0xd8, 0x47, 0xe8, 0xba, 0x47, 0x83, 0x71, 0x67, 0x4a, 0x66, 0x67, 0x8b, 0xb8, 0xad, 0xfc,
	0x8a, 0xee, 0x4e, 0x4b, 0xee, 0x59, 0x2c, 0x81, 0xa1, 0x75, 0xc2, 0xe1, 0x52, 0xca, 0x02, 0xad,
	0x8d, 0xa3, 0x29, 0x99, 0xf5, 0x79, 0x03, 0x63, 0x31, 0xf4, 0x0e, 0x68, 0xad, 0xd8, 0x61, 0xdc,
	0xf5, 0xed, 0x63, 0xc9, 0xa6, 0x30, 0x90, 0x68, 0xd3, 0x22, 0x33, 0xde, 0xec, 0x89, 0xef, 0x86,
	0x50, 0xf2, 0x03, 0x86, 0x97, 0xd5, 0x52, 0xaa, 0x78, 0x17, 0xd0, 0xab, 0x97, 0x54, 0x47, 0x7b,
	0xdd, 0x36, 0x58, 0x0b, 0xf8, 0x91, 0x19, 0x1a, 0xe8, 0x3c, 0x6b, 0x20, 0x6a, 0x1b, 0x40, 0x18,
	0x5f, 0x0b, 0xeb, 0x56, 0xee, 0xee, 0x32, 0xd7, 0xe9, 0xbe, 0x72, 0x11, 0x43, 0x6f, 0x5b, 0x56,
	0x57, 0xd2, 0xbb, 0x88, 0xf8, 0xb1, 0x7c, 0x66, 0xd4, 0x1b, 0xe8, 0xbb, 0xec, 0x80, 0xd6, 0x89,
	0x83, 0xf1, 0x83, 0x22, 0xfe, 0x04, 0x24, 0x3f, 0x09, 0x8c, 0x6f, 0x44, 0x9e, 0xa3, 0xbb, 0xce,
	0xd4, 0x1e, 0x65, 0x35, 0xe7, 0x13, 0x74, 0xcb, 0x1b, 0xa8, 0xa3, 0xbe, 0x6f, 0x47, 0x5d, 0x1a,
	0x93, 0x67, 0xa9, 0x28, 0x9d, 0x6e, 0x2c, 0x16, 0xdc, 0xd3, 0xd9, 0x17, 0x18, 0x7d, 0xf7, 0xdf,
	0x5a, 0xeb, 0x12, 0xb5, 0xde, 0xca, 0x60, 0xf1, 0xae, 0xad, 0xbf, 0x09, 0x68, 0xbc, 0x29, 0x0a,
	0xa3, 0x44, 0x8d, 0x28, 0xc9, 0x2f, 0x02, 0xe7, 0x95, 0x72, 0xa3, 0xf2, 0xff, 0xde, 0xee, 0x87,
	0xdf, 0x04, 0xfa, 0xde, 0xe0, 0xba, 0xbc, 0xd8, 0x31, 0x8c, 0x56, 0xa2, 0x50, 0x99, 0xda, 0x5d,
	0x59, 0x7b, 0x8f, 0x92, 0xbe, 0x60, 0x0c, 0xce, 0x6a, 0x88, 0xe3, 0x83, 0xde, 0xa3, 0xa4, 0x24,
	0xc0, 0xbe, 0xa1, 0x73, 0x39, 0x4a, 0xda, 0x29, 0xb1, 0xfa, 0xb6, 0x36, 0x46, 0x0a, 0x87, 0x92,
	0x46, 0xec, 0x15, 0x9c, 0x87, 0xf7, 0x71, 0x6c, 0x74, 0x19, 0x85, 0x61, 0xf8, 0x43, 0xe9, 0x49,
	0x29, 0x6f, 0x6e, 0x8d, 0x9e, 0xb2, 0x97, 0x8d, 0x27, 0xbc, 0x94, 0x12, 0x25, 0xed, 0x6d, 0x4f,
	0xfd, 0x03, 0xbf, 0xf8, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xca, 0xf9, 0x80, 0x81, 0x40, 0x04, 0x00,
	0x00,
}
