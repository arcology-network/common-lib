package encoding

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"

	evmCommon "github.com/ethereum/go-ethereum/common"
)

const (
	UINT64_LEN = 8
)

type Uint64 uint64

func (v Uint64) Size() uint32 {
	return UINT64_LEN
}

func (v Uint64) Encode() []byte {
	data := make([]byte, UINT64_LEN)
	binary.LittleEndian.PutUint64(data[0:], uint64(v))
	return data
}

func (_ Uint64) Decode(data []byte) uint64 {
	return uint64(binary.LittleEndian.Uint64(data[0:UINT64_LEN]))
}

func (v Uint64) Checksum() evmCommon.Hash {
	return sha256.Sum256(v.Encode())
}

type Uint64s []uint64

func (uint64s Uint64s) Unique() []uint64 {
	sort.SliceStable(uint64s, func(i, j int) bool {
		return uint64s[i] < uint64s[j]
	})

	uniqueV := make([]uint64, 0, len(uint64s))
	current := uint64s[0]
	for i := 0; i < len(uint64s); i++ {
		if current != uint64s[i] {
			uniqueV = append(uniqueV, current)
			current = uint64s[i]
		}
	}

	if current != uniqueV[len(uniqueV)-1] {
		uniqueV = append(uniqueV, current)
	}
	return uniqueV
}

func (uint64s Uint64s) Encode() []byte {
	buffer := make([]byte, len(uint64s)*UINT64_LEN)
	for i := range uint64s {
		binary.LittleEndian.PutUint64(buffer[i*UINT64_LEN:(i+1)*UINT64_LEN], uint64s[i])
	}
	return buffer
}

func (uint64s Uint64s) Decode(data []byte) []uint64 {
	uint64s = make([]uint64, len(data)/UINT64_LEN)
	for i := range uint64s {
		uint64s[i] = binary.LittleEndian.Uint64(data[i*UINT64_LEN : (i+1)*UINT64_LEN])
	}
	return uint64s
}
