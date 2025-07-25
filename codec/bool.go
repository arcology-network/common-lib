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

import common "github.com/arcology-network/common-lib/common"

const (
	BOOL_LEN = 1
)

type Bool bool

func (this *Bool) Clone() interface{} {
	if this == nil {
		return this
	}

	return common.New(*this)
}

func (this *Bool) Get() interface{} {
	return *this
}

func (this *Bool) Set(v interface{}) {
	*this = v.(Bool)
}

func (this Bool) Size() uint64 {
	return uint64(BOOL_LEN)
}

func (this Bool) Encode() []byte {
	buffer := make([]byte, BOOL_LEN)
	this.EncodeTo(buffer)
	return buffer
}

func (this Bool) EncodeTo(buffer []byte) int {
	buffer[0] = uint8(common.IfThen(bool(this), 1, 0))
	return BOOL_LEN
}

func (this Bool) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	this = Bool(buffer[0] > 0)
	return this
}

type Bools []bool

func (this Bools) Size() int {
	return len(this)
}

func (this Bools) Encode() []byte {
	buffer := make([]byte, len(this))
	this.EncodeTo(buffer)
	return buffer
}

func (this Bools) EncodeTo(buffer []byte) int {
	for i := range this {
		if this[i] {
			buffer[i] = 1
		} else {
			buffer[i] = 0
		}
	}
	return len(this) * BOOL_LEN
}

func (Bools) Decode(data []byte) interface{} {
	bools := make([]bool, len(data))
	for i := range data {
		if data[i] == 1 {
			bools[i] = true
		} else {
			bools[i] = false
		}
	}
	return Bools(bools)
}
