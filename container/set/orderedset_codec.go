package orderedset

import (
	codec "github.com/arcology-network/common-lib/codec"
)

func (this *OrderedSet) Size() uint32 {
	return codec.Strings(this.keys).Size()
}

func (this *OrderedSet) Encode() []byte {
	return codec.Strings(this.keys).Encode()
}

func (this *OrderedSet) EncodeToBuffer(buffer []byte) int {
	return codec.Strings(this.keys).EncodeToBuffer(buffer)
}

func (*OrderedSet) Decode(buffer []byte) interface{} {
	keys := []string(codec.Strings{}.Decode(buffer).(codec.Strings))
	return NewOrderedSet(keys)
}
