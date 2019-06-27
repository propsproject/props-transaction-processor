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

func (s *State) GetLastEthBlockData() (*pending_props_pb.LastEthBlock, error) {
	lastEthBlockAddress, _ := LastEthBlockAddress()
	existingLastEthBlockData, err := s.context.GetState([] string{lastEthBlockAddress})
	var lastEthBlockData pending_props_pb.LastEthBlock
	logger.Info("Last block data %v", lastEthBlockAddress)

	if err != nil {
		logger.Infof("Unable to fetch the lastEthBlockData state %s", lastEthBlockAddress)
		return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get eth block state %v (%s)", lastEthBlockAddress, err)}
	}

	if len(string(existingLastEthBlockData[lastEthBlockAddress])) > 0 {
		for _, value := range existingLastEthBlockData {
			err := proto.Unmarshal(value, &lastEthBlockData)
			if err != nil {
				return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto (lastEthBlockData) data (%s)", err)}
			}
		}
	} else {
		logger.Infof("Unable to fetch the lastEthBlockData state %s", lastEthBlockAddress)
		return nil, &processor.InvalidTransactionError{Msg: "Unable to fetch the lastBlockData"}
	}

	return &lastEthBlockData, nil
}

func (s *State) UpdateLastEthBlock(blockUpdate pending_props_pb.LastEthBlock) (error, *pending_props_pb.LastEthBlock) {
	ethBlockAddress, _ := LastEthBlockAddress()
	existingLastEthBlockData, err := s.context.GetState([]string{ethBlockAddress})
	var lastEthBlockData pending_props_pb.LastEthBlock

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


func (s *State) SaveLastEthBlockUpdate(blockUpdates ...pending_props_pb.LastEthBlock) error {
	stateUpdate := make(map[string][]byte)
	for _, blockUpdate := range blockUpdates {
		err, blockUpdateData := s.UpdateLastEthBlock(blockUpdate)

		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not blockUpdateData (%v)", err)}
		}

		blockUpdateBytes, err := proto.Marshal(blockUpdateData)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal blockUpdate proto"}
		}

		blockUpdateAddress, _ := LastEthBlockAddress()
		logger.Infof("Last Block Address: %v", blockUpdateAddress)
		stateUpdate[blockUpdateAddress] = blockUpdateBytes
		receiptBytes, err := json.Marshal(GetLastEthBlockUpdateReceipt(blockUpdateAddress, blockUpdateData.GetId()))
		if err != nil {
			logging.Get().Infof("unable to create new block update receipt (%s)", err)
		}

		err = s.context.AddReceiptData(receiptBytes)
		if err != nil {
			logging.Get().Infof("unable to add new block update receipt (%s)", err)
		}

		lastBlockUpdateEvent := pending_props_pb.LastEthBlockEvent{
			BlockId: blockUpdateData.GetId(),
			Message: fmt.Sprintf("last eth block id updated: %s", blockUpdateAddress),
			Timestamp: blockUpdateData.GetTimestamp(),
		}
		lastBlockUpdateAttr := []processor.Attribute{
			processor.Attribute{"block_id", strconv.FormatInt(blockUpdateData.GetId(), 10)},
			processor.Attribute{"event_type", pending_props_pb.EventType_LastEthBlockUpdated.String()},
			processor.Attribute{"timestamp", strconv.FormatInt(blockUpdateData.GetTimestamp(), 10)},
		}
		s.AddLastEthBlockUpdateEvent(lastBlockUpdateEvent, "pending-props:lastethblockupdate", lastBlockUpdateAttr...)
	}

	_, err := s.context.SetState(stateUpdate)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
	}

	return nil
}