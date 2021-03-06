package state

import (
	"github.com/propsproject/props-transaction-processor/core/eth-utils"
	"github.com/propsproject/props-transaction-processor/core/proto/pending_props_pb"
	"strconv"
	"strings"
)

const MaxLength = 70

type AddressPart struct {
	Data           string
	DigestStartPos int64
	DigestEndPos   int64
}

func NewPart(data string, startPos, endPos int64) *AddressPart {
	return &AddressPart{data, startPos, endPos}
}

type AddressBuilder struct {
	Prefix   string
	Parts []*AddressPart
}

func NewAddress(prefix string) *AddressBuilder {
	return &AddressBuilder{prefix, make([]*AddressPart, 0)}
}

func (a *AddressBuilder) AddParts(part ...*AddressPart) *AddressBuilder {
	a.Parts = append(a.Parts, part...)
	return a
}

func (a *AddressBuilder) Size() int64 {
	size := int64(len(a.Prefix))
	for _, part := range a.Parts {
		size = size + (part.DigestEndPos - part.DigestStartPos)
	}

	return size
}

func (a *AddressBuilder) IsValidSize() int {
	size := a.Size()

	if size == MaxLength {
		return 0
	} else if size < MaxLength {
		return 1
	} else if size > MaxLength {
		return  -1
	}

	return -1
}

func (a *AddressBuilder) PadWithZeros(str string) string {
	for len(str) != MaxLength {
		str = str + "0"
	}

	return str
}

func (a *AddressBuilder) Build(padWithZeros bool) (string, int) {
	var builder strings.Builder
	builder.WriteString(a.Prefix)

	for _, part := range a.Parts {
		distance := part.DigestEndPos - part.DigestStartPos
		digest := NamespaceManager.HexDigest(part.Data, part.DigestStartPos, part.DigestEndPos)
		if part.DigestStartPos > part.DigestEndPos || distance > int64(len(digest)) || distance > int64(len(digest)) {
			return "", -2
		}

		builder.WriteString(digest)
	}

	str := builder.String()
	if padWithZeros {
		str = a.PadWithZeros(str)
	}

	return str, a.IsValidSize()
}

func TransactionAddress(transaction pending_props_pb.Transaction) (string, int) {

	return NewAddress(NamespaceManager.TransactionPrefix()).AddParts(
		NewPart(strconv.FormatInt(int64(pending_props_pb.Method_value[transaction.GetType().String()]),10), 0, 2),
		NewPart(transaction.GetApplicationId(), 0, 10),
		NewPart(transaction.GetUserId(), 0, 42),
		NewPart(strconv.FormatInt(transaction.GetTimestamp(),10), 0, 10),
		).Build(false)
}

func BalanceAddress(balance pending_props_pb.Balance) (string, int) {
	return BalanceAddressByAppUser(balance.GetApplicationId(), balance.GetUserId())
}

func BalanceAddressByAppUser(applicationId, userId string) (string, int) {
	return NewAddress(NamespaceManager.BalancePrefix()).AddParts(NewPart(applicationId, 0, 10), NewPart(userId, 0, 54)).Build(false)
}

func WalletLinkAddress(address string) (string, int) {
	return NewAddress(NamespaceManager.WalletLinkPrefix()).AddParts(NewPart(eth_utils.NormalizeAddress(address), 0, 64)).Build(false)
}

func BalanceUpdatesTransactionHashAddress(ethTxHash string, address string) (string, int) {
	return NewAddress(NamespaceManager.BalanceUpdatesTransactionHashPrefix()).AddParts(NewPart(eth_utils.NormalizeAddress(ethTxHash), 0, 40),NewPart(eth_utils.NormalizeAddress(address), 0, 24)).Build(false)
}

func SettlementAddress(ethTxHash string) (string, int) {
	return NewAddress(NamespaceManager.SettlementPrefix()).AddParts(NewPart(eth_utils.NormalizeAddress(ethTxHash), 0, 64)).Build(false)
}

func LastEthBlockAddress() (string, int) {
	return NewAddress(NamespaceManager.LastEthBlockPrefix()).AddParts(NewPart("LastEthBlockAddress", 0, 64)).Build(false)
}

func ActivityLogAddress(activity pending_props_pb.ActivityLog) (string, int) {
	return NewAddress(NamespaceManager.ActivityLogPrefix()).AddParts(NewPart(activity.GetApplicationId(), 0, 16), NewPart(activity.GetUserId(), 0, 48)).Build(false)
}

func RewardEntityAddressBySidechainAddress(rewardEntity pending_props_pb.RewardEntity) (string, int) {
	return NewAddress(NamespaceManager.RewardEntityPrefix()).AddParts(NewPart(rewardEntity.GetSidechainAddress(), 0, 60),NewPart(rewardEntity.GetType().String(), 0, 4)).Build(false)
}

func RewardEntityAddressByRewardsAddress(rewardEntity pending_props_pb.RewardEntity) (string, int) {
	return NewAddress(NamespaceManager.RewardEntityPrefix()).AddParts(NewPart(rewardEntity.GetRewardsAddress(), 0, 60),NewPart(rewardEntity.GetType().String(), 0, 4)).Build(false)
}