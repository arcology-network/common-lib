package types

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
)

type PartialHeader struct {
	TxRoothash    ethCommon.Hash
	RcptRoothash  ethCommon.Hash
	StateRoothash ethCommon.Hash
}
