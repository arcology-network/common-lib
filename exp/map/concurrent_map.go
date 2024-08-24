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
	"sort"
	"sync"

	"github.com/arcology-network/common-lib/common"
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
	numShards = common.Min(numShards, 256)
	return &ConcurrentMap[K, V]{
		isNilVal:   isNilVal,
		hasher:     partitionHasher,
		shards:     slice.NewDo(numShards, func(i int) map[K]V { return make(map[K]V, 64) }),
		shardLocks: make([]sync.RWMutex, numShards),
	}
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

// Get retrieves the value associated with the specified key from the ConcurrentMap.
// It returns the value and a boolean indicating whether the key was found.
func (this *ConcurrentMap[K, V]) Get(key K, args ...interface{}) (V, bool) {
	shardID := this.hasher(key) % uint64(len(this.shards))

	this.shardLocks[shardID].RLock()
	defer this.shardLocks[shardID].RUnlock()

	v, ok := this.shards[shardID][key]
	return v, ok
}

// Get retrieves the value associated with the specified key from the ConcurrentMap.
// It returns the value and a boolean indicating whether the key was found.
func (this *ConcurrentMap[K, V]) UnsafeGet(key K, args ...interface{}) (V, bool) {
	shardID := this.hasher(key) % uint64(len(this.shards))
	v, ok := this.shards[shardID][key]
	return v, ok
}

// BatchGet retrieves the values associated with the specified keys from the ConcurrentMap.
// It returns a slice of values in the same order as the keys.
func (this *ConcurrentMap[K, V]) BatchGet(keys []K, args ...interface{}) ([]V, []bool) {
	shardIds := slice.NewDo(len(keys), func(i int) uint64 {
		return this.Hash(keys[i])
	})

	// The BatchGet is usually called.
	values, found := make([]V, len(keys)), make([]bool, len(keys))
	slice.ParallelForeach(this.shards, len(keys), func(shard int, _ *map[K]V) {
		this.shardLocks[shard].RLock()
		defer this.shardLocks[shard].RUnlock()

		for i := 0; i < len(keys); i++ {
			if shardIds[i] == uint64(shard) {
				values[i], found[i] = this.shards[shard][keys[i]]
			}
		}
	})
	return values, found
}

func (this *ConcurrentMap[K, V]) delete(shard uint64, key K) {
	delete(this.shards[shard], key)
}

// Set associates the specified value with the specified key in the ConcurrentMap.
// If the value is nil, the key-value pair is deleted from the map.
// It returns an error if the shard ID is out of range.
func (this *ConcurrentMap[K, V]) Set(key K, v V, args ...interface{}) {
	shardID := this.Hash(key)
	if shardID >= uint64(len(this.shards)) {
		return
	}

	this.shardLocks[shardID].Lock()
	defer this.shardLocks[shardID].Unlock()

	if this.isNilVal(v) {
		this.delete(shardID, key)
	} else {
		this.shards[shardID][key] = v
	}
}

func (this *ConcurrentMap[K, V]) RemoveIf(condition func(K, V) bool) {
	slice.ParallelForeach(this.shards, len(this.shards), func(shardNum int, shard *map[K]V) {
		this.shardLocks[shardNum].Lock()
		defer this.shardLocks[shardNum].Unlock()
		RemoveIf(*shard, condition)
	})
}

// Set associates the specified value with the specified key in the ConcurrentMap.
// If the value is nil, the key-value pair is deleted from the map.
// It returns an error if the shard ID is out of range.
func (this *ConcurrentMap[K, V]) UnsafeSet(key K, v V, args ...interface{}) error {
	shardID := this.Hash(key)
	if shardID >= uint64(len(this.shards)) {
		return nil
	}

	if this.isNilVal(v) {
		this.delete(shardID, key)
	} else {
		this.shards[shardID][key] = v
	}
	return nil
}

// BatchUpdate updates the values associated with the specified keys in the ConcurrentMap using the provided updater function.
// The updater function takes the original value, the index of the key in the keys slice, the key, and the new value as arguments.
func (this *ConcurrentMap[K, V]) BatchUpdate(keys []K, values []V, updater func(origin V, index int, key K, value V) V) {
	shards := this.Hash8s(keys)
	slice.ParallelForeach(this.shards, len(this.shards), func(shardNum int, shard *map[K]V) {
		this.shardLocks[shardNum].Lock()
		defer this.shardLocks[shardNum].Unlock()

		for i := 0; i < len(keys); i++ {
			if shards[i] == uint64(shardNum) {
				(*shard)[keys[i]] = updater((*shard)[keys[i]], i, keys[i], values[i])
			}
		}
	})
}

// BatchSet associates the specified values with the specified keys in the ConcurrentMap.
// If the values slice is shorter than the keys slice, the remaining keys are deleted from the map.
// If the values slice is longer than the keys slice, the extra values are ignored.
func (this *ConcurrentMap[K, V]) BatchSet(keys []K, values []V) {
	shardIDs := this.Hash8s(keys)
	this.BatchSetToShards(shardIDs, keys, values)
}

// func (this *ConcurrentMap[K, V]) BatchSetIf(keys []K, setter func(K) (V, bool)) {
// 	values, flags := make([]V, len(keys)), make([]bool, len(keys))
// 	for i := 0; i < len(keys); i++ {
// 		values[i], flags[i] = setter(keys[i])
// 	}

// 	slice.RemoveIf(&keys, func(i int, k K) bool { return flags[i] })
// 	slice.RemoveIf(&values, func(i int, v V) bool { return flags[i] })

// 	this.BatchSetToShards(this.Hash8s(keys), keys, values)
// }

