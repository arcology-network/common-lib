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

import associative "github.com/arcology-network/common-lib/exp/associative"

type IndexedSlice[T0 any, T1 any, K comparable] struct {
	pairs  associative.Pairs[K, T1]
	values T0
	index  map[K]int
	getter func(*T0, int) *T1
	setter func(*T0, *int, T1) int
	length func(*T0) int
}

func NewIndexSlice[T0 any, T1 any, K comparable](
	values T0,
	getter func(*T0, int) *T1,
	setter func(*T0, *int, T1) int,
	length func(*T0) int) *IndexedSlice[T0, T1, K] {
	return &IndexedSlice[T0, T1, K]{
		pairs:  associative.Pairs[K, T1]{},
		values: values,
		index:  make(map[K]int),
		getter: getter,
		setter: setter,
		length: length,
	}
}

func (this *IndexedSlice[T0, T1, K]) GetByKey(key K) *T1 {
	if idx, ok := this.index[key]; ok {
		return this.getter(&this.values, idx)
	}
	return nil
}

func (this *IndexedSlice[T0, T1, K]) SetByKey(key K, v T1) int {
	idx, ok := this.index[key]
	if !ok {
		idx = this.setter(&this.values, nil, v)
		this.index[key] = idx
		return idx
	}
	this.setter(&this.values, &idx, v)
	return idx
}

func (this *IndexedSlice[T0, T1, K]) GetByIndex(index int) *T1 {
	if this.length(&this.values) <= index {
		return nil
	}
	return this.getter(&this.values, index)
}

func (this *IndexedSlice[T0, T1, K]) SetByIndex(index int, v T1) {
	if this.length(&this.values) < index {
		this.setter(&this.values, &index, v)
	}
}
