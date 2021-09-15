package datacompression

import (
	codec "github.com/arcology-network/common-lib/codec"
)

type Key struct {
	id    uint32
	to    uint32
	nonce uint32 // nonce
}

func (this *Key) Encode() []byte {
	return codec.Byteset{
		codec.Uint32(this.id).Encode(),
		codec.Uint32(this.to).Encode(),
		codec.Uint32(this.nonce).Encode(),
	}.Encode()
}

func (*Key) Decode(bytes []byte) interface{} {
	fields := codec.Byteset{}.Decode(bytes)
	return Key{
		id:    codec.Uint32(0).Decode(fields[0]),
		to:    codec.Uint32(0).Decode(fields[1]),
		nonce: codec.Uint32(0).Decode(fields[2]),
	}
}
