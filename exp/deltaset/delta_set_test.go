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
package deltaset

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/common"
)

func TestDeltaSliceBasic(t *testing.T) {
	deltaSet := NewDeltaSet(-1, 100, nil)

	if deltaSet.InsertBatch([]int{11, 12, 13}); !reflect.DeepEqual(deltaSet.committed.Elements(), []int{}) ||
		!reflect.DeepEqual(deltaSet.added.Elements(), []int{11, 12, 13}) {
		t.Error("failed to append", deltaSet.added.IsDirty(), deltaSet.committed.Elements(), deltaSet.added.Elements())
	}

	if deltaSet.Commit(nil); deltaSet.committed.Length() != 3 ||
		(deltaSet.added.Length()) != 0 ||
		(deltaSet.removed.Length()) != 0 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if deltaSet.DeleteBatch([]int{12}); deltaSet.committed.Length() != 3 ||
		(deltaSet.added.Length()) != 0 || (deltaSet.removed.Length()) != 1 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if deltaSet.Commit(nil); deltaSet.committed.Length() != 2 || (deltaSet.added.Length()) != 0 || deltaSet.added.IsDirty() ||
		(deltaSet.removed.Length()) != 0 { // {11, 13}
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// {11, 13} + {15, 16, 17}
	if deltaSet.InsertBatch([]int{15, 16, 17}); deltaSet.committed.Length() != 2 || (deltaSet.added.Length()) != 3 || (deltaSet.removed.Length()) != 0 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// {11, 13} + {15, 16, 17}
	if deltaSet.Delete(16); deltaSet.committed.Length() != 2 ||
		(deltaSet.added.Length()) != 3 ||
		// !deltaSet.added.IsDirty() ||
		!reflect.DeepEqual(deltaSet.added.Elements(), []int{15, 16, 17}) ||
		(deltaSet.removed.Length()) != 1 {
		t.Error("failed to commit", deltaSet.added.Elements(), deltaSet.committed.Elements())
	}

	if deltaSet.DeleteBatch([]int{11}); deltaSet.committed.Length() != 2 ||
		(deltaSet.added.Length()) != 3 ||
		!reflect.DeepEqual(deltaSet.added.Elements(), []int{15, 16, 17}) ||
		(deltaSet.removed.Length()) != 2 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// Re-delete a uncommitted entry, the removed set will grow but neither the committed nor added set not change.
	if deltaSet.Delete(16); deltaSet.committed.Length() != 2 ||
		(deltaSet.added.Length()) != 3 ||
		!reflect.DeepEqual(deltaSet.added.Elements(), []int{15, 16, 17}) ||
		(deltaSet.removed.Length()) != 2 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if !reflect.DeepEqual(deltaSet.Elements(), []int{13, 15, 17}) {
		t.Error("Failed to get Elements()", deltaSet.Elements())
	}

	if deltaSet.Commit(nil); !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17}) || (deltaSet.added.Length()) != 0 || (deltaSet.removed.Length()) != 0 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.InsertBatch([]int{18, 19, 20, 21})

	if k, _, _, ok := deltaSet.Search(0); !ok || k != 13 {
		t.Error("failed to commit", k)
	}

	if k, _, _, ok := deltaSet.Search(1); !ok || k != 15 {
		t.Error("failed to commit", k)
	}

	if k, _, _, ok := deltaSet.Search(2); !ok || k != 17 {
		t.Error("failed to commit", k)
	}

	if k, _, _, ok := deltaSet.Search(3); !ok || k != 18 {
		t.Error("failed to commit", k)
	}
	if k, _, _, ok := deltaSet.Search(4); !ok || k != 19 {
		t.Error("failed to commit", k)
	}
	if k, _, _, ok := deltaSet.Search(5); !ok || k != 20 {
		t.Error("failed to commit", k)
	}

	deltaSet.DeleteByIndex(1) // After { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.DeleteByIndex(4) // After { 13, 15, 17} + { 18, 19, 21}
	deltaSet.DeleteByIndex(5) // will remove { 13, 15, 17} + { 18, 19}

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17}) || !reflect.DeepEqual(deltaSet.added.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.committed.Elements(), deltaSet.added.Elements())
	}
}

