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
	"crypto/sha256"

	slice "github.com/arcology-network/common-lib/exp/slice"
	"golang.org/x/crypto/sha3"
)

type Concatenator struct{}

func (Concatenator) Encode(bytes [][]byte) []byte { return slice.Flatten(bytes) }

type Sha256 struct{}

func (Sha256) Hash(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	hash := sha256.Sum256(data)
	return hash[:]
}

type Keccak256 struct{}

func (Keccak256) Hash(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	return hasher.Sum(nil)
}
