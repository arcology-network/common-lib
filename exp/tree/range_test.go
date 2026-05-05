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

package tree

import (
	"reflect"
	"sort"
	"testing"

	redblacktree "github.com/emirpasic/gods/v2/trees/redblacktree"
)

func TestRangeHelpers(t *testing.T) {
	btree := redblacktree.New[int, int]()
	for _, value := range []int{5, 2, 8, 1, 3, 7, 9} {
		btree.Put(value, value)
	}

	between := []int{}
	Between(btree.Root, 3, 7, &between)
	sort.Ints(between)
	if !reflect.DeepEqual(between, []int{3, 5, 7}) {
		t.Error("Error: Between should collect values within the inclusive range")
	}

	greater := []int{}
	GreaterOrEqualThan(btree.Root, 7, &greater)
	sort.Ints(greater)
	if !reflect.DeepEqual(greater, []int{7, 8, 9}) {
		t.Error("Error: GreaterOrEqualThan should collect values at or above the threshold")
	}

	less := []int{}
	LessOrEqualThan(btree.Root, 3, &less)
	sort.Ints(less)
	if !reflect.DeepEqual(less, []int{1, 2, 3}) {
		t.Error("Error: LessOrEqualThan should collect values at or below the threshold")
	}

	empty := []int{}
	Between[int, int](nil, 0, 1, &empty)
	if !reflect.DeepEqual(empty, []int{}) {
		t.Error("Error: range helpers should ignore nil nodes")
	}
}