func TestDeltaSliceAddThenDelete(t *testing.T) {

	deltaSet := NewDeltaSet(-1, 100, nil)
	deltaSet.InsertBatch([]int{13, 15, 17})
	deltaSet.Commit(nil)

	deltaSet.InsertBatch([]int{18, 19, 20, 21}) // { 13, 15, 17} + { 18, 19, 20, 21}

	// deltaSet.Delete([]int{118, 119, 210, 211}) // non-existing entries, should not affect the set
	// if deltaSet.removed.Length() != 0 {
	// 	t.Error("Deleting non-existing elements should not affect the set", deltaSet.removed.Elements())
	// }

	deltaSet.DeleteByIndex(1) // After { 13, 15, 17} + { 18, 19, 20, 21} removed: {15}
	deltaSet.DeleteByIndex(4) // After { 13, 15, 17} + { 18, 19, 21}  removed { 15, 19}
	deltaSet.DeleteByIndex(5) // will remove { 13, 15, 17} + { 18, 19}

	if deltaSet.removed.Length() != 3 ||
		!reflect.DeepEqual(deltaSet.removed.Elements(), []int{15, 19, 20}) ||
		!reflect.DeepEqual(deltaSet.added.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.Commit(nil)

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 21}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.DeleteBatch([]int{13})
	deltaSet.Delete(17)
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 21}) || !reflect.DeepEqual(deltaSet.removed.Elements(), []int{13, 17}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if common.FilterFirst(deltaSet.Exists(13)) || common.FilterFirst(deltaSet.Exists(17)) || common.FilterFirst(deltaSet.Exists(25)) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if !common.FilterFirst(deltaSet.Exists(18)) || !common.FilterFirst(deltaSet.Exists(21)) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.InsertBatch([]int{13, 17, 22})                                        // Add the deleted entires back to the set
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 21}) || // Won't change until commit.
		// !reflect.DeepEqual(deltaSet.added.Elements(), []int{13, 17, 22}) ||
		!reflect.DeepEqual(deltaSet.removed.Elements(), []int{}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if !common.FilterFirst(deltaSet.Exists(13)) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if v, ok := deltaSet.TryGetKey(0); !ok || v != 13 {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if v, ok := deltaSet.TryGetKey(1); !ok || v != 17 {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if v, ok := deltaSet.TryGetKey(2); !ok || v != 18 {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if v, ok := deltaSet.TryGetKey(3); !ok || v != 21 {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}
}

func TestCascadeDeltaCommit(t *testing.T) {
	deltaSet := NewDeltaSet[int](-1, 100, nil)
	deltaSet.InsertBatch([]int{13, 15, 17})

	deltaSet.InsertBatch([]int{18, 19, 20, 21}) // { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.DeleteByIndex(1)                   // After { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.DeleteByIndex(4)                   // After { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.Commit(nil)

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 20, 21}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
}

func TestCascadeDeltaClone(t *testing.T) {
	deltaSet := NewDeltaSet[int](-1, 100, nil)

	deltaSet.InsertBatch([]int{13, 15, 17})
	deltaSet.Commit(nil)

	deltaSet.InsertBatch([]int{18, 19, 20, 21}) // { 13, 15, 17} + { 18, 19, 20, 21}

	// if deltaSet.NonNilCount() != 7 {
	// 	t.Error("failed to commit", deltaSet.NonNilCount())
	// }

	deltaSet.DeleteByIndex(1) //
	deltaSet.DeleteByIndex(4) //
	deltaSet.DeleteByIndex(5) // will remove {15, 19, 20}

	// if deltaSet.NonNilCount() != 4 {
	// 	t.Error("failed to commit", deltaSet.NonNilCount())
	// }

	if deltaSet.removed.Length() != 3 ||
		!reflect.DeepEqual(deltaSet.removed.Elements(), []int{15, 19, 20}) ||
		!reflect.DeepEqual(deltaSet.added.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	set2 := deltaSet.CloneFull()
	if !deltaSet.Equal(set2) {
		deltaSet.Print()
		set2.Print()
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.Commit(nil)
	set2.Commit(nil)

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 21}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if !deltaSet.Equal(set2) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.DeleteBatch([]int{13})
	if common.FilterFirst(deltaSet.Exists(13)) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if !common.FilterFirst(set2.Exists(13)) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if v, ok := deltaSet.TryGetKey(0); ok {
		t.Error("failed to commit", v) // Should not exist
	}

	if v, _ := deltaSet.TryGetKey(1); v != 17 {
		t.Error("failed to commit", v)
	}
}

func TestDeltaDeleteThenAddBack(t *testing.T) {

	deltaSet := NewDeltaSet(-1, 100, nil)
	deltaSet.InsertBatch([]int{13, 15, 17})
	deltaSet.Commit(nil)

	deltaSet.InsertBatch([]int{18, 19, 20, 21}) // { 13, 15, 17} + { 18, 19, 20, 21}

	deltaSet.DeleteByIndex(1) //
	deltaSet.DeleteByIndex(4) //
	deltaSet.DeleteByIndex(5) // will remove {15, 19, 20}
	if deltaSet.removed.Length() != 3 ||
		!reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17}) ||
		!reflect.DeepEqual(deltaSet.removed.Elements(), []int{15, 19, 20}) ||
		!reflect.DeepEqual(deltaSet.added.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.InsertBatch([]int{15, 19, 20}) // Add the deleted entires back to the set

	if deltaSet.removed.Length() != 0 ||
		// !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17}) ||
		// !reflect.DeepEqual(deltaSet.removed.Elements(), []int{}) ||
		!reflect.DeepEqual(deltaSet.added.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.added.Elements())
	}

	deltaSet.Commit(nil)

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17, 18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	v, ok := deltaSet.PopLast()
	if !ok || v != 21 || deltaSet.NonNilCount() != 6 {
		t.Error("failed to commit", v, deltaSet.NonNilCount())
	}

	v, ok = deltaSet.PopLast()
	if !ok || v != 20 || deltaSet.NonNilCount() != 5 {
		t.Error("failed to commit", v, deltaSet.NonNilCount())
	}

	v, ok = deltaSet.GetByIndex(5)
	if ok || v != 0 || deltaSet.NonNilCount() != 5 {
		t.Error("Should not exist", v, deltaSet.NonNilCount())
	}
}

func TestMultiMerge(t *testing.T) {
	deltaSet := NewDeltaSet(-1, 100, nil)
	deltaSet.InsertBatch([]int{13, 15, 17})
	deltaSet.Commit(nil)

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	_set0 := NewDeltaSet(-1, 100, nil).InsertBatch([]int{58, 59, 20, 51}).DeleteBatch([]int{13})
	_set1 := NewDeltaSet(-1, 100, nil).InsertBatch([]int{78, 59, 70, 71}).DeleteBatch([]int{15, 70})

	// (13, 15, 17) + (58, 59, 20, 51) + (78, 59, 70, 71) - (13, 15, 70) = (17, 58, 59, 20, 51, 78, 59, 71)
	deltaSet.Commit([]*DeltaSet[int]{_set0, _set1})

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17, 58, 59, 20, 51, 78, 71}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
}

