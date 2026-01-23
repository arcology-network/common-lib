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
	"github.com/ethereum/go-ethereum/rlp"
)

func (this *StagedRemovalSet[K]) Size() int {
	size := this.DeltaSet.Size()

	// When CommittedDeleted is true, we only encode the staged removals and the added set.
	if this.CommittedDeleted {
		size = this.New(nil, this.Added(), this.Removed()).Size()
	}
	return size + 1 + 1 // 1 byte for the StagedAddedDeleted flag, 1 byte for the CommittedDeleted flag
}

// The EncodeTo method encodes the staged removal set to the provided buffer.
// Only the staged removals and the added set are encoded, the committed set is not encoded
// to save space. The committed is something that the recipient already has.
func (this *StagedRemovalSet[K]) EncodeTo(buffer []byte) int {
	offset := codec.Bool(this.StagedAddedDeleted).EncodeTo(buffer)
	offset += codec.Bool(this.CommittedDeleted).EncodeTo(buffer[offset:])

	// When committed is deleted, we only encode the staged removals and the added set.
	if this.CommittedDeleted {
		this.New(nil, this.Added(), this.Removed()).EncodeTo(buffer[offset:])
	} else {
		this.DeltaSet.EncodeTo(buffer[offset:])
	}
	return offset
}

func (this *StagedRemovalSet[K]) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeTo(buffer)
	return buffer
}

func (this *StagedRemovalSet[K]) Decode(buffer []byte) any {
	this.StagedAddedDeleted = bool(codec.Bool(this.StagedAddedDeleted).Decode(buffer).(codec.Bool))
	this.CommittedDeleted = bool(codec.Bool(this.CommittedDeleted).Decode(buffer[1:]).(codec.Bool))
	this.DeltaSet.Decode(buffer[2:]) // Skip the header

	return this
}

// StorageEncode encodes the staged removal set to a byte slice for storage.
// It uses RLP encoding to ensure compatibility with Ethereum storage.
// It is different from the Encode method which is used for network transmission.
func (this *StagedRemovalSet[K]) StorageEncode(_ string) []byte {
	buffer, _ := rlp.EncodeToBytes(this.Encode())
	return buffer
}

func (this *StagedRemovalSet[K]) StorageDecode(_ string, buffer []byte) any {
	var decoded []byte
	rlp.DecodeBytes(buffer, &decoded)
	return this.Decode(decoded)
}
