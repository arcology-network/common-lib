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

package indexer

import (
	"github.com/arcology-network/common-lib/exp/array"
	btree "github.com/google/btree"
)

type sortable[T any] struct {
	v       T
	compare *func(T, T) bool
}

func newSortable[T any](v T, compare *func(T, T) bool) *sortable[T] {
	return &sortable[T]{
		v:       v,
		compare: compare,
	}
}

func (this *sortable[T]) Less(other btree.Item) bool {
	return (*this.compare)(this.v, other.(*sortable[T]).v)
}

type Index[T any] struct {
	Name      string
	indexTree *btree.BTree
	compare   func(T, T) bool
}

func NewIndex[T any](name string, compare func(T, T) bool) *Index[T] {
	return &Index[T]{
		Name:      name,
		indexTree: btree.New(4),
		compare:   compare,
	}
}

// Add new values to the index.
func (this *Index[T]) Add(vals []T) {
	sortables := array.Append(vals, func(_ int, v T) *sortable[T] { return newSortable[T](v, &this.compare) })
	for _, v := range sortables {
		this.indexTree.ReplaceOrInsert(v)
	}
}

// Remove the values from the index.
func (this *Index[T]) Remove(vals []T) {
	sortables := array.Append(vals, func(_ int, v T) *sortable[T] { return newSortable[T](v, &this.compare) })
	for _, v := range sortables {
		this.indexTree.Delete(v)
	}
}

// Export the index, return a slice of all the values in the index.
func (this *Index[T]) Export() []T {
	vals := make([]T, 0, this.indexTree.Len())
	this.indexTree.Ascend(func(node btree.Item) bool {
		vals = append(vals, node.(*sortable[T]).v)
		return true
	})
	return vals
}

// Clear the index, return the number of items cleared.
func (this *Index[T]) Clear() int {
	size := this.indexTree.Len()
	this.indexTree.Clear(false)
	return size
}

func (this *Index[T]) GreaterThan(lower T) []T {
	var got []btree.Item
	this.indexTree.AscendGreaterOrEqual(newSortable[T](lower, &this.compare), func(node btree.Item) bool {
		if len(got) == 0 && !this.compare(lower, node.(*sortable[T]).v) {
			return true
		}

		got = append(got, node)
		return true
	})
	return array.Append(got, func(_ int, v btree.Item) T { return v.(*sortable[T]).v })
}

func (this *Index[T]) GreaterEqualThan(lower T) []T {
	var got []btree.Item
	this.indexTree.AscendGreaterOrEqual(newSortable[T](lower, &this.compare), func(node btree.Item) bool {
		got = append(got, node)
		return true
	})
	return array.Append(got, func(_ int, v btree.Item) T { return v.(*sortable[T]).v })
}

func (this *Index[T]) LessThan(upper T) []T {
	var got []btree.Item
	this.indexTree.DescendLessOrEqual(newSortable[T](upper, &this.compare), func(node btree.Item) bool {
		if len(got) == 0 && !this.compare(node.(*sortable[T]).v, upper) {
			return true
		}

		got = append(got, node)
		return true
	})
	return array.Append(got, func(_ int, v btree.Item) T { return v.(*sortable[T]).v }) // Move the result to a slice.
}

func (this *Index[T]) LessEqualThan(lower T) []T {
	var got []btree.Item
	this.indexTree.DescendLessOrEqual(newSortable[T](lower, &this.compare), func(node btree.Item) bool {
		got = append(got, node)
		return true
	})
	return array.Append(got, func(_ int, v btree.Item) T { return v.(*sortable[T]).v })
}

// This function is inclusive, both lower and upper are included.
func (this *Index[T]) Between(lower, upper T) []T {
	var got []btree.Item
	this.indexTree.AscendGreaterOrEqual(newSortable[T](lower, &this.compare), func(node btree.Item) bool {
		if this.compare(upper, node.(*sortable[T]).v) {
			return false
		}
		got = append(got, node)
		return true
	})
	return array.Append(got, func(_ int, v btree.Item) T { return v.(*sortable[T]).v }) // Move the result to a slice.
}

func (this *Index[T]) Find(v T) (T, bool) {
	var target btree.Item
	this.indexTree.AscendGreaterOrEqual(newSortable[T](v, &this.compare), func(node btree.Item) bool {
		target = node
		return false
	})

	if target == nil {
		return *new(T), false
	}
	return target.(*sortable[T]).v, true
}

func (this *Index[T]) BatchFind(vals []T) ([]T, []bool) {
	targets, founds := make([]T, len(vals)), make([]bool, len(vals))
	for i, v := range vals {
		targets[i], founds[i] = this.Find(v)
	}
	return targets, founds
}
