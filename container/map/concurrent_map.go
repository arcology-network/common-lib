// The ConcurrentMap class is a concurrent map implementation allowing
// multiple goroutines to access and modify the map concurrently.

package concurrentmap

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"sync"

	"github.com/arcology-network/common-lib/codec"
	array "github.com/arcology-network/common-lib/exp/array"
)

// ConcurrentMap represents a concurrent map data structure.
type ConcurrentMap struct {
	numShards  uint8
	sharded    []map[string]interface{}
	shardLocks []sync.RWMutex
}

// NewConcurrentMap creates a new instance of ConcurrentMap with the specified number of shards.
// If no number of shards is provided, it defaults to 6.
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

// Size returns the total number of key-value pairs in the ConcurrentMap.
func (this *ConcurrentMap) Size() uint32 {
	total := 0
	for i := 0; i < int(this.numShards); i++ {
		this.shardLocks[i].RLock()
		total += len(this.sharded[i])
		this.shardLocks[i].RUnlock()
	}
	return uint32(total)
}

// Get retrieves the value associated with the specified key from the ConcurrentMap.
// It returns the value and a boolean indicating whether the key was found.
func (this *ConcurrentMap) Get(key string, args ...interface{}) (interface{}, bool) {
	shardID := this.Hash8(key)
	if shardID > uint8(len(this.sharded)) {
		return nil, false
	}

	this.shardLocks[shardID].RLock()
	defer this.shardLocks[shardID].RUnlock()
	k, v := this.sharded[shardID][key]
	return k, v
}

// BatchGet retrieves the values associated with the specified keys from the ConcurrentMap.
// It returns a slice of values in the same order as the keys.
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

// DirectBatchGet retrieves the values associated with the specified shard IDs and keys from the ConcurrentMap.
// It returns a slice of values in the same order as the keys.
func (this *ConcurrentMap) DirectBatchGet(shardIDs []uint8, keys []string, args ...interface{}) []interface{} {
	values := make([]interface{}, len(keys))
	array.ParallelForeach(values, 5, func(i int, _ *interface{}) {
		values[i] = this.sharded[shardIDs[i]][keys[i]]
	})
	return values
}

func (this *ConcurrentMap) delete(shardID uint8, key string) {
	delete(this.sharded[shardID], key)
}

// Set associates the specified value with the specified key in the ConcurrentMap.
// If the value is nil, the key-value pair is deleted from the map.
// It returns an error if the shard ID is out of range.
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

// BatchUpdate updates the values associated with the specified keys in the ConcurrentMap using the provided updater function.
// The updater function takes the original value, the index of the key in the keys slice, the key, and the new value as arguments.
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

// Traverse applies the specified operator function to each key-value pair in the ConcurrentMap.
// The operator function takes a key and a value as arguments and returns a new value and an optional result.
// The new value replaces the original value in the map, and the result is appended to the results slice.
// It returns a slice of result slices, one for each shard.
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

// BatchSet associates the specified values with the specified keys in the ConcurrentMap.
// If the values slice is shorter than the keys slice, the remaining keys are deleted from the map.
// If the values slice is longer than the keys slice, the extra values are ignored.
func (this *ConcurrentMap) BatchSet(keys []string, values []interface{}, args ...interface{}) {
	shardIDs := this.Hash8s(keys)

	if len(args) > 0 {
		this.DirectBatchSet(shardIDs, keys, values, args[0])
	} else {
		this.DirectBatchSet(shardIDs, keys, values)
	}
}

// DirectBatchSet associates the specified values with the specified shard IDs and keys in the ConcurrentMap.
// If the values slice is shorter than the keys slice, the remaining keys are deleted from the map.
// If the values slice is longer than the keys slice, the extra values are ignored.
func (this *ConcurrentMap) DirectBatchSet(shardIDs []uint8, keys []string, values []interface{}, args ...interface{}) {
	if len(keys) != len(values) {
		panic("Lengths don't match")
	}
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

// Keys returns a slice containing all the keys in the ConcurrentMap.
func (this *ConcurrentMap) Keys() []string {
	for i := 0; i < int(this.numShards); i++ {
		this.shardLocks[i].Lock()
		defer this.shardLocks[i].Unlock()
	}

	total := uint32(0)
	offsets := make([]uint32, this.numShards+1)
	for i := 0; i < int(this.numShards); i++ {
		total += uint32(len(this.sharded[i]))
		offsets[i+1] = total
	}

	keys := make([]string, total)
	// worker := func(start, end, index int, args ...interface{}) {
	// 	for i := start; i < end; i++ {
	// 		counter := offsets[i]
	// 		for k := range this.sharded[i] {
	// 			keys[counter] = k
	// 			counter++
	// 		}
	// 	}
	// }
	// common.ParallelWorker(len(this.sharded), len(this.sharded), worker)

	array.ParallelForeach(this.sharded, len(this.sharded), func(i int, _ *map[string]interface{}) {
		counter := offsets[i]
		for k := range this.sharded[i] {
			keys[counter] = k
			counter++
		}
	})

	return keys
}

// Hash8 calculates the shard ID for the specified key using a simple hash function.
func (this *ConcurrentMap) Hash8(key string) uint8 {
	if len(key) == 0 {
		return 0 //math.MaxUint8
	}

	var total uint32 = 0
	for j := 0; j < len(key); j++ {
		total += uint32(key[j])
	}
	return uint8(total % uint32(this.numShards))
}

// Hash8s calculates the shard IDs for the specified keys using the Hash8 function.
// It returns a slice of shard IDs in the same order as the keys.
func (this *ConcurrentMap) Hash8s(keys []string) []uint8 {
	shardIds := make([]uint8, len(keys))
	array.ParallelForeach(keys, 8, func(i int, _ *string) {
		shardIds[i] = this.Hash8(keys[i])
	})

	return shardIds
}

// Shards returns a pointer to the slice of maps representing the shards in the ConcurrentMap.
func (this *ConcurrentMap) Shards() *[]map[string]interface{} {
	return &this.sharded
}

// Find searches for the value in the ConcurrentMap that satisfies the specified comparison function.
// The comparison function takes two values as arguments and returns true if the first value is considered "less" than the second value.
// It returns the value that satisfies the comparison function, or nil if no such value is found.
func (this *ConcurrentMap) Find(Compare func(interface{}, interface{}) bool) interface{} {
	for i := 0; i < int(this.numShards); i++ {
		this.shardLocks[i].Lock()
		defer this.shardLocks[i].Unlock()
	}
	// this.shardLocks[start].RLock()
	// defer this.shardLocks[start].RUnlock()

	values := make([]interface{}, len(this.sharded))
	// worker := func(start, end, index int, args ...interface{}) {
	// 	this.shardLocks[start].RLock()
	// 	defer this.shardLocks[start].RUnlock()

	// 	for i := start; i < end; i++ {
	// 		for _, v := range this.sharded[i] {
	// 			if values[i] == nil || Compare(v, values[i]) {
	// 				values[i] = v
	// 			}
	// 		}
	// 	}
	// }
	// common.ParallelWorker(len(this.sharded), len(this.sharded), worker)

	array.ParallelForeach(this.sharded, len(this.sharded), func(i int, _ *map[string]interface{}) {
		for _, v := range this.sharded[i] {
			if values[i] == nil || Compare(v, values[i]) {
				values[i] = v
			}
		}
	})

	val := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] != nil && (val == nil || Compare(values[i], val)) {
			val = values[i]
		}
	}
	return val
}

