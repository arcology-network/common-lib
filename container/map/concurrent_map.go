package concurrentmap

import (
	"crypto/sha256"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
)

type ConcurrentMap struct {
	numShards  uint8
	sharded    []map[string]interface{}
	shardLocks []sync.RWMutex
}

func NewConcurrentMap(args ...interface{}) *ConcurrentMap {
	defaultShards := uint8(6)
	if len(args) > 0 && args[0] != nil {
		defaultShards = uint8(args[0].(uint8))
		if defaultShards > 254 {
			defaultShards = 254
		}
	}

	ccmap := &ConcurrentMap{
		numShards: defaultShards,
	}

	ccmap.sharded = make([]map[string]interface{}, ccmap.numShards)
	ccmap.shardLocks = make([]sync.RWMutex, ccmap.numShards)
	for i := 0; i < int(ccmap.numShards); i++ {
		ccmap.sharded[i] = make(map[string]interface{}, 64)
	}
	return ccmap
}

func (this *ConcurrentMap) Size() uint32 {
	total := 0
	for i := 0; i < int(this.numShards); i++ {
		this.shardLocks[i].RLock()
		total += len(this.sharded[i])
		this.shardLocks[i].RUnlock()
	}
	return uint32(total)
}

func (this *ConcurrentMap) Get(key string, args ...interface{}) (interface{}, bool) {
	shardID := this.Hash8(key)
	if shardID >= uint8(len(this.sharded)) {
		return nil, true
	}

	this.shardLocks[shardID].RLock()
	defer this.shardLocks[shardID].RUnlock()

	v, ok := this.sharded[shardID][key]
	return v, ok
}

func (this *ConcurrentMap) BatchGet(keys []string, args ...interface{}) []interface{} {
	shardIds := this.Hash8s(keys)
	values := make([]interface{}, len(keys))
	var wg sync.WaitGroup
	for threadID := 0; threadID < int(this.numShards); threadID++ {
		wg.Add(1)
		go func(threadID int) {
			this.shardLocks[threadID].RLock()
			defer this.shardLocks[threadID].RUnlock()

			defer wg.Done()
			for i := 0; i < len(keys); i++ {
				if shardIds[i] == uint8(threadID) {
					values[i] = this.sharded[threadID][keys[i]]
				}
			}

		}(threadID)
	}
	wg.Wait()
	return values
}

func (this *ConcurrentMap) DirectBatchGet(shardIDs []uint8, keys []string, args ...interface{}) []interface{} {
	values := make([]interface{}, len(keys))
	setter := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			values[i] = this.sharded[shardIDs[i]][keys[i]]
		}
	}
	common.ParallelWorker(len(values), 4, setter)
	return values
}

func (this *ConcurrentMap) delete(shardID uint8, key string) {
	delete(this.sharded[shardID], key)
}

func (this *ConcurrentMap) Set(key string, v interface{}, args ...interface{}) error {
	shardID := this.Hash8(key)
	if shardID >= uint8(len(this.sharded)) {
		return nil
	}

	this.shardLocks[shardID].Lock()
	defer this.shardLocks[shardID].Unlock()

	if v == nil {
		this.delete(shardID, key)
	} else {
		this.sharded[shardID][key] = v
	}
	return nil
}

func (this *ConcurrentMap) BatchUpdate(keys []string, values []interface{}, Updater func(origin interface{}, index int, key string, value interface{}) interface{}) {
	shards := this.Hash8s(keys)
	var wg sync.WaitGroup
	for shard := uint8(0); shard < this.numShards; shard++ {
		wg.Add(1)
		go func(shard uint8) {
			this.shardLocks[shard].Lock()
			defer this.shardLocks[shard].Unlock()
			defer wg.Done()

			for i := 0; i < len(keys); i++ {
				if shards[i] != shard {
					continue
				}

				this.sharded[shard][keys[i]] = Updater(this.sharded[shard][keys[i]], i, keys[i], values[i])
			}
		}(shard)
	}
	wg.Wait()
}

func (this *ConcurrentMap) Traverse(Operator func(key string, value interface{}) (interface{}, interface{})) [][]interface{} {
	results := make([][]interface{}, this.numShards)
	var wg sync.WaitGroup
	for shard := uint8(0); shard < this.numShards; shard++ {
		wg.Add(1)
		go func(shard uint8) {
			this.shardLocks[shard].Lock()
			defer this.shardLocks[shard].Unlock()
			defer wg.Done()

			for k, v := range this.sharded[shard] {
				newV, ret := Operator(k, v)
				if ret != nil {
					results[shard] = append(results[shard], ret)
				}

				if newV != nil {
					this.sharded[shard][k] = newV
				} else {
					delete(this.sharded[shard], k)
				}
			}
		}(shard)
	}
	wg.Wait()
	return results
}

func (this *ConcurrentMap) BatchSet(keys []string, values []interface{}, args ...interface{}) {
	shardIDs := this.Hash8s(keys)

	if len(args) > 0 {
		this.DirectBatchSet(shardIDs, keys, values, args[0])
	} else {
		this.DirectBatchSet(shardIDs, keys, values)
	}
}

