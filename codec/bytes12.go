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
	BYTES12_LEN = 12
)

type Bytes12 [BYTES12_LEN]byte

func NewBytes12(v byte) Bytes12 {
	hash16 := [BYTES12_LEN]byte{}
	for i := range hash16 {
		hash16[i] = v
	}
	return hash16
}

func (this *Bytes12) Get() interface{} {
	return *this
}

func (this *Bytes12) Set(v interface{}) {
	*this = v.(Bytes12)
}

func (hash Bytes12) Size() uint64 {
	return uint64(BYTES12_LEN)
}

func (this Bytes12) Clone() interface{} {
	target := Bytes12{}
	copy(target[:], this[:])
	return target
}

func (hash Bytes12) FromBytes(bytes []byte) Bytes12 {
	hash = Bytes12{}
	copy(hash[:], bytes)
	return hash
}

func (hash Bytes12) Encode() []byte {
	return hash[:]
}

func (this Bytes12) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	copy(this[:], buffer)
	return Bytes12(this)
}

type Bytes12s [][12]byte

func (this Bytes12s) Clone() Bytes12s {
	target := make([][BYTES12_LEN]byte, len(this))
	for i := 0; i < len(this); i++ {
		copy(target[i][:], this[i][:])
	}
	return Bytes12s(target)
}

func (this Bytes12s) Encode() []byte {
	return Bytes12s(this).Flatten()
}

func (this Bytes12s) EncodeToBuffer(buffer []byte) int {
	for i := 0; i < len(this); i++ {
		copy(buffer[i*BYTES12_LEN:], this[i][:])
	}
	return len(this) * BYTES12_LEN
}

func (this Bytes12s) Decode(data []byte) interface{} {
	this = make([][12]byte, len(data)/BYTES12_LEN)
	for i := 0; i < len(this); i++ {
		copy(this[i][:], data[i*BYTES12_LEN:(i+1)*BYTES12_LEN])
	}
	return this
}

func (this Bytes12s) Size() uint64 {
	return uint64(len(this) * BYTES12_LEN)
}

func (this Bytes12s) Flatten() []byte {
	buffer := make([]byte, len(this)*BYTES12_LEN)
	for i := 0; i < len(this); i++ {
		copy(buffer[i*BYTES12_LEN:(i+1)*BYTES12_LEN], this[i][:])
	}
	return buffer
}

func (this Bytes12s) Len() int {
	return len(this)
}

func (this Bytes12s) Less(i, j int) bool {
	return bytes.Compare(this[i][:], this[j][:]) < 0
}

func (this Bytes12s) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
