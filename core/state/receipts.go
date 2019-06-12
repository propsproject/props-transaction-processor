package state

import "math/big"

type NewEarningReceipt struct {
	Address     string `json:"address"`
	Recipient   string `json:"recipient"`
	Application string `json:"application"`
	Amount      big.Int `json:"amount"`
}

type EarningRevokedReceipt struct {
	Address     string `json:"address"`
	Recipient   string `json:"recipient"`
	Application string `json:"application"`
	Block       int    `json:"block"`
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

func GetEarningReceipt(address,  recipient, application string, amount big.Int) *NewEarningReceipt {
	return &NewEarningReceipt{address,  recipient, application, amount}
}

func GetBalanceUpdateReceipt(address,  recipient, application string, balance big.Int) *NewBalanceUpdateReceipt {
	return &NewBalanceUpdateReceipt{address,  recipient, application, balance}
}

func GetEarningRevokedReceipt(address,  recipient, application string, block int) *EarningRevokedReceipt {
	return &EarningRevokedReceipt{address, recipient,application, block}
}

func GetLastEthBlockUpdateReceipt(address string,  blockId int64) *LastEthBlockUpdateReceipt {
	return &LastEthBlockUpdateReceipt{address,  blockId}
}
