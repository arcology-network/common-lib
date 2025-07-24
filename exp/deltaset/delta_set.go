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

// DeltaSet represents a mutable view over a base set, allowing staged additions and deletions.
type DeltaSet[K comparable] struct {
	nilVal     K
	committed  *orderedset.OrderedSet[K] // The stable, committed elements
	stagedAdds *orderedset.OrderedSet[K] // New elements staged for addition
	tombstones *orderedset.OrderedSet[K] // Elements marked for deletion
}

// NewIndexedSlice creates a new instance of DeltaSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewDeltaSet[K comparable](nilVal K, preAlloc int, hasher func(K) [32]byte, keys ...K) *DeltaSet[K] {
	deltaSet := &DeltaSet[K]{
		nilVal:     nilVal,
		committed:  orderedset.NewOrderedSet(nilVal, preAlloc, hasher),
		stagedAdds: orderedset.NewOrderedSet(nilVal, preAlloc, hasher),
		tombstones: orderedset.NewOrderedSet(nilVal, preAlloc, hasher),
	}
	deltaSet.InsertBatch(keys)
	return deltaSet
}

func (*DeltaSet[K]) New(
	nilVal K,
	committed *orderedset.OrderedSet[K],
	stagedAdds *orderedset.OrderedSet[K],
	tombstones *orderedset.OrderedSet[K]) *DeltaSet[K] {

	return &DeltaSet[K]{
		nilVal:     nilVal,
		committed:  committed,
		stagedAdds: stagedAdds,
		tombstones: tombstones,
	}
}

func (*DeltaSet[K]) NewFrom(other *DeltaSet[K]) *DeltaSet[K] {
	return &DeltaSet[K]{
		nilVal:     other.GetNilVal(),
		committed:  orderedset.NewFrom(other.committed),
		stagedAdds: orderedset.NewFrom(other.stagedAdds),
		tombstones: orderedset.NewFrom(other.tombstones),
	}
}

// mapTo returns the set and the mapped index of the specified index.
func (this *DeltaSet[K]) mapTo(idx int) (*orderedset.OrderedSet[K], int) {
	if idx >= int(this.Length()) {
		return nil, -1
	}

	// The index is in the stagedAdds  list
	if idx >= this.committed.Length() {
		return this.stagedAdds, idx - this.committed.Length()
	}
	return this.committed, idx
}

// Reindex calculates the index of the element in the DeltaSet by skipping the tombstones  elements.
// func (this *DeltaSet[K]) Reindex(idx int) (*orderedset.OrderedSet[K], int) {
// 	// set, idx := this.mapTo(idx)

// }

// IsEmpty returns true if the DeltaSet is empty, i.e. no committed, stagedAdds , or tombstones  elements.
func (this *DeltaSet[K]) IsEmpty() bool {
	return this.committed.Length() == 0 && this.stagedAdds.Length() == 0 && this.tombstones.Length() == 0
}

func (this *DeltaSet[K]) Committed() *orderedset.OrderedSet[K] { return this.committed }
func (this *DeltaSet[K]) Removed() *orderedset.OrderedSet[K]   { return this.tombstones }
func (this *DeltaSet[K]) Added() *orderedset.OrderedSet[K]     { return this.stagedAdds }

func (this *DeltaSet[K]) SizeRemoved() int { return this.tombstones.Length() }
func (this *DeltaSet[K]) SizeAdded() int   { return this.stagedAdds.Length() }

func (this *DeltaSet[K]) SetCommitted(v *orderedset.OrderedSet[K]) { this.committed = v }
func (this *DeltaSet[K]) SetRemoved(v *orderedset.OrderedSet[K])   { this.tombstones = v }
func (this *DeltaSet[K]) SetAdded(v *orderedset.OrderedSet[K])     { this.stagedAdds = v }

// Elements returns the underlying slice of committed in the DeltaSet,
// Non-nil values are returned in the order they were inserted, equals to committed + stagedAdds  - tombstones.
func (this *DeltaSet[K]) Elements() []K {
	elements := make([]K, 0, this.NonNilCount())
	for i := 0; i < this.committed.Length()+this.stagedAdds.Length(); i++ {
		if v, ok := this.GetByIndex(uint64(i)); ok {
			elements = append(elements, v)
		}
	}
	return elements
}

// Get Byte size of the DeltaSet
func (this *DeltaSet[K]) Size(getter func(K) int) int {
	return common.IfThenDo1st(this.committed != nil, func() int { return this.committed.Size(getter) }, 0) +
		common.IfThenDo1st(this.stagedAdds != nil, func() int { return this.stagedAdds.Size(getter) }, 0) +
		common.IfThenDo1st(this.tombstones != nil, func() int { return this.tombstones.Size(getter) }, 0)
}

func (this *DeltaSet[K]) Clear() {
	this.ResetDelta()
}

func (this *DeltaSet[K]) ClearFull() {
	this.ResetDelta()
	this.committed.Clear()
}

