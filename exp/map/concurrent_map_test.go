/*
 *   Copyright (c) 2024 Arcology Network

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package mapi

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/arcology-network/common-lib/common"
	slice "github.com/arcology-network/common-lib/exp/slice"
)

func TestCcmapBasic(t *testing.T) {
	ccmap := NewConcurrentMap(8, func(v int) bool { return v == -1 }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	ccmap.Set("1", 1)
	ccmap.Set("2", 2)
	ccmap.Set("3", 3)
	ccmap.Set("4", 4)

	if v, ok := ccmap.Get("1"); !ok || v != 1 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("2"); !ok || v != 2 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("3"); !ok || v != 3 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("4"); !ok || v != 4 {
		t.Error("Error: Failed to get")
	}

	ccmap.Set("1", 4)
	ccmap.Set("2", 3)
	ccmap.Set("3", 2)
	ccmap.Set("4", 1)

	if v, ok := ccmap.Get("1"); !ok || v != 4 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("2"); !ok || v != 3 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("3"); !ok || v != 2 {
		t.Error("Error: Failed to get")
	}
	if v, ok := ccmap.Get("4"); !ok || v != 1 {
		t.Error("Error: Failed to get")
	}

	ccmap.Set("3", 3)
	ccmap.Set("4", 4)

	if v, ok := ccmap.Get("3"); !ok || v != 3 {
		t.Error("Error: Failed to get")
	}

	if v, ok := ccmap.Get("4"); !ok || v != 4 {
		t.Error("Error: Failed to get")
	}

	keys := ccmap.Keys()
	sort.SliceStable(keys, func(i, j int) bool {
		return bytes.Compare([]byte(keys[i][:]), []byte(keys[j][:])) < 0
	})

	if v, ok := ccmap.Get("4"); !ok || v != 4 {
		t.Error("Error: Failed to get")
	}

	if !reflect.DeepEqual(keys, []string{"1", "2", "3", "4"}) {
		t.Error("Error: Entries don't match")
	}

	if ccmap.Length() != 4 {
		t.Error("Error: Wrong entry count")
	}
}

func TestCcmapEmptyKeys(t *testing.T) {
	ccmap := NewConcurrentMap(8, func(v int) bool { return v == -1 }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	ccmap.Set("1", 1)
	ccmap.Set("2", 2)
	ccmap.Set("3", 3)
	ccmap.Set("", 4)

	if v, ok := ccmap.Get("1"); !ok || v != 1 {
		t.Error("Error: Failed to get")
	}

	if v, ok := ccmap.Get("2"); !ok || v != 2 {
		t.Error("Error: Failed to get")
	}

	if v, ok := ccmap.Get("3"); !ok || v != 3 {
		t.Error("Error: Failed to get")
	}

	if v, _ := ccmap.Get(""); v != 4 {
		t.Error("Error: Failed to get")
	}

	v, found := ccmap.BatchGet([]string{"1", "2", "3", ""})
	if !found[0] || !reflect.DeepEqual(v, []int{1, 2, 3, 4}) {
		t.Error("Error: Entries don't match")
	}
}

func TestCcmapBatchModeAllEntries(t *testing.T) {
	ccmap := NewConcurrentMap[string, interface{}](8, func(v interface{}) bool { return v == nil }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	keys := []string{"1", "2", "3", "4"}
	values := make([]interface{}, len(keys))
	for i, v := range keys {
		values[i] = v
	}

	ccmap.BatchSet(keys, values)
	outValues := common.First(ccmap.BatchGet(keys))

	if !reflect.DeepEqual(outValues, values) {
		t.Error("Error: Entries don't match")
	}
}

func TestCCmapDump(t *testing.T) {
	m := map[int]int{}
	ky := m[1]
	fmt.Println(ky)

	ccmap := NewConcurrentMap[string, interface{}](8, func(v interface{}) bool { return v == nil }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	keys := []string{"1", "2", "3", "4"}
	values := []interface{}{"1", "2", 3, "4"}

	ccmap.BatchSet(keys, values)
	k, v := ccmap.KVs()
	if !reflect.DeepEqual(k, []string{"1", "2", "3", "4"}) {
		t.Error("Error: Entries don't match")
	}

	if !reflect.DeepEqual(v, []interface{}{"1", "2", 3, "4"}) {
		t.Error("Error: Entries don't match")
	}
}

func TestMinMax(t *testing.T) {
	ccmap := NewConcurrentMap[string, int](8, func(v int) bool { return false }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	keys := []string{"1", "2", "3", "4"}
	values := []int{1, 2, 3, 4}
	ccmap.BatchSet(keys, values)

	minv := math.MaxInt
	less := func(_ string, rhs *int) {
		if minv > *rhs {
			minv = *rhs
		}
	}

	ccmap.Traverse(less)
	if minv != 1 {
		t.Error("Error: Wrong min value")
	}

	maxv := math.MinInt
	greater := func(_ string, v *int) {
		if maxv < *v {
			maxv = *v
		}
	}

	ccmap.Traverse(greater)
	if maxv != 4 {
		t.Error("Error: Wrong max value")
	}
}

func TestForeach(t *testing.T) {
	ccmap := NewConcurrentMap[string, int](8, func(v int) bool { return v < 0 }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	keys := []string{"1", "2", "3", "4"}
	values := []int{1, 2, 3, 4}
	ccmap.BatchSet(keys, values)

	ccmap.Foreach(func(v int) int {
		return v + 10
	})

	_, vs := ccmap.KVs()
	if !reflect.DeepEqual(vs, []int{1 + 10, 2 + 10, 3 + 10, 4 + 10}) {
		t.Error("Error: Checksums don't match")
	}
}

func TestForeachDo(t *testing.T) {
	ccmap := NewConcurrentMap[string, *int](8, func(v *int) bool { return v == nil }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	str0 := 0
	str1 := 1
	str2 := 2
	str3 := 3

	keys := []string{"1", "2", "3", "4"}
	values := []*int{&str0, &str1, &str2, &str3}
	ccmap.BatchSet(keys, values)

	ccmap.ForeachDo(func(k string, v *int) {
		*v += 1
	})

	_, vs := ccmap.KVs()
	if *vs[0] != 1 || *vs[1] != 2 || *vs[2] != 3 || *vs[3] != 4 {
		t.Error("Error: Checksums don't match")
	}

	ccmap.ParallelForeachDo(func(k string, v *int) {
		*v += 1
	})

	if *vs[0] != 2 || *vs[1] != 3 || *vs[2] != 4 || *vs[3] != 5 {
		t.Error("Error: Checksums don't match")
	}
}

func TestParallelDo(t *testing.T) {
	ccmap := NewConcurrentMap[string, interface{}](8, func(v interface{}) bool { return v == nil }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	keys := []string{"1", "2", "3", "4"}
	values := []interface{}{"1", "2", 3, "4"}
	ccmap.BatchSet(keys, values)

	v, found := ccmap.BatchGet([]string{"1", "2", "3", "4"})
	if !found[0] || !reflect.DeepEqual(v, values) {
		t.Error("Error: Entries don't match")
	}

	ccmap.ParallelDo([]string{"1", "2", "3", "5"}, func(i int, k string, v interface{}, found bool) (interface{}, bool) {
		if k == "5" {
			return "5", true
		}
		return nil, false
	})

	if v, ok := ccmap.Get("5"); !ok || v != "5" {
		t.Error("Error: Failed to get")
	}

	ccmap.ParallelDo([]string{"1", "2", "3", "5"}, func(i int, k string, v interface{}, found bool) (interface{}, bool) {
		if k == "5" {
			return "5", true
		}
		return nil, false
	})

	if ccmap.Length() != 5 {
		t.Error("Error: Wrong entry count")
	}

	keys = append(keys, "5")
	ccmap.UnsafeParallelFor(0, len(keys), func(i int) string { return keys[i] }, func(i int, k string, v interface{}, b bool) (interface{}, bool) {
		if k == "5" {
			return "15", true
		}
		return nil, false
	})

	if v, ok := ccmap.Get("5"); !ok || v != "15" {
		t.Error("Error: Failed to get")
	}
}
