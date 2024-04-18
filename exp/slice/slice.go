// The provided code excerpt is from a file named generic.go in the common package.
// It appears to contain a collection of generic functions that can be used across
// different parts of the project.

package slice

import (
	"sort"

	"github.com/arcology-network/common-lib/common"
	"golang.org/x/exp/constraints"
)

// NewArray creates a new slice of a given length and initializes all elements with a given value.
func New[T any](length int, v T) []T {
	array := make([]T, length)
	for i := 0; i < len(array); i++ {
		array[i] = v
	}
	return array
}

// ParallelAppend applies a function to each index in a slice in parallel using multiple threads and returns a new slice with the results.
func NewDo[T any](length int, init func(i int) T) []T {
	values := make([]T, length)
	for i := 0; i < length; i++ {
		values[i] = init(i)
	}
	return values
}

// ToSlice converts a list of values to a slice.
func ToSlice[T any](vals ...T) []T {
	return vals
}

// ParallelAppend applies a function to each index in a slice in parallel using multiple threads and returns a new slice with the results.
func ParallelNew[T any](length int, numThd int, init func(i int) T) []T {
	values := make([]T, length)
	ParallelForeach(values, numThd, func(i int, _ *T) {
		values[i] = init(i)
	})
	return values
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
func RemoveIf[T any](values *[]T, condition func(int, T) bool) []T {
	if len(*values) == 0 {
		return *values
	}

	if len(*values) == 1 && condition(0, (*values)[0]) {
		(*values) = (*values)[:0]
		return *values
	}

	MoveIf(values, condition)
	return *values
}

// RemoveIf removes all elements from a slice that satisfy a given condition.
// It modifies the original slice and returns the modified slice.
func RemoveBothIf[T0, T1 any](values *[]T0, others *[]T1, condition func(int, T0, T1) bool) ([]T0, []T1) {
	MoveBothIf(values, others, condition)
	return *values, *others
}

// MoveIf moves all elements from a slice that satisfy a given condition to a new slice.
// It modifies the original slice and returns the moved elements in a new slice.
func MoveBothIf[T0, T1 any](first *[]T0, second *[]T1, condition func(int, T0, T1) bool) ([]T0, []T1) {
	pos := 0
	for i := 0; i < len(*first); i++ {
		if condition(i, (*first)[i], (*second)[i]) {
			pos = i // Get the first position that satisfies the condition
			break
		}
	}

	for i := pos; i < len(*first); i++ {
		if !condition(i, (*first)[i], (*second)[i]) {
			(*first)[pos], (*first)[i] = (*first)[i], (*first)[pos] // Shift the elements to the front
			(*second)[pos], (*second)[i] = (*second)[i], (*second)[pos]
			pos++
		}
	}

	moved := (*first)[pos:]
	(*first) = (*first)[:pos]

	secondMoved := (*second)[pos:]
	(*second) = (*second)[:pos]

	return moved, secondMoved
}

// MoveIf moves all elements from a slice that satisfy a given condition to a new slice.
// It modifies the original slice and returns the moved elements in a new slice.
func MoveIf[T any](values *[]T, condition func(int, T) bool) []T {
	pos := 0
	// for _, condition := range conditions {
	for i := 0; i < len(*values); i++ {
		if condition(i, (*values)[i]) {
			pos = i
			break
		}
	}

	for i := pos; i < len(*values); i++ {
		if !condition(i, (*values)[i]) {
			(*values)[pos], (*values)[i] = (*values)[i], (*values)[pos]
			pos++
		}
	}
	moved := (*values)[pos:]
	(*values) = (*values)[:pos]
	return moved
}

// Accumulate applies a function to each element in a slice and returns the accumulated result.
func Accumulate[T any, T1 constraints.Integer | constraints.Float](values []T, initialV T1, do func(i int, v T) T1) T1 {
	for i := 0; i < len(values); i++ {
		initialV += do(i, (values)[i])
	}
	return initialV
}

// Foreach applies a function to each element in a slice.
// It modifies the original slice and returns the modified slice.
func Foreach[T any](values []T, do func(idx int, v *T)) {
	for i := 0; i < len(values); i++ {
		do(i, &values[i])
	}
}

// ParallelForeach applies a function to each element in a slice in parallel using multiple threads.
func ParallelForeach[T any](values []T, nThds int, do func(int, *T)) {
	processor := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			do(i, &values[i])
		}
	}
	common.ParallelWorker(len(values), nThds, processor)
}

// Append applies a function to each element in a slice and returns a new slice with the results.
// func Append[T any, T1 any](values []T, do func(i int, v T) T1) []T1 {
// 	vec := make([]T1, len(values))
// 	for i := 0; i < len(values); i++ {
// 		vec[i] = do(i, values[i])
// 	}
// 	return vec
// }

