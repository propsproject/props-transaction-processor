package rpc

import (
	"fmt"
	"github.com/propsproject/pending-props/core/eth-utils"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"math/big"
	"sort"
	"sync"
)

type EarningsProcessor struct {
	earnings Earnings
	settlementData *eth_utils.SettlementData
	ethTransactionHash string
	appEthAddr string
	once sync.Once
	errs []error
}

func (e *EarningsProcessor) Init(earnings Earnings) *EarningsProcessor {
	e.once.Do(func() {
		e.earnings = earnings
		sort.Sort(earnings)
	})

	return e
}

func (e *EarningsProcessor) processSettlementsCB(earning *pending_props_pb.Earning) error {
	logger.Infof("earning: %v\n%v", earning, e.settlementData)
	return e.settleEarning(earning, e.settlementData.Recipient.String(), e.settlementData.From.String(), e.settlementData.Amount)
}

func (e *EarningsProcessor) revokeEarning(earning *pending_props_pb.Earning) error {
	if earning.GetDetails().GetStatus() == pending_props_pb.Status_REVOKED {
		return nil
	}

	check := earning.GetDetails().GetApplicationId()
	if e.appEthAddr != check {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("unauthorized attempt to revoke earning, public keys do not match got (%s) want (%s)", e.appEthAddr, check)}
	}

	earning.Details.Status = pending_props_pb.Status_REVOKED

	return nil
}

func (e *EarningsProcessor) settleEarning(earning *pending_props_pb.Earning, recipient, appEthAddr string, amountTransferred *big.Int) error {
	logger.Infof("%v", amountTransferred)

	earningRec := eth_utils.NormalizeAddress(earning.GetDetails().GetUserId())
	earningApp := eth_utils.NormalizeAddress(earning.GetDetails().GetApplicationId())

	recipient = eth_utils.NormalizeAddress(recipient)
	application := eth_utils.NormalizeAddress(appEthAddr)

	if earningRec != recipient {
		return fmt.Errorf("unauthorized attempt to settle earning, recipient public keys do not match got (%s) want (%s)", recipient, earningRec )
	}

	if application != earningApp {
		return fmt.Errorf("unauthorized attempt to settle earning, public keys do not match got (%s) want (%s)", application, earningApp)
	}

	settlementFrom := eth_utils.NormalizeAddress(e.settlementData.From.String())
	settlementRec := eth_utils.NormalizeAddress(e.settlementData.Recipient.String())

	if settlementFrom != application {
		return fmt.Errorf("unauthorized attempt to settle earning, applciation public keys do not match got earning -> settlement event data (%s) want (%s)", settlementFrom, application)
	}

	if settlementRec != earningRec {
		return fmt.Errorf("unauthorized attempt to settle earning, recipient public keys do not match got earning -> settlement event data (%s) want (%s)", settlementRec, recipient)
	}


	earned, ok := new(big.Int).SetString(earning.GetDetails().AmountEarned, 10)
	if !ok {
		return fmt.Errorf("invalid earning.GetDetails().AmountEarned=%v", earning.GetDetails().AmountEarned)
	}
	settled, ok := new(big.Int).SetString(earning.GetDetails().AmountSettled, 10)
	if !ok {
		return fmt.Errorf("invalid earning.GetDetails().AmountSettled=%v", earning.GetDetails().AmountSettled)
	}
	amountOwed := settled.Abs(settled.Sub(settled, earned))

	if amountTransferred.Cmp(big.NewInt(0)) == 0 {
		return nil
	} else if amountOwed.Cmp(big.NewInt(0)) <= 0 { // nothing owed
		earning.GetDetails().Status = pending_props_pb.Status_SETTLED
	} else if amountTransferred.Cmp(amountOwed) < 0 { // transfer amount does not cover amount owed
		earning.Details.AmountSettled = settled.Sub(settled, amountTransferred).String()
		amountTransferred.SetInt64(0)
	} else if amountTransferred.Cmp(amountOwed) > 0 { // transfer amount exceeds amount owed
		earning.Details.AmountSettled = amountOwed.String()
		amountTransferred.Sub(amountTransferred, amountOwed)
	} else if amountTransferred.Cmp(amountOwed) == 0 { // transfer amount equals amount owed
		earning.Details.AmountSettled = earning.Details.AmountEarned
		amountTransferred.SetInt64(0)
	}

	logger.Infof("%v", earning)
	if earning.GetDetails().GetAmountEarned() == earning.GetDetails().GetAmountSettled() {
		earning.Details.Status = pending_props_pb.Status_SETTLED
		earning.SettledByTransaction = e.ethTransactionHash
	}

	return nil
}

func (e *EarningsProcessor) ProcessSettlements(ethTransactionHash string) *EarningsProcessor {
	e.ethTransactionHash = ethTransactionHash
	settlementData, err := e.GetSettlement()
	if err != nil {
		e.errs = append(e.errs, err)
		return e
	}

	logger.Infof("%v", settlementData)


	err = e.ValidateSettlement(settlementData)
	if err != nil {
		e.errs = append(e.errs, err)
		return e
	}
	e.settlementData = settlementData
	logger.Infof("%v", e.settlementData)

	iterator := e.earnings.NewIterator()
	for iterator.Next(e.processSettlementsCB) {
		if iterator.Err() != nil {
			e.errs = append(e.errs, iterator.Err())
		}
	}

	return e
}

func (e *EarningsProcessor) ProcessRevocations(appEthAddr string) *EarningsProcessor {
	e.appEthAddr = eth_utils.NormalizeAddress(appEthAddr)
	iterator := e.earnings.NewIterator()
	for iterator.Next(e.revokeEarning) {
		if iterator.Err() != nil {
			e.errs = append(e.errs, iterator.Err())
		}
	}

	return e
}

func (e *EarningsProcessor) GetSettlement() (*eth_utils.SettlementData ,error) {
	return eth_utils.GetEthTransactionSettlementData(e.ethTransactionHash)
}

func (e *EarningsProcessor) ValidateSettlement(settlement *eth_utils.SettlementData) error {
	if settlement.Amount.Cmp(big.NewInt(1)) < 0 {
		return fmt.Errorf("invalid settlement amount (%x)", settlement.Amount)
	}

	if settlement.Recipient.String() == "" {
		return fmt.Errorf("invalid recipient address (%s)", settlement.Recipient.String())
	}

	if settlement.From.String() == "" {
		return fmt.Errorf("invalid from address (%s)", settlement.From.String())
	}

	return nil
}

func (e *EarningsProcessor) Errs() []error {
	return e.errs
}

func NewEarningsProcessor() *EarningsProcessor {
	return &EarningsProcessor{}
}

