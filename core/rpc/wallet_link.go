package rpc

import (
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/propsproject/props-transaction-processor/core/state"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
)

var walletLinkHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest, address string) error {
	walletLinks, err := decodeWalletLinkRequest(rpcReq)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	//if address != walletLinks.GetUsers()[0].GetApplicationId() {
	//	logger.Infof("Signer address %v does not match applicationId %v", address, walletLinks.GetUsers()[0].GetApplicationId())
	//	return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Signer address %v does not match applicationId %v", address, walletLinks.GetUsers()[0].GetApplicationId())}
	//}
	return state.NewState(context).SaveWalletLink(walletLinks)
}

var WALLET_LINK = &MethodHandler{walletLinkHandle, pending_props_pb.Method_WALLET_LINK.String()}
