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
	"reflect"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/common"
)

func TestDeltaCommit(t *testing.T) {
	cascadeSet := NewCascadeDeltaSet(-1, 100)
	cascadeSet.Insert(13, 15, 17)

	cascadeSet.Insert(18, 19, 20, 21) // { 13, 15, 17} + { 18, 19, 20, 21}
	cascadeSet.DeleteByIndex(1)       // After { 13, 15, 17} + { 18, 19, 20, 21}
	cascadeSet.DeleteByIndex(4)       // After { 13, 15, 17} + { 18, 19, 20, 21}
	cascadeSet.Commit()

	if !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 17, 18, 20, 21}) {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}
}

func TestCascadeDeltaSliceBasic(t *testing.T) {
	cascadeSet := NewCascadeDeltaSet(-1, 100)

	if cascadeSet.Insert(11, 12, 13); !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{}) ||
		!cascadeSet.updated.First.IsDirty() ||
		!reflect.DeepEqual(cascadeSet.updated.First.Elements(), []int{11, 12, 13}) {
		t.Error("failed to append", cascadeSet.committed.Elements(), cascadeSet.updated.First.Elements())
	}

	if cascadeSet.Commit(); cascadeSet.committed.Length() != 3 ||
		!cascadeSet.updated.First.IsDirty() ||
		(cascadeSet.updated.First.Length()) != 0 ||
		(cascadeSet.removed.Length()) != 0 {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}

	if cascadeSet.Delete(12); cascadeSet.committed.Length() != 3 || !cascadeSet.updated.First.IsDirty() ||
		(cascadeSet.updated.First.Length()) != 0 || (cascadeSet.removed.Length()) != 1 {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}

	if cascadeSet.Commit(); cascadeSet.committed.Length() != 2 || (cascadeSet.updated.First.Length()) != 0 || !cascadeSet.updated.First.IsDirty() ||
		(cascadeSet.removed.Length()) != 0 { // {11, 13}
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}
	// {11, 13} + {15, 16, 17}
	if cascadeSet.Insert(15, 16, 17); cascadeSet.committed.Length() != 2 || (cascadeSet.updated.First.Length()) != 3 || (cascadeSet.removed.Length()) != 0 {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}
	// {11, 13} + {15, 16, 17}
	if cascadeSet.Delete(16); cascadeSet.committed.Length() != 2 ||
		(cascadeSet.updated.First.Length()) != 3 ||
		// !cascadeSet.updated.IsDirty() ||
		!reflect.DeepEqual(cascadeSet.updated.First.Elements(), []int{15, 16, 17}) ||
		(cascadeSet.removed.Length()) != 1 {
		t.Error("failed to commit", cascadeSet.updated.First.Elements(), cascadeSet.committed.Elements())
	}

	if cascadeSet.Delete(11); cascadeSet.committed.Length() != 2 ||
		(cascadeSet.updated.First.Length()) != 3 ||
		!reflect.DeepEqual(cascadeSet.updated.First.Elements(), []int{15, 16, 17}) ||
		(cascadeSet.removed.Length()) != 2 {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}

	if cascadeSet.Delete(16); cascadeSet.committed.Length() != 2 ||
		!cascadeSet.updated.First.IsDirty() ||
		(cascadeSet.updated.First.Length()) != 3 ||
		!reflect.DeepEqual(cascadeSet.updated.First.Elements(), []int{15, 16, 17}) ||
		(cascadeSet.removed.Length()) != 2 {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}

	if !reflect.DeepEqual(cascadeSet.Elements(), []int{13, 15, 17}) {
		t.Error("Failed to get Elements()", cascadeSet.Elements())
	}

	if cascadeSet.Commit(); !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 15, 17}) || (cascadeSet.updated.First.Length()) != 0 || (cascadeSet.removed.Length()) != 0 {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}
	// { 13, 15, 17} + { 18, 19, 20, 21}
	cascadeSet.Insert(18, 19, 20, 21)

	if k, _, _, ok := cascadeSet.Search(0); !ok || k != 13 {
		t.Error("failed to commit", k)
	}

	if k, _, _, ok := cascadeSet.Search(1); !ok || k != 15 {
		t.Error("failed to commit", k)
	}

	if k, _, _, ok := cascadeSet.Search(2); !ok || k != 17 {
		t.Error("failed to commit", k)
	}

	if k, _, _, ok := cascadeSet.Search(3); !ok || k != 18 {
		t.Error("failed to commit", k)
	}
	if k, _, _, ok := cascadeSet.Search(4); !ok || k != 19 {
		t.Error("failed to commit", k)
	}
	if k, _, _, ok := cascadeSet.Search(5); !ok || k != 20 {
		t.Error("failed to commit", k)
	}

	cascadeSet.DeleteByIndex(1) // After { 13, 15, 17} + { 18, 19, 20, 21}
	cascadeSet.DeleteByIndex(4) // After { 13, 15, 17} + { 18, 19, 21}
	cascadeSet.DeleteByIndex(5) // will remove { 13, 15, 17} + { 18, 19}

	if !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 15, 17}) || !reflect.DeepEqual(cascadeSet.updated.First.Elements(), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", cascadeSet.committed.Elements(), cascadeSet.updated.First.Elements())
	}
}