func TestGetNthNonNil(t *testing.T) {
	deltaSet := NewDeltaSet(0, 100, nil)
	deltaSet.InsertBatch([]int{13, 15, 17})

	deltaSet.InsertBatch([]int{18, 19, 20, 21}) //  { 13, 15, 17} +  { 18, 19, 20, 21}
	// DeleteByIndex wouldn't shift the indices
	deltaSet.DeleteByIndex(1) //  { 13, 15, 17} + { 18, 19, 20, 21} - {15} = { 13, 17, 18, 19, 20, 21}
	deltaSet.DeleteByIndex(4) //  { 13, 15, 17} + { 18, 19, 20, 21} - {15, 19}

	if k, idx, ok := deltaSet.GetNthNonNil(0); k != 13 || idx != 0 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(1); k != 17 || idx != 2 || !ok {
		t.Error("failed to commit", k)
	}

	// Check if the deleted entry is still accessible.
	if k, ok := deltaSet.GetByIndex(1); k != deltaSet.nilVal || ok {
		t.Error("A deleted entry shouldn't be available any more", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(2); k != 18 || idx != 3 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(3); k != 20 || idx != 5 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(4); k != 21 || idx != 6 || !ok {
		t.Error("failed to commit", k)
	}

	// Check if the deleted entry is still accessible.
	if k, ok := deltaSet.GetByIndex(4); k != deltaSet.nilVal || ok {
		t.Error("A deleted entry shouldn't be available any more", k)
	}

	deltaSet.Commit(nil)

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 20, 21}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if k, idx, ok := deltaSet.GetNthNonNil(0); k != 13 || idx != 0 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(1); k != 17 || idx != 1 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(2); k != 18 || idx != 2 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(3); k != 20 || idx != 3 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(4); k != 21 || idx != 4 || !ok {
		t.Error("failed to commit", k)
	}

}

func BenchmarkDeltaDeleteThenAddBack(t *testing.B) {
	deltaSet := NewDeltaSet[int](-1, 1000000, nil)
	randoms := make([]int, 1000000)
	for i := 0; i < 1000000; i++ {
		randoms[i] = i //rand.Int()
	}

	t0 := time.Now()
	deltaSet.InsertBatch(randoms)
	fmt.Println("InsertBatch", time.Since(t0))

	t0 = time.Now()

	deltaSet.Commit(nil)

	fmt.Println("Commit", time.Since(t0))

	t0 = time.Now()
	deltaSet.CloneDelta()
	fmt.Println("CloneDelta", time.Since(t0))
}

func BenchmarkGetNthNonNil(b *testing.B) {
	deltaSet := NewDeltaSet[int](-1, 100, nil)
	deltaSet.InsertBatch([]int{13, 15, 17})

	deltaSet.InsertBatch([]int{18, 19, 20, 21}) // { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.DeleteByIndex(1)                   //  { 13, -15, 17} + { 18, 19, 20, 21}
	deltaSet.DeleteByIndex(4)                   // { 13, -15, 17} + { 18, -19, 20, 21}'

	total := 10000
	for i := 0; i < total; i++ {
		deltaSet.InsertBatch([]int{i})
	}

	for i := 0; i < total/2; i++ {
		deltaSet.DeleteByIndex(uint64(rand.Intn(1000000)))
	}

	t0 := time.Now()
	for i := 0; i < total; i++ {
		deltaSet.GetNthNonNil(uint64(i))
	}
	fmt.Println("GetNthNonNil", time.Since(t0))

	// t0 = time.Now()
	// for i := 0; i < total; i++ {
	// 	deltaSet.GetNthNonNil(uint64(i))
	// }
	// fmt.Println("GetNthNonNilv2", time.Since(t0))
}
