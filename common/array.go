// The provided code excerpt is from a file named generic.go in the common package.
// It appears to contain a collection of generic functions that can be used across
// different parts of the project.

package common

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// NewArray creates a new slice of a given length and initializes all elements with a given value.
func NewArray[T any](length int, v T) []T {
	array := make([]T, length)
	for i := 0; i < len(array); i++ {
		array[i] = v
	}
	return array
}

// It modifies the original slice and returns the reversed slice.
func Reverse[T any](values *[]T) []T {
	for i, j := 0, len(*values)-1; i < j; i, j = i+1, j-1 {
		(*values)[i], (*values)[j] = (*values)[j], (*values)[i]
	}
	return *values
}

// Fill fills a slice with a given value.
// It modifies the original slice and returns the filled slice.
func Fill[T any](values []T, v T) []T {
	for i := 0; i < len(values); i++ {
		(values)[i] = v
	}
	return values
}

// PadRight pads a slice with a given value on the right side to reach a target length.
// If the target length is less than or equal to the length of the original slice, it returns the original slice.
// Otherwise, it appends the required number of elements to the original slice and returns the padded slice.
func PadRight[T any](values []T, v T, targetLen int) []T {
	if targetLen <= len(values) {
		return values
	}
	return append(values, make([]T, targetLen-len(values))...)
}

// PadLeft pads a slice with a given value on the left side to reach a target length.
// If the target length is less than the length of the original slice, it returns the original slice.
// Otherwise, it prepends the required number of elements to the original slice and returns the padded slice.
func PadLeft[T any](values []T, v T, targetLen int) []T {
	if targetLen < len(values) {
		return values
	}
	return append(make([]T, targetLen-len(values)), values...)
}

// Remove removes all occurrences of a target value from a slice.
// It modifies the original slice and returns the modified slice.
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

// RemoveAt removes an element at a specific position from a slice.
// It modifies the original slice and returns the modified slice.
func RemoveAt[T any](values *[]T, pos int) []T {
	for i := pos; i < len(*values)-1; i++ {
		(*values)[i] = (*values)[i+1]
	}
	(*values) = (*values)[:len((*values))-1]
	return (*values)
}

// SetByIndices applies a setter function to elements in a slice at specified indices.
// It modifies the original slice and returns the modified slice.
func SetByIndices[T0 any, T1 constraints.Integer](source []T0, indices []T1, setter func(T0) T0) []T0 {
	for _, idx := range indices {
		(source)[idx] = setter((source)[idx])
	}
	return source
}

// RemoveIf removes all elements from a slice that satisfy a given condition.
// It modifies the original slice and returns the modified slice.
func RemoveIf[T any](values *[]T, condition func(T) bool) []T {
	MoveIf(values, condition)
	return *values
}

// MoveIf moves all elements from a slice that satisfy a given condition to a new slice.
// It modifies the original slice and returns the moved elements in a new slice.
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

// Foreach applies a function to each element in a slice.
// It modifies the original slice and returns the modified slice.
func Foreach[T any](values []T, do func(v *T, idx int)) []T {
	for i := 0; i < len(values); i++ {
		do(&values[i], i)
	}
	return values
}

// ParallelForeach applies a function to each element in a slice in parallel using multiple threads.
func ParallelForeach[T any](values []T, nThds int, do func(*T, int)) {
	processor := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			do(&values[i], i)
		}
	}
	ParallelWorker(len(values), nThds, processor)
}

// Accumulate applies a function to each element in a slice and returns the accumulated result.
func Accumulate[T any, T1 constraints.Integer | constraints.Float](values []T, initialV T1, do func(v T) T1) T1 {
	for i := 0; i < len(values); i++ {
		initialV += do((values)[i])
	}
	return initialV
}

// Append applies a function to each element in a slice and returns a new slice with the results.
func Append[T any, T1 any](values []T, do func(v T) T1) []T1 {
	vec := make([]T1, len(values))
	for i := 0; i < len(values); i++ {
		vec[i] = do(values[i])
	}
	return vec
}

// ParallelAppend applies a function to each index in a slice in parallel using multiple threads and returns a new slice with the results.
func ParallelAppend[T any, T1 any](values []T, numThd int, do func(i int) T1) []T1 {
	appended := make([]T1, len(values))
	encoder := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			appended[i] = do(i)
		}
	}
	ParallelWorker(len(values), numThd, encoder)
	return appended
}

// Resize resizes a slice to a new length.
// If the new length is greater than the current length, it appends the required number of elements to the slice.
// If the new length is less than or equal to the current length, it truncates the slice.
func Resize[T any](values []T, newSize int) []T {
	if len(values) >= newSize {
		return values[:newSize]
	}
	return append(values, make([]T, newSize-len(values))...)
}

