package rpc

import (
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/pending-props/core/state"
)

var issueHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest) error {
	earning, err := decodeRequest(rpcReq)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	return state.NewState(context).SavePendingEarnings(earning)
}

var ISSUE = &MethodHandler{issueHandle, pending_props_pb.Method_ISSUE.String()}
