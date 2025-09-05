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
	"fmt"
	"math"

	orderedset "github.com/arcology-network/common-lib/exp/orderedset"
	"github.com/arcology-network/common-lib/exp/slice"
)

// DeltaSet represents a mutable view over a base set, allowing staged additions and deletions.
type DeltaSet[K comparable] struct {
	committed       *orderedset.OrderedSet[K] // The stable, committed elements
	stagedAdditions *orderedset.OrderedSet[K] // New elements staged for addition
	stagedRemovals  *StagedRemovalSet[K]
}

// NewIndexedSlice creates a new instance of DeltaSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewDeltaSet[K comparable](nilVal K, preAlloc int,
	size func(K) int,
	encodeTo func(K, []byte) int,
	decoder func([]byte) K,
	hasher func(K) [32]byte,
	keys ...K) *DeltaSet[K] {
	DeltaSet := &DeltaSet[K]{
		committed:       orderedset.NewOrderedSet(nilVal, preAlloc, size, encodeTo, decoder, hasher),
		stagedAdditions: orderedset.NewOrderedSet(nilVal, preAlloc, size, encodeTo, decoder, hasher),
		stagedRemovals:  NewStagedRemovalSet(nilVal, preAlloc, size, encodeTo, decoder, hasher, keys...),
	}
	DeltaSet.InsertBatch(keys)
	return DeltaSet
}

func (*DeltaSet[K]) New(
	// nilVal K,
	committed *orderedset.OrderedSet[K],
	stagedAdditions *orderedset.OrderedSet[K],
	stagedRemovals *StagedRemovalSet[K]) *DeltaSet[K] {

	return &DeltaSet[K]{
		// nilVal:          nilVal,
		committed:       committed,
		stagedAdditions: stagedAdditions,
		stagedRemovals:  stagedRemovals,
	}
}

func (*DeltaSet[K]) NewFrom(other *DeltaSet[K]) *DeltaSet[K] {
	return &DeltaSet[K]{
		committed:       orderedset.NewFrom(other.committed),
		stagedAdditions: orderedset.NewFrom(other.stagedAdditions),
		stagedRemovals:  new(StagedRemovalSet[K]).NewFrom(other.stagedRemovals),
	}
}

// mapTo returns the set and the mapped index of the specified index.
func (this *DeltaSet[K]) mapTo(idx int) (*orderedset.OrderedSet[K], int) {
	if idx >= int(this.Length()) {
		return nil, -1
	}

	// The index is in the stagedAdditions  list
	if idx >= this.committed.Length() {
		return this.stagedAdditions, idx - this.committed.Length()
	}
	return this.committed, idx
}

// Reindex calculates the index of the element in the DeltaSet by skipping the stagedRemovals  elements.
// func (this *DeltaSet[K]) Reindex(idx int) (*orderedset.OrderedSet[K], int) {
// 	// set, idx := this.mapTo(idx)

// }

// IsEmpty returns true if the DeltaSet is empty, i.e. no committed, stagedAdditions , or stagedRemovals  elements.
func (this *DeltaSet[K]) IsEmpty() bool {
	return this.committed.Length() == 0 && this.stagedAdditions.Length() == 0 && this.stagedRemovals.Length() == 0
}

func (this *DeltaSet[K]) Clear() { this.ResetDelta() }

func (this *DeltaSet[K]) Committed() *orderedset.OrderedSet[K] { return this.committed }
func (this *DeltaSet[K]) Removed() *StagedRemovalSet[K]        { return this.stagedRemovals }
func (this *DeltaSet[K]) Added() *orderedset.OrderedSet[K]     { return this.stagedAdditions }

func (this *DeltaSet[K]) SizeRemoved() int { return int(this.stagedRemovals.Length()) }
func (this *DeltaSet[K]) SizeAdded() int   { return this.stagedAdditions.Length() }

func (this *DeltaSet[K]) SetCommitted(v *orderedset.OrderedSet[K]) { this.committed = v }
func (this *DeltaSet[K]) SetRemoved(v *StagedRemovalSet[K])        { this.stagedRemovals = v }
func (this *DeltaSet[K]) SetAdded(v *orderedset.OrderedSet[K]) {
	this.stagedAdditions = v
}

