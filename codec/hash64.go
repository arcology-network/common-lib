package codec

import (
	"bytes"
)

const (
	HASH64_LEN = 64
)

type Hash64 [HASH64_LEN]byte

func (hash Hash64) Size() uint32 {
	return uint32(HASH64_LEN)
}

func (hash Hash64) Encode() []byte {
	return hash[:]
}

func (hash Hash64) Decode(data []byte) interface{} {
	copy(hash[:], data)
	return Hash64(hash)
}

type Hash64s [][HASH64_LEN]byte

func (hashes Hash64s) Encode() []byte {
	return Hash64s(hashes).Flatten()
}

func (this Hash64s) Decode(data []byte) interface{} {
	this = make([][HASH64_LEN]byte, len(data)/HASH64_LEN)
	for i := 0; i < len(this); i++ {
		copy(this[i][:], data[i*HASH64_LEN:(i+1)*HASH64_LEN])
	}
	return [][HASH64_LEN]byte(this)
}

func (hashes Hash64s) Size() uint32 {
	return uint32(len(hashes) * HASH64_LEN)
}

func (hashes Hash64s) Flatten() []byte {
	buffer := make([]byte, len(hashes)*HASH64_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(buffer[i*HASH64_LEN:(i+1)*HASH64_LEN], hashes[i][:])
	}
	return buffer
}

func (hashes Hash64s) Len() int {
	return len(hashes)
}

func (hashes Hash64s) Less(i, j int) bool {
	return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
}

func (hashes Hash64s) Swap(i, j int) {
	hashes[i], hashes[j] = hashes[j], hashes[i]
}
