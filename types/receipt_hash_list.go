package types

import (
	codec "github.com/arcology-network/common-lib/codec"
	ethCommon "github.com/ethereum/go-ethereum/common"
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
		Hashes(rhl.TxHashList).Encode(),
		Hashes(rhl.ReceiptHashList).Encode(),
		codec.Uint64s(rhl.GasUsedList).Encode(),
	}
	return codec.Byteset(data).Encode(), nil
}
func (rhl *ReceiptHashList) GobDecode(data []byte) error {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	rhl.TxHashList = Hashes(rhl.TxHashList).Decode(fields[0])
	rhl.ReceiptHashList = Hashes(rhl.ReceiptHashList).Decode(fields[1])
	rhl.GasUsedList = codec.Uint64s(rhl.GasUsedList).Decode(fields[2]).(codec.Uint64s)
	return nil
}
