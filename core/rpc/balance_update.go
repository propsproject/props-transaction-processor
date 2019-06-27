package rpc

import (
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/propsproject/props-transaction-processor/core/state"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
)

var balanceUpdateHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest) error {
	balanceUpdate, err := decodeBalanceUpdateRequest(rpcReq)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	return state.NewState(context).SaveBalanceUpdate(balanceUpdate)
}

var BALANCE_UPDATE = &MethodHandler{balanceUpdateHandle, pending_props_pb.Method_BALANCE_UPDATE.String()}
