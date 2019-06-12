package rpc

import (
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/pending-props/core/state"
)

var ethLastBlockUpdateHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest) error {

	lastEthBlockUpdate, err := decodeLastEthBlockRequest(rpcReq)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	return state.NewState(context).SaveLastEthBlockUpdate(lastEthBlockUpdate)
}

var ETH_LAST_BLOCK_UPDATE = &MethodHandler{ethLastBlockUpdateHandle, pending_props_pb.Method_LAST_ETH_BLOCK_UPDATE.String()}
