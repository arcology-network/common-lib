package codec

import (
	"bytes"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
)

const (
	HASH32_LEN = 32
)

type Hash32 [32]byte

func (this *Hash32) Get() interface{} {
	return *this
}

func (this *Hash32) Set(v interface{}) {
	*this = v.(Hash32)
}

func (hash Hash32) Size() uint32 {
	return uint32(HASH32_LEN)
}

func (hash Hash32) Encode() []byte {
	return hash[:]
}

func (hash Hash32) Decode(data []byte) interface{} {
	copy(hash[:], data)
	return hash
}

type Hash32s []ethCommon.Hash

func (hashes Hash32s) Encode() []byte {
	return Hash32s(hashes).Flatten()
}

func (hashes Hash32s) Decode(data []byte) interface{} {
	hashes = make([]ethCommon.Hash, len(data)/HASH32_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(hashes[i][:], data[i*HASH32_LEN:(i+1)*HASH32_LEN])
	}
	return hashes
}

func (hashes Hash32s) Size() uint32 {
	return uint32(len(hashes) * HASH32_LEN)
}

func (hashes Hash32s) Flatten() []byte {
	buffer := make([]byte, len(hashes)*HASH32_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(buffer[i*HASH32_LEN:(i+1)*HASH32_LEN], hashes[i][:])
	}
	return buffer
}

func (hashes Hash32s) Len() int {
	return len(hashes)
}

func (hashes Hash32s) Less(i, j int) bool {
	return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
}

func (hashes Hash32s) Swap(i, j int) {
	hashes[i], hashes[j] = hashes[j], hashes[i]
}
