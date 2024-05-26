package tools

import (
	slice "github.com/arcology-network/common-lib/exp/slice"
	evmCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func CalculateHash(hashes []*evmCommon.Hash) evmCommon.Hash {
	return evmCommon.BytesToHash(crypto.Keccak256(slice.Concate(hashes, func(v *evmCommon.Hash) []byte { return (*v)[:] })))
}
func RlpHash(x interface{}) (h evmCommon.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
