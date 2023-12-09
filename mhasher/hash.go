package mhasher

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
)

// GetTxsHash get transactions roothash
func GetTxsHash(src2d [][]byte) ethCommon.Hash {
	if len(src2d) == 0 {
		return ethCommon.Hash{}
	}
	roothash, err := Roothash(src2d, HashType_256)
	if err != nil {
		return ethCommon.Hash{}
	}
	return ethCommon.BytesToHash(roothash)
}
