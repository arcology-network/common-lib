package encoding

import (
	"crypto/sha256"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
)

const (
	UINT8_LEN = 1
)

type Uint8 uint8

func (v Uint8) Size() int {
	return UINT8_LEN
}

func (v Uint8) Encode() []byte {
	data := make([]byte, UINT8_LEN)
	data[0] = uint8(v)
	return data
}

func (_ Uint8) Decode(data []byte) uint8 {
	return uint8(data[0])
}

func (v Uint8) Checksum() ethCommon.Hash {
	return sha256.Sum256(v.Encode())
}

type Uint8s []uint8

func (uint8s Uint8s) Encode() []byte {
	buffer := make([]byte, len(uint8s)*UINT8_LEN)
	for i := range uint8s {
		buffer[i] = uint8s[i]
	}
	return buffer
}

func (uint8s Uint8s) Decode(data []byte) []uint8 {
	uint8s = make([]uint8, len(data)/UINT8_LEN)
	for i := range uint8s {
		uint8s[i] = data[i]
	}
	return uint8s
}
