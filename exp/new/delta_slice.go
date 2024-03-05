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
	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/exp/associative"
	indexedslice "github.com/arcology-network/common-lib/exp/indexed"
	"github.com/arcology-network/common-lib/exp/slice"
)

// DeltaSlice represents a slice with an index. It is a hybrid combining a slice and a map support fast lookups and iteration.
// Entries with the same key are stored in a slice in the order they were inserted.
type DeltaSlice[K comparable, T any] struct {
	readonlyValues []T                                 // The values are applied to the values only when commit is called.
	modified       *indexedslice.IndexedSlice[K, T, T] // A hybrid structure to record the modified elements.
	appended       []T
	deleter        func(*T) T
	isDeleted      func(T) bool
}

// NewIndexedSlice creates a new instance of DeltaSlice with the specified page size, minimum number of pages, and pre-allocation size.
func NewDeltaSlice[K comparable, T any](modified *indexedslice.IndexedSlice[K, T, T],
	deleter func(*T) T,
	isDeleted func(T) bool,
) *DeltaSlice[K, T] {
	return &DeltaSlice[K, T]{
		readonlyValues: make([]T, 0, 100),
		modified:       modified,
		appended:       []T{},
		deleter:        deleter,
		isDeleted:      isDeleted,
	}
}

func (this *DeltaSlice[K, T]) mapTo(idx int) (*[]T, int) {
	if idx >= this.Length() {
		return nil, -1
	}

	// The index is in the appended list
	if len(this.readonlyValues) <= idx {
		return &this.appended, idx - (len(this.readonlyValues))
	}
	return &this.readonlyValues, idx
}

// Array returns the underlying slice of readonlyValues in the DeltaSlice.
func (this *DeltaSlice[K, T]) Values() []T                                   { return this.readonlyValues }
func (this *DeltaSlice[K, T]) Updated() []T                                  { return this.modified.Values() }
func (this *DeltaSlice[K, T]) Appended() []T                                 { return this.appended } // Returns the appended readonlyValues, some may be removed later.
func (this *DeltaSlice[K, T]) Modified() *indexedslice.IndexedSlice[K, T, T] { return this.modified }

func (this *DeltaSlice[K, T]) Length() int {
	elems := this.modified.Elements()
	numRemoved := slice.CountIf[associative.Pair[*K, T], int](elems, func(_ int, v *associative.Pair[*K, T]) bool {
		return v == nil
	})
	return len(this.readonlyValues) + len(this.appended) - numRemoved
}

// Insert inserts an element into the DeltaSlice and updates the index.
func (this *DeltaSlice[K, T]) Append(elems ...T) int {
	this.appended = append(this.appended, elems...)
	return this.Length() - 1
}

// ToSlice returns the readonlyValues in the DeltaSlice as a slice by removing the removed readonlyValues and adding the appended readonlyValues.
func (this *DeltaSlice[K, T]) ToSlice() []T {
	return append(this.readonlyValues, this.appended...)
}

// Insert inserts an element into the DeltaSlice and updates the index.
func (this *DeltaSlice[K, T]) Delete(indices ...int) {
	for _, idx := range indices {
		// deleter :=
		this.Set(idx, func(v T) T { return this.deleter(&v) })
	}
}

func (this *DeltaSlice[K, T]) Set(idx int, setter func(T) T) (*T, bool) {
	arr, mapped := this.mapTo(idx)
	if mapped < 0 {
		return nil, false
	}

	this.modified.SetByKey(common.ToType[int, K](idx), setter((*arr)[mapped]))
	return nil, false
}

func (this *DeltaSlice[K, T]) Get(idx int) (T, bool) {
	if v, ok := this.modified.GetByKey(common.ToType[int, T](idx)); ok { // Get from the modified first
		return v, true
	}

	if arr, mapped := this.mapTo(idx); arr != nil {
		return (*arr)[mapped], true
	}
	return *new(T), false
}

// DoAt calls the specified function with the element at the specified index.
// The operation is an in-place operation, so it can modify the element in the DeltaSlice, even
// if the element is in the readonlyValues.
func (this *DeltaSlice[K, T]) DoAt(idx int, doer func(*T)) {
	arr, mapped := this.mapTo(idx)
	if mapped < 0 {
		return
	}
	doer(&(*arr)[mapped])
}

func (this *DeltaSlice[K, T]) Commit() *DeltaSlice[K, T] {
	this.readonlyValues = append(this.readonlyValues, this.appended...)
	this.appended = this.appended[:0]

	elements := this.modified.Elements()
	for i := 0; i < len(elements); i++ {
		idx := common.ToType[K, int](*elements[i].First)
		this.readonlyValues[idx] = elements[i].Second
	}

	slice.RemoveIf(&this.readonlyValues, func(_ int, v T) bool { return this.isDeleted(v) })
	this.modified.Clear()
	return this
}
