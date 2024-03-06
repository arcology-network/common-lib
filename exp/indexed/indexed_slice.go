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
	"github.com/arcology-network/common-lib/exp/slice"
)

// IndexedSlice represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type IndexedSlice[K comparable, T0 any, T1 any] struct {
	elements []*associative.Triplet[*K, uint64, T1]
	index    map[K]*associative.Triplet[*K, uint64, T1]

	ToKey       func(T0) K       // Key extractor to get the key from the element.
	initializer func(K, T0) T1   // Initializer when adding a new element.
	updater     func(K, T0, *T1) // Updater when updating an existing element.
	deleter     func(*T1)
	IsEmpty     func(T1) bool // Checker for empty element.
}

// NewIndexedSlice creates a new instance of IndexedSlice with the specified page size, minimum number of pages, and pre-allocation size.
func NewIndexedSlice[K comparable, T0 any, T1 any](
	ToKey func(T0) K,
	initializer func(K, T0) T1,
	updater func(K, T0, *T1),
	deleter func(*T1),
	isEmpty func(T1) bool,
	preAlloc ...int) *IndexedSlice[K, T0, T1] {
	size := 0
	if len(preAlloc) > 0 {
		size = preAlloc[0]
	}

	return &IndexedSlice[K, T0, T1]{
		index:       make(map[K]*associative.Triplet[*K, uint64, T1]),
		elements:    make([]*associative.Triplet[*K, uint64, T1], 0, size),
		ToKey:       ToKey,
		initializer: initializer,
		updater:     updater,
		deleter:     deleter,
		IsEmpty:     isEmpty,
	}
}

func (this *IndexedSlice[K, T0, T1]) New() *IndexedSlice[K, T0, T1] {
	return NewIndexedSlice(this.ToKey, this.initializer, this.updater, this.deleter, this.IsEmpty, len(this.elements))
}

func (this *IndexedSlice[K, T0, T1]) Index() map[K]*associative.Triplet[*K, uint64, T1] {
	return this.index
}
func (this *IndexedSlice[K, T0, T1]) Elements() *[]*associative.Triplet[*K, uint64, T1] {
	return &this.elements
}

func (this *IndexedSlice[K, T0, T1]) Length(getsize func(T1) int) int {
	total := 0
	for _, ele := range this.elements {
		total += getsize(ele.Third)
	}
	return total
}

func (this *IndexedSlice[K, T0, T1]) Append(other *associative.Triplet[*K, uint64, T1]) {
	this.elements = append(this.elements, other)
	this.index[*other.First] = this.elements[len(this.elements)-1]
}

func (this *IndexedSlice[K, T0, T1]) Merge(other *IndexedSlice[K, T0, T1]) {
	for i, ele := range *other.Elements() {
		ele.Second = uint64(len(this.elements) + i)
		this.elements = append(this.elements, ele)
		this.index[*ele.First] = this.elements[len(this.elements)-1]
	}
}

// Insert inserts an unique element into the IndexedSlice and updates the index.
// If the element already exists, it is updated. Otherwise, it is added.
// Returns the index of the element in the slice. The function uses its own key extractor to get the key from the element.
// So the identical element will likely have the same key. Unless the update can handle the case, otherwise the value will be overwritten.
// If that is the case, make sure the key extractor can generate unique keys for the elements even if they are identical.
func (this *IndexedSlice[K, T0, T1]) Add(elems ...T0) {
	for _, ele := range elems {
		k := this.ToKey(ele)
		triplet, ok := this.index[k]
		if ok { // Existing value
			this.updater(k, ele, &(triplet.Third)) // The updater modifies the value in place with the new value.
			continue
		}
		// New value
		this.index[k] = &associative.Triplet[*K, uint64, T1]{
			First:  &k,
			Second: uint64(len(this.elements)),
			Third:  this.initializer(k, ele)} // Added to the lookup index

		this.elements = append(this.elements, this.index[k]) // Added to the slice
	}
}

// Insert inserts an element into the IndexedSlice and updates the index with the specified key.
// If the element already exists, it is updated. Otherwise, it is added.
// Returns the index of the element in the slice.
func (this *IndexedSlice[K, T0, T1]) SetByKey(k K, ele T0) {
	triplet, ok := this.index[k]
	if ok { // Existing value
		this.updater(k, ele, &(triplet.Third)) // The updater modifies the value in place with the new value.
		return
	}
	// New value
	this.index[k] = &associative.Triplet[*K, uint64, T1]{
		First:  &k,
		Second: uint64(len(this.elements)),
		Third:  this.initializer(k, ele)} // Added to the lookup index

	this.elements = append(this.elements, this.index[k]) // Added to the slice
}

func (this *IndexedSlice[K, T0, T1]) GetByKey(k K) (T1, uint64, bool) {
	if triplet, ok := this.index[k]; ok {
		return triplet.Third, triplet.Second, true
	}
	return *new(T1), 0, false
}

func (this *IndexedSlice[K, T0, T1]) KeyToIndex(k K) (uint64, bool) {
	if triplet, ok := this.index[k]; ok {
		return triplet.Second, true
	}
	return 0, false
}

func (this *IndexedSlice[K, T0, T1]) IndexToKey(idx uint64) (K, bool) {
	if idx < 0 || idx >= uint64(len(this.elements)) {
		return *new(K), false
	}
	return *this.elements[idx].First, true
}

func (this *IndexedSlice[K, T0, T1]) GetByIndex(idx int) (T1, bool) {
	if idx < 0 || idx >= len(this.elements) {
		return *new(T1), false
	}
	return this.elements[idx].Third, true
}

func (this *IndexedSlice[K, T0, T1]) SetByIndex(idx int, v T1) bool {
	if idx < 0 || idx >= len(this.elements) {
		return false
	}
	this.elements[idx].Third = v
	return true
}

func (this *IndexedSlice[K, T0, T1]) DeleteByIndex(indices ...int) bool {
	for _, idx := range indices {
		this.deleter(&this.elements[idx].Third)
	}

	slice.RemoveIf(&this.elements, func(i int, v *associative.Triplet[*K, uint64, T1]) bool {
		return this.IsEmpty(v.Third)
	})

	return true
}

func (this *IndexedSlice[K, T0, T1]) Exists(idx int, v T1) bool {
	if idx < 0 || idx >= len(this.elements) {
		return false
	}
	this.elements[idx].Third = v
	return true
}

// CountIf returns the number of elements in the IndexedSlice that satisfy the specified condition.
func (this *IndexedSlice[K, T0, T1]) CountIf(condition func(k *K, v T1) bool) int {
	total := 0
	for i := 0; i < len(this.elements); i++ {
		if condition(this.elements[i].First, this.elements[i].Third) {
			total++
		}
	}
	return total
}

// Array returns the underlying slice of elements in the IndexedSlice.
func (this *IndexedSlice[K, T0, T1]) Values() []T1 {
	elems := make([]T1, len(this.elements))
	for i, ele := range this.elements {
		elems[i] = ele.Third
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
	triplet := this.index[this.ToKey(ele)]
	return triplet.Third
}

func (this *IndexedSlice[K, T0, T1]) Clear() {
	clear(this.index)
	this.elements = this.elements[:0]
}
