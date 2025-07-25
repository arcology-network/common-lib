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
package stringdeltaset

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
)

func TestSoftDeltaSliceBasic(t *testing.T) {
	deltaSet := NewSoftDeltaSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)

	if deltaSet.InsertBatch([]string{"11", "12", "13"}); !reflect.DeepEqual(deltaSet.committed.Elements(), []string{}) ||
		!reflect.DeepEqual(deltaSet.stagedAdditions.Elements(), []string{"11", "12", "13"}) {
		t.Error("failed to append", deltaSet.stagedAdditions.IsDirty(), deltaSet.committed.Elements(), deltaSet.stagedAdditions.Elements())
	}

	if deltaSet.Commit(nil); deltaSet.committed.Length() != 3 ||
		(deltaSet.stagedAdditions.Length()) != 0 ||
		(deltaSet.stagedRemovals.Length()) != 0 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if deltaSet.DeleteBatch([]string{"12"}); deltaSet.committed.Length() != 3 ||
		(deltaSet.stagedAdditions.Length()) != 0 || (deltaSet.stagedRemovals.Length()) != 1 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if deltaSet.Commit(nil); deltaSet.committed.Length() != 2 || (deltaSet.stagedAdditions.Length()) != 0 || deltaSet.stagedAdditions.IsDirty() ||
		(deltaSet.stagedRemovals.Length()) != 0 { // {"11", "13"}
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// {"11", "13"} + {"15", "16", "17"}
	if deltaSet.InsertBatch([]string{"15", "16", "17"}); deltaSet.committed.Length() != 2 || (deltaSet.stagedAdditions.Length()) != 3 || (deltaSet.stagedRemovals.Length()) != 0 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// {"11", "13"} + {"15", "16", "17"}
	if deltaSet.Delete("16"); deltaSet.committed.Length() != 2 ||
		(deltaSet.stagedAdditions.Length()) != 3 ||
		// !deltaSet.stagedAdditions .IsDirty() ||
		!reflect.DeepEqual(deltaSet.stagedAdditions.Elements(), []string{"15", "16", "17"}) ||
		(deltaSet.stagedRemovals.Length()) != 1 {
		t.Error("failed to commit", deltaSet.stagedAdditions.Elements(), deltaSet.committed.Elements())
	}

	if deltaSet.DeleteBatch([]string{"11"}); deltaSet.committed.Length() != 2 ||
		(deltaSet.stagedAdditions.Length()) != 3 ||
		!reflect.DeepEqual(deltaSet.stagedAdditions.Elements(), []string{"15", "16", "17"}) ||
		(deltaSet.stagedRemovals.Length()) != 2 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// Re-delete a uncommitted entry, the stagedRemovals set will grow but neither the committed nor stagedAdditions  set not change.
	if deltaSet.Delete("16"); deltaSet.committed.Length() != 2 ||
		(deltaSet.stagedAdditions.Length()) != 3 ||
		!reflect.DeepEqual(deltaSet.stagedAdditions.Elements(), []string{"15", "16", "17"}) ||
		(deltaSet.stagedRemovals.Length()) != 2 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if !reflect.DeepEqual(deltaSet.Elements(), []string{"13", "15", "17"}) {
		t.Error("Failed to get Elements()", deltaSet.Elements())
	}

	if deltaSet.Commit(nil); !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "15", "17"}) || (deltaSet.stagedAdditions.Length()) != 0 || (deltaSet.stagedRemovals.Length()) != 0 {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
	// { "13", "15", "17"} + { "18", "19", "20", "21"}
	deltaSet.InsertBatch([]string{"18", "19", "20", "21"})

	if k, _, _, ok := deltaSet.Search(0); !ok || *k != "13" {
		t.Error("failed to commit", k)
	}

	if k, _, _, ok := deltaSet.Search(1); !ok || *k != "15" {
		t.Error("failed to commit", k)
	}

	if k, _, _, ok := deltaSet.Search(2); !ok || *k != "17" {
		t.Error("failed to commit", k)
	}

	if k, _, _, ok := deltaSet.Search(3); !ok || *k != "18" {
		t.Error("failed to commit", k)
	}
	if k, _, _, ok := deltaSet.Search(4); !ok || *k != "19" {
		t.Error("failed to commit", k)
	}
	if k, _, _, ok := deltaSet.Search(5); !ok || *k != "20" {
		t.Error("failed to commit", k)
	}

	deltaSet.DeleteByIndex(1) // After { "13", "15", "17"} + { "18", "19", "20", "21"}
	deltaSet.DeleteByIndex(4) // After { "13", "15", "17"} + { "18", "19", "21"}
	deltaSet.DeleteByIndex(5) // will remove { "13", "15", "17"} + { "18", "19"}

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "15", "17"}) || !reflect.DeepEqual(deltaSet.stagedAdditions.Elements(), []string{"18", "19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.committed.Elements(), deltaSet.stagedAdditions.Elements())
	}
}

