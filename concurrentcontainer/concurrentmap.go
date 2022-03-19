package concurrentMap

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
	numShards uint8
	sharded   []map[string]interface{}
	lock      sync.RWMutex
}

func NewConcurrentMap() *ConcurrentMap {
	ccmap := &ConcurrentMap{
		numShards: 6,
	}

	ccmap.sharded = make([]map[string]interface{}, ccmap.numShards)
	for i := 0; i < int(ccmap.numShards); i++ {
		ccmap.sharded[i] = make(map[string]interface{}, 64)
	}
	return ccmap
}

func (this *ConcurrentMap) Size() uint32 {
	total := 0
	for i := 0; i < int(this.numShards); i++ {
		total += len(this.sharded[i])
	}
	return uint32(total)
}

func (this *ConcurrentMap) Get(key string, args ...interface{}) (interface{}, bool) {
	v, ok := this.sharded[this.Hash8(key)%this.numShards][key]
	return v, ok
}

func (this *ConcurrentMap) BatchGet(keys []string, args ...interface{}) []interface{} {
	shardIds := this.Hash8s(keys)
	values := make([]interface{}, len(keys))
	var wg sync.WaitGroup
	for threadID := 0; threadID < int(this.numShards); threadID++ {
		wg.Add(1)
		go func(threadID int) {
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

func (this *ConcurrentMap) delete(key string) {
	shardID := this.Hash8(key) % this.numShards
	delete(this.sharded[shardID], key)
}

func (this *ConcurrentMap) Set(key string, v interface{}, args ...interface{}) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if v == nil {
		this.delete(key)
	} else {
		shardID := this.Hash8(key) % this.numShards
		this.sharded[shardID][key] = v
	}
	return nil
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
	this.lock.Lock()
	defer this.lock.Unlock()

	var flags []bool
	if len(args) > 0 && args[0] != nil {
		flags = args[0].([]bool)
	}

	var wg sync.WaitGroup
	for threadID := 0; threadID < int(this.numShards); threadID++ {
		wg.Add(1)
		go func(threadID int) {
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
						this.delete(keys[i])
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
	total := 0
	for i := 0; i < int(this.numShards); i++ {
		total += len(this.sharded[i])
	}

	keys := make([]string, total)
	counter := 0
	for i := 0; i < int(this.numShards); i++ {
		for k := range this.sharded[i] {
			keys[counter] = k
			counter++
		}
	}
	return keys
}

func (this *ConcurrentMap) Hash8(key string) uint8 {
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
	shards := this.Shards()
	values := make([]interface{}, len(*shards))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			for _, v := range (*shards)[i] {
				if values[i] == nil || Compare(v, values[i]) {
					values[i] = v
				}
			}
		}
	}
	common.ParallelWorker(len(*shards), len(*shards), worker)

	val := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] != nil && (val == nil || Compare(values[i], val)) {
			val = values[i]
		}
	}
	return val
}

func (this *ConcurrentMap) Foreach(predicate func(interface{}) interface{}) {
	shards := this.Shards()
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			for k, v := range (*shards)[i] {
				if v = predicate(v); v == nil {
					delete((*shards)[i], k)
				} else {
					(*shards)[i][k] = v
				}
			}
		}
	}
	common.ParallelWorker(len(*shards), len(*shards), worker)
}

func (this *ConcurrentMap) KVs() ([]string, []interface{}) {
	keys := this.Keys()
	return keys, this.BatchGet(keys)
}

/* -----------------------------------Debug Functions---------------------------------------------------------*/
type Encodeable interface {
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
		vBytes = append(vBytes, v.(Encodeable).Encode()...)
	}

	kSum := sha256.Sum256(codec.Strings(k).Flatten())
	vSum := sha256.Sum256(vBytes)

	return sha256.Sum256(append(kSum[:], vSum[:]...))
}
