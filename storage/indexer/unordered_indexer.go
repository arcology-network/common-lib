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

import "github.com/arcology-network/common-lib/exp/orderedmap"

type UnorderedIndexer[K comparable, T, V any] struct {
	*orderedmap.OrderedMap[K, T, V]
	isIndexable func(T) (K, bool) // If the transition is accepted and the key is returned. An index is only supposed to index the accepted transitions.
}

func NewUnorderedIndexer[K comparable, T, V any](
	nilValue V,
	isIndexable func(T) (K, bool),
	init func(K, T) V,
	setter func(K, T, *V)) *UnorderedIndexer[K, T, V] {

	indexer := &UnorderedIndexer[K, T, V]{
		isIndexable: isIndexable,
	}

	indexer.OrderedMap = orderedmap.NewOrderedMap[K, T, V](
		nilValue,
		1024,
		init,
		setter,
		nil,
	)
	return indexer
}

// New creates a new StateCommitter instance.
func (this *UnorderedIndexer[K, T, V]) Import(transitions []T) {
	for _, t := range transitions {
		if k, ok := this.isIndexable(t); ok { // If the transition is indexable by the index.
			this.Set(k, t)
		}
	}
}

func (this *UnorderedIndexer[K, T, V]) ParallelForeachDo(do func(k K, v *V)) {
	this.OrderedMap.ParallelForeachDo(do)
}
func (this *UnorderedIndexer[K, T, V]) ForeachDo(do func(k K, v V)) { this.OrderedMap.ForeachDo(do) }
func (this *UnorderedIndexer[K, T, V]) Clear()                      { this.OrderedMap.Clear() }
