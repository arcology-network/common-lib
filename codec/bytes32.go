package codec

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
)

const (
	HASH32_LEN = 32
)

type Bytes32 [HASH32_LEN]byte

func (this *Bytes32) Get() interface{} {
	return *this
}

func (this *Bytes32) Set(v interface{}) {
	*this = v.(Bytes32)
}

func (this Bytes32) Size() uint32 {
	return uint32(HASH32_LEN)
}

func (this Bytes32) Sum(offset uint64) uint64 {
	total := uint64(0)
	for j := offset; j < uint64(len(this)); j++ {
		total += uint64((this)[j])
	}
	return total
}

func (this Bytes32) Clone() interface{} {
	target := Bytes32{}
	copy(target[:], this[:])
	return target
}

func (this Bytes32) Encode() []byte {
	return this[:]
}

func (this Bytes32) EncodeToBuffer(buffer []byte) int {
	copy(buffer, this[:])
	return len(this)
}

func (this Bytes32) Decode(buffer []byte) interface{} {
	copy(this[:], buffer)
	return Bytes32(this)
}

// Convert to hex string with 0x prefix.
func (this Bytes32) Hex() string {
	var accHex [2 * len(this)]byte
	hex.Encode(accHex[:], this[:])
	return "0x" + string(accHex[:])
}

func (this Bytes32) UUID(seed uint64) Bytes32 {
	buffer := [HASH32_LEN + 8]byte{}
	copy(this[:], buffer[:])
	Uint64(uint64(seed)).EncodeToBuffer(buffer[len(this):])
	return sha256.Sum256(buffer[:])
}

type Bytes32s [][HASH32_LEN]byte

func (this Bytes32s) Clone() Bytes32s {
	target := make([][HASH32_LEN]byte, len(this))
	for i := 0; i < len(this); i++ {
		copy(target[i][:], this[i][:])
	}
	return Bytes32s(target)
}

func (this Bytes32s) Encode() []byte {
	return Bytes32s(this).Flatten()
}

func (this Bytes32s) EncodeToBuffer(buffer []byte) int {
	for i := 0; i < len(this); i++ {
		copy(buffer[i*HASH32_LEN:], this[i][:])
	}
	return len(this) * HASH32_LEN
}

func (this Bytes32s) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	this = make([][HASH32_LEN]byte, len(buffer)/HASH32_LEN)
	for i := 0; i < len(this); i++ {
		copy(this[i][:], buffer[i*HASH32_LEN:(i+1)*HASH32_LEN])
	}
	return this
}

func (this Bytes32s) Size() uint32 {
	return uint32(len(this) * HASH32_LEN)
}

func (this Bytes32s) Flatten() []byte {
	buffer := make([]byte, len(this)*HASH32_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Bytes32s) Len() int {
	return len(this)
}

func (this Bytes32s) Less(i, j int) bool {
	return bytes.Compare(this[i][:], this[j][:]) < 0
}

func (this Bytes32s) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
