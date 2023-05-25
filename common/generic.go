package common

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func Reverse[T any](values *[]T) {
	for i, j := 0, len(*values)-1; i < j; i, j = i+1, j-1 {
		(*values)[i], (*values)[j] = (*values)[j], (*values)[i]
	}
}

func Fill[T any](values []T, v T) []T {
	for i := 0; i < len(values); i++ {
		(values)[i] = v
	}
	return values
}

func Remove[T comparable](values *[]T, target T) {
	pos := 0
	for i := 0; i < len(*values); i++ {
		if target == (*values)[i] {
			pos = i
			break
		}
	}

	for i := pos; i < len(*values); i++ {
		if target != (*values)[i] {
			(*values)[pos], (*values)[i] = (*values)[i], (*values)[pos]
			pos++
		}
	}
	(*values) = (*values)[:pos]
}

// func RemoveIfMeet[T any](values *[]T, args interface{}, conditions ...func(T, ...interface{}) bool) {
// 	if conditions == nil {
// 		return
// 	}

// 	pos := 0
// 	for _, condition := range conditions {
// 		for i := 0; i < len(*values); i++ {
// 			if condition((*values)[i], args) {
// 				pos = i
// 				break
// 			}
// 		}

// 		for i := pos; i < len(*values); i++ {
// 			if !condition((*values)[i], args) {
// 				(*values)[pos], (*values)[i] = (*values)[i], (*values)[pos]
// 				pos++
// 			}
// 		}
// 		(*values) = (*values)[:pos]
// 	}
// }

func SetIndices[T0 any, T1 constraints.Integer](source []T0, indices []T1, setter func(T0) T0) []T0 {
	for _, idx := range indices {
		(source)[idx] = setter((source)[idx])
	}
	return source
}

func RemoveIf[T any](values *[]T, conditions ...func(T) bool) []T {
	if conditions == nil {
		return *values
	}

	pos := 0
	for _, condition := range conditions {
		for i := 0; i < len(*values); i++ {
			if condition((*values)[i]) {
				pos = i
				break
			}
		}

		for i := pos; i < len(*values); i++ {
			if !condition((*values)[i]) {
				(*values)[pos], (*values)[i] = (*values)[i], (*values)[pos]
				pos++
			}
		}
		(*values) = (*values)[:pos]
	}
	return *values
}

func IfThen[T any](condition bool, v0 T, v1 T) T {
	if condition {
		return v0
	}
	return v1
}

func IfThenDo1st[T any](condition bool, f0 func() T, v1 T) T {
	if condition {
		return f0()
	}
	return v1
}

func IfThenDo2nd[T any](condition bool, f0 func() T, v1 T) T {
	if condition {
		return f0()
	}
	return v1
}

func IfThenDoBoth[T any](condition bool, f0 func() T, f1 func() T) T {
	if condition {
		return f0()
	}
	return f1()
}

func IfThenDo(condition bool, f0 func(), f1 func()) {
	if condition && f0 != nil {
		f0()
		return
	}

	if f1 != nil {
		f1()
	}
}

// None nil
func EitherOf[T any](lhv interface{}, rhv T) T {
	if lhv != nil {
		return lhv.(T)
	}
	return rhv
}

func EitherEqualsTo[T any](lhv interface{}, rhv T, equal func(v interface{}) bool) T {
	if equal(lhv) {
		return lhv.(T)
	}
	return rhv
}

func Foreach[T any](values []T, predicate func(v *T)) {
	for i := 0; i < len(values); i++ {
		predicate(&(values)[i])
	}
}

func Accumulate[T any](values []T, initialV uint64, predicate func(v *T) uint64) uint64 {
	for i := 0; i < len(values); i++ {
		initialV += predicate(&(values)[i])
	}
	return initialV
}

func CopyIf[T any](values *[]T, condition func(v T) bool) []T {
	found := []T{}
	for i := 0; i < len(*values); i++ {
		if condition((*values)[i]) {
			found = append(found, (*values)[i])
		}
	}
	return found
}

// func Unique[T comparable](strs *[]T) []T {
// 	dict := make(map[T]bool)
// 	for i := 0; i < len(*strs); i++ {
// 		dict[(*strs)[i]] = true
// 	}

// 	uniques := make([]T, 0, len(dict))
// 	for k := range dict {
// 		uniques = append(uniques, k)
// 	}
// 	return uniques
// }

func UniqueInts[T constraints.Integer](nums []T) int {
	if len(nums) == 0 {
		return 0
	}

	sort.Slice(nums, func(i, j int) bool {
		return (nums[i] < nums[j])
	})

	current := 0
	for i := 0; i < len(nums); i++ {
		if nums[current] != (nums)[i] {
			nums[current+1] = (nums)[i]
			current++
		}
	}
	return current + 1
}

func Unique[T comparable](nums []T, less func(lhv, rhv T) bool) []T {
	if len(nums) <= 1 {
		return nums
	}

	sort.Slice(nums, func(i, j int) bool {
		return less(nums[i], nums[j])
	})

	current := 0
	for i := 0; i < len(nums); i++ {
		if nums[current] != (nums)[i] {
			nums[current+1] = (nums)[i]
			current++
		}
	}
	return nums[:current+1]
}

