package merkle

import (
	"crypto/sha256"

	"github.com/arcology-network/common-lib/common"
	"golang.org/x/crypto/sha3"
)

type Concatenator struct{}

func (Concatenator) Encode(bytes [][]byte) []byte { return common.Flatten(bytes) }

type Sha256 struct{}

func (Sha256) Hash(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	hash := sha256.Sum256(data)
	return hash[:]
}

type Keccak256 struct{}

func (Keccak256) Hash(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	return sha3.NewLegacyKeccak256().Sum(data)
}
