// The PagedArray class is a custom data structure that represents an array that is divided into multiple pages or pages.
// It is designed to efficiently handle large arrays by storing the elements in a paginated manner to optimize memory usage
// and improve performance.
package array

import (
	"math"

	"github.com/arcology-network/common-lib/common"
)

type PagedArray[T any] struct {
	minPages int
	pageSize int
	length   int
	pages    [][]T
}

// NewPagedArray creates a new instance of PagedArray with the specified block size and minimum number of pages.
func NewPagedArray[T any](pageSize, minPages, preAlloc int) *PagedArray[T] {
	if pageSize*minPages < preAlloc {
		panic("preAlloc must be less than pageSize * minPages")
	}

	array := &PagedArray[T]{
		minPages: common.Max(minPages, 1),
		pageSize: common.Max(pageSize, 1),
		length:   preAlloc,
	}

	for i := 0; i < array.minPages; i++ {
		array.pages = append(array.pages, make([]T, array.pageSize))
	}
	return array
}

func (this *PagedArray[T]) NumPages() int { return len(this.pages) }
func (this *PagedArray[T]) PageSize() int { return this.pageSize }
func (this *PagedArray[T]) MinSize() int  { return this.minPages * this.pageSize }
func (this *PagedArray[T]) Size() int     { return this.length }

// Compact reduces the number of pages in the PagedArray by removing unused pages.
func (this *PagedArray[T]) Compact() {
	usedPages := int(math.Ceil(float64(this.length) / float64(this.pageSize)))
	if usedPages < this.minPages {
		usedPages = this.minPages
	}
	this.pages = this.pages[:usedPages]
}

// Reset reduces the number of pages in the PagedArray to the minimum number of pages set during initialization.
// func (this *PagedArray[T]) Reset() int {
// 	numPages := len(this.pages) - this.minPages
// 	this.pages = this.pages[:this.minPages]

// 	if this.length > this.minPages*this.pageSize { // if the number of elements is greater than the minimum capacity
// 		this.length = this.minPages * this.pageSize
// 	}
// 	return numPages
// }

// Resize changes the size of the PagedArray to the specified new size.
// If the new size is larger than the current capacity, the PagedArray is resized to accommodate the new size.
func (this *PagedArray[T]) Resize(newSize int) {
	if this.Cap() < newSize {
		this.Reserve(newSize)
	}
	this.length = newSize
}

// Append adds the specified values to the end of the PagedArray.
func (this *PagedArray[T]) Concate(values []T) {
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

// PushBack adds a value to the end of the PagedArray.
func (this *PagedArray[T]) PushBack(v T) {
	this.Reserve(1)
	i, j := this.Size()/this.pageSize, this.Size()%this.pageSize
	this.pages[i][j] = v
	this.length++
}

// PopBack removes and returns the value at the end of the PagedArray.
func (this *PagedArray[T]) PopBack() T {
	v := this.Back()
	// if v != nil {
	this.length--
	// }
	return v
}

// Back returns the value at the end of the PagedArray without removing it.
// No bounds checking is performed.
func (this *PagedArray[T]) Back() T {
	return this.Get(this.length - 1)
}

// ToArry copies the elements from the PagedArray to a new slice starting from the specified start index (inclusive)
// and ending at the specified end index (exclusive).
func (this *PagedArray[T]) ToSlice(start int, end int) []T {
	buffer := make([]T, this.length)
	this.ToBuffer(start, end, buffer)
	return buffer
}

// PopBackToBuffer removes and copies the elements from the end of the PagedArray to the specified buffer.
func (this *PagedArray[T]) PopBackToBuffer(buffer []T) {
	start := common.Max(this.length-len(buffer), 0)
	this.ToBuffer(start, common.Min(start+len(buffer), this.Size()), buffer)
	this.length -= common.Min(len(buffer), this.Size())
}

// ToBuffer copies the elements from the PagedArray to the specified buffer starting from the specified start
// index (inclusive) and ending at the specified end index (exclusive).
func (this *PagedArray[T]) ToBuffer(start int, end int, buffer []T) {
	for i := start; i < end; i++ {
		buffer[i-start] = this.Get(i)
	}
}

// Cap returns the total capacity of the PagedArray.
func (this *PagedArray[T]) Cap() int {
	return this.pageSize * len(this.pages)
}

// Clear removes all elements from the PagedArray. It also reduces the number of pages to the minimum number of pages.
func (this *PagedArray[T]) Clear() {
	this.length = 0
	this.Compact()
}

// Set updates the value at the specified position in the PagedArray.
// Returns true if the position is valid and the value is updated, false otherwise.
func (this *PagedArray[T]) Set(pos int, v T) bool {
	if pos >= this.length {
		panic("Index out of bounds")
	}
	this.pages[pos/this.pageSize][pos%this.pageSize] = v
	return true
}

// Get returns the value at the specified position in the PagedArray.
// No bounds checking is performed.
func (this *PagedArray[T]) Get(pos int) T {
	return (this.pages[pos/this.pageSize][pos%this.pageSize])
}

// at returns a pointer to the value at the specified position in the PagedArray.
// No bounds checking is performed.
func (this *PagedArray[T]) at(pos int) *T {
	return &this.pages[pos/this.pageSize][pos%this.pageSize]
}

// Foreach applies the specified functor to each element in the PagedArray[T] starting from the specified
// start index (inclusive) and ending at the specified end index (exclusive).
func (this *PagedArray[T]) ForeachBetween(start int, end int, functor func(int, *T)) *PagedArray[T] {
	for i := start; i < end; i++ {
		functor(i, this.at(i))
	}
	return this
}

// Foreach applies the specified functor to each element in the PagedArray[T] starting from the specified
// start index (inclusive) and ending at the specified end index (exclusive).
func (this *PagedArray[T]) Foreach(functor func(int, *T)) *PagedArray[T] {
	for i := 0; i < this.length; i++ {
		functor(i, this.at(i))
	}
	return this
}

// ParallelForeach applies the specified functor to each element.
func (this *PagedArray[T]) ParallelForeach(functor func(*T)) *PagedArray[T] {
	common.ParallelFor(0, this.length, 6, func(i int) {
		functor(this.at(i))
	})
	return this
}

// Reserve increases the capacity of the PagedArray to accommodate the specified size.
// Returns the number of additional pages reserved.
func (this *PagedArray[T]) Reserve(size int) int {
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

// next returns the block ID and offset for the next element to be added to the PagedArray.
func (this *PagedArray[T]) next() (int, int) {
	if this.length%this.pageSize == 0 {
		return this.length / this.pageSize, 0
	} else {
		return this.length / this.pageSize, this.length % this.pageSize
	}
}
