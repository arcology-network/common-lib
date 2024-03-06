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
	"math"
	"testing"

	indexedslice "github.com/arcology-network/common-lib/exp/indexed"
	"github.com/arcology-network/common-lib/exp/slice"
)

func TestDeltaSliceBasic(t *testing.T) {
	buffer := indexedslice.NewIndexedSlice[int, int, int](
		func(v int) int { return v },
		func(_ int, newV int) int { return newV },
		func(_ int, newV int, old *int) { *old = newV },
		func(v *int) { *v = math.MaxInt },
		func(old int) bool { return old == math.MaxInt },
	)

	deltaSlice := NewDeltaSet(buffer,
		func(v *int) int { return math.MaxInt },
		func(v int) bool { return v == math.MaxInt },
	)

	deltaSlice.Append(func(v int) int { return v }, 11, 12, 13)
	if !slice.Equal(deltaSlice.appended.Keys(), []int{11, 12, 13}) || !slice.Equal(deltaSlice.appended.Keys(), []int{11, 12, 13}) {
		t.Error("failed to append")
	}

	deltaSlice.Commit()
	if deltaSlice.CommitedLen != 3 || len(deltaSlice.appended.Keys()) != 0 {
		t.Error("failed to commit")
	}

	deltaSlice.SetByIndex(0, math.MaxInt)
	if deltaSlice.CommitedLen != 3 || len(deltaSlice.appended.Keys()) != 0 || len(deltaSlice.modified.Index()) != 1 {
		t.Error("failed to delete")
	}

	deltaSlice.SetByIndex(1, math.MaxInt)
	if deltaSlice.CommitedLen != 3 || len(deltaSlice.appended.Keys()) != 0 || !slice.Equal(deltaSlice.modified.Keys(), []int{11, 12}) {
		t.Error("failed to delete")
	}

	deltaSlice.Append(func(v int) int { return v }, 14, 15, 16, 17)
	if deltaSlice.CommitedLen != 3 || !slice.Equal(deltaSlice.modified.Keys(), []int{11, 12}) || !slice.Equal(deltaSlice.appended.Keys(), []int{14, 15, 16, 17}) {
		t.Error("failed to commit", deltaSlice.modified.Keys())
	}
	deltaSlice.Commit()

	if v, _ := deltaSlice.Get(0); v != 13 {
		t.Error("failed to get", v)
	}

	if v, _ := deltaSlice.Get(1); v != 14 {
		t.Error("failed to get")
	}

	if v, _ := deltaSlice.Get(2); v != 15 {
		t.Error("failed to get")
	}

	if deltaSlice.CommitedLen != 5 || len(deltaSlice.appended.Keys()) != 0 || !slice.Equal(deltaSlice.readonlys.Keys(), []int{13, 14, 15, 16, 17}) {
		t.Error("failed to commit", (deltaSlice.CommitedLen))
	}

	deltaSlice.Append(func(v int) int { return v }, 18, 19)
	if !slice.Equal(deltaSlice.ToSlice(), []int{13, 14, 15, 16, 17, 18, 19}) {
		t.Error("failed to commit", deltaSlice.ToSlice())
	}

	deltaSlice.Delete(0)
	if v, _ := deltaSlice.Get(0); v != math.MaxInt {
		t.Error("failed to get", v)
	}

	deltaSlice.SetByIndex(5, 77)
	if v, _ := deltaSlice.Get(5); v != 77 || !slice.Equal(deltaSlice.readonlys.Keys(), []int{13, 14, 15, 16, 17}) { // Take effect but the original values would be altered
		t.Error("failed to get", v)
	}

	deltaSlice.SetByIndex(5, 99)
	if v, _ := deltaSlice.Get(5); v != 99 || !slice.Equal(deltaSlice.readonlys.Keys(), []int{13, 14, 15, 16, 17}) { // Take effect but the original values would be altered
		t.Error("failed to get", v)
	}

	deltaSlice.SetByIndex(5, math.MaxInt)
	if v, _ := deltaSlice.Get(5); v != math.MaxInt {
		t.Error("failed to get", v)
	}
}

// func TestDeltaSliceDuplicateDelete(t *testing.T) {
// 	buffer := indexedslice.NewIndexedSlice[int, int, int](
// 		func(v int) int { return v },
// 		func(_ int, newV int) int { return newV },
// 		func(_ int, newV int, old *int) { *old = newV },
// 		func(old int) bool { return old == math.MaxInt },
// 	)
// 	deltaSlice := NewDeltaSet(buffer, func(v *int) int { return math.MaxInt }, func(v int) bool { return v == math.MaxInt })

