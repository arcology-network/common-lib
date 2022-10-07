package merkle

import (
	"math"

	codec "github.com/HPISTechnologies/common-lib/codec"
	"github.com/HPISTechnologies/common-lib/encoding"
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

func (this *Node) Encode() []byte {
	return codec.Byteset{
		codec.Uint32(this.id).Encode(),
		codec.Uint32(this.level).Encode(),
		codec.Uint32(this.parent).Encode(),
		encoding.Uint32s(this.children).Encode(),
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
