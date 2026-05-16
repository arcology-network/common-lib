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

package cache

import (
	"sync/atomic"

	mapi "github.com/arcology-network/common-lib/exp/map"
	stgintf "github.com/arcology-network/common-lib/storage/interface"
)

var _ stgintf.ReadWriteStore[string, any] = (*Cache[string, any])(nil)

// Cache is a read only cache that is used to store the read values from the storage.
// The cache updates itself when the update is called. The implementation isn't thread safe.
// So, it's the caller's responsibility to ensure that the cache is only accessed by one thread updating it.
// Each entry in the cache holds two values, the first value is the old value, and the second value is the new value.
// The new value will be set to the old value when the Finalize function is called.
type Cache[K stgintf.Key, V any] struct {
	*mapi.ConcurrentMap[K, *entry[V]]
	cachePolicy *CachePolicy[V]
	epoch   atomic.Uint64
	enabled bool
}

func NewCache[K stgintf.Key, V any](
	numShards uint64,
	hasher func(K) uint64,
	cachePolicy *CachePolicy[V],
) *Cache[K, V] {
	newReadCache := &Cache[K, V]{
		ConcurrentMap: mapi.NewConcurrentMap(
			int(numShards),
			func(e *entry[V]) bool { return e == nil },
			hasher,
		),
		cachePolicy: cachePolicy,
		enabled:     true,
	}
	return newReadCache
}

func (this *Cache[K, V]) Has(key K) bool {
	if !this.enabled {
		return false
	}
	_, ok := this.ConcurrentMap.Get(key)
	return ok && this.enabled
}

func (this *Cache[K, V]) Get(key K) (any, error) {
	if this.enabled {
		if record, ok := this.ConcurrentMap.Get(key); ok {
			if record != nil {
				record.visits++
				return record.value, nil
			}
		}
	}
	return nil, stgintf.ErrNotFound
}

func (this *Cache[K, V]) GetBatch(keys []K) ([]any, []error) {
	if !this.enabled {
		errs := make([]error, len(keys))
		for i := range keys {
			errs[i] = stgintf.ErrNotFound
		}
		return nil, errs
	}

	if len(keys) == 0 {
		return nil, nil
	}

	values := make([]any, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		if record, ok := this.ConcurrentMap.Get(key); ok && record != nil {
			record.visits++
			values[i] = record.value
			continue
		}
		errs[i] = stgintf.ErrNotFound
	}
	return values, errs
}

func (this *Cache[K, V]) Set(key K, value V) error {
	if !this.enabled {
		return nil
	}

	origin, ok := this.ConcurrentMap.Get(key)
	if ok {
		oldSize, newSize := origin.Replace(value)
		this.cachePolicy.Update(oldSize, newSize)
		return nil
	}

	newVSize := this.cachePolicy.ValueSize(value)
	if this.cachePolicy.Update(0, newVSize) {
		entry := this.wrap(value)
		this.ConcurrentMap.Set(key, entry)
	}
	return nil
}

func (this *Cache[K, V]) SetBatch(keys []K, values []V) []error {
	errs := make([]error, len(keys))
	if !this.enabled {
		return errs
	}

	for i := 0; i < len(keys); i++ {
		if i < len(values) {
			errs[i] = this.Set(keys[i], values[i])
		}
	}
	return errs
}

func (this *Cache[K, V]) Delete(key K) error {
	if !this.enabled {
		return nil
	}

	if origin, ok := this.ConcurrentMap.Get(key); ok {
		if origin != nil {
			oldSize := this.cachePolicy.ValueSize(origin.value)
			this.cachePolicy.Update(oldSize, 0)
		}
		this.ConcurrentMap.Delete(key)
	}
	return nil
}

func (this *Cache[K, V]) DeleteBatch(keys []K) []error {
	errs := make([]error, len(keys))
	if !this.enabled {
		return errs
	}

	for i, key := range keys {
		errs[i] = this.Delete(key)
	}
	return errs
}

func (this *Cache[K, V]) Query(target K, predicate func(K, V) bool) ([]K, []V, []error) {
	if !this.enabled {
		return nil, nil, nil
	}

	keys, entries := this.ConcurrentMap.KVs()
	matchedKeys := make([]K, 0, len(keys))
	matchedValues := make([]V, 0, len(keys))

	for i, key := range keys {
		entry := entries[i]
		if entry == nil {
			continue
		}

		value := entry.value
		if predicate != nil {
			if !predicate(key, value) {
				continue
			}
		} else if key != target {
			continue
		}

		matchedKeys = append(matchedKeys, key)
		matchedValues = append(matchedValues, value)
	}
	return matchedKeys, matchedValues, nil
}

func (this *Cache[K, V]) Evict() {
	if !this.enabled {
		return
	}

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

		this.cachePolicy.Update(this.cachePolicy.ValueSize(vals[i].value), 0)
		this.ConcurrentMap.Delete(key)
	}
}

func (this *Cache[K, V]) Status() bool            { return this.enabled }
func (this *Cache[K, V]) SetStatus(flag bool)     { this.enabled = flag }
func (this *Cache[K, V]) Hash(k K) uint64         { return this.ConcurrentMap.Hash(k) }
func (this *Cache[K, V]) Cap() uint64             { return this.cachePolicy.Size() }
func (this *Cache[K, V]) Clear()                  { this.ConcurrentMap.Clear() }
func (this *Cache[K, V]) Policy() *CachePolicy[V] { return this.cachePolicy }

// func (this *Cache[K, V]) entrySize(entry *entry[V]) uint64 {
// 	if entry == nil {
// 		return 0
// 	}

// 	if this != nil && this.sizeOf != nil {
// 		return this.sizeOf(entry.value)
// 	}

// 	return entry.Size()
// }

func (this *Cache[K, V]) wrap(value V) *entry[V] {
	entry := &entry[V]{value: value}
	entry.firstLoaded = this.epoch.Load()
	entry.visits++
	return entry
}
