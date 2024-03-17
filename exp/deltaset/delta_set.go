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
	orderedset "github.com/arcology-network/common-lib/exp/orderedset"
	"github.com/arcology-network/common-lib/exp/slice"
)

// DeltaSet represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type DeltaSet[K comparable] struct {
	nilVal K

	committed *orderedset.OrderedSet[K]
	updated   *orderedset.OrderedSet[K] // New entires and updated entries
	removed   *orderedset.OrderedSet[K] // Entries to be removed including the newly
}

// NewIndexedSlice creates a new instance of DeltaSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewDeltaSet[K comparable](nilVal K, preAlloc int, keys ...K) *DeltaSet[K] {
	deltaSet := &DeltaSet[K]{
		nilVal:    nilVal,
		committed: orderedset.NewOrderedSet(nilVal, preAlloc),
		updated:   orderedset.NewOrderedSet(nilVal, preAlloc),
		removed:   orderedset.NewOrderedSet(nilVal, preAlloc),
	}
	deltaSet.Insert(keys...)
	return deltaSet
}

func (*DeltaSet[K]) New(
	nilVal K,
	committed *orderedset.OrderedSet[K],
	updated *orderedset.OrderedSet[K],
	removed *orderedset.OrderedSet[K]) *DeltaSet[K] {

	return &DeltaSet[K]{
		nilVal:    nilVal,
		committed: committed,
		updated:   updated,
		removed:   removed,
	}
}

func (this *DeltaSet[K]) mapTo(idx int) (*orderedset.OrderedSet[K], int) {
	if idx >= int(this.Length()) {
		return nil, -1
	}

	// The index is in the updated list
	if idx >= this.committed.Length() {
		return this.updated, idx - this.committed.Length()
	}
	return this.committed, idx
}

// Array returns the underlying slice of committed in the DeltaSet.
func (this *DeltaSet[K]) Committed() *orderedset.OrderedSet[K] { return this.committed }
func (this *DeltaSet[K]) Removed() *orderedset.OrderedSet[K]   { return this.removed }
func (this *DeltaSet[K]) Added() *orderedset.OrderedSet[K]     { return this.updated }

// Elements returns the underlying slice of committed in the DeltaSet,
// Non-nil values are returned in the order they were inserted, equals to committed + updated - removed.
func (this *DeltaSet[K]) Elements() []K {
	elements := make([]K, 0, this.NonNilCount())
	for i := 0; i < this.committed.Length()+this.updated.Length(); i++ {
		if v, ok := this.GetByIndex(uint64(i)); ok {
			elements = append(elements, v)
		}
	}
	return elements
}

// Get Byte size of the DeltaSet
func (this *DeltaSet[K]) Size(getter func(K) int) int {
	return common.IfThenDo1st(this.committed != nil, func() int { return this.committed.Size(getter) }, 0) +
		common.IfThenDo1st(this.updated != nil, func() int { return this.updated.Size(getter) }, 0) +
		common.IfThenDo1st(this.removed != nil, func() int { return this.removed.Size(getter) }, 0)
}

func (this *DeltaSet[K]) Clear() {
	this.updated.Clear()
	this.removed.Clear()
}

// Debugging only
func (this *DeltaSet[K]) SetCommitted(v []K) { this.committed.Insert(v...) }
func (this *DeltaSet[K]) SetRemoved(v []K)   { this.removed.Insert(v...) }
func (this *DeltaSet[K]) SetAppended(v []K)  { this.updated.Insert(v...) }
func (this *DeltaSet[K]) SetNilVal(v K)      { this.nilVal = v }

// IsDirty returns true if the DeltaSet is up to date, with no
func (this *DeltaSet[K]) IsDirty() bool {
	return this.removed.Length() == 0 && this.updated.Length() == 0
}

// Delta returns a new instance of DeltaSet with the same updated and removed elements only.
func (this *DeltaSet[K]) Delta() *DeltaSet[K] {
	return this.CloneDelta()
}

// SetDelta sets the updated and removed lists to the specified DeltaSet.
func (this *DeltaSet[K]) SetDelta(delta *DeltaSet[K]) {
	this.updated = delta.updated
	this.removed = delta.removed
}

// ResetDelta resets the updated and removed lists to empty.
func (this *DeltaSet[K]) ResetDelta() {
	this.updated.Clear()
	this.removed.Clear()
}

