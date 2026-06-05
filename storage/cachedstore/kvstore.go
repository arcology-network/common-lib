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

package cachedstore

import (
	"fmt"

	cache "github.com/arcology-network/common-lib/storage/cache"
	stgcodec "github.com/arcology-network/common-lib/storage/codec"
	stgintf "github.com/arcology-network/common-lib/storage/interface"
	"github.com/cespare/xxhash"
)

var _ stgintf.ReadWriteStore[string, any] = (*CachedStore[string, any, string, any])(nil)

type CachedStore[K0 stgintf.Key, V0 any, K1 stgintf.Key, V1 any] struct {
	cache     *cache.Cache[K0, V0]
	backend   stgintf.ReadWriteStore[K1, V1]
	converter *stgcodec.StorageCodec[K0, V0, K1, V1]
	decoder   func(K0, any, any) (any, error)
	zero      V0
}

func NewCachedStore[K0 stgintf.Key, V0 any, K1 stgintf.Key, V1 any](
	backend stgintf.ReadWriteStore[K1, V1],
	converter *stgcodec.StorageCodec[K0, V0, K1, V1],
	cacheCap uint64,
	sizeOf func(V0) uint64,
) *CachedStore[K0, V0, K1, V1] {
	store := &CachedStore[K0, V0, K1, V1]{
		cache: cache.NewCache(
			16,
			func(k K0) uint64 {
				return uint64(xxhash.Sum64String(fmt.Sprintf("%v", k)))
			},
			cache.NewCachePolicy(cacheCap, sizeOf),
		),
		backend:   backend,
		converter: converter,
	}
	return store
}

func (this *CachedStore[K0, V0, K1, V1]) Codec() *stgcodec.StorageCodec[K0, V0, K1, V1] {
	return this.converter
}

func (this *CachedStore[K0, V0, K1, V1]) Cache() *cache.Cache[K0, V0] { return this.cache }

func (this *CachedStore[K0, V0, K1, V1]) Preload([]byte) any { return nil }
func (this *CachedStore[K0, V0, K1, V1]) Backend() stgintf.ReadWriteStore[K1, V1] {
	return this.backend
}

// func (this *CachedStore[K0, V0, K1, V1]) Encoder(any) func(string, any) ([]byte, error) { return this.encoder }
// func (this *CachedStore[K0, V0, K1, V1]) Decoder(any) func(string, []byte, any) any     { return this.decoder }

func (this *CachedStore[K0, V0, K1, V1]) Has(key K0) bool {
	if ok := this.cache.Has(key); ok {
		return true
	}

	if this.backend == nil {
		return false
	}

	backendKey, _, err := this.converter.ForwardConvert(key, this.zero)
	if err != nil {
		return false
	}
	return this.backend.Has(backendKey)
}

func (this *CachedStore[K0, V0, K1, V1]) Get(key K0) (any, error) {
	return this.GetAs(key, nil)
}

func (this *CachedStore[K0, V0, K1, V1]) GetAs(key K0, typeHint any) (any, error) {
	if this == nil {
		return nil, stgintf.ErrNotFound
	}

	if record, err := this.cache.Get(key); err == nil {
		if this.decoder != nil {
			return this.decoder(key, record, typeHint)
		}
		return record, nil
	}

	if this.backend == nil {
		return this.zero, stgintf.ErrNotFound
	}

	backendKey, _, err := this.converter.ForwardConvert(key, this.zero)
	if err != nil {
		return this.zero, err
	}

	backendValue, err := this.backend.Get(backendKey)
	if err != nil {
		return this.zero, err
	}

	_, value, err := this.converter.BackwardConvert(backendKey, any(backendValue).(V1))
	if err != nil {
		return this.zero, err
	}

	// this.cache.Set(key, value)
	_ = this.cache.Set(key, value)
	return value, nil
}

func (this *CachedStore[K0, V0, K1, V1]) GetBatch(keys []K0) ([]any, []error) {
	if len(keys) == 0 {
		return nil, nil
	}

	values, errs := this.cache.GetBatch(keys)
	if this.backend == nil {
		return values, errs
	}

	for i := range errs {
		if errs[i] == nil {
			continue
		}

		backendKey, _, err := this.converter.ForwardConvert(keys[i], this.zero)
		if err != nil {
			errs[i] = err
			continue
		}

		backendVal, err := this.backend.Get(backendKey)
		if err != nil {
			errs[i] = err
			continue
		}

		_, value, err := this.converter.BackwardConvert(backendKey, any(backendVal).(V1))
		if err != nil {
			errs[i] = err
			continue
		}

		values[i] = value
		errs[i] = this.cache.Set(keys[i], value)
	}
	return values, errs
}

