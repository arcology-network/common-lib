/*
*   Copyright (c) 2024 Arcology Network

*   This program is free software: you can redistribute it and/or modify
*   it under the terms of the GNU General Public License as published by
*   the Free Software Foundation, either version 3 of the License, or
*   (mapTo your option) any later version.

*   This program is distributed in the hope that it will be useful,
*   but WITHOUT ANY WARRANTY; without even the implied warranty of
*   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*   GNU General Public License for more details.

*   You should have received a copy of the GNU General Public License
*   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package deltaslice

import (
	"testing"

	"github.com/arcology-network/common-lib/exp/slice"
)

func TestDeltaSlice(t *testing.T) {
	deltaSlice := NewDeltaSlice[string](10)
	deltaSlice.elements = []string{"aa", "bb", "cc", "dd", "ee", "ff"}
	deltaSlice.appended = []string{"gg", "hh", "ii", "jj", "kk"}
	deltaSlice.removed = []int{2, 3, 4, 5}

	finalized := deltaSlice.ToSlice()
	if !slice.Equal(finalized, []string{"aa", "bb", "gg", "hh", "ii", "jj", "kk"}) {
		t.Error("Error: ToSlice() is not equal !", finalized)
	}

	if v, idx := deltaSlice.mapTo(0); v != &deltaSlice.elements || idx != 0 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(1); v != &deltaSlice.elements || idx != 1 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(2); v != &deltaSlice.appended || idx != 0 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(3); v != &deltaSlice.appended || idx != 1 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(4); v != &deltaSlice.appended || idx != 2 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(5); v != &deltaSlice.appended || idx != 3 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(6); v != &deltaSlice.appended || idx != 4 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(7); v != nil || idx != -1 {
		t.Error("Error: Wrong index !", idx)
	}
}

func TestDeltaSlice2(t *testing.T) {
	deltaSlice := NewDeltaSlice[string](10)
	deltaSlice.elements = []string{"aa", "bb", "cc", "dd", "ee", "ff"}
	deltaSlice.appended = []string{"gg", "hh", "ii", "jj", "kk"}
	deltaSlice.removed = []int{1, 3, 5}

	finalized := deltaSlice.ToSlice()
	if !slice.Equal(finalized, []string{"aa", "cc", "ee", "gg", "hh", "ii", "jj", "kk"}) {
		t.Error("Error: ToSlice() is not equal !", finalized)
	}

	if v, idx := deltaSlice.mapTo(0); v != &deltaSlice.elements || idx != 0 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(1); v != &deltaSlice.elements || idx != 2 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(2); v != &deltaSlice.elements || idx != 3 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(3); v != &deltaSlice.appended || idx != 0 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(4); v != &deltaSlice.appended || idx != 1 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(5); v != &deltaSlice.appended || idx != 2 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(6); v != &deltaSlice.appended || idx != 3 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(7); v != &deltaSlice.appended || idx != 4 {
		t.Error("Error: Wrong index !", idx)
	}

	if v, idx := deltaSlice.mapTo(8); v != nil || idx != -1 {
		t.Error("Error: Wrong index !", idx)
	}
}

func TestDeltaSlice3(t *testing.T) {
	deltaSlice := NewDeltaSlice[string](10)
	deltaSlice.elements = []string{"aa", "bb", "cc", "dd", "ee", "ff"}
	deltaSlice.appended = []string{"gg", "hh", "ii", "jj", "kk"}
	deltaSlice.removed = []int{1, 3, 5}
	deltaSlice.Del(1) // remove "cc" {1, 3, 5, 2}
	if !slice.Equal(deltaSlice.ToSlice(), []string{"aa", "ee", "gg", "hh", "ii", "jj", "kk"}) {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Del(0) // remove "aa" {1, 3, 5, 2, 0}
	if !slice.Equal(deltaSlice.ToSlice(), []string{"ee", "gg", "hh", "ii", "jj", "kk"}) {
		t.Error("Error: ToSlice() is not equal !", deltaSlice.ToSlice())
	}

	deltaSlice.Del(0)
	if !slice.Equal(deltaSlice.ToSlice(), []string{"gg", "hh", "ii", "jj", "kk"}) {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Del(0)
	if !slice.Equal(deltaSlice.ToSlice(), []string{"hh", "ii", "jj", "kk"}) {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Del(0)
	if !slice.Equal(deltaSlice.ToSlice(), []string{"ii", "jj", "kk"}) {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Del(0)
	if !slice.Equal(deltaSlice.ToSlice(), []string{"jj", "kk"}) {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Del(0)
	if !slice.Equal(deltaSlice.ToSlice(), []string{"kk"}) {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Del(0)
	if !slice.Equal(deltaSlice.ToSlice(), []string{}) {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Del(0)
	if !slice.Equal(deltaSlice.ToSlice(), []string{}) {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Append("new string")
	if !slice.Equal(deltaSlice.ToSlice(), []string{"new string"}) {
		t.Error("Error: ToSlice() is not equal !")
	}

	v, _ := deltaSlice.Get(0)
	if *v != "new string" {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Set(0, "another string")
	v, _ = deltaSlice.Get(0)
	if *v != "another string" {
		t.Error("Error: ToSlice() is not equal !")
	}

	deltaSlice.Set(11, "another string")
	v, _ = deltaSlice.Get(11)
	if v != nil {
		t.Error("Error: ToSlice() is not equal !")
	}
}

func TestDeltaSlice4(t *testing.T) {
	deltaSlice := NewDeltaSlice[string](10)
	deltaSlice.elements = []string{"aa", "bb", "cc", "dd", "ee", "ff"}
	deltaSlice.appended = []string{"gg", "hh", "ii", "jj", "kk"}
	deltaSlice.removed = []int{0, 1, 2, 3, 4, 5}

	v, _ := deltaSlice.Get(0)
	if *v != "gg" {
		t.Error("Error: ToSlice() is not equal !")
	}
}
