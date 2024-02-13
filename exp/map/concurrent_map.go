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
	"sync"

	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/exp/array"
)

// ConcurrentMap represents a concurrent map data structure.
type ConcurrentMap[K comparable, V any] struct {
	hasher     func(K) uint8
	isNilVal   func(V) bool
	shards     []map[K]V
	shardLocks []sync.RWMutex
}

// NewConcurrentMap creates a new instance of ConcurrentMap with the specified number of shards.
// If no number of shards is provided, it defaults to 6.
func NewConcurrentMap[K comparable, V any](numShards int, isNilVal func(V) bool, hasher func(K) uint8) *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		isNilVal:   isNilVal,
		hasher:     hasher,
		shards:     array.NewWith(numShards, func(i int) map[K]V { return make(map[K]V, 64) }),
		shardLocks: make([]sync.RWMutex, numShards),
	}
}

// Size returns the total number of key-value pairs in the ConcurrentMap.
// It isn't exactly accurate since the map is being accessed concurrently.
func (this *ConcurrentMap[K, V]) Size() uint32 {
	v := array.Accumulate[map[K]V, int](this.shards, 0, func(i int, m map[K]V) int {
		return len((m))
	})
	return uint32(v)
}

// Get retrieves the value associated with the specified key from the ConcurrentMap.
// It returns the value and a boolean indicating whether the key was found.
func (this *ConcurrentMap[K, V]) Get(key K, args ...interface{}) (V, bool) {
	shardID := this.hasher(key) % uint8(len(this.shards))

	this.shardLocks[shardID].RLock()
	defer this.shardLocks[shardID].RUnlock()

	v, ok := this.shards[shardID][key]
	return v, ok
}

// BatchGet retrieves the values associated with the specified keys from the ConcurrentMap.
// It returns a slice of values in the same order as the keys.
func (this *ConcurrentMap[K, V]) BatchGet(keys []K, args ...interface{}) []V {
	shardIds := array.NewWith(len(keys), func(i int) uint8 {
		return this.Hash(keys[i])
	})

	values := make([]V, len(keys))
	array.ParallelForeach(this.shards, len(keys), func(shard int, _ *map[K]V) {
		this.shardLocks[shard].RLock()
		defer this.shardLocks[shard].RUnlock()

		for i := 0; i < len(keys); i++ {
			if shardIds[i] == uint8(shard) {
				values[i] = this.shards[uint8(shard)][keys[i]]
			}
		}
	})
	return values
}

// DirectBatchGet retrieves the values associated with the specified shard IDs and keys from the ConcurrentMap.
// It returns a slice of values in the same order as the keys.
func (this *ConcurrentMap[K, V]) DirectBatchGet(shardIDs []uint8, keys []K, args ...interface{}) []V {
	return array.ParallelAppend(keys, 5, func(i int, _ K) V {
		return this.shards[shardIDs[i]][keys[i]]
	})
}

func (this *ConcurrentMap[K, V]) delete(shardID uint8, key K) {
	delete(this.shards[shardID], key)
}

