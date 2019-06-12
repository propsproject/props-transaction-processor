package state

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/propsproject/pending-props/core/eth-utils"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/sawtooth-go-sdk/logging"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"math/big"
)

func (s *State) GetEarnings(address ...string) ([]pending_props_pb.Earning, error) {
	logger.Infof("addresses %s", address,)
	state, err := s.context.GetState(address)
	if err != nil {
		return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err)}
	}

	earnings := make([]pending_props_pb.Earning, 0)
	for _, value := range state {
		var earning pending_props_pb.Earning
		err := proto.Unmarshal(value, &earning)
		if err != nil {
			return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto data (%s)", err)}
		}
		earnings = append(earnings, earning)
	}

	return earnings, nil
}

func (s *State) SavePendingEarnings(earnings ...pending_props_pb.Earning) error {
	stateUpdate := make(map[string][]byte)
	for _, earning := range earnings {
		pendingAddr, settledAddr, revokedAddr := EarningAddress(earning)
		b, err := proto.Marshal(&earning)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal earning proto"}
		}

		stateUpdate[pendingAddr] = b
		earned, ok := new(big.Int).SetString(earning.GetDetails().GetAmountEarned(), 10)
		if !ok {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert earning.GetDetails().GetAmountEarned() to big.Int (%s)",earning.GetDetails().GetAmountEarned())}
		}

		receiptBytes, err := json.Marshal(GetEarningReceipt(pendingAddr, earning.GetDetails().GetUserId(), earning.GetDetails().GetApplicationId(), *earned))
		if err != nil {
			logging.Get().Infof("unable to create new earning receipt (%s)", err)
		}

		err = s.context.AddReceiptData(receiptBytes)
		if err != nil {
			logging.Get().Infof("unable to create new earning receipt (%s)", err)
		}

		e := pending_props_pb.EarningEvent{
			Earning: &earning,
			IssueAddress: pendingAddr,
			RevokeAddress: revokedAddr,
			SettleAddress: settledAddr,
			Message: fmt.Sprintf("earning issued: %s", pendingAddr),
			Description: earning.GetDetails().GetDescription(),
		}
		attr := []processor.Attribute{
			processor.Attribute{"recipient", earning.GetDetails().GetUserId()},
			processor.Attribute{"application", earning.GetDetails().GetApplicationId()},
			processor.Attribute{"event_type", pending_props_pb.EventType_EarningIssued.String()},
			processor.Attribute{"description", earning.GetDetails().GetDescription()},
		}
		s.AddEvent(e, "pending-props:earnings", attr...)


		err1 := s.UpdateBalanceFromEarningsChange(earning.GetDetails().GetUserId(), earning.GetDetails().GetApplicationId(), *earned, earning.GetDetails().GetTimestamp(), stateUpdate)
		if err1 != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not prepare balance data %v", err1)}
		}
	}


	_, err := s.context.SetState(stateUpdate)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
	}

	return nil
}

func (s *State) SaveRevokedEarnings(timestamp int64, earnings ...pending_props_pb.Earning) error {
	stateUpdate := make(map[string][]byte)
	pendingAddresses := make([]string, 0)

	for _, earning := range earnings {
		pendingAddr, settledAddr, revokedAddr := EarningAddress(earning)

		b, err := proto.Marshal(&earning)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal earning proto"}
		}

		stateUpdate[revokedAddr] = b
		pendingAddresses = append(pendingAddresses, pendingAddr)
		earned, ok := new(big.Int).SetString(earning.GetDetails().GetAmountEarned(), 10)
		if !ok {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert earning.GetDetails().GetAmountEarned() to big.Int (%s)",earning.GetDetails().GetAmountEarned())}
		}
		receiptBytes, err := json.Marshal(GetEarningReceipt(revokedAddr, earning.GetDetails().GetUserId(), earning.GetDetails().GetApplicationId(), *earned))
		if err != nil {
			logging.Get().Infof("unable to create new earning receipt (%s)", err)
		}

		err = s.context.AddReceiptData(receiptBytes)
		if err != nil {
			logging.Get().Infof("unable to create new earning receipt (%s)", err)
		}

		e := pending_props_pb.EarningEvent{
			Earning: &earning,
			IssueAddress: pendingAddr,
			RevokeAddress: revokedAddr,
			SettleAddress: settledAddr,
			Message: fmt.Sprintf("earning revoked: %s", revokedAddr),
			Description: earning.GetDetails().GetDescription(),
		}
		attr := []processor.Attribute{
			processor.Attribute{"recipient", earning.GetDetails().GetUserId()},
			processor.Attribute{"application", earning.GetDetails().GetApplicationId()},
			processor.Attribute{"event_type", pending_props_pb.EventType_EarningRevoked.String()},
			processor.Attribute{"description", earning.GetDetails().GetDescription()},
		}
		s.AddEvent(e, "pending-props:earnings", attr...)
		earned = earned.Neg(earned)
		err1 := s.UpdateBalanceFromEarningsChange(earning.GetDetails().GetUserId(), earning.GetDetails().GetApplicationId(), *earned, timestamp, stateUpdate)
		if err1 != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not prepare balance data (%s)", err1)}
		}
	}

	_, err := s.context.SetState(stateUpdate)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
	}

	_, err = s.context.DeleteState(pendingAddresses)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not delete pending earnings from state (%s)", err)}
	}

	return nil
}

