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

package cachedkvstore

import (
	"fmt"

	mapi "github.com/arcology-network/common-lib/exp/map"
	"github.com/cespare/xxhash/v2"
)

type Stat struct {
	sizeInMem   uint64
	firstLoaded uint32
	visits      uint64
}

type Entry[T any] struct {
	Value T
	Stat
}

func (this *Entry[T]) Size() uint64 {
	if this == nil {
		return 0
	}

	if this.sizeInMem != 0 {
		return this.sizeInMem
	}

	sized, ok := any(this.Value).(interface{ MemSize() uint64 })
	if !ok {
		return 0
	}

	this.sizeInMem = sized.MemSize()
	return this.sizeInMem
}

type CachedKVStore[K comparable, T any] struct {
	*mapi.ConcurrentMap[K, *Entry[T]]
	backend          KVStore[K, *Entry[T]]
	currentLayerOnly bool
	cachePolicy      *CachePolicy[*Entry[T]]
	sizeOf           func(T) uint64

	pendingKeys []K
	pendingVals []*Entry[T]
}

func NewCachedKVStore[K comparable, T any](
	backend KVStore[K, *Entry[T]],
	cacheCap uint64,
	sizeOf func(T) uint64,
) *CachedKVStore[K, T] {
	store := &CachedKVStore[K, T]{
		ConcurrentMap: mapi.NewConcurrentMap(
			4096,
			func(v *Entry[T]) bool { return v == nil },
			func(k K) uint64 { return uint64(xxhash.Sum64String(fmt.Sprintf("%v", k))) },
		),
		backend: backend,
		sizeOf:  sizeOf,
	}
	store.cachePolicy = NewCachePolicy(cacheCap, store.entrySize)
	return store
}

func (this *CachedKVStore[K, T]) SetLocalOnly(yes bool) { this.currentLayerOnly = yes }
func (this *CachedKVStore[K, T]) LocalOnly() bool       { return this.currentLayerOnly }

func (this *CachedKVStore[K, T]) entrySize(entry *Entry[T]) uint64 {
	if entry == nil {
		return 0
	}

	if this != nil && this.sizeOf != nil {
		return this.sizeOf(entry.Value)
	}

	return entry.Size()
}

func (this *CachedKVStore[K, T]) Get(key K) (*Entry[T], bool) {
	for i := len(this.pendingKeys) - 1; i >= 0; i-- {
		if this.pendingKeys[i] != key {
			continue
		}
		return this.pendingVals[i], this.pendingVals[i] != nil
	}

	if v, ok := this.ConcurrentMap.Get(key); ok {
		return v, v != nil
	}

	if this.LocalOnly() || this.backend == nil {
		return nil, false
	}

	v, ok := this.backend.Get(key)
	if !ok || v == nil {
		return v, ok
	}

	if this.cachePolicy.Admit(this.cachePolicy.ValueSize(v)) {
		this.ConcurrentMap.Set(key, v)
	}
	return v, ok
}

func (this *CachedKVStore[K, T]) Set(key K, value *Entry[T]) {
	origin, ok := this.ConcurrentMap.Get(key)
	oldSize := uint64(0)
	if ok && origin != nil {
		oldSize = this.cachePolicy.ValueSize(origin)
	}

	if value == nil {
		this.cachePolicy.Remove(oldSize)
		this.ConcurrentMap.Set(key, nil)
		this.pendingKeys = append(this.pendingKeys, key)
		this.pendingVals = append(this.pendingVals, nil)
		return
	}

	this.cachePolicy.Track(oldSize, this.cachePolicy.ValueSize(value))
	this.ConcurrentMap.Set(key, value)
	this.pendingKeys = append(this.pendingKeys, key)
	this.pendingVals = append(this.pendingVals, value)
}

func (this *CachedKVStore[K, T]) Has(key K) bool {
	for i := len(this.pendingKeys) - 1; i >= 0; i-- {
		if this.pendingKeys[i] != key {
			continue
		}
		return this.pendingVals[i] != nil
	}

	if v, ok := this.ConcurrentMap.Get(key); ok {
		return v != nil
	}

	if this.LocalOnly() || this.backend == nil {
		return false
	}
	return this.backend.Has(key)
}

