package encoding

import (
	"encoding/binary"
)

const (
	UINT32_LEN = 4
)

type Uint32 uint32

func (v Uint32) Size() int {
	return UINT32_LEN
}

func (v Uint32) Encode() []byte {
	data := make([]byte, UINT32_LEN)
	binary.LittleEndian.PutUint32(data[0:], uint32(v))
	return data
}

func (_ Uint32) Decode(data []byte) uint32 {
	return uint32(binary.LittleEndian.Uint32(data[0:UINT32_LEN]))
}

type Uint32s []uint32

func (uint32s Uint32s) Encode() []byte {
	buffer := make([]byte, len(uint32s)*UINT32_LEN)
	for i := range uint32s {
		binary.LittleEndian.PutUint32(buffer[i*UINT32_LEN:(i+1)*UINT32_LEN], uint32s[i])
	}
	return buffer
}

func (uint32s Uint32s) Decode(data []byte) []uint32 {
	uint32s = make([]uint32, len(data)/UINT32_LEN)
	for i := range uint32s {
		uint32s[i] = binary.LittleEndian.Uint32(data[i*UINT32_LEN : (i+1)*UINT32_LEN])
	}
	return uint32s
}

func (uint32s Uint32s) Accumulate() []uint32 {
	if len(uint32s) == 0 {
		return []uint32{}
	}

	values := make([]uint32, len(uint32s))
	values[0] = uint32s[0]
	for i := 1; i < len(uint32s); i++ {
		values[i] = values[i-1] + uint32s[i]
	}
	return values
}

func (uint32s Uint32s) Sum() uint32 {
	sum := uint32(0)
	for i := range uint32s {
		sum += uint32s[i]
	}
	return sum
}
