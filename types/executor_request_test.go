package types

import (
	"fmt"
	"math/big"
	"testing"

	evmCommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
)

func TestExecutorRequestEncodingAndDeconing(t *testing.T) {
	from := evmCommon.BytesToAddress([]byte{0, 0, 0, 5, 6, 7, 8, 9})
	to := evmCommon.BytesToAddress([]byte{11, 12, 13})
	from0 := evmCommon.BytesToAddress([]byte{0, 0, 0, 1, 4, 5, 6, 7, 8})
	to0 := evmCommon.BytesToAddress([]byte{11, 12, 13, 14})

	ethMsg_serial_0 := core.NewMessage(from0, &to0, 1, big.NewInt(int64(1)), 100, big.NewInt(int64(8)), []byte{1, 2, 3}, nil, false)
	ethMsg_serial_1 := core.NewMessage(from, &to, 3, big.NewInt(int64(100)), 200, big.NewInt(int64(9)), []byte{4, 5, 6}, nil, false)
	fmt.Printf("ethMsg_serial_0=%v\n", ethMsg_serial_0)
	fmt.Printf("ethMsg_serial_1=%v\n", ethMsg_serial_1)
	hash1 := RlpHash(ethMsg_serial_0)
	hash2 := RlpHash(ethMsg_serial_1)
	fmt.Printf("hash1=%v\n", hash1)
	fmt.Printf("hash2=%v\n", hash2)
}
