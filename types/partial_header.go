package types

import (
	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
)

type PartialHeader struct {
	TxRoothash    ethCommon.Hash
	RcptRoothash  ethCommon.Hash
	StateRoothash ethCommon.Hash
}
