package mhasher

import (
	evmCommon "github.com/arcology-network/evm/common"
)

// GetTxsHash get transactions roothash
func GetTxsHash(src2d [][]byte) evmCommon.Hash {
	if len(src2d) == 0 {
		return evmCommon.Hash{}
	}
	roothash, err := Roothash(src2d, HashType_256)
	if err != nil {
		return evmCommon.Hash{}
	}
	return evmCommon.BytesToHash(roothash)
}
