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
	"github.com/arcology-network/common-lib/exp/slice"
)

// DeltaSlice represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type DeltaSlice[T any] struct {
	elements []*T
	updated  []*T
	appended []*T
	removed  []int
}

// NewIndexedSlice creates a new instance of DeltaSlice with the specified page size, minimum number of pages, and pre-allocation size.
func NewDeltaSlice[T any](size int) *DeltaSlice[T] {
	return &DeltaSlice[T]{
		elements: make([]*T, 0, size),
		updated:  []*T{},
		appended: []*T{},
		removed:  []int{},
	}
}

// mapTo returns the slice and the mapped index of the specified container.
// Some value may have been removed, so there are some "holes" in the index.
// The original indices need to be mapped to the new indices.
func (this *DeltaSlice[T]) mapTo(idx int) (*[]*T, int) {
	if idx >= this.Length() {
		return nil, -1
	}

	// The index is in the appended list
	if len(this.elements)-len(this.removed) <= idx {
		return &this.appended, idx - (len(this.elements) - len(this.removed))
	}

	// The first step is to check if the index is in the removed list. There are two ways:
	// 1. Use a map to store the removed index, and check if the index is in the map.
	// 2. Use a SORTED slice for binary search.
	// Either way, there has to be extra steps to maintain the removed list.
	// So it the removed list is not too long, it is better to use a simple linear search.
	totalOffset := 0
	for i := 0; i < len(this.removed); i++ {
		if idx >= this.removed[i] {
			totalOffset++
		}
	}
	return &this.elements, idx + totalOffset
}

// Array returns the underlying slice of elements in the DeltaSlice.
func (this *DeltaSlice[T]) Values() []*T   { return this.elements }
func (this *DeltaSlice[T]) Updated() []*T  { return this.updated }
func (this *DeltaSlice[T]) Appended() []*T { return this.appended }
func (this *DeltaSlice[T]) Removed() []int { return this.removed }

// Insert inserts an element into the DeltaSlice and updates the index.
func (this *DeltaSlice[T]) Append(elem T) int {
	this.appended = append(this.appended, &elem)
	return this.Length() - 1
}

func (this *DeltaSlice[T]) Length() int {
	return len(this.elements) + len(this.appended) - len(this.removed)
}

// Insert inserts an element into the DeltaSlice and updates the index.
func (this *DeltaSlice[T]) ToSlice() []*T {
	if len(this.elements)-len(this.removed) == 0 { // No elements in the original list
		return this.appended
	}

	elements := make([]*T, len(this.elements)-len(this.removed))
	for i, idx := 0, 0; i < len(this.elements); i++ {
		if loc, _ := slice.FindFirst(this.removed, i); loc < 0 {
			elements[idx] = this.elements[i]
			idx++
		}
	}
	return append(elements, this.appended...)
}

// Insert inserts an element into the DeltaSlice and updates the index.
func (this *DeltaSlice[T]) Del(idx int) (bool, int) {
	from := 0
	if arr, mapped := this.mapTo(idx); arr != nil {
		if arr == &this.appended {
			slice.RemoveAt(arr, mapped) // Remove the element from the appended list
			from = 2
		} else {
			this.removed = append(this.removed, mapped)
			from = 1 // The element is in the original list
		}
		return true, from
	}
	return false, from
}

func (this *DeltaSlice[T]) Set(idx int, v T) (*T, bool) {
	if arr, mapped := this.mapTo(idx); arr != nil {
		if arr == &this.appended {
			(*arr)[mapped] = &v // Directly update the element if it is in the appended list
		} else {
			// Add to the updated list if it is in the original list, replacement is expensive without an index.
			// Because all the updates are in a chronological order, the last updates will be guaranteed to be
			// the final values.
			this.updated = append(this.updated, &v)
		}
		return &v, true
	}
	return nil, false
}

func (this *DeltaSlice[T]) Get(idx int) (*T, bool) {
	if arr, mapped := this.mapTo(idx); arr != nil {
		return (*arr)[mapped], true
	}
	return nil, true
}

func (this *DeltaSlice[T]) Commit() int {
	for i := 0; i < len(this.removed); i++ {
		this.elements[this.removed[i]] = nil
	}
	this.elements = append(slice.Remove(&this.elements, nil), this.appended...)
	this.removed = this.removed[:0]
	this.appended = this.appended[:0]
	return len(this.elements)
}
