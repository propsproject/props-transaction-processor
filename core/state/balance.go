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
	"strconv"
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
	if err == nil && len(string(existingTxStateData[updateBalanceTransactionAddress])) == 0 {
		logger.Infof(fmt.Sprintf("New updateBalanceTransactionAddress %v, (balanceUpdate for %v)", updateBalanceTransactionAddress, balanceUpdate.GetPublicAddress()))
		token, err := propstoken.NewPropsTokenHTTPClient(viper.GetString("props_token_contract_address"), viper.GetString("ethereum_url_tp"))
		if err != nil {
			logger.Infof("Could not connect to main-chain to verify balance update %v",err)
			token.RPC.Close()
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not connect to main-chain to verify balance update (%s) (balanceUpdate for %v)", err,  balanceUpdate.GetPublicAddress())}
		}
		latestHeader, err := token.RPC.HeaderByNumber(context.Background(), nil)
		if err != nil {
			logger.Infof("Could not get current blockId on main-chain to verify balance update %v (balanceUpdate for %v)",err, balanceUpdate.GetPublicAddress())
			token.RPC.Close()
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not get current blockId on main-chain to verify balance update (%s)", err)}
		}
		latestBlockId := latestHeader.Number
		if latestBlockId.Cmp(big.NewInt(0)) <= 0 {
			logger.Infof("Could not get current blockId on main-chain to verify balance update %v (balanceUpdate for %v)",err,balanceUpdate.GetPublicAddress())
			token.RPC.Close()
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not get current blockId on main-chain to verify balance update (%s)", err)}
		}
		logger.Infof("Latest Block on main-chain is %v (balanceUpdate for %v)", latestBlockId.String(), balanceUpdate.GetPublicAddress())
		confirmationBlocks := big.NewInt(viper.GetInt64("ethereum_confirmation_blocks"))
		if latestBlockId.Cmp(big.NewInt(0).Add(confirmationBlocks, big.NewInt(balanceUpdate.GetBlockId()))) >= 0 {
			// check details are correct looking up the transaction transfer details
			_transferDetails, transferBlockId, err := eth_utils.GetEthTransactionTransferDetails(eth_utils.NormalizeAddress(balanceUpdate.GetTxHash()), eth_utils.NormalizeAddress(balanceUpdate.GetPublicAddress()), token)
			token.RPC.Close()
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

				// it's all good we can save the data now

			} else {
				logger.Infof("Could not get TransferDetails from main-chain to verify balance update %v",err)
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could not get TransferDetails from main-chain to verify balance update (%s)", err)}
			}

		} else {
			token.RPC.Close()
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
		balanceWallet.BalanceDetails.Transferable = balanceUpdate.GetOnchainBalance()
		balanceWallet.BalanceDetails.Timestamp = newBalanceDetails.GetTimestamp()
		balanceWallet.BalanceDetails.LastEthBlockId = balanceUpdate.GetBlockId()
		balanceWallet.BalanceDetails.LastUpdateType = pending_props_pb.UpdateType_PROPS_BALANCE

		logger.Infof("Update recipient balance will be %v,%v", balanceWallet.BalanceDetails.Pending, balanceWallet.BalanceDetails.TotalPending)
	}

	// save balance wallet
	s.UpdateBalance(balanceWallet, updates, true)
	// check if it's linked to a wallet with more users and update them as needed
	applicationUsers := make([]*pending_props_pb.ApplicationUser, 0)
	walletLinkAddress, _ := WalletLinkAddress(balanceWallet.GetUserId())
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

	err2 := s.UpdateLinkedWalletBalances(&walletToUserData, nil, applicationUsers, updates, balanceUpdate.GetTimestamp(), pending_props_pb.UpdateType_PROPS_BALANCE, nil, &balanceWallet)
	if err2 != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save balances (%s)", err2)}
	}
	return nil
}