// Elements returns the underlying slice of committed in the DeltaSet,
// Non-nil values are returned in the order they were inserted, equals to committed + stagedAdditions  - stagedRemovals.
func (this *DeltaSet[K]) Elements() []K {
	elements := make([]K, 0, this.NonNilCount())
	for i := 0; i < this.committed.Length()+this.stagedAdditions.Length(); i++ {
		if v, ok := this.GetByIndex(uint64(i)); ok {
			elements = append(elements, *v)
		}
	}
	return elements
}

// Debugging only
func (this *DeltaSet[K]) InsertCommitted(v []K) { this.committed.InsertBatch(v) }
func (this *DeltaSet[K]) InsertRemoved(v []K)   { this.stagedRemovals.InsertBatch(v) }
func (this *DeltaSet[K]) InsertAdded(v []K)     { this.stagedAdditions.InsertBatch(v) }

// func (this *DeltaSet[K]) GetNilVal() K  { return nil }
// func (this *DeltaSet[K]) SetNilVal(v K) { this.nilVal = v }

// IsDirty returns true if the DeltaSet is up to date,
// having no stagedAdditions  or stagedRemovals  elements.
func (this *DeltaSet[K]) IsDirty() bool {
	return this.stagedRemovals.Length() != 0 || this.stagedAdditions.Length() != 0
}

// Only delete elements in the committed set. No additions.
func (this *DeltaSet[K]) CommittedOnly() bool {
	return this.stagedAdditions.Length() == 0
}

// Delta returns a new instance of DeltaSet with the same stagedAdditions  and stagedRemovals elements only.
func (this *DeltaSet[K]) Delta() *DeltaSet[K] {
	return this.CloneDelta()
}

// SetDelta sets the stagedAdditions  and stagedRemovals  lists to the specified DeltaSet.
func (this *DeltaSet[K]) SetDelta(delta *DeltaSet[K]) {
	this.stagedAdditions = delta.stagedAdditions
	this.stagedRemovals = delta.stagedRemovals
}

// ResetDelta resets the stagedAdditions  and stagedRemovals  lists to empty.
func (this *DeltaSet[K]) ResetDelta() {
	this.stagedAdditions.Clear()
	this.stagedRemovals.Clear()
}

// Length returns the number of elements in the DeltaSet, including the NIL values.
func (this *DeltaSet[K]) Length() uint64 {
	return uint64(this.committed.Length() + this.stagedAdditions.Length())
}

// NonNilCount returns the number of NON-NIL elements in the DeltaSet.
func (this *DeltaSet[K]) NonNilCount() uint64 {
	return uint64(this.committed.Length() + this.stagedAdditions.Length() - int(this.stagedRemovals.Length()))
}

// Insert inserts an element into the DeltaSet and updates the index.
func (this *DeltaSet[K]) InsertBatch(elems []K) *DeltaSet[K] {
	for _, elem := range elems {
		this.Insert(elem)
	}
	return this
}

func (this *DeltaSet[K]) Insert(elem K) {
	if ok, _ := this.stagedRemovals.Exists(elem); ok {
		this.stagedRemovals.Delete(elem) // Remove an existing element will take effect only when commit() is called.
		return
	}

	if ok, _ := this.committed.Exists(elem); !ok {
		this.stagedAdditions.Insert(elem) // Either in the stagedAdditions  list or not.
	}
}

// Insert an element into the Delta Set and updates the index.
func (this *DeltaSet[K]) DeleteBatch(elems []K) *DeltaSet[K] {
	for _, elem := range elems {
		this.Delete(elem)
	}
	return this
}

// Insert inserts an element into the Delta Set and updates the index.
func (this *DeltaSet[K]) Delete(elem K) {
	if ok, _ := this.committed.Exists(elem); ok {
		// Remove an existing element will take effect only when commit() is called.
		this.stagedRemovals.Insert(elem)
	} else if ok, _ := this.stagedAdditions.Exists(elem); ok {
		this.stagedRemovals.Insert(elem)
	}
}

