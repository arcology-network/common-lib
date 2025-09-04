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
)

const (
	BYTE16_LEN = 16
)

type Bytes16 [BYTE16_LEN]byte

func NewBytes16(v byte) Bytes16 {
	bytes16 := [BYTE16_LEN]byte{}
	for i := range bytes16 {
		bytes16[i] = v
	}
	return bytes16
}

func (this *Bytes16) Get() any {
	return *this
}

func (this *Bytes16) Set(v any) {
	*this = v.(Bytes16)
}

func (hash Bytes16) Size() uint64 {
	return uint64(BYTE16_LEN)
}

func (this Bytes16) Clone() any {
	target := Bytes16{}
	copy(target[:], this[:])
	return target
}

func (hash Bytes16) FromBytes(bytes []byte) Bytes16 {
	hash = Bytes16{}
	copy(hash[:], bytes)
	return hash
}

func (hash Bytes16) Encode() []byte {
	return hash[:]
}

func (this Bytes16) Decode(buffer []byte) any {
	if len(buffer) == 0 {
		return this
	}

	copy(this[:], buffer)
	return Bytes16(this)
}

type Bytes16s [][16]byte

func (this Bytes16s) Clone() Bytes16s {
	target := make([][BYTE16_LEN]byte, len(this))
	for i := 0; i < len(this); i++ {
		copy(target[i][:], this[i][:])
	}
	return Bytes16s(target)
}

func (this Bytes16s) Encode() []byte {
	return Bytes16s(this).Flatten()
}

func (this Bytes16s) EncodeTo(buffer []byte) int {
	for i := 0; i < len(this); i++ {
		copy(buffer[i*BYTE16_LEN:], this[i][:])
	}
	return len(this) * BYTE16_LEN
}

func (this Bytes16s) Decode(data []byte) any {
	this = make([][16]byte, len(data)/BYTE16_LEN)
	for i := 0; i < len(this); i++ {
		copy(this[i][:], data[i*BYTE16_LEN:(i+1)*BYTE16_LEN])
	}
	return this
}

func (this Bytes16s) Size() uint64 {
	return uint64(len(this) * BYTE16_LEN)
}

func (this Bytes16s) Flatten() []byte {
	buffer := make([]byte, len(this)*BYTE16_LEN)
	for i := 0; i < len(this); i++ {
		copy(buffer[i*BYTE16_LEN:(i+1)*BYTE16_LEN], this[i][:])
	}
	return buffer
}

func (this Bytes16s) Len() int {
	return len(this)
}

func (this Bytes16s) Less(i, j int) bool {
	return bytes.Compare(this[i][:], this[j][:]) < 0
}

func (this Bytes16s) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
