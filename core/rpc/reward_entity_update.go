package rpc

import (
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/propsproject/props-transaction-processor/core/state"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
)

var rewardEntityUpdateHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest, address string) error {
	rewardEntityUpdates, err := decodeRewardEntityUpdateRequest(rpcReq)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	//if address != activities.GetApplicationId() {
	//	logger.Infof("Signer address %v does not match applicationId %v", address, activities.GetApplicationId())
	//	return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Signer address %v does not match applicationId %v", address, activities.GetApplicationId())}
	//}
	return state.NewState(context).SaveRewardEntity(rewardEntityUpdates)
}

var REWARD_ENTITY_UPDATE = &MethodHandler{rewardEntityUpdateHandle, pending_props_pb.Method_REWARD_ENTITY_UPDATE.String()}