//go:build nometri
// +build nometri

package mhasher

import (
	evmCommon "github.com/ethereum/go-ethereum/common"
)

const (
	HashType_160 = 20
	HashType_256 = 32
)

func SortByHash(hashes []evmCommon.Hash) ([]uint64, error) {

	return make([]uint64, len(hashes)), nil
}

func Roothash(ls [][]byte, HashType int) ([]byte, error) {
	if HashType == HashType_160 {
		return make([]byte, HashType_160), nil
	} else {
		return make([]byte, HashType_256), nil
	}
}