// Length returns the number of elements in the DeltaSet, including the NIL values.
func (this *DeltaSet[K]) Length() uint64 {
	return uint64(this.committed.Length() + this.updated.Length())
}

// NonNilCount returns the number of NON-NIL elements in the DeltaSet.
func (this *DeltaSet[K]) NonNilCount() uint64 {
	return uint64(this.committed.Length() + this.updated.Length() - this.removed.Length())
}

// Insert inserts an element into the DeltaSet and updates the index.
func (this *DeltaSet[K]) Insert(elems ...K) *DeltaSet[K] {
	for _, elem := range elems {
		if ok, _ := this.removed.Exists(elem); ok { // If the element is in the removed list, move it to the updated list
			this.removed.Delete(elem)
			this.updated.Insert(elem)
			continue
		}

		// Not in the committed list, add it to the updated list. It is possible
		// that the element is already in the updated list, just add it anyway.
		if ok, _ := this.committed.Exists(elem); !ok {
			this.updated.Insert(elem) // Either in the updated list or not.
		}
	}
	return this
}

// Insert inserts an element into the DeltaSet and updates the index.
func (this *DeltaSet[K]) Delete(elems ...K) *DeltaSet[K] {
	for _, elem := range elems {
		if ok, _ := this.removed.Exists(elem); !ok {
			this.removed.Insert(elem) // Impossible to have duplicate entries possible, since the removed list is a set.
		}
	}
	return this
}

// Clone returns a new instance of DeltaSet with the same elements.
func (this *DeltaSet[K]) CloneFull() *DeltaSet[K] {
	set := this.CloneDelta()
	set.committed = this.committed.Clone()
	return set
}

// Clone returns a new instance with the
// committed list shared original DeltaSet.
func (this *DeltaSet[K]) Clone(v ...*orderedset.OrderedSet[K]) *DeltaSet[K] {
	set := this.CloneDelta(v...)
	set.committed = this.committed //.Clone()
	return set
}

// CloneDelta returns a new instance of DeltaSet with the
// same updated and removed elements only.the committed list is not cloned.
func (this *DeltaSet[K]) CloneDelta(v ...*orderedset.OrderedSet[K]) *DeltaSet[K] {
	set := &DeltaSet[K]{
		nilVal: this.nilVal,
		// committed: orderedset.NewOrderedSet(this.nilVal, 0),
		updated: this.updated.Clone(),
		removed: this.removed.Clone(),
	}

	if len(v) > 0 {
		set.committed = v[0]
	} else {
		set.committed = orderedset.NewOrderedSet(this.nilVal, 0)
	}
	return set
}

// CloneDelta returns a new instance of DeltaSet with the
// same updated and removed elements only.the committed list is not cloned.
// func (this *DeltaSet[K]) CloneDelta() *DeltaSet[K] {
// 	set := &DeltaSet[K]{
// 		nilVal:    this.nilVal,
// 		committed: orderedset.NewOrderedSet(this.nilVal, 0),
// 		updated:   this.updated.Clone(),
// 		removed:   this.removed.Clone(),
// 	}
// 	return set
// }

func (this *DeltaSet[K]) DeleteByIndex(idx uint64) {
	if k, _, _, ok := this.Search(idx); ok {
		this.Delete(k)
	}
}

// GetByIndex returns the element at the specified index.
// It does NOT check if the index is corresponding to a non-nil value.
// If the value at the index is nil, the nil value is returned.
func (this *DeltaSet[K]) GetByIndex(idx uint64) (K, bool) {
	if k, _, _, ok := this.Search(idx); ok {
		if ok, _ := this.removed.Exists(k); !ok {
			return k, true
		}
	}
	return *new(K), false
}

