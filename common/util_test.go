package common

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestRemoveNils(t *testing.T) {
	strs := []interface{}{"1", 2, "3", "4"}
	RemoveNils(&strs)
	if len(strs) != 4 && strs[0] != "1" && strs[1] != 2 && strs[2] != "3" && strs[3] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1"}
	RemoveNils(&strs)
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1", nil, "3", "4"}
	RemoveNils(&strs)
	if len(strs) != 3 && strs[0] != "1" && strs[1] != 3 && strs[2] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, "3", "4"}
	RemoveNils(&strs)
	if len(strs) != 2 && strs[0] != "3" && strs[1] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, "4"}
	RemoveNils(&strs)
	if len(strs) != 1 && strs[0] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, nil}
	RemoveNils(&strs)
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, nil}
	RemoveNils(&strs)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, 2}
	RemoveNils(&strs)
	if len(strs) != 2 && strs[0] != 1 && strs[1] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestRemove(t *testing.T) {
	strs := []interface{}{"1", 2, "3", "4"}
	filter := func(v interface{}) bool { return v == nil }
	Remove(&strs, filter)
	if len(strs) != 4 && strs[0] != "1" && strs[1] != 2 && strs[2] != "3" && strs[3] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1"}
	Remove(&strs, filter)
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{"1", nil, "3", "4"}
	Remove(&strs, filter)
	if len(strs) != 3 && strs[0] != "1" && strs[1] != 3 && strs[2] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, "3", "4"}
	Remove(&strs, filter)
	if len(strs) != 2 && strs[0] != "3" && strs[1] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, "4"}
	Remove(&strs, filter)
	if len(strs) != 1 && strs[0] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{nil, nil, nil, nil}
	Remove(&strs, filter)
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, nil}
	Remove(&strs, filter)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []interface{}{1, nil, nil, 2}
	Remove(&strs, filter)
	if len(strs) != 2 && strs[0] != 1 && strs[1] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestRemoveEmptyStrings(t *testing.T) {
	strs := []string{"1", "2", "3", "4"}
	RemoveEmptyStrings(&strs)
	if len(strs) != 4 && strs[0] != "1" && strs[1] != "2" && strs[2] != "3" && strs[3] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"1"}
	RemoveEmptyStrings(&strs)
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"1", "", "3", "4"}
	RemoveEmptyStrings(&strs)
	if len(strs) != 3 && strs[0] != "1" && strs[1] != "3" && strs[2] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"", "", "3", "4"}
	RemoveEmptyStrings(&strs)
	if len(strs) != 2 && strs[0] != "3" && strs[1] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"", "", "", "4"}
	RemoveEmptyStrings(&strs)
	if len(strs) != 1 && strs[0] != "4" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"1", "", "", ""}
	RemoveEmptyStrings(&strs)
	if len(strs) != 1 && strs[0] != "1" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"1", "", "", "2"}
	RemoveEmptyStrings(&strs)
	if len(strs) != 2 && strs[0] != "1" && strs[1] != "2" {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []string{"", "", "", ""}
	RemoveEmptyStrings(&strs)
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestRemoveUint64s(t *testing.T) {
	strs := []uint64{1, 2, 3, 4}
	RemoveUint64(&strs, math.MaxUint64)
	if len(strs) != 4 && strs[0] != 1 && strs[1] != 2 && strs[2] != 3 && strs[3] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1}
	RemoveUint64(&strs, math.MaxUint64)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, math.MaxUint64, 3, 4}
	RemoveUint64(&strs, math.MaxUint64)
	if len(strs) != 3 && strs[0] != 1 && strs[1] != 3 && strs[2] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{math.MaxUint64, math.MaxUint64, 3, 4}
	RemoveUint64(&strs, math.MaxUint64)
	if len(strs) != 2 && strs[0] != 3 && strs[1] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64, 4}
	RemoveUint64(&strs, math.MaxUint64)
	if len(strs) != 1 && strs[0] != 4 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, math.MaxUint64, math.MaxUint64, math.MaxUint64}
	RemoveUint64(&strs, math.MaxUint64)
	if len(strs) != 1 && strs[0] != 1 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{1, math.MaxUint64, math.MaxUint64, 2}
	RemoveUint64(&strs, math.MaxUint64)
	if len(strs) != 2 && strs[0] != 1 && strs[1] != 2 {
		t.Error("Error: Failed to remove nil values !")
	}

	strs = []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64, math.MaxUint64}
	RemoveUint64(&strs, math.MaxUint64)
	if len(strs) != 0 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestUniqueInts(t *testing.T) {
	nums := []int{4, 5, 5, 6, 1, 4, 2, 3, 3}
	pos := UniqueInts(nums)

	if !reflect.DeepEqual(nums[:pos], []int{1, 2, 3, 4, 5, 6}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = []int{4}
	pos = UniqueInts(nums)
	if !reflect.DeepEqual(nums[:pos], []int{4}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = []int{3, 3}
	pos = UniqueInts(nums)
	if !reflect.DeepEqual(nums[:pos], []int{3}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = []int{}
	pos = UniqueInts(nums)
	if !reflect.DeepEqual(nums[:pos], []int{}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = []int{7, 6, 5, 4, 3, 2, 1}
	pos = UniqueInts(nums)
	if !reflect.DeepEqual(nums[:pos], []int{1, 2, 3, 4, 5, 6, 7}) {
		t.Error("Error: Failed to remove nil values !")
	}

	nums = make([]int, 1000000)
	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Intn(5000000)
	}

	t0 := time.Now()
	UniqueInts(nums)
	for k, v := range nums {
		nums[k] = v
	}
	fmt.Println("UniqueInts: ", len(nums), "leafs in ", time.Now().Sub(t0))

	t0 = time.Now()
	m := map[int]bool{}
	for i := 0; i < len(nums); i++ {
		m[nums[i]] = true
	}

	counter := 0
	for k := range m {
		nums[counter] = k
		counter++
	}
	fmt.Println("UniqueMap: ", len(nums), "leafs in ", time.Now().Sub(t0))
}

func TestRemoveIf(t *testing.T) {
	nums := []interface{}{4, 5, 5, 6, 1, 4, 2, 3, 3}
	RemoveIf(&nums, 3, func(lhv interface{}, rhv interface{}) bool { return lhv.(int) == rhv.(int) })
	if !reflect.DeepEqual(nums, []interface{}{4, 5, 5, 6, 1, 4, 2}) {
		t.Error("Error: Failed to remove the target values !")
	}

	RemoveIf(&nums, 5, func(lhv interface{}, rhv interface{}) bool { return lhv.(int) == rhv.(int) })
	if !reflect.DeepEqual(nums, []interface{}{4, 6, 1, 4, 2}) {
		t.Error("Error: Failed to remove the target values !")
	}
}
