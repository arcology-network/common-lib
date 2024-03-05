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

package indexedslice

import (
	"github.com/arcology-network/common-lib/exp/associative"
)

// IndexedSlice represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type IndexedSlice[K comparable, T0 any, T1 any] struct {
	elements []associative.Pair[*K, T1]
	index    map[K]int

	getkey      func(T0) K
	initializer func(K, T0) T1
	updater     func(K, T0, *T1)
	IsEmpty     func(T1) bool
}

// NewIndexedSlice creates a new instance of IndexedSlice with the specified page size, minimum number of pages, and pre-allocation size.
func NewIndexedSlice[K comparable, T0 any, T1 any](
	getkey func(T0) K,
	initializer func(K, T0) T1,
	updater func(K, T0, *T1),
	isEmpty func(T1) bool,
	preAlloc ...int) *IndexedSlice[K, T0, T1] {
	size := 0
	if len(preAlloc) > 0 {
		size = preAlloc[0]
	}

	return &IndexedSlice[K, T0, T1]{
		index:       make(map[K]int),
		elements:    make([]associative.Pair[*K, T1], 0, size),
		getkey:      getkey,
		initializer: initializer,
		updater:     updater,
		IsEmpty:     isEmpty,
	}
}

func (this *IndexedSlice[K, T0, T1]) Index() map[K]int                     { return this.index }
func (this *IndexedSlice[K, T0, T1]) Elements() []associative.Pair[*K, T1] { return this.elements }
func (this *IndexedSlice[K, T0, T1]) Length(getsize func(T1) int) int {
	total := 0
	for _, ele := range this.elements {
		total += getsize(ele.Second)
	}
	return total
}

// Insert inserts an unique element into the IndexedSlice and updates the index.
// If the element already exists, it is updated. Otherwise, it is added.
// Returns the index of the element in the slice. The function uses its own key extractor to get the key from the element.
// So the identical element will likely have the same key. Unless the update can handle the case, otherwise the value will be overwritten.
// If that is the case, make sure the key extractor can generate unique keys for the elements even if they are identical.
func (this *IndexedSlice[K, T0, T1]) Add(elems ...T0) {
	for _, ele := range elems {
		k := this.getkey(ele)
		idx, ok := this.index[k]
		if ok { // Existing value
			this.updater(k, ele, &(this.elements[idx].Second)) // The updater modifies the value in place with the new value.
			continue
		}
		// New value
		this.index[k] = len(this.elements)                                                                           // Added to the lookup index
		this.elements = append(this.elements, associative.Pair[*K, T1]{First: &k, Second: this.initializer(k, ele)}) // Added to the slice
	}
}

// Insert inserts an element into the IndexedSlice and updates the index with the specified key.
// If the element already exists, it is updated. Otherwise, it is added.
// Returns the index of the element in the slice.
func (this *IndexedSlice[K, T0, T1]) SetByKey(k K, ele T0) {
	idx, ok := this.index[k]
	if ok { // Existing value
		this.updater(k, ele, &(this.elements[idx].Second)) // The updater modifies the value in place with the new value.
		return
	}
	// New value
	this.index[k] = len(this.elements)                                                                           // Added to the lookup index
	this.elements = append(this.elements, associative.Pair[*K, T1]{First: &k, Second: this.initializer(k, ele)}) // Added to the slice
}

func (this *IndexedSlice[K, T0, T1]) GetByKey(t T0) (T1, bool) {
	k := this.getkey(t)
	if idx, ok := this.index[k]; ok {
		return this.elements[idx].Second, true
	}
	return *new(T1), false
}

// CountIf returns the number of elements in the IndexedSlice that satisfy the specified condition.
func (this *IndexedSlice[K, T0, T1]) CountIf(condition func(k *K, v T1) bool) int {
	total := 0
	for i := 0; i < len(this.elements); i++ {
		if condition(this.elements[i].First, this.elements[i].Second) {
			total++
		}
	}
	return total
}

// Array returns the underlying slice of elements in the IndexedSlice.
func (this *IndexedSlice[K, T0, T1]) Values() []T1 {
	elems := make([]T1, len(this.elements))
	for i, ele := range this.elements {
		elems[i] = ele.Second
	}
	return elems
}

// Keys returns the unique keys of the elements in the IndexedSlice.
// Duplicate keys are not included. So it is usally equal or less than the number of elements.
func (this *IndexedSlice[K, T0, T1]) Keys() []K {
	keys := make([]K, len(this.index))
	for i, ele := range this.elements {
		keys[i] = *ele.First
	}
	return keys
}

// Find searches for an element in the IndexedSlice and returns its index.
// Returns -1 if the element is not found.
func (this *IndexedSlice[K, T0, T1]) Find(ele T0) T1 {
	idx := this.index[this.getkey(ele)]
	return this.elements[idx].Second
}

func (this *IndexedSlice[K, T0, T1]) Clear() {
	clear(this.index)
	this.elements = this.elements[:0]
}
