package state

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/propsproject/props-transaction-processor/core/eth-utils"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"math/big"
)


func (s *State) SaveSettlement(settlements ...pending_props_pb.SettlementData) error {
	stateUpdate := make(map[string][]byte)
	for _, settlementData := range settlements {
		settlementAddress, _ := SettlementAddress(settlementData.GetTxHash())
		existingSettlementStateData, err := s.context.GetState([]string{settlementAddress})
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Unable to get state of settlementAddress %v", settlementAddress)}
		}

		if len(string(existingSettlementStateData[settlementAddress])) == 0 {

			transaction := pending_props_pb.Transaction{
				Type:          pending_props_pb.Method_SETTLE,
				UserId:        settlementData.GetUserId(),
				ApplicationId: settlementData.GetApplicationId(),
				Timestamp:     settlementData.GetTimestamp(),
				Amount:        settlementData.GetAmount(),
				Description:   "Settlement",
				TxHash:        settlementData.GetTxHash(),
				Wallet:        eth_utils.NormalizeAddress(settlementData.GetToAddress()),
			}
			transactionAddress, _ := TransactionAddress(transaction)
			b, err := proto.Marshal(&transaction)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: "could not marshal transaction proto"}
			}
			stateUpdate[transactionAddress] = b
			settlementAmount, ok := new(big.Int).SetString(transaction.GetAmount(), 10)
			if !ok {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert settlement transaction.GetAmount() to big.Int (%s)",transaction.GetAmount())}
			}
			receiptBytes, err := json.Marshal(GetTransactionReceipt(transaction.GetType().String(), transactionAddress, transaction.GetUserId(), transaction.GetApplicationId(), *settlementAmount))
			if err != nil {
				logger.Infof("unable to create new transaction receipt (%s)", err)
			}

			err = s.context.AddReceiptData(receiptBytes)
			if err != nil {
				logger.Infof("unable to create new transaction receipt (%s)", err)
			}

			e := pending_props_pb.TransactionEvent{
				Transaction:  &transaction,
				Type:         transaction.GetType(),
				StateAddress: transactionAddress,
				Message:      fmt.Sprintf("transaction added: %s", transactionAddress),
				Description:  transaction.GetDescription(),
			}
			attr := []processor.Attribute{
				processor.Attribute{"recipient", transaction.GetUserId()},
				processor.Attribute{"application", transaction.GetApplicationId()},
				processor.Attribute{"event_type", pending_props_pb.EventType_TransactionAdded.String()},
				processor.Attribute{"transaction_type", transaction.GetType().String()},
				processor.Attribute{"description", transaction.GetDescription()},
			}
			s.AddEvent(e, "pending-props:transaction", attr...)

			settlementToBalance, ok := new(big.Int).SetString(settlementData.GetOnchainBalance(), 10)
			if !ok {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert settlement onchain balance to big.Int (%s)",settlementData.GetOnchainBalance())}
			}
			err1 := s.UpdateBalanceFromTransaction(transaction.GetUserId(), transaction.GetApplicationId(), *settlementAmount.Neg(settlementAmount), transaction.GetTimestamp(), stateUpdate, settlementToBalance)
			if err1 != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not prepare balance data %v", err1)}
			}
		} else {
			logger.Infof("This settlement was already submitted %v, %v", settlementData.GetTxHash(), settlementAddress)
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("This settlement was already submitted %v, %v", settlementData.GetTxHash(), settlementAddress)}
		}

		settlementDataBytes, err := proto.Marshal(&settlementData)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal settlementData proto"}
		}

		stateUpdate[settlementAddress] = settlementDataBytes
	}

	_, err := s.context.SetState(stateUpdate)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
	}

	return nil
}