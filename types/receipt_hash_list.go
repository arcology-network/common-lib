package types

import (
	encoding "github.com/arcology-network/common-lib/encoding"
	evmCommon "github.com/ethereum/go-ethereum/common"
)

type ReceiptHash struct {
	Txhash      *evmCommon.Hash
	Receipthash *evmCommon.Hash
	GasUsed     uint64
}

type ReceiptHashList struct {
	TxHashList      []evmCommon.Hash
	ReceiptHashList []evmCommon.Hash
	GasUsedList     []uint64
}

func (rhl *ReceiptHashList) GobEncode() ([]byte, error) {
	data := [][]byte{
		Hashes(rhl.TxHashList).Encode(),
		Hashes(rhl.ReceiptHashList).Encode(),
		encoding.Uint64s(rhl.GasUsedList).Encode(),
	}
	return encoding.Byteset(data).Encode(), nil
}
func (rhl *ReceiptHashList) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	rhl.TxHashList = Hashes(rhl.TxHashList).Decode(fields[0])
	rhl.ReceiptHashList = Hashes(rhl.ReceiptHashList).Decode(fields[1])
	rhl.GasUsedList = encoding.Uint64s(rhl.GasUsedList).Decode(fields[2])
	return nil
}
