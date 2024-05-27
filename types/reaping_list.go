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
	"math/big"

	codec "github.com/arcology-network/common-lib/codec"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

type ReapingList struct {
	List      []ethCommon.Hash
	Timestamp *big.Int
}

func (rl *ReapingList) GetList() (selectList []ethCommon.Hash, clearList []ethCommon.Hash) {
	selectList = rl.List
	clearList = rl.List
	return
}

func (rl *ReapingList) GobEncode() ([]byte, error) {
	hashArray := rl.List
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
	rl.List = arrs
	if len(fields[1]) > 0 {
		rl.Timestamp = new(big.Int).SetBytes(fields[1])
	}
	return nil
}
