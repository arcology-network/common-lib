package types

import (
	evmCommon "github.com/arcology-network/evm/common"
)

type ParentInfo struct {
	ParentRoot evmCommon.Hash
	ParentHash evmCommon.Hash
}
