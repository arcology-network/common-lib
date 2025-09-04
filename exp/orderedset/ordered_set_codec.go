/*
 *   Copyright (c) 2025 Arcology Network

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

package orderedset

import (
	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/exp/slice"
)

func (this *OrderedSet[K]) Size() int {
	if len(this.elements) == 0 {
		return 0 // Header size only
	}

	size := (len(this.elements)+1)*codec.UINT64_LEN +
		slice.Accumulate(this.elements, 0, func(_ int, k K) int { return this.Sizer(k) }) // Header size
	return size
}

func (this *OrderedSet[K]) EncodeTo(buf []byte) int {
	if this.Size() == 0 {
		return 0
	}

	lengths := slice.Transform(this.elements, func(_ int, k K) uint64 { return uint64(this.Sizer(k)) })
	offset := codec.Encoder{}.FillHeader(
		buf,
		lengths,
	)

	for _, k := range this.elements {
		offset += this.Encoder(k, buf[offset:])
	}
	return offset
}

func (this *OrderedSet[K]) Decode(buf []byte) any {
	fields := codec.Byteset{}.Decode(buf).(codec.Byteset)
	elements := make([]K, len(fields))
	for i := 0; i < len(fields); i++ {
		elements[i] = this.Decoder(fields[i])
	}

	this.elements = elements
	this.Init()
	return this
}

func (this *OrderedSet[K]) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeTo(buffer)
	return buffer
}
