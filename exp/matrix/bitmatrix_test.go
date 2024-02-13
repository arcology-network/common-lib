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

package matrix

import (
	"testing"
)

func TestBitm(t *testing.T) {
	// Example usage
	bm := NewBitMatrix(11, 11, false)

	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if bm.Get(i, j) {
				t.Error("Should be false")
			}
		}
	}

	bm.Set(3, 4, true)
	if !bm.Get(3, 4) {
		t.Error("failed to write")
	}
	bm.Fill(true)

	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if !bm.Get(i, j) {
				t.Error("Should be true")
			}
		}
	}

	bm.Fill(false)
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if bm.Get(i, j) {
				t.Error("Should be false")
			}
		}
	}

	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			bm.Set(i, j, true)
		}
	}

	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if !bm.Get(i, j) {
				t.Error("Should be true")
			}
		}
	}
}
