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

package merkle

import (
	codec "github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
)

func (this *Merkle) Encode() []byte {
	hashes := [][]byte{}

	if common.IsType[Sha256](this.hasher) {
		hashes = append(hashes, codec.Uint8(0).Encode())
	}

	if common.IsType[Keccak256](this.hasher) {
		hashes = append(hashes, codec.Uint8(1).Encode())
	}

	for _, nodes := range this.nodes {
		hashVec := [][]byte{}
		for _, node := range nodes {
			hashVec = append(hashVec, node.Encode())
		}
		hashes = append(hashes, codec.Byteset(hashVec).Encode())
	}
	return codec.Byteset(hashes).Encode()
}

func (*Merkle) Decode(bytes []byte) interface{} {
	merkle := &Merkle{}
	fields := codec.Byteset{}.Decode(bytes).(codec.Byteset)
	switch uint8(codec.Uint8(0).Decode(fields[0]).(codec.Uint8)) {
	case 0:
		merkle.hasher = Sha256{}
	case 1:
		merkle.hasher = Keccak256{}
	}

	for i := 1; i < len(fields); i++ {
		level := []*Node{}
		subFields := codec.Byteset{}.Decode(fields[i]).(codec.Byteset)
		for _, subField := range subFields {
			level = append(level, (&Node{}).Decode(subField).(*Node))
		}
		merkle.nodes = append(merkle.nodes, level)
	}

	// merkle.encoder = Concatenator{}
	return merkle
}
