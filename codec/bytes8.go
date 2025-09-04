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

	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	HASH8_LEN = 8
)

type Hash8 [HASH8_LEN]byte

func NewHash8(v byte) Hash8 {
	hash8 := [HASH8_LEN]byte{}
	for i := range hash8 {
		hash8[i] = v
	}
	return hash8
}

func (this Hash8) Clone() any {
	target := Hash8{}
	copy(target[:], this[:])
	return target
}

func (this *Hash8) Get() any {
	return *this
}

func (this *Hash8) Set(v any) {
	*this = v.(Hash8)
}

func (hash Hash8) Size() uint64 {
	return uint64(HASH8_LEN)
}

func (hash Hash8) FromBytes(bytes []byte) Hash8 {
	hash = Hash8{}
	copy(hash[:], bytes)
	return hash
}

func (this Hash8) Sum(offset uint64) uint64 {
	total := uint64(0)
	for j := offset; j < uint64(len(this)); j++ {
		total += uint64(this[j])
	}
	return total
}

func (hash Hash8) Encode() []byte {
	return hash[:]
}

func (this Hash8) Decode(buffer []byte) any {
	if len(buffer) == 0 {
		return this
	}

	copy(this[:], buffer)
	return Hash8(this)
}

type Hash8s []ethCommon.Hash

func (hashes Hash8s) Encode() []byte {
	return Hash8s(hashes).Flatten()
}

func (hashes Hash8s) Decode(data []byte) any {
	hashes = make([]ethCommon.Hash, len(data)/HASH8_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(hashes[i][:], data[i*HASH8_LEN:(i+1)*HASH8_LEN])
	}
	return hashes
}

func (hashes Hash8s) Size() uint64 {
	return uint64(len(hashes) * HASH8_LEN)
}

func (hashes Hash8s) Flatten() []byte {
	buffer := make([]byte, len(hashes)*HASH8_LEN)
	for i := 0; i < len(hashes); i++ {
		copy(buffer[i*HASH8_LEN:(i+1)*HASH8_LEN], hashes[i][:])
	}
	return buffer
}

func (hashes Hash8s) Len() int {
	return len(hashes)
}

func (hashes Hash8s) Less(i, j int) bool {
	return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
}

func (hashes Hash8s) Swap(i, j int) {
	hashes[i], hashes[j] = hashes[j], hashes[i]
}
