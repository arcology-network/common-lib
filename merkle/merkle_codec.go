package merkle

import (
	"reflect"

	codec "github.com/arcology-network/common-lib/codec"
)

func (this *Merkle) Encode() []byte {
	hashes := [][]byte{}
	if reflect.ValueOf(this.hasher) == reflect.ValueOf(Sha256) {
		hashes = append(hashes, codec.Uint8(0).Encode())
	}

	if reflect.ValueOf(this.hasher) == reflect.ValueOf(Keccak256) {
		hashes = append(hashes, codec.Uint8(1).Encode())
	}

	for _, nodes := range this.nodes {
		hashVec := [][]byte{}
		for _, node := range nodes {
			hashVec = append(hashVec, node.Encode())
		}
		hashes = append(hashes, codec.Byteset(hashVec).Encode())
	}
	return codec.Byteset(hashes).Encode()
}

func (*Merkle) Decode(bytes []byte) interface{} {
	merkle := &Merkle{}
	fields := codec.Byteset{}.Decode(bytes).(codec.Byteset)
	switch uint8(codec.Uint8(0).Decode(fields[0]).(codec.Uint8)) {
	case 0:
		merkle.hasher = Sha256
	case 1:
		merkle.hasher = Keccak256
	}

	for i := 1; i < len(fields); i++ {
		level := []*Node{}
		subFields := codec.Byteset{}.Decode(fields[i]).(codec.Byteset)
		for _, subField := range subFields {
			level = append(level, (&Node{}).Decode(subField).(*Node))
		}
		merkle.nodes = append(merkle.nodes, level)
	}

	merkle.encoder = Concatenator{}.Encode
	return merkle
}
