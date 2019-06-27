package state

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/propsproject/goprops-toolkit/propstoken/bindings/token"
	"github.com/propsproject/props-transaction-processor/core/eth-utils"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/spf13/viper"
	"math/big"
	"strings"
)

func (s *State) UpdateBalanceFromMainchainEvent(balanceUpdate pending_props_pb.BalanceUpdate, updates map[string][]byte) error {
	//1. Check if this transaction was already accounted for -> return success and do nothing
	//2. Check that enough blocks passed -> return error
	//3. Check that balanceOf at that block is indeed what was submitted -> return error
	//4. Store balanceUpdateTransaction
	//5. get current balance
	//6. update current balance

	updateBalanceTransactionAddress, _ := BalanceUpdatesTransactionHashAddress(eth_utils.NormalizeAddress(balanceUpdate.GetTxHash()), balanceUpdate.GetPublicAddress())
	existingTxStateData, err := s.context.GetState([]string{updateBalanceTransactionAddress})
	var settledApplicationUser *pending_props_pb.ApplicationUser
	var settledAmount *big.Int
	if err == nil && len(string(existingTxStateData[updateBalanceTransactionAddress])) == 0 {
		logger.Infof(fmt.Sprintf("Error / Not Found while getting state updateBalanceTransactionAddress %v, %v", updateBalanceTransactionAddress, err))
		token, err := propstoken.NewPropsTokenHTTPClient(viper.GetString("props_token_contract_address"), viper.GetString("ethereum_url"))
		if err != nil {
			logger.Infof("Could not connect to main-chain to verify balance update %v",err)
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not connect to main-chain to verify balance update (%s)", err)}
		}
		latestHeader, err := token.RPC.HeaderByNumber(context.Background(), nil)
		if err != nil {
			logger.Infof("Could not get current blockId on main-chain to verify balance update %v",err)
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not get current blockId on main-chain to verify balance update (%s)", err)}
		}
		latestBlockId := latestHeader.Number
		if latestBlockId.Cmp(big.NewInt(0)) <= 0 {
			logger.Infof("Could not get current blockId on main-chain to verify balance update %v",err)
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not get current blockId on main-chain to verify balance update (%s)", err)}
		}
		logger.Infof("Latest Block on main-chain is %v", latestBlockId.String())
		confirmationBlocks := big.NewInt(viper.GetInt64("ethereum_confirmation_blocks"))
		if latestBlockId.Cmp(big.NewInt(0).Add(confirmationBlocks, big.NewInt(balanceUpdate.GetBlockId()))) >= 0 {
			logger.Infof("Enough blocks (%v) passed submitted=%v,current=%v", confirmationBlocks, balanceUpdate.GetBlockId(), latestBlockId)
			// check details are correct looking up the transaction transfer details
			_transferDetails, transferBlockId, err := eth_utils.GetEthTransactionTransferDetails(eth_utils.NormalizeAddress(balanceUpdate.GetTxHash()), eth_utils.NormalizeAddress(balanceUpdate.GetPublicAddress()), token, false)
			if err == nil && transferBlockId > 0 {
				tdAddress := eth_utils.NormalizeAddress(_transferDetails.Address.String())
				tdBalance := _transferDetails.Balance
				buAddress := eth_utils.NormalizeAddress(balanceUpdate.GetPublicAddress())
				buBalance, ok := new(big.Int).SetString(balanceUpdate.GetOnchainBalance(), 10)
				if !ok {
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert balanceUpdate.GetFromOnchainBalance() to big.Int (%s)", balanceUpdate.GetOnchainBalance())}
				}
				buBlockId := uint64(balanceUpdate.BlockId)

				if tdAddress != buAddress ||
					tdBalance.Cmp(buBalance)!=0 ||
					transferBlockId != buBlockId {
					tdBytes, _ := json.Marshal(_transferDetails)
					logger.Infof("TransferDetails (%v) are different than submitted data (%v) - transferBlockId=%v", string(tdBytes) , balanceUpdate, transferBlockId)
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("TransferDetails (%v) are different than submitted data (%v)", _transferDetails, balanceUpdate)}
				}

				// is the transfer from a settlement address?
				settlementFromAddresses := viper.GetStringMapString("settlement_from_addresses")
				if transferFromAppId, ok := settlementFromAddresses[eth_utils.NormalizeAddress(_transferDetails.From.String())]; ok {
					// does the transfer to address linked to any user?
					walletLinkAddress, _ := WalletLinkAddress(pending_props_pb.WalletToUser{ Address: eth_utils.NormalizeAddress(_transferDetails.To.String())})
					state1, err := s.context.GetState([]string{walletLinkAddress})
					var walletToUserData pending_props_pb.WalletToUser
					if err != nil {
						return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get wallet link data %v (%s)", walletLinkAddress, err)}
					}
					if len(string(state1[walletLinkAddress])) > 0 {
						for _, value := range state1 {
							err := proto.Unmarshal(value, &walletToUserData)
							if err != nil {
								return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal wallet link proto data (%s)", err)}
							}
						}
						applicationUsers := walletToUserData.GetUsers()
						for _, applicationUser := range applicationUsers {
							// element is the element from someSlice for where we are
							if applicationUser.GetApplicationId() == transferFromAppId {
								settledApplicationUser = applicationUser
								settledAmount = _transferDetails.Amount
								transaction := pending_props_pb.Transaction{
									Type: pending_props_pb.Method_SETTLE,
									UserId: applicationUser.GetUserId(),
									ApplicationId: applicationUser.GetApplicationId(),
									Timestamp: balanceUpdate.GetTimestamp(),
									Amount: settledAmount.String(),
									Description: "Settlement",
									TxHash: balanceUpdate.GetTxHash(),
									Wallet: eth_utils.NormalizeAddress(_transferDetails.To.String()),
								}
								transactionAddress, _ := TransactionAddress(transaction)
								b, err := proto.Marshal(&transaction)
								if err != nil {
									return &processor.InvalidTransactionError{Msg: "could not marshal transaction proto"}
								}
								updates[transactionAddress] = b
								receiptBytes, err := json.Marshal(GetTransactionReceipt(transaction.GetType().String(), transactionAddress, transaction.GetUserId(), transaction.GetApplicationId(), *settledAmount))
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

							}
						}
					}

				}

				// it's all good we can save the data now

			} else {
				logger.Infof("Could not get TransferDetails from main-chain to verify balance update %v",err)
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not get TransferDetails from main-chain to verify balance update (%s)", err)}
			}

		} else {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Not enough confirmation blocks latestBlockId=%v submittedBlockId=%v", latestBlockId.String(), balanceUpdate.GetBlockId())}
		}
	} else {
		logger.Infof("This balance update was already submitted %v", updateBalanceTransactionAddress)
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("TransactionHashAlreadyExists %v", updateBalanceTransactionAddress)}
	}

	// if we got here there were no errors and we can return data to be saved
	newBalanceDetails := pending_props_pb.BalanceDetails{
		Pending: big.NewInt(0).String(),
		TotalPending: big.NewInt(0).String(),
		Transferable: balanceUpdate.GetOnchainBalance(),
		Bonded: big.NewInt(0).String(),
		Delegated: big.NewInt(0).String(),
		DelegatedTo: "",
		Timestamp: balanceUpdate.GetTimestamp(),
		LastEthBlockId: balanceUpdate.GetBlockId(),
		LastUpdateType: pending_props_pb.UpdateType_PROPS_BALANCE,
	}

	newBalanceWallet := pending_props_pb.Balance{
		UserId:                 eth_utils.NormalizeAddress(balanceUpdate.GetPublicAddress()),
		BalanceDetails:         &newBalanceDetails,
		PreCutoffDetails:       &newBalanceDetails,
		Type:          pending_props_pb.BalanceType_WALLET,
	}

	balanceAddressWallet, _ := BalanceAddress(newBalanceWallet)
	state, err := s.context.GetState([]string{balanceAddressWallet})
	var balanceWallet pending_props_pb.Balance
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state data %v (%s)", balanceAddressWallet, err)}
	}
	if len(string(state[balanceAddressWallet])) == 0 {
		logger.Infof("Error / Not Found while getting state %v recipient address %v, %v", balanceAddressWallet, eth_utils.NormalizeAddress(balanceUpdate.GetPublicAddress()), err)
		balanceWallet = newBalanceWallet
	} else {
		// update existing balance
		for _, value := range state {

			err := proto.Unmarshal(value, &balanceWallet)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto data (%s)", err)}
			}
		}
		// if the new balance timestamp day is bigger than current one, store current in PreCutOff one
		if balanceWallet.GetPreCutoffDetails().GetTimestamp() == 0 || CalculateRewardsDay(newBalanceDetails.GetTimestamp()) > CalculateRewardsDay(balanceWallet.GetBalanceDetails().GetTimestamp()) {
			logger.Infof("Updating PreCutoff New Balance Timestamp %v, Current Balance Timestamp %v, PreCutoff Balance Timestamp %v",
				newBalanceDetails.GetTimestamp(), balanceWallet.GetBalanceDetails().GetTimestamp(), balanceWallet.GetPreCutoffDetails().GetTimestamp())
			//balanceWallet.PreCutoffDetails = balanceWallet.GetBalanceDetails()
			balanceWallet.PreCutoffDetails = &pending_props_pb.BalanceDetails{
				Pending:  balanceWallet.GetBalanceDetails().GetPending(),
				TotalPending:  balanceWallet.GetBalanceDetails().GetTotalPending(),
				Transferable:  balanceWallet.GetBalanceDetails().GetTransferable(),
				Bonded: balanceWallet.GetBalanceDetails().GetBonded(),
				Delegated: balanceWallet.GetBalanceDetails().GetDelegated(),
				DelegatedTo: balanceWallet.GetBalanceDetails().GetDelegatedTo(),
				Timestamp: balanceWallet.GetBalanceDetails().GetTimestamp(),
				LastUpdateType: balanceWallet.GetBalanceDetails().GetLastUpdateType(),
				LastEthBlockId: balanceWallet.GetBalanceDetails().GetLastEthBlockId(),
			}
		}
		balanceWallet.BalanceDetails.Transferable = balanceUpdate.GetOnchainBalance()
		if settledAmount != nil {
			totalPending, ok := new(big.Int).SetString(balanceWallet.GetBalanceDetails().GetTotalPending(), 10)
			if !ok {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert balanceWallet.GetBalanceDetails().GetTotalPending() to big.Int (%s)", balanceWallet.GetBalanceDetails().GetTotalPending())}
			}
			balanceWallet.BalanceDetails.TotalPending = totalPending.Sub(totalPending, settledAmount).String()
		}
		balanceWallet.BalanceDetails.Timestamp = newBalanceDetails.GetTimestamp()
		balanceWallet.BalanceDetails.LastEthBlockId = balanceUpdate.GetBlockId()
		balanceWallet.BalanceDetails.LastUpdateType = pending_props_pb.UpdateType_PROPS_BALANCE

		logger.Infof("Update recipient balance will be %v,%v", balanceWallet.BalanceDetails.Pending, balanceWallet.BalanceDetails.TotalPending)
	}


	logger.Infof("Current Balance Timestamp %v, PreCutoff Balance Timestamp %v",
		newBalanceDetails.GetTimestamp(), balanceWallet.GetBalanceDetails().GetTimestamp(), balanceWallet.GetPreCutoffDetails().GetTimestamp())
	// save balance wallet
	s.UpdateBalance(balanceWallet, updates, true)
	// check if it's linked to a wallet with more users and update them as needed
	applicationUsers := make([]*pending_props_pb.ApplicationUser, 0)
	walletLinkAddress, _ := WalletLinkAddress(pending_props_pb.WalletToUser{ Address: balanceWallet.GetUserId()})
	state1, err1 := s.context.GetState([]string{walletLinkAddress})
	var walletToUserData pending_props_pb.WalletToUser
	if err1 != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state data %v (%s)", walletLinkAddress, err1)}
	}
	if len(string(state1[walletLinkAddress])) == 0 {

		logger.Infof("Error / Not Found while getting linked wallet data %v from state - it is not linked %v", walletLinkAddress)
		return nil
	} else {
		for _, value := range state1 {
			err := proto.Unmarshal(value, &walletToUserData)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto data (%s)", err)}
			}
		}
		applicationUsers = walletToUserData.GetUsers()
	}

	err2 := s.UpdateLinkedWalletBalances(applicationUsers, balanceWallet, true, updates, settledApplicationUser, settledAmount)
	if err2 != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save balances (%s)", err2)}
	}
	return nil
}

