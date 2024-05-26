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

package associative

import (
	"testing"

	slice "github.com/arcology-network/common-lib/exp/slice"
)

func TestPairs(t *testing.T) {
	pairs := new(Pairs[string, int]).From([]string{"1", "2", "3", "4"}, []int{1, 2, 3, 4}, func(i int, str *string) string { return *str })
	if !slice.EqualSet(pairs.Firsts(), []string{"1", "2", "3", "4"}) || !slice.EqualSet(pairs.Seconds(), []int{1, 2, 3, 4}) {
		t.Error("Error: Values are not equal !")
	}

	_0, _1 := pairs.Split()
	if !slice.EqualSet(_0, pairs.Firsts()) || !slice.EqualSet(_1, pairs.Seconds()) {
		t.Error("Error: Values are not equal !")
	}

	arr := []*Pair[int, string]{
		{First: 1, Second: "str2"},
		// {First: 2, Second: "str3"},
	}

	pairs3 := Pairs[int, string](arr)

	movedp := slice.MoveIf((*[]*Pair[int, string])(&pairs3), func(_ int, v *Pair[int, string]) bool {
		return v.First == 1
	})

	if len(movedp) != 1 || len(pairs3) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}
}
