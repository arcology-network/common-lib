package types

import ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"

type ParentInfo struct {
	ParentRoot ethCommon.Hash
	ParentHash ethCommon.Hash
}
