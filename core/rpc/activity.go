package rpc

import (
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/pending-props/core/state"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
)

var activityLogHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest) error {
	activities, err := decodeActivityLogRequest(rpcReq)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	return state.NewState(context).SaveActivityLog(activities)
}

var ACTIVITY_LOG = &MethodHandler{activityLogHandle, pending_props_pb.Method_ACTIVITY_LOG.String()}