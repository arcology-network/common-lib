package merkle

import (
	"crypto/sha256"

	"golang.org/x/crypto/sha3"
)

func Sha256(bytes []byte) []byte {
	hash := sha256.Sum256(bytes)
	return hash[:]
}

func Keccak256(bytes []byte) []byte {
	keccak := sha3.NewLegacyKeccak256()
	return keccak.Sum(bytes)
}
