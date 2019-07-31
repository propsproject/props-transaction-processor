package state

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/propsproject/goprops-toolkit/propstoken/bindings/token"
	"github.com/propsproject/props-transaction-processor/core/eth-utils"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/spf13/viper"
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
			// check settlement data on ethereum
			logger.Infof(fmt.Sprintf("New settlement for application %v, user %v, wallet %v, amount %v", settlementData.GetApplicationId(), settlementData.GetUserId(), settlementData.GetToAddress(), settlementData.GetAmount()))
			token, err := propstoken.NewPropsTokenHTTPClient(viper.GetString("props_token_contract_address"), viper.GetString("ethereum_url"))
			if err != nil {
				logger.Infof("Could not connect to main-chain to verify settlement %v",err)
				token.RPC.Close()
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not connect to main-chain to verify settlement (%s)", err)}
			}
			latestHeader, err := token.RPC.HeaderByNumber(context.Background(), nil)
			if err != nil {
				logger.Infof("Could not get current blockId on main-chain to verify settlement %v",err)
				token.RPC.Close()
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not get current blockId on main-chain to verify settlement (%s)", err)}
			}
			latestBlockId := latestHeader.Number
			if latestBlockId.Cmp(big.NewInt(0)) <= 0 {
				logger.Infof("Could not get current blockId on main-chain to verify settlement %v",err)
				token.RPC.Close()
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not get current blockId on main-chain to verify settlement (%s)", err)}
			}
			logger.Infof("Latest Block on main-chain is %v", latestBlockId.String())
			confirmationBlocks := big.NewInt(viper.GetInt64("ethereum_confirmation_blocks"))
			if latestBlockId.Cmp(big.NewInt(0).Add(confirmationBlocks, big.NewInt(settlementData.GetBlockId()))) >= 0 {
				// check details are correct looking up the transaction transfer details
				_settlementDetails, settlementBlockId, err := eth_utils.GetEthTransactionSettlementDetails(eth_utils.NormalizeAddress(settlementData.GetTxHash()), token)
				token.RPC.Close()
				if err != nil {
					logger.Infof("Could verify settlement on main-chain %v",err)
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could verify settlement on main-chain %v",err)}
				}
				if settlementBlockId == uint64(settlementData.GetBlockId()) &&
					_settlementDetails.Amount.String() == settlementData.GetAmount() &&
					eth_utils.NormalizeAddress(_settlementDetails.To.String()) == settlementData.GetToAddress() &&
					eth_utils.NormalizeAddress(_settlementDetails.ApplicationId.String()) == settlementData.GetApplicationId() &&
					_settlementDetails.UserId == settlementData.GetUserId() {

					transaction := pending_props_pb.Transaction{
						Type: pending_props_pb.Method_SETTLE,
						UserId: settlementData.GetUserId(),
						ApplicationId: settlementData.GetApplicationId(),
						Timestamp: settlementData.GetTimestamp(),
						Amount: _settlementDetails.Amount.String(),
						Description: "Settlement",
						TxHash: settlementData.GetTxHash(),
						Wallet: eth_utils.NormalizeAddress(_settlementDetails.To.String()),
					}
					transactionAddress, _ := TransactionAddress(transaction)
					b, err := proto.Marshal(&transaction)
					if err != nil {
						return &processor.InvalidTransactionError{Msg: "could not marshal transaction proto"}
					}
					stateUpdate[transactionAddress] = b
					receiptBytes, err := json.Marshal(GetTransactionReceipt(transaction.GetType().String(), transactionAddress, transaction.GetUserId(), transaction.GetApplicationId(), *_settlementDetails.Amount))
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
					err1 := s.UpdateBalanceFromTransaction(transaction.GetUserId(), transaction.GetApplicationId(), *_settlementDetails.Amount.Neg(_settlementDetails.Amount), transaction.GetTimestamp(), stateUpdate, _settlementDetails.Balance)
					if err1 != nil {
						return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not prepare balance data %v", err1)}
					}
				} else {
					logger.Infof("This settlement (%v, %v) could not be verified %v, %v, %v, %v, %v != %v, %v, %v, %v, %v",
						settlementData.GetTxHash(), settlementAddress, settlementBlockId, _settlementDetails.Amount.String(), _settlementDetails.To.String(), _settlementDetails.ApplicationId.String(), _settlementDetails.UserId,
						uint64(settlementData.GetBlockId()), settlementData.GetAmount(), settlementData.GetToAddress(), settlementData.GetApplicationId(), settlementData.GetUserId())
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("This settlement (%v, %v) could not be verified %v, %v, %v, %v, %v != %v, %v, %v, %v, %v",
						settlementData.GetTxHash(), settlementAddress, settlementBlockId, _settlementDetails.Amount.String(), _settlementDetails.To.String(), _settlementDetails.ApplicationId.String(), _settlementDetails.UserId,
						uint64(settlementData.GetBlockId()), settlementData.GetAmount(), settlementData.GetToAddress(), settlementData.GetApplicationId(), settlementData.GetUserId())}
				}
			} else {
				token.RPC.Close()
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Not enough confirmation blocks latestBlockId=%v submittedBlockId=%v", latestBlockId.String(), settlementData.GetBlockId())}
			}
		} else {
			logger.Infof("This settlement was already submitted %v, %v", settlementData.GetTxHash(), settlementAddress)
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("This settlement was already submitted %v, %v", settlementData.GetTxHash(), settlementAddress)}
		}
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not blockUpdateData (%v)", err)}
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