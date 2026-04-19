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
	"reflect"
	"sync/atomic"

	mapi "github.com/arcology-network/common-lib/exp/map"
	"github.com/arcology-network/common-lib/exp/slice"
	"github.com/cespare/xxhash/v2"

	stgintf "github.com/arcology-network/common-lib/storage/interface"
)

type Stat struct {
	sizeInMem   uint64
	firstLoaded uint64
	visits      uint64
}

func (this *Stat) SetLoaded(version uint64) {
	this.firstLoaded = version
}

type entry[T any] struct {
	value T
	Stat
}

func (this *entry[T]) Size() uint64 {
	if this == nil {
		return 0
	}

	if this.sizeInMem != 0 {
		return this.sizeInMem
	}

	sized, ok := any(this.value).(interface{ MemSize() uint64 })
	if !ok {
		return 0
	}

	this.sizeInMem = sized.MemSize()
	return this.sizeInMem
}

type CachedKVStore[K comparable, T any] struct {
	cache            *mapi.ConcurrentMap[K, *entry[T]]
	backend          stgintf.KVStore[K, T]
	currentLayerOnly bool
	cachePolicy      *CachePolicy[*entry[T]]
	sizeOf           func(T) uint64
	version          atomic.Uint64
}

func NewCachedKVStore[K comparable, T any](
	backend stgintf.KVStore[K, T],
	cacheCap uint64,
	sizeOf func(T) uint64,
) *CachedKVStore[K, T] {
	store := &CachedKVStore[K, T]{
		cache: mapi.NewConcurrentMap(
			4096,
			func(v *entry[T]) bool { return v == nil },
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

func (this *CachedKVStore[K, T]) UpdateVersion(version uint64) {
	if this == nil {
		return
	}
	this.version.Store(version)
}

func (this *CachedKVStore[K, T]) entrySize(entry *entry[T]) uint64 {
	if entry == nil {
		return 0
	}

	if this != nil && this.sizeOf != nil {
		return this.sizeOf(entry.value)
	}

	return entry.Size()
}

func (this *CachedKVStore[K, T]) wrap(value T) *entry[T] {
	entry := &entry[T]{value: value}
	entry.firstLoaded = this.version.Load()
	entry.visits++
	return entry
}

func isNil[T any](value T) bool {
	if any(value) == nil {
		return true
	}

	rv := reflect.ValueOf(any(value))
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func (this *CachedKVStore[K, T]) Get(key K) (T, bool) {
	if v, ok := this.cache.Get(key); ok {
		if v != nil {
			v.visits++
			return v.value, true
		}
		var zero T
		return zero, false
	}

	if this.LocalOnly() || this.backend == nil {
		var zero T
		return zero, false
	}

	v, ok := this.backend.Get(key)
	if !ok || isNil(v) {
		var zero T
		return zero, false
	}

	entry := this.wrap(v)
	if this.cachePolicy.Admit(this.cachePolicy.ValueSize(entry)) {
		this.cache.Set(key, entry)
	}
	return entry.value, true
}

func (this *CachedKVStore[K, T]) Set(key K, value T) {
	origin, ok := this.cache.Get(key)
	oldSize := uint64(0)
	if ok && origin != nil {
		oldSize = this.cachePolicy.ValueSize(origin)
	}

	if isNil(value) {
		this.cachePolicy.Remove(oldSize)
		this.cache.Set(key, nil)
		return
	}
	entry := this.wrap(value)
	if origin != nil {
		entry.visits += origin.visits
		entry.firstLoaded = origin.firstLoaded
	}

	this.cachePolicy.Track(oldSize, this.cachePolicy.ValueSize(entry))
	this.cache.Set(key, entry)
}

func (this *CachedKVStore[K, T]) Has(key K) bool {
	if v, ok := this.cache.Get(key); ok {
		return v != nil
	}

	if this.LocalOnly() || this.backend == nil {
		return false
	}
	return this.backend.Has(key)
}

func (this *CachedKVStore[K, T]) GetBatch(keys []K) []T {
	if len(keys) == 0 {
		return nil
	}

	values := make([]T, len(keys))
	missingIdx := make([]int, 0, len(keys))
	missingKeys := make([]K, 0, len(keys))

	for i, key := range keys {
		if v, ok := this.cache.Get(key); ok {
			if v != nil {
				v.visits++
				values[i] = v.value
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
		if isNil(fetched[i]) {
			continue
		}

		entry := this.wrap(fetched[i])
		values[missingIdx[i]] = entry.value
		if this.cachePolicy.Admit(this.cachePolicy.ValueSize(entry)) {
			this.cache.Set(missingKeys[i], entry)
		}
	}
	return values
}

func (this *CachedKVStore[K, T]) SetBatch(keys []K, values []T) {
	for i := 0; i < len(keys); i++ {
		origin, ok := this.cache.Get(keys[i])
		oldSize := uint64(0)
		if ok && origin != nil {
			oldSize = this.cachePolicy.ValueSize(origin)
		}

		if isNil(values[i]) {
			this.cachePolicy.Remove(oldSize)
			this.cache.Set(keys[i], nil)
			continue
		}

		entry := this.wrap(values[i])
		if origin != nil {
			entry.visits += origin.visits
			entry.firstLoaded = origin.firstLoaded
		}

		this.cachePolicy.Track(oldSize, this.cachePolicy.ValueSize(entry))
		this.cache.Set(keys[i], entry)
	}
}

func (this *CachedKVStore[K, T]) deleteFromCache(key K) {
	if origin, ok := this.cache.Get(key); ok && origin != nil {
		this.cachePolicy.Remove(this.cachePolicy.ValueSize(origin))
	}
	this.cache.Set(key, nil)
}

func (this *CachedKVStore[K, T]) Delete(key K) error {
	this.deleteFromCache(key)
	return nil
}

func (this *CachedKVStore[K, T]) DeleteBatch(keys []K) error {
	for _, key := range keys {
		this.deleteFromCache(key)
	}
	return nil
}

func (this *CachedKVStore[K, T]) Len() uint64 {
	if this.backend != nil {
		return this.backend.Len()
	}
	return this.cache.Length()
}

func (this *CachedKVStore[K, T]) Size() uint64 {
	return this.cachePolicy.Size()
}

func (this *CachedKVStore[K, T]) Evict() {
	if !this.cachePolicy.NeedEviction() {
		return
	}

	keys, vals := this.cache.KVs()

	for i, key := range keys {
		if !this.cachePolicy.NeedEviction() {
			return
		}
		if vals[i] == nil {
			continue
		}

		this.cachePolicy.Remove(this.cachePolicy.ValueSize(vals[i]))
		this.cache.Set(key, nil)
	}
}

func (this *CachedKVStore[K, T]) Clear() { this.cache.Clear() }

func (this *CachedKVStore[K, T]) Policy() *CachePolicy[*entry[T]]  { return this.cachePolicy }
func (this *CachedKVStore[K, T]) Profile() *CachePolicy[*entry[T]] { return this.cachePolicy }

func (this *CachedKVStore[K, T]) CacheChecksum() [32]byte {
	encoders := func(k K, v *entry[T]) ([]byte, []byte) {
		key := []byte(fmt.Sprintf("%v", k))
		if v == nil {
			return key, nil
		}

		encoder, ok := any(v.value).(interface{ Encode() []byte })
		if !ok {
			return key, nil
		}
		return key, encoder.Encode()
	}

	less := func(k0, k1 K) bool {
		return fmt.Sprintf("%v", k0) < fmt.Sprintf("%v", k1)
	}
	return this.cache.Checksum(less, encoders)
}

func (this *CachedKVStore[K, T]) getRaw(key K) (T, bool) {
	v, ok := this.cache.Get(key)
	if !ok || v == nil {
		var zero T
		return zero, false
	}
	return v.value, true
}

func (this *CachedKVStore[K, T]) Print() {
	keys, vals := this.cache.KVs()
	slice.SortBy1st(keys, vals, func(k0, k1 K) bool {
		return fmt.Sprintf("%v", k0) < fmt.Sprintf("%v", k1)
	})

	fmt.Println("occupied:", this.cachePolicy.occupied)

	for i, k := range keys {
		println(k, "  =    ")
		printer, ok := any(vals[i].value).(interface{ Print() })
		if ok {
			printer.Print()
		}
	}
}
