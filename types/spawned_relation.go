package types

import (
	ethCommon "github.com/arcology/3rd-party/eth/common"
)

type SpawnedRelation struct {
	Txhash        ethCommon.Hash
	SpawnedTxHash ethCommon.Hash
}
