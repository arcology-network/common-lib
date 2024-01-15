package types

import (
	"math/big"

	codec "github.com/arcology-network/common-lib/codec"
	ethCommon "github.com/ethereum/go-ethereum/common"
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
		Hashes(hashArray).Encode(),
		timeStampData,
	}
	return codec.Byteset(data).Encode(), nil
}
func (rl *ReapingList) GobDecode(data []byte) error {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	arrs := []ethCommon.Hash{}
	arrs = Hashes(arrs).Decode(fields[0])
	rl.List = Arr2Ptr(arrs)
	if len(fields[1]) > 0 {
		rl.Timestamp = new(big.Int).SetBytes(fields[1])
	}
	return nil
}
