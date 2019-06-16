package state

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"strconv"
	"time"
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

	var bufferTimeSec int64 = 180

	// Build the current date
	blockDateStr := time.Unix(lastEthBlockData.Timestamp, 0).Format("20060102")
	blockDate64, _ := strconv.ParseInt(blockDateStr, 10, 32)
	blockDate := int32(blockDate64)

	// Build the current date + bufferTimeSec, this might be the next day, allow it, since the lastEthBlockDate is behind with more or less 2 minutes or so
	// Some apps might already be logging users for the new day
	blockDateBufferStr := time.Unix(lastEthBlockData.Timestamp + bufferTimeSec, 0).Format("20060102")
	blockDateBuffer64, _ := strconv.ParseInt(blockDateBufferStr, 10, 32)
	blockDateBuffer := int32(blockDateBuffer64)


	logger.Info("Block date ", blockDate)

	// Fetch the timestamp of the last eth block address
	for _, activity := range activities {
		activityAddress, _ := ActivityLogAddress(activity)

		activityData, err := s.context.GetState([]string{activityAddress})

		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err)}
		}

		if len(string(activityData[activityAddress])) == 0 {
			b, err := proto.Marshal(&activity)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: "could not marshal activity proto"}
			}

			if blockDate == activity.Date || blockDateBuffer == activity.Date || lastEthBlockData.GetTimestamp() == 0 {
				stateUpdate[activityAddress] = b
			} else {
				logger.Infof("Can't log activityDate=%v while current date is ", activity.Date, blockDate)
			}
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