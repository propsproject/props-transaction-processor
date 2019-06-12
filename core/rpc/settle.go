package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/pending-props/core/state"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
)

var settleHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest) error {
	var reqData SettlementRPCPayload
	//logger.Infof("Inputs=%v",strings.Join(request.Header.Inputs,","))
	//logger.Infof("Outputs=%v",strings.Join(request.Header.Inputs,","))
	err := json.Unmarshal(rpcReq.GetParams().GetData().GetValue(), &reqData)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal earning addresses in rpc request (%s)", err)}
	}
	err = reqData.Validate(context)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("failed payload validation for SETTLE (%s)", err)}
	}

	earnings, err := getAllEarnings(context, reqData.PendingAddresses...)
	if err != nil {
		return err
	}

	if errs := NewEarningsProcessor().Init(earnings).ProcessSettlements(reqData.EthTransactionHash).Errs(); len(errs) > 0 {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("unexpected errors when settling props (%s)", errs)}
	}

	return state.NewState(context).SettleEarnings(reqData.Timestamp, reqData.EthTransactionHash, earnings...)
}
var SETTLE = &MethodHandler{settleHandle, pending_props_pb.Method_SETTLE.String()}
