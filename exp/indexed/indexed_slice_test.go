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

package indexedslice

import (
	"fmt"
	"testing"
)

func TestIndexedSlice(t *testing.T) {
	indexed := NewIndexedSlice[string, int, *[]int](
		func(i int) string { return fmt.Sprint(i) },
		func(_ string, v int) *[]int { return &[]int{v} },
		func(_ string, v int, numbers **[]int) { // Array
			**numbers = append(**numbers, v)
		},
		func(numbers *[]int) bool { return numbers == nil },
	)
	getSize := func(v *[]int) int { return len(*v) }

	// 2 elements per block, 64 blocks
	if indexed.Add(1, 2, 5, 5, 5); indexed.Length(getSize) != 5 {
		t.Error("Error: Size is not equal !")
	}

	if indexed.Add(2, 2, 5, 2); indexed.Length(getSize) != 9 {
		t.Error("Error: Size is not equal !")
	}

	if indexed.Clear(); indexed.Length(getSize) != 0 {
		t.Error("Error: Size is not equal !")
	}

	// Start over again
	if indexed.Add(1, 2, 5, 5, 5); indexed.Length(getSize) != 5 {
		t.Error("Error: Size is not equal !")
	}

	keys := indexed.Keys() // Unique keys only
	if len(keys) != 3 || keys[0] != "1" || keys[1] != "2" || keys[2] != "5" {
		t.Error("Error: Keys are not equal !")
	}

	vals := indexed.Values() // Unique keys only
	if len(*vals[0])+len(*vals[1])+len(*vals[2]) != 5 {
		t.Error("Error: Keys are not equal !")
	}

	indexed.Add(6)
	keys = indexed.Keys() // Unique keys only
	if len(keys) != 4 || keys[3] != "6" {
		t.Error("Error: Keys are not equal !")
	}

	// Get by key
	if v, ok := indexed.GetByKey(1); !ok || len(*v) != 1 {
		t.Error("Error: Value is not equal !", (*v))
	}
}
