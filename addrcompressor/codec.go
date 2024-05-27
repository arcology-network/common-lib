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

func (this *CompressionLut) Encode() []byte {
	return codec.Byteset{
		codec.Strings(this.IdxToKeyLut).Encode(),
	}.Encode()
}

func (*CompressionLut) Decode(bytes []byte) interface{} {
	fields := codec.Byteset{}.Decode(bytes).(codec.Byteset)
	return &CompressionLut{
		IdxToKeyLut: codec.Strings{}.Decode(fields[0]).(codec.Strings),
	}
}
