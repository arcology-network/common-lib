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
	"math"

	paged "github.com/arcology-network/common-lib/exp/pagedslice"
)

type Usage[K comparable] struct {
	key         *K
	score       float64
	sizeInMem   uint64
	lastLoaded  uint32
	firstLoaded uint32
}

func (this *Usage[K]) CalculateScore(maxScore float64) float64 {
	this.score = (float64(this.sizeInMem) / float64(this.lastLoaded-this.firstLoaded)) / maxScore
	return this.score
}

type CacheUsage[K comparable] struct {
	lookup   map[K]uint64 // key and its index in the paged slice
	keys     *paged.PagedSlice[*Usage[K]]
	maxScore float64
	dist     [65536]uint64 // size distribution
}

func NewCacheUsage[K comparable]() *CacheUsage[K] {
	return &CacheUsage[K]{
		maxScore: 0,
		keys:     paged.NewPagedSlice[*Usage[K]](1024, 100, 0),
		dist:     [65536]uint64{},
	}
}

func (this *CacheUsage[K]) Add(k *K, sizeInMem uint64, block uint32) float64 {
	idx, ok := this.lookup[*k]
	if !ok {
		this.keys.PushBack(
			&Usage[K]{
				key:         k,
				sizeInMem:   sizeInMem,
				lastLoaded:  block,
				firstLoaded: block,
			},
		)
		this.lookup[*k] = uint64(this.keys.Size() - 1)
	}

	v := this.keys.Get(int(idx))
	if *v.key != *k {
		panic("key mismatch")
	}

	v.lastLoaded = block
	v.sizeInMem = sizeInMem
	this.maxScore = math.Max(v.CalculateScore(this.maxScore), this.maxScore)
	this.dist[uint32(v.score*float64(len(this.dist)))] += sizeInMem
	return this.maxScore
}

func (this *CacheUsage[K]) FreeCache(sizeToFree uint64) (uint64, []K) {
	threshhold := this.findThreshold(sizeToFree)

	totalSize := uint64(0)
	keysToFree := make([]K, 0, 65536)
	for i := 0; i < this.keys.Size(); i++ {
		k := this.keys.Get(i)
		if k.score >= threshhold {
			totalSize += k.sizeInMem
			keysToFree = append(keysToFree, *k.key)
		}
	}
	return totalSize, keysToFree
}

func (this *CacheUsage[K]) findThreshold(sizeToFree uint64) float64 {
	for i := 0; i < len(this.dist); i++ {
		if sizeToFree -= this.dist[i]; sizeToFree <= 0 {
			return (float64(len(this.dist)) / float64(i)) * this.maxScore
		}
	}
	return this.maxScore
}
