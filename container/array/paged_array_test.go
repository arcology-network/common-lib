package array

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/exp/array"
)

func TestPagedArray(t *testing.T) {
	paged := NewPagedArray[int](2, 64, 0) // 2 elements per block, 64 blocks
	paged.Concate([]int{1, 2, 5, 5, 5})
	buf := paged.ToSlice(0, paged.Size())
	if !reflect.DeepEqual(buf, []int{1, 2, 5, 5, 5}) {
		t.Error("Error: Wrong value")
	}

	paged = NewPagedArray[int](2, 100, 0)
	paged.PushBack(1)
	paged.PushBack(2)
	paged.PushBack(3)
	paged.PushBack(4)

	paged.PushBack(5)
	paged.PushBack(6)
	if paged.Get(paged.Size()-1) != 6 {
		t.Error("Error: Wrong value")
	}

	paged.Concate([]int{7, 8})
	values := paged.ToSlice(0, paged.Size())
	if !reflect.DeepEqual(values, []int{1, 2, 3, 4, 5, 6, 7, 8}) {
		t.Error("Error: Wrong value")
	}

	paged.Concate([]int{9, 10})
	values = paged.ToSlice(0, paged.Size())
	if !reflect.DeepEqual(values, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Error("Error: Wrong value")
	}

	paged.Set(4, 9)
	if paged.Get(4) != 9 {
		t.Error("Error: Wrong value")
	}

	paged.PopBack()
	paged.PopBack()
	paged.PopBack()
	paged.PopBack()

	paged.PopBack()
	paged.PopBack()

	if paged.Size() != 4 {
		t.Error("Error: Wrong size")
	}

	if paged.Get(0) != 1 {
		t.Error("Error: Wrong value")
	}

	if paged.Get(1) != 2 {
		t.Error("Error: Wrong value")
	}

	nums := make([]int, 10)
	for i := 0; i < len(nums); i++ {
		nums[i] = i + 10
	}
	paged.Concate(nums)

	paged.Resize(2)
	if paged.Get(0) != 1 {
		t.Error("Error: Wrong value")
	}

	if paged.Get(1) != 2 {
		t.Error("Error: Wrong value")
	}

	paged.Resize(5)
	paged.ForeachBetween(2, paged.Size(), func(_ int, v *int) {
		(*v) = 5
	})

	values = paged.ToSlice(0, paged.Size())
	if !reflect.DeepEqual(values, []int{1, 2, 5, 5, 5}) {
		t.Error("Error: Wrong value")
	}

	buffer := make([]int, 5)
	paged.PopBackToBuffer(buffer)
	if !reflect.DeepEqual(buffer, []int{1, 2, 5, 5, 5}) {
		t.Error("Error: Wrong value")
	}

	paged.Concate([]int{1, 2, 3, 4, 5})
	buffer = buffer[:2]
	paged.PopBackToBuffer(buffer)
	if !reflect.DeepEqual(buffer, []int{4, 5}) {
		t.Error("Error: Wrong value")
	}

	buffer = paged.ToSlice(0, paged.Size())
	if !reflect.DeepEqual(buffer, []int{1, 2, 3}) {
		t.Error("Error: Wrong value")
	}

	buffer = make([]int, 6)
	paged.PopBackToBuffer(buffer)
	if !reflect.DeepEqual(buffer, []int{1, 2, 3, 0, 0, 0}) {
		t.Error("Error: Wrong value")
	}

	if paged.Size() != 0 {
		t.Error("Error: Wrong length")
	}

	cap := paged.Cap()
	paged.Clear()
	if paged.Size() != 0 || paged.Cap() != cap {
		t.Error("Error: Wrong length or capacity")
	}

	paged = NewPagedArray[int](2, 100, 200)
	paged.Foreach(func(_ int, v *int) {
		*v = 111
	})

	idx, _ := array.FindFirstIf(paged.ToSlice(0, paged.Size()), func(v int) bool {
		return v != 111
	})

	if idx != -1 {
		t.Error("Error: Failed to assign value")
	}
}

