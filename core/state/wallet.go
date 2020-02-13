package state

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/propsproject/props-transaction-processor/core/eth-utils"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"math/big"
)

func (s *State) UpdateWalletLinkData(walletToUserPayload *pending_props_pb.WalletToUser, updates map[string][]byte) (*pending_props_pb.ApplicationUser, []*pending_props_pb.ApplicationUser, error) {

	currentLinkedApplicationUsers := make([]*pending_props_pb.ApplicationUser, 0)
	newLinkedApplicationUsers := make([]*pending_props_pb.ApplicationUser, 0)
	walletToUserAddress, _ := WalletLinkAddress(walletToUserPayload.GetAddress())
	state, err := s.context.GetState([]string{walletToUserAddress})
	var unlinkedApplicationUser *pending_props_pb.ApplicationUser = nil

	if err != nil {
		return nil, nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err)}
	}

	var walletToUserData pending_props_pb.WalletToUser
	if len(string(state[walletToUserAddress])) == 0 {
		logger.Infof("First time this wallet is linked walletLinkAddress:%v, walletAddress:%v", walletToUserAddress, eth_utils.NormalizeAddress(walletToUserPayload.GetAddress()))
		walletToUserData = pending_props_pb.WalletToUser{
			Address: eth_utils.NormalizeAddress(walletToUserPayload.GetAddress()),
		}
	} else {
		// update existing wallet link
		for _, value := range state {
			err := proto.Unmarshal(value, &walletToUserData)
			if err != nil {
				return nil, nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal walletToUserData proto data (%s)", err)}
			}
		}
		currentLinkedApplicationUsers = walletToUserData.GetUsers()
	}
	for _, existingLinkedApplicationUser := range currentLinkedApplicationUsers {
		if existingLinkedApplicationUser.GetApplicationId() == walletToUserPayload.GetUsers()[0].GetApplicationId() {
			if existingLinkedApplicationUser.GetUserId() != walletToUserPayload.GetUsers()[0].GetUserId() {
				unlinkedApplicationUser = existingLinkedApplicationUser
			}
		} else {
			newLinkedApplicationUsers = append(newLinkedApplicationUsers, existingLinkedApplicationUser)
		}
	}

	newLinkedApplicationUsers = append(newLinkedApplicationUsers, walletToUserPayload.GetUsers()[0])
	walletToUserData.Users = newLinkedApplicationUsers
	walletToUserBytes, err := proto.Marshal(&walletToUserData)
	if err != nil {
		return nil, nil, &processor.InvalidTransactionError{Msg: "could not marshal wallet to user proto"}
	}
	updates[walletToUserAddress] = walletToUserBytes
	err2 := s.SaveWalletLinkEvents(&walletToUserData, walletToUserPayload.GetUsers()[0], unlinkedApplicationUser)
	if err2 != nil {
		return nil, nil, err2
	}
	return unlinkedApplicationUser, newLinkedApplicationUsers, nil

}

