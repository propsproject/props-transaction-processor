package rpc

import (
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/propsproject/props-transaction-processor/core/state"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
)

var activityLogHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest, address string) error {
	activities, err := decodeActivityLogRequest(rpcReq)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: err.Error()}
	}
	//if address != activities.GetApplicationId() {
	//	logger.Infof("Signer address %v does not match applicationId %v", address, activities.GetApplicationId())
	//	return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Signer address %v does not match applicationId %v", address, activities.GetApplicationId())}
	//}
	return state.NewState(context).SaveActivityLog(activities)
}

var ACTIVITY_LOG = &MethodHandler{activityLogHandle, pending_props_pb.Method_ACTIVITY_LOG.String()}