// CopyIf copies elements from a slice to a new slice based on a given condition.
func CopyIf[T any](values []T, condition func(v T) bool) []T {
	copied := make([]T, 0, len(values))
	for i := 0; i < len(values); i++ {
		if condition(values[i]) {
			copied = append(copied, values[i])
		}
	}
	return copied
}

// CopyIfDo copies elements from a slice to a new slice based on a given condition and applies a function to each copied element.
func CopyIfDo[T0, T1 any](values []T0, condition func(T0) bool, do func(T0) T1) []T1 {
	copied := make([]T1, 0, len(values))
	for i := 0; i < len(values); i++ {
		if condition(values[i]) {
			copied = append(copied, do(values[i]))
		}
	}
	return copied
}

// UniqueInts removes duplicate elements from a slice of integers.
// It modifies the original slice and returns the modified slice.
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

// Unique removes duplicate elements from a slice of comparable types.
// It modifies the original slice and returns the modified slice.
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

// UniqueDo removes duplicate elements from a slice of comparable types and applies a function to the modified slice.
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

// FindAllIndics finds the indices where a slice of comparable elements changes.
// It returns a slice of indices indicating the positions where the elements change.
func FindAllIndics[T comparable](values []T, equal func(v0, v1 T) bool) []int {
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

// FindFirst finds the first occurrence of a value in a slice and returns its index and a pointer to the value.
func FindFirst[T comparable](values []T, v T) (int, *T) {
	for i := 0; i < len(values); i++ {
		if (values)[i] == v {
			return i, &(values)[i]
		}
	}
	return -1, nil
}

// FindFirstIf finds the first element in a slice that satisfies a given condition and returns its index and a pointer to the element
func FindFirstIf[T any](values []T, condition func(v T) bool) (int, *T) {
	for i := 0; i < len(values); i++ {
		if condition(values[i]) {
			return i, &(values)[i]
		}
	}
	return -1, nil
}

// LocateFirstIf finds the index of the first element in a slice that satisfies a given condition.
// If no element satisfies the condition, it returns -1.
func LocateFirstIf[T any](values []T, condition func(v T) bool) int {
	for i := 0; i < len(values); i++ {
		if condition(values[i]) {
			return i
		}
	}
	return -1
}

// FindLast finds the last occurrence of a value in a slice and returns its index and a pointer to the value.
func FindLast[T comparable](values *[]T, v T) (int, *T) {
	for i := len(*values) - 1; i >= 0; i-- {
		if (*values)[i] == v {
			return i, &(*values)[i]
		}
	}
	return -1, nil
}

// FindLastIf finds the last element in a slice that satisfies a given condition and returns its index and a pointer to the element.
func FindLastIf[T any](values *[]T, condition func(v T) bool) (int, *T) {
	for i := len(*values) - 1; i >= 0; i-- {
		if condition((*values)[i]) {
			return i, &(*values)[i]
		}
	}
	return -1, nil
}

// Contains checks if a slice contains a given value.
// It returns true if the value is found; otherwise, it returns false.
func Contains[T any](values []T, target T, equal func(v0, v1 T) bool) bool {
	for i := 0; i < len(values); i++ {
		if equal(values[i], target) {
			return true
		}
	}
	return false
}

// Clone creates a copy of a slice and returns the copied slice.
func Clone[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

// CloneIf creates a copy of a slice based on a given condition and returns the copied slice.
func CloneIf[T any](src []T, condition func(v T) bool) []T {
	dst := make([]T, 0, len(src))
	for i := range src {
		if condition(src[i]) {
			dst = append(dst, src[i])
		}
	}
	return dst
}

// Concate concatenates multiple slices into a single slice.
// It applies a getter function to each element in the input slice and concatenates the results.
func Concate[T0, T1 any](array []T0, getter func(T0) []T1) []T1 {
	buffer := make([][]T1, len(array))
	for i := 0; i < len(array); i++ {
		buffer[i] = getter(array[i])
	}

	return Flatten(buffer)
}

// ConcateDo concatenates multiple slices into a single slice.
// It applies a sizer function to each element in the input slice to determine the size of the resulting slice.
// It then applies a getter function to each element in the input slice and concatenates the results.
func ConcateDo[T0, T1 any](array []T0, sizer func(T0) uint64, getter func(T0) []T1) []T1 {
	totalSize := uint64(0)
	for i := 0; i < len(array); i++ {
		totalSize += sizer(array[i])
	}

	buffer := make([]T1, totalSize)
	positions := 0
	for i := range array {
		positions += copy(buffer[positions:], getter(array[i]))
	}
	return buffer
}

// ConcateToBuffer concatenates multiple slices into a buffer slice.
// It applies a getter function to each element in the input slice and concatenates the results to the buffer slice.
func ConcateToBuffer[T0, T1 any](array []T0, buffer *[]T1, getter func(T0) []T1) {
	positions := 0
	for i := range array {
		positions += copy((*buffer)[positions:], getter(array[i]))
	}
}

// Flatten flattens a slice of slices into a single slice.
func Flatten[T any](src [][]T) []T {
	totalSize := 0
	for _, data := range src {
		totalSize = totalSize + len(data)
	}
	buffer := make([]T, totalSize)
	positions := 0
	for i := range src {
		positions += copy(buffer[positions:], src[i])
	}
	return buffer
}

// Reshape reshapes a slice into a 2D slice with a given number of columns.
func Reshape[T any](bytes []T, columns int) [][]T {
	hashes := make([][]T, len(bytes)/columns)
	for i := range hashes {
		hashes[i] = bytes[i*columns : (i+1)*columns]
	}
	return hashes
}

// ReorderBy reorders a slice based on a given index slice.
// It returns a new slice with the elements in the original slice reordered according to the index slice.
func ReorderBy[T any, T1 constraints.Integer](src []T, indices []T1) []T {
	reordered := make([]T, len(src))
	for i := range src {
		reordered[i] = src[indices[i]]
	}
	return reordered
}

// SortBy1st sorts two slices based on the values in the first slice.
// It modifies both slices and sorts them in ascending order based on the values in the first slice.
func SortBy1st[T0 any, T1 any](first []T0, second []T1, compare func(T0, T0) bool) {
	array := make([]struct {
		First  T0
		Second T1
	}, len(first))

	for i := range array {
		array[i].First = first[i]
		array[i].Second = second[i]
	}
	sort.SliceStable(array, func(i, j int) bool { return compare(array[i].First, array[j].First) })

	for i := range array {
		first[i] = array[i].First
		second[i] = array[i].Second
	}
}

// Exclude removes elements from a slice that are present in another slice.
// It modifies the original slice and returns the modified slice.
func Exclude[T comparable](source []T, toRemove []T) []T {
	dict := MapFromArray(toRemove, true)
	return RemoveIf(&source, func(v T) bool {
		_, ok := (*dict)[v]
		return ok
	})
}

// To casts each element in a slice to a different type using type assertion.
// It returns a new slice with the casted elements.
func To[T0, T1 any](src []T0) []T1 {
	target := make([]T1, len(src))
	for i := range src {
		target[i] = (interface{}((src[i]))).(T1)
	}
	return target
}

// Count counts the number of occurrences of a value in a slice.
func Count[T comparable](values []T, target T) uint64 {
	total := uint64(0)
	for i := 0; i < len(values); i++ {
		if target == values[i] {
			total++
		}
	}
	return total
}

// EqualArray checks if two slices are equal.
// It returns true if the slices are equal; otherwise, it returns false.
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

// ToPairs converts two arrays into an array of pairs.
// It takes two arrays, arr0 and arr1, and returns an array of structs,
// where each struct contains the corresponding elements from arr0 and arr1.
func ToPairs[T0, T1 any](arr0 []T0, arr1 []T1) []struct {
	First  T0
	Second T1
} {
	pairs := make([]struct {
		First  T0
		Second T1
	}, len(arr0))
	for i := range arr0 {
		pairs[i] = struct {
			First  T0
			Second T1
		}{arr0[i], arr1[i]}
	}
	return pairs
}

// FromPairs converts an array of pairs into two separate arrays.
// It takes an array of structs, where each struct contains two elements,
// and returns two arrays, one containing the first elements and the other containing the second elements.
func FromPairs[T0, T1 any](pairs []struct {
	First  T0
	Second T1
}) ([]T0, []T1) {
	arr0, arr1 := make([]T0, len(pairs)), make([]T1, len(pairs))
	for i, pair := range pairs {
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

// GroupBy groups the elements of an array based on a key getter function.
// It takes an array and a getter function that returns a pointer to the key for each element.
// It returns two slices, one containing the unique keys and the other containing the groups of elements
// that have the same key.
func GroupBy[T0 any, T1 comparable](array []T0, getter func(T0) *T1) ([]T1, [][]T0) {
	if len(array) == 1 {
		return []T1{*getter(array[0])}, [][]T0{array}
	}

	dict := make(map[T1][]T0)
	for _, v := range array {
		if key := getter(v); key != nil {
			vec := dict[*key]
			if vec == nil {
				vec = []T0{}
			}
			dict[*key] = append(vec, v)
		}
	}
	return MapKVs(dict)
}

// GroupIndicesBy groups the elements of an array based on a key getter function and returns the group indices.
func GroupIndicesBy[T0 any, T1 comparable](array []T0, getter func(T0) *T1) ([]int, int) {
	if len(array) == 1 {
		return []int{0}, 1
	}

	indices := make([]int, len(array))
	dict := make(map[T1]int)
	for i, v := range array {
		if key := getter(v); key != nil {
			if v, ok := dict[*key]; ok {
				indices[i] = v
				continue
			}
			indices[i] = len(dict)
			dict[*key] = len(dict)
		}
	}
	return indices, len(dict)
}
