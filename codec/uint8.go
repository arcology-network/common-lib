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

	common "github.com/arcology-network/common-lib/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	UINT8_LEN = 1
)

type Uint8 uint8

func (this *Uint8) Clone() any {
	if this == nil {
		return this
	}

	return common.New(*this)
}

func (this *Uint8) Get() any {
	return *this
}

func (this *Uint8) Set(v any) {
	*this = v.(Uint8)
}

func (v Uint8) Size() uint64 {
	return UINT8_LEN
}

func (v Uint8) Encode() []byte {
	buffer := make([]byte, UINT8_LEN)
	buffer[0] = uint8(v)
	return buffer
}

func (v Uint8) EncodeTo(buffer []byte) int {
	buffer[0] = uint8(v)
	return UINT8_LEN
}

func (this Uint8) Decode(data []byte) any {
	this = Uint8(data[0])
	return this
}

func (v Uint8) Checksum() ethCommon.Hash {
	return sha256.Sum256(v.Encode())
}

type Uint8s []uint8

func (this Uint8s) Get() any {
	return this.Sum()
}

func (this *Uint8s) Set(v any) {
	*this = append(*this, v.(uint8))
}

func (this Uint8s) Sum() int64 {
	sum := int64(0)
	for i := range this {
		sum += int64(this[i])
	}
	return sum
}

func (this Uint8s) Size() uint64 {
	return uint64(len(this))
}

func (this Uint8s) Encode() []byte {
	buffer := make([]byte, len(this)*UINT8_LEN)
	this.EncodeTo(buffer)
	return buffer
}

func (this Uint8s) EncodeTo(buffer []byte) int {
	for i := range this {
		buffer[i] = uint8(this[i])
	}
	return len(this) * UINT8_LEN
}

func (this Uint8s) Decode(buffer []byte) any {
	if len(buffer) == 0 {
		return this
	}

	uint8s := make([]uint8, len(buffer)/UINT8_LEN)
	for i := range uint8s {
		uint8s[i] = buffer[i]
	}
	return Uint8s(uint8s)
}
