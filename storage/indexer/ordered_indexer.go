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

import slice "github.com/arcology-network/common-lib/exp/slice"

// SortedIndexer is a collection of indexes that need to be updated together,
// it is memory only, and is used to speed up the query process.
//
// It is either used with a database, which is used to store the actual data,
// or used alone as a memory database that supports indexing.
type SortedIndexer[T any] struct {
	dict    map[string]int
	indexes []*SortedIndex[T]
}

// NewTable creates a new table with the given indexes.
func NewSortedIndexer[T any](indice ...*SortedIndex[T]) *SortedIndexer[T] {
	table := &SortedIndexer[T]{
		dict: map[string]int{},
	}
	for _, index := range indice {
		if _, ok := table.dict[index.Name]; !ok {
			table.indexes = append(table.indexes, index)
			table.dict[index.Name] = len(table.dict)
		}
	}
	return table
}

// Update updates all indexes in the table, everytime new records are added.
func (this *SortedIndexer[T]) Update(v []T) {
	slice.ParallelForeach(this.indexes, 4, func(i int, index **SortedIndex[T]) {
		(**index).Add(v)
	})
}

// removeIndex removes all the indices in the table specified by the input values.
func (this *SortedIndexer[T]) Remove(v []T) {
	slice.ParallelForeach(this.indexes, 4, func(i int, index **SortedIndex[T]) {
		(**index).Remove(v)
	})
}

// Column returns the index specified by the column name.
func (this *SortedIndexer[T]) Column(name string) *SortedIndex[T] {
	if loc, ok := this.dict[name]; ok {
		return this.indexes[loc]
	}
	return nil
}
