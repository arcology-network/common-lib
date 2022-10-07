package datacompression

import (
	codec "github.com/HPISTechnologies/common-lib/codec"
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
	fields := codec.Byteset{}.Decode(bytes).(codec.Byteset)
	return Key{
		id:    uint32(codec.Uint32(0).Decode(fields[0]).(codec.Uint32)),
		to:    uint32(codec.Uint32(0).Decode(fields[1]).(codec.Uint32)),
		nonce: uint32(codec.Uint32(0).Decode(fields[2]).(codec.Uint32)),
	}
}
