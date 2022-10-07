package cachedstorage

import (
	"fmt"
	"math"
	"math/rand"

	cccontainer "github.com/HPISTechnologies/common-lib/concurrentcontainer"
)

type CachePolicy struct {
	totalAllocated    uint64
	quota             uint64
	threshold         float64
	mincounts         uint32
	scoreboard        *cccontainer.ConcurrentMap
	scoreDistribution *Distribution

	keyBuffer   []string
	sizeBuffer  []uint32
	scoreBuffer []interface{} // access counts
}

func NewCachePolicy(hardQuota uint64, initThreshold float64) *CachePolicy {
	if initThreshold > 0.9 {
		initThreshold = 0.9
	}

	if initThreshold < 0.7 {
		initThreshold = 0.7
	}

	return &CachePolicy{
		totalAllocated:    0,
		quota:             hardQuota,
		threshold:         initThreshold,
		mincounts:         math.MaxUint32,
		scoreboard:        cccontainer.NewConcurrentMap(),
		scoreDistribution: NewDistribution(),

		keyBuffer:   make([]string, 0, 65536),
		sizeBuffer:  make([]uint32, 0, 65536),
		scoreBuffer: make([]interface{}, 0, 65536),
	}
}

func (this *CachePolicy) AdjustThreshold(hardQuota uint64, newthreshold float64) {
	if newthreshold > 0.9 {
		newthreshold = 0.9
	}

	if newthreshold < 0.7 {
		newthreshold = 0.7
	}

	this.threshold = newthreshold
	this.quota = hardQuota
}

func (this *CachePolicy) Size() uint32 {
	return this.scoreboard.Size()
}

func (this *CachePolicy) AddToBuffer(keys []string, vals []interface{}) {
	prevSize := len(this.keyBuffer)
	this.keyBuffer = append(this.keyBuffer, keys...)
	this.sizeBuffer = append(this.sizeBuffer, make([]uint32, len(vals))...)
	this.scoreBuffer = append(this.scoreBuffer, make([]interface{}, len(vals))...)
	for i := prevSize; i < len(this.keyBuffer); i++ {
		this.sizeBuffer[i] = vals[i-prevSize].(MeasurableInterface).Size()
		this.scoreBuffer[i] = vals[i-prevSize].(AccessableInterface).Reads() + vals[i-prevSize].(AccessableInterface).Writes()
	}
}

func (this *CachePolicy) FreeEntries(threshold uint32, probability float64, localCache *cccontainer.ConcurrentMap) (uint64, uint64) {
	scoreShards := this.scoreboard.Shards()
	cacheShards := localCache.Shards()
	freedMem := make([]uint64, len(*scoreShards))
	freedEntries := make([]uint64, len(*scoreShards))

	for i := 0; i < len(*scoreShards); i++ {
		for k := range (*scoreShards)[i] {
			score := (*scoreShards)[i][k].(uint32)
			if score >= threshold {
				continue
			}

			if (math.Abs(probability-1) < 0.05) || (rand.Float64() < probability) {
				v := (*cacheShards)[i][k]
				freedMem[i] += uint64(v.(MeasurableInterface).Size())
				delete((*scoreShards)[i], k)
				delete((*cacheShards)[i], k)
				freedEntries[i]++
			}
		}
	}

	entries := uint64(0)
	memory := uint64(0)
	for i := 0; i < len(freedMem); i++ {
		entries += freedEntries[i]
		memory += freedMem[i]
	}
	return entries, memory
}

func (this *CachePolicy) Refresh(cache *cccontainer.ConcurrentMap) {
	currentScores := this.scoreboard.BatchGet(this.keyBuffer)
	for i := 0; i < len(this.sizeBuffer); i++ {
		if currentScores[i] == nil {
			currentScores[i] = uint32(0)
		}

		this.sizeBuffer[i] = currentScores[i].(uint32) + this.sizeBuffer[i]
	}

	for i := 0; i < len(this.sizeBuffer); i++ {
		this.totalAllocated += uint64(this.sizeBuffer[i])
	}

	this.scoreboard.BatchSet(this.keyBuffer, this.scoreBuffer)                                                           // Update scores
	this.scoreDistribution.UpdateDistribution(this.keyBuffer, this.sizeBuffer, this.scoreBuffer, cache, this.scoreboard) // Update the entry distribution info

	this.keyBuffer = this.keyBuffer[:0]
	this.sizeBuffer = this.sizeBuffer[:0]
	this.scoreBuffer = this.scoreBuffer[:0]
}

func (this *CachePolicy) FreeMemory(localChache *cccontainer.ConcurrentMap) (uint64, uint64) {
	if this.totalAllocated > uint64(math.Round(float64(this.quota)*this.threshold)) {
		target := this.totalAllocated - uint64(math.Round(float64(this.quota)*this.threshold))
		threshold, prob := this.scoreDistribution.EstimateThreshold(target)

		entires, mem := this.FreeEntries(threshold, prob, localChache)
		this.totalAllocated -= mem
		return entires, mem
	}
	return 0, 0
}

func (this *CachePolicy) CheckCapacity(key string, v interface{}) bool {
	return uint64(v.(MeasurableInterface).Size())+this.totalAllocated < this.quota
}

func (this *CachePolicy) BatchCheckCapacity(keys []string, values []interface{}) bool {
	total := uint64(0)
	for i, v := range values {
		if len(keys[i]) != 0 { // Not in cache yet
			if total += uint64(v.(MeasurableInterface).Size()); total+this.totalAllocated >= this.quota {
				keys = keys[:i]
				values = values[:i]
				return false // No enough space for all
			}
		}
	}
	return true // Good for all
}

func (this *CachePolicy) PrintScores() {
	this.scoreboard.Print()
	fmt.Println()
}