func FindRange[T comparable](values []T, equal func(v0, v1 T) bool) []int {
	positions := make([]int, 0, len(values))
	positions = append(positions, 0)
	current := values[0]
	for i := 1; i < len(values); i++ {
		if !equal(current, values[i]) {
			current = values[i]
			positions = append(positions, i)
		}
	}
	positions = append(positions, len(values))
	return positions
}

func FindFirst[T comparable](values *[]T, v T) (int, *T) {
	for i := 0; i < len(*values); i++ {
		if (*values)[i] == v {
			return i, &(*values)[i]
		}
	}
	return -1, nil
}

// Find the leftmost index of the element meeting the criteria
func FindFirstIf[T any](values []T, condition func(v T) bool) (int, *T) {
	for i := 0; i < len(values); i++ {
		if condition(values[i]) {
			return i, &(values)[i]
		}
	}
	return -1, nil
}

func LocateFirstIf[T any](values []T, condition func(v T) bool) int {
	for i := 0; i < len(values); i++ {
		if condition(values[i]) {
			return i
		}
	}
	return -1
}

func FindLast[T comparable](values *[]T, v T) (int, *T) {
	for i := len(*values) - 1; i >= 0; i-- {
		if (*values)[i] == v {
			return i, &(*values)[i]
		}
	}
	return -1, nil
}

// Find the rightmost index of the element meeting the criteria
func FindLastIf[T any](values *[]T, condition func(v T) bool) (int, *T) {
	for i := len(*values) - 1; i >= 0; i-- {
		if condition((*values)[i]) {
			return i, &(*values)[i]
		}
	}
	return -1, nil
}

func New[T any](v T) *T {
	v0 := T(v)
	return &v0
}

func Clone[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

func Flatten[T any](src [][]T) []T {
	totalSize := 0
	for _, data := range src {
		totalSize = totalSize + len(data)
	}
	buffer := make([]T, totalSize)
	positions := 0
	for i := range src {
		positions = positions + copy(buffer[positions:], src[i])
	}
	return buffer
}

func SortBy1st[T0 any, T1 any](first []T0, second []T1, compare func(T0, T0) bool) {
	array := make([]struct {
		_0 T0
		_1 T1
	}, len(first))

	for i := range array {
		array[i]._0 = first[i]
		array[i]._1 = second[i]
	}
	sort.SliceStable(array, func(i, j int) bool { return compare(array[i]._0, array[j]._0) })

	for i := range array {
		first[i] = array[i]._0
		second[i] = array[i]._1
	}
}

// func ConcateFrom[T0, T1 any](array []T0, getter func(T0) []T1) []T1 {
// 	total := 0
// 	for i := 0; i < len(array); i++ {
// 		total += len(getter(array[i]))
// 	}
// 	output := make([]T1, total) // Pre-allocation for better performance

// 	offset := 0
// 	for i := 0; i < total; i++ {
// 		elems := getter(array[i])
// 		copy(output[offset:], elems)
// 		offset += len(elems)
// 	}
// 	return output
// }

func CastTo[T0, T1 any](src []T0, predicate func(T0) T1) []T1 {
	target := make([]T1, len(src))
	for i := range src {
		target[i] = predicate(src[i])
	}
	return target
}

func To[T0, T1 any](src []T0) []T1 {
	target := make([]T1, len(src))
	for i := range src {
		target[i] = (interface{}((src[i]))).(T1)
	}
	return target
}

func Equal[T comparable](lhv, rhv *T, wildcard func(*T) bool) bool {
	return (lhv == rhv) || ((lhv != nil) && (rhv != nil) && (*lhv == *rhv)) || ((lhv == nil && wildcard(rhv)) || (rhv == nil && wildcard(lhv)))
}

func EqualIf[T any](lhv, rhv *T, equal func(*T, *T) bool, wildcard func(*T) bool) bool {
	return (lhv == rhv) || ((lhv != nil) && (rhv != nil) && equal(lhv, rhv)) || ((lhv == nil && wildcard(rhv)) || (rhv == nil && wildcard(lhv)))
}

func EqualArray[T comparable](lhv []T, rhv []T) bool {
	if len(lhv) != len(rhv) {
		return false
	}

	for _, v0 := range lhv {
		flag := false
		for _, v1 := range rhv {
			if v0 == v1 {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}

	for _, v0 := range rhv {
		flag := false
		for _, v1 := range lhv {
			if v0 == v1 {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func MergeMaps[M ~map[K]V, K comparable, V any](from, to M) M {
	for k, v := range to {
		from[k] = v
	}
	return from
}

// func MergeMapsIf[M ~map[K]V, K comparable, V any](from, to M, func()) M {
// 	for k, v := range to {
// 		from[k] = v
// 	}
// 	return from
// }

func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

func MapKVs[M ~map[K]V, K comparable, V any](m M) ([]K, []V) {
	keys := make([]K, len(m))
	values := make([]V, len(m))
	i := 0
	for k, v := range m {
		keys[i] = k
		values[i] = v
		i++
	}
	return keys, values
}
