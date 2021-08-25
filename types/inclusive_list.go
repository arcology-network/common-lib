package types

import (
	ethCommon "github.com/arcology/3rd-party/eth/common"
	"github.com/arcology/common-lib/common"
	encoding "github.com/arcology/common-lib/encoding"
)

const (
	InclusiveMode_Message = byte(0) //used by tx  and message in scheduler/core
	InclusiveMode_Results = byte(1) //used by euresults and receipt in exec/arbitrate/eshing/generic-hashing/storage
)

type InclusiveList struct {
	HashList   []*ethCommon.Hash
	Successful []bool
	Mode       byte
}

func (il *InclusiveList) CopyListAddHeight(height, round uint64) *InclusiveList {
	hashList := make([]*ethCommon.Hash, len(il.HashList))

	for i := range il.HashList {
		hash := il.HashList[i]
		newhash := common.ToNewHash(*hash, height, round)
		hashList[i] = &newhash
	}

	return &InclusiveList{
		Successful: il.Successful,
		Mode:       il.Mode,
		HashList:   hashList,
	}
}
func (il InclusiveList) GetList() (selectList []*ethCommon.Hash, clearList []*ethCommon.Hash) {
	selectList = make([]*ethCommon.Hash, 0, len(il.HashList))
	clearList = make([]*ethCommon.Hash, 0, len(il.HashList))
	for i, hashItem := range il.HashList {
		if hashItem == nil {
			continue
		}

		switch il.Mode {
		case InclusiveMode_Message:
			selectList = append(selectList, hashItem)
		case InclusiveMode_Results:
			if il.Successful[i] {
				selectList = append(selectList, hashItem)
			}
		}

		clearList = append(clearList, hashItem)

	}
	return
}

func (il *InclusiveList) GobEncode() ([]byte, error) {
	hashArray := Ptr2Arr(il.HashList)
	data := [][]byte{
		ethCommon.Hashes(hashArray).Encode(),
		encoding.Bools(il.Successful).Encode(),
	}
	return encoding.Byteset(data).Encode(), nil
}
func (il *InclusiveList) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	arrs := ethCommon.Hashes([]ethCommon.Hash{}).Decode(fields[0])
	il.Successful = encoding.Bools(il.Successful).Decode(fields[1])
	il.HashList = Arr2Ptr(arrs)
	return nil
}

func Ptr2Arr(array []*ethCommon.Hash) []ethCommon.Hash {
	hashArray := make([]ethCommon.Hash, len(array))
	for i := range array {
		hashArray[i] = *array[i]
	}
	return hashArray
}

func Arr2Ptr(array []ethCommon.Hash) []*ethCommon.Hash {
	hashArray := make([]*ethCommon.Hash, len(array))
	for i := range array {
		hashArray[i] = &array[i]
	}
	return hashArray
}