// BatchSetToShard associates the specified values with the specified shard IDs and keys in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) BatchSetToShards(ids []uint64, keys []K, values []V) {
	slice.ParallelForeach(this.shards, 8, func(shardNum int, shard *map[K]V) {
		this.shardLocks[shardNum].Lock()
		defer this.shardLocks[shardNum].Unlock()

		for i := 0; i < len(ids); i++ {
			if ids[i] == uint64(shardNum) { // If the key belongs to this shard
				if this.isNilVal(values[i]) {
					delete(this.shards[shardNum], keys[i]) // Delete the key-value pair from the shard.
					return
				}
				this.shards[shardNum][keys[i]] = values[i] // Update the value in the shard.
			}
		}
	})
}

func (this *ConcurrentMap[K, V]) BatchSetWith(keys []K, setter func(k *K) V) {
	shardIDs := this.Hash8s(keys)
	this.BatchSetToShardWith(shardIDs, keys, setter)
}

// BatchSetToShard associates the specified values with the specified shard IDs and keys in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) BatchSetToShardWith(ids []uint64, keys []K, setter func(k *K) V) {
	slice.ParallelForeach(this.shards, 8, func(shardNum int, shard *map[K]V) {
		this.globalLock.Lock()
		defer this.globalLock.Unlock()

		for i := 0; i < len(ids); i++ {
			if ids[i] == uint64(shardNum) { // If the key belongs to this shard
				v := setter(&keys[i])
				if this.isNilVal(v) {
					delete(this.shards[shardNum], keys[i]) // Delete the key-value pair from the shard.
					return
				}
				this.shards[shardNum][keys[i]] = v // Update the value in the shard.
			}
		}
	})
}

// Keys returns a slice containing all the keys in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) Keys() []K {
	this.globalLock.RLock()
	defer this.globalLock.RUnlock()

	keySet := slice.ParallelTransform(this.shards, 8, func(i int, m map[K]V) []K {
		return common.MapKeys(m)
	})
	return slice.Flatten(keySet)
}

func (this *ConcurrentMap[K, V]) Values() []V {
	this.globalLock.RLock()
	defer this.globalLock.RUnlock()

	valSet := slice.ParallelTransform(this.shards, 8, func(i int, m map[K]V) []V {
		return common.MapValues(m)
	})
	return slice.Flatten(valSet)
}

// Hash8s calculates the shard IDs for the specified keys using the Hash8 function.
// It returns a slice of shard IDs in the same order as the keys.
func (this *ConcurrentMap[K, V]) Hash8s(keys []K) []uint64 {
	return slice.ParallelNew(len(keys), 8, func(i int) uint64 {
		return this.Hash(keys[i])
	})
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
		common.MapForeach(*shard, func(k K, v *V) {
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

func (this *ConcurrentMap[K, V]) ParallelDo(keys []K, do func(i int, k K, v V, b bool) (V, bool)) {
	this.globalLock.Lock()
	defer this.globalLock.Unlock()

	values, found := this.BatchGet(keys)

	assignBack := make([]bool, len(keys))
	slice.ParallelForeach(found, 4, func(i int, _ *bool) {
		values[i], assignBack[i] = do(i, keys[i], values[i], found[i])
	})

	slice.RemoveBothIf(&keys, &values, func(i int, _ K, _ V) bool {
		return !assignBack[i]
	})
	this.BatchSet(keys, values)
}

// ParallelFor applies the specified do function to each key-value pair in a range defined by the first and last indices.
// It is useful for iterating over a slice containing the keys in parallel and updating the values in the map.
func (this *ConcurrentMap[K, V]) ParallelFor(first int, last int, key func(i int) K, do func(i int, k K, v V, b bool) (V, bool)) {
	this.globalLock.Lock()
	defer this.globalLock.Unlock()

	common.ParallelFor(first, last, 4, func(i int) {
		k := key(i)
		v, b := this.UnsafeGet(k)
		if newV, ok := do(i, k, v, b); ok {
			this.Set(k, newV)
		}
	})
}

// ParallelFor applies the specified do function to each key-value pair in a range defined by the first and last indices.
// It is useful for iterating over a slice containing the keys in parallel and updating the values in the map.
func (this *ConcurrentMap[K, V]) UnsafeParallelFor(first int, last int, key func(i int) K, do func(i int, k K, v V, b bool) (V, bool)) {
	common.ParallelFor(first, last, 4, func(i int) {
		k := key(i)
		v, b := this.UnsafeGet(k)
		if newV, ok := do(i, k, v, b); ok {
			this.UnsafeSet(k, newV)
		}
	})
}

// KVs returns two slices: one containing all the keys in the ConcurrentMap, and one containing the corresponding values.
func (this *ConcurrentMap[K, V]) KVs(less ...func(k0, k1 K) bool) ([]K, []V) {
	this.globalLock.Lock()
	defer this.globalLock.Unlock()

	keys := this.Keys()
	if len(less) > 0 {
		sort.Slice(keys, func(i, j int) bool {
			return less[0](keys[i], keys[j])
		})
	}
	return keys, common.FilterFirst(this.BatchGet(keys))
}

func (this *ConcurrentMap[K, V]) Checksum(less func(K, K) bool, encoders func(K, V) ([]byte, []byte)) [32]byte {
	keys, values := this.KVs(less)
	kBytes, vBytes := make([][]byte, len(keys)), make([][]byte, len(values))
	for i, _ := range values {
		kBytes[i], vBytes[i] = encoders(keys[i], values[i])
	}

	kSum := sha256.Sum256(slice.Flatten(kBytes))
	vSum := sha256.Sum256(slice.Flatten(vBytes))
	return sha256.Sum256(append(kSum[:], vSum[:]...))
}

func (this *ConcurrentMap[K, V]) Print() {}
