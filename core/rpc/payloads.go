package rpc

import (
	"github.com/propsproject/pending-props/core/state"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"fmt"
)

type SettlementRPCPayload struct {
	EthTransactionHash string `json:"eth_transaction_hash"`
	Recipient          string `json:"recipient"`
	PendingAddresses []string `json:"pending_addresses"`
	Timestamp          int64  `json:"timestamp"`
}

func (s *SettlementRPCPayload) Validate(context *processor.Context) error {
	settlements, err := state.NewState(context).GetSettlements(s.EthTransactionHash)
	if err != nil {
		return err
	}

	if len(settlements) > 0 {
		if len(settlements[0].GetEarningAddresses()) > 0 {
			return fmt.Errorf("ethereum transaction already processed for earning settlement (%s)", s.EthTransactionHash)
		}
	}

	return nil
}

type RevokeRPCPaylod struct {
	Addresses []string `json:"addresses"`
	Timestamp   int64  `json:"timestamp"`
}

func (r *RevokeRPCPaylod) Validate() error {
	if len(r.Addresses) < 0 {
		return fmt.Errorf("must specify atleast one earning namespace address to revoke earning")
	}

	return nil
}
