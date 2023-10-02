package merkle

import (
	"crypto/sha256"

	"github.com/arcology-network/common-lib/common"
	"golang.org/x/crypto/sha3"
)

type Concatenator struct{}

func (Concatenator) Encode(bytes [][]byte) []byte { return common.Flatten(bytes) }

// func Keccak256(bytes []byte) []byte {
// 	keccak := sha3.NewLegacyKeccak256()
// 	return keccak.Sum(bytes)
// }

func Sha256(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	hash := sha256.Sum256(data)
	return hash[:]
}

func Keccak256(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	return sha3.NewLegacyKeccak256().Sum(data)
}
