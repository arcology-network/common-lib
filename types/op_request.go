package types

import (
	evmTypes "github.com/ethereum/go-ethereum/core/types"
)

type OpRequest struct {
	BlockParam   *BlockParams
	Withdrawals  evmTypes.Withdrawals // List of withdrawals to include in block.
	Transactions []*StandardTransaction
}
