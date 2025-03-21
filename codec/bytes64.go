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
	HASH64_LEN = 64
)

type Hash64 [HASH64_LEN]byte

func (hash Hash64) Size() uint64 {
	return uint64(HASH64_LEN)
}

func (this Hash64) Sum(offset uint64) uint64 {
	total := uint64(0)
	for j := offset; j < uint64(len(this)); j++ {
		total += uint64((this)[j])
	}
	return total
}

func (this Hash64) Clone() interface{} {
	target := Hash64{}
	copy(target[:], this[:])
	return target
}

func (hash Hash64) Encode() []byte {
	return hash[:]
}

func (hash Hash64) Decode(data []byte) interface{} {
	copy(hash[:], data)
	return Hash64(hash)
}

type Hash64s [][HASH64_LEN]byte

func (this Hash64s) Clone() Hash64s {
	target := make([][HASH64_LEN]byte, len(this))
	for i := 0; i < len(this); i++ {
		copy(target[i][:], this[i][:])
	}
	return Hash64s(target)
}

func (hashes Hash64s) Encode() []byte {
	return Hash64s(hashes).Flatten()
}

func (this Hash64s) EncodeToBuffer(buffer []byte) int {
	for i := 0; i < len(this); i++ {
		copy(buffer[i*HASH64_LEN:], this[i][:])
	}
	return len(this) * HASH64_LEN
}

func (this Hash64s) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	this = make([][HASH64_LEN]byte, len(buffer)/HASH64_LEN)
	for i := 0; i < len(this); i++ {
		copy(this[i][:], buffer[i*HASH64_LEN:(i+1)*HASH64_LEN])
	}
	return Hash64s(this)
}

func (hashes Hash64s) Size() uint64 {
	return uint64(len(hashes) * HASH64_LEN)
}

func (hashes Hash64s) Flatten() []byte {
	buffer := make([]byte, len(hashes)*HASH64_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(buffer[i*HASH64_LEN:(i+1)*HASH64_LEN], hashes[i][:])
	}
	return buffer
}

func (hashes Hash64s) Len() int {
	return len(hashes)
}

func (hashes Hash64s) Less(i, j int) bool {
	return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
}

func (hashes Hash64s) Swap(i, j int) {
	hashes[i], hashes[j] = hashes[j], hashes[i]
}
