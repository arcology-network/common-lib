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
package livecache

import (
	"fmt"
	"runtime"
	"sync/atomic"

	crdtcommon "github.com/arcology-network/common-lib/crdt/common"
	statecell "github.com/arcology-network/common-lib/crdt/statecell"
	mapi "github.com/arcology-network/common-lib/exp/map"
	"github.com/arcology-network/common-lib/exp/slice"
)

// Live cache helpers for `CachedKVStore`.

func NewLiveCache(cacheCap uint64) *CachedKVStore[string, crdtcommon.CRDT] {
	store := NewStore[string, crdtcommon.CRDT](nil, cacheCap)
	store.SetLocalOnly(true)
	return store
}

func (this *CachedKVStore[K, T]) Policy() *CachePolicy[*Entry[T]]  { return this.cachePolicy }
func (this *CachedKVStore[K, T]) Profile() *CachePolicy[*Entry[T]] { return this.cachePolicy }

func (this *CachedKVStore[K, T]) CacheChecksum() [32]byte {
	encoders := func(k K, v *Entry[T]) ([]byte, []byte) {
		key := []byte(fmt.Sprintf("%v", k))
		if v == nil {
			return key, nil
		}

		encoder, ok := any(v.Value).(interface{ Encode() []byte })
		if !ok {
			return key, nil
		}
		return key, encoder.Encode()
	}

	less := func(k0, k1 K) bool {
		return fmt.Sprintf("%v", k0) < fmt.Sprintf("%v", k1)
	}
	return this.ConcurrentMap.Checksum(less, encoders)
}

// Get the raw value from the cache with the usage information.
func (this *CachedKVStore[K, T]) GetRaw(key K) (*Entry[T], bool) {
	v, ok := this.ConcurrentMap.Get(key)
	if !ok || v == nil {
		return nil, false
	}
	return v, true
}

func (this *CachedKVStore[K, T]) CommitStateCells(univals []*statecell.StateCell, block uint64) {
	// Prepare the space for the new values in the cache, some univalues may be deleted because of the memory limit.
	this.cachePolicy.PrepareSpace(&univals, this.freeCache)

	// Extract the keys and values from the univalues.
	keys := slice.ParallelTransform(univals, runtime.NumCPU(), func(i int, v *statecell.StateCell) K {
		key, _ := any(*v.GetPath()).(K)
		return key
	})

	entries := slice.ParallelTransform(univals, runtime.NumCPU(), func(i int, v *statecell.StateCell) *Entry[T] {
		if v.Value() == nil {
			return nil
		}

		value, ok := any(v.Value()).(T)
		if !ok {
			return nil
		}

		entry := &Entry[T]{
			Value: value,
			Stat: Stat{
				sizeInMem: func() uint64 {
					sized, ok := any(v.Value()).(interface{ MemSize() uint64 })
					if !ok {
						return 0
					}
					return sized.MemSize()
				}(),
				visits:      uint64(v.Reads()) + uint64(v.Writes()) + uint64(v.DeltaWrites()),
				firstLoaded: uint32(block),
			},
			SizeOf: func(v T) uint64 {
				sized, ok := any(v).(interface{ MemSize() uint64 })
				if !ok {
					return 0
				}
				return sized.MemSize()
			},
		}

		if metav, _ := this.GetRaw(keys[i]); metav != nil {
			entry.visits += metav.visits
			entry.firstLoaded = metav.firstLoaded
		}

		return entry
	})

	this.ConcurrentMap.BatchSet(keys, entries)
}

func (this *CachedKVStore[K, T]) freeCache(sizeToFree uint64) uint64 {
	var totalFreed atomic.Uint64
	shards := this.ConcurrentMap.Shards()

	sizeToFree, shardTarget := this.cachePolicy.AdjustFreeTarget(sizeToFree, len(shards))
	if sizeToFree == 0 || len(shardTarget) == 0 {
		return 0
	}

	slice.ParallelForeach(shards, runtime.NumCPU(), func(i int, _ *map[K]*Entry[T]) {
		if len(shards[i]) == 0 {
			return
		}

		ks, values := mapi.KVs(shards[i])
		scores := slice.ParallelTransform(values, runtime.NumCPU(), func(i int, entry *Entry[T]) float32 {
			return this.cachePolicy.EvictionScore(entry.visits, entry.firstLoaded, sizeToFree)
		})

		slice.SortBy1st(scores, ks, func(v0, v1 float32) bool {
			return v0 < v1
		})

		for j, key := range ks {
			freedSize := values[j].sizeInMem
			delete(shards[i], key)
			totalFreed.Add(freedSize)

			if shardTarget[i] -= float64(freedSize); shardTarget[i] <= 0 {
				break
			}
		}
	})

	return totalFreed.Load()
}

// Print the content of the cache for debugging.
func (this *CachedKVStore[K, T]) Print() {
	keys, vals := this.ConcurrentMap.KVs()
	slice.SortBy1st(keys, vals, func(k0, k1 K) bool {
		return fmt.Sprintf("%v", k0) < fmt.Sprintf("%v", k1)
	})

	fmt.Println("occupied:", this.cachePolicy.occupied)

	for i, k := range keys {
		println(k, "  =    ")
		printer, ok := any(vals[i].Value).(interface{ Print() })
		if ok {
			printer.Print()
		}
	}
}
