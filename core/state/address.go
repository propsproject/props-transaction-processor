package state

import (
	"fmt"
	"github.com/propsproject/pending-props/core/eth-utils"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
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

func EarningAddress(earning pending_props_pb.Earning) (string, string, string) {

	prefixPending := NamespaceManager.EarningsPendingPrefix()
	prefixSettled := NamespaceManager.EarningsSettledPrefix()
	prefixRevoked := NamespaceManager.EarningsRevokedPrefix()

	postfix := fmt.Sprintf("%s%s%s", earning.GetDetails().GetApplicationId(), earning.GetDetails().GetUserId(), earning.GetSignature())
	addressParts := []*AddressPart{
		NewPart(earning.GetDetails().GetApplicationId(), 0, 4),
		NewPart(earning.GetDetails().GetUserId(), 0, 4),
		NewPart(postfix, 0, 56),
	}

	pendingAddr, _ := NewAddress(prefixPending).AddParts(addressParts...).Build(false)
	settledAddr, _ := NewAddress(prefixSettled).AddParts(addressParts...).Build(false)
	revokedAddr, _ := NewAddress(prefixRevoked).AddParts(addressParts...).Build(false)

	return pendingAddr, settledAddr, revokedAddr
}

func SettlementAddress(ethTxtHash string) (string, int) {
	return NewAddress(NamespaceManager.SettlementPrefix()).AddParts(NewPart(eth_utils.NormalizeAddress(ethTxtHash), 0, 64)).Build(false)
}

func NonceAddress(publicKey string) (string, int) {
	return NewAddress(NamespaceManager.NoncePrefix()).AddParts(NewPart(publicKey, 0, 64)).Build(false)
}

func BalanceAddress(balance pending_props_pb.Balance) (string, int) {
	//logger.Infof("BalanceAddress:Prefix=%v, Recipient=%v",NamespaceManager.BalancePrefix(), eth_utils.NormalizeAddress(balance.GetRecipientPublicAddress()))
	return NewAddress(NamespaceManager.BalancePrefix()).AddParts(NewPart(balance.GetApplicationId(), 0, 10), NewPart(balance.GetUserId(), 0, 54)).Build(false)
}

func WalletLinkAddress(walletLink pending_props_pb.WalletToUser) (string, int) {
	return NewAddress(NamespaceManager.WalletLinkPrefix()).AddParts(NewPart(eth_utils.NormalizeAddress(walletLink.GetAddress()), 0, 64)).Build(false)
}

func BalanceUpdatesTransactionHashAddress(ethTxHash string, address string) (string, int) {
	return NewAddress(NamespaceManager.BalanceUpdatesTransactionHashPrefix()).AddParts(NewPart(eth_utils.NormalizeAddress(ethTxHash), 0, 40),NewPart(eth_utils.NormalizeAddress(address), 0, 24)).Build(false)
}

func LastEthBlockAddress() (string, int) {
	return NewAddress(NamespaceManager.LastEthBlockPrefix()).AddParts(NewPart("LastEthBlockAddress", 0, 64)).Build(false)
}

func ActivityLogAddress(activity pending_props_pb.ActivityLog) (string, int) {
	return NewAddress(NamespaceManager.ActivityLogPrefix()).AddParts(NewPart(fmt.Sprint(activity.GetDate()), 0, 8), NewPart(activity.GetApplicationId(), 0, 10), NewPart(activity.GetUserId(), 0, 46)).Build(false)
}