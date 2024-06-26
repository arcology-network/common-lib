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
	"fmt"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
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

	data, err := rcptList.GobEncode()
	if err != nil {
		fmt.Printf(" rcptList.GobEncode err=%v\n", err)
		return

	}
	fmt.Printf(" rcptList.GobEncode result=%x\n", data)

	rcptListResult := ReceiptHashList{}
	err = rcptListResult.GobDecode(data)
	if err != nil {
		fmt.Printf(" rcptList.GobDecode err=%v\n", err)
		return

	}
	fmt.Printf(" rcptList.GobDecode result=%v\n", rcptListResult)
}
