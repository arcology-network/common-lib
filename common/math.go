package common

import (
	"math"

	"golang.org/x/exp/constraints"
)

func Min[T ~int8 | ~int32 | ~int | ~int64 | ~uint8 | ~uint32 | ~uint64 | ~float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T ~int8 | ~int32 | ~int | ~int64 | ~uint8 | ~uint32 | ~uint64 | ~float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

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

func Sum[T0 constraints.Integer | float32 | float64 | byte](values []T0) T0 {
	var sum T0 = 0
	for j := 0; j < len(values); j++ {
		sum += (values[j])
	}
	return sum
}

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

// func Accumulate[T0, T1 constraints.Integer | float32 | float64 | byte](values []T0, Type T1) []T1 {
// 	if len(values) == 0 {
// 		return []T1{}
// 	}

// 	summed := make([]T1, len(values))
// 	summed[0] = T1(values[0])
// 	for i := 1; i < len(values); i++ {
// 		summed[i] = summed[i-1] + T1(values[i])
// 	}
// 	return summed
// }
