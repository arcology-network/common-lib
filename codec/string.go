package codec

import (
	"reflect"
	"sort"
	"unsafe"

	common "github.com/arcology-network/common-lib/common"
)

const (
	CHAR_LEN = 1
)

type String string

func (this String) Clone() interface{} {
	b := make([]byte, len(this))
	copy(b, this)
	return String(*(*string)(unsafe.Pointer(&b)))
}

func (this String) ToBytes() []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&this)).Data),
	)[:len(this):len(this)]
}

func (this String) Reverse() string {
	reversed := []byte(this)
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	return *(*string)(unsafe.Pointer(&reversed))
}

func (this String) Encode() []byte {
	return this.ToBytes()
}

func (this String) EncodeToBuffer(buffer []byte) int {
	if len(this) > 0 {
		copy(buffer, this.ToBytes())
	}
	return len(this) * CHAR_LEN
}

func (this String) Size() uint32 {
	return uint32(len(this))
}

func (String) Decode(buffer []byte) interface{} {
	return String(buffer)
}

type Strings []string

func (this Strings) Concate() string {
	return Bytes(this.Flatten()).ToString()
}

func (this Strings) Sort() {
	sort.SliceStable(this, func(i, j int) bool {
		return this[i] < this[j]
	})
}

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

func (this Strings) EncodeToBuffer(buffer []byte) int {
	if len(buffer) == 0 {
		return 0
	}
	this.FillHeader(buffer)

	offset := this.HeaderSize()
	for i := 0; i < len(this); i++ {
		copy(buffer[offset:offset+uint32(len(this[i]))], this[i])
		offset += uint32(len(this[i]))
	}
	return int(offset)
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
			nStrings[i] = string(String(this[i]).Clone().(String))
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

type Stringset [][]string

func (this Stringset) Size() uint32 {
	length := 0
	for i := 0; i < len(this); i++ {
		length += int(Strings(this[i]).Size())
	}
	return uint32(len(this)+1)*UINT32_LEN + uint32(length)
}

func (this Stringset) Encode() []byte {
	length := int(this.Size())
	buffer := make([]byte, length)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Stringset) EncodeToBuffer(buffer []byte) int {
	lengths := make([]uint32, len(this))
	for i := 0; i < len(this); i++ {
		lengths[i] = Strings(this[i]).Size()
	}

	offset := Encoder{}.FillHeader(buffer, lengths)
	for i := 0; i < len(this); i++ {
		offset += Strings(this[i]).EncodeToBuffer(buffer[offset:])
	}
	return offset
}

func (this Stringset) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	fields := Byteset{}.Decode(buffer).(Byteset)

	stringset := make([][]string, len(fields))
	for i := 0; i < len(fields); i++ {
		stringset[i] = []string(Strings{}.Decode(fields[i]).(Strings))
	}
	return Stringset(stringset)
}

func (this Stringset) Flatten() []string {
	positions := make([]int, len(this)+1)
	positions[0] = 0
	for i := 1; i < len(positions); i++ {
		positions[i] = positions[i-1] + len(this[i-1])
	}

	buffer := make([]string, positions[len(positions)-1])
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			copy(buffer[positions[i]:positions[i+1]], (this[i]))
		}
	}
	common.ParallelWorker(len(this), 4, worker)
	return buffer
}
