package types

import (
	evmCommon "github.com/ethereum/go-ethereum/common"
)

type ParentInfo struct {
	ParentRoot evmCommon.Hash
	ParentHash evmCommon.Hash
}
