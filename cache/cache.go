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

package storage

import (
	"fmt"
	"runtime"

	"github.com/arcology-network/common-lib/exp/associative"
	"github.com/arcology-network/common-lib/exp/slice"
)

// ReadCache is a read only cache that is used to store the read values from the storage.
// The cache updates itself when the update is called. The implementation isn't thread safe.
// So, it's the caller's responsibility to ensure that the cache is only accessed by one thread updating it.
// Each entry in the cache holds two values, the first value is the old value, and the second value is the new value.
// The new value will be set to the old value when the Finalize function is called.
type ReadCache[K comparable, T any] struct {
	mapper  func(K) uint64
	cache   []map[K]*associative.Pair[*T, *T]
	dirties []*associative.Pair[uint64, K] // The buffer that holds the keys that are updated in the current cycle.
	// stats     []associative.Pair[K, uint64]
}

func NewReadCache[K comparable, T any](numShards uint64, mapper func(K) uint64) *ReadCache[K, T] {
	newReadCache := &ReadCache[K, T]{
		mapper:  mapper,
		cache:   make([]map[K]*associative.Pair[*T, *T], 2),
		dirties: make([]*associative.Pair[uint64, K], 0, 1024),
		// stats:     make([]associative.Pair[K, uint64], 0, 1024),
	}

	for i := range newReadCache.cache {
		newReadCache.cache[i] = make(map[K]*associative.Pair[*T, *T])
	}
	return newReadCache
}

func (this *ReadCache[K, T]) Get(key K) (*T, bool) {
	if v, ok := this.cache[this.mapper(key)%uint64(len(this.cache))][key]; ok && v.Second != nil {
		return v.Second, ok
	}
	return nil, false
}

// Raw returns the value pair associated with the given key from the cache.
// It returns the value and a boolean indicating whether the key was found in the cache.
func (this *ReadCache[K, T]) Raw(key K) (*associative.Pair[*T, *T], bool) {
	if v, ok := this.cache[this.mapper(key)%uint64(len(this.cache))][key]; ok {
		return v, ok
	}
	return nil, false
}

// PreAlloc pre-allocates the cache with the given keys. So the values can be updated later using the Update function.
func (this *ReadCache[K, T]) PreAlloc(keys []K, isNil func(k K) bool) {
	for _, k := range keys {
		if isNil(k) {
			continue
		}

		shardId := this.mapper(k) % uint64(len(this.cache))
		if _, ok := this.cache[shardId][k]; !ok {
			this.cache[shardId][k] = &associative.Pair[*T, *T]{}
			v := this.cache[shardId][k]
			this.cache[shardId][k] = v // Create a new pair, waiting for the value to be set later.
		}
		this.dirties = append(this.dirties, &associative.Pair[uint64, K]{First: shardId, Second: k}) // Record the new key and its shard id.
	}
}

// Update updates the cache with the given keys and values. This function isn't thread safe.
func (this *ReadCache[K, T]) Update(keys []K, values []T) {
	slice.ParallelForeach(keys, runtime.NumCPU(), func(i int, k *K) {
		shardId := this.mapper(*k) % uint64(len(this.cache))
		fmt.Print(this.cache[shardId][*k].First, this.cache[shardId][*k].First)
		this.cache[shardId][*k].First = &values[i]
	})
}

// Finalize finalizes the cache by setting the new values to the old values.
func (this *ReadCache[K, T]) Finalize() {
	slice.ParallelForeach(this.dirties, runtime.NumCPU(), func(i int, v **associative.Pair[uint64, K]) {
		val := this.cache[(*v).First][((*v).Second)]
		val.Second = val.First
		val.First = nil
	})

	// Some keys are imported but later removed during conflict resolution.
	for _, k := range this.dirties {
		if vals := this.cache[k.First][k.Second]; vals.First == nil && vals.Second == nil {
			delete(this.cache[k.First], k.Second)
		}
	}
	this.Clear()
}

// Call this function to clear the cache.
func (this *ReadCache[K, T]) Clear() { this.dirties = this.dirties[:0] }
