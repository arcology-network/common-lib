package tools

import (
	"testing"

	evmCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func TestCalculateHash(t *testing.T) {
	h0 := evmCommon.HexToHash("0x01")
	h1 := evmCommon.HexToHash("0x02")

	expected := evmCommon.BytesToHash(crypto.Keccak256(append(append([]byte{}, h0[:]...), h1[:]...)))
	if got := CalculateHash([]*evmCommon.Hash{&h0, &h1}); got != expected {
		t.Error("Error: CalculateHash should hash the concatenated hash bytes")
	}
}

func TestRlpHash(t *testing.T) {
	input := struct {
		Name  string
		Count uint64
	}{Name: "arcology", Count: 9}

	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, input)
	var expected evmCommon.Hash
	hw.Sum(expected[:0])

	if got := RlpHash(input); got != expected {
		t.Error("Error: RlpHash should match the manual RLP keccak hash")
	}
}