func (s *State) SaveUnlinkEvent(walletToUserData *pending_props_pb.WalletToUser, unlinkedApplicationUser *pending_props_pb.ApplicationUser) error {
	walletUnlinkedEvent := pending_props_pb.WalletUnlinkedEvent{
		User:          unlinkedApplicationUser,
		WalletToUsers: walletToUserData,
		Message:       fmt.Sprintf("wallet address %v unlinked from application user %v", walletToUserData.GetAddress(), unlinkedApplicationUser),
	}
	walletUnlinkAttr := []processor.Attribute{
		processor.Attribute{"address", walletToUserData.GetAddress()},
		processor.Attribute{"recipient", unlinkedApplicationUser.GetUserId()},
		processor.Attribute{"application", unlinkedApplicationUser.GetApplicationId()},
		processor.Attribute{"event_type", pending_props_pb.EventType_WalletUnlinked.String()},
		processor.Attribute{"signature", unlinkedApplicationUser.GetSignature()},
	}
	err := s.AddWalletUnlinkEvent(walletUnlinkedEvent, "pending-props:walletl", walletUnlinkAttr...)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) SaveWalletLinkEvents(walletToUserData *pending_props_pb.WalletToUser, newApplicationUser *pending_props_pb.ApplicationUser, unlinkedApplicationUser *pending_props_pb.ApplicationUser) error {
	if unlinkedApplicationUser != nil {
		err := s.SaveUnlinkEvent(walletToUserData, unlinkedApplicationUser)
		if err != nil {
			return err
		}
	}

	walletLinkedEvent := pending_props_pb.WalletLinkedEvent{
		User: newApplicationUser,
		WalletToUsers: walletToUserData,
		Message: fmt.Sprintf("wallet address %v linked to application user %v", walletToUserData.GetAddress(), newApplicationUser),
	}
	walletLinkAttr := []processor.Attribute{
		processor.Attribute{"address", walletToUserData.GetAddress()},
		processor.Attribute{"recipient",  newApplicationUser.GetUserId()},
		processor.Attribute{"application",  newApplicationUser.GetApplicationId()},
		processor.Attribute{"event_type", pending_props_pb.EventType_WalletLinked.String()},
		processor.Attribute{"signature", newApplicationUser.GetSignature()},
	}
	err := s.AddWalletLinkEvent(walletLinkedEvent, "pending-props:walletl", walletLinkAttr...)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) UpdateLinkedWalletBalances(walletToUserData *pending_props_pb.WalletToUser, unlinkedApplicationUser *pending_props_pb.ApplicationUser, linkedWalletApplicationUsers []*pending_props_pb.ApplicationUser, updates map[string][]byte, timestamp int64, updateType pending_props_pb.UpdateType, updatedUserBalance *pending_props_pb.Balance, updateWalletBalance *pending_props_pb.Balance) error {
	if unlinkedApplicationUser != nil {
		unlinkedApplicationUserBalanceAddress, unlinkedApplicationUserBalance, _, err := s.GetBalanceByApplicationUser(*unlinkedApplicationUser)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get unlinked application user balance data from address %v (%s)", unlinkedApplicationUserBalanceAddress, err)}
		}
		unlinkedApplicationUserBalance.LinkedWallet = ""
		unlinkedApplicationUserBalance.BalanceDetails.Transferable = big.NewInt(0).String()
		unlinkedApplicationUserBalance.BalanceDetails.Delegated = big.NewInt(0).String()
		unlinkedApplicationUserBalance.BalanceDetails.TotalPending = unlinkedApplicationUserBalance.GetBalanceDetails().GetPending()
		unlinkedApplicationUserBalance.BalanceUpdateIndex = unlinkedApplicationUserBalance.GetBalanceUpdateIndex() + 1
		err1 := s.UpdateBalance(*unlinkedApplicationUserBalance, updates, true)
		if err1 != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not update balance of unlinked application user %v (%s)", unlinkedApplicationUser.String(), err1)}
		}
	}
	// get wallet balance
	var newWalletBalance *pending_props_pb.Balance
	if updateWalletBalance == nil {
		walletBalanceAddress, walletBalance, newBalanceCreated, err := s.GetBalanceByApplicationUser(pending_props_pb.ApplicationUser{UserId: eth_utils.NormalizeAddress(walletToUserData.GetAddress()), ApplicationId: ""})
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get balance data from address %v (%s)", walletBalanceAddress, err)}
		}

		if newBalanceCreated {
			err := s.UpdateBalance(*walletBalance, updates, false)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save wallet balance data (%s)", err)}
			}
		}
		newWalletBalance = walletBalance
	} else {
		newWalletBalance = updateWalletBalance
	}

	var newTotalPending *big.Int
	if updatedUserBalance != nil {
		totalPending, ok := new(big.Int).SetString(updatedUserBalance.GetBalanceDetails().GetTotalPending(), 10)
		if !ok {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert updatedUserBalance.GetBalanceDetails().GetTotalPending() to big.Int (%s)",updatedUserBalance.GetBalanceDetails().GetTotalPending())}
		}
		newTotalPending = totalPending
	} else {
		totalPending, err := s.CalculateTotalPending(linkedWalletApplicationUsers)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could compute total pending of linked users (%s)",err)}
		}
		newTotalPending = totalPending
	}
	for _, applicationUser := range linkedWalletApplicationUsers {
		if updatedUserBalance != nil && updatedUserBalance.GetApplicationId() == applicationUser.GetApplicationId() && updatedUserBalance.GetUserId() == applicationUser.GetUserId() {
			updatedUserBalance.BalanceUpdateIndex = updatedUserBalance.GetBalanceUpdateIndex() + 1
			s.UpdateBalance(*updatedUserBalance, updates, true)
			continue
		}
		appUserBalanceAddress, appUserBalance, _, err := s.GetBalanceByApplicationUser(*applicationUser)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get application user balance data from address %v (%s)", appUserBalanceAddress, err)}
		}
		appUserBalance.BalanceDetails.TotalPending = newTotalPending.String()
		appUserBalance.LinkedWallet = eth_utils.NormalizeAddress(walletToUserData.GetAddress())
		appUserBalance.BalanceDetails.Transferable = newWalletBalance.GetBalanceDetails().GetTransferable()
		appUserBalance.BalanceDetails.Delegated = newWalletBalance.GetBalanceDetails().GetDelegated()
		appUserBalance.BalanceDetails.Timestamp = timestamp
		appUserBalance.BalanceDetails.LastUpdateType = updateType
		appUserBalance.BalanceUpdateIndex = appUserBalance.GetBalanceUpdateIndex() + 1
		s.UpdateBalance(*appUserBalance, updates, true)
	}
	return nil
}

