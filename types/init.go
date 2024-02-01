package types

import (
	"encoding/gob"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

func init() {
	gob.Register(&InclusiveList{})
	gob.Register(&ReapingList{})
	gob.Register(&ReceiptHashList{})

	gob.Register(&StandardTransaction{})

	gob.Register([]*StandardTransaction{})
	gob.Register([]*StandardTransaction{})

	gob.Register([][]byte{})
	gob.Register([]byte{})

	gob.Register(&big.Int{})

	gob.Register(map[ethCommon.Hash]ethCommon.Hash{})

	gob.Register(&IncomingTxs{})

}