// Debugging only
func (this *DeltaSet[K]) InsertCommitted(v []K) { this.committed.InsertBatch(v) }
func (this *DeltaSet[K]) InsertRemoved(v []K)   { this.tombstones.InsertBatch(v) }
func (this *DeltaSet[K]) InsertAdded(v []K)     { this.stagedAdds.InsertBatch(v) }

func (this *DeltaSet[K]) GetNilVal() K  { return this.nilVal }
func (this *DeltaSet[K]) SetNilVal(v K) { this.nilVal = v }

// IsDirty returns true if the DeltaSet is up to date,
// having no stagedAdds  or tombstones  elements.
func (this *DeltaSet[K]) IsDirty() bool {
	return this.tombstones.Length() != 0 || this.stagedAdds.Length() != 0
}

// Delta returns a new instance of DeltaSet with the same stagedAdds  and tombstones  elements only.
func (this *DeltaSet[K]) Delta() *DeltaSet[K] {
	return this.CloneDelta()
}

// SetDelta sets the stagedAdds  and tombstones  lists to the specified DeltaSet.
func (this *DeltaSet[K]) SetDelta(delta *DeltaSet[K]) {
	this.stagedAdds = delta.stagedAdds
	this.tombstones = delta.tombstones
}

// ResetDelta resets the stagedAdds  and tombstones  lists to empty.
func (this *DeltaSet[K]) ResetDelta() {
	this.stagedAdds.Clear()
	this.tombstones.Clear()
}

// Length returns the number of elements in the DeltaSet, including the NIL values.
func (this *DeltaSet[K]) Length() uint64 {
	return uint64(this.committed.Length() + this.stagedAdds.Length())
}

// NonNilCount returns the number of NON-NIL elements in the DeltaSet.
func (this *DeltaSet[K]) NonNilCount() uint64 {
	return uint64(this.committed.Length() + this.stagedAdds.Length() - this.tombstones.Length())
}

// Insert inserts an element into the DeltaSet and updates the index.
func (this *DeltaSet[K]) InsertBatch(elems []K) *DeltaSet[K] {
	for _, elem := range elems {
		this.Insert(elem)
	}
	return this
}

func (this *DeltaSet[K]) Insert(elem K) any {
	if ok, _ := this.tombstones.Exists(elem); ok {
		this.tombstones.Delete(elem) // Remove an existing element will take effect only when commit() is called.
		return this
	}

	if ok, _ := this.committed.Exists(elem); !ok {
		this.stagedAdds.Insert(elem) // Either in the stagedAdds  list or not.
	}
	return this
}

// Insert an element into the Delta Set and updates the index.
func (this *DeltaSet[K]) DeleteBatch(elems []K) *DeltaSet[K] {
	for _, elem := range elems {
		this.Delete(elem)
	}
	return this
}

// Insert inserts an element into the Del taSet and updates the index.
func (this *DeltaSet[K]) Delete(elem K) any {
	if ok, _ := this.committed.Exists(elem); ok {
		this.tombstones.Insert(elem) // Remove an existing element will take effect only when commit() is called.
	} else if ok, _ := this.stagedAdds.Exists(elem); ok {
		this.tombstones.Insert(elem)
	}

	return this
}

// Clone returns a new instance of DeltaSet with the same elements.
func (this *DeltaSet[K]) CloneFull() *DeltaSet[K] {
	set := this.CloneDelta()
	set.committed = this.committed.Clone()
	return set
}

// Clone returns a new instance with the shared shared committed set.
func (this *DeltaSet[K]) Clone() *DeltaSet[K] {
	set := this.CloneDelta()
	set.committed = this.committed //.Clone()
	return set
}

// CloneDelta returns a new instance of DeltaSet with the
// same stagedAdds  and tombstones  elements only.the committed list is not cloned.
func (this *DeltaSet[K]) CloneDelta() *DeltaSet[K] {
	set := &DeltaSet[K]{
		nilVal: this.nilVal,
		// committed: orderedset.NewOrderedSet(this.nilVal, 0),
		stagedAdds: this.stagedAdds.Clone(),
		tombstones: this.tombstones.Clone(),
	}
	return set
}

// CloneDelta returns a new instance of DeltaSet with the
// same stagedAdds  and tombstones  elements only.the committed list is not cloned.
// func (this *DeltaSet[K]) CloneDelta() *DeltaSet[K] {
// 	set := &DeltaSet[K]{
// 		nilVal:    this.nilVal,
// 		committed: orderedset.NewOrderedSet(this.nilVal, 0),
// 		stagedAdds :   this.stagedAdds .Clone(),
// 		tombstones :   this.tombstones.Clone(),
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
		if ok, _ := this.tombstones.Exists(k); !ok {
			return k, true
		}
	}
	return *new(K), false
}

