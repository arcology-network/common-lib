package types

import (
	"fmt"
	"testing"

	evmCommon "github.com/ethereum/go-ethereum/common"
)

func Test_getList(t *testing.T) {
	hashes := []evmCommon.Hash{
		evmCommon.BytesToHash([]byte{1, 2, 3}),
		evmCommon.BytesToHash([]byte{4, 5, 6}),
		evmCommon.BytesToHash([]byte{7, 8, 9}),
		evmCommon.BytesToHash([]byte{10, 11, 12}),
		evmCommon.BytesToHash([]byte{13, 14, 15}),
		evmCommon.BytesToHash([]byte{16, 17, 18}),
		evmCommon.BytesToHash([]byte{19, 20, 21}),
	}
	succeful := []bool{true, true, true, true, true, false, true}
	inclusive := InclusiveList{
		HashList:   Arr2Ptr(hashes),
		Successful: succeful,
	}
	inclusive.Mode = InclusiveMode_Results
	selList, clearList := inclusive.GetList()
	fmt.Printf(" inclusive.getlist selList=%v\n,clearList=%v\n", selList, clearList)
}
func TestInclusive(t *testing.T) {
	hashes := []evmCommon.Hash{
		evmCommon.BytesToHash([]byte{1, 2, 3}),
		evmCommon.BytesToHash([]byte{4, 5, 6}),
		evmCommon.BytesToHash([]byte{7, 8, 9}),
		evmCommon.BytesToHash([]byte{10, 11, 12}),
		evmCommon.BytesToHash([]byte{13, 14, 15}),
		evmCommon.BytesToHash([]byte{16, 17, 18}),
		evmCommon.BytesToHash([]byte{19, 20, 21}),
	}
	succeful := []bool{true, true, true, true, true, false, true}
	inclusive := InclusiveList{
		HashList:   Arr2Ptr(hashes),
		Successful: succeful,
	}

	datas, err := inclusive.GobEncode()
	if err != nil {
		fmt.Printf(" inclusive.GobEncode err=%v\n", err)
		return

	}
	fmt.Printf(" inclusive.GobEncode result=%x\n", datas)

	inclusiveResult := InclusiveList{}
	err = inclusiveResult.GobDecode(datas)
	if err != nil {
		fmt.Printf(" inclusive.GobDecode err=%v\n", err)
		return

	}
	fmt.Printf(" inclusive.GobDecode result=%v\n", inclusiveResult)
}
