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

// The ConcurrentMap class is a concurrent map implementation allowing
// multiple goroutines to access and modify the map concurrently.

package mapi

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"sync"

	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/exp/associative"
	slice "github.com/arcology-network/common-lib/exp/slice"
)

// ConcurrentMap represents a concurrent map data structure.
type ConcurrentMap[K comparable, V any] struct {
	hasher     func(K) uint64
	isNilVal   func(V) bool
	shards     []map[K]V
	shardLocks []sync.RWMutex
	globalLock sync.RWMutex
}

// NewConcurrentMap creates a new instance of ConcurrentMap with the specified number of shards.
// If no number of shards is provided, it defaults to 6.
func NewConcurrentMap[K comparable, V any](
	numShards int,
	isNilVal func(V) bool,
	partitionHasher func(K) uint64) *ConcurrentMap[K, V] {
	numShards = max(min(numShards, 256), 1)
	return &ConcurrentMap[K, V]{
		isNilVal:   isNilVal,
		hasher:     partitionHasher,
		shards:     slice.NewDo(numShards, func(i int) map[K]V { return make(map[K]V, 64) }),
		shardLocks: make([]sync.RWMutex, numShards),
	}
}

func (this *ConcurrentMap[K, V]) LockShard(shardID uint64) func() {
	this.shardLocks[shardID].Lock()
	return func() {
		this.shardLocks[shardID].Unlock()
	}
}

func (this *ConcurrentMap[K, V]) LockReadShard(shardID uint64) func() {
	this.shardLocks[shardID].RLock()
	return func() {
		this.shardLocks[shardID].RUnlock()
	}
}

func (this *ConcurrentMap[K, V]) NumShards() uint64 {
	return uint64(int(len(this.shards)))
}

func (this *ConcurrentMap[K, V]) ComputeShard(key K) uint64 {
	return this.hasher(key) % uint64(len(this.shards))
}

// Sizes returns a slice containing the number of key-value pairs in each shard.
func (this *ConcurrentMap[K, V]) Sizes() []int {
	sizes := make([]int, len(this.shards))
	for i, shard := range this.shards {
		sizes[i] = len(shard)
	}
	return sizes
}

func (this *ConcurrentMap[K, V]) Clear() {
	this.globalLock.Lock()
	defer this.globalLock.Unlock()

	slice.ParallelForeach(this.shards, len(this.shards), func(i int, _ *map[K]V) {
		clear(this.shards[i])
	})
}

// Size returns the total number of key-value pairs in the ConcurrentMap.
// It isn't exactly accurate since the map is being accessed concurrently.
func (this *ConcurrentMap[K, V]) Length() uint64 {
	this.globalLock.RLock()
	defer this.globalLock.RUnlock()

	v := slice.Accumulate(this.shards, 0, func(i int, m map[K]V) uint64 {
		return uint64(len(m))
	})
	return v
}

// Get the value associated with the specified key from the ConcurrentMap.
// It returns the value and a boolean indicating whether the key was found.
func (this *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	return this.UnsafeGet(key, this.LockShard)
}

// GetBatch retrieves the values associated with the specified keys from the ConcurrentMap.
// It returns a slice of values in the same order as the keys.
func (this *ConcurrentMap[K, V]) GetBatch(keys []K) ([]V, []bool) {
	return this.UnsafeBatchGet(keys, this.LockReadShard)
}

func (this *ConcurrentMap[K, V]) Delete(key K) {
	this.UnsafeDelete(key, this.LockShard)
}

// Set associates the specified value with the specified key in the ConcurrentMap.
// If the value is nil, the key-value pair is deleted from the map.
// It returns an error if the shard ID is out of range.
func (this *ConcurrentMap[K, V]) Set(key K, v V) {
	this.UnsafeSet(key, v, this.LockShard)
}

// SetBatch associates the specified values with the specified keys in the ConcurrentMap.
// If the values slice is shorter than the keys slice, the remaining keys are deleted from the map.
// If the values slice is longer than the keys slice, the extra values are ignored.
func (this *ConcurrentMap[K, V]) SetBatch(keys []K, values []V) {
	this.UnsafeBatchSet(keys, values, this.LockShard)
}

// Keys returns a slice containing all the keys in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) Keys() []K {
	this.globalLock.RLock()
	defer this.globalLock.RUnlock()

	keySet := slice.ParallelTransform(this.shards, 8, func(i int, m map[K]V) []K {
		return Keys(m)
	})
	return slice.Flatten(keySet)
}

func (this *ConcurrentMap[K, V]) Values() []V {
	this.globalLock.RLock()
	defer this.globalLock.RUnlock()

	valSet := slice.ParallelTransform(this.shards, 8, func(i int, m map[K]V) []V {
		return Values(m)
	})
	return slice.Flatten(valSet)
}

// Hash8s calculates the shard IDs for the specified keys using the Hash8 function.
// It returns a slice of shard IDs in the same order as the keys.
func (this *ConcurrentMap[K, V]) Hash8s(keys []K) associative.Pairs[uint64, uint64] {
	pairs := slice.ParallelNew(len(keys), 8, func(i int) *associative.Pair[uint64, uint64] {
		return &associative.Pair[uint64, uint64]{
			First:  uint64(i),
			Second: this.Hash(keys[i]),
		}
	})
	return associative.Pairs[uint64, uint64](pairs)
}

func (this *ConcurrentMap[K, V]) Hash(key K) uint64 {
	return this.hasher(key) % uint64(len(this.shards))
}

