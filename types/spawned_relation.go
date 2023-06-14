package types

import (
	evmCommon "github.com/arcology-network/evm/common"
)

type SpawnedRelation struct {
	Txhash        evmCommon.Hash
	SpawnedTxHash evmCommon.Hash
}
