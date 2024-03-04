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

package deltaslice

import "testing"

func TestIndexedSliceInt(t *testing.T) {
	// deltaSlice := NewDeltaSlice[string](10)
	strings := []string{}
	slice := NewIndexSlice[[]string, string, int](strings,
		func(strs *[]string, i int) *string {
			if i >= len(*strs) {
				return nil
			}
			return &(*strs)[i]
		},
		func(strs *[]string, k *int, v string) int {
			if k == nil { // New element
				*strs = append(*strs, v) // Append to the end
				return len(*strs) - 1    // Return the index of the new element
			}
			(*strs)[*k] = v
			return *k
		},

		func(strs *[]string) int { return len(*strs) },
	)

	slice.SetByKey(0, "aa")
	if str := slice.GetByKey(0); str == nil || *str != "aa" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := slice.GetByKey(1); str != nil {
		t.Errorf("Expected nil, got %d", str)
	}

	slice.SetByKey(0, "bb")
	if str := slice.GetByKey(0); str == nil || *str != "bb" {
		t.Errorf("Expected gg, got %d", str)
	}

	slice.SetByKey(1, "cc")
	if str := slice.GetByKey(1); str == nil || *str != "cc" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := slice.GetByKey(0); str == nil || *str != "bb" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := slice.GetByIndex(0); str == nil || *str != "bb" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := slice.GetByIndex(1); str == nil || *str != "cc" {
		t.Errorf("Expected gg, got %d", str)
	}
}
