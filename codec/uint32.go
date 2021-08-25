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

func (v Uint32) Size() int {
	return UINT32_LEN
}

func (v Uint32) Encode() []byte {
	data := make([]byte, UINT32_LEN)
	binary.LittleEndian.PutUint32(data[0:], uint32(v))
	return data
}

func (Uint32) Decode(data []byte) uint32 {
	return uint32(binary.LittleEndian.Uint32(data[0:UINT32_LEN]))
}

type Uint32s []Uint32

func (uint32s Uint32s) Encode() []byte {
	buffer := make([]byte, uint32(len(uint32s)*UINT32_LEN))
	for i := range uint32s {
		copy(buffer[i*UINT32_LEN:(i+1)*UINT32_LEN], Uint32(uint32s[i]).Encode())
	}
	return buffer
}

func (uint32s Uint32s) Decode(data []byte) interface{} {
	uint32s = make([]Uint32, len(data)/UINT32_LEN)
	for i := range uint32s {
		uint32s[i] = Uint32(uint32s[i].Decode(data[i*UINT32_LEN : (i+1)*UINT32_LEN]))
	}
	return Uint32s(uint32s)
}

func (uint32s Uint32s) Accumulate() []Uint32 {
	if len(uint32s) == 0 {
		return []Uint32{}
	}

	values := make([]Uint32, len(uint32s))
	values[0] = uint32s[0]
	for i := 1; i < len(uint32s); i++ {
		values[i] = values[i-1] + uint32s[i]
	}
	return values
}

func (uint32s Uint32s) Sum() uint32 {
	sum := uint32(0)
	for i := range uint32s {
		sum += uint32(uint32s[i])
	}
	return sum
}
