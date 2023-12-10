package types

import (
	evmCommon "github.com/ethereum/go-ethereum/common"
)

type PartialHeader struct {
	TxRoothash    evmCommon.Hash
	RcptRoothash  evmCommon.Hash
	StateRoothash evmCommon.Hash
}
