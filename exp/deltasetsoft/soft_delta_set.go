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

package softdeltaset

import (
	"math"

	"github.com/arcology-network/common-lib/common"
	orderedset "github.com/arcology-network/common-lib/exp/orderedset"
	"github.com/arcology-network/common-lib/exp/slice"
)

// SoftDeltaSet represents a mutable view over a base set, allowing staged additions and deletions.
type SoftDeltaSet[K comparable] struct {
	nilVal          K
	committed       *orderedset.OrderedSet[K] // The stable, committed elements
	stagedAdditions *orderedset.OrderedSet[K] // New elements staged for addition
	stagedRemovals  *StagedRemovalset[K]
	allDeleted      bool // If true, all elements are deleted, i.e. the set is empty.
}

// NewIndexedSlice creates a new instance of SoftDeltaSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewSoftDeltaSet[K comparable](nilVal K, preAlloc int, hasher func(K) [32]byte, keys ...K) *SoftDeltaSet[K] {
	SoftDeltaSet := &SoftDeltaSet[K]{
		nilVal:          nilVal,
		committed:       orderedset.NewOrderedSet(nilVal, preAlloc, hasher),
		stagedAdditions: orderedset.NewOrderedSet(nilVal, preAlloc, hasher),
		stagedRemovals:  NewStagedRemovalset(nilVal, preAlloc, hasher, keys...),
	}
	SoftDeltaSet.InsertBatch(keys)
	return SoftDeltaSet
}

func (*SoftDeltaSet[K]) New(
	nilVal K,
	committed *orderedset.OrderedSet[K],
	stagedAdditions *orderedset.OrderedSet[K],
	stagedRemovals *StagedRemovalset[K]) *SoftDeltaSet[K] {

	return &SoftDeltaSet[K]{
		nilVal:          nilVal,
		committed:       committed,
		stagedAdditions: stagedAdditions,
		stagedRemovals:  stagedRemovals,
	}
}

func (*SoftDeltaSet[K]) NewFrom(other *SoftDeltaSet[K]) *SoftDeltaSet[K] {
	return &SoftDeltaSet[K]{
		nilVal:          other.GetNilVal(),
		committed:       orderedset.NewFrom(other.committed),
		stagedAdditions: orderedset.NewFrom(other.stagedAdditions),
		stagedRemovals:  new(StagedRemovalset[K]).NewFrom(other.stagedRemovals),
	}
}

// mapTo returns the set and the mapped index of the specified index.
func (this *SoftDeltaSet[K]) mapTo(idx int) (*orderedset.OrderedSet[K], int) {
	if idx >= int(this.Length()) {
		return nil, -1
	}

	// The index is in the stagedAdditions  list
	if idx >= this.committed.Length() {
		return this.stagedAdditions, idx - this.committed.Length()
	}
	return this.committed, idx
}

// Reindex calculates the index of the element in the SoftDeltaSet by skipping the stagedRemovals  elements.
// func (this *SoftDeltaSet[K]) Reindex(idx int) (*orderedset.OrderedSet[K], int) {
// 	// set, idx := this.mapTo(idx)

// }

// IsEmpty returns true if the SoftDeltaSet is empty, i.e. no committed, stagedAdditions , or stagedRemovals  elements.
func (this *SoftDeltaSet[K]) IsEmpty() bool {
	return this.committed.Length() == 0 && this.stagedAdditions.Length() == 0 && this.stagedRemovals.Length() == 0
}

func (this *SoftDeltaSet[K]) Committed() *orderedset.OrderedSet[K] { return this.committed }
func (this *SoftDeltaSet[K]) Removed() *StagedRemovalset[K]        { return this.stagedRemovals }
func (this *SoftDeltaSet[K]) Added() *orderedset.OrderedSet[K]     { return this.stagedAdditions }

func (this *SoftDeltaSet[K]) SizeRemoved() int { return int(this.stagedRemovals.Length()) }
func (this *SoftDeltaSet[K]) SizeAdded() int   { return this.stagedAdditions.Length() }

func (this *SoftDeltaSet[K]) SetCommitted(v *orderedset.OrderedSet[K]) { this.committed = v }
func (this *SoftDeltaSet[K]) SetRemoved(v *StagedRemovalset[K])        { this.stagedRemovals = v }
func (this *SoftDeltaSet[K]) SetAdded(v *orderedset.OrderedSet[K])     { this.stagedAdditions = v }

