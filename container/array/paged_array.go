package pagedarray

import (
	"math"

	"github.com/arcology-network/common-lib/common"
)

type PagedArray struct {
	minBlocks int
	blockSize int
	length    int
	blocks    [][]interface{}
}

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

func (this *PagedArray) Shrink() {
	usedBlocks := int(math.Ceil(float64(this.length) / float64(this.blockSize)))
	if usedBlocks < this.minBlocks {
		usedBlocks = this.minBlocks
	}
	this.blocks = this.blocks[:usedBlocks]
}

func (this *PagedArray) Resize(newSize int) {
	if this.Cap() < newSize {
		this.Reserve(newSize)
	}
	this.length = newSize
}

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

func (this *PagedArray) PushBack(v interface{}) {
	this.Reserve(1)
	i, j := this.Size()/this.blockSize, this.Size()%this.blockSize
	this.blocks[i][j] = v
	this.length++
}

func (this *PagedArray) PopBack() interface{} {
	v := this.Back()
	if v != nil {
		this.length--
	}
	return v
}

func (this *PagedArray) Back() interface{} {
	if this.length > 0 {
		v := this.Get(this.length - 1)
		return v
	}
	return nil
}

func (this *PagedArray) CopyTo(start int, end int) []interface{} {
	buffer := make([]interface{}, this.length)
	this.ToBuffer(start, end, buffer)
	return buffer
}

func (this *PagedArray) PopBackToBuffer(buffer []interface{}) {
	start := common.Max(this.length-len(buffer), 0)
	this.ToBuffer(start, common.Min(start+len(buffer), this.Size()), buffer)
	this.length -= common.Min(len(buffer), this.Size())
}

func (this *PagedArray) ToBuffer(start int, end int, buffer []interface{}) {
	for i := start; i < end; i++ {
		buffer[i-start] = this.Get(i)
	}
}

func (this *PagedArray) Size() int {
	return this.length
}

func (this *PagedArray) Cap() int {
	return this.blockSize * len(this.blocks)
}

func (this *PagedArray) Clear() {
	this.length = 0
	this.Shrink()
}

func (this *PagedArray) Set(pos int, v interface{}) bool {
	if pos >= this.length {
		return false
	}

	this.blocks[pos/this.blockSize][pos%this.blockSize] = v
	return true
}

func (this *PagedArray) Get(pos int) interface{} {
	if pos >= this.length {
		return nil
	}

	return (this.blocks[pos/this.blockSize][pos%this.blockSize])
}

func (this *PagedArray) at(pos int) interface{} {
	if pos >= this.length {
		return nil
	}
	return &(this.blocks[pos/this.blockSize][pos%this.blockSize])
}

func (this *PagedArray) Foreach(start int, end int, functor func(interface{})) {
	for i := start; i < end; i++ {
		functor(this.at(i))
	}
}

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

func (this *PagedArray) next() (int, int) {
	if this.length%this.blockSize == 0 {
		return this.length / this.blockSize, 0
	} else {
		return this.length / this.blockSize, this.length % this.blockSize
	}
}
