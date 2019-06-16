package state

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	globalPrefix = "pending-props"
	prefixEndPos = 6
)

var (
	NamespaceManager *NSMngr

	globalEarningsPrefix = "earnings"

	transactionPrefix = fmt.Sprintf("%s:%s:transaction", globalPrefix, globalEarningsPrefix)

	settlementPrefix 					= fmt.Sprintf("%s:%s:settlements", globalPrefix, globalEarningsPrefix)
	balancePrefix						= fmt.Sprintf("%s:%s:balance", globalPrefix, globalEarningsPrefix)
	noncePrefix 						= fmt.Sprintf("%s:%s:nonce", globalPrefix, globalEarningsPrefix)
	balanceUpdatesTransactionHashPrefix = fmt.Sprintf("%s:%s:bal-rtx", globalPrefix, globalEarningsPrefix)
	lastEthBlockPrefix 					= fmt.Sprintf("%s:%s:lastethblock", globalPrefix, globalEarningsPrefix)

	walletLinkPrefix  = fmt.Sprintf("%s:%s:walletl", globalPrefix, globalEarningsPrefix)
	activityLogPrefix = fmt.Sprintf("%s:%s:activity_log", globalPrefix, globalEarningsPrefix)
)

type NSMngr struct {
	nameSpaces []string
}

func newAddrMngr() *NSMngr {
	return &NSMngr{make([]string, 0)}
}

// ComputePrefix returns namespace prefix of 6 bytes for our default prefixes that will be used to generate various state address's
func (s *NSMngr) computeDefaultPrefix(prefix string) string {
	return s.HexDigest(prefix, 0, prefixEndPos)
}

// HexDigest returns hex digest sub-stringed by designated start and stop pos
func (s *NSMngr) HexDigest(str string, startPos, endPos int64) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))[startPos:endPos]
}

// register a namespace with our state address manager
func (s *NSMngr) registerNamespaces(prefixes ...string) *NSMngr {
	for _, prefix := range prefixes {
		s.nameSpaces = append(s.nameSpaces, s.computeDefaultPrefix(prefix))
	}
	return s
}

func (s *NSMngr) TransactionPrefix() string {
	return s.computeDefaultPrefix(transactionPrefix)
}

func (s *NSMngr) NoncePrefix() string {
	return s.computeDefaultPrefix(noncePrefix)
}

func (s *NSMngr) SettlementPrefix() string {
	return s.computeDefaultPrefix(settlementPrefix)
}

func (s *NSMngr) BalancePrefix() string {
	return s.computeDefaultPrefix(balancePrefix)
}

func (s *NSMngr) WalletLinkPrefix() string {
	return s.computeDefaultPrefix(walletLinkPrefix)
}

func (s *NSMngr) ActivityLogPrefix() string {
	return s.computeDefaultPrefix(activityLogPrefix)
}

func (s *NSMngr) BalanceUpdatesTransactionHashPrefix() string {
	return s.computeDefaultPrefix(balanceUpdatesTransactionHashPrefix)
}

func (s *NSMngr) LastEthBlockPrefix() string {
	return s.computeDefaultPrefix(lastEthBlockPrefix)
}



func (s *NSMngr) Namespaces() []string {
	return s.nameSpaces
}

//// MakeAddress . . .
//func MakeAddress(namespacePrefix, namespaceSuffix string) string {
//	return namespacePrefix + Hexdigest(namespaceSuffix)[:64]
//}
//
//// MakeIdentifierAddress . . .
//func MakeIdentifierAddress(prefix, recipient, application, postfix string) string {
//	return prefix + Hexdigest(recipient)[:4] + Hexdigest(application)[:4] + Hexdigest(postfix)[:56]
//}

func init()  {
	NamespaceManager = newAddrMngr().registerNamespaces(transactionPrefix, settlementPrefix, noncePrefix, balancePrefix, balanceUpdatesTransactionHashPrefix, lastEthBlockPrefix, walletLinkPrefix, activityLogPrefix)
}
