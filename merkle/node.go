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

package merkle

import (
	"math"

	codec "github.com/arcology-network/common-lib/codec"
)

type Node struct {
	id       uint32
	level    uint32
	parent   uint32
	children []uint32
	hash     []byte
}

func NewNode() *Node {
	return &Node{
		children: make([]uint32, 0, 16),
	}
}

func (this *Node) Init(id uint32, level uint32, hash []byte) {
	this.id = id
	this.level = level
	this.parent = math.MaxUint32
	this.children = this.children[:0]
	this.hash = hash
}

func (this *Node) Children() []uint32 { return this.children }

func (this *Node) Encode() []byte {
	return codec.Byteset{
		codec.Uint32(this.id).Encode(),
		codec.Uint32(this.level).Encode(),
		codec.Uint32(this.parent).Encode(),
		codec.Uint32s(this.children).Encode(),
		this.hash[:],
	}.Encode()
}

func (node *Node) Decode(bytes []byte) interface{} {
	fields := codec.Byteset{}.Decode(bytes).(codec.Byteset)
	data := (&codec.Bytes{}).Decode(fields[4]).(codec.Bytes)
	*node = Node{
		uint32(codec.Uint32(0).Decode(fields[0]).(codec.Uint32)),
		uint32(codec.Uint32(0).Decode(fields[1]).(codec.Uint32)),
		uint32(codec.Uint32(0).Decode(fields[2]).(codec.Uint32)),
		[]uint32(codec.Uint32s{}.Decode(fields[3]).(codec.Uint32s)),
		[]byte(data),
	}
	return node
}
