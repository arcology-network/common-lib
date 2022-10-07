package codec

import (
	"crypto/sha256"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
)

const (
	UINT8_LEN = 1
)

type Uint8 uint8

func (this *Uint8) Get() interface{} {
	return *this
}

func (this *Uint8) Set(v interface{}) {
	*this = v.(Uint8)
}

func (v Uint8) Size() uint32 {
	return UINT8_LEN
}

func (v Uint8) Encode() []byte {
	buffer := make([]byte, UINT8_LEN)
	buffer[0] = uint8(v)
	return buffer
}

func (v Uint8) EncodeToBuffer(buffer []byte) {
	buffer[0] = uint8(v)
}

func (this Uint8) Decode(data []byte) interface{} {
	this = Uint8(data[0])
	return this
}

func (v Uint8) Checksum() ethCommon.Hash {
	return sha256.Sum256(v.Encode())
}

type Uint8s []uint8

func (this Uint8s) Get() interface{} {
	return this.Sum()
}

func (this *Uint8s) Set(v interface{}) {
	*this = append(*this, v.(uint8))
}

func (this Uint8s) Sum() int64 {
	sum := int64(0)
	for i := range this {
		sum += int64(this[i])
	}
	return sum
}

func (this Uint8s) Size() uint32 {
	return uint32(len(this))
}

func (this Uint8s) Encode() []byte {
	buffer := make([]byte, len(this)*UINT8_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Uint8s) EncodeToBuffer(buffer []byte) {
	for i := range this {
		buffer[i] = uint8(this[i])
	}
}

func (Uint8s) Decode(data []byte) interface{} {
	uint8s := make([]uint8, len(data)/UINT8_LEN)
	for i := range uint8s {
		uint8s[i] = data[i]
	}
	return Uint8s(uint8s)
}
