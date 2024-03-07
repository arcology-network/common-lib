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

package deltaset

import (
	orderedset "github.com/arcology-network/common-lib/exp/orderedset"
	"github.com/arcology-network/common-lib/exp/slice"
)

// DeltaSet represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type DeltaSet[K comparable] struct {
	nilVal K

	readonlys *orderedset.OrderedSet[K]
	updated   *orderedset.OrderedSet[K] // New entires and updated entries
	removed   *orderedset.OrderedSet[K] // Entries to be removed including the newly
}

// NewIndexedSlice creates a new instance of DeltaSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewDeltaSet[K comparable](nilVal K, preAlloc int) *DeltaSet[K] {
	return &DeltaSet[K]{
		nilVal:    nilVal,
		readonlys: orderedset.NewOrderedSet[K](nilVal, preAlloc),
		updated:   orderedset.NewOrderedSet[K](nilVal, preAlloc),
		removed:   orderedset.NewOrderedSet[K](nilVal, preAlloc),
	}
}

func (this *DeltaSet[K]) mapTo(idx int) (*orderedset.OrderedSet[K], int) {
	if idx >= this.Length() {
		return nil, -1
	}

	// The index is in the updated list
	if this.readonlys.Length() <= idx {
		return this.updated, idx - this.readonlys.Length()
	}
	return this.readonlys, idx
}

// Array returns the underlying slice of readonlys in the DeltaSet.
func (this *DeltaSet[K]) Values() []K   { return this.readonlys.Elements() }
func (this *DeltaSet[K]) Modified() []K { return this.removed.Elements() }
func (this *DeltaSet[K]) Appended() []K { return this.updated.Elements() }

func (this *DeltaSet[K]) Length() int {
	elems := this.removed.Elements()
	numRemoved := slice.CountIf[K, int](elems, func(_ int, v *K) bool {
		return *v == this.nilVal
	})
	return this.readonlys.Length() + this.updated.Length() - numRemoved
}

// Insert inserts an element into the DeltaSet and updates the index.
func (this *DeltaSet[K]) Insert(elems ...K) {
	for _, elem := range elems {
		if this.removed.Exists(elem) { // If the element is in the removed list, remove it from the removed list.
			this.removed.Delete(elem)
			this.removed.Sync()

			this.updated.Insert(elem) // Add it back to the updated list
			continue
		}

		// Not in the readonlys list, add it to the updated list. It is possible
		// that the element is already in the updated list, just add it anyway.
		if !this.readonlys.Exists(elem) {
			this.updated.Insert(elem) // Either in the updated list or not.
		}
	}
}

// Insert inserts an element into the DeltaSet and updates the index.
func (this *DeltaSet[K]) Delete(elems ...K) {
	for _, elem := range elems {
		if this.removed.Exists(elem) {
			continue // already deleted
		}
		this.removed.Insert(elem) // Either in the removed list or not.
	}
}

func (this *DeltaSet[K]) Clone() *DeltaSet[K] {
	set := this.CloneDelta()
	set.readonlys = this.readonlys.Clone()
	return set
}

func (this *DeltaSet[K]) CloneDelta() *DeltaSet[K] {
	set := &DeltaSet[K]{
		nilVal:  this.nilVal,
		updated: this.updated.Clone(),
		removed: this.removed.Clone(),
	}
	return set
}

// Set sets the element at the specified index to the new value.
func (this *DeltaSet[K]) SetByIndex(idx int, newk K) bool {
	_, set, mapped, ok := this.IndexToKey(idx)
	if !ok {
		return false
	}

	// Delete the element if the new value is nil
	if newk == this.nilVal {
		this.Delete(*set.At(mapped))
		return true
	}

	// Already in the updated list
	if set == this.updated {
		set.Replace(mapped, newk)
		return true
	}

	// In the readonlys list
	oldk := *set.At(mapped)             // Get the old value from the readonlys list
	pos, _ := this.updated.Insert(oldk) // Add the old value to the updated list

	// Replace the old value with the new value in the updated list. The key and value are no longer the same.
	// at the point. The key represents the old value and the value is the new value. It can be used to update the value.
	*this.updated.At(pos) = newk

	// oldk := set.IndexToKey(mapped)
	_, ok = this.updated.Insert(newk)
	return ok
}

func (this *DeltaSet[K]) DeleteByIndex(idx int) {
	this.SetByIndex(idx, this.nilVal)
}

func (this *DeltaSet[K]) IndexToKey(idx int) (K, *orderedset.OrderedSet[K], int, bool) {
	set, mapped := this.mapTo(idx)
	if mapped < 0 {
		return this.nilVal, nil, -1, false
	}
	return set.IndexToKey(mapped), set, mapped, true
}

// Get returns the element at the specified index.
func (this *DeltaSet[K]) GetByIndex(idx int) (K, bool) {
	k, _, _, ok := this.IndexToKey(idx)
	if ok {
		if this.removed.Exists(k) { // In the removed set
			return this.nilVal, false
		}
	}
	return k, true
}

func (this *DeltaSet[K]) Exists(k K) bool {
	if this.removed.Exists(k) {
		return false
	}
	return this.readonlys.Exists(k) || this.updated.Exists(k)
}

func (this *DeltaSet[K]) Commit() *DeltaSet[K] {
	this.readonlys.Merge(this.updated)
	this.readonlys.Sub(this.removed)

	this.updated.Clear()
	this.removed.Clear()
	return this
}

func (this *DeltaSet[K]) Equal(other *DeltaSet[K]) bool {
	return this.readonlys.Equal(other.readonlys) &&
		this.updated.Equal(other.updated) &&
		this.removed.Equal(other.removed)
}

func (this *DeltaSet[K]) Print() {
	this.readonlys.Print()
	this.updated.Print()
	this.removed.Print()
}
