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

// IndexedSlice represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type IndexedSlice[T any, K comparable, V any] struct {
	elements []V
	index    map[K]V
	keys     []K

	getkey   func(T) K
	inserter func(T, V) V
	getsize  func(V) int
}

// NewIndexedSlice creates a new instance of IndexedSlice with the specified page size, minimum number of pages, and pre-allocation size.
func NewIndexedSlice[T any, K comparable, V any](
	getkey func(T) K,
	inserter func(T, V) V,
	getsize func(V) int,
	preAlloc ...int) *IndexedSlice[T, K, V] {
	size := 0
	if len(preAlloc) > 0 {
		size = preAlloc[0]
	}

	return &IndexedSlice[T, K, V]{
		index:    make(map[K]V),
		elements: make([]V, 0, size),
		getkey:   getkey,
		inserter: inserter,
		getsize:  getsize,
	}
}

// Insert inserts an element into the IndexedSlice and updates the index.
func (this *IndexedSlice[T, K, V]) InsertSlice(elements []T) {
	for _, ele := range elements {
		this.Insert(ele)
	}
}

func (this *IndexedSlice[T, K, V]) UniqueLength() int { return len(this.index) }
func (this *IndexedSlice[T, K, V]) Length() int {
	total := 0
	for _, ele := range this.elements {
		total += this.getsize(ele)
	}
	return total
}

// Insert inserts an element into the IndexedSlice and updates the index.
func (this *IndexedSlice[T, K, V]) Insert(ele T) {
	k := this.getkey(ele)
	values, ok := this.index[k]
	if !ok {
		values = this.inserter(ele, values)
		this.index[k] = values

		this.keys = append(this.keys, this.getkey(ele))
		this.elements = append(this.elements, values)
		return
	}
	this.inserter(ele, values)
}

// Array returns the underlying slice of elements in the IndexedSlice.
func (this *IndexedSlice[T, K, V]) Elements() []V {
	return this.elements
}

// Array returns the underlying slice of elements in the IndexedSlice.
func (this *IndexedSlice[T, K, V]) Keys() []K {
	return this.keys
}

// Find searches for an element in the IndexedSlice and returns its index.
// Returns -1 if the element is not found.
func (this *IndexedSlice[T, K, V]) Find(ele T) V {
	return this.index[this.getkey(ele)]
}

func (this *IndexedSlice[T, K, V]) Clear() {
	clear(this.index)
	this.elements = this.elements[:0]
	this.keys = this.keys[:0]
}
