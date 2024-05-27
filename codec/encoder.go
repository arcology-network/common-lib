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

func (Encoder) Size(args []interface{}) uint32 {
	length := uint32(0)
	for i := 0; i < len(args); i++ {
		if args[i] != nil {
			length += args[i].(Encodable).Size()
		}
	}
	return UINT32_LEN*uint32(len(args)+1) + uint32(length)
}

func (this Encoder) ToBuffer(buffer []byte, args []interface{}) {
	offset := uint32(0)
	Uint32(len(args)).EncodeToBuffer(buffer)
	for i := 0; i < len(args); i++ {
		Uint32(offset).EncodeToBuffer(buffer[(i+1)*UINT32_LEN:]) // Fill header info
		if args[i] != nil {
			offset += args[i].(Encodable).Size()
		}
	}
	headerSize := uint32((len(args) + 1) * UINT32_LEN)

	offset = uint32(0)
	for i := 0; i < len(args); i++ {
		if args[i] != nil {
			end := headerSize + offset + args[i].(Encodable).Size()
			args[i].(Encodable).EncodeToBuffer(buffer[headerSize+offset : end])
			offset += args[i].(Encodable).Size()
		}
	}
}

func (Encoder) FillHeader(buffer []byte, lengths []uint32) int {
	Uint32(len(lengths)).EncodeToBuffer(buffer[UINT32_LEN*0:])
	offset := uint32(0)
	for i := 0; i < len(lengths); i++ {
		Uint32(offset).EncodeToBuffer(buffer[UINT32_LEN*(i+1):])
		offset += uint32(lengths[i])
	}
	return (len(lengths) + 1) * UINT32_LEN
}
