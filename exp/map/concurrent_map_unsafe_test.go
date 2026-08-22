/*
 *   Copyright (c) 2026 Arcology Network
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 */

package mapi

import (
	"reflect"
	"sync"
	"testing"
)

func newUnsafeTestMap() *ConcurrentMap[int, int] {
	return NewConcurrentMap(6, func(value int) bool { return value < 0 }, func(key int) uint64 {
		return uint64(key)
	})
}

func TestConcurrentMapUnsafeCRUD(t *testing.T) {
	m := newUnsafeTestMap()
	if err := m.UnsafeSet(2, 20, nil); err != nil {
		t.Fatalf("UnsafeSet without lock returned %v", err)
	}
	if value, found := m.UnsafeGet(2, nil); !found || value != 20 {
		t.Fatalf("UnsafeGet without lock = (%d, %t), want (20, true)", value, found)
	}
	if err := m.UnsafeSet(5, 50, m.LockShard); err != nil {
		t.Fatalf("UnsafeSet with lock returned %v", err)
	}
	if value, found := m.UnsafeGet(5, m.LockReadShard); !found || value != 50 {
		t.Fatalf("UnsafeGet with read lock = (%d, %t), want (50, true)", value, found)
	}
	m.UnsafeDelete(5, m.LockShard)
	if _, found := m.UnsafeGet(5, m.LockReadShard); found {
		t.Fatal("UnsafeDelete did not remove key")
	}
}

func TestUnsafeBatchOperations(t *testing.T) {
	m := newUnsafeTestMap()
	keys := []int{14, 5, 2, 11, 8}
	values := []int{140, 50, 20, 110, 80}
	if got := m.UnsafeBatchSet(keys, values, nil); !reflect.DeepEqual(got, []bool{true, true, true, true, true}) {
		t.Fatalf("UnsafeBatchSet results = %v", got)
	}

	query := []int{8, 99, 14, 5, 2, 11}
	wantValues := []int{80, 0, 140, 50, 20, 110}
	wantFound := []bool{true, false, true, true, true, true}
	for name, lock := range map[string]func(uint64) func(){"unlocked": nil, "read locked": m.LockReadShard} {
		t.Run(name, func(t *testing.T) {
			gotValues, gotFound := m.UnsafeBatchGet(query, lock)
			if !reflect.DeepEqual(gotValues, wantValues) || !reflect.DeepEqual(gotFound, wantFound) {
				t.Fatalf("UnsafeBatchGet = (%v, %v), want (%v, %v)", gotValues, gotFound, wantValues, wantFound)
			}
		})
	}

	m.UnsafeBatchSet(keys, []int{1, 2}, m.LockShard)
	gotValues, gotFound := m.UnsafeBatchGet(keys, m.LockReadShard)
	if !reflect.DeepEqual(gotValues, []int{1, 2, 0, 0, 0}) || !reflect.DeepEqual(gotFound, []bool{true, true, false, false, false}) {
		t.Fatalf("short-value batch = (%v, %v)", gotValues, gotFound)
	}

	m.UnsafeBatchDelete([]int{5, 99}, m.LockShard)
	_, gotFound = m.UnsafeBatchGet(keys, nil)
	if !reflect.DeepEqual(gotFound, []bool{true, false, false, false, false}) {
		t.Fatalf("found after delete = %v", gotFound)
	}
}

func TestSetBatchAndGetBatch(t *testing.T) {
	m := newUnsafeTestMap()
	keys := []int{14, 5, 2, 11, 8}
	m.SetBatch(keys, []int{1, 2, 3, 4, 5})
	values, found := m.GetBatch(keys)
	if !reflect.DeepEqual(values, []int{1, 2, 3, 4, 5}) || !reflect.DeepEqual(found, []bool{true, true, true, true, true}) {
		t.Fatalf("SetBatch/GetBatch = (%v, %v)", values, found)
	}
}

func TestUnsafeDeleteMissingKeyIsNoOp(t *testing.T) {
	m := newUnsafeTestMap()
	m.UnsafeSet(2, 20, nil)

	m.UnsafeDelete(99, nil)
	if value, found := m.UnsafeGet(2, nil); !found || value != 20 {
		t.Fatalf("deleting missing key changed existing entry: (%d, %t)", value, found)
	}
}

func TestUnsafeBatchDeleteIgnoresMissingKeys(t *testing.T) {
	m := newUnsafeTestMap()
	m.UnsafeBatchSet([]int{2, 5}, []int{20, 50}, nil)

	m.UnsafeBatchDelete([]int{2, 99, 5, 100}, nil)
	_, found := m.UnsafeBatchGet([]int{2, 5, 99, 100}, nil)
	if !reflect.DeepEqual(found, []bool{false, false, false, false}) {
		t.Fatalf("found after batch delete = %v", found)
	}
}

func TestUnsafeBatchSetDuplicateKeysUsesInputOrder(t *testing.T) {
	m := newUnsafeTestMap()
	m.UnsafeBatchSet([]int{2, 2, 2}, []int{10, 20, 30}, nil)

	value, found := m.UnsafeGet(2, nil)
	if !found || value != 30 {
		t.Fatalf("duplicate-key final value = (%d, %t), want (30, true)", value, found)
	}
}

func TestConcurrentMapRejectsZeroShards(t *testing.T) {
	m := NewConcurrentMap(0, func(int) bool { return false }, func(key int) uint64 {
		return uint64(key)
	})
	if m.NumShards() == 0 {
		t.Fatal("NewConcurrentMap created a map with zero shards")
	}
}

func TestUnsafeBatchLockedConcurrentAccess(t *testing.T) {
	m := newUnsafeTestMap()
	keys := []int{0, 1, 2, 3, 4, 5}

	var workers sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		worker := worker
		workers.Add(1)
		go func() {
			defer workers.Done()
			values := make([]int, len(keys))
			for i := range values {
				values[i] = worker*100 + i
			}
			for iteration := 0; iteration < 100; iteration++ {
				m.UnsafeBatchSet(keys, values, m.LockShard)
				m.UnsafeBatchGet(keys, m.LockReadShard)
			}
		}()
	}
	workers.Wait()

	_, found := m.UnsafeBatchGet(keys, m.LockReadShard)
	if !reflect.DeepEqual(found, []bool{true, true, true, true, true, true}) {
		t.Fatalf("found after concurrent access = %v", found)
	}
}
