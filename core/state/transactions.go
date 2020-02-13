package state

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"math/big"
	"strings"
)

func (s *State) SaveTransactions(transactions ...pending_props_pb.Transaction) error {
	stateUpdate := make(map[string][]byte)
	for _, transaction := range transactions {
		transactionAddress, _ := TransactionAddress(transaction)
		b, err := proto.Marshal(&transaction)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal transaction proto"}
		}
		// settle transaction not allowed - only via external balance updates
		if transaction.GetType() == pending_props_pb.Method_SETTLE {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Illegal operation - settlement transactions happen via external balance updates (%s)", len(b))}
		}

		//stateUpdate[transactionAddress] = b
		var finalAmount *big.Int
		amount, ok := new(big.Int).SetString(transaction.GetAmount(), 10)
		if !ok {
			// check if got a decimal here and round it
			if strings.Index(transaction.GetAmount(),".")>=0 {
				amountSplit := strings.Split(transaction.GetAmount(), ".")
				amount2, ok2 := new(big.Int).SetString(amountSplit[0], 10)
				if !ok2 {
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert transaction.GetAmount() to big.Int (%s) after rounding (%s)",transaction.GetAmount(),amountSplit[0])}
				} else {
					logger.Infof("Successfully rounded (%s) to (%s)", transaction.GetAmount(), amountSplit[0])
					finalAmount = amount2
				}

			} else {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert transaction.GetAmount() to big.Int (%s)",transaction.GetAmount())}
			}
		} else {
			finalAmount = amount
		}

		receiptBytes, err := json.Marshal(GetTransactionReceipt(transaction.GetType().String(), transactionAddress, transaction.GetUserId(), transaction.GetApplicationId(), *finalAmount))
		if err != nil {
			logger.Infof("unable to create new transaction receipt (%s)", err)
		}

		err = s.context.AddReceiptData(receiptBytes)
		if err != nil {
			logger.Infof("unable to create new transaction receipt (%s)", err)
		}

		e := pending_props_pb.TransactionEvent{
			Transaction: &transaction,
			Type: transaction.GetType(),
			StateAddress: transactionAddress,
			Message: fmt.Sprintf("transaction added: %s", transactionAddress),
			Description: transaction.GetDescription(),
		}
		attr := []processor.Attribute{
			processor.Attribute{"recipient", transaction.GetUserId()},
			processor.Attribute{"application", transaction.GetApplicationId()},
			processor.Attribute{"event_type", pending_props_pb.EventType_TransactionAdded.String()},
			processor.Attribute{"transaction_type", transaction.GetType().String()},
			processor.Attribute{"description", transaction.GetDescription()},
		}
		s.AddEvent(e, "pending-props:transaction", attr...)

		if transaction.GetType() == pending_props_pb.Method_REVOKE {
			// this is either revoke or settle which means balance should decrease
			finalAmount = finalAmount.Neg(finalAmount)
		}
		err1 := s.UpdateBalanceFromTransaction(transaction.GetUserId(), transaction.GetApplicationId(), *finalAmount, transaction.GetTimestamp(), stateUpdate, nil)
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