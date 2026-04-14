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

package interfaces

const (
	MEMORY_DB     = 0
	PERSISTENT_DB = 1
)

type Readable[K comparable, V any] interface {
	Get(K) (V, bool)
	GetBatch([]K) []V
	Has(K) bool
	Len() uint64
	Size() uint64
}

type Writeable[K comparable, V any] interface {
	Set(K, V)
	SetBatch([]K, []V)
	Delete(K)
	DeleteBatch([]K)
}

type ReadWriteStore[K comparable, V any] interface {
	Readable[K, V]
	Writeable[K, V]
}

type KVStore[K comparable, V any] interface {
	ReadWriteStore[K, V]
	SetLocalOnly(yes bool)
	LocalOnly() bool
}

type PersistentStorage[K comparable, T any] interface {
	Get(K) (T, error)
	Set(K, T) error
	GetBatch([]K) ([]T, error)
	SetBatch([]K, []T) error
	Query(K, func(K, T) bool) ([]K, []T, error)
}

// ReadOnlyStore defines the interface for a read-only storage source.
type ReadOnlyStore[K comparable, T any] interface {
	Has(K) bool                    // Check if the key exists in the source, which can be a cache or a storage.
	ReadBackend(K, T) (any, error) // Get from persistent storage directly.
	GetAs(K, T) (any, error)       // Get from cache or persistent storage, with cache lookup first.
	Preload([]byte) any
}
