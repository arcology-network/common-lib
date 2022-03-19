package codec

import (
	"encoding/binary"
)

const (
	UINT32_LEN = 4
)

type Uint32 uint32

func (this *Uint32) Get() interface{} {
	return *this
}

func (this *Uint32) Set(v interface{}) {
	*this = v.(Uint32)
}

func (Uint32) Size() uint32 {
	return UINT32_LEN
}

func (this Uint32) Encode() []byte {
	buffer := make([]byte, UINT32_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Uint32) EncodeToBuffer(buffer []byte) {
	binary.LittleEndian.PutUint32(buffer, uint32(this))
}

func (this Uint32) Decode(buffer []byte) interface{} {
	this = Uint32(binary.LittleEndian.Uint32(buffer))
	return Uint32(this)
}

type Uint32s []uint32

func (this Uint32s) Encode() []byte {
	buffer := make([]byte, uint32(len(this)*UINT32_LEN))
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Uint32s) EncodeToBuffer(buffer []byte) {
	for i := range this {
		copy(buffer[i*UINT32_LEN:(i+1)*UINT32_LEN], Uint32(this[i]).Encode())
	}
}

func (this Uint32s) Decode(data []byte) interface{} {
	this = make([]uint32, len(data)/UINT32_LEN)
	for i := range this {
		this[i] = uint32(Uint32(this[i]).Decode(data[i*UINT32_LEN : (i+1)*UINT32_LEN]).(Uint32))
	}
	return Uint32s(this)
}

func (this Uint32s) Accumulate() []uint32 {
	if len(this) == 0 {
		return []uint32{}
	}

	values := make([]uint32, len(this))
	values[0] = this[0]
	for i := 1; i < len(this); i++ {
		values[i] = values[i-1] + this[i]
	}
	return values
}

func (this Uint32s) Sum() uint32 {
	sum := uint32(0)
	for i := range this {
		sum += uint32(this[i])
	}
	return sum
}
