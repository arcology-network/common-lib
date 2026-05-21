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

import (
	"errors"

	"golang.org/x/exp/constraints"
)

var ErrNotFound = errors.New("not found")
var ErrNotInParent = errors.New("not in parent")
var ErrNoFallBack = errors.New("no fallback available")

const (
	MEMORY_DB     = 0
	PERSISTENT_DB = 1
)

type Key interface {
	~string | constraints.Integer
}

// ReadOnlyStore defines the interface for a read-only storage source.
type ReadOnlyStore[K comparable, V any] interface {
	Has(K) bool         // Check if the key exists in the source, which can be a cache or a storage.
	Get(K) (any, error) // Get from cache or persistent storage, with cache lookup first.
	GetAs(K, any) (any, error)
}

type ReadableStore[K comparable, V any] interface {
	ReadOnlyStore[K, V]
	GetBatch([]K) ([]any, []error)
}

type WriteableStore[K comparable, V any] interface {
	Set(K, V) error
	SetBatch([]K, []V) []error
	Delete(K) error
	DeleteBatch([]K) []error
}

type ReadWriteStore[K comparable, V any] interface {
	ReadableStore[K, V]
	WriteableStore[K, V]
	Query(K, func(K, V) bool) ([]K, []V, []error)
}

type StoreWriter[T any] interface {
	Import([]T)
	Precommit(bool) error //should return a error
	Commit(uint64) error  //should return a error

	IsSync() bool // If the writer is synchronous, it will block until the commit is done.
	Name() string
}

type BackendStore[K Key, V any] interface {
	Get(K) (any, error)
	GetAs(K, any) (any, error)
	GetBatch([]K) ([]any, []error)
	Set(K, V) error
	SetBatch([]K, []V) []error
	Delete(K) error
	DeleteBatch([]K) []error
	Has(K) bool
}

// type CacheLike[K comparable, V any] interface {
// 	Len() uint64
// 	Size() uint64
// 	IsEnabled() bool // If the cache is enabled
// 	Enable() bool
// 	Disable() bool
// }
