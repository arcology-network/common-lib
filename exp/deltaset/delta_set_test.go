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
	"math/rand/v2"
	"reflect"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/common"
)

func TestDeltaSliceBasic(t *testing.T) {
	deltaSet := NewDeltaSet[int](-1, 100)

	if deltaSet.Insert(11, 12, 13); !reflect.DeepEqual(deltaSet.committed.Elements(), []int{}) ||
		!deltaSet.updated.IsSynced() ||
		!reflect.DeepEqual(deltaSet.updated.Elements(), []int{11, 12, 13}) {
		t.Error("failed to append", deltaSet.committed.Elements(), deltaSet.updated.Elements())
	}

	if deltaSet.Commit(); deltaSet.committed.Length() != 3 ||
		!deltaSet.updated.IsSynced() ||
		(deltaSet.updated.Length()) != 0 ||
		(deltaSet.removed.Length()) != 0 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if deltaSet.Delete(12); deltaSet.committed.Length() != 3 || !deltaSet.updated.IsSynced() ||
		(deltaSet.updated.Length()) != 0 || (deltaSet.removed.Length()) != 1 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if deltaSet.Commit(); deltaSet.committed.Length() != 2 || (deltaSet.updated.Length()) != 0 || !deltaSet.updated.IsSynced() ||
		(deltaSet.removed.Length()) != 0 { // {11, 13}
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// {11, 13} + {15, 16, 17}
	if deltaSet.Insert(15, 16, 17); deltaSet.committed.Length() != 2 || (deltaSet.updated.Length()) != 3 || (deltaSet.removed.Length()) != 0 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// {11, 13} + {15, 16, 17}
	if deltaSet.Delete(16); deltaSet.committed.Length() != 2 ||
		(deltaSet.updated.Length()) != 3 ||
		// !deltaSet.updated.IsSynced() ||
		!reflect.DeepEqual(deltaSet.updated.Elements(), []int{15, 16, 17}) ||
		(deltaSet.removed.Length()) != 1 {
		t.Error("failed to commit", deltaSet.updated.Elements(), deltaSet.committed.Elements())
	}

	if deltaSet.Delete(11); deltaSet.committed.Length() != 2 ||
		(deltaSet.updated.Length()) != 3 ||
		!reflect.DeepEqual(deltaSet.updated.Elements(), []int{15, 16, 17}) ||
		(deltaSet.removed.Length()) != 2 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if deltaSet.Delete(16); deltaSet.committed.Length() != 2 ||
		!deltaSet.updated.IsSynced() ||
		(deltaSet.updated.Length()) != 3 ||
		!reflect.DeepEqual(deltaSet.updated.Elements(), []int{15, 16, 17}) ||
		(deltaSet.removed.Length()) != 2 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if deltaSet.Commit(); !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17}) || (deltaSet.updated.Length()) != 0 || (deltaSet.removed.Length()) != 0 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.Insert(18, 19, 20, 21)

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

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17}) || !reflect.DeepEqual(deltaSet.updated.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.committed.Elements(), deltaSet.updated.Elements())
	}
}

func TestDeltaSliceAddThenDelete(t *testing.T) {
	deltaSet := NewDeltaSet[int](-1, 100)
	deltaSet.Insert(13, 15, 17)
	deltaSet.Commit()

	deltaSet.Insert(18, 19, 20, 21) // { 13, 15, 17} + { 18, 19, 20, 21}

	deltaSet.DeleteByIndex(1) // After { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.DeleteByIndex(4) // After { 13, 15, 17} + { 18, 19, 21}
	deltaSet.DeleteByIndex(5) // will remove { 13, 15, 17} + { 18, 19}

	if deltaSet.removed.Length() != 3 ||
		!reflect.DeepEqual(deltaSet.removed.Elements(), []int{15, 19, 20}) ||
		!reflect.DeepEqual(deltaSet.updated.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.Commit()
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 21}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.Delete(13)
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

	deltaSet.Insert(13, 17, 22) // Add they deleted entires back to the set
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 21}) ||
		!reflect.DeepEqual(deltaSet.removed.Elements(), []int{}) ||
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

func TestDeltaCommit(t *testing.T) {
	deltaSet := NewDeltaSet[int](-1, 100)
	deltaSet.Insert(13, 15, 17)

	deltaSet.Insert(18, 19, 20, 21) // { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.DeleteByIndex(1)       // After { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.DeleteByIndex(4)       // After { 13, 15, 17} + { 18, 19, 20, 21}
	deltaSet.Commit()

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 20, 21}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
}

