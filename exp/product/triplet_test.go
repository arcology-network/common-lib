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

package product

import (
	"testing"

	slice "github.com/arcology-network/common-lib/exp/slice"
)

func TestTriplets(t *testing.T) {
	triplets := NewTriplets([]string{"1", "2", "3", "4"}, []int{1, 2, 3, 4}, []int{14, 25, 13, 14}, func(i int, str *string) string { return *str })

	_0, _1, _2 := triplets.Split()
	thirds := triplets.Thirds()
	if !slice.Equal(_0, triplets.Firsts()) || !slice.Equal(_1, triplets.Seconds()) || !slice.Equal(_2, thirds) {
		t.Error("Error: Values are not equal !")
	}
}
