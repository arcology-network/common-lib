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

import "github.com/arcology-network/common-lib/common"

// Table is a collection of indexes that need to be updated together,
// it is memory only, and is used to speed up the query process.

// It is supposed to be used with a database, which is used to store the actual data.
type Table[T any] struct {
	dict    map[string]int
	indexes []*Index[T]
}

// NewTable creates a new table with the given indexes.
func NewTable[T any](indice ...*Index[T]) *Table[T] {
	table := &Table[T]{
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

// updateIndex updates all indexes in the table, everytime new records are added.
func (this *Table[T]) Update(v []T) {
	common.ParallelForeach(this.indexes, 4, func(index **Index[T], i int) {
		(**index).Add(v)
	})
}

func (this *Table[T]) Column(name string) *Index[T] {
	if loc, ok := this.dict[name]; ok {
		return this.indexes[loc]
	}
	return nil
}
