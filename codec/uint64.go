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
	"crypto/sha256"
	"encoding/binary"
	"sort"
	"unsafe"

	common "github.com/arcology-network/common-lib/common"
)

const (
	UINT64_LEN = 8
)

type Uint64 uint64

func (this *Uint64) Clone() any {
	if this == nil {
		return this
	}
	return common.New(*this)
}

func (this *Uint64) Get() any {
	return *this
}

func (this *Uint64) Set(v any) {
	*this = v.(Uint64)
}

func (Uint64) Size() uint64 {
	return UINT64_LEN
}

func (this Uint64) Encode() []byte {
	buffer := make([]byte, UINT64_LEN)
	this.EncodeTo(buffer)
	return buffer
}

func (this Uint64) EncodeTo(buffer []byte) int {
	binary.LittleEndian.PutUint64(buffer, uint64(this))
	return UINT64_LEN
}

func (this Uint64) Decode(data []byte) any {
	if len(data) == 0 {
		return this
	}

	this = Uint64(binary.LittleEndian.Uint64(data))
	return Uint64(this)
}

func (v Uint64) Checksum() [32]byte {
	return sha256.Sum256(v.Encode())
}

func (v Uint64) ToInt64() int64 {
	return *(*int64)(unsafe.Pointer(&v))
}

type Uint64s []uint64

func (this Uint64s) Get() any {
	return this.Sum()
}

func (this Uint64s) Set(v any) {
	this = append(this, v.(uint64))
}

func (this Uint64s) Sum() int64 {
	sum := int64(0)
	for i := range this {
		sum += int64(this[i])
	}
	return sum
}

func (this Uint64s) Accumulate() []uint64 {
	if len(this) == 0 {
		return []uint64{}
	}

	values := make([]uint64, len(this))
	values[0] = this[0]
	for i := 1; i < len(this); i++ {
		values[i] = values[i-1] + this[i]
	}
	return values
}

func (this Uint64s) Unique() []uint64 {
	sort.SliceStable(this, func(i, j int) bool {
		return this[i] < this[j]
	})

	uniqueV := make([]uint64, 0, len(this))
	current := uint64(this[0])
	for i := 0; i < len(this); i++ {
		if current != uint64(this[i]) {
			uniqueV = append(uniqueV, current)
			current = uint64(this[i])
		}
	}

	if current != uniqueV[len(uniqueV)-1] {
		uniqueV = append(uniqueV, current)
	}

	return uniqueV
}

func (this Uint64s) Encode() []byte {
	buffer := make([]byte, len(this)*UINT64_LEN)
	this.EncodeTo(buffer)
	return buffer
}

func (this Uint64s) EncodeTo(buffer []byte) int {
	offset := 0
	for i := range this {
		offset += Uint64(this[i]).EncodeTo(buffer[offset:])
	}
	return len(this) * UINT64_LEN
}

func (this Uint64s) Decode(buffer []byte) any {
	if len(buffer) == 0 {
		return this
	}

	this = make([]uint64, len(buffer)/UINT64_LEN)
	for i := range this {
		this[i] = uint64(Uint64(this[i]).Decode(buffer[i*UINT64_LEN : (i+1)*UINT64_LEN]).(Uint64))
	}
	return Uint64s(this)
}
