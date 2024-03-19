package slice

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/common"
	// array "github.com/arcology-network/common-lib/exp/array"
	// "github.com/HPISTechnologies/common-lib/common"
)

func TestRemoveNils(t *testing.T) {
	encoded := [][]byte{{1}, {1}, {3}, {2}, nil}
	RemoveIf(&encoded, func(_ int, v []byte) bool { return v == nil })
	if len(encoded) != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	encoded = make([][]byte, 3)
	encoded[0] = []byte{1}
	encoded[1] = []byte{2}

	RemoveIf(&encoded, func(_ int, v []byte) bool { return v == nil })
	if len(encoded) != 2 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs := []interface{}{"1", 2, "3", "4"}
	RemoveIf(&strs, func(_ int, v interface{}) bool { return v == nil })

	if len(strs) != 4 && strs[0] != "1" && strs[1] != 2 && strs[2] != "3" && strs[3] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1"}
	RemoveIf(&strs, func(_ int, v interface{}) bool { return v == nil })
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1", nil, "3", "4"}
	RemoveIf(&strs, func(_ int, v interface{}) bool { return v == nil })
	if len(strs) != 3 && strs[0] != "1" && strs[1] != 3 && strs[2] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, "3", "4"}
	RemoveIf(&strs, func(_ int, v interface{}) bool { return v == nil })
	if len(strs) != 2 && strs[0] != "3" && strs[1] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, "4"}
	RemoveIf(&strs, func(_ int, v interface{}) bool { return v == nil })
	if len(strs) != 1 && strs[0] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, nil}
	RemoveIf(&strs, func(_ int, v interface{}) bool { return v == nil })
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, nil}
	RemoveIf(&strs, func(_ int, v interface{}) bool { return v == nil })
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, 2}
	RemoveIf(&strs, func(_ int, v interface{}) bool { return v == nil })
	if len(strs) != 2 && strs[0] != 1 && strs[1] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestRemoveIf(t *testing.T) {
	strs := []interface{}{"1", 2, "3", "4"}
	filter := func(_ int, v interface{}) bool { return v == nil }
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

	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Intn(100)
	}

	t0 = time.Now()
	UniqueInts(nums)
	fmt.Println("Unique: ", 1000000, " entries in:", time.Now().Sub(t0))

	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Intn(100)
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
	fmt.Println("UniqueInts: ", len(nums), "in ", time.Now().Sub(t0))

	for i := 0; i < 1000000; i++ {
		nums[i] = rand.Intn(5000000)
	}

	m := map[int]bool{}
	t0 = time.Now()
	for i := 0; i < len(nums); i++ {
		m[nums[i]] = true
	}
	common.MapKeys(m)
	fmt.Println("UniqueMap: ", len(nums), "in ", time.Now().Sub(t0))
}

