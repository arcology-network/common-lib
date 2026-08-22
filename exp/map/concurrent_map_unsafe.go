/*
 *   Copyright (c) 2026 Arcology Network

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
	"sort"

	slice "github.com/arcology-network/common-lib/exp/slice"
)

// UnsafeGet returns the value for key.
func (this *ConcurrentMap[K, V]) UnsafeGet(key K, lock func(uint64) func()) (V, bool) {
	shardID := this.ComputeShard(key)
	if lock != nil {
		defer lock(shardID)()
	}

	v, ok := this.shards[shardID][key]
	return v, ok
}

// UnsafeSet stores or removes a value.
func (this *ConcurrentMap[K, V]) UnsafeSet(key K, value V, lock func(uint64) func()) error {
	shardID := this.ComputeShard(key)
	if lock != nil {
		defer lock(shardID)()
	}

	this.shards[shardID][key] = value
	return nil
}

// UnsafeDelete removes key.
func (this *ConcurrentMap[K, V]) UnsafeDelete(key K, lock func(uint64) func()) {
	shardID := this.ComputeShard(key)
	if lock != nil {
		defer lock(shardID)()
	}
	delete(this.shards[shardID], key)
}

// UnsafeBatchUpdate updates keys using updater.
func (this *ConcurrentMap[K, V]) unsafeBatchDo(
	keys []K,
	lock func(uint64) func(),
	doer func(i int, m map[K]V) bool,
) []bool {
	shardIndices := this.Hash8s(keys)
	sort.Slice(shardIndices, func(i, j int) bool {
		_0, _1 := shardIndices[i].Second, shardIndices[j].Second
		if _0 == _1 {
			return shardIndices[i].First < shardIndices[j].First
		}
		return _0 < _1
	})

	// Each shard occupies [boundaries[shard], boundaries[shard+1]).
	boundaries := make([]int, len(this.shards)+1)
	offset := 0
	for shard := range this.shards {
		boundaries[shard] = offset
		for offset < len(shardIndices) && shardIndices[offset].Second == uint64(shard) {
			offset++
		}
	}
	boundaries[len(this.shards)] = offset

	successes := make([]bool, len(keys))
	slice.ParallelForeach(this.shards, len(this.shards), func(shard int, _ *map[K]V) {
		if boundaries[shard] == boundaries[shard+1] {
			return
		}
		if lock != nil {
			defer lock(uint64(shard))()
		}
		for i := boundaries[shard]; i < boundaries[shard+1]; i++ {
			pair := shardIndices[i]
			successes[pair.First] = doer(int(pair.First), this.shards[shard])
		}
	})
	return successes
}

// UnsafeBatchGet returns values for keys.
func (this *ConcurrentMap[K, V]) UnsafeBatchGet(keys []K, lock func(uint64) func()) ([]V, []bool) {
	vals := make([]V, len(keys))
	getter := func(i int, shard map[K]V) bool {
		var ok bool
		vals[i], ok = shard[keys[i]]
		return ok
	}

	return vals, this.unsafeBatchDo(keys, lock, getter)
}

func (this *ConcurrentMap[K, V]) UnsafeBatchSet(keys []K, values []V, lock func(uint64) func()) []bool {
	setter := func(i int, shard map[K]V) bool {
		if i >= len(values) || this.isNilVal(values[i]) {
			delete(shard, keys[i])
			return true
		}
		shard[keys[i]] = values[i]
		return true
	}
	return this.unsafeBatchDo(keys, lock, setter)
}

// UnsafeDeleteBatch removes keys.
func (this *ConcurrentMap[K, V]) UnsafeBatchDelete(keys []K, lock func(uint64) func()) {
	deleter := func(i int, shard map[K]V) bool {
		delete(shard, keys[i])
		return true
	}
	this.unsafeBatchDo(keys, lock, deleter)
}
