/*
 *   Copyright (c) 2024 Arcology Network

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package codec

import (
	"encoding/hex"
	"sort"
	"unsafe"

	common "github.com/arcology-network/common-lib/common"
)

const (
	CHAR_LEN = 1
)

type String string

// func UUnsafeStringToBytes(s *string) *[]byte {
// 	return (*[]byte)(unsafe.Pointer(s))
// }

// Avoid copying the data.
// func UnsafeBytesToString(b *[]byte) *string {
// 	return (*string)(unsafe.Pointer(b))
// }

func (this String) Clone() interface{} {
	b := make([]byte, len(this))
	copy(b, this)
	return String(*(*string)(unsafe.Pointer(&b)))
}

func (this String) ToBytes() []byte {
	return []byte(this)
}

func (this String) Reverse() string {
	reversed := []byte(this)
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	return *(*string)(unsafe.Pointer(&reversed))
}

func (this String) Sum(offset uint64) uint64 {
	total := uint64(0)
	for j := offset; j < uint64(len(this)); j++ {
		total += uint64(this[j])
	}
	return total
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

func (this String) Size() uint64 {
	return uint64(len(this))
}

func (String) Decode(buffer []byte) interface{} {
	return String(buffer)
}

type Strings []string

func (this Strings) Concate() string {
	return Bytes(this.Flatten()).ToString()
}

func (this Strings) Sort() Strings {
	sort.Slice(this, func(i, j int) bool {
		return this[i] < this[j]
	})
	return this
}

func (this Strings) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Strings) HeaderSize() uint64 {
	if len(this) == 0 {
		return 0
	}
	return uint64(len(this)+1) * UINT64_LEN
}

func (this Strings) Size() uint64 {
	total := uint64(0)
	for i := 0; i < len(this); i++ {
		total += uint64(len(this[i]))
	}
	return this.HeaderSize() + total
}

func (this Strings) FillHeader(buffer []byte) {
	if len(this) == 0 {
		return
	}
	Uint64(len(this)).EncodeToBuffer(buffer)

	offset := 0
	for i := 0; i < len(this); i++ {
		Uint64(offset).EncodeToBuffer(buffer[UINT64_LEN*(i+1):])
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
		copy(buffer[offset:offset+uint64(len(this[i]))], this[i])
		offset += uint64(len(this[i]))
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

func (this Strings) ToHex() []string {
	hexStrings := make([]string, len(this))
	for i := 0; i < len(hexStrings); i++ {
		hexStrings[i] = hex.EncodeToString([]byte(this[i]))
	}
	return hexStrings
}

type Stringset [][]string

func (this Stringset) Size() uint64 {
	length := 0
	for i := 0; i < len(this); i++ {
		length += int(Strings(this[i]).Size())
	}
	return uint64(len(this)+1)*UINT64_LEN + uint64(length)
}

func (this Stringset) Encode() []byte {
	length := int(this.Size())
	buffer := make([]byte, length)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Stringset) EncodeToBuffer(buffer []byte) int {
	lengths := make([]uint64, len(this))
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
