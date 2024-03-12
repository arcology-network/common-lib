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
	nilK        K
	mapper      func(K) uint64
	singleCache map[K]*T
}

func NewReadCache[K comparable, T any](numShards uint64, mapper func(K) uint64, nilK K) *ReadCache[K, T] {
	newReadCache := &ReadCache[K, T]{
		nilK:        nilK,
		mapper:      mapper,
		singleCache: make(map[K]*T),
	}
	return newReadCache
}

func (this *ReadCache[K, T]) Length() int { return len(this.singleCache) }

func (this *ReadCache[K, T]) Get(k K) (*T, bool) {
	v, ok := this.singleCache[k]
	if !ok {
		return nil, false
	}
	return v, ok
}

func (this *ReadCache[K, T]) Commit(keys []K, values []T) {
	for i, k := range keys {
		this.singleCache[k] = &values[i]
	}
	this.Clear()
}

// Call this function to clear the cache.
func (this *ReadCache[K, T]) Clear() {}