func TestUnique(t *testing.T) {
	nums := []int{4, 5, 5, 6, 1, 3, 4, 2, 3}
	nums = Unique(nums, func(lhv, rhv int) bool { return lhv < rhv })

	if !reflect.DeepEqual(nums, []int{1, 2, 3, 4, 5, 6}) {
		t.Error("Error: Failed to remove nil values !")
	}

	strs := []string{"4", "5", "5", "6", "1", "3", "4", "2", "3"}
	strs = Unique(strs, func(lhv, rhv string) bool { return lhv < rhv })

	if !reflect.DeepEqual(strs, []string{"1", "2", "3", "4", "5", "6"}) {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestForeach(t *testing.T) {
	nums := [][]int{{4}, {5}, {5}, {6}}
	Foreach(nums, func(_ int, lhv *[]int) { (*lhv)[0] += 1 })

	if nums[0][0] != 5 || nums[1][0] != 6 || nums[2][0] != 6 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestParallelForeach(t *testing.T) {
	nums := []int{3, 5, 5, 6, 6}
	ParallelForeach(nums, 120, func(_ int, lhv *int) { (*lhv) = (*lhv) + 1 })

	if nums[0] != 4 || nums[1] != 6 || nums[2] != 6 || nums[3] != 7 || nums[4] != 7 {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = make([]int, 1000000)
	for i := 0; i < len(nums); i++ {
		nums[i] = i
	}

	t0 := time.Now()
	ParallelForeach(nums, 32, func(i int, v *int) {
		sha256.Sum256([]byte(strconv.Itoa(*v)))
	})
	fmt.Println("Time: ", time.Since(t0))
}

func TestFindLastIf(t *testing.T) {
	nums := []int{4, '/', 5, '/', 6}

	idx, _ := FindLastIf(nums, func(v int) bool { return v == '/' })
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

	idx, _ = FindLastIf(charArr, func(v byte) bool { return v == '/' })
	if idx != 3 {
		t.Error("Error: FindLastIf() Failed")
	}
}

func TestEqual(t *testing.T) {
	array0 := []int{1, 2, 3}
	array1 := []int{1, 2, 3}

	if !EqualSet(array0, array1) {
		t.Error("Error: Not equal")
	}

	array0 = []int{}
	array1 = []int{}
	if !EqualSet(array0, array1) {
		t.Error("Error: Not equal")
	}

	array0 = []int{1, 1, 2, 3}
	array1 = []int{1, 2, 3}
	if EqualSet(array0, array1) {
		t.Error("Error: Not equal")
	}

	array0 = []int{1, 1, 3}
	array1 = []int{1, 2, 3}
	if EqualSet(array0, array1) {
		t.Error("Error: Not equal")
	}

	if EqualSet(array0, nil) {
		t.Error("Error: Not equal")
	}

	if EqualSet(nil, array0) {
		t.Error("Error: Not equal")
	}
}

// type testStruct struct {
// 	num int
// }

func TestArrayEqual(t *testing.T) {
	v0 := uint64(1)
	v1 := uint64(1)
	if !common.Equal(&v0, &v1, func(v *uint64) bool { return *v == 1 }) {
		t.Error("Error: Not equal")
	}

	if !common.Equal(&v0, &v1, func(v *uint64) bool { return *v == 10 }) {
		t.Error("Error: Not equal")
	}

	v0 = uint64(1)
	v1 = uint64(2)

	if common.Equal(&v0, &v1, func(v *uint64) bool { return *v == 1 }) {
		t.Error("Error: Should not be equal")
	}

	if !common.Equal(&v0, nil, func(v *uint64) bool { return *v == 1 }) {
		t.Error("Error: Should not be equal")
	}

	if !common.Equal(nil, &v1, func(v *uint64) bool { return *v == 2 }) {
		t.Error("Error: Should not be equal")
	}

	if !common.Equal(nil, nil, func(v *uint64) bool { return *v == 2 }) {
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
	if Count[interface{}, int](v, nil) != 1 {
		t.Error("Error: Not equal")
	}
}

func TestAppend(t *testing.T) {
	src := []int{4, 2, 6, 3, 1}

	target := Transform(src, func(_ int, v int) int { return v + 1 })
	if !reflect.DeepEqual(target, []int{5, 3, 7, 4, 2}) {
		t.Error("Expected: ", target)
	}

	target = ParallelTransform(target, 4, func(i int, v int) int { return target[i] + 1 })
	if !reflect.DeepEqual(target, []int{6, 4, 8, 5, 3}) {
		t.Error("Expected: ", target)
	}
}

func TestInsert(t *testing.T) {
	src := []int{4, 2, 6, 3, 1}

	Insert(&src, 1, int(10))
	if !EqualSet(src, []int{4, 10, 2, 6, 3, 1}) {
		t.Error("Expected: ", "{4, 10, 2, 6, 3, 1}", "actual: ", src)
	}

	Insert(&src, 0, int(10))
	if !EqualSet(src, []int{10, 4, 10, 2, 6, 3, 1}) {
		t.Error("Expected: ", "{10, 4, 10, 2, 6, 3, 1}", "actual: ", src)
	}

	Insert(&src, 7, int(11))
	if !EqualSet(src, []int{10, 4, 10, 2, 6, 3, 1, 11}) {
		t.Error("Expected: ", "{10, 4, 10, 2, 6, 3, 1, 11}", "actual: ", src)
	}

	Insert(&src, 9, int(11))
	if !EqualSet(src, []int{10, 4, 10, 2, 6, 3, 1, 11}) {
		t.Error("Expected: ", "{10, 4, 10, 2, 6, 3, 1, 11}", "actual: ", src)
	}
}

func TestReferenceAndDereference(t *testing.T) {
	src := []int{4, 2, 6, 3, 1}

	refs := Reference(src)
	target := Dereference(refs)

	if !reflect.DeepEqual(target, src) {
		t.Error("Expected: ", target)
	}
}

func TestMinMaxElem(t *testing.T) {
	src := []int{4, 2, 86, 3, 1}
	minidx, min := Min(src, func(a, b int) bool { return a < b })
	if min != 1 || minidx != 4 {
		t.Error("Expected: ", min, minidx)
	}

	maxIdx, max := Max(src, func(a, b int) bool { return a > b })
	if max != 86 || maxIdx != 2 {
		t.Error("Expected: ", min, minidx)
	}
}

func TestSum(t *testing.T) {
	v := []uint{1, 2, 3, 4}
	if Sum[uint, uint](v) != 10 {
		t.Error("Error: Should be 10")
	}

	bytes := []byte{1, 2, 3, 4}
	if Sum[byte, uint](bytes) != 10 {
		t.Error("Error: Should be 10")
	}
}

func TestMoveIf(t *testing.T) {
	strs := []interface{}{"1", 2, "3", "4"}
	filter := func(_ int, v interface{}) bool { return v == 2 }
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

	filter = func(_ int, v interface{}) bool { return v == nil }
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

	strs = []interface{}{"1"}

	moved = MoveIf(&strs, func(_ int, v any) bool {
		return v == "1"
	})

	if len(moved) != 1 || moved[0] != "1" || len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestGroupBy(t *testing.T) {
	strs := []string{"1-1", "2-2", "13", "2-4"}

	keys, v := GroupBy(strs, func(_ int, v string) *string {
		str := (v[0:1])
		return &str
	})

	if len(keys) != 2 || len(v) != 2 || len(v[0]) != 2 || len(v[1]) != 2 {
		t.Error("Error: Failed to remove nil values !")
	}
}

// func TestGroupIndicesBy(t *testing.T) {
// 	keys := []string{"1", "2", "1", "2"}
// 	vals := []string{"1-1", "2-2", "13", "2-4"}
// 	indices, v := GroupIndicesBy(keys, func(i int, v string) *string {
// 		str := vals[i]
// 		return &str
// 	})

// 	if keys[indices[0]] != "1-1" || keys[indices[1]] != "13" || len(indices) != 2 || v != 2 {
// 		t.Error("Error: Failed to remove nil values !")
// 	}
// }

func TestConcate(t *testing.T) {
	strs := []string{"1-1", "2-2", "1-3", "2-4"}

	strVec := Concate(strs, func(v string) []string {
		return []string{v[0:1], v[2:]}
	})

	if len(strVec) != 8 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestJoin(t *testing.T) {
	target := Join([]string{}, []string{"1", "2"}, []string{}, []string{"3", "4"}, []string{})
	if !reflect.DeepEqual(target, []string{"1", "2", "3", "4"}) {
		t.Error("Error: should be equal", target)
	}

	buffer := Join([][]byte{}, [][]byte{{1}, {2}}, [][]byte{{3}, {4}}, [][]byte{})
	if !reflect.DeepEqual(buffer, [][]byte{{1}, {2}, {3}, {4}}) {
		t.Error("Error: should be equal", buffer)
	}

	buffer = Join([][]byte{}, [][]byte{})
	if !reflect.DeepEqual(buffer, [][]byte{}) {
		t.Error("Error: should be equal", buffer)
	}
}

func TestRemoveBothIf(t *testing.T) {
	first := []string{"1", "2", "3", "4"}
	second := []int{1, 2, 3, 4}

	RemoveBothIf(&first, &second, func(_ int, v string, _ int) bool { return v == "1" })
	if len(first) != 3 || len(second) != 3 || first[0] != "2" || second[0] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestResize(t *testing.T) {
	first := []string{"1", "2", "3", "4"}
	Resize(&first, 2)

	if len(first) != 2 || first[0] != "1" || first[1] != "2" {
		t.Error("Error: Failed to remove nil values !")
	}

	Resize(&first, 4)
	if len(first) != 4 || first[0] != "1" || first[1] != "2" || first[2] != "" || first[3] != "" {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestSortBy1st(t *testing.T) {
	nums := []int{4, 5, 3, 2}
	first := []string{"1", "2", "3", "4"}

	SortBy1st(nums, first, func(a, b int) bool { return a < b })
	if !reflect.DeepEqual(first, []string{"4", "3", "1", "2"}) {
		t.Error("Error: Failed to remove nil values !")
	}
}

func BenchmarkTestUniqueInts(t *testing.B) {
	t0 := time.Now()
	arr := NewDo(1000000, func(i int) int { return rand.Int() })
	UniqueInts(arr)
	fmt.Println("UniqueInts: ", 1000000, " entries in:", time.Now().Sub(t0))

	t0 = time.Now()
	arr = NewDo(1000000, func(i int) int { return rand.Int() })
	m := map[int]bool{}
	for i := 0; i < len(arr); i++ {
		m[arr[i]] = true
	}
	fmt.Println("Map unique: ", 1000000, " entries in:", time.Now().Sub(t0))
}

func BenchmarkGroupBy(t *testing.B) {
	randKeys := make([]string, 1000000)
	values := make([]int, len(randKeys))
	for i := 0; i < len(randKeys); i++ {
		randKeys[i] = strconv.Itoa(rand.Intn(len(randKeys) / 4))
		values[i] = rand.Intn(1000000)
	}
	t0 := time.Now()
	k, v := GroupBy(randKeys, func(i int, v string) *string { return &randKeys[i] }, 4)
	fmt.Println("GroupBy: ", len(randKeys), " entries in :", len(k), len(v), time.Now().Sub(t0))

	t0 = time.Now()
	sort.Slice(randKeys, func(i, j int) bool {
		return randKeys[i] < randKeys[j]
	})
	fmt.Println("Unique: ", len(randKeys), " entries in:", time.Now().Sub(t0))
}

// func BenchmarkGroupINdicesBy(t *testing.B) {
// 	randKeys := make([]string, 1000000)
// 	values := make([]int, 1000000)
// 	for i := 0; i < 1000000; i++ {
// 		randKeys[i] = strconv.Itoa(rand.Intn(1000 / 4))
// 		values[i] = rand.Intn(1000000)
// 	}
// 	t0 := time.Now()
// 	GroupIndicesBy(randKeys, func(i int, v string) *int { return &values[i] })
// 	fmt.Println("GroupBy: ", len(randKeys), " entries in:", time.Now().Sub(t0))
// }
