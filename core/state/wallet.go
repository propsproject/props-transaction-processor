package state

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/propsproject/pending-props/core/eth-utils"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"math/big"
)

func (s *State) SaveWalletLink(walletToUsers ...pending_props_pb.WalletToUser) error {
	stateUpdate := make(map[string][]byte)
	for _, walletToUser := range walletToUsers {

		// verify signature
		if !eth_utils.VerifySig(walletToUser.GetAddress(), walletToUser.GetUsers()[0].GetSignature(), []byte(fmt.Sprintf("%v_%v", walletToUser.GetUsers()[0].GetApplicationId(), walletToUser.GetUsers()[0].GetUserId()))) {
			logger.Infof(fmt.Sprintf("Wallet verification %v, %v, %v, %v", walletToUser.GetAddress(), walletToUser.GetUsers()[0].GetSignature(), walletToUser.GetUsers()[0].GetApplicationId(),  walletToUser.GetUsers()[0].GetUserId()))
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Signature verification fail %v, %v, %v", walletToUser.GetAddress(), walletToUser.GetUsers()[0].GetSignature(), []byte(fmt.Sprintf("%v_%v", walletToUser.GetUsers()[0].GetApplicationId(), walletToUser.GetUsers()[0].GetUserId())))}
		}

		// get wallet balance
		walletBalanceAddress, walletBalance, newBalanceCreated, err := s.GetBalanceByApplicationUser(pending_props_pb.ApplicationUser{UserId:eth_utils.NormalizeAddress(walletToUser.GetAddress()), ApplicationId:""})
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get balance data from address %v (%s)", walletBalanceAddress, err)}
		}
		if newBalanceCreated {
			err := s.UpdateBalance(*walletBalance, stateUpdate, true)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save balance data (%s)", err)}
			}
		}

		currentLinkedApplicationUsers := make([]*pending_props_pb.ApplicationUser, 0)
		walletToUserAddress, _ := WalletLinkAddress(walletToUser)
		state1, err1 := s.context.GetState([]string{walletToUserAddress})

		if err1 != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err1)}
		}

		var walletToUserData pending_props_pb.WalletToUser
		if len(string(state1[walletToUserAddress])) == 0 {
			logger.Infof("Error / Not Found while getting state %v, %v, %v", walletToUserAddress, eth_utils.NormalizeAddress(walletToUser.GetAddress()), err1)
			walletToUserData = pending_props_pb.WalletToUser{
				Address: eth_utils.NormalizeAddress(walletToUser.GetAddress()),
			}
		} else {
			// update existing wallet link
			for _, value := range state1 {

				err := proto.Unmarshal(value, &walletToUserData)
				if err != nil {
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal walletToUserData proto data (%s)", err)}
				}
			}
		}

		newApplicationUserBalanceAddress, newApplicationUserBalance, newBalanceCreated, err := s.GetBalanceByApplicationUser(pending_props_pb.ApplicationUser{UserId:walletToUser.GetUsers()[0].GetUserId(), ApplicationId:walletToUser.GetUsers()[0].GetApplicationId()})
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get balance data from address %v (%s)", newApplicationUserBalanceAddress, err)}
		}
		if newBalanceCreated {
			newApplicationUserBalance.BalanceDetails = walletBalance.GetBalanceDetails()
			newApplicationUserBalance.LinkedWallet = eth_utils.NormalizeAddress(walletToUser.GetAddress())
		}

		currentLinkedApplicationUsers = append(currentLinkedApplicationUsers, walletToUser.GetUsers()[0])
		var unlinkedBalance pending_props_pb.Balance
		var existingApplicationUserWasUnlinked = false
		for _, existingLinkedApplicationUser := range walletToUserData.Users {
			if existingLinkedApplicationUser.GetApplicationId() == walletToUser.GetUsers()[0].GetApplicationId() {
				// unlink wallet from previous application user
				unlinkedBalanceAddress, _ := BalanceAddress(pending_props_pb.Balance{ UserId: existingLinkedApplicationUser.GetUserId(), ApplicationId: existingLinkedApplicationUser.GetApplicationId()})
				state, err := s.context.GetState([]string{unlinkedBalanceAddress})

				if err != nil || len(string(state[unlinkedBalanceAddress])) == 0 {
					logger.Infof(fmt.Sprintf("Error / Not Found while getting state balance address=%v, applicationId=%v, userId=%v (%s)", unlinkedBalanceAddress, existingLinkedApplicationUser.GetApplicationId(), existingLinkedApplicationUser.GetUserId(), err1))
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Existing application user (%v, %v) to unlink must have an existing balance object %v (%s)", existingLinkedApplicationUser.GetApplicationId(), existingLinkedApplicationUser.GetUserId(), unlinkedBalanceAddress, err)}
				} else {
					for _, value := range state {

						err := proto.Unmarshal(value, &unlinkedBalance)
						if err != nil {
							return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal unlinkedBalance proto data (%s)", err)}
						}
					}
					existingApplicationUserWasUnlinked = true
					unlinkedBalance.LinkedWallet = ""
					unlinkedBalance.BalanceDetails.Transferable = big.NewInt(0).String()
					unlinkedBalance.BalanceDetails.Delegated = big.NewInt(0).String()
					unlinkedBalance.BalanceDetails.TotalPending = unlinkedBalance.GetBalanceDetails().GetPending()

					s.UpdateBalance(
						unlinkedBalance,
						stateUpdate,
						existingLinkedApplicationUser.GetUserId() != walletToUser.GetUsers()[0].GetUserId())

					walletUnlinkedEvent := pending_props_pb.WalletUnlinkedEvent{
						User: existingLinkedApplicationUser,
						WalletToUsers: &walletToUserData,
						Message: fmt.Sprintf("wallet address %v unlinked from application user %v", walletToUser.GetAddress(),  existingLinkedApplicationUser),
					}
					walletUnlinkAttr := []processor.Attribute{
						processor.Attribute{"address", walletToUser.GetAddress()},
						processor.Attribute{"recipient",  existingLinkedApplicationUser.GetUserId()},
						processor.Attribute{"application",  existingLinkedApplicationUser.GetApplicationId()},
						processor.Attribute{"event_type", pending_props_pb.EventType_WalletUnlinked.String()},
						processor.Attribute{"signature", existingLinkedApplicationUser.GetSignature()},
					}
					s.AddWalletUnlinkEvent(walletUnlinkedEvent, "pending-props:walletl", walletUnlinkAttr...)
				}
			} else {
				currentLinkedApplicationUsers = append(currentLinkedApplicationUsers, existingLinkedApplicationUser)
			}
		}
		walletToUserData.Users = currentLinkedApplicationUsers

		if len(walletToUserData.Users) == 1 { // there's only one which is the new one
			newApplicationUserBalance.LinkedWallet = walletBalance.GetUserId()
			newApplicationUserBalance.BalanceDetails.Transferable = walletBalance.GetBalanceDetails().GetTransferable()
			newApplicationUserBalance.BalanceDetails.Delegated = walletBalance.GetBalanceDetails().GetDelegated()
			err := s.UpdateBalance(*newApplicationUserBalance, stateUpdate, true)
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save balance data (%s)", err)}
			}
			//walletBalance.BalanceDetails.TotalPending = newApplicationUserBalance.BalanceDetails.GetTotalPending()
			//err2 := s.UpdateBalance(*walletBalance, stateUpdate)
			//if err2 != nil {
			//	return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save updated wallet balance data (%s)", err)}
			//}
		} else {
			// There's 1 or more existing applicationUserBalance objects
			// Get the first one, add the new pending and update all of the application users to have the same total and totalPending
			existingUserBalanceAddress, existingUserBalance, newBalanceCreated, err := s.GetBalanceByApplicationUser(pending_props_pb.ApplicationUser{UserId:walletToUserData.GetUsers()[1].GetUserId(), ApplicationId:walletToUserData.GetUsers()[1].GetApplicationId()})
			if err != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get balance data from address %v (%s)", existingUserBalanceAddress, err)}
			}
			if newBalanceCreated {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("There must be a record for existingUserBalance before linking it to external wallet (%s)", err)}
			}

			existingTotalPending, ok := new(big.Int).SetString(existingUserBalance.GetBalanceDetails().GetTotalPending(), 10)
			if !ok {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert existingUserBalance.GetBalanceDetails().GetTotalPending() to big.Int (%s)",existingUserBalance.GetBalanceDetails().GetTotalPending())}
			}

			newApplicationUserPending, ok := new(big.Int).SetString(newApplicationUserBalance.GetBalanceDetails().GetPending(), 10)
			if !ok {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert newApplicationUserBalance.GetBalanceDetails().GetPending() to big.Int (%s)",newApplicationUserBalance.GetBalanceDetails().GetPending())}
			}

			if existingApplicationUserWasUnlinked {
				pendingToBeRemoved, ok := new(big.Int).SetString(unlinkedBalance.GetBalanceDetails().GetPending(), 10)
				if !ok {
					return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Could convert unlinkedBalance.GetBalanceDetails().GetPending() to big.Int (%s)",unlinkedBalance.GetBalanceDetails().GetPending())}
				}
				newTotalPending := new(big.Int).Add(existingTotalPending, newApplicationUserPending)
				newTotalPendingMinusUnlinked := new(big.Int).Sub(newTotalPending, pendingToBeRemoved)
				newApplicationUserBalance.BalanceDetails.TotalPending = newTotalPendingMinusUnlinked.String()
			} else {
				newTotalPending := new(big.Int).Add(existingTotalPending, newApplicationUserPending)
				newApplicationUserBalance.BalanceDetails.TotalPending = newTotalPending.String()
			}


			newApplicationUserBalance.BalanceDetails.Timestamp = walletToUser.GetUsers()[0].Timestamp
			newApplicationUserBalance.LinkedWallet = walletBalance.GetUserId()
			err1 := s.UpdateLinkedWalletBalances(walletToUserData.Users, *newApplicationUserBalance, false, stateUpdate)
			if err1 != nil {
				return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save balances of linked wallet data (%s)", err1)}
			}
			//walletBalance.BalanceDetails.TotalPending = newApplicationUserBalance.BalanceDetails.GetTotalPending()
			//err2 := s.UpdateBalance(*walletBalance, stateUpdate)
			//if err2 != nil {
			//	return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save updated wallet balance data (%s)", err)}
			//}
		}

		walletLinkedEvent := pending_props_pb.WalletLinkedEvent{
			User: walletToUser.GetUsers()[0],
			WalletToUsers: &walletToUserData,
			Message: fmt.Sprintf("wallet address %v linked to application user %v", walletToUser.GetAddress(),  walletToUser.GetUsers()[0]),
		}
		walletLinkAttr := []processor.Attribute{
			processor.Attribute{"address", walletToUser.GetAddress()},
			processor.Attribute{"recipient",  walletToUser.GetUsers()[0].GetUserId()},
			processor.Attribute{"application",  walletToUser.GetUsers()[0].GetApplicationId()},
			processor.Attribute{"event_type", pending_props_pb.EventType_WalletLinked.String()},
			processor.Attribute{"signature", walletToUser.GetUsers()[0].GetSignature()},
		}
		s.AddWalletLinkEvent(walletLinkedEvent, "pending-props:walletl", walletLinkAttr...)

		walletToUserBytes, err := proto.Marshal(&walletToUserData)
		if err != nil {
			return &processor.InvalidTransactionError{Msg: "could not marshal wallet to user proto"}
		}
		stateUpdate[walletToUserAddress] = walletToUserBytes
	}
	_, err := s.context.SetState(stateUpdate)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not set state (%s)", err)}
	}

	return nil
}

