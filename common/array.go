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

// The provided code excerpt is from a file named generic.go in the common package.
// It appears to contain a collection of generic functions that can be used across
// different parts of the project.

package common

// Resize resizes a slice to a new length.
// If the new length is greater than the current length, it appends the required number of elements to the slice.
// If the new length is less than or equal to the current length, it truncates the slice.
func Resize[T any](values []T, newSize int) []T {
	if len(values) >= newSize {
		return values[:newSize]
	}
	return append(values, make([]T, newSize-len(values))...)
}

// Reshape reshapes a slice into a 2D slice with a given number of columns.
func Reshape[T any](bytes []T, columns int) [][]T {
	hashes := make([][]T, len(bytes)/columns)
	for i := range hashes {
		hashes[i] = bytes[i*columns : (i+1)*columns]
	}
	return hashes
}

// MinElement returns the minimum element in a slice, if there are multiple minimum elements, it returns the first one.
func MinElement[T0 any](array []T0, less func(T0, T0) bool) (int, T0) {
	idx := 0
	minv := array[idx]
	for i := idx; i < len(array); i++ {
		if less(array[i], minv) {
			idx = i
			minv = array[i]
		}
	}
	return idx, minv
}

// MaxElement returns the index and the maximum element in a slice. If there are multiple maximum elements, it returns the first one.
func MaxElement[T0 any](array []T0, greater func(T0, T0) bool) (int, T0) {
	idx := 0
	maxv := array[idx]
	for i := idx; i < len(array); i++ {
		if greater(array[i], maxv) {
			idx = i
			maxv = array[i]
		}
	}
	return idx, maxv
}

// Append applies a function to each element in a slice and returns a new slice with the results.
func Append[T any, T1 any](values []T, do func(i int, v T) T1) []T1 {
	vec := make([]T1, len(values))
	for i := 0; i < len(values); i++ {
		vec[i] = do(i, values[i])
	}
	return vec
}

// TrimTrail removes trailing elements from a slice that are equal to a specified value.
func TrimTrail[T comparable](s []T, zero T) []T {
	i := len(s) - 1
	for i >= 0 && s[i] == zero {
		i--
	}
	return s[:i+1]
}
