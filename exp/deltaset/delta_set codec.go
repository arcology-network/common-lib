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

package deltaset

import (
	"github.com/arcology-network/common-lib/codec"
	orderedset "github.com/arcology-network/common-lib/exp/orderedset"
)

func (this *DeltaSet[K]) Size() int {
	return 4*codec.UINT64_LEN +
		this.committed.Size() +
		this.stagedAdditions.Size() +
		this.stagedRemovals.Size()
}

func (this *DeltaSet[K]) EncodeTo(buf []byte) int {
	offset := codec.Encoder{}.FillHeader(buf, []uint64{
		uint64(this.committed.Size()),
		uint64(this.stagedAdditions.Size()),
		uint64(this.stagedRemovals.Size()),
	})

	offset += this.committed.EncodeTo(buf[offset:]) // allDeleted
	offset += this.stagedAdditions.EncodeTo(buf[offset:])
	offset += this.stagedRemovals.EncodeTo(buf[offset:])
	return offset
}

func (this *DeltaSet[K]) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeTo(buffer)
	return buffer
}

func (this *DeltaSet[K]) Decode(buffer []byte) *DeltaSet[K] {
	fields := codec.Byteset{}.Decode(buffer).(codec.Byteset) // Decode header
	return &DeltaSet[K]{
		committed:       this.committed.Decode(fields[0]).(*orderedset.OrderedSet[K]),
		stagedAdditions: this.stagedAdditions.Decode(fields[1]).(*orderedset.OrderedSet[K]),
		stagedRemovals:  this.stagedRemovals.Decode(fields[2]).(*orderedset.OrderedSet[K]),
	}
}