// Elements returns the underlying slice of committed in the SoftDeltaSet,
// Non-nil values are returned in the order they were inserted, equals to committed + stagedAdditions  - stagedRemovals.
func (this *SoftDeltaSet[K]) Elements() []K {
	elements := make([]K, 0, this.NonNilCount())
	for i := 0; i < this.committed.Length()+this.stagedAdditions.Length(); i++ {
		if v, ok := this.GetByIndex(uint64(i)); ok {
			elements = append(elements, v)
		}
	}
	return elements
}

// Get Byte size of the SoftDeltaSet
func (this *SoftDeltaSet[K]) Size(getter func(K) int) int {
	return common.IfThenDo1st(this.committed != nil, func() int { return this.committed.Size(getter) }, 0) +
		common.IfThenDo1st(this.stagedAdditions != nil, func() int { return this.stagedAdditions.Size(getter) }, 0) +
		common.IfThenDo1st(this.stagedRemovals != nil, func() int { return this.stagedRemovals.Size(getter) }, 0)
}

func (this *SoftDeltaSet[K]) Clear() {
	this.ResetDelta()
}

// Debugging only
func (this *SoftDeltaSet[K]) InsertCommitted(v []K) { this.committed.InsertBatch(v) }
func (this *SoftDeltaSet[K]) InsertRemoved(v []K)   { this.stagedRemovals.InsertBatch(v) }
func (this *SoftDeltaSet[K]) InsertAdded(v []K)     { this.stagedAdditions.InsertBatch(v) }

func (this *SoftDeltaSet[K]) GetNilVal() K  { return this.nilVal }
func (this *SoftDeltaSet[K]) SetNilVal(v K) { this.nilVal = v }

// IsDirty returns true if the SoftDeltaSet is up to date,
// having no stagedAdditions  or stagedRemovals  elements.
func (this *SoftDeltaSet[K]) IsDirty() bool {
	return this.stagedRemovals.Length() != 0 || this.stagedAdditions.Length() != 0
}

// Delta returns a new instance of SoftDeltaSet with the same stagedAdditions  and stagedRemovals elements only.
func (this *SoftDeltaSet[K]) Delta() *SoftDeltaSet[K] {
	return this.CloneDelta()
}

// SetDelta sets the stagedAdditions  and stagedRemovals  lists to the specified SoftDeltaSet.
func (this *SoftDeltaSet[K]) SetDelta(delta *SoftDeltaSet[K]) {
	this.stagedAdditions = delta.stagedAdditions
	this.stagedRemovals = delta.stagedRemovals
}

// ResetDelta resets the stagedAdditions  and stagedRemovals  lists to empty.
func (this *SoftDeltaSet[K]) ResetDelta() {
	this.stagedAdditions.Clear()
	this.stagedRemovals.Clear()
}

// Length returns the number of elements in the SoftDeltaSet, including the NIL values.
func (this *SoftDeltaSet[K]) Length() uint64 {
	return uint64(this.committed.Length() + this.stagedAdditions.Length())
}

// NonNilCount returns the number of NON-NIL elements in the SoftDeltaSet.
func (this *SoftDeltaSet[K]) NonNilCount() uint64 {
	return uint64(this.committed.Length() + this.stagedAdditions.Length() - int(this.stagedRemovals.Length()))
}

// Insert inserts an element into the SoftDeltaSet and updates the index.
func (this *SoftDeltaSet[K]) InsertBatch(elems []K) *SoftDeltaSet[K] {
	for _, elem := range elems {
		this.Insert(elem)
	}
	return this
}

func (this *SoftDeltaSet[K]) Insert(elem K) {
	if ok, _ := this.stagedRemovals.Exists(elem); ok {
		this.stagedRemovals.Delete(elem) // Remove an existing element will take effect only when commit() is called.
		return
	}

	if ok, _ := this.committed.Exists(elem); !ok {
		this.stagedAdditions.Insert(elem) // Either in the stagedAdditions  list or not.
	}
}

// Insert an element into the Delta Set and updates the index.
func (this *SoftDeltaSet[K]) DeleteBatch(elems []K) *SoftDeltaSet[K] {
	for _, elem := range elems {
		this.Delete(elem)
	}
	return this
}

// Insert inserts an element into the Del taSet and updates the index.
func (this *SoftDeltaSet[K]) Delete(elem K) {
	if ok, _ := this.committed.Exists(elem); ok {
		this.stagedRemovals.Insert(elem) // Remove an existing element will take effect only when commit() is called.
	} else if ok, _ := this.stagedAdditions.Exists(elem); ok {
		this.stagedRemovals.Insert(elem)
	}
}

