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

package softdeltaset

import (
	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	orderedset "github.com/arcology-network/common-lib/exp/orderedset"
	"github.com/ethereum/go-ethereum/rlp"
)

func (this *DeltaSet[K]) Size() int {
	// return 4*codec.UINT64_LEN +
	// 	this.committed.Size() +
	// 	this.stagedAdditions.Size() +
	// 	this.stagedRemovals.Size()

	return 4*codec.UINT64_LEN + common.IfThenDo1st(this.committed != nil, func() int { return this.committed.Size() }, 0) +
		common.IfThenDo1st(this.stagedAdditions != nil, func() int { return this.stagedAdditions.Size() }, 0) +
		common.IfThenDo1st(this.stagedRemovals != nil, func() int { return int(this.stagedRemovals.Size()) }, 0)

}

func (this *DeltaSet[K]) EncodeTo(buffer []byte) int {
	// Some sets may be empty to save space, so we only encode the sizes of non-empty sets.
	offset := codec.Encoder{}.FillHeader(buffer,
		[]uint64{
			common.IfThenDo1st(this.committed != nil, func() uint64 { return uint64(this.committed.Size()) }, 0),
			common.IfThenDo1st(this.stagedAdditions != nil, func() uint64 { return uint64(this.stagedAdditions.Size()) }, 0),
			common.IfThenDo1st(this.stagedRemovals != nil, func() uint64 { return uint64(this.stagedRemovals.Size()) }, 0),
		},
	)

	offset += common.IfThenDo1st(this.committed != nil, func() int {
		return this.committed.EncodeTo(buffer[offset:])
	}, 0)

	offset += common.IfThenDo1st(this.stagedAdditions != nil, func() int {
		return this.stagedAdditions.EncodeTo(buffer[offset:])
	}, 0)

	offset += common.IfThenDo1st(this.stagedRemovals != nil, func() int {
		return this.stagedRemovals.EncodeTo(buffer[offset:])
	}, 0)
	return offset
}

func (this *DeltaSet[K]) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeTo(buffer)
	return buffer
}

func (this *DeltaSet[K]) Decode(buffer []byte) any {
	fields := codec.Byteset{}.Decode(buffer).(codec.Byteset) // Decode header
	this.committed = this.committed.Decode(fields[0]).(*orderedset.OrderedSet[K])
	this.stagedAdditions = this.stagedAdditions.Decode(fields[1]).(*orderedset.OrderedSet[K])
	this.stagedRemovals = this.stagedRemovals.Decode(fields[2]).(*StagedRemovalSet[K])

	// When the StagedAddedDeleted flag is set, it means that all elements in the committed set
	// are considered deleted. But the stagedRemovals does not have the
	// committed elements for saving space, so we need to set the committed elements
	// to the stagedRemovals.
	if this.stagedRemovals.StagedAddedDeleted {
		this.stagedRemovals.SetCommitted(this.committed) // Set the committed elements to the stagedRemovals.
	}
	return this
}

// // func (this *DeltaSet[string]) Print() {
// // 	fmt.Println("TotalSize: ", this.TotalSize)
// // 	fmt.Println("IsTransient: ", this.IsTransient)
// // 	fmt.Println("Committed: ", codec.Strings(this.DeltaSet.Committed().Elements()).ToHex())
// // 	fmt.Println("Updated  ", codec.Strings(this.DeltaSet.Added().Elements()).ToHex())
// // 	fmt.Println("Removed: ", codec.Strings(this.DeltaSet.Removed().Elements()).ToHex())
// // 	fmt.Println("Type: ", codec.Strings(this.DeltaSet.Removed().Elements()).ToHex())
// // 	fmt.Println()
// // }

func (this *DeltaSet[K]) StorageEncode(_ string) []byte {
	buffer, _ := rlp.EncodeToBytes(this.Encode())
	return buffer
}

func (this *DeltaSet[K]) StorageDecode(_ string, buffer []byte) any {
	var decoded []byte
	rlp.DecodeBytes(buffer, &decoded)
	return this.Decode(decoded)
}