// Foreach applies the specified predicate function to each value in the ConcurrentMap.
// The predicate function takes a value as an argument and returns a new value.
// If the new value is nil, the key-value pair is deleted from the map.
// Otherwise, the new value replaces the original value in the map.
func (this *ConcurrentMap) Foreach(predicate func(interface{}) interface{}) {
	for i := 0; i < int(this.numShards); i++ {
		this.shardLocks[i].Lock()
		defer this.shardLocks[i].Unlock()
	}

	array.ParallelForeach(this.sharded, len(this.sharded), func(i int, _ *map[string]interface{}) {
		for k, v := range this.sharded[i] {
			if v = predicate(v); v == nil {
				delete(this.sharded[i], k)
			} else {
				this.sharded[i][k] = v
			}
		}
	})

	// worker := func(start, end, index int, args ...interface{}) {
	// 	this.shardLocks[start].RLock()
	// 	defer this.shardLocks[start].RUnlock()

	// 	for i := start; i < end; i++ {
	// 		for k, v := range this.sharded[i] {
	// 			if v = predicate(v); v == nil {
	// 				delete(this.sharded[i], k)
	// 			} else {
	// 				this.sharded[i][k] = v
	// 			}
	// 		}
	// 	}
	// }
	// common.ParallelWorker(len(this.sharded), len(this.sharded), worker)
}

// ForeachDo applies the specified do function to each key-value pair in the ConcurrentMap.
// The do function takes a key and a value as arguments and performs some action.
func (this *ConcurrentMap) ForeachDo(do func(interface{}, interface{})) {
	for i := 0; i < len(this.sharded); i++ {
		this.shardLocks[i].RLock()
		for k, v := range this.sharded[i] {
			do(k, v)
		}
		this.shardLocks[i].RUnlock()
	}
}

// ParallelForeachDo applies the specified do function to each key-value pair in the ConcurrentMap in parallel.
// The do function takes a key and a value as arguments and performs some action.
func (this *ConcurrentMap) ParallelForeachDo(do func(interface{}, interface{})) {
	array.ParallelForeach(this.sharded, len(this.sharded), func(_ int, shard *map[string]interface{}) {
		for k, v := range *shard {
			do(k, v)
		}
	})
}

// KVs returns two slices: one containing all the keys in the ConcurrentMap, and one containing the corresponding values.
func (this *ConcurrentMap) KVs() ([]string, []interface{}) {
	keys := this.Keys()
	return keys, this.BatchGet(keys)
}

// Clear removes all key-value pairs from the ConcurrentMap.
// func (this *ConcurrentMap) Clear() {
// 	cleaner := func(start, end, index int, args ...interface{}) {
// 		this.shardLocks[start].RLock()
// 		defer this.shardLocks[start].RUnlock()
// 		this.sharded[start] = make(map[string]interface{}, 64)
// 	}
// 	common.ParallelWorker(len(this.sharded), len(this.sharded), cleaner)

// 	array.ParallelForeach(buffers, 6, func(i int, _ *[]byte) {
// 		v := (&Univalue{}).Decode(buffers[i])
// 		univalues[i] = v.(*Univalue)
// 	})
// }

/* -----------------------------------Debug Functions---------------------------------------------------------*/
type Encodable interface {
	Encode() []byte
}

// Dump returns two slices: one containing all the keys in the ConcurrentMap, and one containing the corresponding values.
// The keys are sorted in ascending order.
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

// Print prints all the key-value pairs in the ConcurrentMap to the standard output.
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

// Checksum calculates the checksum of the ConcurrentMap by concatenating the SHA256 hashes of the keys and values.
// It returns a 32-byte checksum.
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
