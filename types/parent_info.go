package types

import ethCommon "github.com/arcology-network/3rd-party/eth/common"

type ParentInfo struct {
	ParentRoot ethCommon.Hash
	ParentHash ethCommon.Hash
}
