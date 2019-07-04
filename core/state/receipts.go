package state

import (
	"math/big"
)

type NewTransactionReceipt struct {
	TransactionType 	string	`json:"type"`
	Address     		string	`json:"address"`
	Recipient   		string	`json:"recipient"`
	Application 		string 	`json:"application"`
	Amount      		big.Int	`json:"amount"`
}

type NewBalanceUpdateReceipt struct {
	Address     string `json:"address"`
	Recipient   string `json:"recipient"`
	Application string `json:"application"`
	Balance     big.Int `json:"balance"`
}

type LastEthBlockUpdateReceipt struct {
	Address     string `json:"address"`
	BlockId     int64  `json:"blockId"`
}

type RewardEntityUpdateReceipt struct {
	Name               string `json:"name"`
	Address            string `json:"address"`
	RewardsAddress     string `json:"rewardsAddress"`
	SidechainAddress   string `json:"sidechainAddress"`
}

func GetTransactionReceipt(transactionType, address,  recipient, application string, amount big.Int) *NewTransactionReceipt {
	return &NewTransactionReceipt{transactionType, address,  recipient, application, amount}
}

func GetBalanceUpdateReceipt(address,  recipient, application string, balance big.Int) *NewBalanceUpdateReceipt {
	return &NewBalanceUpdateReceipt{address,  recipient, application, balance}
}

func GetLastEthBlockUpdateReceipt(address string,  blockId int64) *LastEthBlockUpdateReceipt {
	return &LastEthBlockUpdateReceipt{address,  blockId}
}

func GetRewardEntityUpdateReceipt(name, address, rewardsAddress, sidechainAddress string) *RewardEntityUpdateReceipt {
	return &RewardEntityUpdateReceipt{name, address, rewardsAddress, sidechainAddress }
}
