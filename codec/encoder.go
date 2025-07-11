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

type Encoder struct{}

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
	Uint32(len(args)).EncodeToBuffer(buffer)
	for i := 0; i < len(args); i++ {
		Uint32(offset).EncodeToBuffer(buffer[(i+1)*UINT64_LEN:]) // Fill header info
		if args[i] != nil {
			offset += args[i].(Encodable).Size()
		}
	}
	headerSize := uint64((len(args) + 1) * UINT64_LEN)

	offset = uint64(0)
	for i := 0; i < len(args); i++ {
		if args[i] != nil {
			end := headerSize + offset + args[i].(Encodable).Size()
			args[i].(Encodable).EncodeToBuffer(buffer[headerSize+offset : end])
			offset += args[i].(Encodable).Size()
		}
	}
}

func (Encoder) FillHeader(buffer []byte, lengths []uint64) int {
	Uint32(len(lengths)).EncodeToBuffer(buffer[UINT64_LEN*0:])
	offset := uint64(0)
	for i := 0; i < len(lengths); i++ {
		Uint32(offset).EncodeToBuffer(buffer[UINT64_LEN*(i+1):])
		offset += uint64(lengths[i])
	}
	return (len(lengths) + 1) * UINT64_LEN
}
