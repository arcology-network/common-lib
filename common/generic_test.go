package common

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
	// "github.com/HPISTechnologies/common-lib/common"
)

func TestUniqueSorted(t *testing.T) {
	nums := []int{3, 1, 1, 1, 1, 1, 1, 3, 2}
	nums = Unique(nums, func(lhv, rhv int) bool { return lhv < rhv })
	if !reflect.DeepEqual(nums, []int{1, 2, 3}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = []int{1, 1, 1, 1, 1, 1}
	nums = Unique(nums, func(lhv, rhv int) bool { return lhv < rhv })
	if !reflect.DeepEqual(nums, []int{1}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = make([]int, 1000000)
	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Intn(100)
	}

	t0 := time.Now()
	Unique(nums, func(lhv, rhv int) bool { return lhv < rhv })
	fmt.Println("Unique: ", 1000000, " entries in:", time.Now().Sub(t0))
}

func TestRemoveNils(t *testing.T) {
	encoded := [][]byte{{1}, {1}, {3}, {2}, nil}
	RemoveIf(&encoded, func(v []byte) bool { return v == nil })
	if len(encoded) != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	encoded = make([][]byte, 3)
	encoded[0] = []byte{1}
	encoded[1] = []byte{2}

	RemoveIf(&encoded, func(v []byte) bool { return v == nil })
	if len(encoded) != 2 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs := []interface{}{"1", 2, "3", "4"}
	RemoveIf(&strs, func(v interface{}) bool { return v == nil })

	if len(strs) != 4 && strs[0] != "1" && strs[1] != 2 && strs[2] != "3" && strs[3] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1"}
	RemoveIf(&strs, func(v interface{}) bool { return v == nil })
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1", nil, "3", "4"}
	RemoveIf(&strs, func(v interface{}) bool { return v == nil })
	if len(strs) != 3 && strs[0] != "1" && strs[1] != 3 && strs[2] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, "3", "4"}
	RemoveIf(&strs, func(v interface{}) bool { return v == nil })
	if len(strs) != 2 && strs[0] != "3" && strs[1] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, "4"}
	RemoveIf(&strs, func(v interface{}) bool { return v == nil })
	if len(strs) != 1 && strs[0] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, nil}
	RemoveIf(&strs, func(v interface{}) bool { return v == nil })
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, nil}
	RemoveIf(&strs, func(v interface{}) bool { return v == nil })
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, 2}
	RemoveIf(&strs, func(v interface{}) bool { return v == nil })
	if len(strs) != 2 && strs[0] != 1 && strs[1] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestRemoveIf(t *testing.T) {
	strs := []interface{}{"1", 2, "3", "4"}
	filter := func(v interface{}) bool { return v == nil }
	RemoveIf(&strs, filter)
	if len(strs) != 4 && strs[0] != "1" && strs[1] != 2 && strs[2] != "3" && strs[3] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1"}
	RemoveIf(&strs, filter)
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1", nil, "3", "4"}
	RemoveIf(&strs, filter)
	if len(strs) != 3 && strs[0] != "1" && strs[1] != 3 && strs[2] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, "3", "4"}
	RemoveIf(&strs, filter)
	if len(strs) != 2 && strs[0] != "3" && strs[1] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, "4"}
	RemoveIf(&strs, filter)
	if len(strs) != 1 && strs[0] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, nil}
	RemoveIf(&strs, filter)
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, nil}
	RemoveIf(&strs, filter)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, 2}
	RemoveIf(&strs, filter)
	if len(strs) != 2 && strs[0] != 1 && strs[1] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestMapMoveIf(t *testing.T) {
	m := map[string]bool{
		"1": true,
		"2": false,
		"3": true,
		"4": false,
	}

	MapRemoveIf(m, func(k string, _ bool) bool { return k == "1" })
	if len(m) != 3 {
		t.Error("Error: Failed to remove nil values !")
	}

	target := map[string]bool{}
	MapMoveIf(m, target, func(k string, _ bool) bool { return k == "2" })
	if len(m) != 2 || len(target) != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

}