// // ParallelAppend applies a function to each index in a slice in parallel using multiple threads
// // and returns a new slice with the results.
// func ParallelAppend[T any, T1 any](values []T, numThd int, do func(i int, v T) T1) []T1 {
// 	appended := make([]T1, len(values))
// 	worker := func(start, end, index int, args ...interface{}) {
// 		for i := start; i < end; i++ {
// 			appended[i] = do(i, values[i])
// 		}
// 	}
// 	common.ParallelWorker(len(values), numThd, worker)
// 	return appended
// }

// Transform applies a function to each element in a slice and returns a new slice with the results.
func Transform[T any, T1 any](values []T, do func(i int, v T) T1) []T1 {
	vec := make([]T1, len(values))
	for i := 0; i < len(values); i++ {
		vec[i] = do(i, values[i])
	}
	return vec
}

// TransformIf applies a function to each element that satisfies a given condition in a slice and returns a new slice with the results.
func TransformIf[T any, T1 any](values []T, do func(i int, v T) (bool, T1)) []T1 {
	vec := make([]T1, 0, len(values))
	for i := 0; i < len(values); i++ {
		if ok, v := do(i, values[i]); ok {
			vec = append(vec, v)
		}
	}
	return vec
}

// ParallelTransform applies a function to each index in a slice in parallel using multiple threads
// and returns a new slice with the results.
func ParallelTransform[T any, T1 any](values []T, numThd int, do func(i int, v T) T1) []T1 {
	appended := make([]T1, len(values))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			appended[i] = do(i, values[i])
		}
	}
	common.ParallelWorker(len(values), numThd, worker)
	return appended
}

// Insert inserts a value at a specific position in a slice.
func Insert[T any](values *[]T, pos int, v T) []T {
	if pos <= len(*values) { // if pos is the last element
		*values = append(*values, v)
		copy((*values)[pos+1:], (*values)[pos:])
		(*values)[pos] = v
	}
	return *values
}

// Insert inserts a value at a specific position in a slice.
func PushFront[T any](v T, values *[]T) []T {
	return Insert(values, 0, v)
}

