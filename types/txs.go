package types

import (
	codec "github.com/arcology-network/common-lib/codec"
)

type Txs struct {
	Data [][]byte
}

func (txs Txs) GobEncode() ([]byte, error) {
	return codec.Byteset(txs.Data).Encode(), nil
}
func (txs *Txs) GobDecode(data []byte) error {
	txs.Data = codec.Byteset{}.Decode(data).(codec.Byteset)
	return nil
}