func TestCascadeDeltaSliceAddThenDelete(t *testing.T) {
	cascadeSet := NewCascadeDeltaSet[int](-1, 100)
	cascadeSet.Insert(13, 15, 17)
	cascadeSet.Commit()

	cascadeSet.Insert(18, 19, 20, 21) // { 13, 15, 17} + { 18, 19, 20, 21}

	cascadeSet.DeleteByIndex(1) // After { 13, 15, 17} + { 18, 19, 20, 21}
	cascadeSet.DeleteByIndex(4) // After { 13, 15, 17} + { 18, 19, 21}
	cascadeSet.DeleteByIndex(5) // will remove { 13, 15, 17} + { 18, 19}

	if cascadeSet.removed.Length() != 3 ||
		!reflect.DeepEqual(cascadeSet.removed.Elements(), []int{15, 19, 20}) ||
		!reflect.DeepEqual(append(cascadeSet.updated.Second.Elements(), cascadeSet.updated.First.Elements()...), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	cascadeSet.Commit()
	if !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 17, 18, 21}) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	cascadeSet.Delete(13)
	cascadeSet.Delete(17)
	if !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 17, 18, 21}) || !reflect.DeepEqual(cascadeSet.removed.Elements(), []int{13, 17}) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	if common.FilterFirst(cascadeSet.Exists(13)) || common.FilterFirst(cascadeSet.Exists(17)) || common.FilterFirst(cascadeSet.Exists(25)) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	if !common.FilterFirst(cascadeSet.Exists(18)) || !common.FilterFirst(cascadeSet.Exists(21)) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	cascadeSet.Insert(13, 17, 22) // Add they deleted entires back to the set
	if !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 17, 18, 21}) ||
		!reflect.DeepEqual(cascadeSet.removed.Elements(), []int{}) ||
		!reflect.DeepEqual(cascadeSet.removed.Elements(), []int{}) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	if !common.FilterFirst(cascadeSet.Exists(13)) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	if v, ok := cascadeSet.TryGetKey(0); !ok || v != 13 {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	if v, ok := cascadeSet.TryGetKey(1); !ok || v != 17 {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	if v, ok := cascadeSet.TryGetKey(2); !ok || v != 18 {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	if v, ok := cascadeSet.TryGetKey(3); !ok || v != 21 {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}
}

func TestDeltaClone(t *testing.T) {
	cascadeSet := NewCascadeDeltaSet[int](-1, 100)
	cascadeSet.Insert(13, 15, 17)
	cascadeSet.Commit()

	cascadeSet.Insert(18, 19, 20, 21) // { 13, 15, 17} + { 18, 19, 20, 21}

	// if cascadeSet.NonNilCount() != 7 {
	// 	t.Error("failed to commit", cascadeSet.NonNilCount())
	// }

	cascadeSet.DeleteByIndex(1) //
	cascadeSet.DeleteByIndex(4) //
	cascadeSet.DeleteByIndex(5) // will remove {15, 19, 20}

	// if cascadeSet.NonNilCount() != 4 {
	// 	t.Error("failed to commit", cascadeSet.NonNilCount())
	// }

	if cascadeSet.removed.Length() != 3 ||
		!reflect.DeepEqual(cascadeSet.removed.Elements(), []int{15, 19, 20}) ||
		!reflect.DeepEqual(append(cascadeSet.updated.Second.Elements(), cascadeSet.updated.First.Elements()...), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	set2 := cascadeSet.CloneFull()
	if !cascadeSet.Equal(set2) {
		cascadeSet.Print()
		set2.Print()
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	cascadeSet.Commit()
	set2.Commit()

	if !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 17, 18, 21}) {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}

	if !cascadeSet.Equal(set2) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	cascadeSet.Delete(13)
	if common.FilterFirst(cascadeSet.Exists(13)) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	if !common.FilterFirst(set2.Exists(13)) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	if v, ok := cascadeSet.TryGetKey(0); ok {
		t.Error("failed to commit", v) // Should not exist
	}

	if v, _ := cascadeSet.TryGetKey(1); v != 17 {
		t.Error("failed to commit", v)
	}
}

func TestCascadeDeltaDeleteThenAddBack(t *testing.T) {
	cascadeSet := NewCascadeDeltaSet[int](-1, 100)
	cascadeSet.Insert(13, 15, 17)
	cascadeSet.Commit()

	cascadeSet.Insert(18, 19, 20, 21) // { 13, 15, 17} + { 18, 19, 20, 21}

	cascadeSet.DeleteByIndex(1) //
	cascadeSet.DeleteByIndex(4) //
	cascadeSet.DeleteByIndex(5) // will remove {15, 19, 20}
	if cascadeSet.removed.Length() != 3 ||
		!reflect.DeepEqual(cascadeSet.removed.Elements(), []int{15, 19, 20}) ||
		!reflect.DeepEqual(append(cascadeSet.updated.Second.Elements(), cascadeSet.updated.First.Elements()...), []int{18, 19, 20, 21}) {
		t.Error("failed to commit", cascadeSet.removed.Elements())
	}

	cascadeSet.Insert(15, 19, 20) // Add the deleted entires back to the set

	if cascadeSet.removed.NonNilCount() != 0 ||
		!reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 15, 17}) ||
		!reflect.DeepEqual(cascadeSet.removed.Elements(), []int{}) ||
		!reflect.DeepEqual(append(cascadeSet.updated.Second.Elements(), cascadeSet.updated.First.Elements()...), []int{18, 19, 20, 21, 15}) {
		t.Error("failed to commit",
			cascadeSet.removed.Length(),
			cascadeSet.committed.Elements(),
			cascadeSet.removed.Elements(),
			append(cascadeSet.updated.Second.Elements(), cascadeSet.updated.First.Elements()...),
		)
	}

	cascadeSet.Commit()
	if !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 15, 17, 18, 19, 20, 21}) {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}

	v, ok := cascadeSet.PopLast()
	if !ok || v != 21 || cascadeSet.NonNilCount() != 6 {
		t.Error("failed to commit", v, cascadeSet.NonNilCount())
	}

	v, ok = cascadeSet.PopLast()
	if !ok || v != 20 || cascadeSet.NonNilCount() != 5 {
		t.Error("failed to commit", v, cascadeSet.NonNilCount())
	}

	v, ok = cascadeSet.GetByIndex(5)
	if ok || v != 0 || cascadeSet.NonNilCount() != 5 {
		t.Error("Should not exist", v, cascadeSet.NonNilCount())
	}
}

