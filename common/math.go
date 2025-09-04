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
