package mhasher

import (
	"bytes"
	"testing"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
	"github.com/arcology-network/3rd-party/eth/crypto/sha3"
)

func TestUniqueStringss(t *testing.T) {
	str := []string{"1", "2", "3", ""}
	result, _ := UniqueStrings(str)
	for i := range result {
		if str[i] != (result[i]) {
			t.Error("Error")
		}
	}
}

func TestUniqueSingleString(t *testing.T) {
	str := []string{"1234"}
	result, _ := UniqueStrings(str)
	for i := range result {
		if str[i] != (result[i]) {
			t.Error("Error")
		}
	}
}

func TestUniqueEmptyStrings(t *testing.T) {
	str := []string{"", "", ""}
	result, _ := UniqueStrings(str)
	for i := range result {
		if str[i] != (result[i]) {
			t.Error("Error")
		}
	}
}

func TestCalculateRoothash(t *testing.T) {
	//datas := [][]byte{{4, 5, 6}, {1, 2, 3}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}, {16, 17, 18}, {19, 20, 21}}
	byteset := [][]byte{
		ethCommon.HexToHash("11686a2cfd1c30b6aed43e93c9254e7d819d0e58866ec7d0b50e59a75865a0bf").Bytes(),
		ethCommon.HexToHash("6131d51e229570f83c87ca48ee6d789b52dac9e7df08da8b9e4704dfe305568d").Bytes(),
		ethCommon.HexToHash("0730158bb52019425e1063b40be465a6014aad7686a0793a6a13cd1bd11fd0f3").Bytes(),
		ethCommon.HexToHash("f1f83e1448a720f073c542c905da761940ebd16ad28d920fe8785a6cccf5bb24").Bytes(),
	}
	//

	sha3256data := Sha3256(byteset)
	for i := range byteset {
		sha256 := sha3.Sum256(byteset[i][:])
		if !bytes.Equal(sha256[:], sha3256data[i]) {
			t.Error("Error")
		}
	}
}
