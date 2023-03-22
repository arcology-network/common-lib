package common

import "math"

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
