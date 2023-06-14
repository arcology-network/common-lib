package types

import (
	"math/big"
	"reflect"
	"testing"

	evmCommon "github.com/arcology-network/evm/common"
)

func TestAccountClone(t *testing.T) {
	a := Account{
		Nonce:    uint64(10),
		Balance:  big.NewInt(20),
		Root:     evmCommon.BytesToHash([]byte{1, 2, 3, 4, 5, 6}),
		CodeHash: []byte{11, 12, 13, 14, 15, 16},
	}
	b := a.Clone()

	b.Nonce = uint64(100)
	b.Balance = big.NewInt(2000)
	b.Root = evmCommon.BytesToHash([]byte{21, 22, 23, 24, 25, 26})
	b.CodeHash[3] = byte(90)

	if reflect.DeepEqual(a.Nonce, b.Nonce) {
		t.Error("field nonce err in account clone!", a.Nonce, b.Nonce)
	}
	if reflect.DeepEqual(a.Balance, b.Balance) {
		t.Error("field Balance err in account clone!", a.Balance, b.Balance)
	}
	if reflect.DeepEqual(a.Root, b.Root) {
		t.Error("field Root err in account clone!", a.Root, b.Root)
	}
	if reflect.DeepEqual(a.CodeHash, b.CodeHash) {
		t.Error("field CodeHash err in account clone!", a.CodeHash, b.CodeHash)
	}
}
