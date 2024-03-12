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
	"math"

	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/exp/associative"
	orderedset "github.com/arcology-network/common-lib/exp/orderedset"
	"github.com/arcology-network/common-lib/exp/slice"
)

// CascadeDeltaSet represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type CascadeDeltaSet[K comparable] struct {
	nilVal K

	committed *orderedset.OrderedSet[K]
	updated   *associative.Pair[*orderedset.OrderedSet[K], *orderedset.OrderedSet[K]] // New entires and updated entries
	removed   *DeltaSet[K]                                                            // Entries to be removed including the newly
}

// NewIndexedSlice creates a new instance of CascadeDeltaSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewCascadeDeltaSet[K comparable](nilVal K, preAlloc int, keys ...K) *CascadeDeltaSet[K] {
	deltaSet := &CascadeDeltaSet[K]{
		nilVal:    nilVal,
		committed: orderedset.NewOrderedSet(nilVal, preAlloc),
		updated: &associative.Pair[*orderedset.OrderedSet[K], *orderedset.OrderedSet[K]]{
			First:  orderedset.NewOrderedSet(nilVal, preAlloc),
			Second: orderedset.NewOrderedSet(nilVal, preAlloc)},
		removed: NewDeltaSet(nilVal, preAlloc),
	}
	deltaSet.Insert(keys...)
	return deltaSet
}

func (*CascadeDeltaSet[K]) New(
	nilVal K,
	committed *orderedset.OrderedSet[K],
	updated *associative.Pair[*orderedset.OrderedSet[K], *orderedset.OrderedSet[K]],
	removed *DeltaSet[K]) *CascadeDeltaSet[K] {

	return &CascadeDeltaSet[K]{
		nilVal:    nilVal,
		committed: committed,
		updated:   updated,
		removed:   removed,
	}
}

func (this *CascadeDeltaSet[K]) mapTo(idx int) (*orderedset.OrderedSet[K], int) {
	if idx >= int(this.Length()) {
		return nil, -1
	}

	// The index is in the committed list
	if idx < this.committed.Length() {
		return this.committed, idx
	}

	// The index is in the first updated list
	if idx >= this.committed.Length() && idx < this.committed.Length()+this.updated.First.Length() {
		return this.updated.First, idx - this.committed.Length()
	}

	// The index is in the second updated list
	return this.updated.Second, idx - this.committed.Length() - -this.updated.First.Length()
}

// Array returns the underlying slice of committed in the CascadeDeltaSet.
func (this *CascadeDeltaSet[K]) Committed() *orderedset.OrderedSet[K] { return this.committed }

func (this *CascadeDeltaSet[K]) Removed() *orderedset.OrderedSet[K] {
	return this.Removed().Clone().Merge(this.removed.Elements())
}

func (this *CascadeDeltaSet[K]) Added() *orderedset.OrderedSet[K] {
	return this.updated.First.Clone().Merge(this.updated.Second.Elements())
}

// Elements returns the underlying slice of committed in the CascadeDeltaSet,
// Non-nil values are returned in the order they were inserted, equals to committed + updated - removed.
func (this *CascadeDeltaSet[K]) Elements() []K {
	elements := make([]K, 0, this.NonNilCount())
	for i := 0; i < this.committed.Length()+this.updated.First.Length()+this.updated.Second.Length(); i++ {
		if v, ok := this.GetByIndex(uint64(i)); ok {
			elements = append(elements, v)
		}
	}
	return elements
}

// Get Byte size of the CascadeDeltaSet
func (this *CascadeDeltaSet[K]) Size(getter func(K) int) int {
	return common.IfThenDo1st(this.committed != nil, func() int { return this.committed.Size(getter) }, 0) +
		common.IfThenDo1st(this.updated != nil, func() int { return this.updated.First.Size(getter) + this.updated.Second.Size(getter) }, 0) +
		common.IfThenDo1st(this.removed != nil, func() int { return this.removed.Size(getter) }, 0)
}

func (this *CascadeDeltaSet[K]) Clear(v K) {
	this.updated.First.Clear()
	this.updated.Second.Clear()
	this.removed.Clear()
}

// Debugging only
func (this *CascadeDeltaSet[K]) SetCommitted(v []K) { this.committed.Insert(v...) }
func (this *CascadeDeltaSet[K]) SetRemoved(v []K)   { this.removed.Insert(v...) }
func (this *CascadeDeltaSet[K]) SetAppended(v []K)  { this.updated.Second.Insert(v...) }
func (this *CascadeDeltaSet[K]) SetNilVal(v K)      { this.nilVal = v }

// SetDelta sets the updated and removed lists to the specified CascadeDeltaSet.
func (this *CascadeDeltaSet[K]) SetDelta(delta *CascadeDeltaSet[K]) {
	this.updated = delta.updated
	this.removed = delta.removed
}