func (s *State) UpdateBalanceFromTransaction(userId, applicationId string, amount big.Int, timestamp int64, updates map[string][]byte) error {
	//1. get current balance
	//2. update current balance

	newBalanceDetails := pending_props_pb.BalanceDetails{
		Pending: amount.String(),
		TotalPending: amount.String(),
		Transferable: big.NewInt(0).String(),
		Bonded: big.NewInt(0).String(),
		Delegated: big.NewInt(0).String(),
		DelegatedTo: "",
		Timestamp: timestamp,
		LastUpdateType: pending_props_pb.UpdateType_PENDING_PROPS_BALANCE,
	}

	newBalanceUser := pending_props_pb.Balance{
		UserId: userId,
		ApplicationId: applicationId,
		BalanceDetails: &newBalanceDetails,
		PreCutoffDetails: &newBalanceDetails,
		Type: pending_props_pb.BalanceType_USER,
	}

	balanceAddressUser, _ := BalanceAddress(newBalanceUser)
	logger.Infof("BalanceAddress = %v", balanceAddressUser)
	state, err := s.context.GetState([]string{balanceAddressUser})
	var balanceUser pending_props_pb.Balance
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state data %v (%s)", balanceAddressUser, err)}
	}
	if len(string(state[balanceAddressUser])) == 0{
		// how to differentiate between error and not found?
		// assume error caused by not found
		logger.Infof("Error / Not Found while getting previous balance from state %v", err)
		balanceUser = newBalanceUser
	} else {
		// update existing balance
		for _, value := range state {

			err := proto.Unmarshal(value, &balanceUser)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto data (%s)", err)}
			}
		}

		if balanceUser.GetPreCutoffDetails().GetTimestamp() == 0 || CalculateRewardsDay(newBalanceDetails.GetTimestamp()) > CalculateRewardsDay(balanceUser.GetBalanceDetails().GetTimestamp()) {
			logger.Infof("Updating PreCutoff New Balance Timestamp %v, Current Balance Timestamp %v, PreCutoff Balance Timestamp %v",
				newBalanceDetails.GetTimestamp(), balanceUser.GetBalanceDetails().GetTimestamp(), balanceUser.GetPreCutoffDetails().GetTimestamp())
			//balanceUser.PreCutoffDetails = balanceUser.GetBalanceDetails()
			balanceUser.PreCutoffDetails = &pending_props_pb.BalanceDetails{
				Pending:  balanceUser.GetBalanceDetails().GetPending(),
				TotalPending:  balanceUser.GetBalanceDetails().GetTotalPending(),
				Transferable:  balanceUser.GetBalanceDetails().GetTransferable(),
				Bonded: balanceUser.GetBalanceDetails().GetBonded(),
				Delegated: balanceUser.GetBalanceDetails().GetDelegated(),
				DelegatedTo: balanceUser.GetBalanceDetails().GetDelegatedTo(),
				Timestamp: balanceUser.GetBalanceDetails().GetTimestamp(),
				LastUpdateType: balanceUser.GetBalanceDetails().GetLastUpdateType(),
				LastEthBlockId: balanceUser.GetBalanceDetails().GetLastEthBlockId(),
			}
		}

		pending, ok := new(big.Int).SetString(balanceUser.GetBalanceDetails().GetPending(), 10)
		if !ok {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert balanceUser.GetBalanceDetails().GetPending() to big.Int (%s)", balanceUser.GetBalanceDetails().GetPending())}
		}

		totalPending, ok := new(big.Int).SetString(balanceUser.GetBalanceDetails().GetTotalPending(), 10)
		if !ok {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert balanceUser.GetBalanceDetails().GetTotalPending() to big.Int (%s)", balanceUser.GetBalanceDetails().GetTotalPending())}
		}

		newPending, ok := new(big.Int).SetString(newBalanceDetails.GetPending(), 10)
		if !ok {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert newBalanceDetails.GetPending() to big.Int (%s)", newBalanceDetails.GetPending())}
		}

		newTotalPending, ok := new(big.Int).SetString(newBalanceDetails.GetTotalPending(), 10)
		if !ok {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert newBalanceDetails.GetTotalPending() to big.Int (%s)", newBalanceDetails.GetTotalPending())}
		}

		balanceUser.BalanceDetails.Pending = pending.Add(pending, newPending).String()
		balanceUser.BalanceDetails.TotalPending = totalPending.Add(totalPending, newTotalPending).String()
		balanceUser.BalanceDetails.Timestamp = newBalanceDetails.GetTimestamp()
	}

	// check if it's linked to a wallet with more users
	applicationUsers := make([]*pending_props_pb.ApplicationUser, 0)
	if len(balanceUser.GetLinkedWallet())>0 {
		walletLinkAddress, _ := WalletLinkAddress(pending_props_pb.WalletToUser{ Address: balanceUser.GetLinkedWallet()})
		state, err := s.context.GetState([]string{walletLinkAddress})
		var walletToUserData pending_props_pb.WalletToUser
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state data %v (%s)", walletLinkAddress, err)}
		}
		if len(string(state[walletLinkAddress])) == 0 {

			logger.Infof("Error / Not Found while getting previous linked wallet data %v from state %v", walletLinkAddress, err)
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Wallet %v is marked as link but no such object exists", walletLinkAddress)}
		} else {
			for _, value := range state {
				err := proto.Unmarshal(value, &walletToUserData)
				if err != nil {
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto data (%s)", err)}
				}
			}
			applicationUsers = walletToUserData.GetUsers()
			walletBalanceAddress, walletBalance, newBalanceCreated, err1 := s.GetBalanceByApplicationUser(pending_props_pb.ApplicationUser{UserId:eth_utils.NormalizeAddress(balanceUser.GetLinkedWallet()), ApplicationId:""})
			if err1 != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get linked wallet balance object (%s)", err1)}
			}
			if newBalanceCreated {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("if wallet is linked walletBalance object must exist at (%v)", walletBalanceAddress)}
			}

			if walletBalance.GetPreCutoffDetails().GetTimestamp() == 0 || CalculateRewardsDay(newBalanceDetails.GetTimestamp()) > CalculateRewardsDay(walletBalance.GetBalanceDetails().GetTimestamp()) {
				logger.Infof("Updating walletBalance PreCutoff New Balance Timestamp %v, Current Balance Timestamp %v, PreCutoff Balance Timestamp %v",
					newBalanceDetails.GetTimestamp(), walletBalance.GetBalanceDetails().GetTimestamp(), walletBalance.GetPreCutoffDetails().GetTimestamp())
				//walletBalance.PreCutoffDetails = walletBalance.GetBalanceDetails()
				walletBalance.PreCutoffDetails = &pending_props_pb.BalanceDetails{
					Pending:  balanceUser.GetPreCutoffDetails().GetPending(),
					TotalPending:  balanceUser.GetPreCutoffDetails().GetTotalPending(),
					Transferable:  balanceUser.GetPreCutoffDetails().GetTransferable(),
					Bonded: balanceUser.GetPreCutoffDetails().GetBonded(),
					Delegated: balanceUser.GetPreCutoffDetails().GetDelegated(),
					DelegatedTo: balanceUser.GetPreCutoffDetails().GetDelegatedTo(),
					Timestamp: balanceUser.GetPreCutoffDetails().GetTimestamp(),
					LastUpdateType: balanceUser.GetPreCutoffDetails().GetLastUpdateType(),
					LastEthBlockId: balanceUser.GetPreCutoffDetails().GetLastEthBlockId(),
				}
			}

			walletBalance.BalanceDetails.TotalPending = balanceUser.GetBalanceDetails().GetTotalPending()
			walletBalance.BalanceDetails.LastUpdateType = balanceUser.GetBalanceDetails().GetLastUpdateType()
			walletBalance.BalanceDetails.Timestamp = balanceUser.GetBalanceDetails().GetTimestamp()
			s.UpdateBalance(*walletBalance, updates, true)
		}
	} else {
		applicationUsers = append(applicationUsers, &pending_props_pb.ApplicationUser{ ApplicationId: balanceUser.GetApplicationId(), UserId: balanceUser.GetUserId()})
	}
	err1 := s.UpdateLinkedWalletBalances(applicationUsers, balanceUser, false, updates, nil, nil)
	if err1 != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save balances (%s)", err1)}
	}
	logger.Infof("Update balance will be %v,%v", balanceUser.BalanceDetails.GetPending(), balanceUser.BalanceDetails.GetTotalPending())

	return nil
}