func TestSoftDeltaSliceAddThenDelete(t *testing.T) {
	deltaSet := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})
	deltaSet.Commit(nil)

	deltaSet.InsertBatch([]string{"18", "19", "20", "21"}) // { "13", "15", "17"} + { "18", "19", "20", "21"}

	// deltaSet.Delete([]string{118, 119, 210, 211}) // non-existing entries, should not affect the set
	// if deltaSet.stagedRemovals.Length() != 0 {
	// 	t.Error("Deleting non-existing elements should not affect the set", deltaSet.stagedRemovals.Elements())
	// }

	deltaSet.DeleteByIndex(1) // After { "13", "15", "17"} + { "18", "19", "20", "21"} stagedRemovals: {"15"}
	deltaSet.DeleteByIndex(4) // After { "13", "15", "17"} + { "18", "19", "21"}  stagedRemovals { "15", "19"}
	deltaSet.DeleteByIndex(5) // will remove { "13", "15", "17"} + { "18", "19"}

	if deltaSet.stagedRemovals.Length() != 3 ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Elements(), []string{"15", "19", "20"}) ||
		!reflect.DeepEqual(deltaSet.stagedAdditions.Elements(), []string{"18", "19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	deltaSet.Commit(nil)
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "17", "18", "21"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	deltaSet.DeleteBatch([]string{"13"})
	deltaSet.Delete("17")
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "17", "18", "21"}) || !reflect.DeepEqual(deltaSet.stagedRemovals.Elements(), []string{"13", "17"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if common.FilterFirst(deltaSet.Exists("13")).(bool) || common.FilterFirst(deltaSet.Exists("17")).(bool) || common.FilterFirst(deltaSet.Exists("25")).(bool) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if !common.FilterFirst(deltaSet.Exists("18")).(bool) || !common.FilterFirst(deltaSet.Exists("21")).(bool) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	deltaSet.InsertBatch([]string{"13", "17", "22"})                                          // Add the deleted entires back to the set
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "17", "18", "21"}) || // Won't change until commit.
		// !reflect.DeepEqual(deltaSet.stagedAdditions .Elements(), []string{"13", "17", "22"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Elements(), []string{}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if !common.FilterFirst(deltaSet.Exists("13")).(bool) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if v, ok := deltaSet.TryGetKey(0); !ok || *v != "13" {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if v, ok := deltaSet.TryGetKey(1); !ok || *v != "17" {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if v, ok := deltaSet.TryGetKey(2); !ok || *v != "18" {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if v, ok := deltaSet.TryGetKey(3); !ok || *v != "21" {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}
}

