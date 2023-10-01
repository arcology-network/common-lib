package codec

import (
	"bytes"
	"encoding/hex"
)

const (
	BYTES20_LEN = 20
)

type Bytes20 [BYTES20_LEN]byte

func (this *Bytes20) Get() interface{} {
	return *this
}

func (this *Bytes20) Set(v interface{}) {
	*this = v.(Bytes20)
}

func (this Bytes20) Size() uint32 {
	return uint32(BYTES20_LEN)
}

func (this Bytes20) Sum(offset uint64) uint64 {
	total := uint64(0)
	for j := offset; j < uint64(len(this)); j++ {
		total += uint64((this)[j])
	}
	return total
}

func (this Bytes20) Clone() interface{} {
	target := Bytes20{}
	copy(target[:], this[:])
	return target
}

func (this Bytes20) Encode() []byte {
	return this[:]
}

func (this Bytes20) EncodeToBuffer(buffer []byte) int {
	copy(buffer, this[:])
	return len(this)
}

func (this Bytes20) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	copy(this[:], buffer)
	return Bytes20(this)
}

func (this Bytes20) Hex() string {
	var bytes [2 * len(this)]byte
	hex.Encode(bytes[:], this[:])
	return string(bytes[:])
}

// func (this Bytes20) UUID(seed uint64) Bytes20 {
// 	buffer := [BYTES20_LEN + 8]byte{}
// 	copy(this[:], buffer[:])
// 	Uint64(uint64(seed)).EncodeToBuffer(buffer[len(this):])
// 	v := sha256.Sum256(buffer[:])

// 	return v[:BYTES20_LEN]
// }

type Byte20s [][BYTES20_LEN]byte

func (this Byte20s) Clone() Byte20s {
	target := make([][BYTES20_LEN]byte, len(this))
	for i := 0; i < len(this); i++ {
		copy(target[i][:], this[i][:])
	}
	return Byte20s(target)
}

func (this Byte20s) Encode() []byte {
	return Byte20s(this).Flatten()
}

func (this Byte20s) EncodeToBuffer(buffer []byte) int {
	for i := 0; i < len(this); i++ {
		copy(buffer[i*BYTES20_LEN:], this[i][:])
	}
	return len(this) * BYTES20_LEN
}

func (this Byte20s) Decode(data []byte) interface{} {
	this = make([][BYTES20_LEN]byte, len(data)/BYTES20_LEN)
	for i := 0; i < len(this); i++ {
		copy(this[i][:], data[i*BYTES20_LEN:(i+1)*BYTES20_LEN])
	}
	return this
}

func (this Byte20s) Size() uint32 {
	return uint32(len(this) * BYTES20_LEN)
}

func (this Byte20s) Flatten() []byte {
	buffer := make([]byte, len(this)*BYTES20_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Byte20s) Len() int {
	return len(this)
}

func (this Byte20s) Less(i, j int) bool {
	return bytes.Compare(this[i][:], this[j][:]) < 0
}

func (this Byte20s) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
