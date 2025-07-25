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

package stringdeltaset

import (
	"github.com/arcology-network/common-lib/codec"
	orderedset "github.com/arcology-network/common-lib/exp/orderedset"
)

func (this *SoftDeltaSet[K]) Size() int {
	return 4*codec.UINT64_LEN +
		this.committed.Size() +
		this.stagedAdditions.Size() +
		this.stagedRemovals.Size()
}

func (this *SoftDeltaSet[K]) EncodeTo(buffer []byte) int {
	offset := codec.Encoder{}.FillHeader(buffer,
		[]uint64{
			uint64(this.committed.Size()),
			uint64(this.stagedAdditions.Size()),
			uint64(this.stagedRemovals.Size()),
		},
	)

	offset += this.committed.EncodeTo(buffer[offset:])       // allDeleted
	offset += this.stagedAdditions.EncodeTo(buffer[offset:]) // stagedAdditions
	offset += this.stagedRemovals.EncodeTo(buffer[offset:])  // stagedRemovals
	return offset
}

func (this *SoftDeltaSet[K]) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeTo(buffer)
	return buffer
}

func (this *SoftDeltaSet[K]) Decode(buffer []byte) any {
	fields := codec.Byteset{}.Decode(buffer).(codec.Byteset) // Decode header
	this.committed = this.committed.Decode(fields[0]).(*orderedset.OrderedSet[K])
	this.stagedAdditions = this.stagedAdditions.Decode(fields[1]).(*orderedset.OrderedSet[K])
	this.stagedRemovals = this.stagedRemovals.Decode(fields[2]).(*StagedRemovalSet[K])
	return this
}

// // func (this *SoftDeltaSet[string]) Print() {
// // 	fmt.Println("TotalSize: ", this.TotalSize)
// // 	fmt.Println("IsTransient: ", this.IsTransient)
// // 	fmt.Println("Committed: ", codec.Strings(this.DeltaSet.Committed().Elements()).ToHex())
// // 	fmt.Println("Updated  ", codec.Strings(this.DeltaSet.Added().Elements()).ToHex())
// // 	fmt.Println("Removed: ", codec.Strings(this.DeltaSet.Removed().Elements()).ToHex())
// // 	fmt.Println("Type: ", codec.Strings(this.DeltaSet.Removed().Elements()).ToHex())
// // 	fmt.Println()
// // }

// func (this *SoftDeltaSet) StorageEncode(_ string) []byte {
// 	buffer, _ := rlp.EncodeToBytes(this.Encode())
// 	return buffer
// }

// func (this *SoftDeltaSet) StorageDecode(_ string, buffer []byte) any {
// 	var decoded []byte
// 	rlp.DecodeBytes(buffer, &decoded)
// 	return this.Decode(decoded)
// }
