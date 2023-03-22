package codec

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
)

const (
	HASH32_LEN = 32
)

type Hash32 [HASH32_LEN]byte

func (this *Hash32) Get() interface{} {
	return *this
}

func (this *Hash32) Set(v interface{}) {
	*this = v.(Hash32)
}

func (this Hash32) Size() uint32 {
	return uint32(HASH32_LEN)
}

func (this Hash32) Clone() Hash32 {
	target := Hash32{}
	copy(target[:], this[:])
	return target
}

func (this Hash32) Encode() []byte {
	return this[:]
}

func (this Hash32) EncodeToBuffer(buffer []byte) int {
	copy(buffer, this[:])
	return len(this)
}

func (this Hash32) Decode(buffer []byte) interface{} {
	copy(this[:], buffer)
	return Hash32(this)
}

func (this Hash32) Hex() string {
	var accHex [2 * len(this)]byte
	hex.Encode(accHex[:], this[:])
	return string(accHex[:])
}

func (this Hash32) UUID(seed uint64) Hash32 {
	buffer := [HASH32_LEN + 8]byte{}
	copy(this[:], buffer[:])
	Uint64(uint64(seed)).EncodeToBuffer(buffer[len(this):])
	return sha256.Sum256(buffer[:])
}

type Hash32s [][HASH32_LEN]byte

func (this Hash32s) Clone() Hash32s {
	target := make([][HASH32_LEN]byte, len(this))
	for i := 0; i < len(this); i++ {
		copy(target[i][:], this[i][:])
	}
	return Hash32s(target)
}

func (this Hash32s) Encode() []byte {
	return Hash32s(this).Flatten()
}

func (this Hash32s) EncodeToBuffer(buffer []byte) int {
	for i := 0; i < len(this); i++ {
		copy(buffer[i*HASH32_LEN:], this[i][:])
	}
	return len(this) * HASH32_LEN
}

func (this Hash32s) Decode(data []byte) interface{} {
	this = make([][HASH32_LEN]byte, len(data)/HASH32_LEN)
	for i := 0; i < len(this); i++ {
		copy(this[i][:], data[i*HASH32_LEN:(i+1)*HASH32_LEN])
	}
	return this
}

func (this Hash32s) Size() uint32 {
	return uint32(len(this) * HASH32_LEN)
}

func (this Hash32s) Flatten() []byte {
	buffer := make([]byte, len(this)*HASH32_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Hash32s) Len() int {
	return len(this)
}

func (this Hash32s) Less(i, j int) bool {
	return bytes.Compare(this[i][:], this[j][:]) < 0
}

func (this Hash32s) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
