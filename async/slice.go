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

package async

import "sync"

type Slice[T any] struct {
	values []T
	lock   sync.Mutex
}

// Start starts the goroutines.
func NewSlice[T any]() *Slice[T] {
	return &Slice[T]{
		values: []T{},
		lock:   sync.Mutex{},
	}
}

// Start starts the goroutines.
func (this *Slice[T]) Get(idx int) T {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.values[idx]
}

// Start starts the goroutines.
func (this *Slice[T]) Append(v T) *Slice[T] {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.values = append(this.values, v)
	return this
}

func (this *Slice[T]) MoveToSlice() *Slice[T] {
	this.lock.Lock()
	defer this.lock.Unlock()

	slice := &Slice[T]{
		values: this.values,
	}
	this.values = this.values[:0]
	return slice
}