// Add both committed and the added elements to the deletion list.
// The return value indicates if some pending entries were involved.
func (this *DeltaSet[K]) DeleteCommitted() { this.stagedRemovals.SetCommitted(this.Committed()) }
func (this *DeltaSet[K]) DeleteAdded()     { this.stagedRemovals.SetAdded(this.stagedAdditions.Clone()) }

// Add both committed and the added elements to the deletion list.
func (this *DeltaSet[K]) DeleteAll() {
	// Shallow copy is enough, because the Committed elements won't change.
	this.stagedRemovals.SetCommitted(this.Committed())
	this.stagedRemovals.CommittedDeleted = true

	// No further changes to the stagedAdditions won't affect the stagedRemovals
	// from this point on.
	this.stagedRemovals.SetAdded(this.stagedAdditions.Clone())
	this.stagedRemovals.StagedAddedDeleted = true
}

func (this *DeltaSet[K]) DeleteByIndex(idx uint64) {
	if k, _, _, ok := this.Search(idx); ok {
		this.Delete(*k)
	}
}

// Clone returns a new instance of DeltaSet with the same elements.
// under no circumstances the committed should be deeply copied.
func (this *DeltaSet[K]) CloneFull() *DeltaSet[K] {
	return this.Clone()
}

// Clone returns a new instance with the
// committed list shared original DeltaSet.
func (this *DeltaSet[K]) Clone() *DeltaSet[K] {
	set := this.CloneDelta()
	set.committed = this.committed //Share the committed set.
	return set
}

// CloneDelta returns a new instance of DeltaSet with the
// same stagedAdditions  and stagedRemovals  elements only.the committed list is not cloned.
func (this *DeltaSet[K]) CloneDelta() *DeltaSet[K] {
	return &DeltaSet[K]{
		committed: orderedset.NewOrderedSet(*new(K), 0,
			this.committed.Sizer,
			this.committed.Encoder,
			this.committed.Decoder,
			nil),
		stagedAdditions: this.stagedAdditions.Clone(),
		stagedRemovals:  this.stagedRemovals.Clone(),
	}
}

// GetByIndex returns the element at the specified index.
// It does NOT check if the index is corresponding to a non-nil value.
// If the value at the index is nil, the nil value is returned.
func (this *DeltaSet[K]) GetByIndex(idx uint64) (*K, bool) {
	if k, _, _, ok := this.Search(idx); ok {
		if ok, _ := this.stagedRemovals.Exists(*k); !ok {
			return k, true
		}
	}
	return nil, false
}

// NthNonNil returns the nth non-nil value from the DeltaSet.
// The nth non-nil value isn't necessarily the nth value in the DeltaSet, but the nth non-nil value.
// func (this *DeltaSet[K]) GetNthNonNil(nth uint64) ([K], int, bool) {
// 	// If the nth value is out of range, no need to search. The nil value is returned.
// 	if nth >= this.NonNilCount() {
// 		return *new([K]), -1, false
// 	}

// 	// This isn't efficient; it is better to search from the beginning or the min and max indices of the stagedRemovals  list to
// 	// narrow down the search range. However, that requires some extra code to keep track of the corresponding indices of
// 	// the keys in the stagedRemovals  list.
// 	cnt := 0
// 	for i := 0; i < int(this.Length()); i++ {
// 		if [K], ok := this.GetByIndex(uint64(i)); ok {
// 			if uint64(cnt) == nth {
// 				return [K], i, true
// 			}
// 			cnt++
// 		}
// 	}
// 	return *new([K]), -1, false
// }

