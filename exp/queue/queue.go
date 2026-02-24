/*
 *   Copyright (c) 2025 Arcology Network

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

import (
	"sort"

	"github.com/arcology-network/common-lib/exp/slice"
)

type Queue[T any] []T

func NewQueue[T any]() *Queue[T] {
	q := &Queue[T]{}
	return q
}

func NewSortedQueueFromSlice[T any](items []T, less func(a, b T) bool) *Queue[T] {
	q := &Queue[T]{}
	*q = append(*q, items...)
	q.Sort(less)
	return q
}

func NewQueueFromSlice[T any](items []T) *Queue[T] {
	q := &Queue[T]{}
	*q = append(*q, items...)
	return q
}

func (q *Queue[T]) Enqueue(item T) {
	*q = append(*q, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var nilElem T
	if q.IsEmpty() {
		return nilElem, false
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item, true
}

func (q *Queue[T]) Peek() (T, bool) {
	var nilElem T
	if q.IsEmpty() {
		return nilElem, false
	}
	item := (*q)[0]
	return item, true
}

func (q *Queue[T]) Back() (T, bool) {
	var nilElem T
	if q.IsEmpty() {
		return nilElem, false
	}
	item := (*q)[len(*q)-1]
	return item, true
}

func (q *Queue[T]) Clear() {
	*q = (*q)[:0]
}

func (q *Queue[T]) ToSlice() []T {
	return *q
}

func (q *Queue[T]) Sort(less func(a, b T) bool) {
	sort.Slice(*q, func(i, j int) bool {
		return less((*q)[i], (*q)[j])
	})
}

func (q *Queue[T]) IsEmpty() bool {
	return len(*q) == 0
}

func (q *Queue[T]) Size() int {
	return len(*q)
}

func (q *Queue[T]) RemoveIf(pred func(idx int, item T) bool) {
	s := []T(*q)
	slice.RemoveIf(&s, pred)
	*q = Queue[T](s)
}

func (q *Queue[T]) Clone() *Queue[T] {
	cloned := make(Queue[T], len(*q))
	copy(cloned, *q)
	return &cloned
}

func (q *Queue[T]) CloneDo(cloneFunc func(T) T) *Queue[T] {
	cloned := make(Queue[T], len(*q))
	for i, item := range *q {
		cloned[i] = cloneFunc(item)
	}
	return &cloned
}