func (this *CachedKVStore[K, T]) GetBatch(keys []K) []*Entry[T] {
	if len(keys) == 0 {
		return nil
	}

	values := make([]*Entry[T], len(keys))
	missingIdx := make([]int, 0, len(keys))
	missingKeys := make([]K, 0, len(keys))

	for i, key := range keys {
		foundLocal := false
		for j := len(this.pendingKeys) - 1; j >= 0; j-- {
			if this.pendingKeys[j] != key {
				continue
			}
			values[i] = this.pendingVals[j]
			foundLocal = true
			break
		}

		if foundLocal {
			continue
		}

		if v, ok := this.ConcurrentMap.Get(key); ok {
			if v != nil {
				values[i] = v
			}
			continue
		}
		missingIdx = append(missingIdx, i)
		missingKeys = append(missingKeys, key)
	}

	if len(missingKeys) == 0 || this.LocalOnly() || this.backend == nil {
		return values
	}

	fetched := this.backend.GetBatch(missingKeys)
	limit := len(missingIdx)
	if len(fetched) < limit {
		limit = len(fetched)
	}

	for i := 0; i < limit; i++ {
		values[missingIdx[i]] = fetched[i]
		if fetched[i] == nil {
			continue
		}

		if this.cachePolicy.Admit(this.cachePolicy.ValueSize(fetched[i])) {
			this.ConcurrentMap.Set(missingKeys[i], fetched[i])
		}
	}
	return values
}

func (this *CachedKVStore[K, T]) SetBatch(keys []K, values []*Entry[T]) {
	limit := len(keys)
	if len(values) < limit {
		limit = len(values)
	}

	for i := 0; i < limit; i++ {
		origin, ok := this.ConcurrentMap.Get(keys[i])
		oldSize := uint64(0)
		if ok && origin != nil {
			oldSize = this.cachePolicy.ValueSize(origin)
		}

		if values[i] == nil {
			this.cachePolicy.Remove(oldSize)
			this.ConcurrentMap.Set(keys[i], nil)
			this.pendingKeys = append(this.pendingKeys, keys[i])
			this.pendingVals = append(this.pendingVals, nil)
			continue
		}

		this.cachePolicy.Track(oldSize, this.cachePolicy.ValueSize(values[i]))
		this.ConcurrentMap.Set(keys[i], values[i])
		this.pendingKeys = append(this.pendingKeys, keys[i])
		this.pendingVals = append(this.pendingVals, values[i])
	}
}

func (this *CachedKVStore[K, T]) Delete(key K) {
	this.Set(key, nil)
}

func (this *CachedKVStore[K, T]) DeleteBatch(keys []K) {
	this.SetBatch(keys, make([]*Entry[T], len(keys)))
}

func (this *CachedKVStore[K, T]) Len() uint64 {
	if this.backend != nil {
		return this.backend.Len()
	}
	return this.ConcurrentMap.Length()
}

func (this *CachedKVStore[K, T]) Size() uint64 {
	return this.cachePolicy.Size()
}

func (this *CachedKVStore[K, T]) evict() {
	if !this.cachePolicy.NeedEviction() {
		return
	}

	keys, vals := this.ConcurrentMap.KVs()

	for i, key := range keys {
		if !this.cachePolicy.NeedEviction() {
			return
		}
		if vals[i] == nil {
			continue
		}

		this.cachePolicy.Remove(this.cachePolicy.ValueSize(vals[i]))
		this.ConcurrentMap.Set(key, nil)
	}
}

func (this *CachedKVStore[K, T]) Precommit() error {
	return nil
}

func (this *CachedKVStore[K, T]) Commit(successful bool, version uint64) error {
	defer func() {
		this.pendingKeys = nil
		this.pendingVals = nil
	}()

	if !successful || this.backend == nil || len(this.pendingKeys) == 0 {
		return nil
	}

	this.backend.SetBatch(this.pendingKeys, this.pendingVals)
	this.evict()
	return nil
}
