package types

import (
	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
)

type SpawnedRelation struct {
	Txhash        ethCommon.Hash
	SpawnedTxHash ethCommon.Hash
}
