package state

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"strconv"
)

func (s *State) SaveRewardEntity(rewardEntityUpdate pending_props_pb.RewardEntity) (error, *pending_props_pb.RewardEntity) {
	rewardEntityAddress, _ := RewardEntityAddress(rewardEntityUpdate)


	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Unable to get state of ethblockaddress %v", ethBlockAddress)}, nil
	}

	if len(string(existingLastEthBlockData[ethBlockAddress])) == 0 {
		logger.Infof("Error / Not Found while getting state ethBlockAddress %v, %v", ethBlockAddress, err)
		newLastEthBlock := pending_props_pb.LastEthBlock{
			Id: blockUpdate.GetId(),
			Timestamp: blockUpdate.GetTimestamp(),
		}
		lastEthBlockData = newLastEthBlock
	} else {
		// update existing last eth block
		for _, value := range existingLastEthBlockData {

			err := proto.Unmarshal(value, &lastEthBlockData)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto (lastEthBlockData) data (%s)", err)}, nil
			}
		}
		if lastEthBlockData.GetId() > blockUpdate.GetId() {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("last block can be smaller than previouly stored one %v > %v", lastEthBlockData.GetId(), blockUpdate.GetId())}, nil
		}
	}
	return nil, &blockUpdate
}


func (s *State) SaveRewardEntity(rewardEntityUpdates ...pending_props_pb.RewardEntity) error {
	stateUpdate := make(map[string][]byte)
	for _, rewardEntityUpdate := range rewardEntityUpdates {

		rewardEntityBytes, err := proto.Marshal(&rewardEntityUpdate)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal reward entity proto"}
		}
		rewardEntityAddressBySidechainAddress, _ := RewardEntityAddressBySidechainAddress(rewardEntityUpdate)
		rewardEntityAddressByRewardsAddress, _ := RewardEntityAddressByRewardsAddress(rewardEntityUpdate)
		logger.Infof("Reward Entity Addresses: %v, %v", rewardEntityAddressBySidechainAddress, rewardEntityAddressByRewardsAddress)
		stateUpdate[rewardEntityAddressBySidechainAddress] = rewardEntityBytes
		stateUpdate[rewardEntityAddressByRewardsAddress] = rewardEntityBytes
		receiptBytes, err := json.Marshal(GetRewardEntityUpdateReceipt(rewardEntityUpdate.GetName(), rewardEntityUpdate.GetAddress(), rewardEntityUpdate.GetRewardsAddress(), rewardEntityUpdate.GetSidechainAddress() ))
		if err != nil {
			logging.Get().Infof("unable to create new reward entity update receipt (%s)", err)
		}

		err = s.context.AddReceiptData(receiptBytes)
		if err != nil {
			logging.Get().Infof("unable to add new reward entity update receipt (%s)", err)
		}

		rewardEntityUpdateEvent := pending_props_pb.RewardEntityUpdateEvent{
			Entity: &rewardEntityUpdate,
			Message: fmt.Sprintf("reward entity %s updated: addresses %v, %v", rewardEntityAddressBySidechainAddress, rewardEntityAddressByRewardsAddress),
		}
		rewardEntityUpdateAttr := []processor.Attribute{
			processor.Attribute{"name", rewardEntityUpdate.GetName()},
			processor.Attribute{"event_type", pending_props_pb.EventType_RewardEntityUpdated.String()},
			processor.Attribute{"address", rewardEntityUpdate.GetAddress()},
			processor.Attribute{"rewards_address", rewardEntityUpdate.GetRewardsAddress()},
			processor.Attribute{"sidechain_address", rewardEntityUpdate.GetSidechainAddress()},
		}
		s.AddLastEthBlockUpdateEvent(rewardEntityUpdateEvent, "pending-props:rewardentityupdate", rewardEntityUpdateAttr...)
	}

	_, err := s.context.SetState(stateUpdate)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
	}

	return nil
}