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

type Triplet[T0, T1, T2 any] struct {
	First  T0
	Second T1
	Third  T2
}

type Triplets[T0, T1, T2 any] []*Triplet[T0, T1, T2]

// ToTuples converts three arrays into an array of tuples.
// It takes three arrays, arr0, arr1, and arr2, and returns an array of structs,
// where each struct contains the corresponding elements from arr0, arr1, and arr2.
func NewTriplets[T0, T1, T2 any](arr0 []T0, arr1 []T1, arr2 []T2, getter func(int, *T0) T0) Triplets[T0, T1, T2] {
	triplets := make([]*Triplet[T0, T1, T2], len(arr0))
	for i := range arr0 {
		triplets[i] = &Triplet[T0, T1, T2]{
			First:  arr0[i],
			Second: arr1[i],
			Third:  arr2[i],
		}
	}
	return triplets
}

// Firsts extracts the first elements from an array of pairs and returns a new array.
func (this *Triplets[T0, T1, T2]) Array() *[]*Triplet[T0, T1, T2] {
	return (*[]*Triplet[T0, T1, T2])(this)
}

// Firsts extracts the first elements from an array of pairs and returns a new array.
func (this *Triplets[T0, T1, T2]) Firsts() []T0 {
	return array.ParallelAppend(*this, 4, func(i int, triplet *Triplet[T0, T1, T2]) T0 {
		return triplet.First
	})
}

// Seconds extracts the second elements from an array of pairs and returns a new array.
func (this *Triplets[T0, T1, T2]) Seconds() []T1 {
	return array.ParallelAppend(*this, 4, func(i int, triplet *Triplet[T0, T1, T2]) T1 {
		return triplet.Second
	})
}

// Seconds extracts the second elements from an array of pairs and returns a new array.
func (this *Triplets[T0, T1, T2]) Thirds() []T2 {
	return array.ParallelAppend(*this, 4, func(i int, triplet *Triplet[T0, T1, T2]) T2 {
		return triplet.Third
	})
}

func (this *Triplets[T0, T1, T2]) Split() ([]T0, []T1, []T2) {
	seconds, thirds := make([]T1, len(*this)), make([]T2, len(*this))
	return array.ParallelAppend(*this, 4, func(i int, triplet *Triplet[T0, T1, T2]) T0 {
		seconds[i] = triplet.Second
		thirds[i] = triplet.Third
		return triplet.First
	}), seconds, thirds
}

// From converts two arrays into an array of pairs.
// It takes two arrays, arr0 and arr1, and returns an array of structs,
// where each struct contains the corresponding elements from arr0 and arr1.
// func (this *Triplets[T0, T1, T2]) From(arr0 []T0, arr1 []T1, arr2 []T2, getter func(int, *T0) T0) *Triplets[T0, T1, T2] {
// 	// (*this) = make([]*Triplet[T0, T1, T2], len(arr0))

// 	if len(arr0) > 8192 {
// 		(*this) = Triplets[T0, T1, T2](array.Append(arr0, func(i int, v T0) *Triplet[T0, T1, T2] {
// 			return &Triplet[T0, T1, T2]{
// 				First:  getter(i, &v),
// 				Second: arr1[i],
// 			}
// 		}))
// 		return this
// 	}

// 	(*this) = Triplets[T0, T1, T2](array.ParallelAppend(arr0, 8, func(i int, v T0) *Triplet[T0, T1, T2] {
// 		return &Triplet[T0, T1, T2]{
// 			First:  getter(i, &v),
// 			Second: arr1[i],
// 		}
// 	}))
// 	// v := Triplets[T0, T1, T2](pairs)
// 	// return &v

// 	// for i := range arr0 {
// 	// 	(*this)[i] = &Triplet[T0, T1, T2]{
// 	// 		First:  getter(i, &arr0[i]),
// 	// 		Second: arr1[i],
// 	// 	}
// 	// }
// 	return this
// }

// From converts two arrays into an array of pairs.
// It takes two arrays, arr0 and arr1, and returns an array of structs,
// where each struct contains the corresponding elements from arr0 and arr1.
// func (this *Triplets[T0, T1, T2]) FromSlice(arr0 []T0, arr1 []T1) *Triplets[T0, T1, T2] {
// 	(*this) = make([]*Triplet[T0, T1, T2], len(arr0))

// 	if len(arr0) > 8192 {
// 		(*this) = Triplets[T0, T1, T2](array.Append(arr0, func(i int, v T0) *Triplet[T0, T1, T2] {
// 			return &Triplet[T0, T1, T2]{
// 				First:  arr0[i],
// 				Second: arr1[i],
// 			}
// 		}))
// 		return this
// 	}

// 	(*this) = Triplets[T0, T1, T2](array.ParallelAppend(arr0, 8, func(i int, v T0) *Triplet[T0, T1, T2] {
// 		return &Triplet[T0, T1, T2]{
// 			First:  arr0[i],
// 			Second: arr1[i],
// 		}
// 	}))
// 	return this
// }

// To converts an array of pairs into two separate arrays.
// It takes an array of structs, where each struct contains two elements,
// and returns two arrays, one containing the first elements and the other containing the second elements.
func (this *Triplets[T0, T1, T2]) To() ([]T0, []T1) {
	arr0, arr1 := make([]T0, len(*this)), make([]T1, len(*this))
	for i, triplet := range *this {
		arr0[i] = triplet.First
		arr1[i] = triplet.Second
	}
	return arr0, arr1
}
