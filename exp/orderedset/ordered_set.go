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

package orderedset

import (
	"fmt"
	"sort"

	mapi "github.com/arcology-network/common-lib/exp/map"
	"github.com/arcology-network/common-lib/exp/slice"
)

// OrderedSet represents a slice with an dict. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type OrderedSet[K comparable, T any] struct {
	elements []T
	dict     map[K]int
	nilValue K
	getter   func(T) K
}

// NewIndexedSlice creates a new instance of OrderedSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewOrderedSet[K comparable, T any](nilValue K, getter func(T) K, preAlloc int, vals ...T) *OrderedSet[K, T] {
	set := &OrderedSet[K, T]{
		dict:     make(map[K]int),
		elements: append(make([]T, 0, preAlloc+len(vals)), vals...),
		nilValue: nilValue,
		getter:   getter,
	}
	return set.Init()
}

func (this *OrderedSet[K, T]) Init() *OrderedSet[K, T] {
	for i, v := range this.elements {
		this.dict[this.getter(v)] = i
	}
	return this
}

func (this *OrderedSet[K, T]) Getter() func(T) K { return this.getter }
func (this *OrderedSet[K, T]) Dict() map[K]int   { return this.dict }
func (this *OrderedSet[K, T]) Elements() []T     { return this.elements }
func (this *OrderedSet[K, T]) Length() int       { return len(this.elements) }
func (this *OrderedSet[K, T]) Clone() *OrderedSet[K, T] {
	return NewOrderedSet[K, T](this.nilValue, this.getter, len(this.elements), this.elements...)
}

func (this *OrderedSet[K, T]) Size(getter func(K) int) int { // For encoding
	return slice.Accumulate(this.elements, 0, func(acc int, v T) int { return acc + getter(this.getter(v)) })
}

func (this *OrderedSet[K, T]) Merge(elements []T) *OrderedSet[K, T] {
	for _, ele := range elements {
		this.Insert(ele)
	}
	return this
}

func (this *OrderedSet[K, T]) Sub(elements []T) *OrderedSet[K, T] {
	for _, ele := range elements {
		this.Delete(ele)
	}
	return this
}

// Insert inserts an element into the OrderedSet and updates the dict with the specified key.
// If the element already exists, it is updated. Otherwise, it is added.
// Returns the dict of the element in the slice.
func (this *OrderedSet[K, T]) Insert(vals ...T) {
	for _, v := range vals {
		k := this.getter(v)
		if _, ok := this.dict[k]; !ok { // New entries
			this.elements = append(this.elements, v)
			this.dict[k] = len(this.elements) - 1
		}
	}
}

func (this *OrderedSet[K, T]) At(idx int) *T {
	return &this.elements[idx]
}

func (this *OrderedSet[K, T]) IndexToKey(k K) int {
	if idx, ok := this.dict[k]; ok {
		return idx
	}
	return -1
}

func (this *OrderedSet[K, T]) KeyToIndex(idx int) T {
	return this.elements[idx]
}

func (this *OrderedSet[K, T]) DeleteByIndex(indices ...int) {
	for _, idx := range indices {
		delete(this.dict, this.getter(this.elements[idx])) // remove the old key
		slice.RemoveAt(&this.elements, idx)
	}
}

func (this *OrderedSet[K, T]) Delete(vals ...T) bool {
	removed := make([]int, len(vals))
	for i, v := range vals {
		k := this.getter(v)
		if idx, ok := this.dict[k]; ok {
			slice.RemoveAt(&this.elements, idx)
			delete(this.dict, k)
			removed[i] = idx
		}
	}
	this.Reorder(removed...)
	return false
}

// Reorder order the indexes of the elements with the dict.
func (this *OrderedSet[K, T]) Reorder(offsets ...int) {
	sort.Ints(offsets)
	offsets = append(offsets, len(this.elements))
	for i := 0; i < len(offsets)-1; i++ {
		for j := offsets[i]; j < offsets[i+1]; j++ {
			k := this.getter(this.elements[j])
			this.dict[k] = this.dict[k] - 1
		}
	}
}

func (this *OrderedSet[K, T]) Exists(v T) (bool, int) {
	outv, ok := this.dict[this.getter(v)]
	return ok, outv
}

func (this *OrderedSet[K, T]) Clear() {
	clear(this.dict)
	this.elements = this.elements[:0]
}

// Debugging function to check if the dict is in sync with the slice.
func (this *OrderedSet[K, T]) IsDirty() bool {
	if len(this.elements) != len(this.dict) {
		return false
	}

	for i, v := range this.elements {
		k := this.getter(this.elements[i])
		if this.dict[k] != i {
			fmt.Printf("Index out of sync: %v, %v, %v\n", i, v, this.dict[k])
			return false
		}
	}
	return true
}

func (this *OrderedSet[K, T]) Equal(other *OrderedSet[K, T]) bool {
	equal := func(v0 T, v1 T) bool {
		return this.getter(v0) == this.getter(v1)
	}
	return slice.EqualSetIf(this.elements, other.elements, equal) && mapi.EqualIf(this.dict, other.dict, func(v0 int, v1 int) bool { return v0 == v1 })
}

func (this *OrderedSet[K, T]) Print() {
	fmt.Println(this.dict, this.elements)
}

// This is for debug purpose only !!, don't use it in production
// since it has some quite complicated consequences. !!!
func (this *OrderedSet[K, T]) replace(idx int, v T) T {
	old := this.elements[idx]
	delete(this.dict, this.getter(this.elements[idx])) // remove the old key
	this.elements[idx] = v                             // update the value
	this.dict[this.getter(this.elements[idx])] = idx   // update the dict
	return old
}
