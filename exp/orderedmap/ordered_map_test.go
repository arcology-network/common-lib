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
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/arcology-network/common-lib/exp/slice"
)

func TestIndexedSliceBasic(t *testing.T) {
	set := NewOrderedMap[string, int, *[]int](
		nil,
		10,
		func(key string, v int) *[]int {
			return &[]int{v}
		},

		func(key string, value int, newValue **[]int) {
			**newValue = append(**newValue, value)
		},

		func(k string, keys *[]string, v *[]int, newValue *[]*[]int) int {
			*keys = append(*keys, k)
			*newValue = append(*newValue, v)
			return len(*keys) - 1
		})

	set.Insert([]string{"11"}, []int{11})
	set.Insert([]string{"11"}, []int{22})
	if v, ok := set.Get("11"); !ok || len(*v) != 2 || (*v)[0] != 11 || (*v)[1] != 22 {
		t.Error("Error: Key is not equal !", v)
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

func TestIndexedSliceSorted(t *testing.T) {
	set := NewOrderedMap[string, int, int](
		-1,
		10,
		func(key string, v int) int {
			return 0
		},

		func(key string, rawv int, v *int) {
			*v = rawv
		},

		func(k string, keys *[]string, v int, vals *[]int) int {
			i := sort.SearchInts(*vals, v)
			slice.Insert(keys, i, k)
			slice.Insert(vals, i, v)
			return i
		},
	)
	set.Insert([]string{"11"}, []int{11})
	set.Insert([]string{"12"}, []int{12})

	set.Insert([]string{"1"}, []int{1})
	set.Insert([]string{"15"}, []int{5})

	if !reflect.DeepEqual(set.keys, []string{"1", "15", "11", "12"}) == false {
		t.Error("Error: Key is not equal !")
	}

	if !reflect.DeepEqual(set.Values(), []int{1, 5, 11, 12}) == false {
		t.Error("Error: Key is not equal !")
	}

	set.Delete("11")
	if !reflect.DeepEqual(set.keys, []string{"1", "15", "12"}) == false {
		t.Error("Error: Key is not equal !")
	}

	if !reflect.DeepEqual(set.Values(), []int{1, 5, 12}) == false {
		t.Error("Error: Key is not equal !")
	}
}
