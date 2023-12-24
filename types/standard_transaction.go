package types

import (
	codec "github.com/arcology-network/common-lib/codec"
	ethCommon "github.com/ethereum/go-ethereum/common"
	evmTypes "github.com/ethereum/go-ethereum/core/types"
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
	Native    *evmTypes.Transaction
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
	return codec.Byteset(data).Encode(), nil
}
func (stdTx *StandardTransaction) GobDecode(data []byte) error {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	stdTx.TxHash = ethCommon.BytesToHash(fields[0])
	stdTx.Source = uint8(fields[1][0])
	stdTx.TxRawData = fields[2]
	return nil
}
