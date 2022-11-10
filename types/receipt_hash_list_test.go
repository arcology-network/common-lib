package types

import (
	"fmt"
	"testing"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
)

func TestReceiptHash(t *testing.T) {
	txhashes := []ethCommon.Hash{
		ethCommon.BytesToHash([]byte{1, 2, 3}),
		ethCommon.BytesToHash([]byte{4, 5, 6}),
		ethCommon.BytesToHash([]byte{7, 8, 9}),
		ethCommon.BytesToHash([]byte{10, 11, 12}),
		ethCommon.BytesToHash([]byte{13, 14, 15}),
		ethCommon.BytesToHash([]byte{16, 17, 18}),
		ethCommon.BytesToHash([]byte{19, 20, 21}),
	}
	rcpthashes := []ethCommon.Hash{
		ethCommon.BytesToHash([]byte{3, 2, 1}),
		ethCommon.BytesToHash([]byte{6, 5, 4}),
		ethCommon.BytesToHash([]byte{9, 8, 7}),
		ethCommon.BytesToHash([]byte{12, 11, 10}),
		ethCommon.BytesToHash([]byte{15, 14, 13}),
		ethCommon.BytesToHash([]byte{18, 17, 16}),
		ethCommon.BytesToHash([]byte{21, 20, 19}),
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
