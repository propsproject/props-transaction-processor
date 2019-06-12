package rpc

import (
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"math/big"
	"time"
	"math/rand"
	"testing"
	"fmt"
)

func randomTimestamp() time.Time {
	randomTime := rand.Int63n(time.Now().Unix() - 94608000) + 94608000

	randomNow := time.Unix(randomTime, 0)

	return randomNow
}

func getTestEarnings(address string, max int, amount big.Int) Earnings {
	address = eth_utils.NormalizeAddress(address)
	earnings := make([]pending_props_pb.Earning, 0)
	for i := 0; i < max; i++ {
		earnings = append(earnings, pending_props_pb.Earning{
			Details: &pending_props_pb.EarningDetails{
				Timestamp: randomTimestamp().Unix(),
				AmountEarned: amount.String(),
				RecipientPublicAddress: address,
				ApplicationPublicAddress: address,
			},
		})
	}

	return earnings
}

func TestEarningsIterator_Next(t *testing.T) {
	iterator := getTestEarnings("0x42EB768f2244C8811C63729A21A3569731535f06", 5, *big.NewInt(10)).NewIterator()
	for iterator.Next(func(earning *pending_props_pb.Earning) error {
		fmt.Printf("%v\n", earning.Details.GetTimestamp())
		return nil
	}) {
		if iterator.Err() != nil {
			t.Fatalf("expected error to be nil got (%s)", iterator.Err())
		}
	}

	hashNext := iterator.Next(func(earning *pending_props_pb.Earning) error {
		fmt.Printf("%v\n", earning.Details.GetTimestamp())
		return nil
	})
	if hashNext == true {
		t.Fatalf("expected hashNext to be false got (%v)", hashNext)
	}
}