// 	deltaSlice.Append(func(i int) int { return i }, 11, 12, 13)

// 	if deltaSlice.Delete(0); !slice.Equal(deltaSlice.ToSlice(), []int{11, 12, 13}) {
// 		t.Error("failed to append", deltaSlice.ToSlice())
// 	}

// 	if deltaSlice.Delete(0, 1, 2); !slice.Equal(deltaSlice.ToSlice(), []int{11, 12, 13}) {
// 		t.Error("failed to append", deltaSlice.ToSlice())
// 	}

// 	if !slice.Equal(deltaSlice.Values(), []int{}) || !slice.Equal(deltaSlice.Appended(), []int{11, 12, 13}) || len(deltaSlice.Modified().Elements()) != 3 {
// 		t.Error("failed to append", deltaSlice.ToSlice())
// 	}

// 	deltaSlice.SetByIndex(0, 21)
// 	deltaSlice.SetByIndex(1, 31)
// 	deltaSlice.SetByIndex(2, 41)
// 	deltaSlice.SetByIndex(3, 51) // this should take no effect

// 	// Modify the newly appended elements.
// 	if !slice.Equal(deltaSlice.Values(), []int{}) ||
// 		!slice.Equal(deltaSlice.ToSlice(), []int{11, 12, 13}) ||
// 		!slice.Equal(deltaSlice.Appended(), []int{11, 12, 13}) ||
// 		len(deltaSlice.Modified().Elements()) != 3 {
// 		t.Error("failed to append", deltaSlice.ToSlice(), deltaSlice.Appended())
// 	}

// 	deltaSlice.Commit()
// 	if !slice.Equal(deltaSlice.ToSlice(), []int{21, 31, 41}) || !slice.Equal(deltaSlice.Values(), []int{21, 31, 41}) {
// 		t.Error("failed to append", deltaSlice.ToSlice())
// 	}

// 	deltaSlice.Append(func(v int) int { return v }, 51, 61) // {21, 31, 41, 51, 61}
// 	deltaSlice.SetByIndex(0, math.MaxInt)                   // {MaxInt, 31, 41, 51, 61}
// 	deltaSlice.SetByIndex(4, math.MaxInt)                   // {MaxInt, 31, 41, 51, MaxInt}

// 	if !slice.Equal(deltaSlice.ToSlice(), []int{21, 31, 41, 51, 61}) ||
// 		len(deltaSlice.Modified().Elements()) != 2 {
// 		t.Error("failed to append", deltaSlice.ToSlice())
// 	}

// 	deltaSlice.Commit()
// 	if !slice.Equal(deltaSlice.ToSlice(), []int{31, 41, 51}) {
// 		t.Error("failed to append", deltaSlice.ToSlice())
// 	}
// }

// func TestDeltaSliceString(t *testing.T) {
// 	buffer := indexedslice.NewIndexedSlice(
// 		func(k string) string { return k },
// 		func(_ string, v string) string { return v },
// 		func(_ string, v string, old *string) { *old = v },
// 		func(v string) bool { return v == "" },
// 	)

// 	deltaSlice := NewDeltaSet(buffer,
// 		func(v *string) string { return "" },
// 		func(v string) bool { return v == "" })

// 	if deltaSlice.Append(func(k string) string { return k }, "11", "12", "13"); !slice.Equal(deltaSlice.Appended(), []string{"11", "12", "13"}) ||
// 		!slice.Equal(deltaSlice.ToSlice(), []string{"11", "12", "13"}) {
// 		t.Error("failed to append")
// 	}

// 	deltaSlice.SetByIndex(0, "")
// 	if deltaSlice.CommitedLen!= 3 || len(deltaSlice.appended.Keys()) != 0 || len(deltaSlice.modified.Index()) != 1 {
// 		t.Error("failed to delete")
// 	}

// 	deltaSlice.SetByIndex(1, "")
// 	if deltaSlice.CommitedLen!= 3 || len(deltaSlice.appended.Keys()) != 0 || len(deltaSlice.modified.Index()) != 2 {
// 		t.Error("failed to delete")
// 	}

// 	deltaSlice.Append(func(k string) string { return k }, "14", "15", "16", "17")
// 	if deltaSlice.CommitedLen!= 3 || len(deltaSlice.appended.Keys()) != 4 {
// 		t.Error("failed to commit")
// 	}
// 	deltaSlice.Commit()
// }
