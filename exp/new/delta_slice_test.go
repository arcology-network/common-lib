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
		func(i int) int { return i },
		func(_ int, newV int) int { return newV },
		func(_ int, newV int, old *int) { *old = newV },
		func(old int) bool { return old == math.MaxInt },
	)
	deltaSlice := NewDeltaSlice(buffer,
		func(v *int) { *v = math.MaxInt },
		func(v int) bool { return v == math.MaxInt })

	if deltaSlice.Append(11, 12, 13); !slice.Equal(deltaSlice.Appended(), []int{11, 12, 13}) || !slice.Equal(deltaSlice.ToSlice(), []int{11, 12, 13}) {
		t.Error("failed to append")
	}

	if deltaSlice.Commit(); len(deltaSlice.readonlyValues) != 3 || len(deltaSlice.appended) != 0 {
		t.Error("failed to commit")
	}

	deltaSlice.Set(0, func(v *int) { *v = math.MaxInt })
	if len(deltaSlice.readonlyValues) != 3 || len(deltaSlice.appended) != 0 || len(deltaSlice.modified.Index()) != 1 {
		t.Error("failed to delete")
	}

	deltaSlice.Set(1, func(v *int) { *v = math.MaxInt })
	if len(deltaSlice.readonlyValues) != 3 || len(deltaSlice.appended) != 0 || len(deltaSlice.modified.Index()) != 2 {
		t.Error("failed to delete")
	}

	deltaSlice.Append(14, 15, 16, 17)
	if len(deltaSlice.readonlyValues) != 3 || len(deltaSlice.appended) != 4 {
		t.Error("failed to commit")
	}
	deltaSlice.Commit()

	if v, _ := deltaSlice.Get(0); v != 13 {
		t.Error("failed to get")
	}
	if v, _ := deltaSlice.Get(1); v != 14 {
		t.Error("failed to get")
	}
	if v, _ := deltaSlice.Get(2); v != 15 {
		t.Error("failed to get")
	}

	if len(deltaSlice.readonlyValues) != 5 || len(deltaSlice.appended) != 0 {
		t.Error("failed to commit", len(deltaSlice.readonlyValues))
	}

	if !slice.Equal(deltaSlice.readonlyValues, []int{13, 14, 15, 16, 17}) {
		t.Error("failed to commit", deltaSlice.readonlyValues)
	}

	deltaSlice.Append(18, 19)
	if !slice.Equal(deltaSlice.ToSlice(), []int{13, 14, 15, 16, 17, 18, 19}) {
		t.Error("failed to commit", deltaSlice.ToSlice())
	}

	if !slice.Equal(deltaSlice.Values(), []int{13, 14, 15, 16, 17}) {
		t.Error("failed to commit", deltaSlice.ToSlice())
	}

	deltaSlice.Delete(0)
	if v, _ := deltaSlice.Get(0); v != math.MaxInt {
		t.Error("failed to get", v)
	}

	deltaSlice.Set(5, func(v *int) { *v = 77 })
	if v, _ := deltaSlice.Get(5); v != 77 || !slice.Equal(deltaSlice.readonlyValues, []int{13, 14, 15, 16, 17}) { // Take effect but the original values would be altered
		t.Error("failed to get", v)
	}

	deltaSlice.Delete(5)
	if v, _ := deltaSlice.Get(5); v != math.MaxInt {
		t.Error("failed to get", v)
	}
}

func TestDeltaSliceDuplicateDelete(t *testing.T) {
	buffer := indexedslice.NewIndexedSlice[int, int, int](
		func(i int) int { return i },
		func(_ int, newV int) int { return newV },
		func(_ int, newV int, old *int) { *old = newV },
		func(old int) bool { return old == math.MaxInt },
	)
	deltaSlice := NewDeltaSlice(buffer, func(v *int) { *v = math.MaxInt }, func(v int) bool { return v == math.MaxInt })

	deltaSlice.Append(11, 12, 13)
	if deltaSlice.Delete(0); !slice.Equal(deltaSlice.ToSlice(), []int{11, 12, 13}) {
		t.Error("failed to append", deltaSlice.ToSlice())
	}

	if deltaSlice.Delete(0, 1, 2); !slice.Equal(deltaSlice.ToSlice(), []int{11, 12, 13}) {
		t.Error("failed to append", deltaSlice.ToSlice())
	}

	if !slice.Equal(deltaSlice.Values(), []int{}) || !slice.Equal(deltaSlice.Appended(), []int{11, 12, 13}) || len(deltaSlice.Modified().Elements()) != 3 {
		t.Error("failed to append", deltaSlice.ToSlice())
	}

	deltaSlice.Set(0, func(v *int) { *v = 21 })
	deltaSlice.Set(1, func(v *int) { *v = 31 })
	deltaSlice.Set(2, func(v *int) { *v = 41 })
	deltaSlice.Set(3, func(v *int) { *v = 51 }) // this should take no effect

	// Modify the newly appended elements.
	if !slice.Equal(deltaSlice.Values(), []int{}) ||
		!slice.Equal(deltaSlice.ToSlice(), []int{11, 12, 13}) ||
		!slice.Equal(deltaSlice.Appended(), []int{11, 12, 13}) ||
		len(deltaSlice.Modified().Elements()) != 3 {
		t.Error("failed to append", deltaSlice.ToSlice(), deltaSlice.Appended())
	}

	deltaSlice.Commit()
	if !slice.Equal(deltaSlice.ToSlice(), []int{21, 31, 41}) || !slice.Equal(deltaSlice.Values(), []int{21, 31, 41}) {
		t.Error("failed to append", deltaSlice.ToSlice())
	}

	deltaSlice.Append(51, 61)                            // {21, 31, 41, 51, 61}
	deltaSlice.Set(0, func(v *int) { *v = math.MaxInt }) // {MaxInt, 31, 41, 51, 61}
	deltaSlice.Set(4, func(v *int) { *v = math.MaxInt }) // {MaxInt, 31, 41, 51, MaxInt}

	if !slice.Equal(deltaSlice.ToSlice(), []int{21, 31, 41, 51, 61}) ||
		len(deltaSlice.Modified().Elements()) != 2 {
		t.Error("failed to append", deltaSlice.ToSlice())
	}

	deltaSlice.Commit()
	if !slice.Equal(deltaSlice.ToSlice(), []int{31, 41, 51}) {
		t.Error("failed to append", deltaSlice.ToSlice())
	}
}
