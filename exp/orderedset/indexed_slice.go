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

// OrderedSet represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type OrderedSet[K comparable] struct {
	elements []K
	index    map[K]int
	nilValue K
}

// NewIndexedSlice creates a new instance of OrderedSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewOrderedSet[K comparable](nilValue K, preAlloc int, vals ...K) *OrderedSet[K] {
	set := &OrderedSet[K]{
		index:    make(map[K]int),
		elements: make([]K, 0, preAlloc+len(vals)),
		nilValue: nilValue,
	}
	set.Append(vals...)
	return set
}

func (this *OrderedSet[K]) Index() map[K]int { return this.index }
func (this *OrderedSet[K]) Elements() []K    { return this.elements }
func (this *OrderedSet[K]) Length() int      { return len(this.elements) }
func (this *OrderedSet[K]) Clone() *OrderedSet[K] {
	return NewOrderedSet[K](this.nilValue, len(this.elements), this.elements...)
}

func (this *OrderedSet[K]) Append(other ...K) *OrderedSet[K] {
	this.elements = append(this.elements, other...)
	for i := len(this.elements) - len(other); i < len(this.elements); i++ {
		this.index[this.elements[i]] = i
	}
	return this
}

func (this *OrderedSet[K]) Merge(other *OrderedSet[K]) {
	for _, ele := range other.Elements() {
		this.Insert(ele)
	}
}

func (this *OrderedSet[K]) Sub(other *OrderedSet[K]) {
	for _, ele := range other.Elements() {
		this.Delete(ele)
	}
}

// Insert inserts an element into the OrderedSet and updates the index with the specified key.
// If the element already exists, it is updated. Otherwise, it is added.
// Returns the index of the element in the slice.
func (this *OrderedSet[K]) Insert(k K) (int, bool) {
	if _, ok := this.index[k]; !ok {
		this.elements = append(this.elements, k)
		this.index[k] = len(this.elements) - 1
		return len(this.elements) - 1, true
	}
	return -1, false
}

func (this *OrderedSet[K]) At(idx int) *K {
	return &this.elements[idx]
}

func (this *OrderedSet[K]) KeyToIndex(k K) int {
	if idx, ok := this.index[k]; ok {
		return idx
	}
	return -1
}

func (this *OrderedSet[K]) IndexToKey(idx int) K {
	return this.elements[idx]
}

func (this *OrderedSet[K]) Replace(idx int, v K) K {
	old := this.elements[idx]
	delete(this.index, this.elements[idx]) // remove the old key
	this.elements[idx] = v                 // update the value
	this.index[this.elements[idx]] = idx   // update the index
	return old
}

func (this *OrderedSet[K]) DeleteByIndex(indices ...int) {
	for _, idx := range indices {
		delete(this.index, this.elements[idx]) // remove the old key
		slice.RemoveAt(&this.elements, idx)
	}
}

func (this *OrderedSet[K]) Delete(keys ...K) bool {
	removed := make([]int, len(keys))
	for i, k := range keys {
		if idx, ok := this.index[k]; ok {
			slice.RemoveAt(&this.elements, idx)
			delete(this.index, k)
			removed[i] = idx
		}
	}
	this.Sync(removed...)
	return false
}

func (this *OrderedSet[K]) Sync(offsets ...int) {
	sort.Ints(offsets)
	offsets = append(offsets, len(this.elements))
	for i := 0; i < len(offsets)-1; i++ {
		for j := offsets[i]; j < offsets[i+1]; j++ {
			k := this.elements[j]
			this.index[k] = this.index[k] - 1
		}
	}
}

func (this *OrderedSet[K]) Exists(k K) bool {
	_, ok := this.index[k]
	return ok
}

func (this *OrderedSet[K]) Clear() {
	clear(this.index)
	this.elements = this.elements[:0]
}

// Debugging function to check if the index is in sync with the slice.
func (this *OrderedSet[K]) IsSynced() bool {
	if len(this.elements) != len(this.index) {
		return false
	}

	for i, v := range this.elements {
		if this.index[v] != i {
			fmt.Printf("Index out of sync: %v, %v, %v\n", i, v, this.index[v])
			return false
		}
	}
	return true
}

func (this *OrderedSet[K]) Equal(other *OrderedSet[K]) bool {
	return slice.Equal(this.elements, other.elements) && mapi.EqualIf(this.index, other.index, func(v0 int, v1 int) bool { return v0 == v1 })
}

func (this *OrderedSet[K]) Print() {
	fmt.Println(this.index, this.elements)
}
