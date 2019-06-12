package core

import (
	"github.com/propsproject/pending-props/core/rpc"
	"github.com/propsproject/pending-props/core/state"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
)

// TransactionHandler ...
type TransactionHandler struct {
	FName     string   `json:"familyName"`
	FVersions []string `json:"familyVersions"`
	NSpace    []string `json:"nameSpace"`
	RPCClient *rpc.RPCClient
}

// FamilyName ...
const FamilyName = "pending-earnings" // move to configuration

// FamilyVersions ...
var FamilyVersions = []string{"1.0"} // move to configuration

// FamilyName ...
func (t *TransactionHandler) FamilyName() string {
	return t.FName
}

// FamilyVersions ...
func (t *TransactionHandler) FamilyVersions() []string {
	return t.FVersions
}

// Namespaces ...
func (t *TransactionHandler) Namespaces() []string {
	return t.NSpace
}

// Apply ...
func (t *TransactionHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	err := t.RPCClient.DelegateMethod(request, context)
	return err
}

// NewTransactionHandler returns a new transaction handler
func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{
		FName:     FamilyName,
		FVersions: FamilyVersions,
		NSpace:    state.NamespaceManager.Namespaces(),
		RPCClient: rpc.NewClient(),
	}
}
