package eth_utils

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/propsproject/goprops-toolkit/propstoken/bindings/token"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
)

type TransferDetails struct {
	TimeStamp         *big.Int					`json:"timestamp"`
	Address           common.Address			`json:"address"`
	Balance           *big.Int					`json:"balance"`
	Amount            *big.Int					`json:"amount"`
	From              common.Address			`json:"address"`
	To                common.Address			`json:"address"`
}

func GetEthTransactionTransferDetails(transactionHash string, address string, client *propstoken.Client, settlement bool) (*TransferDetails, uint64, error) {
	logging.Get().Infof("GetEthTransactionTransferDetails TransactionHash %s", transactionHash)
	transaction, err := client.RPC.TransactionReceipt(context.Background(), common.HexToHash(transactionHash))
	if err != nil {
		return nil, 0, fmt.Errorf("unable to get transaction receipt for hash (%s)", err)
	}

	for _, log := range transaction.Logs {
		topics := log.Topics
		sig := topics[0].Hex()
		if sig == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" { // only if matches Transfer signature keccak256(Transfer(address,address,uint256))
			from := fmt.Sprintf("0x%v", topics[1].Hex()[26:])
			to := fmt.Sprintf("0x%v", topics[2].Hex()[26:])
			if (address == from && !settlement) || address == to {

			} else {
				if (from == "0x0000000000000000000000000000000000000000" && to != address) || (to == "0x0000000000000000000000000000000000000000" && from != address) {
					continue
				}

				return nil, 0, fmt.Errorf("transaction %v address does not match reported balance update address %v, %v do not match %v", transactionHash, from, to, address)
			}

			callOptions := bind.CallOpts{
				Pending:     false,
				BlockNumber: new(big.Int).SetUint64(log.BlockNumber),
			}
			balance, err := client.Token.BalanceOf(&callOptions, common.HexToAddress(address))
			if err != nil {
				return nil, 0, fmt.Errorf("unable to get balanceOf %v on block %v (%s)", address, log.BlockNumber, err)
			} else {
				logging.Get().Infof("GetEthTransactionTransferDetails balance %v blockNumber %v (log.BlockNumber %v)", balance.String(), callOptions.BlockNumber.String(), log.BlockNumber)
			}

			transferDetails := TransferDetails{
				Address:        common.HexToAddress(address),
				Balance: balance,
				Amount: new(big.Int).SetBytes(log.Data),
				From: common.HexToAddress(from),
				To: common.HexToAddress(to),
			}
			return &transferDetails, log.BlockNumber, nil
		}
	}

	return nil, 0, fmt.Errorf("unable to get TransferDetails data from transaction (%s)", transactionHash)
}