func (this *ConcurrentMap) DirectBatchSet(shardIDs []uint8, keys []string, values []interface{}, args ...interface{}) {
	var flags []bool
	if len(args) > 0 && args[0] != nil {
		flags = args[0].([]bool)
	}

	var wg sync.WaitGroup
	for threadID := 0; threadID < int(this.numShards); threadID++ {
		wg.Add(1)
		go func(threadID int) {
			this.shardLocks[threadID].Lock()
			defer this.shardLocks[threadID].Unlock()
			defer wg.Done()
			for i := 0; i < len(keys); i++ {
				if shardIDs[i] == uint8(threadID) {
					if len(flags) > 0 && !flags[i] {
						continue
					}

					if len(keys[i]) == 0 {
						continue
					}

					if values[i] == nil {
						this.delete(shardIDs[i], keys[i])
					} else {
						this.sharded[threadID][keys[i]] = values[i]
					}
				}
			}

		}(threadID)
	}
	wg.Wait()
}

func (this *ConcurrentMap) Keys() []string {
	total := uint32(0)
	offsets := make([]uint32, this.numShards+1)
	for i := 0; i < int(this.numShards); i++ {
		total += uint32(len(this.sharded[i]))
		offsets[i+1] = total
	}

	keys := make([]string, total)
	worker := func(start, end, index int, args ...interface{}) {
		this.shardLocks[start].Lock()
		defer this.shardLocks[start].Unlock()

		for i := start; i < end; i++ {
			counter := offsets[i]
			for k := range this.sharded[i] {
				keys[counter] = k
				counter++
			}
		}
	}
	common.ParallelWorker(len(this.sharded), len(this.sharded), worker)
	return keys
}

func (this *ConcurrentMap) Hash8(key string) uint8 {
	if len(key) == 0 {
		return math.MaxUint8
	}

	var total uint32 = 0
	for j := 0; j < len(key); j++ {
		total += uint32(key[j])
	}
	return uint8(total % uint32(this.numShards))
}

func (this *ConcurrentMap) Hash8s(keys []string) []uint8 {
	shardIds := make([]uint8, len(keys))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			if len(keys[i]) > 0 {
				shardIds[i] = this.Hash8(keys[i])
			} else {
				shardIds[i] = math.MaxUint8
			}
		}
	}
	common.ParallelWorker(len(keys), 8, worker)
	return shardIds
}

func (this *ConcurrentMap) Shards() *[]map[string]interface{} {
	return &this.sharded
}

func (this *ConcurrentMap) Find(Compare func(interface{}, interface{}) bool) interface{} {
	values := make([]interface{}, len(this.sharded))
	worker := func(start, end, index int, args ...interface{}) {
		this.shardLocks[start].RLock()
		defer this.shardLocks[start].RUnlock()

		for i := start; i < end; i++ {
			for _, v := range this.sharded[i] {
				if values[i] == nil || Compare(v, values[i]) {
					values[i] = v
				}
			}
		}
	}
	common.ParallelWorker(len(this.sharded), len(this.sharded), worker)

	val := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] != nil && (val == nil || Compare(values[i], val)) {
			val = values[i]
		}
	}
	return val
}

func (this *ConcurrentMap) Foreach(predicate func(interface{}) interface{}) {
	worker := func(start, end, index int, args ...interface{}) {
		this.shardLocks[start].RLock()
		defer this.shardLocks[start].RUnlock()

		for i := start; i < end; i++ {
			for k, v := range this.sharded[i] {
				if v = predicate(v); v == nil {
					delete(this.sharded[i], k)
				} else {
					this.sharded[i][k] = v
				}
			}
		}
	}
	common.ParallelWorker(len(this.sharded), len(this.sharded), worker)
}

func (this *ConcurrentMap) KVs() ([]string, []interface{}) {
	keys := this.Keys()
	return keys, this.BatchGet(keys)
}

func (this *ConcurrentMap) Clear() {
	cleaner := func(start, end, index int, args ...interface{}) {
		this.shardLocks[start].RLock()
		defer this.shardLocks[start].RUnlock()
		this.sharded[start] = make(map[string]interface{}, 64)
	}
	common.ParallelWorker(len(this.sharded), len(this.sharded), cleaner)
}

/* -----------------------------------Debug Functions---------------------------------------------------------*/
type Encodable interface {
	Encode() []byte
}

func (this *ConcurrentMap) Dump() ([]string, []interface{}) {
	keys := this.Keys()
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	vStr := make([]interface{}, len(keys))
	for i, v := range this.BatchGet(keys) {
		vStr[i] = v
	}
	return keys, vStr
}

func (this *ConcurrentMap) Print() {
	keys := []string{}
	partitionID := []int{}
	values := make([]interface{}, 0)
	for i, shard := range this.sharded {
		for k, v := range shard {
			// if v != nil {
			keys = append(keys, k)
			partitionID = append(partitionID, i)
			if v != nil {
				values = append(values, v)
			} else {
				values = append(values, "nil")
			}
			// }
		}
	}

	for i := range keys {
		fmt.Println("Partition: ", partitionID[i], "| Key: ", keys[i], "| Value: ", values[i])
	}
}

func (this *ConcurrentMap) Checksum() [32]byte {
	k, values := this.Dump()
	vBytes := []byte{}
	for _, v := range values {
		vBytes = append(vBytes, v.(Encodable).Encode()...)
	}

	kSum := sha256.Sum256(codec.Strings(k).Flatten())
	vSum := sha256.Sum256(vBytes)

	return sha256.Sum256(append(kSum[:], vSum[:]...))
}
