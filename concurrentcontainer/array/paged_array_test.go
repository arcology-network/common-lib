package pagedarray

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestPagedArray(t *testing.T) {
	array := NewPagedArray(2, 64)
	array.Append([]interface{}{1, 2, 5, 5, 5})
	buf := array.CopyTo(0, array.Size())
	if !reflect.DeepEqual(buf, []interface{}{1, 2, 5, 5, 5}) {
		t.Error("Error: Wrong value")
	}

	array = NewPagedArray(2, 100)
	array.PushBack(1)
	array.PushBack(2)
	array.PushBack(3)
	array.PushBack(4)

	array.PushBack(5)
	array.PushBack(6)
	if array.Get(array.Size()-1).(int) != 6 {
		t.Error("Error: Wrong value")
	}

	array.Append([]interface{}{7, 8})
	values := array.CopyTo(0, array.Size())
	if !reflect.DeepEqual(values, []interface{}{1, 2, 3, 4, 5, 6, 7, 8}) {
		t.Error("Error: Wrong value")
	}

	array.Append([]interface{}{9, 10})
	values = array.CopyTo(0, array.Size())
	if !reflect.DeepEqual(values, []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Error("Error: Wrong value")
	}

	array.Set(4, 9)
	if array.Get(4).(int) != 9 {
		t.Error("Error: Wrong value")
	}

	array.PopBack()
	array.PopBack()
	array.PopBack()
	array.PopBack()

	array.PopBack()
	array.PopBack()

	if array.Size() != 4 {
		t.Error("Error: Wrong size")
	}

	if array.Get(0).(int) != 1 {
		t.Error("Error: Wrong value")
	}

	if array.Get(1).(int) != 2 {
		t.Error("Error: Wrong value")
	}

	nums := make([]interface{}, 10)
	for i := 0; i < len(nums); i++ {
		nums[i] = i + 10
	}
	array.Append(nums)

	array.Resize(2)
	if array.Get(0).(int) != 1 {
		t.Error("Error: Wrong value")
	}

	if array.Get(1).(int) != 2 {
		t.Error("Error: Wrong value")
	}

	if array.Get(2) != nil {
		t.Error("Error: Wrong value")
	}

	array.Resize(5)
	array.Foreach(2, array.Size(), func(v interface{}) {
		*(v.(*interface{})) = 5
	})

	values = array.CopyTo(0, array.Size())
	if !reflect.DeepEqual(values, []interface{}{1, 2, 5, 5, 5}) {
		t.Error("Error: Wrong value")
	}

	buffer := make([]interface{}, 5)
	array.PopBackToBuffer(buffer)
	if !reflect.DeepEqual(buffer, []interface{}{1, 2, 5, 5, 5}) {
		t.Error("Error: Wrong value")
	}

	array.Append([]interface{}{1, 2, 3, 4, 5})
	buffer = buffer[:2]
	array.PopBackToBuffer(buffer)
	if !reflect.DeepEqual(buffer, []interface{}{4, 5}) {
		t.Error("Error: Wrong value")
	}

	buffer = array.CopyTo(0, array.Size())
	if !reflect.DeepEqual(buffer, []interface{}{1, 2, 3}) {
		t.Error("Error: Wrong value")
	}

	buffer = make([]interface{}, 6)
	array.PopBackToBuffer(buffer)
	if !reflect.DeepEqual(buffer, []interface{}{1, 2, 3, nil, nil, nil}) {
		t.Error("Error: Wrong value")
	}

	if array.Size() != 0 {
		t.Error("Error: Wrong length")
	}
}

func BenchmarkTestInterfaceArray(b *testing.B) {
	nums := make([]interface{}, 10000000)
	for i := 0; i < len(nums); i++ {
		nums[i] = i
	}

	t0 := time.Now()
	arr := make([]interface{}, 10000)
	for i := 0; i < 10000000; i++ {
		arr = append(arr, i)
	}
	fmt.Println("slice.append(): ", len(nums), "leafs in ", time.Now().Sub(t0))

	array := NewPagedArray(4096, 100)
	t0 = time.Now()
	for i := 0; i < len(nums); i++ {
		array.PushBack(i)
	}
	fmt.Println("array.PushBack(): ", len(nums), "leafs in ", time.Now().Sub(t0))

	t0 = time.Now()
	array.Foreach(0, array.Size(), func(v interface{}) {
		*(v.(*interface{})) = (*(v.(*interface{}))).(int) + 10
	})
	fmt.Println("array.Foreach(): ", len(nums), "leafs in ", time.Now().Sub(t0))

	t0 = time.Now()
	for i := 0; i < len(nums); i++ {
		array.PopBack()
	}
	fmt.Println("array.PopBack(): ", len(nums), "leafs in ", time.Now().Sub(t0))

	array2 := NewPagedArray(4096, 100)
	t0 = time.Now()
	array2.Append(nums)
	fmt.Println("array.Append(): ", len(nums), "leafs in ", time.Now().Sub(t0))
}