// IsDirty returns true if the CascadeDeltaSet is up to date, with no
func (this *CascadeDeltaSet[K]) IsDirty() bool {
	return !this.removed.IsDirty() && this.updated.First.Length() == 0 && this.updated.Second.Length() == 0
}

// Delta returns a new instance of CascadeDeltaSet with the same updated and removed elements only.
func (this *CascadeDeltaSet[K]) Delta() *CascadeDeltaSet[K] {
	return this.CloneDelta()
}

// ResetDelta resets the updated and removed lists to empty.
func (this *CascadeDeltaSet[K]) ResetDelta() {
	this.updated.First.Clear()
	this.updated.Second.Clear()
	this.removed.Clear()
}

// Length returns the number of elements in the CascadeDeltaSet, including the NIL values.
func (this *CascadeDeltaSet[K]) Length() uint64 {
	return uint64(this.committed.Length() + this.updated.First.Length() + this.updated.Second.Length())
}

// NonNilCount returns the number of NON-NIL elements in the CascadeDeltaSet.
func (this *CascadeDeltaSet[K]) NonNilCount() uint64 {
	return uint64(
		this.committed.Length() +
			this.updated.First.Length() +
			this.updated.Second.Length() -
			int(this.removed.NonNilCount()))
}

// Insert inserts an element into the CascadeDeltaSet and updates the index.
func (this *CascadeDeltaSet[K]) Insert(elems ...K) {
	for _, elem := range elems {
		if ok, _ := this.removed.Exists(elem); ok { // If the element is in the removed list, move it to the updated list
			this.removed.Delete(elem)
			this.updated.First.Insert(elem)
			continue
		}

		// Not in the committed list, add it to the updated list. It is possible
		// that the element is already in the updated list, just add it anyway.
		if ok, _ := this.committed.Exists(elem); !ok {
			this.updated.First.Insert(elem) // Either in the updated list or not.
		}
	}
}

// Insert inserts an element into the CascadeDeltaSet and updates the index.
func (this *CascadeDeltaSet[K]) Delete(elems ...K) {
	for _, elem := range elems {
		if ok, _ := this.removed.Exists(elem); !ok {
			this.removed.Insert(elem) // Impossible to have duplicate entries possible, since the removed list is a set.
		}
	}
}

// Clone returns a new instance of CascadeDeltaSet with the same elements.
func (this *CascadeDeltaSet[K]) CloneFull() *CascadeDeltaSet[K] {
	set := this.CloneDelta()
	set.committed = this.committed.Clone()
	return set
}

// Clone returns a new instance with the
// committed list shared original CascadeDeltaSet.
func (this *CascadeDeltaSet[K]) Clone() *CascadeDeltaSet[K] {
	set := this.CloneDelta()
	set.committed = this.committed //.Clone()
	return set
}

// CloneDelta returns a new instance of CascadeDeltaSet with the
// same updated and removed elements only.the committed list is not cloned.
func (this *CascadeDeltaSet[K]) CloneDelta() *CascadeDeltaSet[K] {
	set := &CascadeDeltaSet[K]{
		nilVal:    this.nilVal,
		committed: orderedset.NewOrderedSet(this.nilVal, 0),
		updated: &associative.Pair[*orderedset.OrderedSet[K], *orderedset.OrderedSet[K]]{
			First:  this.updated.First.Clone(),
			Second: this.updated.Second.Clone(),
		},
		removed: this.removed.Clone(),
	}
	return set
}

func (this *CascadeDeltaSet[K]) DeleteByIndex(idx uint64) {
	if k, _, _, ok := this.Search(idx); ok {
		this.Delete(k)
	}
}

// GetByIndex returns the element at the specified index.
// It does NOT check if the index is corresponding to a non-nil value.
// If the value at the index is nil, the nil value is returned.
func (this *CascadeDeltaSet[K]) GetByIndex(idx uint64) (K, bool) {
	if k, _, _, ok := this.Search(idx); ok {
		if ok, _ := this.removed.Exists(k); !ok {
			return k, true
		}
	}
	return *new(K), false
}

// NthNonNil returns the nth non-nil value from the CascadeDeltaSet.
// The nth non-nil value isn't necessarily the nth value in the CascadeDeltaSet, but the nth non-nil value.
func (this *CascadeDeltaSet[K]) GetNthNonNil(nth uint64) (K, int, bool) {
	// If the nth value is out of range, no need to search. The nil value is returned.
	if nth >= this.NonNilCount() {
		return *new(K), -1, false
	}

	// This isn't efficient; it is better to search from the beginning or the min and max indices of the removed list to
	// narrow down the search range. However, that requires some extra code to keep track of the corresponding indices of
	// the keys in the removed list.
	cnt := 0
	for i := 0; i < int(this.Length()); i++ {
		if K, ok := this.GetByIndex(uint64(i)); ok {
			if uint64(cnt) == nth {
				return K, i, true
			}
			cnt++
		}
	}
	return *new(K), -1, false

	// dict := this.removed.Dict()
	// _, minv := mapi.MinValue(dict, func(l int, r int) bool { return l < r })
	// _, maxv := mapi.MaxValue(dict, func(l int, r int) bool { return l > r })

	// if nth < uint64(minv) {
	// 	if k, ok := this.GetByIndex(nth); ok {
	// 		return k, int(nth), true
	// 	}
	// }

	// if nth > uint64(maxv) {
	// 	if k, ok := this.GetByIndex(nth - uint64(this.removed.Length())); ok {
	// 		return k, int(nth), true
	// 	}
	// }

}

