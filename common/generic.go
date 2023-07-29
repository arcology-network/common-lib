package common

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func Reverse[T any](values *[]T) []T {
	for i, j := 0, len(*values)-1; i < j; i, j = i+1, j-1 {
		(*values)[i], (*values)[j] = (*values)[j], (*values)[i]
	}
	return *values
}

func Fill[T any](values []T, v T) []T {
	for i := 0; i < len(values); i++ {
		(values)[i] = v
	}
	return values
}

func PadRight[T any](values []T, v T, targetLen int) []T {
	if targetLen <= len(values) {
		return values
	}
	return append(values, make([]T, targetLen-len(values))...)
}

func PadLeft[T any](values []T, v T, targetLen int) []T {
	if targetLen < len(values) {
		return values
	}
	return append(make([]T, targetLen-len(values)), values...)
}

func Remove[T comparable](values *[]T, target T) []T {
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
	return (*values)
}

func SetByIndices[T0 any, T1 constraints.Integer](source []T0, indices []T1, setter func(T0) T0) []T0 {
	for _, idx := range indices {
		(source)[idx] = setter((source)[idx])
	}
	return source
}

func RemoveIf[T any](values *[]T, condition func(T) bool) []T {
	MoveIf(values, condition)
	return *values
}

func MoveIf[T any](values *[]T, condition func(T) bool) []T {
	pos := 0
	// for _, condition := range conditions {
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
	moved := (*values)[pos:]
	(*values) = (*values)[:pos]
	return moved
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

func IfThenDo2nd[T any](condition bool, v1 T, f0 func() T) T {
	if condition {
		return f0()
	}
	return v1
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

func IfThenDoEither[T any](condition bool, f0 func() T, f1 func() T) T {
	if condition {
		return f0()
	}
	return f1()
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

func Foreach[T any](values []T, predicate func(v *T)) []T {
	for i := 0; i < len(values); i++ {
		predicate(&(values)[i])
	}
	return values
}

func Accumulate[T any, T1 constraints.Integer | constraints.Float](values []T, initialV T1, predicate func(v T) T1) T1 {
	for i := 0; i < len(values); i++ {
		initialV += predicate((values)[i])
	}
	return initialV
}

func CopyIf[T any](values []T, condition func(v T) bool) []T {
	copied := make([]T, 0, len(values))
	for i := 0; i < len(values); i++ {
		if condition(values[i]) {
			copied = append(copied, values[i])
		}
	}
	return copied
}

func CopyIfDo[T0, T1 any](values []T0, condition func(T0) bool, do func(T0) T1) []T1 {
	copied := make([]T1, 0, len(values))
	for i := 0; i < len(values); i++ {
		if condition(values[i]) {
			copied = append(copied, do(values[i]))
		}
	}
	return copied
}

func UniqueInts[T constraints.Integer](nums []T) []T {
	if len(nums) <= 1 {
		return nums
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
	return nums[:current+1]
}

func Unique[T comparable](src []T, less func(lhv, rhv T) bool) []T {
	if len(src) <= 1 {
		return src
	}

	sort.Slice(src, func(i, j int) bool {
		return less(src[i], src[j])
	})

	current := 0
	for i := 0; i < len(src); i++ {
		if src[current] != (src)[i] {
			src[current+1] = (src)[i]
			current++
		}
	}

	var uniqueElems []T
	UniqueDo(src, less, func(offset int) { uniqueElems = src[:current+1] })
	return uniqueElems
}

func UniqueDo[T comparable](nums []T, less func(lhv, rhv T) bool, do func(int)) {
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
	do(current + 1)
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

func CloneIf[T any](src []T, condition func(v T) bool) []T {
	dst := make([]T, 0, len(src))
	for i := range src {
		if condition(src[i]) {
			dst = append(dst, src[i])
		}
	}
	return dst
}

func Concate[T0, T1 any](array []T0, getter func(T0) []T1) []T1 {
	buffer := make([][]T1, len(array))
	for i := 0; i < len(array); i++ {
		buffer[i] = getter(array[i])
	}

	return Flatten(buffer)
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

func Exclude[T comparable](source []T, toRemove []T) []T {
	dict := MapFromArray(toRemove, true)
	return CopyIf(source, func(v T) bool { return (*dict)[v] })
}

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
	return (lhv == rhv) ||
		((lhv != nil) && (rhv != nil) && (*lhv == *rhv)) ||
		((lhv == nil && wildcard(rhv)) || (rhv == nil && wildcard(lhv)))
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

func IsType[T any](v interface{}) bool {
	switch v.(type) {
	case T:
		return true
	}
	return false
}

func GroupBy[T0 any, T1 comparable](array []T0, getter func(T0) *T1) [][]T0 {
	if len(array) == 1 {
		return [][]T0{array}
	}

	dict := make(map[T1][]T0)
	for _, v := range array {
		if key := getter(v); key != nil {
			vec := dict[*key]
			if vec != nil {
				vec = []T0{}
			}
			dict[*key] = append(vec, v)
		}
	}
	return MapValues(dict)
}

func MapRemoveIf[M ~map[K]V, K comparable, V any](source M, condition func(k K, v V) bool) {
	for k, v := range source {
		if condition(k, v) {
			delete(source, k)
		}
	}
}

func MapMoveIf[M ~map[K]V, K comparable, V any](source, target M, condition func(k K, v V) bool) {
	for k, v := range source {
		if condition(k, v) {
			target[k] = v
			delete(source, k)
		}
	}
}

func MergeMaps[M ~map[K]V, K comparable, V any](from, to M) M {
	for k, v := range to {
		from[k] = v
	}
	return from
}

func MapFromArray[K comparable, V any](keys []K, v V) *map[K]V {
	M := make(map[K]V)
	for _, k := range keys {
		M[k] = v
	}
	return &M
}

func MapFromArrayBy[K comparable, T, V any](keys []T, initv V, getter func(t T) K) *map[K]V {
	M := make(map[K]V)
	for _, k := range keys {
		M[getter(k)] = initv
	}
	return &M
}

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