func TestDeltaClone(t *testing.T) {
	deltaSet := NewDeltaSet[int](-1, 100)
	deltaSet.Insert(13, 15, 17)
	deltaSet.Commit()

	deltaSet.Insert(18, 19, 20, 21) // { 13, 15, 17} + { 18, 19, 20, 21}

	// if deltaSet.Length() != 7 {
	// 	t.Error("failed to commit", deltaSet.Length())
	// }

	deltaSet.DeleteByIndex(1) //
	deltaSet.DeleteByIndex(4) //
	deltaSet.DeleteByIndex(5) // will remove {15, 19, 20}

	// if deltaSet.Length() != 4 {
	// 	t.Error("failed to commit", deltaSet.Length())
	// }

	if deltaSet.removed.Length() != 3 ||
		!reflect.DeepEqual(deltaSet.removed.Elements(), []int{15, 19, 20}) ||
		!reflect.DeepEqual(deltaSet.updated.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	set2 := deltaSet.Clone()
	if !deltaSet.Equal(set2) {
		deltaSet.Print()
		set2.Print()
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.Commit()
	set2.Commit()

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 17, 18, 21}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	if !deltaSet.Equal(set2) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.Delete(13)
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
	deltaSet := NewDeltaSet[int](-1, 100)
	deltaSet.Insert(13, 15, 17)
	deltaSet.Commit()

	deltaSet.Insert(18, 19, 20, 21) // { 13, 15, 17} + { 18, 19, 20, 21}

	deltaSet.DeleteByIndex(1) //
	deltaSet.DeleteByIndex(4) //
	deltaSet.DeleteByIndex(5) // will remove {15, 19, 20}
	if deltaSet.removed.Length() != 3 ||
		!reflect.DeepEqual(deltaSet.removed.Elements(), []int{15, 19, 20}) ||
		!reflect.DeepEqual(deltaSet.updated.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.removed.Elements())
	}

	deltaSet.Insert(15, 19, 20) // Add the deleted entires back to the set

	if deltaSet.removed.Length() != 0 ||
		!reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17}) ||
		!reflect.DeepEqual(deltaSet.removed.Elements(), []int{}) ||
		!reflect.DeepEqual(deltaSet.updated.Elements(), []int{18, 19, 20, 21, 15}) {
		t.Error("failed to commit", deltaSet.updated.Elements())
	}

	deltaSet.Commit()
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{13, 15, 17, 18, 19, 20, 21}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	v, ok := deltaSet.PopBack()
	if !ok || v != 21 || deltaSet.Length() != 6 {
		t.Error("failed to commit", v, deltaSet.Length())
	}
}

func TestMultiMerge(t *testing.T) {
	deltaSet := NewDeltaSet[int](-1, 100)
	deltaSet.Insert(13, 15, 17)
	deltaSet.Commit()

	_set0 := NewDeltaSet[int](-1, 100)
	_set1 := NewDeltaSet[int](-1, 100)

	_set0.Insert(58, 59, 20, 51) // { 13, 15, 17} + { 18, 19, 20, 21}
	_set1.Insert(78, 59, 70, 71) // { 13, 15, 17} + { 18, 19, 20, 21}

	_set0.Delete(13)
	_set1.Delete(15, 70)

	// (13, 15, 17) + (58, 59, 20, 51) + (78, 59, 70, 71) - (13, 15, 70) = (17, 58, 59, 20, 51, 78, 59, 71)
	deltaSet.Commit(_set0, _set1)

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []int{17, 58, 59, 20, 51, 78, 71}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
}

func BenchmarkDeltaDeleteThenAddBack(t *testing.B) {
	deltaSet := NewDeltaSet[int](-1, 1000000)
	randoms := make([]int, 1000000)
	for i := 0; i < 1000000; i++ {
		randoms[i] = rand.Int()
	}

	t0 := time.Now()
	deltaSet.Insert(randoms...)
	fmt.Println("Insert", time.Since(t0))

	t0 = time.Now()
	deltaSet.Commit()
	fmt.Println("Commit", time.Since(t0))

	t0 = time.Now()
	deltaSet.CloneDelta()
	fmt.Println("CloneDelta", time.Since(t0))
}
