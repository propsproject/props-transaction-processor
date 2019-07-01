package rpc

import (
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"fmt"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

type MethodHandler struct {
	Handle func(*processor_pb2.TpProcessRequest, *processor.Context, *pending_props_pb.RPCRequest, string) error
	Method string
}

type RPCClient struct {
	MethodHandlers map[string]*MethodHandler
}

func (r *RPCClient) registerMethod(handler *MethodHandler) *RPCClient {
	r.MethodHandlers[handler.Method] = handler
	return r
}

func (r *RPCClient) DelegateMethod(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	var rpcRequest pending_props_pb.RPCRequest
	err := proto.Unmarshal(request.GetPayload(), &rpcRequest)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: "malformed payload data"}
	}

	return r.delegate(request, context, rpcRequest)
}

func (r *RPCClient) delegate(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcRequest pending_props_pb.RPCRequest) error {
	method := rpcRequest.GetMethod().String()
	// only allowed signers
	publicKey := request.GetHeader().GetSignerPublicKey()
	address, err := pubToEthAddr(publicKey)
	//logger.Infof("Address: %v", address)
	if err != nil {
		logger.Infof("Invalid key used to sign can't convert to Ethereum address: %v", err)
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid key used to sign can't convert to Ethereum address: %v", err)}
	} else {
		allowedSignersMap := viper.GetStringMapString("valid_signers_addresses")
		if _, ok := allowedSignersMap[address]; !ok {
			logger.Infof("Invalid signer address: %v", address)
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid signer address: %v", address)}
		}
	}
	if methodHandler, exists := r.MethodHandlers[method]; exists {
		return methodHandler.Handle(request, context, &rpcRequest, address)
	}

	return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not determine RPC method: %v", method)}
}

func NewClient() *RPCClient {
	client := &RPCClient{make(map[string]*MethodHandler)}
	return client.registerMethod(ISSUE).registerMethod(REVOKE).registerMethod(SETTLE).registerMethod(BALANCE_UPDATE).registerMethod(ETH_LAST_BLOCK_UPDATE).registerMethod(WALLET_LINK).registerMethod(ACTIVITY_LOG)
}