func TestCascadeGetNthNonNil(t *testing.T) {
	cascadeSet := NewCascadeDeltaSet(-1, 100)
	cascadeSet.Insert(13, 15, 17)

	cascadeSet.Insert(18, 19, 20, 21) // { 13, 15, 17} + { 18, 19, 20, 21}
	cascadeSet.DeleteByIndex(1)       //  { 13, -15, 17} + { 18, 19, 20, 21}
	cascadeSet.DeleteByIndex(4)       // { 13, -15, 17} + { 18, -19, 20, 21}'

	if k, idx, ok := cascadeSet.GetNthNonNil(0); k != 13 || idx != 0 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := cascadeSet.GetNthNonNil(1); k != 17 || idx != 2 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := cascadeSet.GetNthNonNil(2); k != 18 || idx != 3 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := cascadeSet.GetNthNonNil(3); k != 20 || idx != 5 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := cascadeSet.GetNthNonNil(4); k != 21 || idx != 6 || !ok {
		t.Error("failed to commit", k)
	}

	cascadeSet.Commit()
	if !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{13, 17, 18, 20, 21}) {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}

	if k, idx, ok := cascadeSet.GetNthNonNil(0); k != 13 || idx != 0 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := cascadeSet.GetNthNonNil(1); k != 17 || idx != 1 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := cascadeSet.GetNthNonNil(2); k != 18 || idx != 2 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := cascadeSet.GetNthNonNil(3); k != 20 || idx != 3 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := cascadeSet.GetNthNonNil(4); k != 21 || idx != 4 || !ok {
		t.Error("failed to commit", k)
	}

}

func TestCascadeMultiMerge(t *testing.T) {
	cascadeSet := NewCascadeDeltaSet[int](-1, 100)
	cascadeSet.Insert(13, 15, 17)
	cascadeSet.Commit()

	_set0 := NewCascadeDeltaSet[int](-1, 100)
	_set1 := NewCascadeDeltaSet[int](-1, 100)

	_set0.Insert(58, 59, 20, 51) // { 13, 15, 17} + { 18, 19, 20, 21}
	_set1.Insert(78, 59, 70, 71) // { 13, 15, 17} + { 18, 19, 20, 21}

	_set0.Delete(13)
	_set1.Delete(15, 70)

	// (13, 15, 17) + (58, 59, 20, 51) + (78, 59, 70, 71) - (13, 15, 70) = (17, 58, 59, 20, 51, 78, 59, 71)
	cascadeSet.Commit(_set0, _set1)

	if !reflect.DeepEqual(cascadeSet.committed.Elements(), []int{17, 58, 59, 20, 51, 78, 71}) {
		t.Error("failed to commit", cascadeSet.committed.Elements())
	}
}

func BenchmarkCascadeDeltaDeleteThenAddBack(t *testing.B) {
	cascadeSet := NewDeltaSet[int](-1, 1000000)
	randoms := make([]int, 1000000)
	for i := 0; i < 1000000; i++ {
		randoms[i] = i
	}

	t0 := time.Now()
	cascadeSet.Insert(randoms...)
	fmt.Println("Insert", time.Since(t0))

	t0 = time.Now()
	cascadeSet.Commit()
	fmt.Println("Commit", time.Since(t0))

	t0 = time.Now()
	cascadeSet.CloneDelta()
	fmt.Println("CloneDelta", time.Since(t0))
}
