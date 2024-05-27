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
