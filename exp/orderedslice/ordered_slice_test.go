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

package orderedslice

import (
	"reflect"
	"testing"
)

func TestOrderedSlice(t *testing.T) {
	dict := make(map[string]int)
	dict["13"] = 7
	dict["12"] = 2
	dict["11"] = 1
	dict["5"] = 3
	dict["2"] = 13
	dict["1"] = 5

	set := NewOrderedSlice[string, int](10, nil, nil, "1", "2", "5", "11", "12", "13")
	set.Append("15")

	if !reflect.DeepEqual(set.Elements, []string{"1", "2", "5", "11", "12", "13", "15"}) {
		t.Error("Error: Key is not equal !", set)
	}

	getter := func(str string) int { return dict[str] }
	greaterEqual := func(v0, v1 int) bool { return v0 >= v1 }

	set = NewOrderedSlice(10, getter, greaterEqual, "1", "2", "5", "11", "12", "13")
	set.Append("15")

	if !reflect.DeepEqual(set.Elements, []string{"15", "11", "12", "5", "1", "13", "2"}) {
		t.Error("Error: Key is not equal !", set)
	}

}

// func TestIndexedSliceDelet(t *testing.T) {
// 	set := NewOrderedSet[string]("", 10, func(str string) [32]byte { return [32]byte{} }, "1", "2", "5", "11", "12", "13")
// 	set.Delete("2", "11")
// 	if !reflect.DeepEqual(set.Elements(), []string{"1", "5", "12", "13"}) {
// 		t.Error("Error: Key is not equal !", set.Elements())
// 	}

// 	k, v := mapi.FindValue(set.dict, func(v0 *int, v1 *int) bool { return *v0 < *v1 })
// 	if len(set.dict) != 4 || *v != 0 || k != "1" {
// 		t.Error("Error: Key is not equal !", set.dict)
// 	}

// 	k, v = mapi.FindValue(set.dict, func(v0 *int, v1 *int) bool { return *v0 > *v1 })
// 	if len(set.dict) != 4 || *v != 3 || k != "13" {
// 		t.Error("Error: Key is not equal !", set.dict)
// 	}

// 	set.Delete("1")
// 	if !reflect.DeepEqual(set.Elements(), []string{"5", "12", "13"}) {
// 		t.Error("Error: Key is not equal !", set.Elements())
// 	}

// 	k, v = mapi.FindValue(set.dict, func(v0 *int, v1 *int) bool { return *v0 < *v1 })
// 	if len(set.dict) != 3 || *v != 0 || k != "5" {
// 		t.Error("Error: Key is not equal !", set.dict)
// 	}

// 	k, v = mapi.FindValue(set.dict, func(v0 *int, v1 *int) bool { return *v0 > *v1 })
// 	if len(set.dict) != 3 || *v != 2 || k != "13" {
// 		t.Error("Error: Key is not equal !", set.dict)
// 	}

// 	set.SetAt(1, "15")
// 	if !reflect.DeepEqual(set.Elements(), []string{"5", "15", "13"}) {
// 		t.Error("Error: Key is not equal !", set.Elements())
// 	}
// }

// type OrderedSlice[T any] struct {
// 	elements  []T
// 	idxGetter func(T) int
// }
