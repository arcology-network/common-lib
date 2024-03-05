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
	associative "github.com/arcology-network/common-lib/exp/associative"
	mapi "github.com/arcology-network/common-lib/exp/map"
)

type IndexDeltaSlice[T any, K comparable] struct {
	*DeltaSlice[associative.Pair[K, T]]
	index map[K]int

	addedIndex   map[K]int
	removedIndex map[K]int
}

func NewIndexDeltaSlice[T any, K comparable]() *IndexDeltaSlice[T, K] {
	return &IndexDeltaSlice[T, K]{
		DeltaSlice:   NewDeltaSlice[associative.Pair[K, T]](100),
		index:        make(map[K]int),
		addedIndex:   make(map[K]int),
		removedIndex: make(map[K]int),
	}
}

func (this *IndexDeltaSlice[T, K]) GetByKey(key K) *T {
	if idx, ok := this.index[key]; ok {
		if v, ok := this.DeltaSlice.Get(idx); ok {
			return &v.Second
		}
	}
	return nil
}

// SetByKey set value by key, if key not exist, append to delta slice
// It is equal to a deltaSlice append operation with a index map insertion.
func (this *IndexDeltaSlice[T, K]) SetByKey(key K, newV T) {
	if idx, ok := this.addedIndex[key]; ok { // Already added, it is effectively a value update.
		if v, ok := this.DeltaSlice.Get(idx); ok {
			v.Second = newV
			return
		}
		panic("Invalid added index")
	}

	idx, ok := this.index[key]
	if !ok {
		pair := associative.Pair[K, T]{First: key, Second: newV}
		idx = this.DeltaSlice.Append(pair) // Add to delta slice
	}
	this.index[key] = idx // Add to index
}

func (this *IndexDeltaSlice[T, K]) GetByIndex(index int) *T {
	if index >= this.DeltaSlice.Length() {
		return nil
	}

	v, _ := this.DeltaSlice.Get(index)
	return &v.Second
}

func (this *IndexDeltaSlice[T, K]) SetByIndex(index int, newV T) bool {
	if index >= this.DeltaSlice.Length() {
		return false
	}

	if pair, ok := this.DeltaSlice.Get(index); ok {
		pair.Second = newV
	}
	return true
}

func (this *IndexDeltaSlice[T, K]) IndexToKey(index int) *K {
	if index >= this.DeltaSlice.Length() {
		return nil
	}
	pair, _ := this.DeltaSlice.Get(index)
	return &pair.First // The map key of the index
}

func (this *IndexDeltaSlice[T, K]) KeyToIndex(key K) int {
	idx, ok := this.index[key]
	if ok { // In the main index
		if _, ok := this.removedIndex[key]; ok { // Check if the key has been deleted
			return -1 // THe key is in the removed index so it has been deleted already.
		}
		return idx
	}
	return -1
}

func (this *IndexDeltaSlice[T, K]) DeleteByKey(key K) bool {
	if _, ok := this.removedIndex[key]; ok { // The key has been deleted already. No need to delete again.
		return true
	}

	idx := this.KeyToIndex(key)
	if idx == -1 {
		return false // Not found
	}

	ok, from := this.DeltaSlice.Del(idx) // Delete from delta slice
	if !ok {
		return false // Seomthing wrong, this really should not happen.
	}

	switch {
	case from == 1:
		this.removedIndex[key] = idx // Delete from committed set
	case from == 2:
		delete(this.addedIndex, key) // Remove from added index
	default:
		panic("Invalid from value")
	}
	return true
}

func (this *IndexDeltaSlice[T, K]) DeleteByIndex(index int) bool {
	if index >= this.DeltaSlice.Length() {
		return false
	}

	success, _ := this.DeltaSlice.Del(index)
	return success
}

func (this *IndexDeltaSlice[T, K]) Commit() {
	this.DeltaSlice.Commit()

	mapi.Sub(this.index, this.removedIndex)
	mapi.Merge(this.index, this.addedIndex)

	clear(this.addedIndex)
	clear(this.removedIndex)
}
