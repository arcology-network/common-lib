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

package common

import (
	"testing"
)

func TestMinMax(t *testing.T) {
	if Min(1, 9) != 1 {
		t.Error("Error: Should be 1")
	}

	if Max(1, 9) != 9 {
		t.Error("Error: Should be 9")
	}
}

func TestDoMax(t *testing.T) {
	v, idx := DoMax([]int{1, 7, 3, 7}, func(a, b int) bool { return a < b })
	if v != 7 || idx != 1 {
		t.Errorf("Error: expected (7,1), got (%d,%d)", v, idx)
	}

	v, idx = DoMax([]int{}, func(a, b int) bool { return a < b })
	if v != 0 || idx != -1 {
		t.Errorf("Error: expected (0,-1), got (%d,%d)", v, idx)
	}
}

func TestDoMin(t *testing.T) {
	v, idx := DoMin([]int{9, 1, 3, 1}, func(a, b int) bool { return a < b })
	if v != 1 || idx != 1 {
		t.Errorf("Error: expected (1,1), got (%d,%d)", v, idx)
	}

	v, idx = DoMin([]int{}, func(a, b int) bool { return a < b })
	if v != 0 || idx != -1 {
		t.Errorf("Error: expected (0,-1), got (%d,%d)", v, idx)
	}
}

func TestDoMedian(t *testing.T) {
	entries := []int{9, 1, 7, 3, 5}
	original := append([]int(nil), entries...)

	v, idx := DoMedian(entries, func(a, b int) bool { return a < b })
	if v != 5 || idx != 4 {
		t.Errorf("Error: expected (5,4), got (%d,%d)", v, idx)
	}

	for i := range entries {
		if entries[i] != original[i] {
			t.Errorf("Error: DoMedian mutated input at %d", i)
			break
		}
	}

	v, idx = DoMedian([]int{8, 2, 6, 4}, func(a, b int) bool { return a < b })
	if v != 4 || idx != 3 {
		t.Errorf("Error: expected (4,3), got (%d,%d)", v, idx)
	}

	v, idx = DoMedian([]int{10, 2, 6, 4, 8, 12}, func(a, b int) bool { return a < b })
	if v != 6 || idx != 2 {
		t.Errorf("Error: expected (6,2), got (%d,%d)", v, idx)
	}

	v, idx = DoMedian([]int{}, func(a, b int) bool { return a < b })
	if v != 0 || idx != -1 {
		t.Errorf("Error: expected (0,-1), got (%d,%d)", v, idx)
	}
}
