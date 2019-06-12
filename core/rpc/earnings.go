package rpc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/propsproject/pending-props/core/eth-utils"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/pending-props/core/state"
	"github.com/propsproject/sawtooth-go-sdk/logging"
	"github.com/propsproject/sawtooth-go-sdk/processor"
)

var logger = logging.Get()

func decodeRequest(rpcReq *pending_props_pb.RPCRequest) (pending_props_pb.Earning, error) {
	var earning pending_props_pb.Earning
	err := getReqData(rpcReq.Params.GetData(), &earning)
	if err != nil {
		return earning, err
	}
	return earning, nil
}

func decodeActivityLogRequest(rpcReq *pending_props_pb.RPCRequest) (pending_props_pb.ActivityLog, error) {
	var activity pending_props_pb.ActivityLog
	err := getActivityLogReqData(rpcReq.Params.GetData(), &activity)
	if err != nil {
		return activity, err
	}
	return activity, nil
}

func decodeBalanceUpdateRequest(rpcReq *pending_props_pb.RPCRequest) (pending_props_pb.BalanceUpdate, error) {
	var balanceUpdate pending_props_pb.BalanceUpdate
	err := getBalanceUpdateReqData(rpcReq.Params.GetData(), &balanceUpdate)
	if err != nil {
		return balanceUpdate, err
	}
	return balanceUpdate, nil
}

func decodeLastEthBlockRequest(rpcReq *pending_props_pb.RPCRequest) (pending_props_pb.LastEthBlock, error) {
	var lastEthBlock pending_props_pb.LastEthBlock
	err := getLastBlockReqData(rpcReq.Params.GetData(), &lastEthBlock)
	if err != nil {
		return lastEthBlock, err
	}
	return lastEthBlock, nil
}

func decodeWalletLinkRequest(rpcReq *pending_props_pb.RPCRequest) (pending_props_pb.WalletToUser, error) {
	var walletToUser pending_props_pb.WalletToUser
	err := getLinkedWalletReqData(rpcReq.Params.GetData(), &walletToUser)
	if err != nil {
		return walletToUser, err
	}
	return walletToUser, nil
}

func getLinkedWalletReqData(data *any.Any, walletToUser *pending_props_pb.WalletToUser) error {
	err := ptypes.UnmarshalAny(data, walletToUser)
	if err != nil {
		return errors.New(fmt.Sprintf("could not unmarshal walletToUser proto (%s)", err.Error()))
	}
	return nil
}

func getLastBlockReqData(data *any.Any, lastEthBlock *pending_props_pb.LastEthBlock) error {
	err := ptypes.UnmarshalAny(data, lastEthBlock)
	if err != nil {
		return errors.New(fmt.Sprintf("could not unmarshal lastEthBlock proto (%s)", err.Error()))
	}
	return nil
}

func getBalanceUpdateReqData(data *any.Any, balanceUpdate *pending_props_pb.BalanceUpdate) error {
	err := ptypes.UnmarshalAny(data, balanceUpdate)
	if err != nil {
		return errors.New(fmt.Sprintf("could not unmarshal balanceUpdate proto (%s)", err.Error()))
	}
	return nil
}

func getActivityLogReqData(data *any.Any, activity *pending_props_pb.ActivityLog) error {
	err := ptypes.UnmarshalAny(data, activity)
	if err != nil {
		return errors.New(fmt.Sprintf("could not unmarshal activity log proto (%s)", err.Error()))
	}
	return nil
}

func getReqData(data *any.Any, earning *pending_props_pb.Earning) error {
	err := ptypes.UnmarshalAny(data, earning)
	if err != nil {
		return errors.New(fmt.Sprintf("could not unmarshal earning proto (%s)", err.Error()))
	}
	return nil
}

func pubToEthAddr(publickey string) (string, error) {
	var address string

	hexDecoded, err := hex.DecodeString(publickey)
	if err != nil {
		return "", fmt.Errorf("publickey string should be hex encoded (%s)", err)
	}

	if keyLength := len(hexDecoded); keyLength == 33 {
		ecdsaPub, err := crypto.DecompressPubkey(hexDecoded)
		if err != nil {
			return "", err
		}

		address = crypto.PubkeyToAddress(*ecdsaPub).String()
	} else if keyLength == 120 {
		ecdsaPub, err := crypto.UnmarshalPubkey(hexDecoded)
		if err != nil {
			return "", err
		}
		address = crypto.PubkeyToAddress(*ecdsaPub).String()
	} else {
		return "", errors.New("invalid key length")
	}

	return eth_utils.NormalizeAddress(address), nil
}

/*
func withHexPrefix(s string) string {
	if s[:2] != "0x" {
		return fmt.Sprintf("%s%s", "0x", s)
	}
	return s
}

func stripHexPrefix(s string) string {
	if s[:2] == "0x" {
		return s[2:]
	}
	return s
}
*/

func getAllEarnings(context *processor.Context, addresses ...string) (Earnings, error) {
	s := state.NewState(context)
	earnings, err := s.GetEarnings(addresses...)
	if err != nil {
		return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get earnings (%s)", err)}
	}

	if len(earnings) == 0 {
		return nil, &processor.InvalidTransactionError{Msg: "no earnings found"}
	}

	return earnings,  nil
}