func (this *CachedStore[K0, V0, K1, V1]) Set(key K0, value V0) error {
	if err := this.cache.Set(key, value); err != nil {
		return err
	}

	backendKey, backendValue, err := this.converter.ForwardConvert(key, value)
	if err == nil && this.backend != nil {
		return this.backend.Set(backendKey, backendValue)
	}
	return nil
}

func (this *CachedStore[K0, V0, K1, V1]) SetBatch(keys []K0, values []V0) []error {
	this.cache.SetBatch(keys, values)
	errs := make([]error, len(keys))
	if this.backend == nil {
		return errs
	}

	backendKeys := make([]K1, 0, len(keys))
	backendVals := make([]V1, 0, len(keys))
	idxMap := make([]int, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		backendKey, backendVal, err := this.converter.ForwardConvert(keys[i], values[i])
		if err != nil {
			errs[i] = err
			continue
		}
		backendKeys = append(backendKeys, backendKey)
		backendVals = append(backendVals, backendVal)
		idxMap = append(idxMap, i)
	}

	if len(backendKeys) == 0 {
		return errs
	}
	backendErrs := this.backend.SetBatch(backendKeys, backendVals)
	for j, i := range idxMap {
		if j < len(backendErrs) {
			errs[i] = backendErrs[j]
		}
	}
	return errs
}

func (this *CachedStore[K0, V0, K1, V1]) Delete(key K0) error {
	if err := this.cache.Delete(key); err != nil {
		return err
	}

	backendKey, _, err := this.converter.ForwardConvert(key, this.zero)
	if err == nil && this.backend != nil {
		return this.backend.Delete(backendKey)
	}
	return nil
}

func (this *CachedStore[K0, V0, K1, V1]) DeleteBatch(keys []K0) []error {
	errs := make([]error, len(keys))
	for i, key := range keys {
		errs[i] = this.cache.Delete(key)
		backendKey, _, err := this.converter.ForwardConvert(key, this.zero)
		if err == nil && this.backend != nil {
			errs[i] = this.backend.Delete(backendKey)
		}
	}
	return errs
}

func (this *CachedStore[K0, V0, K1, V1]) Query(target K0, predicate func(K0, V0) bool) ([]K0, []V0, []error) {
	cacheKeys, cacheValues, cacheErrs := this.cache.Query(target, predicate)

	type backendQuerier[K stgintf.Key, V any] interface {
		Query(K, func(K, V) bool) ([]K, []V, []error)
	}

	backendQuery, ok := any(this.backend).(backendQuerier[K1, V1])
	if !ok || this.backend == nil {
		return cacheKeys, cacheValues, cacheErrs
	}

	backendPredicate := func(key K1, value V1) bool {
		decodedKey, decodedValue, err := this.converter.BackwardConvert(key, value)
		if err != nil {
			return false
		}
		if predicate == nil {
			return decodedKey == target
		}
		return predicate(decodedKey, decodedValue)
	}

	backendTarget, _, err := this.converter.ForwardConvert(target, this.zero)
	if err != nil {
		return cacheKeys, cacheValues, append(cacheErrs, err)
	}

	backendKeys, backendValues, backendErrs := backendQuery.Query(backendTarget, backendPredicate)
	seen := make(map[K0]struct{}, len(cacheKeys))
	for _, key := range cacheKeys {
		seen[key] = struct{}{}
	}

	for i := 0; i < len(backendKeys) && i < len(backendValues); i++ {
		decodedKey, decodedValue, decodeErr := this.converter.BackwardConvert(backendKeys[i], backendValues[i])
		if decodeErr != nil {
			backendErrs = append(backendErrs, decodeErr)
			continue
		}
		if _, exists := seen[decodedKey]; exists {
			continue
		}
		cacheKeys = append(cacheKeys, decodedKey)
		cacheValues = append(cacheValues, decodedValue)
		seen[decodedKey] = struct{}{}
	}

	if len(cacheErrs) == 0 {
		return cacheKeys, cacheValues, backendErrs
	}
	return cacheKeys, cacheValues, append(cacheErrs, backendErrs...)
}