func (this *SoftDeltaSet[K]) DeleteAll() {
	this.stagedRemovals.SetCommitted(this.Committed())
	this.stagedRemovals.SetAdded(this.stagedAdditions)
	this.allDeleted = true
}

func (this *SoftDeltaSet[K]) DeleteByIndex(idx uint64) {
	if k, _, _, ok := this.Search(idx); ok {
		this.Delete(k)
	}
}

// Clone returns a new instance of SoftDeltaSet with the same elements.
func (this *SoftDeltaSet[K]) CloneFull() *SoftDeltaSet[K] {
	set := this.CloneDelta()
	set.committed = this.committed.Clone()
	return set
}

// Clone returns a new instance with the
// committed list shared original SoftDeltaSet.
func (this *SoftDeltaSet[K]) Clone() *SoftDeltaSet[K] {
	set := this.CloneDelta()
	set.committed = this.committed //.Clone()
	return set
}

// CloneDelta returns a new instance of SoftDeltaSet with the
// same stagedAdditions  and stagedRemovals  elements only.the committed list is not cloned.
func (this *SoftDeltaSet[K]) CloneDelta() *SoftDeltaSet[K] {
	return &SoftDeltaSet[K]{
		nilVal: this.nilVal,
		// committed: orderedset.NewOrderedSet(this.nilVal, 0),
		stagedAdditions: this.stagedAdditions.Clone(),
		stagedRemovals:  this.stagedRemovals.CloneFull(),
	}
}

// CloneDelta returns a new instance of SoftDeltaSet with the
// same stagedAdditions  and stagedRemovals  elements only.the committed list is not cloned.
// func (this *SoftDeltaSet[K]) CloneDelta() *SoftDeltaSet[K] {
// 	set := &SoftDeltaSet[K]{
// 		nilVal:    this.nilVal,
// 		committed: orderedset.NewOrderedSet(this.nilVal, 0),
// 		stagedAdditions :   this.stagedAdditions .Clone(),
// 		stagedRemovals :   this.stagedRemovals.Clone(),
// 	}
// 	return set
// }

// GetByIndex returns the element at the specified index.
// It does NOT check if the index is corresponding to a non-nil value.
// If the value at the index is nil, the nil value is returned.
func (this *SoftDeltaSet[K]) GetByIndex(idx uint64) (K, bool) {
	if k, _, _, ok := this.Search(idx); ok {
		if ok, _ := this.stagedRemovals.Exists(k); !ok {
			return k, true
		}
	}
	return *new(K), false
}

// NthNonNil returns the nth non-nil value from the SoftDeltaSet.
// The nth non-nil value isn't necessarily the nth value in the SoftDeltaSet, but the nth non-nil value.
// func (this *SoftDeltaSet[K]) GetNthNonNil(nth uint64) (K, int, bool) {
// 	// If the nth value is out of range, no need to search. The nil value is returned.
// 	if nth >= this.NonNilCount() {
// 		return *new(K), -1, false
// 	}

// 	// This isn't efficient; it is better to search from the beginning or the min and max indices of the stagedRemovals  list to
// 	// narrow down the search range. However, that requires some extra code to keep track of the corresponding indices of
// 	// the keys in the stagedRemovals  list.
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

