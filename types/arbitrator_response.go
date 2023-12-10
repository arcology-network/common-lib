package types

import evmCommon "github.com/ethereum/go-ethereum/common"

type ArbitratorResponse struct {
	ConflictedList []*evmCommon.Hash
	CPairLeft      []uint32
	CPairRight     []uint32
}