// Shards returns a pointer to the slice of maps representing the shards in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) Shards() []map[K]V {
	return this.shards
}

// Traverse applies the specified operator function to each key-value pair in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) Traverse(processor func(K, *V)) {
	this.globalLock.RLock()
	defer this.globalLock.RUnlock()

	slice.ParallelForeach(this.shards, 8, func(i int, shard *map[K]V) {
		Foreach(*shard, func(k K, v *V) {
			processor(k, v)
		})
	})
}

// Foreach applies the specified predicate function to each value in the ConcurrentMap.
// The predicate function takes a value as an argument and returns a new value.
// If the new value is nil, the key-value pair is deleted from the map.
// Otherwise, the new value replaces the original value in the map.
func (this *ConcurrentMap[K, V]) Foreach(predicate func(V) V) {
	this.globalLock.RLock()
	defer this.globalLock.RUnlock()

	slice.ParallelForeach(this.shards, len(this.shards), func(i int, _ *map[K]V) {
		this.shardLocks[i].Lock()
		defer this.shardLocks[i].Unlock()

		for k, v := range this.shards[i] {
			if v = predicate(v); this.isNilVal(v) {
				delete(this.shards[i], k)
				continue
			}
			this.shards[i][k] = v
		}
	})
}

// ForeachDo applies the specified do function to each key-value pair in the ConcurrentMap.
// The do function takes a key and a value as arguments and performs some action.
func (this *ConcurrentMap[K, V]) ForeachDo(do func(K, V)) {
	for i := 0; i < len(this.shards); i++ {
		this.shardLocks[i].Lock()
		for k, v := range this.shards[i] {
			do(k, v)
		}
		this.shardLocks[i].Unlock()
	}
}

// ParallelForeachDo applies the specified do function to each key-value pair in the ConcurrentMap in parallel.
// The do function takes a key and a value as arguments and performs some action.
func (this *ConcurrentMap[K, V]) ParallelForeachDo(do func(K, V)) {
	slice.ParallelForeach(this.shards, len(this.shards), func(shardNum int, shard *map[K]V) {
		this.shardLocks[shardNum].RLock()
		defer this.shardLocks[shardNum].RUnlock()

		for k, v := range *shard {
			do(k, v)
		}
	})
}

// Parallel deletes the key-value pairs for which the specified delete function returns true.
func (this *ConcurrentMap[K, V]) ParallelDelete(del func(K, V) bool) {
	this.globalLock.Lock()
	defer this.globalLock.Unlock()

	slice.ParallelForeach(this.shards, len(this.shards), func(_ int, shard *map[K]V) {
		for k, v := range *shard {
			if del(k, v) {
				delete(*shard, k)
			}
		}
	})
}

func (this *ConcurrentMap[K, V]) ParallelDo(keys []K, do func(i int, k K, v V, b bool) (V, bool)) {
	this.globalLock.Lock()
	defer this.globalLock.Unlock()

	values, found := this.GetBatch(keys)

	assignBack := make([]bool, len(keys))
	slice.ParallelForeach(found, 4, func(i int, _ *bool) {
		values[i], assignBack[i] = do(i, keys[i], values[i], found[i])
	})

	slice.RemoveBothIf(&keys, &values, func(i int, _ K, _ V) bool {
		return !assignBack[i]
	})
	this.SetBatch(keys, values)
}

// ParallelFor applies the specified do function to each key-value pair in a range defined by the first and last indices.
// It is useful for iterating over a slice containing the keys in parallel and updating the values in the map.
func (this *ConcurrentMap[K, V]) ParallelFor(first int, last int, key func(i int) K, do func(i int, k K, v V, b bool) (V, bool)) {
	this.globalLock.Lock()
	defer this.globalLock.Unlock()

	common.ParallelFor(first, last, 4, func(i int) {
		k := key(i)
		v, b := this.UnsafeGet(k, this.LockShard)
		if newV, ok := do(i, k, v, b); ok {
			this.Set(k, newV)
		}
	})
}

// ParallelFor applies the specified do function to each key-value pair in a range defined by the first and last indices.
// It is useful for iterating over a slice containing the keys in parallel and updating the values in the map.
// KVs returns two slices: one containing all the keys in the ConcurrentMap, and one containing the corresponding values.
func (this *ConcurrentMap[K, V]) KVs(less ...func(k0, k1 K) bool) ([]K, []V) {
	this.globalLock.RLock()
	defer this.globalLock.RUnlock()

	keys := this.Keys()
	if len(less) > 0 {
		sort.Slice(keys, func(i, j int) bool {
			return less[0](keys[i], keys[j])
		})
	}
	return keys, common.First(this.GetBatch(keys)).([]V)
}

func (this *ConcurrentMap[K, V]) Checksum(less func(K, K) bool, encoders func(K, V) ([]byte, []byte)) [32]byte {
	keys, values := this.KVs(less)
	kBytes, vBytes := make([][]byte, len(keys)), make([][]byte, len(values))
	for i := range values {
		kBytes[i], vBytes[i] = encoders(keys[i], values[i])
	}

	kSum := sha256.Sum256(slice.Flatten(kBytes))
	vSum := sha256.Sum256(slice.Flatten(vBytes))
	return sha256.Sum256(append(kSum[:], vSum[:]...))
}

func (this *ConcurrentMap[K, V]) Print(toString func(K, V) string) {
	slice.Foreach(this.shards, func(i int, shard *map[K]V) {
		for k, v := range *shard {
			fmt.Println("Shard", i, toString(k, v))
		}
	})
}
