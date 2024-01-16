package datastore

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/arcology-network/common-lib/exp/array"
	ccmap "github.com/arcology-network/common-lib/exp/map"
	intf "github.com/arcology-network/common-lib/storage/interface"
	memdb "github.com/arcology-network/common-lib/storage/memdb"
)

const (
	Cache_Quota_Full = math.MaxUint64
)

type TypeAccessibleInterface interface {
	Value() interface{}
	MemSize() uint32
}

type CachePolicy struct {
	totalAllocated    uint64
	quota             uint64
	threshold         float64
	scoreboard        *ccmap.ConcurrentMap[string, any]
	scoreDistribution *Distribution

	keys   []string
	sizes  []uint32
	scores []interface{} // access counts
	lock   sync.RWMutex
}

// Memory hard Quota
func NewCachePolicy(hardQuota uint64, threshold float64) *CachePolicy {
	m := ccmap.NewConcurrentMap[string, any](8, func(v any) bool { return v == nil }, func(k string) uint8 {
		return array.Sum[byte, uint8]([]byte(k))
	})

	policy := &CachePolicy{
		totalAllocated:    0,
		quota:             hardQuota,
		threshold:         threshold,
		scoreboard:        m,
		scoreDistribution: NewDistribution(),

		keys:   make([]string, 0, 65536),
		sizes:  make([]uint32, 0, 65536),
		scores: make([]interface{}, 0, 65536),
	}
	policy.adjustThreshold(hardQuota, threshold)
	return policy
}

func (this *CachePolicy) Customize(db intf.PersistentStorage) *CachePolicy {
	if this != nil {
		if _, ok := db.(*memdb.MemoryDB); ok { // A memory DB doesn't need a in-memory cache
			this.quota = 0
		}
	}
	return this
}

func (this *CachePolicy) IsFull() bool {
	return this.quota == 0 || this.totalAllocated >= this.quota
}

func (this *CachePolicy) InfinitCache() bool {
	return this.quota == Cache_Quota_Full
}

func (this *CachePolicy) adjustThreshold(hardQuota uint64, threshold float64) {
	if this.isFixed() {
		return
	}

	this.threshold = this.guardRange(threshold)
	this.quota = hardQuota
}

func (this *CachePolicy) guardRange(newthreshold float64) float64 {
	if newthreshold > 0.9 {
		newthreshold = 0.9
	}

	if newthreshold < 0.7 {
		newthreshold = 0.7
	}
	return newthreshold
}

func (this *CachePolicy) isFixed() bool {
	return this.quota == 0 || this.quota == math.MaxUint64
}

func (this *CachePolicy) Size() uint32 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.scoreboard.Size()
}

func (this *CachePolicy) AddToStats(keys []string, vals []intf.Accessible) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.isFixed() {
		return
	}

	prevSize := len(this.keys)
	this.keys = append(this.keys, keys...)
	this.sizes = append(this.sizes, make([]uint32, len(vals))...)
	this.scores = append(this.scores, make([]interface{}, len(vals))...)

	if len(this.keys) != len(this.sizes) || len(this.keys) != len(this.scores) {
		fmt.Println("this.keys: ", len(this.keys))
		fmt.Println("this.sizes: ", len(this.sizes))
		fmt.Println("this.scores: ", len(this.scores))
		panic("Error: Sizes don't match !")
	}

	for i := 0; i < len(keys); i++ {
		this.sizes[i+prevSize] = vals[i].Size()
		this.scores[i+prevSize] = vals[i].Reads() + vals[i].Writes()
	}

	//there are some dupilicate keys
}

func (this *CachePolicy) freeEntries(threshold uint32, probability float64, localCache *ccmap.ConcurrentMap[string, any]) (uint64, uint64) {
	if this.isFixed() {
		return 0, 0
	}

	scoreShards := this.scoreboard.Shards()
	cacheShards := localCache.Shards()
	freedMem := make([]uint64, len(scoreShards))
	freedEntries := make([]uint64, len(scoreShards))

	for i := 0; i < len(scoreShards); i++ {
		for k := range (scoreShards)[i] {
			score := (scoreShards)[i][k].(uint32)
			if score >= threshold {
				continue
			}

			if (math.Abs(probability-1) < 0.05) || (rand.Float64() < probability) {
				v := (cacheShards)[i][k]
				freedMem[i] += uint64(v.(intf.Accessible).Size())
				delete((scoreShards)[i], k)
				delete((cacheShards)[i], k)
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

func (this *CachePolicy) Refresh(cache *ccmap.ConcurrentMap[string, any]) (uint64, uint64) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.isFixed() {
		return 0, 0
	}

	currentScores := this.scoreboard.BatchGet(this.keys)
	for i := 0; i < len(this.sizes); i++ {
		if currentScores[i] == nil {
			currentScores[i] = uint32(0)
		}

		this.sizes[i] = currentScores[i].(uint32) + this.sizes[i]
	}

	for i := 0; i < len(this.sizes); i++ {
		this.totalAllocated += uint64(this.sizes[i])
	}

	this.scoreboard.BatchSet(this.keys, this.scores)                                                      // Update scores
	this.scoreDistribution.updateDistribution(this.keys, this.sizes, this.scores, cache, this.scoreboard) // Update the entry distribution info

	this.keys = this.keys[:0]
	this.sizes = this.sizes[:0]
	this.scores = this.scores[:0]

	return this.freeMemory(cache)
}

func (this *CachePolicy) freeMemory(localChache *ccmap.ConcurrentMap[string, any]) (uint64, uint64) {
	if this.isFixed() {
		return 0, 0
	}

	if this.totalAllocated > uint64(math.Round(float64(this.quota)*this.threshold)) {
		target := this.totalAllocated - uint64(math.Round(float64(this.quota)*this.threshold))
		threshold, prob := this.scoreDistribution.EstimateThreshold(target)

		entires, mem := this.freeEntries(threshold, prob, localChache)
		this.totalAllocated -= mem
		return entires, mem
	}
	return 0, 0
}

func (this *CachePolicy) CheckCapacity(key string, v interface{}) bool {
	if this.quota == math.MaxUint64 || v == nil { // All in cache
		return true
	}

	if this.quota == 0 { // Non in cache
		return false
	}

	m := uint64(v.(TypeAccessibleInterface).MemSize())
	return m+this.totalAllocated < this.quota
}

func (this *CachePolicy) BatchCheckCapacity(keys []string, values []interface{}) ([]bool, uint32, bool) {
	if this.quota == math.MaxUint64 {
		return nil, 0, true
	}

	flags := make([]bool, len(keys))
	count := uint32(0)

	if this.quota == 0 {
		return flags, 0, false // Non in cache
	}

	if this.quota == math.MaxUint64 {
		for i := 0; i < len(flags); i++ {
			flags[i] = true
		}
		return flags, uint32(len(keys)), false // All in the cache
	}

	total := uint64(0)
	for i, v := range values {
		if len(keys[i]) != 0 && v != nil { // Not in the cache yet
			total += uint64(v.(TypeAccessibleInterface).MemSize())
			if flags[i] = total+this.totalAllocated < this.quota; flags[i] { // If no enough space for all
				count++
			} else {
				break
			}
		}

		if v == nil {
			flags[i] = true // Delete is always fine
			count++
		}
	}
	return flags, count, false // Good for all entries to stay in the memory
}

func (this *CachePolicy) PrintScores() {
	this.scoreboard.Print()
	fmt.Println()
}
