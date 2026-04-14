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
	backend          stgintf.KVStore[K, T]
	currentLayerOnly bool
	cachePolicy      *CachePolicy[*Entry[T]]
	sizeOf           func(T) uint64
	version          atomic.Uint64
}

func NewCachedKVStore[K comparable, T any](
	backend stgintf.KVStore[K, T],
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

func (this *CachedKVStore[K, T]) UpdateVersion(version uint64) {
	if this == nil {
		return
	}
	this.version.Store(version)
}

func (this *CachedKVStore[K, T]) entrySize(entry *Entry[T]) uint64 {
	if entry == nil {
		return 0
	}

	if this != nil && this.sizeOf != nil {
		return this.sizeOf(entry.Value)
	}

	return entry.Size()
}

func (this *CachedKVStore[K, T]) wrap(value T) *Entry[T] {
	entry := &Entry[T]{Value: value}
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

func (this *CachedKVStore[K, T]) Get(key K) (*Entry[T], bool) {
	if v, ok := this.ConcurrentMap.Get(key); ok {
		if v != nil {
			v.visits++
		}
		return v, v != nil
	}

	if this.LocalOnly() || this.backend == nil {
		return nil, false
	}

	v, ok := this.backend.Get(key)
	if !ok || isNil(v) {
		return nil, false
	}

	entry := this.wrap(v)
	if this.cachePolicy.Admit(this.cachePolicy.ValueSize(entry)) {
		this.ConcurrentMap.Set(key, entry)
	}
	return entry, true
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
		return
	}
	if origin != nil && origin != value {
		value.visits += origin.visits
		if value.firstLoaded == 0 {
			value.firstLoaded = origin.firstLoaded
		}
	}
	if origin == nil && value.firstLoaded == 0 {
		value.firstLoaded = this.version.Load()
	}
	value.visits++

	this.cachePolicy.Track(oldSize, this.cachePolicy.ValueSize(value))
	this.ConcurrentMap.Set(key, value)
}

func (this *CachedKVStore[K, T]) Has(key K) bool {
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
		if v, ok := this.ConcurrentMap.Get(key); ok {
			if v != nil {
				v.visits++
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
		if isNil(fetched[i]) {
			continue
		}

		entry := this.wrap(fetched[i])
		values[missingIdx[i]] = entry
		if this.cachePolicy.Admit(this.cachePolicy.ValueSize(entry)) {
			this.ConcurrentMap.Set(missingKeys[i], entry)
		}
	}
	return values
}

func (this *CachedKVStore[K, T]) SetBatch(keys []K, values []*Entry[T]) {
	for i := 0; i < len(keys); i++ {
		origin, ok := this.ConcurrentMap.Get(keys[i])
		oldSize := uint64(0)
		if ok && origin != nil {
			oldSize = this.cachePolicy.ValueSize(origin)
		}

		if values[i] == nil {
			this.cachePolicy.Remove(oldSize)
			this.ConcurrentMap.Set(keys[i], nil)
			continue
		}
		if origin != nil && origin != values[i] {
			values[i].visits += origin.visits
			if values[i].firstLoaded == 0 {
				values[i].firstLoaded = origin.firstLoaded
			}
		}
		if origin == nil && values[i].firstLoaded == 0 {
			values[i].firstLoaded = this.version.Load()
		}
		values[i].visits++

		this.cachePolicy.Track(oldSize, this.cachePolicy.ValueSize(values[i]))
		this.ConcurrentMap.Set(keys[i], values[i])
	}
}

func (this *CachedKVStore[K, T]) Delete(key K) error {
	this.Set(key, nil)
	return nil
}

func (this *CachedKVStore[K, T]) DeleteBatch(keys []K) error {
	this.SetBatch(keys, make([]*Entry[T], len(keys)))
	return nil
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

func (this *CachedKVStore[K, T]) Evict() {
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

func (this *CachedKVStore[K, T]) GetRaw(key K) (*Entry[T], bool) {
	v, ok := this.ConcurrentMap.Get(key)
	if !ok || v == nil {
		return nil, false
	}
	return v, true
}

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