func (s *State) UpdateBalance(balance pending_props_pb.Balance, updates map[string][]byte, sendEvent bool) error {
	balanceAddress, _ := BalanceAddress(balance)
	balanceBytes, err := proto.Marshal(&balance)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: "could not marshal balance to balance proto"}
	}

	if sendEvent {
		balanceEvent := pending_props_pb.BalanceEvent{
			Balance: &balance,
			Message: fmt.Sprintf("balance updated: %s", balance),
		}
		balanceUpdateAttr := []processor.Attribute{
			processor.Attribute{"recipient", balance.GetUserId()},
			processor.Attribute{"application", balance.GetApplicationId()},
			processor.Attribute{"event_type", pending_props_pb.EventType_BalanceUpdated.String()},
			processor.Attribute{"balance_type", balance.GetType().String()},
		}
		s.AddBalanceEvent(balanceEvent, "pending-props:balance", balanceUpdateAttr...)
	}

	updates[balanceAddress] = balanceBytes
	return nil
}

func (s * State) GetBalanceByApplicationUser(applicationUser pending_props_pb.ApplicationUser) (string, *pending_props_pb.Balance, bool, error) {
	var newBalanceCreated bool = false
	balanceAddress, _ := BalanceAddressByAppUser(applicationUser.GetApplicationId(), applicationUser.GetUserId())
	state, err := s.context.GetState([]string{balanceAddress})
	var balance pending_props_pb.Balance
	if err != nil {
		return balanceAddress, nil, newBalanceCreated, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state data %v (%s)", balanceAddress, err)}
	}
	if len(string(state[balanceAddress])) == 0 {
		// balance does not exist yet
		newBalanceCreated = true
		logger.Infof(fmt.Sprintf("Error / Not Found while getting balance address=%v,applicationId=%v,userId=%v,err=%v", balanceAddress, applicationUser.GetApplicationId(), applicationUser.GetUserId(), err))
		var balanceType = pending_props_pb.BalanceType_USER
		if applicationUser.GetApplicationId() == "" {
			balanceType = pending_props_pb.BalanceType_WALLET
		}
		balance = pending_props_pb.Balance{
			UserId:         applicationUser.GetUserId(),
			ApplicationId:  applicationUser.GetApplicationId(),
			BalanceDetails: &pending_props_pb.BalanceDetails{
				Pending: big.NewInt(0).String(),
				TotalPending: big.NewInt(0).String(),
				Transferable: big.NewInt(0).String(),
				Bonded: big.NewInt(0).String(),
				Delegated: big.NewInt(0).String(),
				DelegatedTo: "",
			},
			Type: balanceType,
		}
	} else {
		for _, value := range state {

			err := proto.Unmarshal(value, &balance)
			if err != nil {
				return balanceAddress, nil, newBalanceCreated, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal balance proto data (%s)", err)}
			}
		}
	}
	return balanceAddress, &balance, newBalanceCreated, nil
}

