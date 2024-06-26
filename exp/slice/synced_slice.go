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

package slice

import "sync"

type SyncedSlice[T any] struct {
	values []T
	lock   sync.Mutex
}

// Start starts the goroutines.
func NewSlice[T any]() *SyncedSlice[T] {
	return &SyncedSlice[T]{
		values: []T{},
		lock:   sync.Mutex{},
	}
}

// Start starts the goroutines.
func (this *SyncedSlice[T]) Get(idx int) T {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.values[idx]
}

// Start starts the goroutines.
func (this *SyncedSlice[T]) Append(v ...T) *SyncedSlice[T] {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.values = append(this.values, v...)
	return this
}

func (this *SyncedSlice[T]) MoveToSlice() []T {
	this.lock.Lock()
	defer this.lock.Unlock()

	values := this.values
	this.values = this.values[:0]
	return values
}

func (this *SyncedSlice[T]) ToSlice() []T {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.values
}

func (this *SyncedSlice[T]) Length() int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return len(this.values)
}

func (this *SyncedSlice[T]) Clear() {
	this.values = this.values[:0]
}