// NthNonNil returns the nth non-nil value from the DeltaSet.
// The nth non-nil value isn't necessarily the nth value in the DeltaSet, but the nth non-nil value.
func (this *DeltaSet[K]) GetNthNonNil(nth uint64) (K, int, bool) {
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

// Back get the last non-nil value from the DeltaSet.
// The last non-nil value isn't necessarily the last value in the DeltaSet, but the last non-nil value.
func (this *DeltaSet[K]) Last() (K, bool) {
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

func (this *DeltaSet[K]) Back() (K, bool) {
	if this.NonNilCount() == 0 {
		return *new(K), false
	}
	return this.GetByIndex(uint64(this.Length() - 1))
}

// Return then remove the last non-nil value from the DeltaSet.
func (this *DeltaSet[K]) PopLast() (K, bool) {
	k, ok := this.Last() // Get the last non-nil value
	if ok {
		this.Delete(k) // Re
	}
	return k, ok
}

// Search returns the element at the specified index and the set it is in.
func (this *DeltaSet[K]) Search(idx uint64) (K, *orderedset.OrderedSet[K], int, bool) {
	set, mapped := this.mapTo(int(idx))
	if mapped < 0 {
		return this.nilVal, nil, -1, false
	}
	return set.IndexToKey(mapped), set, mapped, true
}

func (this *DeltaSet[K]) KeyAt(idx uint64) K {
	k, _, _, _ := this.Search(idx)
	return k
}

// Get the index of the element in the DeltaSet by key.
func (this *DeltaSet[K]) IdxOf(k K) uint64 {
	if ok, idx := this.Exists(k); ok {
		return uint64(idx)
	}
	return math.MaxUint64
}

// Get returns the element at the specified index.
func (this *DeltaSet[K]) TryGetKey(idx uint64) (K, bool) {
	k, _, _, ok := this.Search(idx)
	if ok {
		if ok, _ := this.removed.Exists(k); ok { // In the removed set
			return this.nilVal, false
		}
		return k, true
	}
	return k, false
}

// Get returns the element at the specified index. If the index is out of range, the nil value is returned.
func (this *DeltaSet[K]) Exists(k K) (bool, int) {
	if ok, _ := this.removed.Exists(k); ok {
		return false, -1 // Has been removed already
	}

	if ok, v := this.committed.Exists(k); ok {
		return ok, v
	}

	if ok, v := this.updated.Exists(k); ok {
		return ok, v
	}

	return false, -1
}

// Commit commits the updated and removed lists to the committed list.
// Commit assumes that the updated and removed lists are disjoint.
// func (this *DeltaSet[K]) Commit(other ...*DeltaSet[K]) *DeltaSet[K] {
// 	updatedElems := slice.Transform(other, func(_ int, v *DeltaSet[K]) []K {
// 		return v.updated.Elements()
// 	})
// 	slice.PushFront(&updatedElems, this.updated.Elements())

// 	removed := slice.Transform(other, func(_ int, v *DeltaSet[K]) []K {
// 		return v.removed.Elements()
// 	})
// 	slice.PushFront(&removed, this.removed.Elements())

// 	this.committed.Merge(slice.Flatten(updatedElems)) // Merge the updated list to the committed list
// 	this.committed.Sub(slice.Flatten(removed))        // Remove the removed list from the committed list
// 	this.ResetDelta()                                 // Reset the updated and removed lists to empty
// 	fmt.Println(updatedElems)
// 	fmt.Println(removed)
// 	return this
// }

func (this *DeltaSet[K]) Commit(other ...*DeltaSet[K]) *DeltaSet[K] {
	updateBuffer := make([]K, slice.CountDo(other, func(_ int, v **DeltaSet[K]) int { return (*v).updated.Length() })+this.updated.Length())
	slice.ConcateToBuffer(other, &updateBuffer, func(v *DeltaSet[K]) []K { return v.updated.Elements() })
	copy(updateBuffer[len(updateBuffer)-this.updated.Length():], this.updated.Elements())

	removedBuffer := make([]K, slice.CountDo(other, func(_ int, v **DeltaSet[K]) int { return (*v).removed.Length() })+this.removed.Length())
	slice.ConcateToBuffer(other, &removedBuffer, func(v *DeltaSet[K]) []K { return v.removed.Elements() })
	copy(removedBuffer[len(removedBuffer)-this.removed.Length():], this.removed.Elements())

	this.committed.Merge(updateBuffer) // Merge the updated list to the committed list
	this.committed.Sub(removedBuffer)  // Remove the removed list from the committed list
	this.ResetDelta()                  // Reset the updated and removed lists to empty	// return this

	return this
}

func (this *DeltaSet[K]) Equal(other *DeltaSet[K]) bool {
	return this.committed.Equal(other.committed) &&
		this.updated.Equal(other.updated) &&
		this.removed.Equal(other.removed)
}

func (this *DeltaSet[K]) Print() {
	this.committed.Print()
	this.updated.Print()
	this.removed.Print()
}
