package codec

import (
	"bytes"
	"math"

	evmCommon "github.com/arcology-network/evm/common"
)

const (
	HASH8_LEN = 8
)

type Hash8 [HASH8_LEN]byte

func NewHash8(v byte) Hash8 {
	hash8 := [HASH8_LEN]byte{}
	for i := range hash8 {
		hash8[i] = v
	}
	return hash8
}

func (Hash8) FromSlice(v []byte) Hash8 {
	hash8 := [HASH8_LEN]byte{}
	length := math.Min(float64(HASH8_LEN), float64(len(v)))
	for i := 0; i < int(length); i++ {
		hash8[i] = v[i]
	}
	return hash8
}

func (this Hash8) Clone() interface{} {
	target := Hash8{}
	copy(target[:], this[:])
	return target
}

func (this *Hash8) Get() interface{} {
	return *this
}

func (this *Hash8) Set(v interface{}) {
	*this = v.(Hash8)
}

func (hash Hash8) Size() uint32 {
	return uint32(HASH8_LEN)
}

func (hash Hash8) FromBytes(bytes []byte) Hash8 {
	hash = Hash8{}
	copy(hash[:], bytes)
	return hash
}

func (this Hash8) Sum(offset uint64) uint64 {
	total := uint64(0)
	for j := offset; j < uint64(len(this)); j++ {
		total += uint64(this[j])
	}
	return total
}

func (hash Hash8) Encode() []byte {
	return hash[:]
}

func (this Hash8) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	copy(this[:], buffer)
	return Hash8(this)
}

type Hash8s []evmCommon.Hash

func (hashes Hash8s) Encode() []byte {
	return Hash8s(hashes).Flatten()
}

func (hashes Hash8s) Decode(data []byte) interface{} {
	hashes = make([]evmCommon.Hash, len(data)/HASH8_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(hashes[i][:], data[i*HASH8_LEN:(i+1)*HASH8_LEN])
	}
	return hashes
}

func (hashes Hash8s) Size() uint32 {
	return uint32(len(hashes) * HASH8_LEN)
}

func (hashes Hash8s) Flatten() []byte {
	buffer := make([]byte, len(hashes)*HASH8_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(buffer[i*HASH8_LEN:(i+1)*HASH8_LEN], hashes[i][:])
	}
	return buffer
}

func (hashes Hash8s) Len() int {
	return len(hashes)
}

func (hashes Hash8s) Less(i, j int) bool {
	return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
}

func (hashes Hash8s) Swap(i, j int) {
	hashes[i], hashes[j] = hashes[j], hashes[i]
}
