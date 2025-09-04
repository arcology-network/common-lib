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

package orderedmap

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/arcology-network/common-lib/exp/slice"
)

func TestIndexedSliceInt(t *testing.T) {
	set := NewOrderedMap(
		nil,
		10,
		func(key string, v int) *[]int {
			return &[]int{v}
		},
		func(key string, value int, newValue **[]int) {
			**newValue = append(**newValue, value)
		})

	set.Insert([]string{"11"}, []int{11})
	set.Insert([]string{"11"}, []int{22})
	if v, ok := set.Get("11"); !ok || len(*v) != 2 || (*v)[0] != 11 || (*v)[1] != 22 {
		t.Error("Error: Key is not equal !")
	}

	set.Insert([]string{"33"}, []int{33})
	set.Insert([]string{"33"}, []int{66})
	if v, ok := set.Get("33"); !ok || len(*v) != 2 || (*v)[0] != 33 || (*v)[1] != 66 {
		t.Error("Error: Key is not equal !")
	}

	if set.keys[0] != "11" || set.keys[1] != "33" || len(set.keys) != 2 {
		t.Error("Error: Key is not equal !")
	}

	length := 0
	set.ParallelForeachDo(func(k string, v **[]int) {
		length += len(**v)
	})

	if length != 4 {
		t.Error("Error: Length is not equal !")
	}

	set.ParallelForeachDo(func(k string, v **[]int) {
		slice.Foreach(**v, func(i int, value *int) {
			(**v)[i] = *value * 2
		})
	})

	if v, ok := set.Get("11"); !ok || len(*v) != 2 || (*v)[0] != 22 || (*v)[1] != 44 {
		t.Error("Error: Key is not equal !")
	}

	if v, ok := set.Get("33"); !ok || len(*v) != 2 || (*v)[0] != 66 || (*v)[1] != 132 {
		t.Error("Error: Key is not equal !")
	}

	x := 6
	a := []int{1, 2, 3, 4, 55}
	i := sort.Search(len(a), func(i int) bool { return a[i] >= x })
	slice.Insert(&a, i, x)
	fmt.Println(a)
}

func TestIndexedSliceBytes(t *testing.T) {
	set := NewOrderedMap(
		nil,
		10,
		func(key string, v []byte) *[][]byte {
			return &[][]byte{v}
		},
		func(key string, value []byte, newValue **[][]byte) {
			**newValue = append(**newValue, value)
		})

	set.Insert([]string{"11"}, [][]byte{{11}})
	set.Insert([]string{"11"}, [][]byte{{22}})
	if v, ok := set.Get("11"); !ok || len(*v) != 2 || bytes.Equal((*v)[0], []byte("11")) || bytes.Equal((*v)[1], []byte("22")) {
		t.Error("Error: Key is not equal !")
	}
}
