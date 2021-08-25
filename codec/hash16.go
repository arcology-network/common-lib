package codec

import (
	"bytes"
	"math"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
)

const (
	HASH16_LEN = 16
)

type Hash16 [HASH16_LEN]byte

func NewHash16(v byte) Hash16 {
	hash16 := [HASH16_LEN]byte{}
	for i := range hash16 {
		hash16[i] = v
	}
	return hash16
}

func (Hash16) FromSlice(v []byte) Hash16 {
	hash16 := [HASH16_LEN]byte{}
	length := math.Min(float64(HASH16_LEN), float64(len(v)))
	for i := 0; i < int(length); i++ {
		hash16[i] = v[i]
	}
	return hash16
}

func (this *Hash16) Get() interface{} {
	return *this
}

func (this *Hash16) Set(v interface{}) {
	*this = v.(Hash16)
}

func (hash Hash16) Size() uint32 {
	return uint32(HASH16_LEN)
}

func (hash Hash16) FromBytes(bytes []byte) Hash16 {
	hash = Hash16{}
	copy(hash[:], bytes)
	return hash
}

func (hash Hash16) Encode() []byte {
	return hash[:]
}

func (hash Hash16) Decode(data []byte) interface{} {
	copy(hash[:], data)
	return hash
}

type Hash16s []ethCommon.Hash

func (hashes Hash16s) Encode() []byte {
	return Hash16s(hashes).Flatten()
}

func (hashes Hash16s) Decode(data []byte) interface{} {
	hashes = make([]ethCommon.Hash, len(data)/HASH16_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(hashes[i][:], data[i*HASH16_LEN:(i+1)*HASH16_LEN])
	}
	return hashes
}

func (hashes Hash16s) Size() uint32 {
	return uint32(len(hashes) * HASH16_LEN)
}

func (hashes Hash16s) Flatten() []byte {
	buffer := make([]byte, len(hashes)*HASH16_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(buffer[i*HASH16_LEN:(i+1)*HASH16_LEN], hashes[i][:])
	}
	return buffer
}

func (hashes Hash16s) Len() int {
	return len(hashes)
}

func (hashes Hash16s) Less(i, j int) bool {
	return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
}

func (hashes Hash16s) Swap(i, j int) {
	hashes[i], hashes[j] = hashes[j], hashes[i]
}
