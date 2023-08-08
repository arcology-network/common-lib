package types

import (
	"math/big"

	encoding "github.com/arcology-network/common-lib/encoding"
	evmCommon "github.com/arcology-network/evm/common"
)

type ReapingList struct {
	List      []*evmCommon.Hash
	Timestamp *big.Int
}

func (rl *ReapingList) GetList() (selectList []*evmCommon.Hash, clearList []*evmCommon.Hash) {
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
	return encoding.Byteset(data).Encode(), nil
}
func (rl *ReapingList) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	arrs := []evmCommon.Hash{}
	arrs = Hashes(arrs).Decode(fields[0])
	rl.List = Arr2Ptr(arrs)
	if len(fields[1]) > 0 {
		rl.Timestamp = new(big.Int).SetBytes(fields[1])
	}
	return nil
}
