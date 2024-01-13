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

// MemoryPool is responsible for managing a pool of allocated of the same type. It is thread-safe.
package mempool

import (
	"sync"

	"github.com/arcology-network/common-lib/container/array"
)

// MemoryPool represents a pool of allocated of the same type.
type MemoryPool[T any] struct {
	new       func() T
	allocated *array.PagedArray[T]
	counter   int
	lock      sync.Mutex
}

// NewMemoryPool creates a new MemoryPool instance with the given ID and object creation function.
func NewMemoryPool[T any](perPage, pages int, new func() T) *MemoryPool[T] {
	pool := &MemoryPool[T]{
		new:       new,
		allocated: array.NewPagedArray[T](perPage, pages, perPage*pages),
		counter:   0,
	}

	pool.allocated.Foreach(func(i int, v *T) {
		pool.allocated.Set(i, new())
	})
	return pool
}

// New returns an object from the MemoryPool.
func (this *MemoryPool[T]) New() T {
	if this.counter >= this.allocated.Size() {
		this.allocated.PushBack(this.new())
	}
	v := this.allocated.Get(this.counter)
	this.counter++
	return v
}

// Reclaim resets the MemoryPool, allowing the allocated to be reused.
func (this *MemoryPool[T]) Reclaim() {
	this.counter = 0
	this.allocated.Resize(this.allocated.MinSize())
}
