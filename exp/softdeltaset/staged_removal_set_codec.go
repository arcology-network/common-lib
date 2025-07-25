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
	return 1 + this.DeltaSet.Size()
}

func (this *StagedRemovalSet[K]) EncodeTo(buffer []byte) int {
	offset := codec.Bool(this.allDeleted).EncodeTo(buffer)
	this.DeltaSet.EncodeTo(buffer[offset:])
	return offset
}

func (this *StagedRemovalSet[K]) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeTo(buffer)
	return buffer
}

func (this *StagedRemovalSet[K]) Decode(buffer []byte) any {
	this.allDeleted = bool(codec.Bool(this.allDeleted).Decode(buffer).(codec.Bool))
	this.DeltaSet.Decode(buffer[1:]) // Skip the header
	return this
}

func (this *StagedRemovalSet[K]) StorageEncode(_ string) []byte {
	buffer, _ := rlp.EncodeToBytes(this.Encode())
	return buffer
}

func (this *StagedRemovalSet[K]) StorageDecode(_ string, buffer []byte) any {
	var decoded []byte
	rlp.DecodeBytes(buffer, &decoded)
	return this.Decode(decoded)
}
