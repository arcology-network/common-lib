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

	"github.com/arcology-network/common-lib/exp/slice"
	"golang.org/x/crypto/sha3"
)

// OrderedMap represents a slice with an dict. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type OrderedMap[K comparable, T, V any] struct {
	dict     map[K]*int
	keys     []K
	values   []V
	init     func(K, T) V
	setter   func(K, T, *V)
	nilValue V
}

// NewIndexedSlice creates a new instance of OrderedMap with the specified page size, minimum number of pages, and pre-allocation size.
func NewOrderedMap[K comparable, T, V any](
	nilValue V,
	preAlloc int,
	init func(K, T) V,
	setter func(K, T, *V)) *OrderedMap[K, T, V] {
	set := &OrderedMap[K, T, V]{
		dict:     make(map[K]*int),
		keys:     make([]K, 0, preAlloc),
		values:   make([]V, 0, preAlloc),
		init:     init,
		setter:   setter,
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
	cloned := NewOrderedMap[K, T, V](this.nilValue, len(this.values), this.init, this.setter)
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

func (this *OrderedMap[K, T, V]) InsertDo(vals []T, getter func(int, T) K) *OrderedMap[K, T, V] {
	for i, v := range vals {
		this.Set(getter(i, v), v)
	}
	return this
}

func (this *OrderedMap[K, T, V]) Set(k K, v T) {
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

func (this *OrderedMap[K, T, V]) KeyToIndex(k K) int {
	if idx, ok := this.dict[k]; ok {
		return *idx
	}
	return -1
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

func (this *OrderedMap[K, T, V]) Checksum(less func(K, K), encoder func(K, V) ([]byte, []byte)) [32]byte {
	kByteArr, vByteArr := make([][]byte, len(this.keys)), make([][]byte, len(this.values))
	for i := 0; i < len(this.keys); i++ {
		kByteArr[i], vByteArr[i] = encoder(this.keys[i], this.values[i])
	}
	return sha3.Sum256(append(slice.Flatten(kByteArr), slice.Flatten(vByteArr)...))
}

// Debugging function to check if the dict is in sync with the slice.
func (this *OrderedMap[K, T, V]) IsDirty() bool {
	return len(this.values) != len(this.dict)
}

// func (this *OrderedMap[K, T, V]) Equal(other *OrderedMap[K, T, V]) bool {
// 	return slice.EqualSet(this.values, other.values) && mapi.EqualIf(this.dict, other.dict, func(v0 int, v1 int) bool { return v0 == v1 })
// }

// func (this *OrderedMap[K, T, V]) Print() {
// 	fmt.Println(this.dict, this.values)
// }

// This is for debug purpose only !!, don't use it in production
// since it has some quite complicated consequences. !!!
// func (this *OrderedMap[K, T, V]) replace(idx int, v K) K {
// 	old := this.values[idx]
// 	delete(this.dict, this.values[idx]) // remove the old key
// 	this.values[idx] = v                // update the value
// 	this.dict[this.values[idx]] = idx   // update the dict
// 	return old
// }
