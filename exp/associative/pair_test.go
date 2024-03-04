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
	if !slice.Equal(pairs.Firsts(), []string{"1", "2", "3", "4"}) || !slice.Equal(pairs.Seconds(), []int{1, 2, 3, 4}) {
		t.Error("Error: Values are not equal !")
	}

	_0, _1 := pairs.Split()
	if !slice.Equal(_0, pairs.Firsts()) || !slice.Equal(_1, pairs.Seconds()) {
		t.Error("Error: Values are not equal !")
	}
}
