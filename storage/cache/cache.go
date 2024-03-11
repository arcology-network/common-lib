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
	"runtime"

	"github.com/arcology-network/common-lib/exp/associative"
	"github.com/arcology-network/common-lib/exp/slice"
)

// ReadCache is a concurrent map of K to T.
type ReadCache[K comparable, T any] struct {
	mapper    func(K) uint64
	cache     []map[K]*associative.Pair[*T, *T]
	newBuffer []*associative.Pair[uint64, K]
	// stats     []associative.Pair[K, uint64]
}

func NewReadCache[K comparable, T any](numShards uint64, mapper func(K) uint64) *ReadCache[K, T] {
	newReadCache := &ReadCache[K, T]{
		mapper:    mapper,
		cache:     make([]map[K]*associative.Pair[*T, *T], 2),
		newBuffer: make([]*associative.Pair[uint64, K], 0, 1024),
		// stats:     make([]associative.Pair[K, uint64], 0, 1024),
	}

	for i := range newReadCache.cache {
		newReadCache.cache[i] = make(map[K]*associative.Pair[*T, *T])
	}
	return newReadCache
}

func (this *ReadCache[K, T]) Get(key K) (*associative.Pair[*T, *T], bool) {
	if v, ok := this.cache[this.mapper(key)%uint64(len(this.cache))][key]; ok {
		return v, ok
	}
	return nil, false
}

func (this *ReadCache[K, T]) PreAlloc(keys []K, isNil func(k K) bool) {
	for _, k := range keys {
		if isNil(k) {
			continue
		}

		shardId := this.mapper(k) % uint64(len(this.cache))
		if _, ok := this.cache[shardId][k]; !ok {
			this.cache[shardId][k] = &associative.Pair[*T, *T]{}
			v := this.cache[shardId][k]
			this.cache[shardId][k] = v                                                                       // Create a new pair, waiting for the value to be set later.
			this.newBuffer = append(this.newBuffer, &associative.Pair[uint64, K]{First: shardId, Second: k}) // Record the new key and its shard id.

			// pair := this.newBuffer[len(this.newBuffer)-1]
			// val := this.cache[pair.First][*pair.Second]
			// val.First = val.First

			for _, v := range this.newBuffer {
				val := this.cache[(*v).First][((*v).Second)]
				val.First = val.First
			}
		}
	}

	// slice.ParallelForeach(this.newBuffer, runtime.NumCPU(), func(i int, v **associative.Pair[uint64, *K]) {
	// 	val := this.cache[(*v).First][*((*v).Second)]
	// 	val.First = val.First
	// })
	// for _, v := range this.newBuffer {
	// 	val := this.cache[(*v).First][((*v).Second)]
	// 	val.First = val.First
	// }
}

// source M, keys []T, getter func(T) K, do func(K) V)
func (this *ReadCache[K, T]) Update(keys []K, values []T) {
	slice.ParallelForeach(keys, runtime.NumCPU(), func(i int, k *K) {
		shardId := this.mapper(*k) % uint64(len(this.cache))
		this.cache[shardId][*k].First = &values[i]
	})

	// for _, v := range this.newBuffer {
	// 	val := this.cache[(v).First][*((v).Second)]
	// 	val.Second = val.First
	// 	// k := getter(*v)
	// 	// shardId := this.mapper(k) % uint64(len(this.cache))
	// 	// this.cache[shardId][k].First = values[i]
	// }
}

// source M, keys []T, getter func(T) K, do func(K) V)
func (this *ReadCache[K, T]) Finalize() {
	slice.ParallelForeach(this.newBuffer, runtime.NumCPU(), func(i int, v **associative.Pair[uint64, K]) {
		val := this.cache[(*v).First][((*v).Second)]
		val.Second = val.First
		val.First = nil
	})

	// Some keys are imported but later removed during conflict resolution.
	for _, k := range this.newBuffer {
		if vals := this.cache[k.First][k.Second]; vals.First == nil && vals.Second == nil {
			delete(this.cache[k.First], k.Second)
		}
	}
	this.Clear()
}

// Call this function to clear the cache.
func (this *ReadCache[K, T]) Clear() { this.newBuffer = this.newBuffer[:0] }
