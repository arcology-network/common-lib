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
	"encoding/binary"
	"math"
)

type Encoder struct{}

func (Encoder) FillHeader(buffer []byte, lengths []uint64) int {
	Uint32(len(lengths)).EncodeTo(buffer[UINT64_LEN*0:])
	offset := uint64(0)
	for i := 0; i < len(lengths); i++ {
		Uint32(offset).EncodeTo(buffer[UINT64_LEN*(i+1):])
		offset += uint64(lengths[i])
	}
	return (len(lengths) + 1) * UINT64_LEN
}

func Sizer[T any](t T) int {
	switch v := any(t).(type) {
	case string:
		return len(v)
	case []byte:
		return len(v)
	case int64:
		return INT64_LEN
	case uint64:
		return UINT64_LEN
	case float32:
		return 4
	case float64:
		return 8
	case bool:
		return BOOL_LEN
	case Hash64:
		return HASH64_LEN
	case Hash64s:
		return len(v) * HASH64_LEN
	case Int64:
		return INT64_LEN
	case Int64s:
		return len(v) * INT64_LEN
	case Encodable:
		return int(v.Size())
	default:
		return 0
	}
}
func EncodeTo[T any](t T, buf []byte) int {
	switch v := any(t).(type) {
	case string:
		return String(v).EncodeTo(buf)
	case []byte:
		return Bytes(v).EncodeTo(buf)
	case int64:
		return Int64(v).EncodeTo(buf)
	case uint64:
		return Uint64(v).EncodeTo(buf)
	case float32:
		binary.LittleEndian.PutUint32(buf, math.Float32bits(v))
		return 4
	case float64:
		binary.LittleEndian.PutUint64(buf, math.Float64bits(v))
		return 8
	case bool:
		if v {
			buf[0] = 1
		} else {
			buf[0] = 0
		}
		return BOOL_LEN
	// case Hash64:
	// 	return Hash64(v).EncodeTo(buf)
	case Hash64s:
		return Hash64s(v).EncodeTo(buf)
	case Int64:
		return Int64(v).EncodeTo(buf)
	case Int64s:
		return Int64s(v).EncodeTo(buf)
	case Encodable:
		return v.EncodeTo(buf)
	default:
		return 0
	}
}
