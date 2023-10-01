package codec

import (
	"bytes"
	"math"
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

func (this Hash16) Clone() interface{} {
	target := Hash16{}
	copy(target[:], this[:])
	return target
}

func (hash Hash16) FromBytes(bytes []byte) Hash16 {
	hash = Hash16{}
	copy(hash[:], bytes)
	return hash
}

func (hash Hash16) Encode() []byte {
	return hash[:]
}

func (this Hash16) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	copy(this[:], buffer)
	return Hash16(this)
}

type Hash16s [][16]byte

func (this Hash16s) Clone() Hash16s {
	target := make([][HASH16_LEN]byte, len(this))
	for i := 0; i < len(this); i++ {
		copy(target[i][:], this[i][:])
	}
	return Hash16s(target)
}

func (this Hash16s) Encode() []byte {
	return Hash16s(this).Flatten()
}

func (this Hash16s) EncodeToBuffer(buffer []byte) int {
	for i := 0; i < len(this); i++ {
		copy(buffer[i*HASH16_LEN:], this[i][:])
	}
	return len(this) * HASH16_LEN
}

func (this Hash16s) Decode(data []byte) interface{} {
	this = make([][16]byte, len(data)/HASH16_LEN)
	for i := 0; i < len(this); i++ {
		copy(this[i][:], data[i*HASH16_LEN:(i+1)*HASH16_LEN])
	}
	return this
}

func (this Hash16s) Size() uint32 {
	return uint32(len(this) * HASH16_LEN)
}

func (this Hash16s) Flatten() []byte {
	buffer := make([]byte, len(this)*HASH16_LEN)
	for i := 0; i < len(this); i++ {
		copy(buffer[i*HASH16_LEN:(i+1)*HASH16_LEN], this[i][:])
	}
	return buffer
}

func (this Hash16s) Len() int {
	return len(this)
}

func (this Hash16s) Less(i, j int) bool {
	return bytes.Compare(this[i][:], this[j][:]) < 0
}

func (this Hash16s) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
