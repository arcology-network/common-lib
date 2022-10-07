package types

import (
	encoding "github.com/HPISTechnologies/common-lib/encoding"
)

type Txs struct {
	Data [][]byte
}

func (txs Txs) GobEncode() ([]byte, error) {
	return encoding.Byteset(txs.Data).Encode(), nil
}
func (txs *Txs) GobDecode(data []byte) error {
	txs.Data = encoding.Byteset{}.Decode(data)
	return nil
}
