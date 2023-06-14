package types

import evmCommon "github.com/arcology-network/evm/common"

type ArbitratorResponse struct {
	ConflictedList []*evmCommon.Hash
	CPairLeft      []uint32
	CPairRight     []uint32
}
