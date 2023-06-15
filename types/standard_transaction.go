package types

import (
	ethCommon "github.com/arcology-network/3rd-party/eth/common"
	ethTypes "github.com/arcology-network/3rd-party/eth/types"
	encoding "github.com/arcology-network/common-lib/encoding"
)

const (
	TxType_Eth  = 0
	TxType_Coin = 1

	TxFrom_Remote = 1
	TxFrom_Local  = 2
	TxFrom_Block  = 3
)

const (
	TX_SOURCE_REMOTE = iota
	TX_SOURCE_LOCAL
	TX_SOURCE_BLOCK
	TX_SOURCE_DEFERRED
)

type StandardTransaction struct {
	TxHash    ethCommon.Hash
	Native    *ethTypes.Transaction
	TxRawData []byte
	Source    uint8
}

func (stdTx *StandardTransaction) Hash() ethCommon.Hash { return stdTx.TxHash }

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
	stdTx.TxHash = ethCommon.BytesToHash(fields[0])
	stdTx.Source = uint8(fields[1][0])
	stdTx.TxRawData = fields[2]
	return nil
}
