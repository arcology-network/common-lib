package codec

import (
	"bytes"
	"math"

	ethCommon "github.com/arcology/3rd-party/eth/common"
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

func (hash Hash8) Encode() []byte {
	return hash[:]
}

func (hash Hash8) Decode(data []byte) interface{} {
	copy(hash[:], data)
	return hash
}

type Hash8s []ethCommon.Hash

func (hashes Hash8s) Encode() []byte {
	return Hash8s(hashes).Flatten()
}

func (hashes Hash8s) Decode(data []byte) interface{} {
	hashes = make([]ethCommon.Hash, len(data)/HASH8_LEN)
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
