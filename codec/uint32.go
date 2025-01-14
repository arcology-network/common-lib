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
	"encoding/binary"

	common "github.com/arcology-network/common-lib/common"
)

const (
	UINT32_LEN = uint64(4)
)

type Uint32 uint32

func (this *Uint32) Clone() interface{} {
	if this == nil {
		return this
	}

	return common.New(*this)
}

func (this *Uint32) Get() interface{} {
	return *this
}

func (this *Uint32) Set(v interface{}) {
	*this = v.(Uint32)
}

func (Uint32) Size() uint64 {
	return UINT64_LEN
}

func (this Uint32) Encode() []byte {
	buffer := make([]byte, UINT64_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Uint32) EncodeToBuffer(buffer []byte) int {
	binary.LittleEndian.PutUint32(buffer, uint32(this))
	return UINT64_LEN
}

func (this Uint32) Decode(buffer []byte) interface{} {
	this = Uint32(binary.LittleEndian.Uint32(buffer))
	return Uint32(this)
}

type Uint32s []uint32

func (this Uint32s) Encode() []byte {
	buffer := make([]byte, uint64(len(this)*UINT64_LEN))
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Uint32s) EncodeToBuffer(buffer []byte) int {
	offset := 0
	for i := range this {
		offset += Uint32(this[i]).EncodeToBuffer(buffer[offset:])
	}
	return len(this) * UINT64_LEN
}

func (this Uint32s) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	this = make([]uint32, len(buffer)/UINT64_LEN)
	for i := range this {
		this[i] = uint32(Uint32(this[i]).Decode(buffer[i*UINT64_LEN : (i+1)*UINT64_LEN]).(Uint32))
	}
	return Uint32s(this)
}

func (this Uint32s) Accumulate() []uint64 {
	if len(this) == 0 {
		return []uint64{}
	}

	values := make([]uint64, len(this))
	values[0] = uint64(this[0])
	for i := 1; i < len(this); i++ {
		values[i] = values[i-1] + uint64(this[i])
	}
	return values
}

func (this Uint32s) Sum() uint64 {
	sum := uint64(0)
	for i := range this {
		sum += uint64(this[i])
	}
	return sum
}
