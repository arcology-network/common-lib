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

package matrix

import (
	"bytes"
	"fmt"
	"os"

	"github.com/arcology-network/common-lib/codec"
)

type BitMatrix struct {
	width  int
	height int
	data   []byte
}

func NewBitMatrix(width, height int, initv bool) *BitMatrix {
	size := (width * height) / 8
	if (width*height)%8 != 0 {
		size++
	}

	matrix := &BitMatrix{
		width:  width,
		height: height,
		data:   make([]byte, size),
	}

	matrix.Fill(initv)
	return matrix
}

func (this *BitMatrix) Get(x, y int) bool {
	index := (y*this.width + x) / 8
	offset := uint((y*this.width + x) % 8)
	return (this.data[index] & (1 << offset)) != 0
}

func (this *BitMatrix) Set(x, y int, value bool) {
	index := (y*this.width + x) / 8
	offset := uint((y*this.width + x) % 8)
	if value {
		this.data[index] |= (1 << offset)
	} else {
		this.data[index] &= ^(1 << offset)
	}
}

func (this *BitMatrix) Foreach(fun func(x, y int, v bool) bool) {
	for i := 0; i < this.width; i++ {
		for j := 0; j < this.height; j++ {
			this.Set(i, j, fun(i, j, this.Get(i, j)))
		}
	}
}

func (this *BitMatrix) Fill(value bool) *BitMatrix {
	v := byte(0)
	if value == true {
		v = 0xff
	}

	for i := 0; i < len(this.data); i++ {
		this.data[i] = v
	}
	return this
}

func (this *BitMatrix) CountInCol(col int, v bool) int {
	total := 0
	for i := 0; i < this.height; i++ {
		if this.Get(col, i) == v {
			total++
		}
	}
	return total
}

func (this *BitMatrix) CountInRow(row int, v bool) int {
	total := 0
	for i := 0; i < this.width; i++ {
		if this.Get(i, row) == v {
			total++
		}
	}
	return total
}

func (this *BitMatrix) FillCol(col int, v bool) {
	for i := 0; i < this.height; i++ {
		this.Set(col, i, v)
	}
}

func (this *BitMatrix) FillRow(row int, v bool) {
	for i := 0; i < this.width; i++ {
		this.Set(i, row, v)
	}
}

func (this *BitMatrix) Width() int  { return this.width }
func (this *BitMatrix) Height() int { return this.height }
func (this *BitMatrix) Raw() []byte { return this.data }

func (this *BitMatrix) Equal(other *BitMatrix) bool {
	return this.width == other.width && this.height == other.height && bytes.Equal(this.data, other.data)
}

func (this *BitMatrix) Print() {
	for i := 0; i < this.height; i++ {
		for j := 0; j < this.width; j++ {
			if this.Get(j, i) {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
}

func (this *BitMatrix) WriteToFile(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := codec.Byteset([][]byte{
		codec.Uint32(this.width).Encode(),
		codec.Uint32(this.height).Encode(),
		this.data,
	}).Encode()
	_, err = file.Write(buffer)
	return err
}

func (this *BitMatrix) ReadFromFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	buffer := make([]byte, stat.Size())
	_, err = file.Read(buffer)

	buffers := [][]byte(new(codec.Byteset).Decode(buffer).(codec.Byteset))
	this.width = int(new(codec.Uint32).Decode(buffers[0]).(codec.Uint32))
	this.height = int(new(codec.Uint32).Decode(buffers[1]).(codec.Uint32))
	this.data = buffers[2]
	return nil
}
