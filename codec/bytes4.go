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
	"bytes"
	"math"
)

const (
	BYTES4_LEN = 4
)

type Bytes4 [BYTES4_LEN]byte

func NewBytes4(v byte) Bytes4 {
	data := [BYTES4_LEN]byte{}
	for i := range data {
		data[i] = v
	}
	return data
}

func (Bytes4) FromSlice(v []byte) Bytes4 {
	data := [BYTES4_LEN]byte{}
	length := math.Min(float64(BYTES4_LEN), float64(len(v)))
	for i := 0; i < int(length); i++ {
		data[i] = v[i]
	}
	return data
}

func (this *Bytes4) Get() interface{} {
	return *this
}

func (this *Bytes4) Set(v interface{}) {
	*this = v.(Bytes4)
}

func (hash Bytes4) Size() uint32 {
	return uint32(BYTES4_LEN)
}

func (this Bytes4) Clone() interface{} {
	target := Bytes4{}
	copy(target[:], this[:])
	return target
}

func (hash Bytes4) FromBytes(bytes []byte) Bytes4 {
	hash = Bytes4{}
	copy(hash[:], bytes)
	return hash
}

func (hash Bytes4) Encode() []byte {
	return hash[:]
}

func (this Bytes4) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	copy(this[:], buffer)
	return Bytes4(this)
}

type Bytes4s [][4]byte

func (this Bytes4s) Clone() Bytes4s {
	target := make([][BYTES4_LEN]byte, len(this))
	for i := 0; i < len(this); i++ {
		copy(target[i][:], this[i][:])
	}
	return Bytes4s(target)
}

func (this Bytes4s) Encode() []byte {
	return Bytes4s(this).Flatten()
}

func (this Bytes4s) EncodeToBuffer(buffer []byte) int {
	for i := 0; i < len(this); i++ {
		copy(buffer[i*BYTES4_LEN:], this[i][:])
	}
	return len(this) * BYTES4_LEN
}

func (this Bytes4s) Decode(data []byte) interface{} {
	this = make([][4]byte, len(data)/BYTES4_LEN)
	for i := 0; i < len(this); i++ {
		copy(this[i][:], data[i*BYTES4_LEN:(i+1)*BYTES4_LEN])
	}
	return this
}

func (this Bytes4s) Size() uint32 {
	return uint32(len(this) * BYTES4_LEN)
}

func (this Bytes4s) Flatten() []byte {
	buffer := make([]byte, len(this)*BYTES4_LEN)
	for i := 0; i < len(this); i++ {
		copy(buffer[i*BYTES4_LEN:(i+1)*BYTES4_LEN], this[i][:])
	}
	return buffer
}

func (this Bytes4s) Len() int {
	return len(this)
}

func (this Bytes4s) Less(i, j int) bool {
	return bytes.Compare(this[i][:], this[j][:]) < 0
}

func (this Bytes4s) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
