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

type Encodable interface {
	Clone() interface{}
	Size() uint64
	Encode() []byte
	EncodeTo([]byte) int
	Decode([]byte) interface{}
}

type Encodables []Encodable

func (this Encodables) Size() uint64 {
	length := uint64(0)
	for i := 0; i < len(this); i++ {
		if this[i] != nil {
			length += this[i].Size()
		}
	}
	return UINT64_LEN*uint64(len(this)+1) + uint64(length)
}

func (this Encodables) Sizes() []uint64 {
	lengths := make([]uint64, len(this))
	for i := 0; i < len(lengths); i++ {
		if this[i] != nil {
			lengths[i] += this[i].Size()
		}
	}
	return lengths
}

func (this Encodables) FillHeader(buffer []byte) int {
	lengths := this.Sizes()
	Uint32(len(lengths)).EncodeTo(buffer[UINT64_LEN*0:])
	offset := uint64(0)
	for i := 0; i < len(lengths); i++ {
		Uint32(offset).EncodeTo(buffer[UINT64_LEN*(i+1):])
		offset += uint64(lengths[i])
	}
	return (len(lengths) + 1) * UINT64_LEN
}

func (this Encodables) Encode() []byte {
	total := this.Size()
	buffer := make([]byte, total)
	this.EncodeTo(buffer)
	return buffer
}

func (this Encodables) EncodeTo(buffer []byte) int {
	offset := this.FillHeader(buffer)
	for i := 0; i < len(this); i++ {
		// if selectors[i] {
		offset += this[i].EncodeTo(buffer[offset:])
		// }
	}
	return offset
}

func (this Encodables) Decode(buffer []byte, decoders ...func([]byte) interface{}) []interface{} {
	fields := Byteset{}.Decode(buffer).(Byteset)
	values := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		values[i] = decoders[i](fields[i])
	}
	return values
}

func (Encoder) Size(args []any) uint64 {
	length := uint64(0)
	for i := 0; i < len(args); i++ {
		if args[i] != nil {
			length += args[i].(Encodable).Size()
		}
	}
	return UINT64_LEN*uint64(len(args)+1) + uint64(length)
}

func (this Encoder) ToBuffer(buffer []byte, args []any) {
	offset := uint64(0)
	Uint32(len(args)).EncodeTo(buffer)
	for i := 0; i < len(args); i++ {
		Uint32(offset).EncodeTo(buffer[(i+1)*UINT64_LEN:]) // Fill header info
		if args[i] != nil {
			offset += args[i].(Encodable).Size()
		}
	}
	headerSize := uint64((len(args) + 1) * UINT64_LEN)

	offset = uint64(0)
	for i := 0; i < len(args); i++ {
		if args[i] != nil {
			end := headerSize + offset + args[i].(Encodable).Size()
			args[i].(Encodable).EncodeTo(buffer[headerSize+offset : end])
			offset += args[i].(Encodable).Size()
		}
	}
}
