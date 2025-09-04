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
	"math/big"
)

type Bigint big.Int

func (this *Bigint) Clone() interface{} {
	if this == nil {
		return this
	}

	return (*Bigint)((*big.Int)(this).Set(new(big.Int)))
}

func (this *Bigint) Size() uint64 {
	return BOOL_LEN + uint64((*big.Int)(this).BitLen())
}

func (this *Bigint) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeTo(buffer)
	return buffer
}

func (this *Bigint) EncodeTo(buffer []byte) int {
	flag := (*big.Int)(this).Cmp(big.NewInt(0)) >= 0

	Bool(flag).EncodeTo(buffer)
	val := (big.Int)(*this)
	val.FillBytes(buffer[1 : val.BitLen()+1])
	return (val.BitLen() + 1)
}

func (this *Bigint) Decode(buffer []byte) interface{} {
	if len(buffer) > 0 {
		v := new(big.Int)
		*this = *(*Bigint)(v.SetBytes(buffer[1:]))
		if !Bool(true).Decode(buffer[:1]).(Bool) { // negative value
			return (*Bigint)((&big.Int{}).Neg(v))
		}
	}
	return (*Bigint)(this)
}
