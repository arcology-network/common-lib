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

package addrcompressor

import (
	codec "github.com/arcology-network/common-lib/codec"
)

type Key struct {
	id    uint32
	to    uint32
	nonce uint32 // nonce
}

func (this *Key) Encode() []byte {
	return codec.Byteset{
		codec.Uint32(this.id).Encode(),
		codec.Uint32(this.to).Encode(),
		codec.Uint32(this.nonce).Encode(),
	}.Encode()
}

func (*Key) Decode(bytes []byte) interface{} {
	fields := codec.Byteset{}.Decode(bytes).(codec.Byteset)
	return Key{
		id:    uint32(codec.Uint32(0).Decode(fields[0]).(codec.Uint32)),
		to:    uint32(codec.Uint32(0).Decode(fields[1]).(codec.Uint32)),
		nonce: uint32(codec.Uint32(0).Decode(fields[2]).(codec.Uint32)),
	}
}