func (s *State) UpdateBalanceFromTransaction(userId, applicationId string, amount big.Int, timestamp int64, updates map[string][]byte, newBalanceAmount *big.Int) error {
	//1. get current balance
	//2. update current balance

	var newBalanceDetails pending_props_pb.BalanceDetails;
	if newBalanceAmount != nil {
		newBalanceDetails = pending_props_pb.BalanceDetails{
			Pending: amount.String(),
			TotalPending: amount.String(),
			Transferable: newBalanceAmount.String(),
			Bonded: big.NewInt(0).String(),
			Delegated: big.NewInt(0).String(),
			DelegatedTo: "",
			Timestamp: timestamp,
			LastUpdateType: pending_props_pb.UpdateType_PENDING_PROPS_BALANCE,
		}
	} else {
		newBalanceDetails = pending_props_pb.BalanceDetails{
			Pending: amount.String(),
			TotalPending: amount.String(),
			Transferable: big.NewInt(0).String(),
			Bonded: big.NewInt(0).String(),
			Delegated: big.NewInt(0).String(),
			DelegatedTo: "",
			Timestamp: timestamp,
			LastUpdateType: pending_props_pb.UpdateType_PENDING_PROPS_BALANCE,
		}
	}


	newBalanceUser := pending_props_pb.Balance{
		UserId: userId,
		ApplicationId: applicationId,
		BalanceDetails: &newBalanceDetails,
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
		if newBalanceAmount != nil {
			balanceUser.BalanceDetails.Transferable = newBalanceAmount.String()
		}
	}

	// check if it's linked to a wallet with more users
	applicationUsers := make([]*pending_props_pb.ApplicationUser, 0)
	if len(balanceUser.GetLinkedWallet())>0 {
		walletLinkAddress, _ := WalletLinkAddress(balanceUser.GetLinkedWallet())
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

			if newBalanceAmount != nil {
				walletBalance.BalanceDetails.Transferable = newBalanceAmount.String()
				s.UpdateBalance(*walletBalance, updates, true)
			}
		}
		err1 := s.UpdateLinkedWalletBalances(&walletToUserData, nil, applicationUsers, updates, timestamp, pending_props_pb.UpdateType_PENDING_PROPS_BALANCE, &balanceUser, nil)
		if err1 != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Failed to update linked balances (%v)", walletToUserData.String())}
		}
	} else {
		balanceUser.BalanceUpdateIndex = balanceUser.GetBalanceUpdateIndex() + 1
		err := s.UpdateBalance(balanceUser, updates, true)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save balance %v (%s)", balanceUser.String(), err)}
		}
	}
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
			processor.Attribute{ "balance_update_index", strconv.FormatInt(balance.GetBalanceUpdateIndex(), 10)},
		}
		s.AddBalanceEvent(balanceEvent, "pending-props:balance", balanceUpdateAttr...)
	}

	// if there's an activity object for this day update it with new balance - (unless it's a balance of a wallet)
	if balance.Type == pending_props_pb.BalanceType_USER {
		rewardsDay := eth_utils.CalculateRewardsDay(balance.GetBalanceDetails().GetTimestamp())
		activityLog := pending_props_pb.ActivityLog{
			UserId:        balance.GetUserId(),
			ApplicationId: balance.GetApplicationId(),
			Date:          int32(rewardsDay),
		}
		activityAddress, _ := ActivityLogAddress(activityLog)
		state, err := s.context.GetState([]string{activityAddress})
		if err != nil {
			logger.Infof("Could not get state data %v rewardsDay=%v, timestamp=%v, balance.userId=%v, balance.applicationId=%v (%s)", activityAddress, rewardsDay, balance.GetBalanceDetails().GetTimestamp(), balance.GetUserId(), balance.GetApplicationId(), err)
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state data %v (%s)", activityAddress, err)}
		}
		if len(string(state[activityAddress])) > 0 {
			for _, value := range state {

				err := proto.Unmarshal(value, &activityLog)
				if err != nil {
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal activity log proto data (%s)", err)}
				}
			}
			logger.Infof("******* activity.Date=%v, rewardsDay=%v", activityLog.GetDate(), int32(rewardsDay))
			if activityLog.GetDate() == int32(rewardsDay) { // update activity balance object only if user has activity the day of the update
				activityLog.Balance = &balance
				activityLogBytes, err := proto.Marshal(&activityLog)
				if err != nil {
					return &processor.InvalidTransactionError{Msg: "could not marshal activityLog update to activityLog proto"}
				}
				updates[activityAddress] = activityLogBytes
			}

		}
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
		logger.Infof("SaveBalanceUpdate: Received balance update for %v, %v", balanceUpdate.GetPublicAddress(), balanceUpdate.GetOnchainBalance())
		err := s.UpdateBalanceFromMainchainEvent(balanceUpdate, stateUpdate)

		if err != nil {
			errMsg := err.Error()
			if strings.Index(errMsg,"TransactionHashAlreadyExists") >= 0 {
				logger.Infof("SaveBalanceUpdate: TransactionHashAlreadyExists: %v", err)
				return nil
			} else
			{
				logger.Infof("SaveBalanceUpdate: Invalid transaction: %v", err)
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