func (s *State) UpdateLinkedWalletBalances(applicationUsers []*pending_props_pb.ApplicationUser, balance pending_props_pb.Balance, onchainOnly bool, updates map[string][]byte) error {
	for _, applicationUser := range applicationUsers {
		applicationUserBalanceAddress, applicationUserBalance, newBalanceCreated, err := s.GetBalanceByApplicationUser(pending_props_pb.ApplicationUser{UserId:applicationUser.GetUserId(), ApplicationId:applicationUser.GetApplicationId()})
		if err != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get balance data from address %v (%s)", applicationUserBalanceAddress, err)}
		}

		if newBalanceCreated && !(applicationUser.GetUserId() == balance.GetUserId() && applicationUser.GetApplicationId() == balance.GetApplicationId()) {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("There must be a record for applicationUserBalance before linking it to external wallet (%s)", err)}
		} else if onchainOnly {
			applicationUserBalance.BalanceDetails.Delegated = balance.GetBalanceDetails().GetDelegated()
			applicationUserBalance.BalanceDetails.Timestamp = balance.GetBalanceDetails().GetTimestamp()
			applicationUserBalance.BalanceDetails.Transferable = balance.GetBalanceDetails().GetTransferable()
		} else {
			if applicationUser.GetUserId() == balance.GetUserId() && applicationUser.GetApplicationId() == balance.GetApplicationId() {
				applicationUserBalance = &balance
			}
			applicationUserBalance.BalanceDetails.TotalPending = balance.BalanceDetails.GetTotalPending()
			applicationUserBalance.BalanceDetails.Timestamp = balance.GetBalanceDetails().GetTimestamp()
		}

		err1 := s.UpdateBalance(*applicationUserBalance, updates, true)

		if err1 != nil {
			return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not save balance data (%s)", err1)}
		}
	}
	return nil
}