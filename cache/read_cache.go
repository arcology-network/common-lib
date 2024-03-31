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

// ReadCache is a read only cache that is used to store the read values from the storage.
// The cache updates itself when the update is called. The implementation isn't thread safe.
// So, it's the caller's responsibility to ensure that the cache is only accessed by one thread updating it.
// Each entry in the cache holds two values, the first value is the old value, and the second value is the new value.
// The new value will be set to the old value when the Finalize function is called.
type ReadCache[K comparable, T any] struct {
	isNil   func(T) bool
	cache   map[K]*T
	enabled bool
}

func NewReadCache[K comparable, T any](numShards uint64, isNil func(T) bool) *ReadCache[K, T] {
	newReadCache := &ReadCache[K, T]{
		isNil:   isNil,
		cache:   make(map[K]*T),
		enabled: true,
	}
	return newReadCache
}

func (this *ReadCache[K, T]) Status() bool { return this.enabled }
func (this *ReadCache[K, T]) Enable()      { this.enabled = true }
func (this *ReadCache[K, T]) Disable()     { this.enabled = false }

func (this *ReadCache[K, T]) Length() int {
	if !this.enabled {
		return 0
	}
	return len(this.cache)
}

func (this *ReadCache[K, T]) Get(k K) (*T, bool) {
	if !this.enabled {
		return nil, false
	}

	v, ok := this.cache[k]
	if !ok {
		return nil, false
	}
	return v, ok
}

func (this *ReadCache[K, T]) Precommit(keys []K, values []T) {
	if !this.enabled {
		return
	}

	for i, k := range keys {
		if this.isNil(values[i]) {
			delete(this.cache, k)
		} else {
			this.cache[k] = &values[i]
		}
	}
}

// Call this function to clear the cache completely.
func (this *ReadCache[K, T]) Clear() { clear(this.cache) }
