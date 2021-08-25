package merkle

import (
	codec "github.com/arcology/common-lib/codec"
	encoding "github.com/arcology/common-lib/encoding"
)

type Node struct {
	id       uint32
	level    uint32
	parent   uint32
	children []uint32
	hash     []byte
}

func NewNode(id uint32, level uint32, hash []byte) *Node {
	return &Node{
		id:       id,
		level:    level,
		parent:   UINT32_MAX,
		children: []uint32{},
		hash:     hash,
	}
}

func (this *Node) Encode() []byte {
	return codec.Byteset{
		codec.Uint32(this.id).Encode(),
		codec.Uint32(this.level).Encode(),
		codec.Uint32(this.parent).Encode(),
		encoding.Uint32s(this.children).Encode(),
		this.hash[:],
	}.Encode()
}

func (*Node) Decode(bytes []byte) interface{} {
	fields := codec.Byteset{}.Decode(bytes)
	return &Node{
		codec.Uint32(0).Decode(fields[0]),
		codec.Uint32(0).Decode(fields[1]),
		codec.Uint32(0).Decode(fields[2]),
		(&encoding.Uint32s{}).Decode(fields[3]),
		(&codec.Bytes{}).Decode(fields[4]).(codec.Bytes),
	}
}
