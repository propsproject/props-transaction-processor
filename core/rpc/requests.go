package rpc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/propsproject/props-transaction-processor/core/eth-utils"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
)

var logger = logging.Get()

func decodeRequest(rpcReq *pending_props_pb.RPCRequest) (pending_props_pb.Transaction, error) {
	var transaction pending_props_pb.Transaction
	err := getReqData(rpcReq.Params.GetData(), &transaction)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
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

func decodeRewardEntityUpdateRequest(rpcReq *pending_props_pb.RPCRequest) (pending_props_pb.RewardEntity, error) {
	var rewardEntity pending_props_pb.RewardEntity
	err := getRewardEntityReqData(rpcReq.Params.GetData(), &rewardEntity)
	if err != nil {
		return rewardEntity, err
	}
	return rewardEntity, nil
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

func getRewardEntityReqData(data *any.Any, rewardEntity *pending_props_pb.RewardEntity) error {
	err := ptypes.UnmarshalAny(data, rewardEntity)
	if err != nil {
		return errors.New(fmt.Sprintf("could not unmarshal reward entity proto (%s)", err.Error()))
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

func getReqData(data *any.Any, transaction *pending_props_pb.Transaction) error {
	err := ptypes.UnmarshalAny(data, transaction)
	if err != nil {
		return errors.New(fmt.Sprintf("could not unmarshal transaction proto (%s)", err.Error()))
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
