package codec

import (
	"encoding/binary"
)

const (
	INT_LEN = 8
)

type Int64 int64

func (this *Int64) Get() interface{} {
	return *this
}

func (this *Int64) Set(v interface{}) {
	*this = v.(Int64)
}

func (this Int64) Size() int64 {
	return int64(INT_LEN)
}

func (this Int64) Encode() []byte {
	data := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(data, int64(this))
	return data
}

func (Int64) Decode(data []byte) int64 {
	v, _ := binary.Varint(data)
	return v
}

type Int64s []Int64

func (this Int64s) Sum() int64 {
	sum := int64(0)
	for i := range this {
		sum += int64(this[i])
	}
	return sum
}

func (this Int64s) Accumulate() []Int64 {
	if len(this) == 0 {
		return []Int64{}
	}

	values := make([]Int64, len(this))
	values[0] = this[0]
	for i := 1; i < len(this); i++ {
		values[i] = values[i-1] + this[i]
	}
	return values
}
