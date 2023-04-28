package orderedset

import (
	codec "github.com/arcology-network/common-lib/codec"
)

func (this *OrderedSet) Encode(processors ...func(interface{}) interface{}) []byte {
	return codec.Strings(this.keys).Encode()
}

func (this *OrderedSet) EncodeToBuffer(buffer []byte, processors ...func(interface{}) interface{}) int {
	return codec.Strings(this.keys).EncodeToBuffer(buffer)
}

func (*OrderedSet) Decode(buffer []byte) interface{} {
	keys := []string(codec.Strings{}.Decode(buffer).(codec.Strings))
	return NewOrderedSet(keys)
}
