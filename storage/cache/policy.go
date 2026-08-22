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
package cache

import (
	"github.com/arcology-network/common-lib/common"
)

type CachePolicy[T any] struct {
	occupied     uint64
	maxSize      uint64
	sizeOf       func(T) uint64
	triggerRatio float64
}

// NewCachePolicy returns a CachePolicy with a memory cap and value sizer.
func NewCachePolicy[T any](maxSize uint64, sizeOf func(T) uint64) *CachePolicy[T] {
	usage := &CachePolicy[T]{
		occupied:     0,
		maxSize:      uint64(24 * 1024 * 1024 * 1024),
		sizeOf:       sizeOf,
		triggerRatio: 1.2,
	}

	if v, err := common.GetAvailableMemory(); err == nil {
		usage.maxSize = min(maxSize, uint64(float64(v)*0.8))
	}
	return usage
}

// Size returns total memory used by the cache.
func (this *CachePolicy[T]) Size() uint64 { return this.occupied }

// ValueSize computes the memory size of a value.
func (this *CachePolicy[T]) ValueSize(value T) uint64 {
	if this == nil || this.sizeOf == nil {
		return 0
	}
	return this.sizeOf(value)
}

// Update adjusts memory usage by delta; returns false if a positive delta exceeds the limit.
func (this *CachePolicy[T]) Update(delta int64) bool {
	if this == nil {
		return true
	}

	if delta <= 0 {
		decrease := uint64(-(delta + 1)) + 1
		this.occupied -= min(decrease, this.occupied)
		return true
	}

	growth := uint64(delta)
	if this.maxSize > 0 && this.occupied+growth > uint64(float64(this.maxSize)*this.triggerRatio) {
		return false
	}

	this.occupied += growth
	return true
}

// NeedEviction returns true if cache exceeds its memory cap.
func (this *CachePolicy[T]) NeedEviction() bool {
	if this == nil || this.maxSize == 0 {
		return false
	}
	return this.occupied > this.maxSize
}
