package types

import (
	evmCommon "github.com/arcology-network/evm/common"
)

type PartialHeader struct {
	TxRoothash    evmCommon.Hash
	RcptRoothash  evmCommon.Hash
	StateRoothash evmCommon.Hash
}
