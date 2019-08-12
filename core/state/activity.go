package state

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/propsproject/props-transaction-processor/core/eth-utils"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
)

// logDate is in format YYYYMMDD
// blockTimestamp is a unix timestamp (etc: 1558607755) of the latest eth block
func ValidateActivityTimestamp(logDate int, blockDate int64) error {

	return nil
}

func (s *State) SaveActivityLog(activities ...pending_props_pb.ActivityLog) error {
	stateUpdate := make(map[string][]byte)

	lastEthBlockData, err := s.GetLastEthBlockData()

	if err != nil {
		logger.Infof("Unable to get the last eth block data will use 0...")
		lastEthBlockData = &pending_props_pb.LastEthBlock{
			Id: 0,
			Timestamp:0,
		}
	}

	// Fetch the timestamp of the last eth block address
	for _, activity := range activities {
		calculatedActivityRewardsDay := eth_utils.CalculateRewardsDay(activity.GetTimestamp())
		if int64(activity.GetDate()) != calculatedActivityRewardsDay {
			logger.Infof("Activity rewards day does not match timestamp rewards day %v, %v", activity.GetDate(), calculatedActivityRewardsDay)
		}
		calculatedEthRewardsDay := eth_utils.CalculateRewardsDay(lastEthBlockData.GetTimestamp())
		// allow for deviation by 1 day in the future max
		if (int64(activity.GetDate()) - calculatedEthRewardsDay) <= 1 {
			logger.Infof("Activity rewards day does not match eth block rewards day %v, %v", activity.GetDate(), calculatedEthRewardsDay)
		}
		activityAddress, _ := ActivityLogAddress(activity)

		activityData, err := s.context.GetState([]string{activityAddress})

		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err)}
		}

		if len(string(activityData[activityAddress])) == 0 {
			// check if there's a balance object for this user
			balanceAddress, _ := BalanceAddressByAppUser(activity.GetApplicationId(), activity.GetUserId())
			balanceStateData, err := s.context.GetState([]string{balanceAddress})
			var balanceData pending_props_pb.Balance
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err)}
			}
			if len(string(balanceStateData[balanceAddress])) == 0 {
				logger.Infof("Balance object was not found while getting state %v for %v, %v", balanceAddress, activity.GetApplicationId(), activity.GetUserId())
			} else {
				// update existing balance
				for _, value := range balanceStateData {

					err := proto.Unmarshal(value, &balanceData)
					if err != nil {
						return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto data (%s)", err)}
					}
				}
			}
			activity.Balance = &balanceData
			b, err := proto.Marshal(&activity)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: "could not marshal activity proto"}
			}
			stateUpdate[activityAddress] = b
		} else {
			logger.Infof(fmt.Sprintf("Activity log already exists for applicationId=%v, userId=%v, date=%v, skipping...", activity.ApplicationId, activity.UserId, activity.Date))
		}
	}

	if len(stateUpdate) > 0 {
		_, err := s.context.SetState(stateUpdate)
		logger.Info("Updating the activity state ...")
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
		}
	} else {
		logger.Infof("No activity states to update")
	}

	return nil
}