package datacompression

import (
	codec "github.com/arcology-network/common-lib/codec"
)

func (this *CompressionLut) Encode() []byte {
	return codec.Byteset{
		codec.Strings(this.KeyLut).Encode(),
		codec.Uint32(this.offset).Encode(),
	}.Encode()
}

func (*CompressionLut) Decode(bytes []byte) interface{} {
	fields := codec.Byteset{}.Decode(bytes)
	return &CompressionLut{
		KeyLut: codec.Strings{}.Decode(fields[0]),
		offset: codec.Uint32(0).Decode(fields[1]),
	}
}
