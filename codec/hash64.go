package codec

import (
	"bytes"
)

type Hash64 [64]byte

func (hash Hash64) Size() uint64 {
	return uint64(64)
}

func (hash Hash64) Encode() []byte {
	return hash[:]
}

func (hash Hash64) Decode(data []byte) interface{} {
	copy(hash[:], data)
	return hash
}

type Hash64s [][64]byte

func (hashes Hash64s) Encode() []byte {
	return Hash64s(hashes).Flatten()
}

func (hashes Hash64s) Decode(data []byte) interface{} {
	hashes = make([][64]byte, len(data)/64)
	for i := 0; i < len(hashes); i++ {
		copy(hashes[i][:], data[i*64:(i+1)*64])
	}
	return hashes
}

func (hashes Hash64s) Size() uint64 {
	return uint64(len(hashes) * 64)
}

func (hashes Hash64s) Flatten() []byte {
	buffer := make([]byte, len(hashes)*64)
	for i := 0; i < len(hashes); i++ {
		copy(buffer[i*64:(i+1)*64], hashes[i][:])
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
