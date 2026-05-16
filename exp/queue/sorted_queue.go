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

package queue

import "sort"

type SortedQueue[T any] struct {
	items Queue[T]
	less  func(a, b T) bool
}

func NewSortedQueue[T any](less func(a, b T) bool) *SortedQueue[T] {
	if less == nil {
		panic("less func is nil")
	}
	return &SortedQueue[T]{less: less}
}

func NewSortedQueueWithItems[T any](items []T, less func(a, b T) bool) *SortedQueue[T] {
	q := NewSortedQueue(less)
	q.items = append(q.items, items...)
	q.items.Sort(less)
	return q
}

func (q *SortedQueue[T]) Enqueue(item T) int {
	idx := q.Greater(item)
	q.items = append(q.items, item)
	copy(q.items[idx+1:], q.items[idx:])
	q.items[idx] = item
	return idx
}

func (q *SortedQueue[T]) Dequeue() (T, bool) {
	return q.items.Dequeue()
}

func (q *SortedQueue[T]) Peek() (T, bool) {
	return q.items.Peek()
}

func (q *SortedQueue[T]) Back() (T, bool) {
	return q.items.Back()
}

func (q *SortedQueue[T]) Clear() {
	q.items.Clear()
}

func (q *SortedQueue[T]) ToSlice() []T {
	return q.items.ToSlice()
}

func (q *SortedQueue[T]) IsEmpty() bool {
	return q.items.IsEmpty()
}

func (q *SortedQueue[T]) Size() int {
	return q.items.Size()
}

func (q *SortedQueue[T]) Greater(item T) int {
	return sort.Search(len(q.items), func(i int) bool {
		return q.less(item, q.items[i])
	})
}

func (q *SortedQueue[T]) GreaterOrEqual(item T) int {
	return sort.Search(len(q.items), func(i int) bool {
		return !q.less(q.items[i], item)
	})
}

func (q *SortedQueue[T]) Less(item T) int {
	return q.GreaterOrEqual(item) - 1
}

func (q *SortedQueue[T]) LessOrEqual(item T) int {
	return q.Greater(item) - 1
}

func (q *SortedQueue[T]) Index(item T) (int, bool) {
	idx := q.GreaterOrEqual(item)
	if idx >= len(q.items) {
		return -1, false
	}

	if q.less(item, q.items[idx]) || q.less(q.items[idx], item) {
		return -1, false
	}
	return idx, true
}