// NthNonNil returns the nth non-nil value from the DeltaSet.
// The nth non-nil value isn't necessarily the nth value in the DeltaSet, but the nth non-nil value.
func (this *DeltaSet[K]) GetNthNonNil(nth uint64) (*K, int, bool) {
	// If the nth value is out of range, no need to search. The nil value is returned.
	if nth >= this.NonNilCount() {
		return nil, -1, false
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

	return nil, -1, false
}

// Back get the last non-nil value from the DeltaSet.
// The last non-nil value isn't necessarily the last value in the DeltaSet, but the last non-nil value.
func (this *DeltaSet[K]) Last() (*K, bool) {
	if this.NonNilCount() == 0 {
		return nil, false
	}

	for i := int(this.Length() - 1); i >= 0; i-- {
		if k, ok := this.GetByIndex(uint64(i)); ok {
			return k, true
		}
	}
	return nil, false
}

func (this *DeltaSet[K]) Back() (*K, bool) {
	if this.NonNilCount() == 0 {
		return nil, false
	}
	return this.GetByIndex(uint64(this.Length() - 1))
}

// Return then remove the last non-nil value from the DeltaSet.
func (this *DeltaSet[K]) PopLast() (*K, bool) {
	// Get the last non-nil value
	if k, ok := this.Last(); ok {
		this.Delete(*k) // Remove the last non-nil value
		return k, true
	}
	return nil, false
}

// Search returns the element at the specified index and the set it is in.
func (this *DeltaSet[K]) Search(idx uint64) (*K, *orderedset.OrderedSet[K], int, bool) {
	set, mapped := this.mapTo(int(idx))
	if mapped < 0 {
		return nil, nil, -1, false
	}
	return set.IndexToKey(mapped), set, mapped, true
}

func (this *DeltaSet[K]) KeyAt(idx uint64) (*K, bool) {
	k, _, _, ok := this.Search(idx)
	return k, ok
}

// Get the index of the element in the DeltaSet by key.
func (this *DeltaSet[K]) IdxOf(k K) uint64 {
	if ok, idx := this.Exists(k); ok {
		return uint64(idx)
	}
	return math.MaxUint64
}

// Get returns the element at the specified index.
func (this *DeltaSet[K]) TryGetKey(idx uint64) (*K, bool) {
	if k, _, _, ok := this.Search(idx); ok {
		if ok, _ := this.stagedRemovals.Exists(*k); !ok { // In the stagedRemovals  set
			return k, true
		}
	}
	return nil, false
}

// Get returns the element at the specified index. If the index is out of range, the nil value is returned.
func (this *DeltaSet[K]) Exists(k K) (bool, int) {
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

func (this *DeltaSet[K]) Commit(other []*DeltaSet[K]) *DeltaSet[K] {
	// Merge stagedAdditions  elements from all the delta sets
	stagedAdditionsBuffer := make([]K, slice.CountDo(other, func(_ int, v **DeltaSet[K]) int { return (*v).stagedAdditions.Length() })+this.stagedAdditions.Length())
	slice.ConcateToBuffer(other, &stagedAdditionsBuffer, func(v *DeltaSet[K]) []K { return v.stagedAdditions.Elements() })
	copy(stagedAdditionsBuffer[len(stagedAdditionsBuffer)-this.stagedAdditions.Length():], this.stagedAdditions.Elements())

	// Merge deleted elements from all the delta sets
	stagedRemovalsBuffer := make([]K, slice.CountDo(other, func(_ int, v **DeltaSet[K]) int {
		return int((*v).stagedRemovals.Length())
	})+int(this.stagedRemovals.Length()))

	slice.ConcateToBuffer(other, &stagedRemovalsBuffer, func(v *DeltaSet[K]) []K { return v.stagedRemovals.Elements() })
	copy(stagedRemovalsBuffer[len(stagedRemovalsBuffer)-int(this.stagedRemovals.Length()):], this.stagedRemovals.Elements())

	this.committed.Merge(stagedAdditionsBuffer) // Merge the stagedAdditions  list to the committed list
	this.committed.Sub(stagedRemovalsBuffer)    // Remove the stagedRemovals  list from the committed list
	this.ResetDelta()                           // Reset the stagedAdditions  and stagedRemovals  lists to empty	// return this

	return this
}

func (this *DeltaSet[K]) Equal(other *DeltaSet[K]) bool {
	return this.committed.Equal(other.committed) &&
		this.stagedAdditions.Equal(other.stagedAdditions) &&
		this.stagedRemovals.Equal(other.stagedRemovals)
}

func (this *DeltaSet[K]) Print() {
	fmt.Print("Committed: ")
	this.committed.Print()

	fmt.Print("Staged Added: ")
	this.stagedAdditions.Print()

	fmt.Print("Staged Removed: ")
	this.stagedRemovals.Print()
}
