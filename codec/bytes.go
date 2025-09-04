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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"unsafe"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	BYTE_LEN = 1
)

type Bytes []byte

func (*Bytes) LessAsUint64(first, second []byte) bool {
	return *(*uint64)(unsafe.Pointer((*[8]byte)(unsafe.Pointer(&first)))) <
		*(*uint64)(unsafe.Pointer((*[8]byte)(unsafe.Pointer(&second))))
}

func (this *Bytes) Get() interface{} {
	return *this
}

func (this *Bytes) Set(v interface{}) {
	*this = v.(Bytes)
}

func (this *Bytes) Sum(offset uint64) uint64 {
	total := uint64(0)
	for j := offset; j < uint64(len(*this)); j++ {
		total += uint64((*this)[j])
	}
	return total
}

func (this *Bytes) Hex() string {
	bytes := make([]byte, 2*len(*this))
	hex.Encode(bytes[:], (*this)[:])
	return string(bytes)
}

func (this Bytes) Encode() []byte {
	return []byte(this)
}

func (this Bytes) Size() uint64 {
	return uint64(len(this))
}

func (this Bytes) Clone() interface{} {
	if this == nil {
		return this
	}

	target := make([]byte, len(this))
	copy(target, this)
	return Bytes(target)
}

func (this Bytes) EncodeTo(buffer []byte) int {
	copy(buffer, this)
	return len(this)
}

func (this Bytes) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}
	return Bytes(buffer)
}

func (this Bytes) ToString() string {
	return *(*string)(unsafe.Pointer(&this))
}

type Byteset [][]byte

func (this Byteset) Clone() interface{} {
	if this == nil {
		return this
	}

	target := make([][]byte, len(this))
	for i := range this {
		target[i] = make([]byte, len(this[i]))
		copy(target[i], this[i])
	}
	return Byteset(target)
}

func (this Byteset) Size() uint64 {
	if len(this) == 0 {
		return 0
	}

	total := (len(this) + 1) * UINT64_LEN // Header size
	for i := 0; i < len(this); i++ {
		total += len(this[i])
	}
	return uint64(total)
}

func (this Byteset) Sizes() Uint64s {
	sizes := make([]uint64, len(this))
	for i := range this {
		sizes[i] = uint64(len(this[i]))
	}
	return sizes
}

func (this Byteset) Flatten() []byte {
	total := 0
	for i := range this {
		total += len(this[i])
	}
	buffer := make([]byte, total)

	offset := 0
	for i := 0; i < len(this); i++ {
		copy(buffer[offset:], this[i])
		offset += len(this[i])
	}
	return buffer
}

func (this Byteset) Equal(other Byteset) bool {
	if len(this) != len(other) {
		return false
	}

	for i := range this {
		if len(this[i]) != len(other[i]) {
			return false
		}
		if !bytes.Equal(this[i], other[i]) {
			return false
		}
	}
	return true
}

func (this Byteset) Checksum() ethCommon.Hash {
	return sha256.Sum256(this.Flatten())
}

func (this Byteset) Hash(hasher func([]byte) []byte) []byte {
	return hasher(this.Flatten())
}

func (this Byteset) Encode() []byte {
	total := this.Size()
	buffer := make([]byte, total)
	this.EncodeTo(buffer)
	return buffer
}

func (this Byteset) HeaderSize() uint64 {
	if len(this) == 0 {
		return 0
	}
	return uint64(len(this)+1) * UINT64_LEN
}

func (this Byteset) FillHeader(buffer []byte) {
	if len(this) == 0 {
		return
	}

	Uint64(len(this)).EncodeTo(buffer)

	offset := uint64(0)
	for i := 0; i < len(this); i++ {
		Uint64(offset).EncodeTo(buffer[(i+1)*UINT64_LEN:])
		offset += uint64(len(this[i]))
	}
}

func (this Byteset) EncodeTo(buffer []byte) int {
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

func (this Byteset) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return Byteset{}
	}

	count := uint64(Uint64(0).Decode(buffer[:UINT64_LEN]).(Uint64))
	this = make([][]byte, count)

	headerLen := (count + 1) * UINT64_LEN
	prev := uint64(Uint64(0).Decode(buffer[UINT64_LEN : UINT64_LEN+UINT64_LEN]).(Uint64))
	next := uint64(0)
	for i := 0; i < int(count); i++ {
		if i == int(count)-1 {
			next = uint64(len(buffer)) - headerLen
		} else {
			next = uint64(Uint64(0).Decode(buffer[UINT64_LEN+(i+1)*UINT64_LEN : UINT64_LEN+(i+2)*UINT64_LEN]).(Uint64))
		}

		this[i] = buffer[headerLen+prev : headerLen+next]
		prev = next
	}
	return Byteset(this)
}

func (this Byteset) Print() {
	for i, b := range this {
		fmt.Printf("Byteset[%d]: %s\n", i, b)
	}
	fmt.Println()
}

type Bytegroup [][][]byte

func (this Bytegroup) Clone() Bytegroup {
	target := make([][][]byte, len(this))
	for i := range this {
		target[i] = make([][]byte, len(this[i]))
		for j := range this[i] {
			target[i][j] = make([]byte, len(this[i][j]))
			copy(target[i][j], this[i][j])
		}
	}
	return Bytegroup(target)
}

func (bytegroup Bytegroup) Sizes() []uint64 {
	sizes := make([]uint64, len(bytegroup))
	for i := range bytegroup {
		sizes[i] = uint64(len(bytegroup[i]))
	}
	return sizes
}

func (bytegroup Bytegroup) Flatten() [][]byte {
	lengths := bytegroup.Sizes()
	buffer := make([][]byte, Uint64s(lengths).Sum())

	positions := append([]uint64{0}, Uint64s(lengths).Accumulate()...)
	for i := 0; i < len(positions)-1; i++ {
		copy(buffer[positions[i]:positions[i+1]], bytegroup[i])
	}
	return buffer
}
