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

// The PagedSlice class is a custom data structure that represents an array that is divided into multiple pages or pages.
// It is designed to efficiently handle large arrays by storing the elements in a paginated manner to optimize memory usage
// and improve performance.
package indexedslice

import (
	"math"

	"github.com/arcology-network/common-lib/common"
)

type PagedSlice[T any] struct {
	minPages int
	pageSize int
	length   int
	pages    [][]T
}

// NewPagedSlice creates a new instance of PagedSlice with the specified block size and minimum number of pages.
func NewPagedSlice[T any](pageSize, minPages, preAlloc int) *PagedSlice[T] {
	if pageSize*minPages < preAlloc {
		panic("preAlloc must be less than pageSize * minPages")
	}

	array := &PagedSlice[T]{
		minPages: common.Max(minPages, 1),
		pageSize: common.Max(pageSize, 1),
		length:   preAlloc,
	}

	for i := 0; i < array.minPages; i++ {
		array.pages = append(array.pages, make([]T, array.pageSize))
	}
	return array
}

func (this *PagedSlice[T]) NumPages() int { return len(this.pages) }
func (this *PagedSlice[T]) PageSize() int { return this.pageSize }
func (this *PagedSlice[T]) MinSize() int  { return this.minPages * this.pageSize }
func (this *PagedSlice[T]) Size() int     { return this.length }

// Compact reduces the number of pages in the PagedSlice by removing unused pages.
func (this *PagedSlice[T]) Compact() {
	usedPages := int(math.Ceil(float64(this.length) / float64(this.pageSize)))
	if usedPages < this.minPages {
		usedPages = this.minPages
	}
	this.pages = this.pages[:usedPages]
}

// Reset reduces the number of pages in the PagedSlice to the minimum number of pages set during initialization.
// func (this *PagedSlice[T]) Reset() int {
// 	numPages := len(this.pages) - this.minPages
// 	this.pages = this.pages[:this.minPages]

// 	if this.length > this.minPages*this.pageSize { // if the number of elements is greater than the minimum capacity
// 		this.length = this.minPages * this.pageSize
// 	}
// 	return numPages
// }

// Resize changes the size of the PagedSlice to the specified new size.
// If the new size is larger than the current capacity, the PagedSlice is resized to accommodate the new size.
func (this *PagedSlice[T]) Resize(newSize int) {
	if this.Cap() < newSize {
		this.Reserve(newSize)
	}
	this.length = newSize
}

// Append adds the specified values to the end of the PagedSlice.
func (this *PagedSlice[T]) Concate(values []T) {
	nextBlockID, offset := this.next()
	copy(this.pages[nextBlockID][offset:], values)

	minLen := common.Min(len(values), this.pageSize-offset)
	values = values[minLen:]
	this.length += minLen

	this.Reserve(len(values))
	for i := nextBlockID + 1; i < nextBlockID+int(math.Ceil(float64(len(values))/float64(this.pageSize)))+1; i++ {
		offset := (i - (nextBlockID + 1)) * this.pageSize
		copy(this.pages[i][:], values[offset:])
	}
	this.length += len(values)
}

// PushBack adds a value to the end of the PagedSlice.
func (this *PagedSlice[T]) PushBack(v T) {
	this.Reserve(1)
	i, j := this.Size()/this.pageSize, this.Size()%this.pageSize
	this.pages[i][j] = v
	this.length++
}

// PopBack removes and returns the value at the end of the PagedSlice.
func (this *PagedSlice[T]) PopBack() T {
	v := this.Back()
	this.length--
	return v
}

// Back returns the value at the end of the PagedSlice without removing it.
// No bounds checking is performed.
func (this *PagedSlice[T]) Back() T {
	return this.Get(this.length - 1)
}

