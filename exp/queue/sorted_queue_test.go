package queue

import (
	"reflect"
	"testing"
)

func TestSortedQueueEnqueueAndSearch(t *testing.T) {
	type job struct {
		priority int
		id       string
	}

	less := func(a, b job) bool {
		if a.priority != b.priority {
			return a.priority < b.priority
		}
		return a.id < b.id
	}

	this := NewSortedQueue(less)
	if idx := this.Enqueue(job{priority: 20, id: "b"}); idx != 0 {
		t.Fatalf("expected first insert at index 0, got %d", idx)
	}
	this.Enqueue(job{priority: 10, id: "a"})
	inserted := this.Enqueue(job{priority: 20, id: "a"})
	this.Enqueue(job{priority: 30, id: "z"})

	if inserted != 1 {
		t.Fatalf("expected duplicate-priority insert at index 1, got %d", inserted)
	}

	want := []job{{priority: 10, id: "a"}, {priority: 20, id: "a"}, {priority: 20, id: "b"}, {priority: 30, id: "z"}}
	if !reflect.DeepEqual(this.ToSlice(), want) {
		t.Fatalf("sorted queue contents mismatch: got %#v want %#v", this.ToSlice(), want)
	}

	probe := job{priority: 20, id: "a"}
	if idx := this.Greater(probe); idx != 2 {
		t.Fatalf("Greater returned %d, want 2", idx)
	}
	if idx := this.GreaterOrEqual(probe); idx != 1 {
		t.Fatalf("GreaterOrEqual returned %d, want 1", idx)
	}
	if idx := this.Less(probe); idx != 0 {
		t.Fatalf("Less returned %d, want 0", idx)
	}
	if idx := this.LessOrEqual(probe); idx != 1 {
		t.Fatalf("LessOrEqual returned %d, want 1", idx)
	}

	if idx, ok := this.Index(probe); !ok || idx != 1 {
		t.Fatalf("Index returned (%d, %v), want (1, true)", idx, ok)
	}

	if idx, ok := this.Index(job{priority: 25, id: "x"}); ok || idx != -1 {
		t.Fatalf("Index returned (%d, %v), want (-1, false)", idx, ok)
	}

	if front, ok := this.Dequeue(); !ok || front != (job{priority: 10, id: "a"}) {
		t.Fatalf("Dequeue returned (%#v, %v), want first sorted element", front, ok)
	}
	if back, ok := this.Back(); !ok || back != (job{priority: 30, id: "z"}) {
		t.Fatalf("Back returned (%#v, %v), want highest sorted element", back, ok)
	}
	if head, ok := this.Peek(); !ok || head != (job{priority: 20, id: "a"}) {
		t.Fatalf("Peek returned (%#v, %v), want next sorted element", head, ok)
	}
	if this.Size() != 3 {
		t.Fatalf("Size returned %d, want 3", this.Size())
	}
	if this.IsEmpty() {
		t.Fatal("sorted queue should not be empty")
	}

	this.Clear()
	if !this.IsEmpty() || this.Size() != 0 {
		t.Fatal("Clear should empty the sorted queue")
	}

	if _, ok := this.Dequeue(); ok {
		t.Fatal("Dequeue on an empty sorted queue should fail")
	}
	if _, ok := this.Peek(); ok {
		t.Fatal("Peek on an empty sorted queue should fail")
	}
	if _, ok := this.Back(); ok {
		t.Fatal("Back on an empty sorted queue should fail")
	}
}

func TestNewSortedQueueWithItems(t *testing.T) {
	this := NewSortedQueueWithItems([]int{4, 2, 5, 1, 3}, func(a, b int) bool { return a < b })
	if !reflect.DeepEqual(this.ToSlice(), []int{1, 2, 3, 4, 5}) {
		t.Fatalf("NewSortedQueueWithItems should sort the input, got %#v", this.ToSlice())
	}

	if idx := this.Greater(3); idx != 3 {
		t.Fatalf("Greater returned %d, want 3", idx)
	}
	if idx := this.GreaterOrEqual(3); idx != 2 {
		t.Fatalf("GreaterOrEqual returned %d, want 2", idx)
	}
	if idx := this.Less(3); idx != 1 {
		t.Fatalf("Less returned %d, want 1", idx)
	}
	if idx := this.LessOrEqual(3); idx != 2 {
		t.Fatalf("LessOrEqual returned %d, want 2", idx)
	}
}
