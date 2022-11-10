package datacompression

import (
	codec "github.com/arcology-network/common-lib/codec"
)

func (this *CompressionLut) Encode() []byte {
	return codec.Byteset{
		codec.Strings(this.IdxToKeyLut).Encode(),
	}.Encode()
}

func (*CompressionLut) Decode(bytes []byte) interface{} {
	fields := codec.Byteset{}.Decode(bytes).(codec.Byteset)
	return &CompressionLut{
		IdxToKeyLut: codec.Strings{}.Decode(fields[0]).(codec.Strings),
	}
}