func (s *State) CalculateTotalPending(linkedWalletApplicationUsers []*pending_props_pb.ApplicationUser) (*big.Int, error) {
	totalPending := big.NewInt(0)
	for _, applicationUser := range linkedWalletApplicationUsers {
		appUserBalanceAddress, appUserBalance, _, err := s.GetBalanceByApplicationUser(*applicationUser)
		if err != nil {
			return totalPending, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get application user balance data from address %v (%s)", appUserBalanceAddress, err)}
		}
		pending, ok := new(big.Int).SetString(appUserBalance.GetBalanceDetails().GetPending(), 10)
		if !ok {
			return totalPending, &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert appUserBalance.GetBalanceDetails().GetPending() to big.Int (%s)",appUserBalance.GetBalanceDetails().GetPending())}
		}
		totalPending = totalPending.Add(totalPending, pending)
	}
	return totalPending, nil
}

func (s *State) unlinkUserFromWallet(balanceUserAddress string, balanceUser *pending_props_pb.Balance, updates map[string][]byte) error {
	newLinkedApplicationUsers := make([]*pending_props_pb.ApplicationUser, 0)

	oldWalletLinkAddress, _ := WalletLinkAddress(balanceUser.GetLinkedWallet())
	state, err := s.context.GetState([]string{oldWalletLinkAddress})
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could get user old wallet link state data %v (%s)",oldWalletLinkAddress, err)}
	}
	if len(string(state[oldWalletLinkAddress])) > 0{
		var walletLinkData pending_props_pb.WalletToUser
		for _, value := range state {

			err := proto.Unmarshal(value, &walletLinkData)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal wallet link proto data (%s)", err)}
			}
		}
		currentLinkedApplicationUsers := walletLinkData.GetUsers()
		var unlinkedAppUser *pending_props_pb.ApplicationUser
		for _, existingLinkedApplicationUser := range currentLinkedApplicationUsers {
			if existingLinkedApplicationUser.GetApplicationId() != balanceUser.GetApplicationId() {
				newLinkedApplicationUsers = append(newLinkedApplicationUsers, existingLinkedApplicationUser)
			} else {
				unlinkedAppUser = existingLinkedApplicationUser
			}
		}
		walletLinkData.Users = newLinkedApplicationUsers
		walletToUserBytes, err := proto.Marshal(&walletLinkData)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal wallet link data to proto"}
		}
		updates[oldWalletLinkAddress] = walletToUserBytes
		err1 := s.SaveUnlinkEvent(&walletLinkData, unlinkedAppUser)
		if err1 != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could save unlink event user old wallet balance state data %v (%s)",unlinkedAppUser.String(), err1)}
		}
	}
	return nil
}

func (s *State) SaveWalletLink(walletToUsers ...pending_props_pb.WalletToUser) error {
	stateUpdate := make(map[string][]byte)
	for _, walletToUser := range walletToUsers {

		// verify signature
		if !eth_utils.VerifySig(walletToUser.GetAddress(), walletToUser.GetUsers()[0].GetSignature(), []byte(fmt.Sprintf("%v_%v", walletToUser.GetUsers()[0].GetApplicationId(), walletToUser.GetUsers()[0].GetUserId()))) {
			logger.Infof(fmt.Sprintf("Wallet verification %v, %v, %v, %v", walletToUser.GetAddress(), walletToUser.GetUsers()[0].GetSignature(), walletToUser.GetUsers()[0].GetApplicationId(), walletToUser.GetUsers()[0].GetUserId()))
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Signature verification fail %v, %v, %v", walletToUser.GetAddress(), walletToUser.GetUsers()[0].GetSignature(), []byte(fmt.Sprintf("%v_%v", walletToUser.GetUsers()[0].GetApplicationId(), walletToUser.GetUsers()[0].GetUserId())))}
		}
		//// check if user to be linked is
		balanceAddressUser, _ := BalanceAddressByAppUser( walletToUser.GetUsers()[0].GetApplicationId(), walletToUser.GetUsers()[0].GetUserId())
		state, err := s.context.GetState([]string{balanceAddressUser})
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could get linked user state data %v (%s)", walletToUser.String(), err)}
		}
		if len(string(state[balanceAddressUser])) > 0{
			var balanceUser pending_props_pb.Balance
			for _, value := range state {

				err := proto.Unmarshal(value, &balanceUser)
				if err != nil {
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal user balance proto data (%s)", err)}
				}
			}
			if len(balanceUser.GetLinkedWallet()) > 0 && balanceUser.GetLinkedWallet() != walletToUser.GetAddress() { // user is linked to a different wallet ==> unlink it
				err := s.unlinkUserFromWallet(balanceAddressUser, &balanceUser, stateUpdate)
				if err != nil {
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get unlink wallet balance (%s)", err)}
				}
			}
		}

		unlinkedApplicationUser, linkedApplicationUsers, err := s.UpdateWalletLinkData(&walletToUser, stateUpdate)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not update wallet link data %v (%s)", walletToUser.String(), err)}
		}
		err1 := s.UpdateLinkedWalletBalances(&walletToUser, unlinkedApplicationUser, linkedApplicationUsers, stateUpdate, walletToUser.GetUsers()[0].GetTimestamp(), pending_props_pb.UpdateType_WALLET_LINK_BALANCE, nil, nil)
		if err1 != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not update linked balances %v (%s)", walletToUser.String(), err1)}
		}
	}

	_, err := s.context.SetState(stateUpdate)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
	}
	return nil
}
