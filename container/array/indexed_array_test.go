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
package array

import (
	"fmt"
	"testing"
)

func TestIndexedArray(t *testing.T) {
	indexed := NewIndexedArray[int, string, *[]int](
		func(i int) string { return fmt.Sprint(i) },

		func(i int, v *[]int) *[]int {
			if v == nil {
				return &[]int{i}
			}
			*v = append(*v, i)
			return v
		},

		func(v *[]int) int {
			return len(*v)
		},
	)

	// 2 elements per block, 64 blocks
	indexed.InsertSlice([]int{1, 2, 5, 5, 5})

	if indexed.Length() != 5 || indexed.UniqueLength() != 3 {
		t.Error("Error: Size is not equal !")
	}

	indexed.Insert(int(2))
	if indexed.Length() != 6 {
		t.Error("Error: Size is not equal !")
	}

	indexed.Insert(int(2))
	if indexed.Length() != 7 || indexed.UniqueLength() != 3 {
		t.Error("Error: Size is not equal !")
	}

	arr := indexed.Find(int(5))
	if len(*arr) != 3 {
		t.Error("Error: Size is not equal !")
	}

	arr = indexed.Find(int(2))
	if len(*arr) != 3 {
		t.Error("Error: Size is not equal !")
	}
}
