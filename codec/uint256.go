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

package codec

import (
	uint256 "github.com/holiman/uint256"
)

type Uint256 uint256.Int

func (this *Uint256) Clone() interface{} {
	if this == nil {
		return this
	}

	return (*Uint256)((*uint256.Int)(this).Clone())
}

func (*Uint256) NewInt(v uint64) *Uint256 {
	return (*Uint256)(uint256.NewInt(v))
}

func (this *Uint256) Eq(v *Uint256) bool {
	return (*uint256.Int)(this).Eq((*uint256.Int)(v))
}

func (this *Uint256) Cmp(v *Uint256) int {
	return (*uint256.Int)(this).Cmp((*uint256.Int)(v))
}

func (this *Uint256) Add(lhv, rhv *Uint256) *Uint256 {
	return (*Uint256)((*uint256.Int)(this).Add((*uint256.Int)(lhv), (*uint256.Int)(rhv)))
}

func (this *Uint256) Sub(lhv, rhv *Uint256) *Uint256 {
	return (*Uint256)((*uint256.Int)(this).Sub((*uint256.Int)(lhv), (*uint256.Int)(rhv)))
}

func (this *Uint256) Uint64() uint64 {
	return (*uint256.Int)(this).ToBig().Uint64()
}

func (this *Uint256) Size() uint32 {
	return 32
}

func (this *Uint256) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this *Uint256) EncodeToBuffer(buffer []byte) int {
	return Uint64s((*uint256.Int)(this)[:]).EncodeToBuffer(buffer)
}

func (this *Uint256) Decode(buffer []byte) interface{} {
	copy(this[:], Uint64s{}.Decode(buffer).(Uint64s))
	return this
}
