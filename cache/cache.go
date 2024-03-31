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
	"github.com/arcology-network/common-lib/common"
	ccmap "github.com/arcology-network/common-lib/exp/map"
)

// RWCache is a read only cache that is used to store the read values from the storage.
// The cache updates itself when the update is called. The implementation isn't thread safe.
// So, it's the caller's responsibility to ensure that the cache is only accessed by one thread updating it.
// Each entry in the cache holds two values, the first value is the old value, and the second value is the new value.
// The new value will be set to the old value when the Finalize function is called.
type RWCache[K comparable, V any] struct {
	isNil   func(V) bool
	cache   *ccmap.ConcurrentMap[K, V]
	enabled bool
}

func NewRWCache[K comparable, V any](numShards uint64, isNil func(V) bool, hasher func(K) uint64) *RWCache[K, V] {
	return &RWCache[K, V]{
		isNil:   isNil,
		cache:   ccmap.NewConcurrentMap[K, V](common.Min(int(numShards), 256), isNil, hasher),
		enabled: true,
	}
}

func (this *RWCache[K, V]) Status() bool { return this.enabled }
func (this *RWCache[K, V]) Enable()      { this.enabled = true }
func (this *RWCache[K, V]) Disable()     { this.enabled = false }

func (this *RWCache[K, V]) Length() int {
	if !this.enabled {
		return 0
	}
	return int(this.cache.Length())
}

func (this *RWCache[K, V]) Get(k K) (*V, bool) {
	if !this.enabled {
		return nil, false
	}

	v, ok := this.cache.Get(k)
	if !ok {
		return nil, false
	}
	return &v, ok
}

func (this *RWCache[K, V]) Commit(keys []K, values []V) {
	if !this.enabled {
		return
	}
	this.cache.BatchSet(keys, values)
}

// Call this function to clear the cache completely.
func (this *RWCache[K, V]) Clear() {
	this.cache.Clear()
}
