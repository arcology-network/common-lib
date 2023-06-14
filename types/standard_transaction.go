package types

import (
	encoding "github.com/arcology-network/common-lib/encoding"
	evmCommon "github.com/arcology-network/evm/common"
	evmTypes "github.com/arcology-network/evm/core/types"
)

const (
	TxType_Eth  = 0
	TxType_Coin = 1

	TxFrom_Remote = 1
	TxFrom_Local  = 2
	TxFrom_Block  = 3
)

type StandardTransaction struct {
	TxHash    evmCommon.Hash
	Native    *evmTypes.Transaction
	TxRawData []byte
	Source    uint8
}

func (stdTx *StandardTransaction) Hash() evmCommon.Hash { return stdTx.TxHash }

func (stdTx *StandardTransaction) GobEncode() ([]byte, error) {
	data := [][]byte{
		stdTx.TxHash.Bytes(),
		[]byte{stdTx.Source},
		stdTx.TxRawData,
	}
	return encoding.Byteset(data).Encode(), nil
}
func (stdTx *StandardTransaction) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	stdTx.TxHash = evmCommon.BytesToHash(fields[0])
	stdTx.Source = uint8(fields[1][0])
	stdTx.TxRawData = fields[2]
	return nil
}
