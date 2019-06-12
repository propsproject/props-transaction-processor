package rpc

import (
	"testing"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"math/big"
)

func TestSettlementProcessor_Init(t *testing.T) {
	s := NewEarningsProcessor().Init([]pending_props_pb.Earning{})

	if len(s.earnings) != 0 || s.ethTransactionHash != "id" {
		t.Fatalf("processor data not initialized correctly")
	}
}

func TestSettlementProcessor_settleEarning(t *testing.T) {
	testAddr := "0x42EB768f2244C8811C63729A21A3569731535f06"

	t.Run("ExpectNoChanges", func(t *testing.T) {
		amount := big.NewInt(100)
		earnings := getTestEarnings(testAddr, 2, *big.NewInt(20))
		earnings[0].GetDetails().AmountSettled = earnings[0].GetDetails().AmountEarned
		earnings[1].GetDetails().AmountSettled = earnings[1].GetDetails().AmountEarned
		proc := NewEarningsProcessor().Init(earnings)
		iterator := proc.earnings.NewIterator()
		for iterator.Next(func(earning *pending_props_pb.Earning) error {
			return proc.settleEarning(earning, testAddr, testAddr, amount)
		}) {
			if iterator.Err() != nil {
				t.Fatalf("expected error to be nil got %s", iterator.Err())
			}
		}

		for _, earning := range earnings {
			if earning.GetDetails().GetStatus().String() != pending_props_pb.Status_SETTLED.String() {
				t.Fatalf("expected all earnings to be settled got (%s)", earning.GetDetails().GetStatus().String())
			}
		}

		if amount.Cmp(big.NewInt(100)) != 0 {
			t.Fatalf("expected amount to stay the same want 100, got (%v)", amount)
		}
	})

	t.Run("ExpectCorrectSettledAmounts", func(t *testing.T) {
		amount := big.NewInt(100)
		earnings := getTestEarnings(testAddr, 6, *big.NewInt(20))
		proc := NewEarningsProcessor().Init(earnings)

		iterator := proc.earnings.NewIterator()
		for iterator.Next(func(earning *pending_props_pb.Earning) error {
			return proc.settleEarning(earning, testAddr, testAddr, amount)
		}) {
			if iterator.Err() != nil {
				t.Fatalf("expected error to be nil got %s", iterator.Err())
			}
		}

		for i := 0; i < 5; i++ {
			if earnings[i].GetDetails().GetStatus().String() != pending_props_pb.Status_SETTLED.String() {
				t.Fatalf("expected first 5 earnings to be settled got (%s)", earnings[i].GetDetails().GetStatus().String())
			}
		}

		if earnings[5].GetDetails().GetStatus().String() == pending_props_pb.Status_SETTLED.String() {
			t.Fatalf("expected last earnings to nothing settled got (%v)", earnings[5].GetDetails().GetAmountSettled())
		}

		if amount.Cmp(big.NewInt(0)) != 0 {
			t.Fatalf("expected amount to be 0, got (%v)", amount)
		}
	})
}

func TestSettlementProcessor_revokeEarning(t *testing.T) {
	testAddr := "0x42EB768f2244C8811C63729A21A3569731535f06"
	earnings := getTestEarnings(testAddr, 2, *big.NewInt(20))
	proc := NewEarningsProcessor().Init(earnings).ProcessRevocations(testAddr)
	if len(proc.Errs()) > 0 {
		t.Fatalf("expected errors to be empty got %v", proc.Errs())
	}

	for _, earning := range earnings {
		if earning.GetDetails().GetStatus() != pending_props_pb.Status_REVOKED {
			t.Fatalf("expected all earnings to be revoked got (%s)", earning.GetDetails().GetStatus().String())
		}
	}
}
