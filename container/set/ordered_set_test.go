package orderedset

import (
	"fmt"
	"testing"
	"time"
)

func TestIndexedSet(t *testing.T) {
	set := NewOrderedSet([]string{"1", "2"})
	set.Insert("3")
	set.Insert("3")
	set.Insert("4")

	if set.Size() != 4 {
		t.Error("Error: Wrong entry count")
	}

	if set.DeleteByIdx(4) {
		t.Error("Error: Wrong entry count")
	}

	if set.Size() != 4 {
		t.Error("Error: Wrong entry count")
	}

	if v, ok := set.IdxOf("1"); !ok || v != 0 {
		t.Error("Error: Failed to 1")
	}
	if v, ok := set.IdxOf("2"); !ok || v != 1 {
		t.Error("Error: Failed to 2")
	}
	if v, ok := set.IdxOf("3"); !ok || v != 2 {
		t.Error("Error: Failed to 3")
	}
	if v, ok := set.IdxOf("4"); !ok || v != 3 {
		t.Error("Error: Failed to 4")
	}

	if v, ok := set.KeyOf(0); !ok || v.(string) != "1" {
		t.Error("Error: Failed to 1")
	}
	if v, ok := set.KeyOf(1); !ok || v.(string) != "2" {
		t.Error("Error: Failed to 2")
	}
	if v, ok := set.KeyOf(2); !ok || v.(string) != "3" {
		t.Error("Error: Failed to 3")
	}

	if !set.DeleteByIdx(0) {
		t.Error("Error: Failed to delete by index 0")
	}

	if v, ok := set.KeyOf(0); !ok || v.(string) != "2" {
		t.Error("Error: Failed to get")
	}
	if v, ok := set.KeyOf(1); !ok || v.(string) != "3" {
		t.Error("Error: Failed to get")
	}
	if v, ok := set.KeyOf(2); !ok || v.(string) != "4" {
		t.Error("Error: Failed to get")
	}

	if !set.DeleteByKey("3") {
		t.Error("Error: Failed to delete")
	}

	if set.Size() != 2 {
		t.Error("Error: Wrong entry count")
	}

	if v, ok := set.KeyOf(0); !ok || v.(string) != "2" {
		t.Error("Error: Failed to get")
	}

	if v, ok := set.KeyOf(1); !ok || v.(string) != "4" {
		t.Error("Error: Failed to get")
	}

	if set.Size() != 2 {
		t.Error("Error: Wrong entry count")
	}

	if !set.DeleteByIdx(set.Size() - 1) {
		t.Error("Error: Failed to delete by index 0")
	}

	if v, ok := set.KeyOf(0); !ok || v.(string) != "2" {
		t.Error("Error: Failed to get")
	}

	if _, ok := set.KeyOf(1); ok {
		t.Error("Error: should failed")
	}

	if !set.DeleteByIdx(0) {
		t.Error("Error: Failed to delete by index 0")
	}

	if set.Size() != 0 {
		t.Error("Error: Wrong entry count")
	}

	if _, ok := set.KeyOf(0); ok {
		t.Error("Error: should failed")
	}

	if _, ok := set.IdxOf("1"); ok {
		t.Error("Error: should failed")
	}

	for i := 0; i < 100; i++ {
		set.Insert(fmt.Sprint(i))
	}

	for i := 0; i < 100; i++ {
		if !set.DeleteByIdx(0) {
			t.Error("Error: Failed to delete")
		}
	}

	if set.Size() != 0 {
		t.Error("Error: Wrong entry count")
	}

	for i := 0; i < 100; i++ {
		set.Insert(fmt.Sprint(i))
	}

	for {
		if set.Size() == 0 {
			break
		}

		if !set.DeleteByIdx(set.Size() - 1) {
			t.Error("Error: Failed to delete")
		}
	}

	set.Insert("3")
	set.Insert("3")
	set.Insert("4")

	set.DeleteByKey("3")
	if set.DeleteByKey("3") {
		t.Error("Error: Should fail")
	}

	if set.Size() != 1 {
		t.Error("Error: Wrong entry count")
	}

	set.DeleteByKey("4")
	if set.Size() != 0 {
		t.Error("Error: Wrong entry count")
	}
}

func BenchmarkSetInsertion(b *testing.B) {
	set := NewOrderedSet([]string{})
	t0 := time.Now()
	for i := 0; i < 1000000; i++ {
		set.Insert(fmt.Sprint(i))
	}
	fmt.Println("set.Insert "+fmt.Sprint(1000000), " in ", time.Since(t0))
}

func BenchmarkSetPopFront(b *testing.B) {
	set := NewOrderedSet([]string{})

	for i := 0; i < 10000; i++ {
		set.Insert(fmt.Sprint(i))
	}

	t0 := time.Now()
	for i := 0; i < 10000; i++ {
		set.DeleteByIdx(0)
	}
	fmt.Println("set.Insert "+fmt.Sprint(1000), " in ", time.Since(t0))

}

func BenchmarkSetPopBack(b *testing.B) {
	set := NewOrderedSet([]string{})
	for i := 0; i < 1000000; i++ {
		set.Insert(fmt.Sprint(i))
	}

	t0 := time.Now()
	for {
		if set.Size() == 0 {
			break
		}

		if !set.DeleteByIdx(set.Size() - 1) {
			b.Error("Error: Failed to delete")
		}
	}
	fmt.Println("set.Insert "+fmt.Sprint(1000000), " in ", time.Since(t0))
}
