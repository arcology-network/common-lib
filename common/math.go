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

package common

import (
	"math"
	"math/big"
	"reflect"
	"sort"
)

// Min returns the minimum value between two values of type T.
func Min[T ~int8 | ~int32 | ~int | ~int64 | ~uint8 | ~uint32 | ~uint64 | ~float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func MinIf[T any](a, b T, less func(T, T) bool) T {
	if less(a, b) {
		return a
	}
	return b
}

// Max returns the maximum value between two values of type T.
func Max[T ~int8 | ~int32 | ~int | ~int64 | ~uint8 | ~uint32 | ~uint64 | ~float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func DoMax[T any](entries []T, less func(T, T) bool) (T, int) {
	if len(entries) == 0 {
		var zero T
		return zero, -1
	}

	max := entries[0]
	maxIndex := 0
	for i, entry := range entries[1:] {
		if less(max, entry) {
			max = entry
			maxIndex = i + 1
		}
	}
	return max, maxIndex
}

func DoMin[T any](entries []T, less func(T, T) bool) (T, int) {
	if len(entries) == 0 {
		var zero T
		return zero, -1
	}

	min := entries[0]
	minIndex := 0
	for i, entry := range entries[1:] {
		if less(entry, min) {
			min = entry
			minIndex = i + 1
		}
	}
	return min, minIndex
}

func DoMedian[T any](entries []T, less func(T, T) bool) (T, int) {
	if len(entries) == 0 {
		var zero T
		return zero, -1
	}

	type indexedEntry struct {
		value T
		index int
	}

	sorted := make([]indexedEntry, len(entries))
	for i, entry := range entries {
		sorted[i] = indexedEntry{value: entry, index: i}
	}

	sort.SliceStable(sorted, func(i, j int) bool {
		return less(sorted[i].value, sorted[j].value)
	})

	medianPos := (len(sorted) - 1) / 2
	return sorted[medianPos].value, sorted[medianPos].index
}

// IsBetween checks if a value is within the range defined by min and max, inclusive.
func IsBetween[T ~int8 | ~int32 | ~int | ~int64 | ~uint8 | ~uint32 | ~uint64 | ~float64](value, min, max T) bool {
	return value >= min && value <= max
}

// Remainder calculates the remainder of dividing the total sum of the ASCII values of the characters in the key by numShards.
func Remainder(numShards int, key string) int {
	if len(key) == 0 {
		return math.MaxInt
	}

	var total int = 0
	for j := 0; j < len(key); j++ {
		total += int(key[j])
	}
	return total % numShards
}

// IsHex checks if the given byte slice represents a valid hexadecimal string.
func IsHex(bytes []byte) bool {
	if len(bytes)%2 != 0 {
		return false
	}
	for _, c := range bytes {
		if !(('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')) {
			return false
		}
	}
	return true
}

// Equal checks if two values are equal.
// It returns true if the values are equal; otherwise, it returns false.
func Equal[T comparable](lhv, rhv *T, pred func(*T) bool) bool {
	return (lhv == rhv) ||
		((lhv != nil) && (rhv != nil) && (*lhv == *rhv)) ||
		((lhv == nil && pred(rhv)) || (rhv == nil && pred(lhv)))
}

// EqualIf checks if two values are equal based on a given equality function.
// It returns true if the values are equal; otherwise, it returns false.
func EqualIf[T any](lhv, rhv *T, equal func(*T, *T) bool, wildcard func(*T) bool) bool {
	return (lhv == rhv) || ((lhv != nil) && (rhv != nil) && equal(lhv, rhv)) || ((lhv == nil && wildcard(rhv)) || (rhv == nil && wildcard(lhv)))
}

func NumericEqual(got, want any) bool {
	switch g := got.(type) {

	case nil:
		return want == nil

	// ---------- big.Int ----------
	case *big.Int:
		switch w := want.(type) {
		case *big.Int:
			return g.Cmp(w) == 0
		case big.Int:
			return g.Cmp(&w) == 0
		}

	case big.Int:
		switch w := want.(type) {
		case *big.Int:
			return g.Cmp(w) == 0
		case big.Int:
			return g.Cmp(&w) == 0
		}

	// ---------- signed integers ----------
	case int:
		w, ok := want.(int)
		return ok && g == w
	case *int:
		w, ok := want.(*int)
		return ok && w != nil && *g == *w

	case int8:
		w, ok := want.(int8)
		return ok && g == w
	case *int8:
		w, ok := want.(*int8)
		return ok && w != nil && *g == *w

	case int16:
		w, ok := want.(int16)
		return ok && g == w
	case *int16:
		w, ok := want.(*int16)
		return ok && w != nil && *g == *w

	case int32:
		w, ok := want.(int32)
		return ok && g == w
	case *int32:
		w, ok := want.(*int32)
		return ok && w != nil && *g == *w

	case int64:
		w, ok := want.(int64)
		return ok && g == w
	case *int64:
		w, ok := want.(*int64)
		return ok && w != nil && *g == *w

	// ---------- unsigned integers ----------
	case uint:
		w, ok := want.(uint)
		return ok && g == w
	case *uint:
		w, ok := want.(*uint)
		return ok && w != nil && *g == *w

	case uint8:
		w, ok := want.(uint8)
		return ok && g == w
	case *uint8:
		w, ok := want.(*uint8)
		return ok && w != nil && *g == *w

	case uint16:
		w, ok := want.(uint16)
		return ok && g == w
	case *uint16:
		w, ok := want.(*uint16)
		return ok && w != nil && *g == *w

	case uint32:
		w, ok := want.(uint32)
		return ok && g == w
	case *uint32:
		w, ok := want.(*uint32)
		return ok && w != nil && *g == *w

	case uint64:
		w, ok := want.(uint64)
		return ok && g == w
	case *uint64:
		w, ok := want.(*uint64)
		return ok && w != nil && *g == *w

	// ---------- floats ----------
	case float32:
		w, ok := want.(float32)
		return ok && g == w
	case *float32:
		w, ok := want.(*float32)
		return ok && w != nil && *g == *w

	case float64:
		w, ok := want.(float64)
		return ok && g == w
	case *float64:
		w, ok := want.(*float64)
		return ok && w != nil && *g == *w

	// ---------- numeric slices ----------
	case []int:
		w, ok := want.([]int)
		return ok && reflect.DeepEqual(g, w)
	case []*int:
		w, ok := want.([]*int)
		return ok && reflect.DeepEqual(g, w)

	case []int64:
		w, ok := want.([]int64)
		return ok && reflect.DeepEqual(g, w)
	case []*int64:
		w, ok := want.([]*int64)
		return ok && reflect.DeepEqual(g, w)

	case []uint64:
		w, ok := want.([]uint64)
		return ok && reflect.DeepEqual(g, w)
	case []*uint64:
		w, ok := want.([]*uint64)
		return ok && reflect.DeepEqual(g, w)

	case []float32:
		w, ok := want.([]float32)
		return ok && reflect.DeepEqual(g, w)
	case []*float32:
		w, ok := want.([]*float32)
		return ok && reflect.DeepEqual(g, w)

	case []float64:
		w, ok := want.([]float64)
		return ok && reflect.DeepEqual(g, w)
	case []*float64:
		w, ok := want.([]*float64)
		return ok && reflect.DeepEqual(g, w)

	case []*big.Int:
		w, ok := want.([]*big.Int)
		if !ok || len(g) != len(w) {
			return false
		}
		for i := range g {
			if g[i] == nil || w[i] == nil || g[i].Cmp(w[i]) != 0 {
				return false
			}
		}
		return true

	case []big.Int:
		w, ok := want.([]big.Int)
		if !ok || len(g) != len(w) {
			return false
		}
		for i := range g {
			if g[i].Cmp(&w[i]) != 0 {
				return false
			}
		}
		return true
	}
	return reflect.DeepEqual(got, want)
}
