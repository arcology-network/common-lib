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

package deltaslice

import (
	"github.com/arcology-network/common-lib/exp/associative"
	indexedslice "github.com/arcology-network/common-lib/exp/indexed"
	"github.com/arcology-network/common-lib/exp/slice"
)

// DeltaSet represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type DeltaSet[K comparable] struct {
	modified *indexedslice.IndexedSlice[K, K, K] // A hybrid structure to record the modified elements.
	deleter  func(*K) K
	isEmpty  func(K) bool

	appended    *indexedslice.IndexedSlice[K, K, K]
	readonlys   *indexedslice.IndexedSlice[K, K, K]
	CommitedLen int
}

// NewIndexedSlice creates a new instance of DeltaSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewDeltaSet[K comparable](modified *indexedslice.IndexedSlice[K, K, K],
	deleter func(*K) K,
	isEmpty func(K) bool,
) *DeltaSet[K] {
	return &DeltaSet[K]{
		modified: modified,
		deleter:  deleter,
		isEmpty:  isEmpty,
		// lookup:         map[K]*K{},
		readonlys: modified.New(),
		appended:  modified.New(),
	}
}

func (this *DeltaSet[K]) mapTo(idx int) (*indexedslice.IndexedSlice[K, K, K], int) {
	if idx >= this.Length() {
		return nil, -1
	}

	// The index is in the appended list
	if len(*this.readonlys.Elements()) <= idx {
		return this.appended, idx - len(*this.readonlys.Elements())
	}
	return this.readonlys, idx
}

// Array returns the underlying slice of readonlys in the DeltaSet.
// func (this *DeltaSet[K]) Values() []K                                   { return this.readonlys }
// func (this *DeltaSet[K]) Updated() []K                                  { return this.modified.Values() }
func (this *DeltaSet[K]) Appended() []K {
	v := associative.Triplets[*K, uint64, K](*this.appended.Elements()).Thirds()
	return v

} // Returns the appended readonlys, some may be removed later.
// func (this *DeltaSet[K]) Modified() *indexedslice.IndexedSlice[K, K, K] { return this.modified }
// func (this *DeltaSet[K]) Dict() map[K]*K { return this.lookup }

func (this *DeltaSet[K]) Length() int {
	elems := *this.modified.Elements()
	numRemoved := slice.CountIf[*associative.Triplet[*K, uint64, K], int](elems, func(_ int, v **associative.Triplet[*K, uint64, K]) bool {
		return *v == nil
	})
	return len(*this.readonlys.Elements()) + len(*this.appended.Elements()) - numRemoved
}

// Insert inserts an element into the DeltaSet and updates the index.
func (this *DeltaSet[K]) Append(keyGetter func(K) K, elems ...K) int {
	for i := 0; i < len(elems); i++ {
		newv := &associative.Triplet[*K, uint64, K]{
			First:  &elems[i],
			Second: uint64(this.Length()),
			Third:  elems[i],
		}
		this.appended.Append(newv)
	}
	return this.Length()
}

// ToSlice returns the readonlys in the DeltaSet as a slice by removing the removed readonlys and adding the appended readonlys.
func (this *DeltaSet[K]) ToSlice() []K {
	appLen := this.appended.Length(func(K) int { return 1 })
	readonlyLen := len(*this.readonlys.Elements())

	elems := make([]K, appLen+readonlyLen)
	for i := 0; i < readonlyLen; i++ {
		elems[i], _ = this.readonlys.GetByIndex(i)
	}

	for i := 0; i < appLen; i++ {
		elems[i+readonlyLen], _ = this.appended.GetByIndex(i)
	}
	return elems
}

// Insert inserts an element into the DeltaSet and updates the index.
func (this *DeltaSet[K]) Delete(indices ...int) {
	for _, idx := range indices {
		this.SetByIndex(idx, this.deleter(new(K)))
	}
}

// Set sets the element at the specified index to the new value.
func (this *DeltaSet[K]) SetByIndex(idx int, newk K) bool {
	if idx >= this.Length() {
		return false
	}

	arr, mapped := this.mapTo(idx)
	if mapped < 0 {
		return false
	}

	if arr == this.readonlys {
		k, _ := this.readonlys.GetByIndex(idx)
		this.modified.SetByKey(k, newk)
		return true
	}

	k, _ := this.appended.GetByIndex(mapped)
	this.modified.SetByKey(k, newk)

	// if mapped >= 0 {
	// 	this.modified.SetByKey(common.ToType[int, K](idx), newk)
	// 	(*arr.Elements())[mapped].Third = newk
	// 	return true
	// }
	return false
}

// Get returns the element at the specified index.
func (this *DeltaSet[K]) Get(idx int) (K, bool) {
	arr, mapped := this.mapTo(idx)
	if mapped < 0 {
		return *new(K), false
	}

	k, _ := arr.IndexToKey(uint64(mapped))
	if v, _, ok := this.modified.GetByKey(k); ok {
		return v, ok
	}
	return arr.GetByIndex(mapped)
}

// GetByKey returns the element with the specified key.
func (this *DeltaSet[K]) GetByKey(k K) (K, bool) {
	if v, _, ok := this.modified.GetByKey(k); ok { // Get from the modified first
		return v, true
	}

	if v, _, ok := this.appended.GetByKey(k); ok { // Get from the appended first
		return v, true
	}

	if v, _, ok := this.readonlys.GetByKey(k); ok { // Get from the modified first
		return v, true
	}
	return *new(K), false
}

// SetByKey sets the element with the specified key to the new value.
// New elements will be appended. Only the modified elements will be updated.
func (this *DeltaSet[K]) SetByKey(k K) {
	if _, _, ok := this.modified.GetByKey(k); ok {
		this.modified.SetByKey(k, k) // Previously modified elements will be updated.
		return
	}

	if _, idx, ok := this.appended.GetByKey(k); ok {
		this.modified.SetByIndex(int(idx)+len(*this.readonlys.Elements()), k) // Previously apppended elements will be save to modified set.
		return
	}

	if _, idx, ok := this.readonlys.GetByKey(k); ok {
		this.modified.SetByIndex(int(idx), k) // First time modified elements will be save to modified set.
		return
	}
}

// DoAt calls the specified function with the element at the specified index.
// The operation is an in-place operation, so it can modify the element in the DeltaSet, even
// if the element is in the readonlys.
// func (this *DeltaSet[K]) DoAt(idx int, doer func(*K)) {
// 	arr, mapped := this.mapTo(idx)
// 	if mapped < 0 {
// 		return
// 	}
// 	doer(&(*arr)[mapped])
// }

func (this *DeltaSet[K]) Commit() *DeltaSet[K] {
	this.readonlys.Merge(this.appended)
	this.appended.Clear()

	elements := this.modified.Elements()
	for i := 0; i < len(*elements); i++ {
		this.readonlys.SetByIndex(int((*elements)[i].Second), (*elements)[i].Third)
	}

	dict := this.readonlys.Index()
	slice.RemoveIf(this.readonlys.Elements(), func(_ int, v *associative.Triplet[*K, uint64, K]) bool {
		if this.isEmpty(v.Third) {
			delete(dict, *v.First)
			return true
		}
		return false
	})

	this.modified.Clear()
	this.CommitedLen = len(*this.readonlys.Elements())
	return this
}
