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
	elements []T
	appended []T
	removed  []int
}

// NewIndexedSlice creates a new instance of DeltaSlice with the specified page size, minimum number of pages, and pre-allocation size.
func NewDeltaSlice[T any](size int) *DeltaSlice[T] {
	return &DeltaSlice[T]{
		elements: make([]T, 0, size),
		appended: []T{},
		removed:  []int{},
	}
}

// mapTo returns the slice and the mapped index of the specified container.
func (this *DeltaSlice[T]) mapTo(idx int) (*[]T, int) {
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
	// Either way, there has to be extra steps to maintain the removed list. So it the removed list is not too long, it is better to use a simple linear search.
	totalOffset := 0
	for i := 0; i < len(this.removed); i++ {
		if this.removed[i] <= idx {
			totalOffset++
		}
	}
	return &this.elements, idx + totalOffset
}

// Array returns the underlying slice of elements in the DeltaSlice.
func (this *DeltaSlice[T]) Elements() []T  { return this.elements }
func (this *DeltaSlice[T]) Appended() []T  { return this.appended }
func (this *DeltaSlice[T]) Removed() []int { return this.removed }

// Insert inserts an element into the DeltaSlice and updates the index.
func (this *DeltaSlice[T]) Append(elem T) { this.appended = append(this.appended, elem) }
func (this *DeltaSlice[T]) Length() int {
	return len(this.elements) + len(this.appended) - len(this.removed)
}

// Insert inserts an element into the DeltaSlice and updates the index.
func (this *DeltaSlice[T]) ToSlice() []T {
	if len(this.elements)-len(this.removed) == 0 { // No elements in the original list
		return this.appended
	}

	elements := make([]T, len(this.elements)-len(this.removed))
	for i, idx := 0, 0; i < len(this.elements); i++ {
		if loc, _ := slice.FindFirst(this.removed, i); loc < 0 {
			elements[idx] = this.elements[i]
			idx++
		}
	}
	return append(elements, this.appended...)
}

// Insert inserts an element into the DeltaSlice and up^dates the index.
func (this *DeltaSlice[T]) Del(idx int) bool {
	if arr, mapped := this.mapTo(idx); arr != nil {
		if arr == &this.appended {
			slice.RemoveAt(arr, mapped)
		} else {
			this.removed = append(this.removed, mapped)
		}
		return true
	}
	return false
}

func (this *DeltaSlice[T]) Set(idx int, v T) bool {
	if arr, mapped := this.mapTo(idx); arr != nil {
		(*arr)[mapped] = v
		return true
	}
	return false
}

func (this *DeltaSlice[T]) Get(idx int) (*T, bool) {
	if arr, mapped := this.mapTo(idx); arr != nil {
		return &(*arr)[mapped], true
	}
	return nil, true
}
