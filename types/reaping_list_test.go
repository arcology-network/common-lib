package types

import (
	"fmt"
	"math/big"
	"testing"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
)

func TestReapinglist(t *testing.T) {
	hashes := []ethCommon.Hash{
		ethCommon.BytesToHash([]byte{1, 2, 3}),
		ethCommon.BytesToHash([]byte{4, 5, 6}),
		ethCommon.BytesToHash([]byte{7, 8, 9}),
		ethCommon.BytesToHash([]byte{10, 11, 12}),
		ethCommon.BytesToHash([]byte{13, 14, 15}),
		ethCommon.BytesToHash([]byte{16, 17, 18}),
		ethCommon.BytesToHash([]byte{19, 20, 21}),
	}
	reapinglist := ReapingList{
		List:      Arr2Ptr(hashes),
		Timestamp: big.NewInt(12),
	}

	datas, err := reapinglist.GobEncode()
	if err != nil {
		fmt.Printf(" reapinglist.GobEncode err=%v\n", err)
		return

	}
	fmt.Printf(" reapinglist.GobEncode result=%x\n", datas)

	reapinglistResult := ReapingList{}
	err = reapinglistResult.GobDecode(datas)
	if err != nil {
		fmt.Printf(" reapinglist.GobDecode err=%v\n", err)
		return

	}
	fmt.Printf(" reapinglist.GobDecode result=%v\n", reapinglistResult)
}