func TestCascadeSoftDeltaCommit(t *testing.T) {
	deltaSet := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})

	deltaSet.InsertBatch([]string{"18", "19", "20", "21"}) // { "13", "15", "17"} + { "18", "19", "20", "21"}
	deltaSet.DeleteByIndex(1)                              // After { "13", "15", "17"} + { "18", "19", "20", "21"}
	deltaSet.DeleteByIndex(4)                              // After { "13", "15", "17"} + { "18", "19", "20", "21"}
	deltaSet.Commit(nil)

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "17", "18", "20", "21"}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
}

func TestCascadeSoftDeltaClone(t *testing.T) {
	deltaSet := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})
	deltaSet.Commit(nil)

	deltaSet.InsertBatch([]string{"18", "19", "20", "21"}) // { "13", "15", "17"} + { "18", "19", "20", "21"}

	deltaSet.DeleteByIndex(1) //
	deltaSet.DeleteByIndex(4) //
	deltaSet.DeleteByIndex(5) // will remove {"15", "19", "20"} left: {"13", "17", "18", "21"}

	if deltaSet.stagedRemovals.Length() != 3 ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Elements(), []string{"15", "19", "20"}) ||
		!reflect.DeepEqual(deltaSet.stagedAdditions.Elements(), []string{"18", "19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if !reflect.DeepEqual(deltaSet.Elements(), []string{"13", "17", "18", "21"}) {
		t.Error("failed to commit", deltaSet.Elements())
	}

	set2 := deltaSet.CloneFull()
	if !deltaSet.Equal(set2) {
		deltaSet.Print()
		set2.Print()
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if !deltaSet.Equal(set2) {
		t.Error("Mismatch", "expected", set2, "actual:", deltaSet.stagedRemovals.Elements())
		deltaSet.Print()
		fmt.Println("-----------------------")
		set2.Print()
	}
	set2.Commit(nil)
	deltaSet.Commit(nil)

	if !deltaSet.Equal(set2) {
		t.Error("Mismatch", "expected", set2, "actual:", deltaSet.stagedRemovals.Elements())
	}

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "17", "18", "21"}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	deltaSet.DeleteBatch([]string{"13"})
	if common.FilterFirst(deltaSet.Exists("13")).(bool) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if !common.FilterFirst(set2.Exists("13")).(bool) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	if v, ok := deltaSet.TryGetKey(0); ok {
		t.Error("failed to commit", v) // Should not exist
	}

	if v, _ := deltaSet.TryGetKey(1); *v != "17" {
		t.Error("failed to commit", v)
	}
}

