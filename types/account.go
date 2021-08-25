package types

import (
	"math/big"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
	"github.com/HPISTechnologies/common-lib/common"
)

type AccountInfo struct {
	Address ethCommon.Address
	Account Account
}

// Account is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     ethCommon.Hash // merkle root of the storage trie
	CodeHash []byte
}

func (a *Account) Clone() *Account {
	return &Account{
		Nonce:    a.Nonce,
		Balance:  a.Balance,
		Root:     a.Root,
		CodeHash: common.ArrayCopy(a.CodeHash),
	}
}