// Set associates the specified value with the specified key in the ConcurrentMap.
// If the value is nil, the key-value pair is deleted from the map.
// It returns an error if the shard ID is out of range.
func (this *ConcurrentMap[K, V]) Set(key K, v V, args ...interface{}) error {
	shardID := this.Hash(key)
	if shardID >= uint8(uint8(len(this.shards))) {
		return nil
	}

	this.shardLocks[shardID].Lock()
	defer this.shardLocks[shardID].Unlock()

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
	array.ParallelForeach(this.shards, len(this.shards), func(shardNum int, shard *map[K]V) {
		for i := 0; i < len(keys); i++ {
			if shards[i] == uint8(shardNum) {
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
	this.DirectBatchSet(shardIDs, keys, values)
}

func (this *ConcurrentMap[K, V]) BatchSetIf(keys []K, setter func(K) (V, bool)) {
	values, flags := make([]V, len(keys)), make([]bool, len(keys))
	for i := 0; i < len(keys); i++ {
		values[i], flags[i] = setter(keys[i])
	}

	array.RemoveIf(&keys, func(i int, k K) bool { return flags[i] })
	array.RemoveIf(&values, func(i int, v V) bool { return flags[i] })

	this.DirectBatchSet(this.Hash8s(keys), keys, values)
}

// DirectBatchSet associates the specified values with the specified shard IDs and keys in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) DirectBatchSet(ids []uint8, keys []K, values []V) {
	array.ParallelForeach(this.shards, 8, func(shardNum int, shard *map[K]V) {
		for i := 0; i < len(ids); i++ {
			if ids[i] == uint8(shardNum) { // If the key belongs to this shard
				if this.isNilVal(values[i]) {
					delete(this.shards[shardNum], keys[i]) // Delete the key-value pair from the shard.
					return
				}
				this.shards[shardNum][keys[i]] = values[i] // Update the value in the shard.
			}
		}
	})
}

// Keys returns a slice containing all the keys in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) Keys() []K {
	for i := 0; i < int(uint8(len(this.shards))); i++ {
		this.shardLocks[i].Lock()
		defer this.shardLocks[i].Unlock()
	}

	keySet := array.ParallelAppend[map[K]V](this.shards, 8, func(i int, m map[K]V) []K {
		return common.MapKeys(m)
	})
	return array.Flatten(keySet)
}

// Hash8s calculates the shard IDs for the specified keys using the Hash8 function.
// It returns a slice of shard IDs in the same order as the keys.
func (this *ConcurrentMap[K, V]) Hash8s(keys []K) []uint8 {
	return array.ParallelNew(len(keys), 8, func(i int) uint8 {
		return this.Hash(keys[i])
	})
}

func (this *ConcurrentMap[K, V]) Hash(key K) uint8 {
	return this.hasher(key) % uint8(len(this.shards))
}

// Shards returns a pointer to the slice of maps representing the shards in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) Shards() []map[K]V {
	return this.shards
}

// Traverse applies the specified operator function to each key-value pair in the ConcurrentMap.
func (this *ConcurrentMap[K, V]) Traverse(processor func(K, *V)) {
	array.ParallelForeach(this.shards, 8, func(i int, shard *map[K]V) {
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
	for i := 0; i < int(uint8(len(this.shards))); i++ {
		this.shardLocks[i].Lock()
		defer this.shardLocks[i].Unlock()
	}

	array.ParallelForeach(this.shards, len(this.shards), func(i int, _ *map[K]V) {
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
		this.shardLocks[i].RLock()
		for k, v := range this.shards[i] {
			do(k, v)
		}
		this.shardLocks[i].RUnlock()
	}
}

// ParallelForeachDo applies the specified do function to each key-value pair in the ConcurrentMap in parallel.
// The do function takes a key and a value as arguments and performs some action.
func (this *ConcurrentMap[K, V]) ParallelForeachDo(do func(K, V)) {
	array.ParallelForeach(this.shards, len(this.shards), func(_ int, shard *map[K]V) {
		for k, v := range *shard {
			do(k, v)
		}
	})
}

// KVs returns two slices: one containing all the keys in the ConcurrentMap, and one containing the corresponding values.
func (this *ConcurrentMap[K, V]) KVs() ([]K, []V) {
	keys := this.Keys()
	return keys, this.BatchGet(keys)
}

func (this *ConcurrentMap[K, V]) Checksum() [32]byte {
	// k, values := this.KVs()
	// vBytes := []byte{}
	// for _, v := range values {
	// 	vBytes = append(vBytes, v.(Encodable).Encode()...)
	// }

	// kSum := sha256.Sum256(codec.Strings(k).Flatten())
	// vSum := sha256.Sum256(vBytes)

	// return sha256.Sum256(append(kSum[:], vSum[:]...))
	return [32]byte{}
}

func (this *ConcurrentMap[K, V]) Print() {}
