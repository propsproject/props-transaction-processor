package rpc

import (
	"testing"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/propsproject/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

func TestRPCClient(t *testing.T) {
	t.Run("NewClient", newClient)
	t.Run("registerMethod", registerMethodHandler)
	t.Run("delegate", delegateReq)
}

func newClient(t *testing.T) {
	client := NewClient()

	if _, ok := client.MethodHandlers["ISSUE"]; !ok {
		t.Fatalf("expected issue handler to registered")
	}
}

func registerMethodHandler(t *testing.T) {
	handler := new(MethodHandler)
	handler.Method = "test"

	client := NewClient()
	client.registerMethod(handler)

	if _, ok := client.MethodHandlers[handler.Method]; !ok {
		t.Fatalf("expected handler to register method (test)")
	}
}

func delegateReq(t *testing.T) {
	testSig := "a2703d9c0d4c996c993fd7611380bbda914bc6ee4a90d655225211cbfce31d3779f0c357e87d95c31685ca5428f57fbd970c5830b34b76db598c38cda3f1049c"
	//testSigPubKey := "02ccb8bc17397c55242a27d1681bf48b5b40a734205760882cd83f92aca4f1cf45"
	//testSigHash := "02abf9b2da335be9d5200f976d8ce7354f1dd3e7d092c0ba981af9431da03c251b806f5c1a07b62bdd700f140dcb89fef57f242435fec5cb7b4984f152ab6945"

	client := NewClient()
	handler := new(MethodHandler)
	handler.Method = pending_props_pb.Method_REVOKE.String()

	rpcRequest := getTestRpcReq(testSig)
	rpcRequest.Method = pending_props_pb.Method_REVOKE
	tpReq := new(processor_pb2.TpProcessRequest)
	tpReq.Payload, _ = proto.Marshal(rpcRequest)

	handler.Handle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcRequest *pending_props_pb.RPCRequest) error {
		return errors.New("INVOKED")
	}

	client.registerMethod(handler)
	res := client.delegate(tpReq, new(processor.Context), *rpcRequest)

	if res.Error() != "INVOKED" {
		t.Fatalf("expected handler to be invoked 1 times")
	}
}