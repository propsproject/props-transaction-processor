package eth_utils

import (
	"context"
	"encoding/hex"
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
	From              common.Address			`json:"from"`
	To                common.Address			`json:"to"`
}

type SettlementDetails struct {
	TimeStamp         *big.Int					`json:"timestamp"`
	Amount            *big.Int					`json:"amount"`
	From              common.Address			`json:"from"`
	To                common.Address			`json:"to"`
	ApplicationId     common.Address			`json:"applicationId"`
	UserId            string					`json:"userId"`
}

type SettlementEvent struct {
	From              common.Address
	UserId            string
	To                common.Address
	Amount            *big.Int
	RewardsAddress    common.Address
}

func GetEthTransactionTransferDetails(transactionHash string, address string, client *propstoken.Client) (*TransferDetails, uint64, error) {
	logging.Get().Infof("GetEthTransactionTransferDetails TransactionHash %s ((balanceUpdate for %v)", transactionHash, address)
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
			logging.Get().Infof("GetEthTransactionTransferDetails checking addresses from=0x%v, to=0x%v address=%v", from, to, address)
			if address == from || address == to {

			} else {
				if (from == "0x0000000000000000000000000000000000000000" && to != address) || (to == "0x0000000000000000000000000000000000000000" && from != address) {
					logging.Get().Infof("GetEthTransactionTransferDetails skipping due to from/to address being 0 (balanceUpdate for %v)", address)
					continue
				}
				logging.Get().Infof("transaction %v address does not match reported balance update address %v, %v do not match %v", transactionHash, from, to, address)
				return nil, 0, fmt.Errorf("transaction %v address does not match reported balance update address %v, %v do not match %v", transactionHash, from, to, address)
			}

			callOptions := bind.CallOpts{
				Pending:     false,
				BlockNumber: new(big.Int).SetUint64(log.BlockNumber),
			}
			balance, err := client.Token.BalanceOf(&callOptions, common.HexToAddress(address))
			if err != nil {
				logging.Get().Infof("unable to get balanceOf %v on block %v (%s)", address, log.BlockNumber, err)
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
			logging.Get().Infof("GetEthTransactionTransferDetails got transfer details amount=%v (balanceUpdate for %v)", transferDetails.Amount.String(), address)
			return &transferDetails, log.BlockNumber, nil
		} else {
			logging.Get().Infof("GetEthTransactionTransferDetails signature didn't match transfer got sig=%v (balanceUpdate for %v)", sig, address)
		}
	}
	logging.Get().Infof("unable to get TransferDetails data from transaction (%s) (balanceUpdate for %v)", transactionHash, address)
	return nil, 0, fmt.Errorf("unable to get TransferDetails data from transaction (%s)", transactionHash)
}


func GetEthTransactionSettlementDetails(transactionHash string, client *propstoken.Client) (*SettlementDetails, uint64, error) {
	logging.Get().Infof("GetEthTransactionSettlementDetails TransactionHash %s", transactionHash)
	transaction, err := client.RPC.TransactionReceipt(context.Background(), common.HexToHash(transactionHash))
	if err != nil {
		return nil, 0, fmt.Errorf("unable to get transaction receipt for hash (%s)", err)
	}

	for _, log := range transaction.Logs {
		topics := log.Topics
		sig := topics[0].Hex()
		logging.Get().Infof("GetEthTransactionSettlementDetails sig=%v", sig)
		if sig == "53b5073ff19aef23b167e83c6be14817da210375bec35b4c0ccfc0cded9a23f8" { // only if matches Transfer signature keccak256(Settlement(address,bytes32,address,uint256,address))
			var settlementEvent SettlementEvent
			applicationId := common.HexToAddress(topics[1].Hex())
			userId, err := hex.DecodeString(topics[2].Hex())
			if err != nil {
				return nil, 0, fmt.Errorf("unable to decode userId %v (%s)", topics[2].Hex(), err)
			}
			to := common.HexToAddress(topics[3].Hex())
			err1 := client.ABI.Unpack(&settlementEvent, "Settlement", log.Data)
			if err1 != nil {
				return nil, 0, fmt.Errorf("unable to unpack log.Data %v (%s)", log.Data, err1)
			}

			settlementDetails := SettlementDetails{
				Amount:			settlementEvent.Amount,
				From:			settlementEvent.RewardsAddress,
				To:				to,
				ApplicationId: 	applicationId,
				UserId:			string(userId),

			}

			logging.Get().Infof("GetEthTransactionSettlementDetails got %v ", settlementDetails)
			return &settlementDetails, log.BlockNumber, nil
		} else {
			logging.Get().Infof("GetEthTransactionSettlementDetails signature didn't match settlement got sig=%v", sig)
		}
	}
	logging.Get().Infof("unable to get SettlementDetails data from transaction (%s)", transactionHash)
	return nil, 0, fmt.Errorf("unable to get TransferDetails data from transaction (%s)", transactionHash)
}