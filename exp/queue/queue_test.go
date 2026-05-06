package queue

import (
	"reflect"
	"testing"
)

func TestQueueOperations(t *testing.T) {
	this := NewQueue[int]()
	if !this.IsEmpty() || this.Size() != 0 {
		t.Error("Error: new queue should be empty")
	}

	this.Enqueue(2)
	this.Enqueue(3)
	this.Enqueue(1)

	if value, ok := this.Peek(); !ok || value != 2 {
		t.Error("Error: Peek should return the front element")
	}

	if value, ok := this.Back(); !ok || value != 1 {
		t.Error("Error: Back should return the last element")
	}

	if value, ok := this.Dequeue(); !ok || value != 2 {
		t.Error("Error: Dequeue should return the first enqueued element")
	}

	if !reflect.DeepEqual(this.ToSlice(), []int{3, 1}) {
		t.Error("Error: queue contents should shift after Dequeue")
	}

	this.Clear()
	if !this.IsEmpty() || this.Size() != 0 {
		t.Error("Error: Clear should empty the queue")
	}

	if _, ok := this.Dequeue(); ok {
		t.Error("Error: Dequeue on an empty queue should fail")
	}
	if _, ok := this.Peek(); ok {
		t.Error("Error: Peek on an empty queue should fail")
	}
	if _, ok := this.Back(); ok {
		t.Error("Error: Back on an empty queue should fail")
	}
}

func TestQueueSortingAndClone(t *testing.T) {
	this := NewSortedQueueFromSlice([]int{4, 1, 3, 2}, func(a, b int) bool { return a < b })
	if !reflect.DeepEqual(this.ToSlice(), []int{1, 2, 3, 4}) {
		t.Error("Error: NewSortedQueueFromSlice should sort the items")
	}

	this.RemoveIf(func(_ int, item int) bool { return item%2 == 0 })
	if !reflect.DeepEqual(this.ToSlice(), []int{1, 3}) {
		t.Error("Error: RemoveIf should remove matching items")
	}

	clone := this.Clone()
	this.Enqueue(5)
	if !reflect.DeepEqual(clone.ToSlice(), []int{1, 3}) {
		t.Error("Error: Clone should produce an independent queue copy")
	}

	fromSlice := NewQueueFromSlice([]int{5, 6})
	cloned := fromSlice.CloneDo(func(v int) int { return v * 10 })
	if !reflect.DeepEqual(cloned.ToSlice(), []int{50, 60}) {
		t.Error("Error: CloneDo should transform each copied element")
	}
}