package types

import (
	"fmt"
	"math/big"
	"testing"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
	ethTypes "github.com/HPISTechnologies/3rd-party/eth/types"
)

func TestExecutorRequestEncodingAndDeconing(t *testing.T) {
	from := ethCommon.BytesToAddress([]byte{0, 0, 0, 5, 6, 7, 8, 9})
	to := ethCommon.BytesToAddress([]byte{11, 12, 13})
	from0 := ethCommon.BytesToAddress([]byte{0, 0, 0, 1, 4, 5, 6, 7, 8})
	to0 := ethCommon.BytesToAddress([]byte{11, 12, 13, 14})
	ethMsg_serial_0 := ethTypes.NewMessage(from0, &to0, 1, big.NewInt(int64(1)), 100, big.NewInt(int64(8)), []byte{1, 2, 3}, false)
	ethMsg_serial_1 := ethTypes.NewMessage(from, &to, 3, big.NewInt(int64(100)), 200, big.NewInt(int64(9)), []byte{4, 5, 6}, false)
	fmt.Printf("ethMsg_serial_0=%v\n", ethMsg_serial_0)
	fmt.Printf("ethMsg_serial_1=%v\n", ethMsg_serial_1)
	hash1 := ethCommon.RlpHash(ethMsg_serial_0)
	hash2 := ethCommon.RlpHash(ethMsg_serial_1)
	fmt.Printf("hash1=%v\n", hash1)
	fmt.Printf("hash2=%v\n", hash2)
}
