package codec

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"

	ethCommon "github.com/arcology/3rd-party/eth/common"
)

const (
	UINT64_LEN = 8
)

type Uint64 uint64

func (this *Uint64) Get() interface{} {
	return *this
}

func (this *Uint64) Set(v interface{}) {
	*this = v.(Uint64)
}

func (v Uint64) Size() uint32 {
	return UINT64_LEN
}

func (v Uint64) Encode() []byte {
	data := make([]byte, UINT64_LEN)
	binary.LittleEndian.PutUint64(data[0:], uint64(v))
	return data
}

func (Uint64) Decode(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data[0:UINT64_LEN])
}

func (v Uint64) Checksum() ethCommon.Hash {
	return sha256.Sum256(v.Encode())
}

type Uint64s []Uint64

func (this Uint64s) Get() interface{} {
	return this.Sum()
}

func (this Uint64s) Set(v interface{}) {
	this = append(this, v.(Uint64))
}

func (this Uint64s) Sum() int64 {
	sum := int64(0)
	for i := range this {
		sum += int64(this[i])
	}
	return sum
}

func (uint64s Uint64s) Unique() []uint64 {
	sort.SliceStable(uint64s, func(i, j int) bool {
		return uint64s[i] < uint64s[j]
	})

	uniqueV := make([]uint64, 0, len(uint64s))
	current := uint64(uint64s[0])
	for i := 0; i < len(uint64s); i++ {
		if current != uint64(uint64s[i]) {
			uniqueV = append(uniqueV, current)
			current = uint64(uint64s[i])
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
		copy(buffer[i*UINT64_LEN:(i+1)*UINT64_LEN], Uint64(uint64s[i]).Encode())
	}
	return buffer
}

func (uint64s Uint64s) Decode(data []byte) interface{} {
	uint64s = make([]Uint64, len(data)/UINT64_LEN)
	for i := range uint64s {
		uint64s[i] = Uint64(uint64s[i].Decode(data[i*UINT64_LEN : (i+1)*UINT64_LEN]))
	}
	return Uint64s(uint64s)
}
