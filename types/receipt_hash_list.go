package types

import (
	ethCommon "github.com/arcology/3rd-party/eth/common"
	encoding "github.com/arcology/common-lib/encoding"
)

type ReceiptHash struct {
	Txhash      *ethCommon.Hash
	Receipthash *ethCommon.Hash
	GasUsed     uint64
}

type ReceiptHashList struct {
	TxHashList      []ethCommon.Hash
	ReceiptHashList []ethCommon.Hash
	GasUsedList     []uint64
}

func (rhl *ReceiptHashList) GobEncode() ([]byte, error) {
	data := [][]byte{
		ethCommon.Hashes(rhl.TxHashList).Encode(),
		ethCommon.Hashes(rhl.ReceiptHashList).Encode(),
		encoding.Uint64s(rhl.GasUsedList).Encode(),
	}
	return encoding.Byteset(data).Encode(), nil
}
func (rhl *ReceiptHashList) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	rhl.TxHashList = ethCommon.Hashes(rhl.TxHashList).Decode(fields[0])
	rhl.ReceiptHashList = ethCommon.Hashes(rhl.ReceiptHashList).Decode(fields[1])
	rhl.GasUsedList = encoding.Uint64s(rhl.GasUsedList).Decode(fields[2])
	return nil
}
