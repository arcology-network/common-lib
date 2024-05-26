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

package associative

type ipArray interface {
	[8]byte | [16]byte | [20]byte | [32]byte | [64]byte
}

type Bytesn[T ipArray] struct{}

// type Hashn8 Bytesn[[8]byte]

// func Clone[T ipArray](bytes T) T {
// 	target := *new(T)
// 	target = bytes
// 	return target
// }

// func Hex[T ipArray](bytes T) string {
// 	// var accHex [2 * len(bytes)]byte
// 	// accHex := make([]byte, 2*len(bytes))
// 	// hex.Encode(accHex[:], bytes[:])
// 	return string(bytes[:])
// }

// // type Bytesn[T] [HASH32_LEN]byte

// func (this *Bytesn[T]) Get() interface{} {
// 	return *this
// }

// func (this *Bytesn[T]) Set(v interface{}) {
// 	*this = v.(Bytesn[T])
// }

// func (this Bytesn[T]) Size() uint32 {
// 	return uint32(HASH32_LEN)
// }

// func (this Bytesn[T]) Sum(offset uint64) uint64 {
// 	total := uint64(0)
// 	for j := offset; j < uint64(len(this)); j++ {
// 		total += uint64((this)[j])
// 	}
// 	return total
// }

// func (this Bytesn[T]) Clone() interface{} {
// 	target := Bytesn[T]{}
// 	copy(target[:], this[:])
// 	return target
// }

// func (this Bytesn[T]) Encode() []byte {
// 	return this[:]
// }

// func (this Bytesn[T]) EncodeToBuffer(buffer []byte) int {
// 	copy(buffer, this[:])
// 	return len(this)
// }

// func (this Bytesn[T]) Decode(buffer []byte) interface{} {
// 	copy(this[:], buffer)
// 	return Bytesn[T](this)
// }

// func (this Bytesn[T]) Hex() string {
// 	var accHex [2 * len(this)]byte
// 	hex.Encode(accHex[:], this[:])
// 	return string(accHex[:])
// }

// func (this Bytesn[T]) UUID(seed uint64) Bytesn[T] {
// 	buffer := [HASH32_LEN + 8]byte{}
// 	copy(this[:], buffer[:])
// 	Uint64(uint64(seed)).EncodeToBuffer(buffer[len(this):])
// 	return sha256.Sum256(buffer[:])
// }

// type Bytesn[T]s [HASH32_LEN]byte

// func (this Bytesn[T]s) Clone() Bytesn[T]s {
// 	target := make([][HASH32_LEN]byte, len(this))
// 	for i := 0; i < len(this); i++ {
// 		copy(target[i][:], this[i][:])
// 	}
// 	return Bytesn[T]s(target)
// }

// func (this Bytesn[T]s) Encode() []byte {
// 	return Bytesn[T]s(this).Flatten()
// }

// func (this Bytesn[T]s) EncodeToBuffer(buffer []byte) int {
// 	for i := 0; i < len(this); i++ {
// 		copy(buffer[i*HASH32_LEN:], this[i][:])
// 	}
// 	return len(this) * HASH32_LEN
// }

// func (this Bytesn[T]s) Decode(buffer []byte) interface{} {
// 	if len(buffer) == 0 {
// 		return this
// 	}

// 	this = make([][HASH32_LEN]byte, len(buffer)/HASH32_LEN)
// 	for i := 0; i < len(this); i++ {
// 		copy(this[i][:], buffer[i*HASH32_LEN:(i+1)*HASH32_LEN])
// 	}
// 	return this
// }

// func (this Bytesn[T]s) Size() uint32 {
// 	return uint32(len(this) * HASH32_LEN)
// }

// func (this Bytesn[T]s) Flatten() []byte {
// 	buffer := make([]byte, len(this)*HASH32_LEN)
// 	this.EncodeToBuffer(buffer)
// 	return buffer
// }

// func (this Bytesn[T]s) Len() int {
// 	return len(this)
// }

// func (this Bytesn[T]s) Less(i, j int) bool {
// 	return bytes.Compare(this[i][:], this[j][:]) < 0
// }

// func (this Bytesn[T]s) Swap(i, j int) {
// 	this[i], this[j] = this[j], this[i]
// }
