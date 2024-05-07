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

package orderedslice

import (
	"sort"

	slice "github.com/arcology-network/common-lib/exp/slice"
)

type OrderedSlice[T0, T1 any] struct {
	Elements     []T0
	indices      []T1
	getter       func(T0) T1
	greaterEqual func(T1, T1) bool
	cloner       func(T0) T0
}

func NewOrderedSlice[T0 any, T1 comparable](
	preAlloc int,
	getter func(T0) T1,
	greaterEqual func(T1, T1) bool,
	cloner func(T0) T0,
	vals ...T0) *OrderedSlice[T0, T1] {
	orderedSlice := &OrderedSlice[T0, T1]{
		Elements:     make([]T0, 0, preAlloc),
		indices:      []T1{},
		getter:       getter,
		greaterEqual: greaterEqual,
		cloner:       cloner,
	}

	for _, v := range vals {
		orderedSlice.Append(v)
	}
	return orderedSlice
}

func (this *OrderedSlice[T, T1]) Append(v T) *OrderedSlice[T, T1] {
	if this.getter == nil {
		this.Elements = append(this.Elements, v)
	} else {
		idx := this.getter(v)
		nPos := sort.Search(len(this.indices), func(i int) bool {
			return this.greaterEqual(this.indices[i], idx)
		})

		slice.Insert(&this.indices, nPos, idx)
		slice.Insert(&this.Elements, nPos, v)
	}
	return this
}

func (this *OrderedSlice[T, T1]) Clone() *OrderedSlice[T, T1] {
	cloned := &OrderedSlice[T, T1]{
		Elements:     make([]T, len(this.Elements)),
		indices:      make([]T1, len(this.indices)),
		getter:       this.getter,
		greaterEqual: this.greaterEqual,
		cloner:       this.cloner,
	}

	for i, v := range this.Elements {
		cloned.Elements[i] = this.cloner(v)
		cloned.indices[i] = this.getter(v)
	}
	return this
}

func (this *OrderedSlice[T, T1]) Equal(other *OrderedSlice[T, T1]) bool {
	cloned := &OrderedSlice[T, T1]{
		Elements:     make([]T, len(this.Elements)),
		indices:      make([]T1, len(this.indices)),
		getter:       this.getter,
		greaterEqual: this.greaterEqual,
		cloner:       this.cloner,
	}

	for i, v := range this.Elements {
		cloned.Elements[i] = this.cloner(v)
		cloned.indices[i] = this.getter(v)
	}
	return true
}