func (s *State) SettleEarnings(timestamp int64, ethTxtID string, earnings ...pending_props_pb.Earning) error {
	stateUpdate := make(map[string][]byte)
	pendingAddrs := make([]string, 0)
	ethTxtID = eth_utils.NormalizeAddress(ethTxtID)
	settlements := new(pending_props_pb.Settlements)
	address, _ := SettlementAddress(ethTxtID)

	//@Todo BUG FIX not saving partially settled earnings, accumalate in separate slice and save
	for _, earning := range earnings {
		if earning.GetSettledByTransaction() != ethTxtID {
			continue
		}
		pendingAddr, settledAddr, revokedAddr := EarningAddress(earning)
		b, err := proto.Marshal(&earning)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal earning proto"}
		}

		stateUpdate[settledAddr] = b
		pendingAddrs = append(pendingAddrs, pendingAddr)

		settlements.EarningAddresses = append(settlements.EarningAddresses, settledAddr)

		e := pending_props_pb.EarningEvent{
			Earning: &earning,
			IssueAddress: pendingAddr,
			RevokeAddress: revokedAddr,
			SettleAddress: settledAddr,
			Message: fmt.Sprintf("earning settled: %s", settledAddr),
			Description: earning.GetDetails().GetDescription(),
		}
		attr := []processor.Attribute{
			processor.Attribute{"recipient", eth_utils.NormalizeAddress(earning.GetDetails().GetUserId())},
			processor.Attribute{"application", eth_utils.NormalizeAddress(earning.GetDetails().GetApplicationId())},
			processor.Attribute{"event_type", pending_props_pb.EventType_EarningSettled.String()},
			processor.Attribute{"description", earning.GetDetails().GetDescription()},

		}
		s.AddEvent(e, "pending-props:earnings", attr...)
		//receiptBytes, err := json.Marshal(GetEarningReceipt(address, earning.GetDetails().GetRecipientPublicAddress(),  earning.GetDetails().GetApplicationPublicAddress(),  earning.GetDetails().GetAmountEarned()))
		//if err != nil {
		//	logging.Get().Infof("unable to create new earning receipt (%s)", err)
		//}
		//
		//err = s.context.AddReceiptData(receiptBytes)
		//if err != nil {d
		//	logging.Get().Infof("unable to create new earning receipt (%s)", err)
		//}
		settled, ok := new(big.Int).SetString(earning.GetDetails().GetAmountSettled(), 10)
		if !ok {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert earning.GetDetails().GetAmountSettled() to big.Int (%s)",earning.GetDetails().GetAmountSettled())}
		}
		settled = settled.Neg(settled)
		err1 := s.UpdateBalanceFromEarningsChange(earning.GetDetails().GetUserId(), earning.GetDetails().GetApplicationId(), *settled, timestamp, stateUpdate)
		if err1 != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not prepare balance and balanceTimestamp data %s", err1)}
		}
	}

	b, err := proto.Marshal(settlements)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: "could not marshal settlements proto"}
	}

	stateUpdate[address] = b

	_, err = s.context.SetState(stateUpdate)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
	}

	_, err = s.context.DeleteState(pendingAddrs)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not remove earnings from pending state (%s)", err)}
	}

	return nil
}