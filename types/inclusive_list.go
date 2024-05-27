/*
 *   Copyright (c) 2024 Arcology Network

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package types

import (
	codec "github.com/arcology-network/common-lib/codec"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	InclusiveMode_Message = byte(0) //used by tx  and message in scheduler/core
	InclusiveMode_Results = byte(1) //used by euresults and receipt in exec/arbitrate/eshing/generic-hashing/storage
)

type InclusiveList struct {
	HashList      []ethCommon.Hash
	Successful    []bool
	Mode          byte
	GenerationIdx uint32
}

func (il *InclusiveList) CopyListAddHeight(height, round uint64) *InclusiveList {
	hashList := make([]ethCommon.Hash, len(il.HashList))

	// for i := range il.HashList {
	// 	// hash := il.HashList[i]
	// 	//newhash := common.ToNewHash(*hash, height, round)
	// 	// hashList[i] = &newhash
	// }

	return &InclusiveList{
		Successful: il.Successful,
		Mode:       il.Mode,
		HashList:   hashList,
	}
}
func (il InclusiveList) GetList() (selectList []ethCommon.Hash, clearList []ethCommon.Hash) {
	selectList = make([]ethCommon.Hash, 0, len(il.HashList))
	clearList = make([]ethCommon.Hash, 0, len(il.HashList))
	for i, hashItem := range il.HashList {
		// if hashItem == nil {
		// 	continue
		// }

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
	hashArray := il.HashList
	data := [][]byte{
		Hashes(hashArray).Encode(),
		codec.Bools(il.Successful).Encode(),
		codec.Uint32(il.GenerationIdx).Encode(),
	}
	return codec.Byteset(data).Encode(), nil
}
func (il *InclusiveList) GobDecode(data []byte) error {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	arrs := Hashes([]ethCommon.Hash{}).Decode(fields[0])
	il.Successful = codec.Bools(il.Successful).Decode(fields[1]).(codec.Bools)
	il.GenerationIdx = uint32(codec.Uint32(il.GenerationIdx).Decode(fields[2]).(codec.Uint32))
	il.HashList = arrs
	return nil
}
