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

	"github.com/arcology-network/common-lib/common"
	mapi "github.com/arcology-network/common-lib/exp/map"
	"github.com/arcology-network/common-lib/exp/slice"
)

// OrderedSet represents a slice with an dict. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type OrderedSet[K comparable] struct {
	elements []K
	index    []int
	dict     map[K]*int
	nilValue K
	hasher   func(K) [32]byte
}

// NewIndexedSlice creates a new instance of OrderedSet with the specified page size, minimum number of pages, and pre-allocation size.
func NewOrderedSet[K comparable](
	nilValue K,
	preAlloc int,
	hasher func(K) [32]byte,
	vals ...K) *OrderedSet[K] {
	set := &OrderedSet[K]{
		dict:     make(map[K]*int),
		elements: append(make([]K, 0, preAlloc+len(vals)), vals...),
		nilValue: nilValue,
		hasher:   hasher,
	}

	return set.Init()
}

func (this *OrderedSet[K]) Init() *OrderedSet[K] {
	for i, idx := range this.elements {
		this.dict[idx] = common.New(i)
	}
	return this
}

func (this *OrderedSet[K]) Hasher() func(K) [32]byte { return this.hasher }
func (this *OrderedSet[K]) Dict() map[K]*int         { return this.dict }
func (this *OrderedSet[K]) Elements() []K            { return this.elements }
func (this *OrderedSet[K]) Length() int              { return len(this.elements) }
func (this *OrderedSet[K]) Clone() *OrderedSet[K] {
	return NewOrderedSet(this.nilValue, len(this.elements), nil, this.elements...)
}

func (this *OrderedSet[K]) Size(getter func(K) int) int { // For encoding
	return slice.Accumulate(this.elements, 0, func(acc int, k K) int { return acc + getter(k) })
}

func (this *OrderedSet[K]) Merge(elements []K) *OrderedSet[K] {
	for _, ele := range elements {
		this.Insert(ele)
	}
	return this
}

func (this *OrderedSet[K]) Sub(elements []K) *OrderedSet[K] {
	for _, ele := range elements {
		this.Delete(ele)
	}
	return this
}

// Insert inserts an element into the OrderedSet and updates the dict with the specified key.
// If the element already exists, it is updated. Otherwise, it is added.
// Returns the dict of the element in the slice.
func (this *OrderedSet[K]) Insert(keys ...K) {
	for _, k := range keys {
		if _, ok := this.dict[k]; !ok { // New entries
			this.dict[k] = common.New(len(this.elements))
			this.elements = append(this.elements, k)
		}
	}
}

// Insert inserts an element into the OrderedSet and updates the dict with the specified key.
// If the element already exists, it is updated. Otherwise, it is added.
// Returns the dict of the element in the slice.
// func (this *OrderedSet[K]) InsertAfter(k K) {
// 	pos := this.getter(k)

// 	idx := sort.Search(len(this.elements), func(i int) bool { return a[i] >= x })
// 	slice.Insert(&a, idx, x)

// 	if _, ok := this.dict[k]; !ok { // New entries
// 		this.dict[k] = common.New(len(this.elements))
// 		this.elements = append(this.elements, k)
// 	}
// }

// SetAt sets the element at the specified index to the new value.
// The dict is updated with the new key.
func (this *OrderedSet[K]) SetAt(idx int, newv K) bool {
	if idx < 0 || idx >= len(this.elements) {
		return false
	}

	delete(this.dict, this.elements[idx]) // remove the old key from the dict
	this.elements[idx] = newv             // Replace the old key with the new key
	this.dict[newv] = common.New(idx)     // Add the new key to the dict
	return true
}

func (this *OrderedSet[K]) At(idx int) *K {
	return &this.elements[idx]
}

func (this *OrderedSet[K]) KeyToIndex(k K) int {
	if idx, ok := this.dict[k]; ok {
		return *idx
	}
	return -1
}

func (this *OrderedSet[K]) IndexToKey(idx int) K {
	return this.elements[idx]
}

func (this *OrderedSet[K]) DeleteByIndex(indices ...int) {
	for _, idx := range indices {
		delete(this.dict, this.elements[idx]) // remove the old key
		slice.RemoveAt(&this.elements, idx)
	}

	idx, _ := slice.Min(indices)
	for i, k := range this.elements[idx:] {
		*this.dict[k] = i + idx
	}
}

func (this *OrderedSet[K]) Delete(keys ...K) bool {
	dict := mapi.FromSlice(keys, func(k K) *int { return this.dict[k] })
	for _, k := range keys {
		delete(this.dict, k)
	}

	minIdx := len(this.elements)
	slice.RemoveIf(&this.elements, func(i int, k K) bool {
		idx, ok := dict[k]
		if ok && *idx < minIdx {
			minIdx = *idx
		}
		return ok
	})

	// Some elements may have been removed, so there are some gaps in the slice. The dictionary
	// no longer reflects the correct index of the elements. This function will reorder the elements
	// in the slice and update the dict accordingly.
	for i, k := range this.elements {
		*this.dict[k] = i
	}
	return false
}

func (this *OrderedSet[K]) Exists(k K) (bool, int) {
	if v, ok := this.dict[k]; ok {
		return ok, *v
	}
	return false, -1
}

func (this *OrderedSet[K]) Clear() {
	clear(this.dict)
	this.elements = this.elements[:0]
}

// Debugging function to check if the dict is in sync with the slice.
func (this *OrderedSet[K]) IsDirty() bool { return len(this.elements) != len(this.dict) }

func (this *OrderedSet[K]) Equal(other *OrderedSet[K]) bool {
	return slice.EqualSet(this.elements, other.elements) && mapi.EqualIf(this.dict, other.dict, func(v0 *int, v1 *int) bool { return *v0 == *v1 })
}

func (this *OrderedSet[K]) Print() {
	fmt.Println(this.dict, this.elements)
}

// Count the number of elements BEFORE the specified key, not including the key itself.
func (this *OrderedSet[K]) CountBefore(key K) int {
	if idx, ok := this.dict[key]; ok {
		return *idx
	}
	return -1
}

// Count the number of elements AFTER the specified index, not including the key itself.
func (this *OrderedSet[K]) CountAfter(key K) int {
	if idx, ok := this.dict[key]; ok {
		return len(this.dict) - *idx - 1
	}
	return -1
}
