package orderedset

import (
	"fmt"
	"testing"
	"time"

	slice "github.com/arcology-network/common-lib/exp/slice"
)

func TestIndexedSet(t *testing.T) {
	set := NewOrderedSet([]string{"1", "2"})
	set.Insert("3")
	set.Insert("3")
	set.Insert("4")

	if set.Len() != 4 {
		t.Error("Error: Wrong entry count")
	}

	if set.DeleteByIdx(4) {
		t.Error("Error: Wrong entry count")
	}

	if set.Len() != 4 {
		t.Error("Error: Wrong entry count")
	}

	if v := set.IdxOf("1"); v != 0 {
		t.Error("Error: Failed to 1")
	}
	if v := set.IdxOf("2"); v != 1 {
		t.Error("Error: Failed to 2")
	}
	if v := set.IdxOf("3"); v != 2 {
		t.Error("Error: Failed to 3")
	}
	if v := set.IdxOf("4"); v != 3 {
		t.Error("Error: Failed to 4")
	}

	if v := set.KeyAt(0); v != "1" {
		t.Error("Error: Failed to 1")
	}
	if v := set.KeyAt(1); v != "2" {
		t.Error("Error: Failed to 2")
	}
	if v := set.KeyAt(2); v != "3" {
		t.Error("Error: Failed to 3")
	}

	if !set.DeleteByIdx(0) {
		t.Error("Error: Failed to delete by index 0")
	}

	if v := set.KeyAt(0); v != "2" {
		t.Error("Error: Failed to get")
	}
	if v := set.KeyAt(1); v != "3" {
		t.Error("Error: Failed to get")
	}
	if v := set.KeyAt(2); v != "4" {
		t.Error("Error: Failed to get")
	}

	if !set.DeleteByKey("3") {
		t.Error("Error: Failed to delete")
	}

	if set.Len() != 2 {
		t.Error("Error: Wrong entry count")
	}

	if v := set.KeyAt(0); v != "2" {
		t.Error("Error: Failed to get")
	}

	if v := set.KeyAt(1); v != "4" {
		t.Error("Error: Failed to get")
	}

	if set.Len() != 2 {
		t.Error("Error: Wrong entry count")
	}

	if !set.DeleteByIdx(set.Len() - 1) {
		t.Error("Error: Failed to delete by index 0")
	}

	if v := set.KeyAt(0); v != "2" {
		t.Error("Error: Failed to get")
	}

	if len(set.KeyAt(1)) != 0 {
		t.Error("Error: should failed")
	}

	if !set.DeleteByIdx(0) {
		t.Error("Error: Failed to delete by index 0")
	}

	if set.Len() != 0 {
		t.Error("Error: Wrong entry count")
	}

	if len(set.KeyAt(0)) != 0 {
		t.Error("Error: should failed")
	}

	if len(set.KeyAt(1)) != 0 {
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

	if set.Len() != 0 {
		t.Error("Error: Wrong entry count")
	}

	for i := 0; i < 100; i++ {
		set.Insert(fmt.Sprint(i))
	}

	for {
		if set.Len() == 0 {
			break
		}

		if !set.DeleteByIdx(set.Len() - 1) {
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

	if set.Len() != 1 {
		t.Error("Error: Wrong entry count")
	}

	set.DeleteByKey("4")
	if set.Len() != 0 {
		t.Error("Error: Wrong entry count")
	}

	var emptySet *OrderedSet
	emptySet.Clone()
}

func TestSetOperations(t *testing.T) {
	set := NewOrderedSet([]string{"1", "2"})
	set.Insert("3")
	set.Insert("3")
	set.Insert("4")

	set1 := NewOrderedSet([]string{"1", "2"})
	set1.Insert("3")
	set1.Insert("3")
	set1.Insert("4")
	set1.Insert("5")

	set.Union(set1)

	if set.Length() != 5 {
		t.Error("Error: Wrong entry count")
	}
}

func TestOrderedSetCodec(t *testing.T) {
	set := NewOrderedSet([]string{"1", "2"})
	set.Insert("3")
	set.Insert("3")
	set.Insert("4")

	buffer := set.Encode()
	out := (&OrderedSet{}).Decode(buffer).(*OrderedSet)

	if !slice.Equal(set.Keys(), out.Keys()) {
		t.Error("Error: Lookup Mismatch")
	}

	// if !common.EqualMap(set.dict, out.dict) {
	// 	t.Error("Error: Dic Mismatch")
	// }
}

func BenchmarkSetInsertion(b *testing.B) {
	set := NewOrderedSet([]string{})
	t0 := time.Now()
	for i := 0; i < 1000000; i++ {
		set.Insert(fmt.Sprint(i))
	}
	fmt.Println("set.Insert "+fmt.Sprint(1000000), " in ", time.Since(t0))

	t0 = time.Now()
	m := make(map[string]int)
	for i := 0; i < 1000000; i++ {
		m[fmt.Sprint(i)] = i
	}
	fmt.Println("golang native map Insert "+fmt.Sprint(1000000), " in ", time.Since(t0))
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

	t0 = time.Now()
	m := make(map[string]int)
	for i := 0; i < 1000000; i++ {
		m[fmt.Sprint(i)] = i
	}
	fmt.Println("golang native map deletion "+fmt.Sprint(1000000), " in ", time.Since(t0))

}

func BenchmarkSetPopBack(b *testing.B) {
	set := NewOrderedSet([]string{})
	for i := 0; i < 1000000; i++ {
		set.Insert(fmt.Sprint(i))
	}

	t0 := time.Now()
	for {
		if set.Len() == 0 {
			break
		}

		if !set.DeleteByIdx(set.Len() - 1) {
			b.Error("Error: Failed to delete")
		}
	}
	fmt.Println("set.Insert "+fmt.Sprint(1000000), " in ", time.Since(t0))
}
