package types

import (
	"fmt"
	"testing"

	evmCommon "github.com/arcology-network/evm/common"
)

func TestReceiptHash(t *testing.T) {
	txhashes := []evmCommon.Hash{
		evmCommon.BytesToHash([]byte{1, 2, 3}),
		evmCommon.BytesToHash([]byte{4, 5, 6}),
		evmCommon.BytesToHash([]byte{7, 8, 9}),
		evmCommon.BytesToHash([]byte{10, 11, 12}),
		evmCommon.BytesToHash([]byte{13, 14, 15}),
		evmCommon.BytesToHash([]byte{16, 17, 18}),
		evmCommon.BytesToHash([]byte{19, 20, 21}),
	}
	rcpthashes := []evmCommon.Hash{
		evmCommon.BytesToHash([]byte{3, 2, 1}),
		evmCommon.BytesToHash([]byte{6, 5, 4}),
		evmCommon.BytesToHash([]byte{9, 8, 7}),
		evmCommon.BytesToHash([]byte{12, 11, 10}),
		evmCommon.BytesToHash([]byte{15, 14, 13}),
		evmCommon.BytesToHash([]byte{18, 17, 16}),
		evmCommon.BytesToHash([]byte{21, 20, 19}),
	}
	gasUsedList := []uint64{
		uint64(10), uint64(11), uint64(12), uint64(13), uint64(14), uint64(15), uint64(16),
	}

	rcptList := ReceiptHashList{
		TxHashList:      txhashes,
		ReceiptHashList: rcpthashes,
		GasUsedList:     gasUsedList,
	}

	datas, err := rcptList.GobEncode()
	if err != nil {
		fmt.Printf(" rcptList.GobEncode err=%v\n", err)
		return

	}
	fmt.Printf(" rcptList.GobEncode result=%x\n", datas)

	rcptListResult := ReceiptHashList{}
	err = rcptListResult.GobDecode(datas)
	if err != nil {
		fmt.Printf(" rcptList.GobDecode err=%v\n", err)
		return

	}
	fmt.Printf(" rcptList.GobDecode result=%v\n", rcptListResult)
}
