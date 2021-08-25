package types

import (
	ethCommon "github.com/arcology/3rd-party/eth/common"
	encoding "github.com/arcology/common-lib/encoding"
)

type MetaBlock struct {
	Txs      [][]byte
	Hashlist []*ethCommon.Hash
}

func (mb MetaBlock) GobEncode() ([]byte, error) {
	hashArray := Ptr2Arr(mb.Hashlist)
	data := [][]byte{
		encoding.Byteset(mb.Txs).Encode(),
		ethCommon.Hashes(hashArray).Encode(),
	}
	return encoding.Byteset(data).Encode(), nil
}
func (mb *MetaBlock) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	mb.Txs = encoding.Byteset{}.Decode(fields[0])
	arrs := ethCommon.Hashes([]ethCommon.Hash{}).Decode(fields[1])
	mb.Hashlist = Arr2Ptr(arrs)
	return nil
}
