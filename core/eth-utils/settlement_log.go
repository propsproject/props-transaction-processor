package eth_utils

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/propsproject/goprops-toolkit/propstoken/bindings/token"
	"github.com/propsproject/sawtooth-go-sdk/logging"
	"github.com/spf13/viper"
)

const SettlementEvent = "Settlement"

type SettlementData struct {
	TimeStamp *big.Int
	From      common.Address
	Recipient common.Address
	Amount    *big.Int
}

func GetEthTransactionSettlementData(transactionHash string) (*SettlementData, error) {
	client, err := propstoken.NewPropsTokenHTTPClient(viper.GetString("props_token_contract_address"), viper.GetString("ethereum_url"))
	if err != nil {
		return nil, fmt.Errorf("could not instantiate token client (%v)", err)
	}
	logging.Get().Infof("%s", transactionHash)
	transaction, err := client.RPC.TransactionReceipt(context.Background(), common.HexToHash(transactionHash))
	if err != nil {
		return nil, fmt.Errorf("unable to get transaction receipt for hash (%s)", err)
	}

	for _, log := range transaction.Logs {
		var settlement SettlementData
		err := client.ABI.Unpack(&settlement, SettlementEvent, log.Data)
		if err == nil {
			return &settlement, nil
		}
	}

	return nil, fmt.Errorf("unable to get Settlement data from transaction (%s)", transactionHash)
}
