package types

import (
	"math/big"

	ethCommon "github.com/arcology/3rd-party/eth/common"
	encoding "github.com/arcology/common-lib/encoding"
)

type ReapingList struct {
	List      []*ethCommon.Hash
	Timestamp *big.Int
}

func (rl *ReapingList) GetList() (selectList []*ethCommon.Hash, clearList []*ethCommon.Hash) {
	selectList = rl.List
	clearList = rl.List
	return
}

func (rl *ReapingList) GobEncode() ([]byte, error) {
	hashArray := Ptr2Arr(rl.List)
	timeStampData := []byte{}
	if rl.Timestamp != nil {
		timeStampData = rl.Timestamp.Bytes()
	}

	data := [][]byte{
		ethCommon.Hashes(hashArray).Encode(),
		timeStampData,
	}
	return encoding.Byteset(data).Encode(), nil
}
func (rl *ReapingList) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	arrs := []ethCommon.Hash{}
	arrs = ethCommon.Hashes(arrs).Decode(fields[0])
	rl.List = Arr2Ptr(arrs)
	if len(fields[1]) > 0 {
		rl.Timestamp = new(big.Int).SetBytes(fields[1])
	}
	return nil
}
