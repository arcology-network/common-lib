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
package livecache

import (
	"math"
	"sort"

	"github.com/arcology-network/common-lib/common"
	crdtcommon "github.com/arcology-network/common-lib/crdt/common"
	statecell "github.com/arcology-network/common-lib/crdt/statecell"
	"github.com/arcology-network/common-lib/exp/slice"
)



type CachePolicy[T any] struct {
	occupied uint64 // The total memory used by the cache.
	maxSize  uint64
	sizeOf   func(T) uint64
}

func NewCachePolicy[T any](maxSize uint64, sizeOf func(T) uint64) *CachePolicy[T] {
	usage := &CachePolicy[T]{
		occupied: 0,
		maxSize:  uint64(24 * 1024 * 1024 * 1024), // 0.8 of the minimum memory required.
		sizeOf:   sizeOf,
	}

	if v, err := common.GetAvailableMemory(); err == nil {
		usage.maxSize = common.Min(maxSize, uint64(float64(v)*0.8))
	}
	return usage
}

func (this *CachePolicy[T]) Size() uint64 {
	if this == nil {
		return 0
	}
	return this.occupied
}

func (this *CachePolicy[T]) ValueSize(value T) uint64 {
	if this == nil || this.sizeOf == nil {
		return 0
	}
	return this.sizeOf(value)
}

func (this *CachePolicy[T]) Admit(size uint64) bool {
	return this.Update(0, size)
}

func (this *CachePolicy[T]) Update(oldSize, newSize uint64) bool {
	if this == nil {
		return true
	}
	if oldSize > this.occupied {
		oldSize = this.occupied
	}
	if newSize <= oldSize {
		this.occupied -= oldSize - newSize
		return true
	}

	growth := newSize - oldSize
	if this.maxSize > 0 && this.occupied+growth > this.maxSize {
		return false
	}
	this.occupied += growth
	return true
}

func (this *CachePolicy[T]) Track(oldSize, newSize uint64) {
	if this == nil {
		return
	}
	if oldSize > this.occupied {
		oldSize = this.occupied
	}
	if newSize <= oldSize {
		this.occupied -= oldSize - newSize
		return
	}
	this.occupied += newSize - oldSize
}

func (this *CachePolicy[T]) NeedEviction() bool {
	if this == nil || this.maxSize == 0 {
		return false
	}
	return this.occupied > this.maxSize
}

func (this *CachePolicy[T]) Remove(size uint64) {
	if this == nil {
		return
	}
	if size >= this.occupied {
		this.occupied = 0
		return
	}
	this.occupied -= size
}

// Check if the cache has enough space to store the new values.
// If not, the cache will be cleared. If still not enough space,
// some new values won't be stored.
func (this *CachePolicy[T]) PrepareSpace(stCells *[]*statecell.StateCell, freeCache func(uint64) uint64) {
	if this == nil {
		return
	}

	// The total memory required to store the new values.
	totalRequired := slice.Accumulate(*stCells, uint64(0), func(_ int, v *statecell.StateCell) uint64 {
		if v.Value() == nil {
			return 0
		}
		return v.Value().(crdtcommon.CRDT).MemSize()
	})

	availableMemory, err := common.GetAvailableMemory()
	if err != nil {
		return
	}
	actualCap := common.Min(this.maxSize, availableMemory)

	toFree := int(totalRequired) - (int(actualCap) - int(this.occupied))
	if toFree <= 0 {
		this.occupied += totalRequired
		return
	}

	freedMemory := uint64(0)
	if freeCache != nil {
		freedMemory = freeCache(uint64(toFree))
	}

	remainingOccupied := this.occupied
	if freedMemory >= remainingOccupied {
		remainingOccupied = 0
	} else {
		remainingOccupied -= freedMemory
	}

	totalAvailable := int(actualCap) - int(remainingOccupied)
	if int(totalRequired) > totalAvailable {
		sort.Slice(*stCells, func(i, j int) bool {
			return (*stCells)[i].Value().(crdtcommon.CRDT).MemSize() < (*stCells)[j].Value().(crdtcommon.CRDT).MemSize()
		})

		idx := len(*stCells)
		accumSize := uint64(0)
		for i, v := range *stCells {
			size := v.Value().(crdtcommon.CRDT).MemSize()
			if accumSize += size; int(accumSize) > totalAvailable {
				idx = i
				accumSize -= size
				break
			}
		}

		this.occupied = remainingOccupied + accumSize
		*stCells = (*stCells)[:idx]
		return
	}

	this.occupied = remainingOccupied + totalRequired
}

func (this *CachePolicy[T]) AdjustFreeTarget(sizeToFree uint64, shardCount int) (uint64, []float64) {
	if this == nil || shardCount == 0 {
		return 0, nil
	}

	sizeToFree = common.Max(
		uint64(math.Ceil(float64(sizeToFree)*1.15)),
		uint64(shardCount*64),
	)
	return sizeToFree, slice.New(shardCount, math.Ceil(float64(sizeToFree)/float64(shardCount)))
}

func (this *CachePolicy[T]) EvictionScore(visits uint64, firstLoaded uint32, sizeToFree uint64) float32 {
	if sizeToFree <= uint64(firstLoaded) {
		return float32(visits)
	}
	return float32(visits) / float32(sizeToFree-uint64(firstLoaded))
}