func TestSoftDeltaDeleteThenAddBack(t *testing.T) {
	deltaSet := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})
	deltaSet.Commit(nil)

	deltaSet.InsertBatch([]string{"18", "19", "20", "21"}) // { "13", "15", "17"} + { "18", "19", "20", "21"}

	deltaSet.DeleteByIndex(1) //
	deltaSet.DeleteByIndex(4) //
	deltaSet.DeleteByIndex(5) // will remove {"15", "19", "20"}
	if deltaSet.stagedRemovals.Length() != 3 ||
		!reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "15", "17"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Elements(), []string{"15", "19", "20"}) ||
		!reflect.DeepEqual(deltaSet.stagedAdditions.Elements(), []string{"18", "19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	deltaSet.InsertBatch([]string{"15", "19", "20"}) // Add the deleted entires back to the set

	if deltaSet.stagedRemovals.Length() != 0 ||
		// !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "15", "17"}) ||
		// !reflect.DeepEqual(deltaSet.stagedRemovals.Elements(), []string{}) ||
		!reflect.DeepEqual(deltaSet.stagedAdditions.Elements(), []string{"18", "19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.stagedAdditions.Elements())
	}

	deltaSet.Commit(nil)
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "15", "17", "18", "19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	v, ok := deltaSet.PopLast()
	if !ok || v != "21" || deltaSet.NonNilCount() != 6 {
		t.Error("failed to commit", v, deltaSet.NonNilCount())
	}

	v, ok = deltaSet.PopLast()
	if !ok || v != "20" || deltaSet.NonNilCount() != 5 {
		t.Error("failed to commit", v, deltaSet.NonNilCount())
	}

	v, ok = deltaSet.GetByIndex(5)
	if ok || v != "" || deltaSet.NonNilCount() != 5 {
		t.Error("Should not exist", v, deltaSet.NonNilCount())
	}
}

func TestDeleteAllThenAddBack(t *testing.T) {
	deltaSet := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})
	deltaSet.Commit(nil) // The strings are in the committed set already

	deltaSet.InsertBatch([]string{"18", "19", "20", "21"}) // { "13", "15", "17"} + { "18", "19", "20", "21"}

	deltaSet.DeleteByIndex(1) //
	deltaSet.DeleteByIndex(4) //
	deltaSet.DeleteByIndex(5) //  {"15", "19", "20"} are in the stagedRemovals set

	deltaSet.DeleteAll() // This will move the stagedRemovals to committed set

	if common.FilterFirst(deltaSet.Exists("13")).(bool) ||
		common.FilterFirst(deltaSet.Exists("17")).(bool) ||
		common.FilterFirst(deltaSet.Exists("15")).(bool) ||
		common.FilterFirst(deltaSet.Exists("18")).(bool) ||
		common.FilterFirst(deltaSet.Exists("19")).(bool) ||
		common.FilterFirst(deltaSet.Exists("20")).(bool) ||
		common.FilterFirst(deltaSet.Exists("21")).(bool) {
		t.Error("18 or 21 should not exist after DeleteAll", deltaSet.stagedRemovals.Elements())
	}

	if !reflect.DeepEqual(deltaSet.stagedRemovals.Elements(), []string{"13", "15", "17", "18", "19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Elements())
	}

	// Make sure the stagedRemovals is as expected.
	if !reflect.DeepEqual(deltaSet.stagedRemovals.Added().Elements(), []string{"18", "19", "20", "21"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Committed().Elements(), []string{"13", "15", "17"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Removed().Elements(), []string{}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Removed())
	}

	deltaSet.InsertBatch([]string{"15", "19", "20"}) // Add the deleted entires back to the set

	if !reflect.DeepEqual(deltaSet.stagedRemovals.Added().Elements(), []string{"18", "19", "20", "21"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Committed().Elements(), []string{"13", "15", "17"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Removed().Elements(), []string{"15", "19", "20"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Removed(), []string{"15", "19", "20"})
	}

	deltaSet.InsertBatch([]string{"15", "19", "20", "21"})
	if !reflect.DeepEqual(deltaSet.stagedRemovals.Added().Elements(), []string{"18", "19", "20", "21"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Committed().Elements(), []string{"13", "15", "17"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Removed().Elements(), []string{"15", "19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals.Removed(), []string{"15", "19", "20"})
	}

	deltaSet.Insert("22")
	if !reflect.DeepEqual(deltaSet.stagedRemovals.Added().Elements(), []string{"18", "19", "20", "21"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Committed().Elements(), []string{"13", "15", "17"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Removed().Elements(), []string{"15", "19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals)
	}

	if common.FilterFirst(deltaSet.Exists("13")).(bool) {
		t.Error("failed to commit", deltaSet.Elements())
	}

	deltaSet.Delete("15")
	if !reflect.DeepEqual(deltaSet.stagedRemovals.Added().Elements(), []string{"18", "19", "20", "21"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Committed().Elements(), []string{"13", "15", "17"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Removed().Elements(), []string{"19", "20", "21"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals)
	}

	deltaSet.Delete("20")
	if !reflect.DeepEqual(deltaSet.stagedRemovals.Added().Elements(), []string{"18", "19", "20", "21"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Committed().Elements(), []string{"13", "15", "17"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Removed().Elements(), []string{"19", "21"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals)
	}

	deltaSet.Insert("20")
	if !reflect.DeepEqual(deltaSet.stagedRemovals.Added().Elements(), []string{"18", "19", "20", "21"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Committed().Elements(), []string{"13", "15", "17"}) ||
		!reflect.DeepEqual(deltaSet.stagedRemovals.Removed().Elements(), []string{"19", "21", "20"}) {
		t.Error("failed to commit", deltaSet.stagedRemovals)
	}
}

func TestMultiSoftSetMerge(t *testing.T) {
	deltaSet := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})
	deltaSet.Commit(nil)

	_set0 := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).InsertBatch([]string{"58", "59", "20", "51"}).DeleteBatch([]string{"13"})
	_set1 := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).InsertBatch([]string{"78", "59", "70", "71"}).DeleteBatch([]string{"15", "70"})

	// ("13", "15", "17") + ("58", "59", "20", "51") + ("78", "59", "70", "71") - ("13", "15", "70") = ("17", "58", "59", "20", "51", "78", "59", "71")
	deltaSet.Commit([]*SoftDeltaSet[string]{_set0, _set1})

	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "15", "17", "58", "59", "20", "51", "78", "71"}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
}

func TestMultiSoftSetMergeWithStagedRemovals(t *testing.T) {
	_set0 := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).InsertBatch([]string{"58", "59", "20", "51"}).DeleteBatch([]string{"13"})       // "13" will be ignored.
	_set1 := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).InsertBatch([]string{"78", "59", "70", "71"}).DeleteBatch([]string{"15", "70"}) // 15 will be ignored.

	deltaSet := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})
	deltaSet.Commit(nil)
	deltaSet.Commit([]*SoftDeltaSet[string]{_set0, _set1})

	// ("13", "15", "17") + ("58", "59", "20", "51") + ("78", "59", "70", "71") - ("13", "15", "70") = ("17", "58", "59", "20", "51", "78", "59", "71")
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "15", "17", "58", "59", "20", "51", "78", "71"}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	_set2 := deltaSet.CloneFull()
	_set3 := deltaSet.CloneFull()

	_set2.InsertBatch([]string{"100", "101", "102"})
	_set3.InsertBatch([]string{"777", "888"})

	_set2.DeleteAll()                                 // will remove {"13", "15", "17", "58", "59", "20", "51", "78", "71"}
	_set2.InsertBatch([]string{"15", "++++", "9999"}) // Add "15" back to the set, Add "++++", "9999" to the set

	deltaSet.Commit([]*SoftDeltaSet[string]{_set2, _set3})
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"15", "++++", "9999", "777", "888"}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

}

func TestSoftDeltaGetNthNonNil(t *testing.T) {
	deltaSet := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})

	deltaSet.InsertBatch([]string{"18", "19", "20", "21"}) //  { "13", "15", "17"} +  { "18", "19", "20", "21"}
	// DeleteByIndex wouldn't shift the indices
	deltaSet.DeleteByIndex(1) //  { "13", "15", "17"} + { "18", "19", "20", "21"} - {"15"} = { "13", "17", "18", "19", "20", "21"}
	deltaSet.DeleteByIndex(4) //  { "13", "15", "17"} + { "18", "19", "20", "21"} - {"15", "19"}

	if k, idx, ok := deltaSet.GetNthNonNil(0); k != "13" || idx != 0 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(1); k != "17" || idx != 2 || !ok {
		t.Error("failed to commit", k)
	}

	// Check if the deleted entry is still accessible.
	if k, ok := deltaSet.GetByIndex(1); k != "" || ok {
		t.Error("A deleted entry shouldn't be available any more", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(2); k != "18" || idx != 3 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(3); k != "20" || idx != 5 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(4); k != "21" || idx != 6 || !ok {
		t.Error("failed to commit", k)
	}

	// Check if the deleted entry is still accessible.
	if k, ok := deltaSet.GetByIndex(4); k != "" || ok {
		t.Error("A deleted entry shouldn't be available any more", k)
	}

	deltaSet.Commit(nil)
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "17", "18", "20", "21"}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	if k, idx, ok := deltaSet.GetNthNonNil(0); k != "13" || idx != 0 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(1); k != "17" || idx != 1 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(2); k != "18" || idx != 2 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(3); k != "20" || idx != 3 || !ok {
		t.Error("failed to commit", k)
	}

	if k, idx, ok := deltaSet.GetNthNonNil(4); k != "21" || idx != 4 || !ok {
		t.Error("failed to commit", k)
	}
}

func TestSoftDeltaSetCodec(t *testing.T) {
	_set0 := NewSoftDeltaSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).InsertBatch([]string{"58", "59", "20", "51"}).DeleteBatch([]string{"13"})       // "13" will be ignored.
	_set1 := NewSoftDeltaSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).InsertBatch([]string{"78", "59", "70", "71"}).DeleteBatch([]string{"15", "70"}) // 15 will be ignored.

	deltaSet := NewSoftDeltaSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})
	deltaSet.Commit(nil)
	deltaSet.Commit([]*SoftDeltaSet[string]{_set0, _set1})

	// ("13", "15", "17") + ("58", "59", "20", "51") + ("78", "59", "70", "71") - ("13", "15", "70") = ("17", "58", "59", "20", "51", "78", "59", "71")
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"13", "15", "17", "58", "59", "20", "51", "78", "71"}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}

	_set2 := deltaSet.CloneFull()
	_set3 := deltaSet.CloneFull()

	_set2.InsertBatch([]string{"100", "101", "102"})
	_set3.InsertBatch([]string{"777", "888"})

	_set2.DeleteAll()                                 // will remove {"13", "15", "17", "58", "59", "20", "51", "78", "71"}
	_set2.InsertBatch([]string{"15", "++++", "9999"}) // Add "15" back to the set, Add "++++", "9999" to the set

	buffer2 := _set2.Encode()
	buffer3 := _set3.Encode()

	out2 := NewSoftDeltaSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).Decode(buffer2).(*SoftDeltaSet[string])
	out3 := NewSoftDeltaSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).Decode(buffer3).(*SoftDeltaSet[string])

	deltaSet.Commit([]*SoftDeltaSet[string]{out2, out3})
	if !reflect.DeepEqual(deltaSet.committed.Elements(), []string{"15", "++++", "9999", "777", "888"}) {
		t.Error("failed to commit", deltaSet.committed.Elements())
	}
}

