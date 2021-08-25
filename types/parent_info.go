package types

import ethCommon "github.com/arcology/3rd-party/eth/common"

type ParentInfo struct {
	ParentRoot ethCommon.Hash
	ParentHash ethCommon.Hash
}
