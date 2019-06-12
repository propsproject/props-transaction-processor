package rpc

import (
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"encoding/json"
	"fmt"
	"github.com/propsproject/pending-props/core/state"
)

var revokeHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest) error {
	var reqData RevokeRPCPaylod
	err := json.Unmarshal(rpcReq.GetParams().GetData().GetValue(), &reqData)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal earning addresses in rpc request (%s)", err)}
	}

	err = reqData.Validate()
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("failed payload validation for REVOKE (%s)", err)}
	}

	appEthAddr, err := pubToEthAddr(request.GetHeader().GetSignerPublicKey())
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not ethereum address from publickey in rpc request (%s)", err)}
	}

	earnings, err := getAllEarnings(context, reqData.Addresses...)
	if err != nil {
		return err
	}

	_ = NewEarningsProcessor().Init(earnings).ProcessRevocations(appEthAddr)
	return state.NewState(context).SaveRevokedEarnings(reqData.Timestamp, earnings...)
}
var REVOKE = &MethodHandler{revokeHandle, pending_props_pb.Method_REVOKE.String()}