// func sizer(K string) int                      { return 8 }
// func encodeToBuffer(K string, buf []byte) int { return codec.String(K).EncodeTo(buf) }
// func decoder(buf []byte) string               { return string(codec.String(buf).Decode(buf).(codec.String)) }

func BenchmarkDeltaSoftDeleteThenAddBack(t *testing.B) {
	deltaSet := NewSoftDeltaSet("", 1000000, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	randoms := make([]string, 1000000)
	for i := 0; i < 1000000; i++ {
		randoms[i] = fmt.Sprintf("%d", i) //rand.Int()
		// randoms[i] = i //rand.Int()
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

func BenchmarkSoftGetNthNonNil(b *testing.B) {
	deltaSet := NewSoftDeltaSet[string]("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	deltaSet.InsertBatch([]string{"13", "15", "17"})

	deltaSet.InsertBatch([]string{"18", "19", "20", "21"}) // { "13", "15", "17"} + { "18", "19", "20", "21"}
	deltaSet.DeleteByIndex(1)                              //  { "13", -"15", "17"} + { "18", "19", "20", "21"}
	deltaSet.DeleteByIndex(4)                              // { "13", -"15", "17"} + { "18", -"19", "20", "21"}'

	total := 10000
	for i := 0; i < total; i++ {
		deltaSet.InsertBatch([]string{fmt.Sprintf("%d", i)})
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
