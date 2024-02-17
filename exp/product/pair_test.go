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

	"github.com/arcology-network/common-lib/exp/array"
)

func TestPairs(t *testing.T) {
	pairs := new(Pairs[string, int]).From([]string{"1", "2", "3", "4"}, []int{1, 2, 3, 4})
	if !array.Equal(pairs.Firsts(), []string{"1", "2", "3", "4"}) || !array.Equal(pairs.Seconds(), []int{1, 2, 3, 4}) {
		t.Error("Error: Failed to remove nil values !")
	}
}
