package rpc

import (
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/propsproject/props-transaction-processor/core/state"
	"strings"
)

var transactionHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest) error {
	logger.Infof("Inputs=%v",strings.Join(request.Header.Inputs,","))
	logger.Infof("Outputs=%v",strings.Join(request.Header.Outputs,","))
	transaction, err := decodeRequest(rpcReq)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	return state.NewState(context).SaveTransactions(transaction)
}

var ISSUE = &MethodHandler{transactionHandle, pending_props_pb.Method_ISSUE.String()}
var REVOKE = &MethodHandler{transactionHandle, pending_props_pb.Method_REVOKE.String()}
var SETTLE = &MethodHandler{transactionHandle, pending_props_pb.Method_SETTLE.String()}
