package rpc

import (
	"testing"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/golang/protobuf/ptypes"
)

const (
	earningTypeUrl = "github.com/propsproject/pending-props/protos/pending_props_pb.Earning"
)

func getTestRpcReq(testSig string) *pending_props_pb.RPCRequest {
	rpcRequest := new(pending_props_pb.RPCRequest)
	rpcRequest.Params = new(pending_props_pb.Params)
	earning := pending_props_pb.Earning{
		Details:   new(pending_props_pb.EarningDetails),
		Signature: testSig,
	}
	rpcRequest.Params.Data, _ = ptypes.MarshalAny(&earning)
	rpcRequest.Params.Data.TypeUrl = earningTypeUrl

	return rpcRequest
}

func TestPubToEthAddr(t *testing.T)  {
	//testSig := "76edeb9966af0d79dd4e369cf02ccd9ce15f3947ef03bd674629bb3e0cf64f24034d06ccc73bd2f12605c197a2a042dd4e012d5b29ab46fd82f37c59ca6b6566"
	testPubKey := "02ccb8bc17397c55242a27d1681bf48b5b40a734205760882cd83f92aca4f1cf45"
	expectedAddr := "0x4a3595ddb0dee4a0f053bae4734312ec4f9863e9"

	addr, err := pubToEthAddr(testPubKey)
	if err != nil {
		t.Fatalf("expected err to be nil got (%s)", err)
	}

	if addr != expectedAddr {
		t.Fatalf("expected eth address to be (%s) got (%s)", expectedAddr, addr)
	}

}

func TestDecodeRequest(t *testing.T) {
	t.Run("ExpectErrorNil", expectErrIsNil)
	t.Run("ExpectError", expectErrorInvalidReq)
}

func expectErrIsNil(t *testing.T) {
	testSig := "a2703d9c0d4c996c993fd7611380bbda914bc6ee4a90d655225211cbfce31d3779f0c357e87d95c31685ca5428f57fbd970c5830b34b76db598c38cda3f1049c"
	rpcRequest := getTestRpcReq(testSig)
	_, err := decodeRequest(rpcRequest)
	if err != nil {
		t.Fatalf("expected err to be nil got (%s)", err)
	}
}

func expectErrorInvalidReq(t *testing.T) {
	rpcRequest := new(pending_props_pb.RPCRequest)
	rpcRequest.Params = new(pending_props_pb.Params)
	rpcRequest.Method = pending_props_pb.Method_ISSUE
	_, err := decodeRequest(rpcRequest)
	if err == nil {
		t.Fatalf("expected err to be defined got (%s)", err)
	}
}