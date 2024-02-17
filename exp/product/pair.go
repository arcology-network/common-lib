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

package product

import "github.com/arcology-network/common-lib/exp/array"

type Pair[T0, T1 any] struct {
	First  T0
	Second T1
}

type Pairs[T0, T1 any] []*Pair[T0, T1]

// Firsts extracts the first elements from an array of pairs and returns a new array.
func (this *Pairs[T0, T1]) Array() *[]*Pair[T0, T1] {
	return (*[]*Pair[T0, T1])(this)
}

// Firsts extracts the first elements from an array of pairs and returns a new array.
func (this *Pairs[T0, T1]) Firsts() []T0 {
	return array.ParallelAppend(*this, 4, func(i int, pair *Pair[T0, T1]) T0 {
		return pair.First
	})
}

// Seconds extracts the second elements from an array of pairs and returns a new array.
func (this *Pairs[T0, T1]) Seconds() []T1 {
	return array.ParallelAppend(*this, 4, func(i int, pair *Pair[T0, T1]) T1 {
		return pair.Second
	})
}

// From converts two arrays into an array of pairs.
// It takes two arrays, arr0 and arr1, and returns an array of structs,
// where each struct contains the corresponding elements from arr0 and arr1.
func (this *Pairs[T0, T1]) From(arr0 []T0, arr1 []T1) *Pairs[T0, T1] {
	(*this) = make([]*Pair[T0, T1], len(arr0))
	for i := range arr0 {
		(*this)[i] = &Pair[T0, T1]{
			First:  arr0[i],
			Second: arr1[i],
		}
	}
	return this
}

// To converts an array of pairs into two separate arrays.
// It takes an array of structs, where each struct contains two elements,
// and returns two arrays, one containing the first elements and the other containing the second elements.
func (this *Pairs[T0, T1]) To() ([]T0, []T1) {
	arr0, arr1 := make([]T0, len(*this)), make([]T1, len(*this))
	for i, pair := range *this {
		arr0[i] = pair.First
		arr1[i] = pair.Second
	}
	return arr0, arr1
}

// ToTuples converts three arrays into an array of tuples.
// It takes three arrays, arr0, arr1, and arr2, and returns an array of structs,
// where each struct contains the corresponding elements from arr0, arr1, and arr2.
func ToTuples[T0, T1, T2 any](arr0 []T0, arr1 []T1, arr2 []T2) []struct {
	First  T0
	Second T1
	Third  T2
} {
	pairs := make([]struct {
		First  T0
		Second T1
		Third  T2
	}, len(arr0))

	for i := range arr0 {
		pairs[i] = struct {
			First  T0
			Second T1
			Third  T2
		}{arr0[i], arr1[i], arr2[i]}
	}
	return pairs
}

// FromTuples converts an array of tuples into three separate arrays.
// It takes an array of structs, where each struct contains three elements,
// and returns three arrays, one containing the first elements, one containing the second elements,
// and one containing the third elements.
func FromTuples[T0, T1, T2 any](tuples []struct {
	First  T0
	Second T1
	Third  T2
}) ([]T0, []T1, []T2) {
	arr0, arr1, arr2 := make([]T0, len(tuples)), make([]T1, len(tuples)), make([]T2, len(tuples))
	for i, pair := range tuples {
		arr0[i] = pair.First
		arr1[i] = pair.Second
		arr2[i] = pair.Third
	}
	return arr0, arr1, arr2
}
