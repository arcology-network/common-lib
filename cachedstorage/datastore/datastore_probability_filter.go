package cachedstorage

import (
	"sync"
	"time"

	ccmap "github.com/arcology-network/common-lib/container/map"
)

type ProbFilter struct {
	timeouts []string
	lookup   *ccmap.ConcurrentMap
	tLock    sync.Mutex
	locks    []sync.RWMutex
}

func NewProbFilter(shards uint8) *ProbFilter {
	return &ProbFilter{
		timeouts: []string{},
		lookup:   ccmap.NewConcurrentMap(shards),
		locks:    make([]sync.RWMutex, shards),
	}
}

func (this *ProbFilter) Checkin(key string) (bool, error) {
	idx := this.lookup.Hash8(key)

	this.locks[idx].Lock()
	if v, _ := this.lookup.Get(key); v == nil {
		this.lookup.Set(key, time.Now())
		this.locks[idx].Unlock()
		return true, nil
	}
	this.locks[idx].Unlock()

	counter := 0
	for { // Wait until the entry is avilable in the cache
		if v, _ := this.lookup.Get(key); v == nil {
			counter++
		} else {
			if counter > 10000 { // 10s
				this.tLock.Lock()
				this.timeouts = append(this.timeouts, key)
				this.tLock.Unlock()
			}
		}
		time.Sleep(time.Millisecond)
	}
	return false, nil
}

func (this *ProbFilter) Checkout(key string) error {
	idx := this.lookup.Hash8(key)
	this.locks[idx].Lock()
	defer this.locks[idx].Unlock()
	return this.lookup.Set(key, nil)
}

func (this *ProbFilter) Clear() {
	this.timeouts = this.timeouts[:0]
	this.lookup = ccmap.NewConcurrentMap(len(this.locks))
	this.locks = make([]sync.RWMutex, len(this.locks))
}
