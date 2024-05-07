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

package orderedmap

import (
	"runtime"

	mapi "github.com/arcology-network/common-lib/exp/map"
	"github.com/arcology-network/common-lib/exp/slice"
)

// OrderedMap represents a slice with an dict. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type OrderedMap[K comparable, T, V any] struct {
	dict     map[K]*int
	keys     []K
	values   []V
	init     func(K, T) V
	setter   func(K, T, *V)
	inserter func(K, *[]K, V, *[]V) int // Insert a new element into the slice
	nilValue V
}

// NewIndexedSlice creates a new instance of OrderedMap with the specified page size, minimum number of pages, and pre-allocation size.
func NewOrderedMap[K comparable, T, V any](
	nilValue V,
	preAlloc int,
	init func(K, T) V,
	setter func(K, T, *V),
	inserter func(K, *[]K, V, *[]V) int) *OrderedMap[K, T, V] {
	set := &OrderedMap[K, T, V]{
		dict:     make(map[K]*int),
		keys:     make([]K, 0, preAlloc),
		values:   make([]V, 0, preAlloc),
		init:     init,
		setter:   setter,
		inserter: inserter,
		nilValue: nilValue,
	}
	return set.Init()
}

func (this *OrderedMap[K, T, V]) Init() *OrderedMap[K, T, V] {
	clear(this.dict)
	for i, idx := range this.keys {
		this.dict[idx] = &i
	}
	return this
}

func (this *OrderedMap[K, T, V]) Dict() map[K]*int { return this.dict }
func (this *OrderedMap[K, T, V]) Keys() []K        { return this.keys }
func (this *OrderedMap[K, T, V]) Values() []V      { return this.values }
func (this *OrderedMap[K, T, V]) KVs() ([]K, []V)  { return this.keys, this.values }
func (this *OrderedMap[K, T, V]) Length() int      { return len(this.values) }

func (this *OrderedMap[K, T, V]) Clone() *OrderedMap[K, T, V] {
	cloned := NewOrderedMap[K, T, V](this.nilValue, len(this.values), this.init, this.setter, this.inserter)
	cloned.keys = slice.Clone(this.keys)
	cloned.values = slice.Clone(this.values)
	return cloned.Init()
}

func (this *OrderedMap[K, T, V]) ForeachDo(do func(K, V)) { // For encoding
	slice.Foreach(this.keys, func(i int, k *K) {
		do(*k, this.values[i])
	})
}

func (this *OrderedMap[K, T, V]) ParallelForeachDo(do func(K, *V)) { // For encoding
	slice.ParallelForeach(this.keys, runtime.NumCPU(), func(i int, k *K) {
		do(*k, &(this.values[i]))
	})
}

// Insert inserts an element into the OrderedMap and updates the dict with the specified key.
// If the element already exists, it is updated. Otherwise, it is added.
// Returns the dict of the element in the slice.
func (this *OrderedMap[K, T, V]) Insert(keys []K, vals []T) *OrderedMap[K, T, V] {
	for i, k := range keys {
		this.Set(k, vals[i])
	}
	return this
}

func (this *OrderedMap[K, T, V]) KeyToIndex(k K) int {
	if idx, ok := this.dict[k]; ok {
		return *idx
	}
	return -1
}

func (this *OrderedMap[K, T, V]) IndexToKey(idx int) K {
	return this.keys[idx]
}

func (this *OrderedMap[K, T, V]) Set(k K, v T) int {
	idx, ok := this.dict[k]
	if ok { // Existing entry
		this.setter(k, v, &this.values[*idx])
		return *idx
	}
	newv := this.init(k, v)

	// Default inserter, simply append the new element to the slice.
	if this.inserter == nil {
		this.values = append(this.values, newv)
		this.keys = append(this.keys, k)
		length := len(this.values) - 1
		this.dict[k] = &length
		return length
	}

	// Custom inserter.
	newIdx := this.inserter(k, &this.keys, newv, &this.values)
	this.dict[k] = &newIdx

	// Indices are shifted after the new element is inserted,so we need to update the dict.
	for i := newIdx + 1; i < len(this.keys); i++ {
		v := this.dict[this.keys[i]]
		*v = i
	}
	return newIdx
}

func (this *OrderedMap[K, T, V]) DoSet(k K, v T, setter func(*[]K, *[]V)) {
	idx, ok := this.dict[k]
	if !ok { // New entries
		this.values = append(this.values, this.init(k, v))
		this.keys = append(this.keys, k)
		length := len(this.values) - 1
		this.dict[k] = &length
		return
	}
	this.setter(k, v, &this.values[*idx])
}

func (this *OrderedMap[K, T, V]) Get(k K) (V, bool) {
	if idx, ok := this.dict[k]; ok {
		return this.values[*idx], ok
	}
	return this.nilValue, false
}

func (this *OrderedMap[K, T, V]) At(idx int) (K, V) {
	return this.keys[idx], this.values[idx]
}

func (this *OrderedMap[K, T, V]) Exists(k K) bool {
	_, ok := this.dict[k]
	return ok
}

func (this *OrderedMap[K, T, V]) Clear() {
	clear(this.dict)
	this.keys = this.keys[:0]
	this.values = this.values[:0]
}

// Debugging function to check if the dict is in sync with the slice.
func (this *OrderedMap[K, T, V]) IsDirty() bool { return len(this.values) != len(this.dict) }

func (this *OrderedMap[K, T, V]) DeleteByIndex(indices ...int) {
	for _, idx := range indices {
		delete(this.dict, this.keys[idx]) // remove the old key
		slice.RemoveAt(&this.keys, idx)
		slice.RemoveAt(&this.values, idx)
	}

	idx, _ := slice.Min(indices)
	// Only need to update the elements after the last deleted index, previous elements' indices are not affected.
	for i, k := range this.keys[idx:] {
		*this.dict[k] = i + idx
	}
}

func (this *OrderedMap[K, T, V]) Delete(k K) bool {
	idx, ok := this.dict[k]
	if !ok {
		return false
	}
	slice.RemoveAt(&this.keys, *idx)
	slice.RemoveAt(&this.values, *idx)
	delete(this.dict, k)

	// Only update the dict for the elements after the deleted element.
	for i := *idx; i < len(this.keys); i++ {
		*this.dict[this.keys[i]] = i
	}
	return true
}

func (this *OrderedMap[K, T, V]) DeleteBatch(keys ...K) bool {
	dict := mapi.FromSlice(keys, func(k K) *int { return this.dict[k] })
	for _, k := range keys {
		idx := this.dict[k]
		slice.RemoveAt(&this.keys, *idx)

	}

	minIdx := len(this.values)
	slice.RemoveBothIf(&this.keys, &this.values, func(i int, k K, _ V) bool {
		idx, ok := dict[k]
		if ok && *idx < minIdx {
			minIdx = *idx
		}
		return ok
	})

	//Remove the keys from the dictionary
	for _, k := range keys {
		delete(this.dict, k)
	}

	// Some elements may have been removed, so there are some gaps in the slice. The dictionary
	// no longer reflects the correct index of the elements. This function will reorder the elements
	// in the slice and update the dict accordingly.
	for i, k := range this.keys {
		*this.dict[k] = i
	}
	return false
}
