package codec

import (
	"encoding/binary"
	"unsafe"

	common "github.com/arcology-network/common-lib/common"
)

const (
	INT64_LEN = 8
)

type Int64 int64

func (this *Int64) Clone() interface{} {
	return common.New(*this)
}

func (this *Int64) Get() interface{} {
	return *this
}

func (this *Int64) Set(v interface{}) {
	*this = v.(Int64)
}

func (this Int64) Size() uint32 {
	return uint32(INT64_LEN)
}

func (this Int64) Encode() []byte {
	buffer := make([]byte, INT64_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Int64) EncodeToBuffer(buffer []byte) int {
	binary.LittleEndian.PutUint64(buffer, uint64(this))
	return INT64_LEN
}

func (this Int64) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	return Int64(int64(binary.LittleEndian.Uint64(buffer)))
}

func (this Int64) ToUint64(src1 int64) uint64 {
	return *(*uint64)(unsafe.Pointer(&src1))
}

type Int64s []Int64

func (this Int64s) Encode() []byte {
	buffer := make([]byte, len(this)*INT64_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Int64s) EncodeToBuffer(buffer []byte) int {
	for i := 0; i < len(this); i++ {
		binary.LittleEndian.PutUint64(buffer[i*INT64_LEN:], uint64(this[i]))
	}
	return len(this) * INT64_LEN
}

func (this Int64s) Decode(buffer []byte) Int64s {
	for i := 0; i < len(this); i++ {
		this[i] = Int64(int64(binary.LittleEndian.Uint64(buffer)))
	}
	return Int64s(this)
}

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
