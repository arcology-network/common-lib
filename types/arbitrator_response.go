package types

import ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"

type ArbitratorResponse struct {
	ConflictedList []*ethCommon.Hash
	CPairLeft      []uint32
	CPairRight     []uint32
}