func TestPagedArrayPreAlloc(t *testing.T) {
	paged := NewPagedArray[int](2, 64, 2) // 2 elements per block, 64 blocks
	paged.Concate([]int{1, 2, 5, 5, 5})
	buf := paged.ToSlice(0, paged.Size())

	if !reflect.DeepEqual(buf, []int{0, 0, 1, 2, 5, 5, 5}) {
		t.Error("Error: Wrong value")
	}

	// paged = NewPagedArray[int](2, 100, 0)
	paged.PushBack(1)
	buf = paged.ToSlice(0, paged.Size())
	if !reflect.DeepEqual(buf, []int{0, 0, 1, 2, 5, 5, 5, 1}) {
		t.Error("Error: Wrong value")
	}
}

func TestCustomType(t *testing.T) {
	type CustomType struct {
		a int
		b [20]byte
		e string
	}

	t0 := time.Now()
	paged := NewPagedArray[CustomType](4096, 1000, 4096*1000)
	paged.ParallelForeach(func(v *CustomType) {
		*v = CustomType{
			a: 1,
			b: [20]byte{1, 2, 3},
			e: "hello",
		}
	})

	fmt.Println("Paged array Initlaized: ", paged.Size(), paged.Cap(), "objects in ", time.Now().Sub(t0))

	t0 = time.Now()
	slice := make([]CustomType, paged.Cap())
	array.ParallelAppend(slice, 4, func(i int, _ CustomType) CustomType {
		return CustomType{
			a: 1,
			b: [20]byte{1, 2, 3},
			e: "hello",
		}
	})

	if paged.Size() != 4096*1000 {
		t.Error("Error: Wrong length")
	}

	fmt.Println("Array Initlaized: ", paged.Cap(), "objects in ", time.Now().Sub(t0))
	paged.PushBack(CustomType{})
	paged.Clear()

	paged.Resize(2)
	if paged.Size() != 2 {
		t.Error("Error: Wrong length")
	}

	paged.Foreach(func(_ int, v *CustomType) {
		v.a = 999
		v.b = [20]byte{3, 2, 1}
		v.e = "hi hello"
	})

	vec := paged.ToSlice(0, paged.Size())
	idx, _ := array.FindFirstIf(vec, func(v CustomType) bool {
		return (v).a != 999 || v.b != [20]byte{3, 2, 1} || v.e != "hi hello"
	})

	if idx != -1 {
		t.Error("Error: Failed to assign value")
	}
}

func BenchmarkTestIntArray(b *testing.B) {
	nums := make([]int, 10000000)
	for i := 0; i < len(nums); i++ {
		nums[i] = i
	}

	t0 := time.Now()
	arr := make([]int, 10000)
	for i := 0; i < 10000000; i++ {
		arr = append(arr, i)
	}
	fmt.Println("slice.append(): ", len(nums), "leafs in ", time.Now().Sub(t0))

	paged := NewPagedArray[int](4096, 100, 0)
	t0 = time.Now()
	for i := 0; i < len(nums); i++ {
		paged.PushBack(i)
	}
	fmt.Println("paged.PushBack(): ", len(nums), "leafs in ", time.Now().Sub(t0))

	t0 = time.Now()
	paged.ForeachBetween(0, paged.Size(), func(_ int, v *int) {
		(*v) = (*v) + 10
	})
	fmt.Println("paged.Foreach(): ", len(nums), "leafs in ", time.Now().Sub(t0))

	t0 = time.Now()
	for i := 0; i < len(nums); i++ {
		paged.PopBack()
	}
	fmt.Println("paged.PopBack(): ", len(nums), "leafs in ", time.Now().Sub(t0))

	array2 := NewPagedArray[int](4096, 100, 0)
	t0 = time.Now()
	array2.Concate(nums)
	fmt.Println("paged.Concate(): ", len(nums), "leafs in ", time.Now().Sub(t0))

	paged.Clear()
	if paged.Size() != 0 || paged.Cap() != 4096*100 {
		b.Error("Error: Wrong length or capacity")
	}
}