func (s *State) SaveBalanceUpdate(balanceUpdates ...pending_props_pb.BalanceUpdate) error {
	stateUpdate := make(map[string][]byte)
	for _, balanceUpdate := range balanceUpdates {
		err := s.UpdateBalanceFromMainchainEvent(balanceUpdate, stateUpdate)

		if err != nil {
			errMsg := err.Error()
			if strings.Index(errMsg,"TransactionHashAlreadyExists") >= 0 {
				return nil
			} else
			{
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Error verifying transaction  (%v)", err)}
			}
		}

		updateBalanceTransactionAddress, _ := BalanceUpdatesTransactionHashAddress(eth_utils.NormalizeAddress(balanceUpdate.GetTxHash()), balanceUpdate.GetPublicAddress())
		balanceUpdateBytes, err := proto.Marshal(&balanceUpdate)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal balance update proto"}
		}
		stateUpdate[updateBalanceTransactionAddress] = balanceUpdateBytes
	}


	_, err := s.context.SetState(stateUpdate)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
	}

	return nil
}

func CalculateRewardsDay(timestamp int64) int64 {
	// (block.timestamp.sub(_self.rewardsStartTimestamp)).div(_self.minSecondsBetweenDays).add(1);
	secondsInDay := viper.GetInt64("seconds_in_day")
	return timestamp / secondsInDay

}