// Resize resizes a slice to a new length.
// If the new length is greater than the current length, it appends the required number of elements to the slice.
// If the new length is less than or equal to the current length, it truncates the slice.
func Resize[T any](values *[]T, newSize int) []T {
	if len(*values) >= newSize {
		*values = (*values)[:newSize]
	} else {
		*values = append(*values, make([]T, newSize-len(*values))...)
	}
	return *values
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
	// if common.IsType[[]int](nums) {
	// 	sort.Ints(To[T, int](nums))
	// } else {
	sort.Slice(nums, func(i, j int) bool {
		return (nums[i] < nums[j])
	})
	// }

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
func FindLastIf[T any](values []T, condition func(v T) bool) (int, *T) {
	for i := len(values) - 1; i >= 0; i-- {
		if condition((values)[i]) {
			return i, &(values)[i]
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
func Clone[T any](src []T, fun ...func(T) T) []T {
	dst := make([]T, len(src))
	if len(fun) > 0 {
		for i := range src {
			dst[i] = fun[0](src[i])
		}
		return dst
	} else {
		copy(dst, src)
	}
	return dst
}

// CloneIf creates a copy of a slice based on a given condition and returns the copied slice.
func CloneIf[T any](src []T, condition func(v T) bool, fun ...func(T) T) []T {
	dst := make([]T, 0, len(src))
	for i := range src {
		if condition(src[i]) {
			if len(fun) > 0 {
				dst = append(dst, fun[0](src[i]))
				continue
			}
			dst = append(dst, src[i])
		}
	}
	return dst
}

// Concate concatenates multiple slices into a single slice.
// It applies a getter function to each element in the input slice and concatenates the results.
func Concate[T0, T1 any](array []T0, getter func(T0) []T1) []T1 {
	// buffer := make([][]T1, len(array))
	// for i := 0; i < len(array); i++ {
	// 	buffer[i] = getter(array[i])
	// }
	buffer := Transform(array, func(_ int, v T0) []T1 { return getter(v) })
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

// Flatten flattens a slice of slices into a single slice.
func Join[T any](src ...[]T) []T {
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
	dict := make(map[T]bool)
	for _, k := range toRemove {
		dict[k] = true
	}

	// This is low efficient, because it will scan the whole array for each element to be removed.
	return RemoveIf(&source, func(_ int, v T) bool {
		_, ok := (dict)[v]
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

// Sum calculates the sum of all values in the given slice.
func Sum[T0 constraints.Integer | constraints.Float | byte, T1 constraints.Float | constraints.Integer](values []T0) T1 {
	var sum T1 = 0
	for j := 0; j < len(values); j++ {
		sum += T1(values[j])
	}
	return sum
}

// Count counts the number of occurrences of a value in a slice.
func Count[T comparable, T1 constraints.Integer](values []T, target T) T1 {
	total := T1(0)
	for i := 0; i < len(values); i++ {
		if target == values[i] {
			total++
		}
	}
	return total
}

// Count counts the number of occurrences of a value in a slice.
func CountIf[T0 any, T1 constraints.Integer](values []T0, condition func(int, *T0) bool) T1 {
	total := T1(0)
	for i := 0; i < len(values); i++ {
		if condition(i, &values[i]) {
			total++
		}
	}
	return total
}

// Count counts the number of occurrences of a value in a slice.
func CountDo[T0 any, T1 constraints.Integer](values []T0, getter func(int, *T0) T1) T1 {
	total := T1(0)
	for i := 0; i < len(values); i++ {
		total += getter(i, &values[i]) // Call the
	}
	return total
}

// Equal checks if two slices have the same elements, but the order of the elements doesn't matter.
// It returns true if the slices are equal; otherwise, it returns false.
func EqualSet[T comparable](lhv []T, rhv []T) bool {
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

func EqualSetIf[T any](lhv []T, rhv []T, equal func(T, T) bool) bool {
	if len(lhv) != len(rhv) {
		return false
	}

	for _, v0 := range lhv {
		flag := false
		for _, v1 := range rhv {
			if equal(v0, v1) {
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
			if equal(v0, v1) {
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

// SelectN selects the nth element from each slice in the given 2D slice and returns a new 1D slice containing the selected elements.
// The type parameter T represents the type of elements in the slices. The function returns a slice of type T.
func SelectN[T any](vals [][]T, n int) []T {
	selected := make([]T, len(vals))
	for i, v := range vals {
		selected[i] = v[n]
	}
	return selected
}

// GroupBy groups the elements of an array based on a key getter function.
// It takes an array and a getter function that returns a pointer to the key for each element.
// It returns two slices, one containing the unique keys and the other containing the groups of elements
// that have the same key.
func GroupBy[T0 any, T1 comparable](array []T0, getter func(int, T0) *T1, reserved ...int) ([]T1, [][]T0) {
	if len(array) == 1 {
		return []T1{*getter(0, array[0])}, [][]T0{array}
	}

	length := len(reserved)
	if len(reserved) > 0 {
		length = reserved[0]
	}
	inkeys := ParallelTransform(array, 4, func(i int, v T0) *T1 { return getter(i, v) })

	dict := make(map[T1]*[]T0)
	for i, v := range array {
		if key := inkeys[i]; key != nil {
			vec := dict[*key]
			if vec == nil {
				vec = common.New(make([]T0, 0, length))
				dict[*key] = vec
			}
			*vec = append(*vec, v)
		}
	}
	// fmt.Println("range array:", len(array), "in ", time.Since(t0))
	keys := make([]T1, len(dict))
	values := make([][]T0, len(dict))
	i := 0
	for k, v := range dict {
		keys[i] = k
		values[i] = *v
		i++
	}
	// fmt.Println("k, v := range dict :", len(array), "in ", time.Since(t0))
	return keys, values
}

// GroupIndicesBy groups the elements of an array based on a key getter function and returns the group indices.
// func GroupIndicesBy[T0 any, T1 comparable](array []T0, getter func(int, T0) *T1) ([]int, int) {
// 	if len(array) == 1 {
// 		return []int{0}, 1
// 	}

// 	indices := make([]int, len(array))
// 	dict := make(map[T1]int)
// 	for i, v := range array {
// 		if key := getter(i, v); key != nil {
// 			if v, ok := dict[*key]; ok {
// 				indices[i] = v
// 				continue
// 			}
// 			indices[i] = len(dict)
// 			dict[*key] = len(dict)
// 		}
// 	}
// 	return indices, len(dict)
// }

// Reference returns a new slice containing pointers to the elements in the original slice.
func Reference[T any](array []T) []*T {
	return Transform(array, func(i int, v T) *T { return &v })
}

// Dereference returns a new slice containing the values that the pointers in the original slice point to.
func Dereference[T any](array []*T) []T {
	return Transform(array, func(i int, v *T) T { return *v })
}

func Extreme[T0 any](array []T0, compare func(T0, T0) bool) (int, T0) {
	if len(array) == 0 {
		return -1, *new(T0)
	}

	idx := 0
	minv := array[idx]
	for i := idx; i < len(array); i++ {
		if compare(array[i], minv) {
			idx = i
			minv = array[i]
		}
	}
	return idx, minv
}

func Min[T constraints.Float | constraints.Integer](array []T) (int, T) {
	if len(array) == 0 {
		return -1, 0
	}

	idx := 0
	minv := array[idx]
	for i := idx; i < len(array); i++ {
		if array[i] < minv {
			idx = i
			minv = array[i]
		}
	}
	return idx, minv
}

func Max[T constraints.Float | constraints.Integer](array []T) (int, T) {
	if len(array) == 0 {
		return -1, 0
	}

	idx := 0
	maxv := array[idx]
	for i := idx; i < len(array); i++ {
		if array[i] > maxv {
			idx = i
			maxv = array[i]
		}
	}
	return idx, maxv
}
