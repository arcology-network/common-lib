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

// "github.com/HPISTechnologies/common-lib/common"

func TestParallelFor(t *testing.T) {
	vals := []int{1, 2, 3, 4}
	ParallelFor(0, len(vals), 1000, func(i int) {
		vals[i] = vals[i] * 2
	})

	if vals[0] != 2 || vals[1] != 4 || vals[2] != 6 || vals[3] != 8 {
		t.Error("ParallelFor Failed")
	}
}

func TestTrimLeft(t *testing.T) {
	s := []int{0, 0, 1, 2, 0}
	result := TrimLeft(s, 0)
	expected := []int{1, 2, 0}
	if len(result) != len(expected) {
		t.Errorf("TrimLeft failed: expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("TrimLeft failed at index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}

	// Test with no cutset at the start
	s2 := []int{1, 2, 0}
	result2 := TrimLeft(s2, 0)
	if len(result2) != len(s2) {
		t.Errorf("TrimLeft failed: expected length %d, got %d", len(s2), len(result2))
	}

	s = []int{0}
	result = TrimLeft(s, 0)
	if len(result) != 0 {
		t.Errorf("TrimLeft failed: expected empty slice, got %d elements", len(result))
	}
}

func TestTrimRight(t *testing.T) {
	s := []int{0, 1, 2, 0, 0}
	result := TrimRight(s, 0)
	expected := []int{0, 1, 2}
	if len(result) != len(expected) {
		t.Errorf("TrimRight failed: expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("TrimRight failed at index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}

	// Test with no cutset at the end
	s2 := []int{0, 1, 2}
	result2 := TrimRight(s2, 0)
	if len(result2) != len(s2) {
		t.Errorf("TrimRight failed: expected length %d, got %d", len(s2), len(result2))
	}
}
