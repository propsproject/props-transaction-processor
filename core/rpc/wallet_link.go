package rpc

import (
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/pending-props/core/state"
)

var walletLinkHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest) error {
	walletLinks, err := decodeWalletLinkRequest(rpcReq)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}

	return state.NewState(context).SaveWalletLink(walletLinks)
}

var WALLET_LINK = &MethodHandler{walletLinkHandle, pending_props_pb.Method_WALLET_LINK.String()}
