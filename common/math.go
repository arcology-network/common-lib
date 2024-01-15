package common

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Min returns the minimum value between two values of type T.
func Min[T ~int8 | ~int32 | ~int | ~int64 | ~uint8 | ~uint32 | ~uint64 | ~float64](a, b T) T {
	if a < b {
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

// Sum calculates the sum of all values in the given slice.
func Sum[T0 constraints.Integer | constraints.Float | byte, T1 constraints.Float | constraints.Integer](values []T0) T1 {
	var sum T1 = 0
	for j := 0; j < len(values); j++ {
		sum += T1(values[j])
	}
	return sum
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
