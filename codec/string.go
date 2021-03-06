package codec

import (
	"reflect"
	"unsafe"

	"github.com/arcology-network/common-lib/common"
)

const (
	CHAR_LEN = 1
)

type String string

func (this String) Clone() String {
	b := make([]byte, len(this))
	copy(b, this)
	return String(*(*string)(unsafe.Pointer(&b)))
}

func (this String) ToBytes() []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&this)).Data),
	)[:len(this):len(this)]
}

func (this String) Encode() []byte {
	return this.ToBytes()
}

func (this String) EncodeToBuffer(buffer []byte) {
	if len(this) > 0 {
		copy(buffer, this.ToBytes())
	}
}

func (this String) Size() uint32 {
	return uint32(len(this))
}

func (String) Decode(bytes []byte) interface{} {
	var s string
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	stringHeader.Data = sliceHeader.Data
	stringHeader.Len = sliceHeader.Len
	return String(s)
}

type Strings []string

func (this Strings) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Strings) HeaderSize() uint32 {
	if len(this) == 0 {
		return 0
	}
	return uint32(len(this)+1) * UINT32_LEN
}

func (this Strings) Size() uint32 {
	total := uint32(0)
	for i := 0; i < len(this); i++ {
		total += uint32(len(this[i]))
	}
	return this.HeaderSize() + total
}

func (this Strings) FillHeader(buffer []byte) {
	if len(this) == 0 {
		return
	}
	Uint32(len(this)).EncodeToBuffer(buffer)

	offset := 0
	for i := 0; i < len(this); i++ {
		Uint32(offset).EncodeToBuffer(buffer[UINT32_LEN*(i+1):])
		offset += len(this[i])
	}
}

func (this Strings) EncodeToBuffer(buffer []byte) {
	if len(buffer) == 0 {
		return
	}
	this.FillHeader(buffer)

	offset := uint32(0)
	headerLen := this.HeaderSize()
	for i := 0; i < len(this); i++ {
		copy(buffer[headerLen+offset:headerLen+offset+uint32(len(this[i]))], this[i])
		offset += uint32(len(this[i]))
	}
}

func (this Strings) Decode(bytes []byte) interface{} {
	if len(bytes) == 0 {
		return Strings{}
	}

	fields := Byteset{}.Decode(bytes).(Byteset)
	if len(bytes) < 1024 {
		return Strings(this.singleThreadDecode(fields))
	}
	return Strings(this.multiThreadDecode(fields))
}

func (Strings) singleThreadDecode(fields [][]byte) []string {
	this := make([]string, len(fields))
	for i := range fields {
		this[i] = string(String("").Decode(fields[i]).(String))
	}
	return this
}

func (Strings) multiThreadDecode(fields [][]byte) []string {
	this := make([]string, len(fields))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			this[i] = string(String("").Decode(fields[i]).(String))
		}
	}
	common.ParallelWorker(len(fields), 4, worker)
	return this
}

func (this Strings) Flatten() []byte {
	positions := make([]int, len(this)+1)
	positions[0] = 0
	for i := 1; i < len(positions); i++ {
		positions[i] = positions[i-1] + len(this[i-1])
	}

	buffer := make([]byte, positions[len(positions)-1])
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			copy(buffer[positions[i]:positions[i+1]], []byte(this[i]))
		}
	}
	common.ParallelWorker(len(this), 4, worker)
	return buffer
}

func (this Strings) Clone() Strings {
	nStrings := make([]string, len(this))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			nStrings[i] = string(String(this[i]).Clone())
		}
	}
	common.ParallelWorker(len(this), 4, worker)
	return nStrings
}

func (this Strings) ToBytes() [][]byte {
	bytes := make([][]byte, len(this))
	for i := 0; i < len(this); i++ {
		bytes[i] = String(this[i]).ToBytes()
	}
	return bytes
}

func (Strings) FromBytes(byteSet [][]byte) []string {
	strings := make([]string, len(byteSet))
	for i := 0; i < len(byteSet); i++ {
		strings[i] = String("").Decode(byteSet[i]).(string)
	}
	return strings
}
