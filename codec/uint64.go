package codec

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
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

func (Uint64) Size() uint32 {
	return UINT64_LEN
}

func (this Uint64) Encode() []byte {
	buffer := make([]byte, UINT64_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Uint64) EncodeToBuffer(buffer []byte) {
	binary.LittleEndian.PutUint64(buffer, uint64(this))
}

func (this Uint64) Decode(data []byte) interface{} {
	this = Uint64(binary.LittleEndian.Uint64(data))
	return Uint64(this)
}

func (v Uint64) Checksum() ethCommon.Hash {
	return sha256.Sum256(v.Encode())
}

type Uint64s []uint64

func (this Uint64s) Get() interface{} {
	return this.Sum()
}

func (this Uint64s) Set(v interface{}) {
	this = append(this, v.(uint64))
}

func (this Uint64s) Sum() int64 {
	sum := int64(0)
	for i := range this {
		sum += int64(this[i])
	}
	return sum
}

func (this Uint64s) Unique() []uint64 {
	sort.SliceStable(this, func(i, j int) bool {
		return this[i] < this[j]
	})

	uniqueV := make([]uint64, 0, len(this))
	current := uint64(this[0])
	for i := 0; i < len(this); i++ {
		if current != uint64(this[i]) {
			uniqueV = append(uniqueV, current)
			current = uint64(this[i])
		}
	}

	if current != uniqueV[len(uniqueV)-1] {
		uniqueV = append(uniqueV, current)
	}

	return uniqueV
}

func (this Uint64s) Encode() []byte {
	buffer := make([]byte, len(this)*UINT64_LEN)
	for i := range this {
		copy(buffer[i*UINT64_LEN:(i+1)*UINT64_LEN], Uint64(this[i]).Encode())
	}
	return buffer
}

func (this Uint64s) Decode(data []byte) interface{} {
	this = make([]uint64, len(data)/UINT64_LEN)
	for i := range this {
		this[i] = uint64(Uint64(this[i]).Decode(data[i*UINT64_LEN : (i+1)*UINT64_LEN]).(Uint64))
	}
	return Uint64s(this)
}