// Back get the last non-nil value from the CascadeDeltaSet.
// The last non-nil value isn't necessarily the last value in the CascadeDeltaSet, but the last non-nil value.
func (this *CascadeDeltaSet[K]) Last() (K, bool) {
	if this.NonNilCount() == 0 {
		return *new(K), false
	}

	for i := this.Length() - 1; i >= 0; i-- {
		if k, ok := this.GetByIndex(uint64(i)); ok {
			return k, true
		}
	}
	return *new(K), false
}

func (this *CascadeDeltaSet[K]) Back() (K, bool) {
	if this.NonNilCount() == 0 {
		return *new(K), false
	}
	return this.GetByIndex(uint64(this.Length() - 1))
}

// Return then remove the last non-nil value from the CascadeDeltaSet.
func (this *CascadeDeltaSet[K]) PopLast() (K, bool) {
	k, ok := this.Last() // Get the last non-nil value
	if ok {
		this.Delete(k) // Re
	}
	return k, ok
}

// Search returns the element at the specified index and the set it is in.
func (this *CascadeDeltaSet[K]) Search(idx uint64) (K, *orderedset.OrderedSet[K], int, bool) {
	set, mapped := this.mapTo(int(idx))
	if mapped < 0 {
		return this.nilVal, nil, -1, false
	}
	return set.KeyToIndex(mapped), set, mapped, true
}

func (this *CascadeDeltaSet[K]) KeyAt(idx uint64) K {
	k, _, _, _ := this.Search(idx)
	return k
}

// Get the index of the element in the CascadeDeltaSet by key.
func (this *CascadeDeltaSet[K]) IdxOf(k K) uint64 {
	if ok, idx := this.Exists(k); ok {
		return uint64(idx)
	}
	return math.MaxUint64
}

// Get returns the element at the specified index.
func (this *CascadeDeltaSet[K]) TryGetKey(idx uint64) (K, bool) {
	k, _, _, ok := this.Search(idx)
	if ok {
		if ok, _ := this.removed.Exists(k); ok { // In the removed set
			return this.nilVal, false
		}
		return k, true
	}
	return k, false
}

// Exists returns true if the element is in the CascadeDeltaSet.
func (this *CascadeDeltaSet[K]) Exists(k K) (bool, int) {
	if ok, _ := this.removed.Exists(k); ok {
		return false, -1 // Has been removed already
	}

	if ok, v := this.committed.Exists(k); ok {
		return ok, v
	}

	if ok, v := this.updated.First.Exists(k); ok {
		return ok, v
	}

	if ok, v := this.updated.Second.Exists(k); ok {
		return ok, v
	}

	return false, -1
}

// Merge the temporary changes to the final change list.
func (this *CascadeDeltaSet[K]) Flush() *CascadeDeltaSet[K] {
	this.updated.Second.Merge(this.updated.First.Elements())
	this.updated.First.Clear()
	this.removed.Commit()
	return this
}

// Commit commits the updated and removed lists to the committed list.
// Commit assumes that the updated and removed lists are disjoint.
func (this *CascadeDeltaSet[K]) Commit(other ...*CascadeDeltaSet[K]) *CascadeDeltaSet[K] {
	updatedElems := slice.Transform(other, func(_ int, v *CascadeDeltaSet[K]) []K {
		return append(v.updated.First.Elements(), v.updated.Second.Elements()...)
	})

	removed := slice.Transform(other, func(_ int, v *CascadeDeltaSet[K]) []K {
		return this.removed.Commit().Elements()
	})

	updatedElems = append(updatedElems, this.updated.First.Elements(), this.updated.Second.Elements())
	removed = append(removed, this.removed.Elements())

	this.committed.Merge(slice.Flatten(updatedElems)) // Merge the updated list to the committed list
	this.committed.Sub(slice.Flatten(removed))        // Remove the removed list from the committed list
	this.ResetDelta()                                 // Reset the updated and removed lists to empty
	return this
}

func (this *CascadeDeltaSet[K]) Equal(other *CascadeDeltaSet[K]) bool {
	return this.committed.Equal(other.committed) &&
		this.updated.First.Equal(other.updated.First) &&
		this.updated.Second.Equal(other.updated.Second) &&
		this.removed.Equal(other.removed)
}

func (this *CascadeDeltaSet[K]) Print() {
	this.committed.Print()
	this.updated.First.Print()
	this.updated.Second.Print()
	this.removed.Print()
}