// NthNonNil returns the nth non-nil value from the SoftDeltaSet.
// The nth non-nil value isn't necessarily the nth value in the SoftDeltaSet, but the nth non-nil value.
func (this *SoftDeltaSet[K]) GetNthNonNil(nth uint64) (K, int, bool) {
	// If the nth value is out of range, no need to search. The nil value is returned.
	if nth >= this.NonNilCount() {
		return *new(K), -1, false
	}

	// removedElems := this.removed.Elements()
	start := 0
	stagedRemovalsElems := this.stagedRemovals.Elements()
	// if this.Length() > uint64(len(stagedRemovals Elems)) {
	for _, k := range stagedRemovalsElems {
		if idx := this.IdxOf(k); idx <= uint64(nth) {
			start++
		}
	}
	// }

	// This isn't efficient; it is better to search from the beginning or the min and max indices of the stagedRemovals  list to
	// narrow down the search range. However, that requires some extra code to keep track of the corresponding indices of
	// the keys in the stagedRemovals  list.
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

// Back get the last non-nil value from the SoftDeltaSet.
// The last non-nil value isn't necessarily the last value in the SoftDeltaSet, but the last non-nil value.
func (this *SoftDeltaSet[K]) Last() (K, bool) {
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

func (this *SoftDeltaSet[K]) Back() (K, bool) {
	if this.NonNilCount() == 0 {
		return *new(K), false
	}
	return this.GetByIndex(uint64(this.Length() - 1))
}

// Return then remove the last non-nil value from the SoftDeltaSet.
func (this *SoftDeltaSet[K]) PopLast() (K, bool) {
	k, ok := this.Last() // Get the last non-nil value
	if ok {
		this.Delete(k) // Re
	}
	return k, ok
}

// Search returns the element at the specified index and the set it is in.
func (this *SoftDeltaSet[K]) Search(idx uint64) (K, *orderedset.OrderedSet[K], int, bool) {
	set, mapped := this.mapTo(int(idx))
	if mapped < 0 {
		return this.nilVal, nil, -1, false
	}
	return set.IndexToKey(mapped), set, mapped, true
}

func (this *SoftDeltaSet[K]) KeyAt(idx uint64) K {
	k, _, _, _ := this.Search(idx)
	return k
}

// Get the index of the element in the SoftDeltaSet by key.
func (this *SoftDeltaSet[K]) IdxOf(k K) uint64 {
	if ok, idx := this.Exists(k); ok {
		return uint64(idx)
	}
	return math.MaxUint64
}

// Get returns the element at the specified index.
func (this *SoftDeltaSet[K]) TryGetKey(idx uint64) (K, bool) {
	k, _, _, ok := this.Search(idx)
	if ok {
		if ok, _ := this.stagedRemovals.Exists(k); ok { // In the stagedRemovals  set
			return this.nilVal, false
		}
		return k, true
	}
	return k, false
}

// Get returns the element at the specified index. If the index is out of range, the nil value is returned.
func (this *SoftDeltaSet[K]) Exists(k K) (bool, int) {
	if ok, _ := this.stagedRemovals.Exists(k); ok {
		return false, -1 // Has been stagedRemovals  already
	}

	if ok, v := this.committed.Exists(k); ok {
		return ok, v
	}

	if ok, v := this.stagedAdditions.Exists(k); ok {
		return ok, v
	}

	return false, -1
}

func (this *SoftDeltaSet[K]) Commit(other []*SoftDeltaSet[K]) *SoftDeltaSet[K] {
	// Merge stagedAdditions  elements from all the delta sets
	stagedAdditionsBuffer := make([]K, slice.CountDo(other, func(_ int, v **SoftDeltaSet[K]) int { return (*v).stagedAdditions.Length() })+this.stagedAdditions.Length())
	slice.ConcateToBuffer(other, &stagedAdditionsBuffer, func(v *SoftDeltaSet[K]) []K { return v.stagedAdditions.Elements() })
	copy(stagedAdditionsBuffer[len(stagedAdditionsBuffer)-this.stagedAdditions.Length():], this.stagedAdditions.Elements())

	// Merge deleted elements from all the delta sets
	stagedRemovalsBuffer := make([]K, slice.CountDo(other, func(_ int, v **SoftDeltaSet[K]) int {
		return int((*v).stagedRemovals.Length())
	})+int(this.stagedRemovals.Length()))

	slice.ConcateToBuffer(other, &stagedRemovalsBuffer, func(v *SoftDeltaSet[K]) []K { return v.stagedRemovals.Elements() })
	copy(stagedRemovalsBuffer[len(stagedRemovalsBuffer)-int(this.stagedRemovals.Length()):], this.stagedRemovals.Elements())

	this.committed.Merge(stagedAdditionsBuffer) // Merge the stagedAdditions  list to the committed list
	this.committed.Sub(stagedRemovalsBuffer)    // Remove the stagedRemovals  list from the committed list
	this.ResetDelta()                           // Reset the stagedAdditions  and stagedRemovals  lists to empty	// return this

	return this
}

func (this *SoftDeltaSet[K]) Equal(other *SoftDeltaSet[K]) bool {
	return this.committed.Equal(other.committed) &&
		this.stagedAdditions.Equal(other.stagedAdditions) &&
		this.stagedRemovals.Equal(&other.stagedRemovals.DeltaSet)
}

func (this *SoftDeltaSet[K]) Print() {
	this.committed.Print()
	this.stagedAdditions.Print()
	this.stagedRemovals.Print()
}
