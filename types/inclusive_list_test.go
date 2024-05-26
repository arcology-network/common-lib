package types

import (
	"fmt"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

func Test_getList(t *testing.T) {
	hashes := []ethCommon.Hash{
		ethCommon.BytesToHash([]byte{1, 2, 3}),
		ethCommon.BytesToHash([]byte{4, 5, 6}),
		ethCommon.BytesToHash([]byte{7, 8, 9}),
		ethCommon.BytesToHash([]byte{10, 11, 12}),
		ethCommon.BytesToHash([]byte{13, 14, 15}),
		ethCommon.BytesToHash([]byte{16, 17, 18}),
		ethCommon.BytesToHash([]byte{19, 20, 21}),
	}
	succeful := []bool{true, true, true, true, true, false, true}
	inclusive := InclusiveList{
		HashList:   hashes,
		Successful: succeful,
	}
	inclusive.Mode = InclusiveMode_Results
	selList, clearList := inclusive.GetList()
	fmt.Printf(" inclusive.getlist selList=%v\n,clearList=%v\n", selList, clearList)
}
func TestInclusive(t *testing.T) {
	hashes := []ethCommon.Hash{
		ethCommon.BytesToHash([]byte{1, 2, 3}),
		ethCommon.BytesToHash([]byte{4, 5, 6}),
		ethCommon.BytesToHash([]byte{7, 8, 9}),
		ethCommon.BytesToHash([]byte{10, 11, 12}),
		ethCommon.BytesToHash([]byte{13, 14, 15}),
		ethCommon.BytesToHash([]byte{16, 17, 18}),
		ethCommon.BytesToHash([]byte{19, 20, 21}),
	}
	succeful := []bool{true, true, true, true, true, false, true}
	inclusive := InclusiveList{
		HashList:   hashes,
		Successful: succeful,
	}

	data, err := inclusive.GobEncode()
	if err != nil {
		fmt.Printf(" inclusive.GobEncode err=%v\n", err)
		return

	}
	fmt.Printf(" inclusive.GobEncode result=%x\n", data)

	inclusiveResult := InclusiveList{}
	err = inclusiveResult.GobDecode(data)
	if err != nil {
		fmt.Printf(" inclusive.GobDecode err=%v\n", err)
		return

	}
	fmt.Printf(" inclusive.GobDecode result=%v\n", inclusiveResult)
}
