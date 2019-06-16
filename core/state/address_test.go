package state

import (
	"fmt"
	"github.com/propsproject/props-transaction-processor/core/eth-utils"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"testing"
)

func TestAddressBuilder(t *testing.T) {
	t.Run("addParts", TestAddParts)
	t.Run("checkSize", TestCheckSize)
	t.Run("buildAddress", TestBuild)
}

//func TestPendingEarningsAddress(t *testing.T) {
//	earningDetails := pending_props_pb.EarningDetails{
//		Timestamp: 1,
//		AmountEarned: 1,
//		RecipientPublicAddress:
//
//	}
//	earning := pending_props_pb.Earning{
//
//	}
//}
func TestAddParts(t *testing.T)  {
	part1 := NewPart("hello", 0, 4)
	part2 := NewPart("hello", 4, 8)

	builder := NewAddress("big").AddParts(part1, part2)

	if len(builder.Parts) != 2 {
		t.Fatalf("expected builder to have 2 parts but got (%v)", len(builder.Parts))
	}
}

func TestCheckSize(t *testing.T)  {
	part1 := NewPart("hello", 0, 25)
	part2 := NewPart("world", 0, 25)
	part3 := NewPart("universe", 0, 14)

	builder := NewAddress(NamespaceManager.HexDigest("hello", 0, 6)).AddParts(part1, part2)

	if builder.IsValidSize()>0 {
		t.Fatalf("expected builder to have invalid address size")
	}

	builder.AddParts(part3)

	if builder.IsValidSize()==0 {
		t.Fatalf("expected builder to have valid address size")
	}
}

func TestBuild(t *testing.T)  {
	prefix := NamespaceManager.HexDigest("what", 0, 6)
	middle := NamespaceManager.HexDigest("are", 0, 32)
	postfix := NamespaceManager.HexDigest("those", 0, 32)

	expectedAddr := fmt.Sprintf("%s%s%s", prefix, middle, postfix)

	middlePart := NewPart("are", 0, 32)
	postFixPart := NewPart("those", 0, 32)

	address, ok := NewAddress(prefix).AddParts(middlePart, postFixPart).Build()
	if ok!=0 {
		t.Fatalf("expected true for ok value")
	}

	if address != expectedAddr {
		t.Fatalf("expected %s to be %s", expectedAddr, address)
	}
}