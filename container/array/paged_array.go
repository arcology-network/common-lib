// The PagedArray class is a custom data structure that represents an array that is divided into multiple blocks or pages.
// It is designed to efficiently handle large arrays by storing the elements in a paginated manner.

// The PagedArray class provides methods and operations to manipulate and access the elements stored in the array.
// It allows you to perform operations such as adding elements, retrieving elements by index, updating elements, and more.
// By dividing the array into blocks, it can optimize memory usage and improve performance.

package pagedarray

import (
	"math"

	"github.com/arcology-network/common-lib/common"
)

// PagedArray represents an array that is divided into multiple blocks or pages.
type PagedArray struct {
	minBlocks int
	blockSize int
	length    int
	blocks    [][]interface{}
}

// NewPagedArray creates a new instance of PagedArray with the specified block size and minimum number of blocks.
func NewPagedArray(blockSize int, minBlocks int) *PagedArray {
	ccArray := &PagedArray{
		minBlocks: common.Max(minBlocks, 1),
		blockSize: common.Max(blockSize, 1),
		length:    0,
	}

	for i := 0; i < ccArray.minBlocks; i++ {
		ccArray.blocks = append(ccArray.blocks, make([]interface{}, ccArray.blockSize))
	}
	return ccArray
}

// Shrink reduces the number of blocks in the PagedArray to match the number of blocks required
// to store the current number of elements.
func (this *PagedArray) Shrink() {
	usedBlocks := int(math.Ceil(float64(this.length) / float64(this.blockSize)))
	if usedBlocks < this.minBlocks {
		usedBlocks = this.minBlocks
	}
	this.blocks = this.blocks[:usedBlocks]
}

// Resize changes the size of the PagedArray to the specified new size.
// If the new size is larger than the current capacity, the PagedArray is resized to accommodate the new size.
func (this *PagedArray) Resize(newSize int) {
	if this.Cap() < newSize {
		this.Reserve(newSize)
	}
	this.length = newSize
}

// Append adds the specified values to the end of the PagedArray.
func (this *PagedArray) Append(values []interface{}) {
	nextBlockID, offset := this.next()
	copy(this.blocks[nextBlockID][offset:], values)

	minLen := common.Min(len(values), this.blockSize-offset)
	values = values[minLen:]
	this.length += minLen

	this.Reserve(len(values))
	for i := nextBlockID + 1; i < nextBlockID+int(math.Ceil(float64(len(values))/float64(this.blockSize)))+1; i++ {
		offset := (i - (nextBlockID + 1)) * this.blockSize
		copy(this.blocks[i][:], values[offset:])
	}
	this.length += len(values)
}

// PushBack adds a value to the end of the PagedArray.
func (this *PagedArray) PushBack(v interface{}) {
	this.Reserve(1)
	i, j := this.Size()/this.blockSize, this.Size()%this.blockSize
	this.blocks[i][j] = v
	this.length++
}

// PopBack removes and returns the value at the end of the PagedArray.
func (this *PagedArray) PopBack() interface{} {
	v := this.Back()
	if v != nil {
		this.length--
	}
	return v
}

// Back returns the value at the end of the PagedArray without removing it.
func (this *PagedArray) Back() interface{} {
	if this.length > 0 {
		v := this.Get(this.length - 1)
		return v
	}
	return nil
}

// CopyTo copies the elements from the PagedArray to a new slice starting from the specified start index (inclusive)
// and ending at the specified end index (exclusive).
func (this *PagedArray) CopyTo(start int, end int) []interface{} {
	buffer := make([]interface{}, this.length)
	this.ToBuffer(start, end, buffer)
	return buffer
}

// PopBackToBuffer removes and copies the elements from the end of the PagedArray to the specified buffer.
func (this *PagedArray) PopBackToBuffer(buffer []interface{}) {
	start := common.Max(this.length-len(buffer), 0)
	this.ToBuffer(start, common.Min(start+len(buffer), this.Size()), buffer)
	this.length -= common.Min(len(buffer), this.Size())
}

// ToBuffer copies the elements from the PagedArray to the specified buffer starting from the specified start
// index (inclusive) and ending at the specified end index (exclusive).
func (this *PagedArray) ToBuffer(start int, end int, buffer []interface{}) {
	for i := start; i < end; i++ {
		buffer[i-start] = this.Get(i)
	}
}

// Size returns the number of elements in the PagedArray.
func (this *PagedArray) Size() int {
	return this.length
}

// Cap returns the total capacity of the PagedArray.
func (this *PagedArray) Cap() int {
	return this.blockSize * len(this.blocks)
}

// Clear removes all elements from the PagedArray.
func (this *PagedArray) Clear() {
	this.length = 0
	this.Shrink()
}

// Set updates the value at the specified position in the PagedArray.
// Returns true if the position is valid and the value is updated, false otherwise.
func (this *PagedArray) Set(pos int, v interface{}) bool {
	if pos >= this.length {
		return false
	}

	this.blocks[pos/this.blockSize][pos%this.blockSize] = v
	return true
}

// Get returns the value at the specified position in the PagedArray.
// Returns nil if the position is invalid.
func (this *PagedArray) Get(pos int) interface{} {
	if pos >= this.length {
		return nil
	}

	return (this.blocks[pos/this.blockSize][pos%this.blockSize])
}

// at returns a pointer to the value at the specified position in the PagedArray.
// Returns nil if the position is invalid.
func (this *PagedArray) at(pos int) interface{} {
	if pos >= this.length {
		return nil
	}
	return &(this.blocks[pos/this.blockSize][pos%this.blockSize])
}

// Foreach applies the specified functor to each element in the PagedArray starting from the specified
// start index (inclusive) and ending at the specified end index (exclusive).
func (this *PagedArray) Foreach(start int, end int, functor func(interface{})) {
	for i := start; i < end; i++ {
		functor(this.at(i))
	}
}

// Reserve increases the capacity of the PagedArray to accommodate the specified size.
// Returns the number of additional blocks reserved.
func (this *PagedArray) Reserve(size int) int {
	if this.Cap()-this.length >= size {
		return 0
	}

	offset := this.Cap() - this.Size()
	numBlocks := int(math.Ceil(float64(size-offset) / float64(this.blockSize)))
	for i := 0; i < numBlocks; i++ {
		newBlock := make([]interface{}, this.blockSize)
		this.blocks = append(this.blocks, newBlock)
	}
	return numBlocks
}

// next returns the block ID and offset for the next element to be added to the PagedArray.
func (this *PagedArray) next() (int, int) {
	if this.length%this.blockSize == 0 {
		return this.length / this.blockSize, 0
	} else {
		return this.length / this.blockSize, this.length % this.blockSize
	}
}
