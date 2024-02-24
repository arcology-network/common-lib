package datastore

import (
	"math"

	common "github.com/arcology-network/common-lib/common"
	ccmap "github.com/arcology-network/common-lib/exp/map"
	intf "github.com/arcology-network/common-lib/storage/interface"
)

type Score struct {
	totalMemory  uint64
	avgMemory    uint64
	totalEntries uint64
}

func (this *Score) From(memSize uint32) {
	this.avgMemory = 0
	this.totalMemory = uint64(math.Max(float64(this.totalMemory)-float64(memSize), 0))
	if this.totalEntries > 0 {
		this.totalEntries--
	}
}

func (this *Score) To(newMemSize uint32) {
	this.totalMemory += uint64(newMemSize)
	this.totalEntries++
}

type Distribution struct {
	numCat    uint32
	scoreInfo []Score
}

func NewDistribution() *Distribution {
	return &Distribution{
		numCat:    32768,
		scoreInfo: make([]Score, 32768),
	}
}

func (this *Distribution) updateDistribution(keys []string, nSizes []uint32, newScores []interface{}, cache *ccmap.ConcurrentMap[string, any], scoreBoard *ccmap.ConcurrentMap[string, any]) {
	curtSizes := this.getCurrentSizes(keys, cache)
	curtScores := common.FilterFirst(scoreBoard.BatchGet(keys)) // Get the scores of the existing values.

	for i := 0; i < len(keys); i++ {
		curtCat := curtScores[i].(uint32) % this.numCat // curt category ID
		this.scoreInfo[curtCat].From(curtSizes[i])
	}

	for i := 0; i < len(keys); i++ {
		newCat := newScores[i].(uint32) % this.numCat
		this.scoreInfo[newCat].To(nSizes[i])
	}
}

func (this *Distribution) getCurrentSizes(keys []string, cache *ccmap.ConcurrentMap[string, any]) []uint32 {
	curtSizes := make([]uint32, len(keys))
	curtValues := common.FilterFirst(cache.BatchGet(keys))
	for i := 0; i < len(keys); i++ {
		if curtValues[i] == nil {
			curtSizes[i] = 0
		} else {
			curtSizes[i] = curtValues[i].(intf.Accessible).Size() // Get the sizes of the existing values.
		}
	}
	return curtSizes
}

func (this *Distribution) EstimateThreshold(target uint64) (uint32, float64) {
	mem := uint64(0)
	for i := 0; i < len(this.scoreInfo); i++ {
		mem += this.scoreInfo[i].totalMemory
		if mem > target {
			return uint32(i+1) * (this.numCat), float64(target) / float64(mem)
		}
	}
	return math.MaxUint32, 0.0
}