func TestMoveIf(t *testing.T) {
	strs := []interface{}{"1", 2, "3", "4"}
	filter := func(v interface{}) bool { return v == 2 }
	moved := MoveIf(&strs, filter)

	if len(moved) != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	if len(strs) != 3 || strs[0] != "1" || strs[1] != "3" || strs[2] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1"}
	MoveIf(&strs, filter)
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	filter = func(v interface{}) bool { return v == nil }
	strs = []interface{}{"1", nil, "3", "4"}
	MoveIf(&strs, filter)
	if len(strs) != 3 || strs[0] != "1" || strs[1] != "3" || strs[2] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, "3", "4"}
	moved = MoveIf(&strs, filter)
	if len(strs) != 2 && strs[0] != "3" && strs[1] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	if len(moved) != 2 && moved[0] != nil && moved[1] != nil {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, "4"}
	moved = MoveIf(&strs, filter)
	if len(strs) != 1 && strs[0] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	if len(moved) != 3 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, nil}
	MoveIf(&strs, filter)
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, nil}
	MoveIf(&strs, filter)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, 2}
	MoveIf(&strs, filter)
	if len(strs) != 2 && strs[0] != 1 && strs[1] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestRemoveEmptyStrings(t *testing.T) {
	strs := []string{"1", "2", "3", "4"}
	Remove(&strs, "")
	if len(strs) != 4 && strs[0] != "1" && strs[1] != "2" && strs[2] != "3" && strs[3] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"1"}
	Remove(&strs, "")
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"1", "", "3", "4"}
	Remove(&strs, "")
	if len(strs) != 3 && strs[0] != "1" && strs[1] != "3" && strs[2] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"", "", "3", "4"}
	Remove(&strs, "")
	if len(strs) != 2 && strs[0] != "3" && strs[1] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"", "", "", "4"}
	Remove(&strs, "")
	if len(strs) != 1 && strs[0] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"1", "", "", ""}
	Remove(&strs, "")
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"1", "", "", "2"}
	Remove(&strs, "")
	if len(strs) != 2 && strs[0] != "1" && strs[1] != "2" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"", "", "", ""}
	Remove(&strs, "")
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestRemoveUint64s(t *testing.T) {
	strs := []uint64{1, 2, 3, 4}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 4 && strs[0] != 1 && strs[1] != 2 && strs[2] != 3 && strs[3] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, math.MaxUint64, 3, 4}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 3 && strs[0] != 1 && strs[1] != 3 && strs[2] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{math.MaxUint64, math.MaxUint64, 3, 4}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 2 && strs[0] != 3 && strs[1] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64, 4}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 1 && strs[0] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, math.MaxUint64, math.MaxUint64, math.MaxUint64}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, math.MaxUint64, math.MaxUint64, 2}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 2 && strs[0] != 1 && strs[1] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestRemove(t *testing.T) {
	strs := []uint64{1, 2, 3, 4}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 4 && strs[0] != 1 && strs[1] != 2 && strs[2] != 3 && strs[3] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, math.MaxUint64, 3, 4}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 3 && strs[0] != 1 && strs[1] != 3 && strs[2] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{math.MaxUint64, math.MaxUint64, 3, 4}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 2 && strs[0] != 3 && strs[1] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64, 4}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 1 && strs[0] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, math.MaxUint64, math.MaxUint64, math.MaxUint64}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, math.MaxUint64, math.MaxUint64, 2}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 2 && strs[0] != 1 && strs[1] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64}
	Remove(&strs, math.MaxUint64)
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, 2, 3, 4}
	RemoveAt(&strs, 1)
	if !reflect.DeepEqual(strs, []uint64{1, 3, 4}) {
		t.Error("Error: Failed to remove nil values !")
	}

	RemoveAt(&strs, 2)
	if !reflect.DeepEqual(strs, []uint64{1, 3}) {
		t.Error("Error: Failed to remove nil values !")
	}

	RemoveAt(&strs, 0)
	if !reflect.DeepEqual(strs, []uint64{3}) {
		t.Error("Error: Failed to remove nil values !")
	}

	RemoveAt(&strs, 0)
	if !reflect.DeepEqual(strs, []uint64{}) {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestUniqueInts(t *testing.T) {
	nums := []int{4, 5, 5, 6, 1, 4, 2, 3, 3}
	nums = UniqueInts(nums)

	if !reflect.DeepEqual(nums, []int{1, 2, 3, 4, 5, 6}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = []int{4}
	nums = UniqueInts(nums)
	if !reflect.DeepEqual(nums, []int{4}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = []int{3, 3}
	nums = UniqueInts(nums)
	if !reflect.DeepEqual(nums, []int{3}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = []int{}
	nums = UniqueInts(nums)
	if !reflect.DeepEqual(nums, []int{}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = []int{7, 6, 5, 4, 3, 2, 1}
	nums = UniqueInts(nums)
	if !reflect.DeepEqual(nums, []int{1, 2, 3, 4, 5, 6, 7}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = make([]int, 1000000)
	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Intn(5000000)
	}

	t0 := time.Now()
	UniqueInts(nums)
	fmt.Println("UniqueInts: ", len(nums), "leafs in ", time.Now().Sub(t0))

	for i := 0; i < 1000000; i++ {
		nums[i] = rand.Intn(5000000)
	}

	m := map[int]bool{}
	t0 = time.Now()
	for i := 0; i < len(nums); i++ {
		m[nums[i]] = true
	}
	MapKeys(m)
	fmt.Println("UniqueMap: ", len(nums), "leafs in ", time.Now().Sub(t0))
}

func TestForeach(t *testing.T) {
	// nums := [][]int{{4}, {5}, {5}, {6}}
	// Foreach(nums, func(lhv *[]int) { (*lhv)[0] += 1 })

	// if nums[0][0] != 5 || nums[1][0] != 6 || nums[2][0] != 6 {
	// 	t.Error("Error: Failed to remove nil values !")
	// }
}

func TestParallelForeach(t *testing.T) {
	// nums := []int{3, 5, 5, 6, 6}
	// ParallelForeach(nums, 120, func(lhv *int) int { return (*lhv) + 1 })

	// if nums[0] != 4 || nums[1] != 6 || nums[2] != 6 || nums[3] != 7 || nums[4] != 7 {
	// 	t.Error("Error: Failed to remove nil values !")
	// }

	// nums = make([]int, 1000000)
	// for i := 0; i < len(nums); i++ {
	// 	nums[i] = i
	// }

	// t0 := time.Now()
	// ParallelForeach(nums, 32, func(v *int) int {
	// 	sha256.Sum256([]byte(strconv.Itoa(*v)))
	// 	return *v
	// })
	// fmt.Println("Time: ", time.Since(t0))
}

func TestFindLastIf(t *testing.T) {
	nums := []int{4, '/', 5, '/', 6}

	idx, _ := FindLastIf(&nums, func(v int) bool { return v == '/' })
	if idx != 3 {
		t.Error("Error: Failed to remove nil values !")
	}

	idx, _ = FindLast(&nums, '/')
	if idx != 3 {
		t.Error("Error: Failed to remove nil values !")
	}

	idx, _ = FindFirstIf(nums, func(v int) bool { return v == '/' })
	if idx != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	idx, _ = FindFirst(nums, '/')
	if idx != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	str := "4/5/6"
	charArr := []byte(str)

	idx, _ = FindLastIf(&charArr, func(v byte) bool { return v == '/' })
	if idx != 3 {
		t.Error("Error: FindLastIf() Failed")
	}
}

func TestMapKeys(t *testing.T) {
	_map := map[uint32]int{}
	_map[11] = 99
	_map[21] = 25

	keys := MapKeys(_map)
	if len(keys) != 2 || (keys[0] != 11 && keys[0] != 21) {
		t.Error("Error: Not equal")
	}
}

func TestValues(t *testing.T) {
	_map := map[uint32]int{}
	_map[11] = 99
	_map[21] = 25

	keys := MapValues(_map)
	if keys[0] != 99 || keys[1] != 25 {
		t.Error("Error: Not equal")
	}
}

func TestEqualArray(t *testing.T) {
	array0 := []int{1, 2, 3}
	array1 := []int{1, 2, 3}

	if !EqualArray(array0, array1) {
		t.Error("Error: Not equal")
	}

	array0 = []int{}
	array1 = []int{}
	if !EqualArray(array0, array1) {
		t.Error("Error: Not equal")
	}

	array0 = []int{1, 1, 2, 3}
	array1 = []int{1, 2, 3}
	if EqualArray(array0, array1) {
		t.Error("Error: Not equal")
	}

	array0 = []int{1, 1, 3}
	array1 = []int{1, 2, 3}
	if EqualArray(array0, array1) {
		t.Error("Error: Not equal")
	}

	if EqualArray(array0, nil) {
		t.Error("Error: Not equal")
	}

	if EqualArray(nil, array0) {
		t.Error("Error: Not equal")
	}
}

// type testStruct struct {
// 	num int
// }

func TestEqual(t *testing.T) {
	v0 := uint64(1)
	v1 := uint64(1)
	if !Equal(&v0, &v1, func(v *uint64) bool { return *v == 1 }) {
		t.Error("Error: Not equal")
	}

	if !Equal(&v0, &v1, func(v *uint64) bool { return *v == 10 }) {
		t.Error("Error: Not equal")
	}

	v0 = uint64(1)
	v1 = uint64(2)

	if Equal(&v0, &v1, func(v *uint64) bool { return *v == 1 }) {
		t.Error("Error: Should not be equal")
	}

	if !Equal(&v0, nil, func(v *uint64) bool { return *v == 1 }) {
		t.Error("Error: Should not be equal")
	}

	if !Equal(nil, &v1, func(v *uint64) bool { return *v == 2 }) {
		t.Error("Error: Should not be equal")
	}

	if !Equal(nil, nil, func(v *uint64) bool { return *v == 2 }) {
		t.Error("Error: Should not be equal")
	}
}

func TestReorderBy(t *testing.T) {
	src := []int{4, 2, 6, 3, 1}
	indices := []int{4, 1, 3, 0, 2}
	reordered := ReorderBy(src, indices)

	if !reflect.DeepEqual(reordered, []int{1, 2, 3, 4, 6}) {
		t.Error("Error: Not equal")
	}
}

func TestAnyIs(t *testing.T) {
	var e []byte
	v := []interface{}{1, 2, e, nil}
	if Count(v, nil) != 1 {
		t.Error("Error: Not equal")
	}
}

func TestAppend(t *testing.T) {
	src := []int{4, 2, 6, 3, 1}

	target := Append(src, func(v int) int { return v + 1 })
	if !reflect.DeepEqual(target, []int{5, 3, 7, 4, 2}) {
		t.Error("Expected: ", target)
	}

	// target = ParallelAppend(target, func(i int) int { return target[i] + 1 })
	// if !reflect.DeepEqual(target, []int{6, 4, 8, 5, 3}) {
	// 	t.Error("Expected: ", target)
	// }
}