// ToArry copies the elements from the PagedSlice to a new slice starting from the specified start index (inclusive)
// and ending at the specified end index (exclusive).
func (this *PagedSlice[T]) ToSlice(start int, end int) []T {
	buffer := make([]T, this.length)
	this.ToBuffer(start, end, buffer)
	return buffer
}

// PopBackToBuffer removes and copies the elements from the end of the PagedSlice to a specified buffer.
func (this *PagedSlice[T]) PopBackToBuffer(buffer []T) {
	start := common.Max(this.length-len(buffer), 0)
	this.ToBuffer(start, common.Min(start+len(buffer), this.Size()), buffer)
	this.length -= common.Min(len(buffer), this.Size())
}

// ToBuffer copies the elements from the PagedSlice to the specified buffer starting from the specified start
// index (inclusive) and ending at the specified end index (exclusive).
func (this *PagedSlice[T]) ToBuffer(start int, end int, buffer []T) {
	for i := start; i < end; i++ {
		buffer[i-start] = this.Get(i)
	}
}

// Cap returns the total capacity of the PagedSlice.
func (this *PagedSlice[T]) Cap() int {
	return this.pageSize * len(this.pages)
}

// Clear removes all elements from the PagedSlice. It also reduces the number of pages to the minimum number of pages.
func (this *PagedSlice[T]) Clear() {
	this.length = 0
	this.Compact()
}

// Set updates the value at the specified position in the PagedSlice.
// Returns true if the position is valid and the value is updated, false otherwise.
func (this *PagedSlice[T]) Set(pos int, v T) bool {
	if pos >= this.length {
		panic("Index out of bounds")
	}
	this.pages[pos/this.pageSize][pos%this.pageSize] = v
	return true
}

// Get returns the value at the specified position in the PagedSlice.
// No bounds checking is performed.
func (this *PagedSlice[T]) Get(pos int) T {
	return (this.pages[pos/this.pageSize][pos%this.pageSize])
}

// at returns a pointer to the value at the specified position in the PagedSlice.
// No bounds checking is performed.
func (this *PagedSlice[T]) at(pos int) *T {
	return &this.pages[pos/this.pageSize][pos%this.pageSize]
}

// Foreach applies the specified functor to each element in the PagedSlice[T] starting from the specified
// start index (inclusive) and ending at the specified end index (exclusive).
func (this *PagedSlice[T]) ForeachBetween(start int, end int, functor func(int, *T)) *PagedSlice[T] {
	for i := start; i < end; i++ {
		functor(i, this.at(i))
	}
	return this
}

// Foreach applies the specified functor to each element in the PagedSlice[T] starting from the specified
// start index (inclusive) and ending at the specified end index (exclusive).
func (this *PagedSlice[T]) Foreach(functor func(int, *T)) *PagedSlice[T] {
	for i := 0; i < this.length; i++ {
		functor(i, this.at(i))
	}
	return this
}

// ParallelForeach applies the specified functor to each element.
func (this *PagedSlice[T]) ParallelForeach(functor func(*T)) *PagedSlice[T] {
	common.ParallelFor(0, this.length, 6, func(i int) {
		functor(this.at(i))
	})
	return this
}

// Reserve increases the capacity of the PagedSlice to accommodate the specified size.
// Returns the number of additional pages reserved.
func (this *PagedSlice[T]) Reserve(size int) int {
	if this.Cap()-this.length >= size {
		return 0
	}

	offset := this.Cap() - this.Size()
	numBlocks := int(math.Ceil(float64(size-offset) / float64(this.pageSize)))
	for i := 0; i < numBlocks; i++ {
		newBlock := make([]T, this.pageSize)
		this.pages = append(this.pages, newBlock)
	}
	return numBlocks
}

// next returns the block ID and offset for the next element to be added to the PagedSlice.
func (this *PagedSlice[T]) next() (int, int) {
	if this.length%this.pageSize == 0 {
		return this.length / this.pageSize, 0
	} else {
		return this.length / this.pageSize, this.length % this.pageSize
	}
}