// NthNonNil returns the nth non-nil value from the DeltaSet.
// The nth non-nil value isn't necessarily the nth value in the DeltaSet, but the nth non-nil value.
// func (this *DeltaSet[K]) GetNthNonNil(nth uint64) (K, int, bool) {
// 	// If the nth value is out of range, no need to search. The nil value is returned.
// 	if nth >= this.NonNilCount() {
// 		return *new(K), -1, false
// 	}

// 	// This isn't efficient; it is better to search from the beginning or the min and max indices of the tombstones  list to
// 	// narrow down the search range. However, that requires some extra code to keep track of the corresponding indices of
// 	// the keys in the tombstones  list.
// 	cnt := 0
// 	for i := 0; i < int(this.Length()); i++ {
// 		if K, ok := this.GetByIndex(uint64(i)); ok {
// 			if uint64(cnt) == nth {
// 				return K, i, true
// 			}
// 			cnt++
// 		}
// 	}
// 	return *new(K), -1, false
// }

// NthNonNil returns the nth non-nil value from the DeltaSet.
// The nth non-nil value isn't necessarily the nth value in the DeltaSet, but the nth non-nil value.
func (this *DeltaSet[K]) GetNthNonNil(nth uint64) (K, int, bool) {
	// If the nth value is out of range, no need to search. The nil value is returned.
	if nth >= this.NonNilCount() {
		return *new(K), -1, false
	}

	start := 0
	tombstonesElems := this.tombstones.Elements()
	// if this.Length() > uint64(len(tombstonesElems)) {
	for _, k := range tombstonesElems {
		if idx := this.IdxOf(k); idx <= uint64(nth) {
			start++
		}
	}
	// }

	// This isn't efficient; it is better to search from the beginning or the min and max indices of the removed list to
	// narrow down the search range. However, that requires some extra code to keep track of the corresponding indices of
	// the keys in the removed list.
	for i, cnt := start, 0; i < int(this.Length()); i++ {
		if K, ok := this.GetByIndex(uint64(i)); ok {
			if uint64(cnt) == nth+uint64(start) {
				return K, i, true
			}
			cnt++
		}
	}

	return *new(K), -1, false
}

// Back get the last non-nil value from the DeltaSet.
// The last non-nil value isn't necessarily the last value in the DeltaSet, but the last non-nil value.
func (this *DeltaSet[K]) Last() (K, bool) {
	if this.NonNilCount() == 0 {
		return *new(K), false
	}

	for i := int(this.Length() - 1); i >= 0; i-- {
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
		if ok, _ := this.tombstones.Exists(k); ok { // In the tombstones  set
			return this.nilVal, false
		}
		return k, true
	}
	return k, false
}

// Get returns the element at the specified index. If the index is out of range, the nil value is returned.
func (this *DeltaSet[K]) Exists(k K) (bool, int) {
	if ok, _ := this.tombstones.Exists(k); ok {
		return false, -1 // Has been tombstones  already
	}

	if ok, v := this.committed.Exists(k); ok {
		return ok, v
	}

	if ok, v := this.stagedAdds.Exists(k); ok {
		return ok, v
	}

	return false, -1
}

func (this *DeltaSet[K]) Commit(other []*DeltaSet[K]) *DeltaSet[K] {
	// Merge stagedAdds  elements from all the delta sets
	stagedAddsBuffer := make([]K, slice.CountDo(other, func(_ int, v **DeltaSet[K]) int { return (*v).stagedAdds.Length() })+this.stagedAdds.Length())
	slice.ConcateToBuffer(other, &stagedAddsBuffer, func(v *DeltaSet[K]) []K { return v.stagedAdds.Elements() })
	copy(stagedAddsBuffer[len(stagedAddsBuffer)-this.stagedAdds.Length():], this.stagedAdds.Elements())

	// Merge deleted elements from all the delta sets
	tombstonesBuffer := make([]K, slice.CountDo(other, func(_ int, v **DeltaSet[K]) int { return (*v).tombstones.Length() })+this.tombstones.Length())
	slice.ConcateToBuffer(other, &tombstonesBuffer, func(v *DeltaSet[K]) []K { return v.tombstones.Elements() })
	copy(tombstonesBuffer[len(tombstonesBuffer)-this.tombstones.Length():], this.tombstones.Elements())

	this.committed.Merge(stagedAddsBuffer) // Merge the stagedAdds  list to the committed list
	this.committed.Sub(tombstonesBuffer)   // Remove the tombstones  list from the committed list
	this.ResetDelta()                      // Reset the stagedAdds  and tombstones  lists to empty	// return this

	return this
}

func (this *DeltaSet[K]) Equal(other *DeltaSet[K]) bool {
	return this.committed.Equal(other.committed) &&
		this.stagedAdds.Equal(other.stagedAdds) &&
		this.tombstones.Equal(other.tombstones)
}

func (this *DeltaSet[K]) Print() {
	this.committed.Print()
	this.stagedAdds.Print()
	this.tombstones.Print()
}
