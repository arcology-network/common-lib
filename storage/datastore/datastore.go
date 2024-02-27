package datastore

import (
	"crypto/sha256"
	"errors"
	"sync"

	addrcompressor "github.com/arcology-network/common-lib/addrcompressor"
	codec "github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	"github.com/cespare/xxhash/v2"

	// expmap "github.com/arcology-network/common-lib/container/map"
	expmap "github.com/arcology-network/common-lib/exp/map"
	slice "github.com/arcology-network/common-lib/exp/slice"
	intf "github.com/arcology-network/common-lib/storage/interface"
)

type DataStore struct {
	db   intf.PersistentStorage
	lock sync.RWMutex

	cachePolicy      *CachePolicy
	compressionLut   *addrcompressor.CompressionLut
	localCache       *expmap.ConcurrentMap[string, any]
	maxCacheCapacity int

	encoder func(string, interface{}) []byte
	decoder func(string, []byte, any) interface{}

	partitionIDs []uint64

	keyBuffer     []string
	valueBuffer   []interface{} //this should be binary
	encodedBuffer [][]byte      //this should be binary

	//dbfilter     DbFilter
	commitLock sync.RWMutex

	globalCache map[string]interface{}
	cacheGuard  sync.RWMutex
}

func NewDataStore(
	compressionLut *addrcompressor.CompressionLut,
	cachePolicy *CachePolicy,
	db intf.PersistentStorage,
	encoder func(string, any) []byte,
	decoder func(string, []byte, any) interface{},
) *DataStore {
	dataStore := &DataStore{
		partitionIDs: make([]uint64, 0, 65536),
		localCache: expmap.NewConcurrentMap(8, func(v any) bool { return v == nil }, func(k string) uint64 {
			return xxhash.Sum64([]byte(k))
		}),
		globalCache:    make(map[string]interface{}),
		compressionLut: compressionLut,
		cachePolicy:    cachePolicy,
		db:             db,
		encoder:        encoder,
		decoder:        decoder,
	}

	dataStore.cachePolicy.Customize(dataStore.db)
	return dataStore
}

// Pleaseholder only
func (this *DataStore) WriteEthTries(...interface{}) [32]byte {
	this.commitLock.Lock()
	return [32]byte{}
}

func (this *DataStore) Encoder() func(string, interface{}) []byte {
	return this.encoder
}

func (this *DataStore) Decoder() func(string, []byte, any) interface{} {
	return this.decoder
}

func (this *DataStore) Size() uint32 {
	return this.localCache.Size()
}

func (this *DataStore) Cache() *expmap.ConcurrentMap[string, any] {
	return this.localCache
}

func (this *DataStore) GetMaxCacheCapacity() int {
	return this.maxCacheCapacity
}

func (this *DataStore) SetMaxCacheCapacity(size int) {
	this.maxCacheCapacity = 0
}

func (this *DataStore) Checksum() [32]byte {
	return this.localCache.Checksum()
}

func (this *DataStore) Query(pattern string, condition func(string, string) bool) ([]string, [][]byte, error) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.db.Query(pattern, condition)
}

func (this *DataStore) IfExists(key string) bool {
	v, _ := this.Retrive(key, nil)
	return v != nil
}

// Inject directly to the local cache.
func (this *DataStore) Inject(key string, v interface{}) error {
	if this.compressionLut != nil {
		key = this.compressionLut.CompressOnTemp([]string{key})[0]
		this.compressionLut.Commit()
	}

	this.addToCache(key, v)

	// err := this.localCache.Set(key, v)
	// if err == nil {
	return this.batchWritePersistentStorage([]string{key}, [][]byte{this.encoder(key, v)})
	// }
	// return err
}

// Inject directly to the local cache.
func (this *DataStore) BatchInject(keys []string, values []interface{}) error {
	if this.compressionLut != nil {
		this.compressionLut.CompressOnTemp(keys)
		this.compressionLut.Commit()
	}

	this.batchAddToCache(this.GetParitions(keys), keys, values)
	encoded := make([][]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		encoded[i] = this.encoder(keys[i], values[i])
	}
	return this.batchWritePersistentStorage(keys, encoded)
}

// func (this *DataStore) prefetch(key string) (uint32, uint32, error) {
// 	if this.db == nil {
// 		return 0, 0, errors.New("Error: DB not found !")
// 	}

// 	pattern := filepath.Dir(key)

// 	this.lock.RLock()
// 	prefetchedKeys, valBytes, err := this.db.Query(pattern, Under)
// 	this.lock.RUnlock()

// 	prefetchedValues := make([]interface{}, len(valBytes))
// 	for i := 0; i < len(valBytes); i++ {
// 		prefetchedValues[i] = this.decoder(valBytes[i])
// 	}

// 	flags, count := this.cachePolicy.BatchCheckCapacity(prefetchedKeys, prefetchedValues) // need to check the cache status first
// 	if count > 0 {
// 		this.localCache.BatchSet(prefetchedKeys, prefetchedValues, flags) // Save to the local cache
// 	}
// 	return uint32(len(prefetchedKeys)), count, err
// }

func (this *DataStore) fetchPersistentStorage(key string, T any) (interface{}, error) {
	if this.db == nil {
		return nil, errors.New("Error: DB not found")
	}

	this.lock.RLock()
	bytes, err := this.db.Get(key)
	this.lock.RUnlock()

	if len(bytes) > 0 && err == nil { // Get from the cache
		if T == nil {
			return bytes, nil
		}
		return this.decoder(key, bytes, T), nil
	}
	return nil, err
}

func (this *DataStore) batchFetchPersistentStorage(keys []string) ([][]byte, error) {
	if this.db == nil {
		return nil, errors.New("Error: DB not found")
	}

	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.db.BatchGet(keys) // Get from the cache
}

func (this *DataStore) batchWritePersistentStorage(keys []string, encodedValues [][]byte) error {
	if this.db == nil {
		return errors.New("Error: DB not found")
	}

	this.lock.Lock()
	defer this.lock.Unlock()
	return this.db.BatchSet(keys, encodedValues)
}

func (this *DataStore) addToCache(key string, value interface{}) {
	if this.cachePolicy == nil || this.cachePolicy.IsFull() {
		return
	}

	if this.cachePolicy.InfinitCache() {
		this.localCache.Set(key, value)
		return
	}

	if !this.cachePolicy.CheckCapacity(key, value) { // Not full yet
		this.localCache.Set(key, value)
	}
}

func (this *DataStore) batchAddToCache(ids []uint64, keys []string, values []interface{}) {
	if this.cachePolicy == nil {
		return
	}

	if this.cachePolicy.InfinitCache() {
		this.localCache.DirectBatchSet(ids, keys, values)
		return
	}

	if _, count, all := this.cachePolicy.BatchCheckCapacity(keys, values); all || count > 0 { // need to check the cache status first
		this.localCache.DirectBatchSet(ids, keys, values)
	}
}

func (this *DataStore) Buffers() ([]string, []interface{}, [][]byte) {
	return this.keyBuffer, this.valueBuffer, this.encodedBuffer
}

func (this *DataStore) FillCache(path string) {

}

func (this *DataStore) Retrive(key string, T any) (interface{}, error) {
	if this.compressionLut != nil {
		key = this.compressionLut.TryCompress(key) // Convert the key
	}

	// Read the local cache first
	if v, _ := this.localCache.Get(key); v != nil {
		return v, nil
	}

	// if v == nil && this.cachePolicy != nil && !this.cachePolicy.InfinitCache() {
	v, err := this.fetchPersistentStorage(key, T)
	if err == nil {
		// if this.cachePolicy.CheckCapacity(key, v) { // need to check the cache status first
		// if err = this.localCache.Set(key, v); err != nil { // Save to the local cache
		// 	return nil, err
		// }
		this.addToCache(key, v) //update to the local cache and add all the missing values to the cache
		// }
	}
	return v, err
}

func (this *DataStore) BatchRetrive(keys []string, T []any) []interface{} {
	this.commitLock.RLock()
	defer this.commitLock.RUnlock()
	if this.compressionLut != nil {
		keys = this.compressionLut.TryBatchCompress(keys)
	}

	values := common.FilterFirst(this.localCache.BatchGet(keys)) // From the local cache first
	if slice.Count(values, nil) == 0 {                           // All found
		return values
	}

	/* Find the values missing from the local cache*/
	queryKeys, queryIdxes := make([]string, 0, len(keys)), make([]int, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		if values[i] == nil {
			queryKeys = append(queryKeys, keys[i])
			queryIdxes = append(queryIdxes, i)
		}
	}

	if data, err := this.batchFetchPersistentStorage(queryKeys); err == nil { // search for the values that aren't in the cache
		for i, idx := range queryIdxes {
			if data[i] != nil {
				if len(T) > 0 {
					values[idx] = this.decoder(queryKeys[i], data[i], T[i])
				} else {
					values[idx] = this.decoder(queryKeys[i], data[i], nil)
				}
			}
		}
		this.batchAddToCache(this.GetParitions(keys), keys, values) //update to the local cache and add all the missing values to the cache
	}
	return values
}

func (this *DataStore) CacheRetrive(key string, valueTransformer func(interface{}) interface{}) (interface{}, error) {
	this.cacheGuard.RLock()
	if v, ok := this.globalCache[key]; ok {
		this.cacheGuard.RUnlock()
		return v, nil
	}

	if v, err := this.Retrive(key, nil); err != nil {
		this.cacheGuard.RUnlock()
		return nil, err
	} else {
		this.cacheGuard.RUnlock()
		this.cacheGuard.Lock()
		tv := valueTransformer(v)
		this.globalCache[key] = tv
		this.cacheGuard.Unlock()
		return tv, nil
	}
}

func (this *DataStore) Clear() {
	this.partitionIDs = this.partitionIDs[:0]
	this.keyBuffer = this.keyBuffer[:0]
	this.valueBuffer = this.valueBuffer[:0]
	this.globalCache = make(map[string]interface{})
}

// Get the shard ids, values, and preupdate the compression dict
func (this *DataStore) Precommit(args ...interface{}) [32]byte {
	keys, values := args[0].([]string), args[1].([]interface{})

	this.commitLock.Lock()
	this.keyBuffer = common.IfThenDo1st(
		this.compressionLut != nil,
		func() []string { return this.compressionLut.CompressOnTemp(codec.Strings(keys).Clone()) },
		keys)

	this.valueBuffer = values
	this.encodedBuffer = make([][]byte, len(this.valueBuffer))
	for i := 0; i < len(this.valueBuffer); i++ {
		if this.valueBuffer[i] != nil {
			this.valueBuffer[i] = this.valueBuffer[i].(intf.Accessible).Value() // Strip access info
			this.encodedBuffer[i] = this.encoder(keys[i], this.valueBuffer[i])
		}
	}
	this.partitionIDs = this.GetParitions(keys)
	return [32]byte{}
}

func (this *DataStore) GetParitions(keys []string) []uint64 {
	partitionIDs := make([]uint64, len(keys))
	// worker := func(start, end, index int, args ...interface{}) {
	// 	for i := start; i < end; i++ {
	// 		partitionIDs[i] = this.localCache.Hash8(keys[i]) //Must use the compressed ky to compute the shard
	// 	}
	// }
	// common.ParallelWorker(len(keys), 4, worker)

	slice.ParallelForeach(keys, 4, func(i int, _ *string) {
		partitionIDs[i] = this.localCache.Hash(keys[i]) //Must use the compressed ky to compute the shard
	})

	return partitionIDs
}

func (this *DataStore) Commit(_ uint64) error {
	defer this.commitLock.Unlock()

	this.batchAddToCache(this.partitionIDs, this.keyBuffer, this.valueBuffer)
	// this.localCache.DirectBatchSet(this.partitionIDs, this.keyBuffer, this.valueBuffer) // update the local cache

	var err error
	if this.compressionLut != nil {
		common.ParallelExecute(
			//func() { this.localCache.DirectBatchSet(this.partitionIDs, this.keyBuffer, this.valueBuffer) },
			func() { err = this.batchWritePersistentStorage(this.keyBuffer, this.encodedBuffer) }, // Write data back
			func() { this.compressionLut.Commit() })

	} else {
		// common.ParallelExecute(
		// 	func() { this.localCache.DirectBatchSet(this.partitionIDs, this.keyBuffer, this.valueBuffer) },
		// 	func() { err = this.batchWritePersistentStorage(this.keyBuffer, this.valueBuffer) },
		// )
		err = this.batchWritePersistentStorage(this.keyBuffer, this.encodedBuffer)
	}
	this.Clear()
	return err
}

func (this *DataStore) UpdateCacheStats(nVals []interface{}) {
	// if this.cachePolicy != nil {
	// 	objs := make([]Accessible, len(nVals))
	// 	for i := range nVals {
	// 		objs[i] = nVals[i].(Accessible)
	// 	}
	// 	this.CachePolicy().AddToStats(keys, objs)
	// }
}

func (this *DataStore) RefreshCache() (uint64, uint64) {
	return this.CachePolicy().Refresh(this.Cache())
}

func (this *DataStore) Print() {
	this.localCache.Print()
}

func (this *DataStore) CheckSum() [32]byte {
	k, vs := this.KVs()
	kData := codec.Strings(k).Flatten()
	vData := make([][]byte, len(vs))
	for i, v := range vs {
		vData[i] = this.encoder(k[i], v)
	}
	vData = append(vData, kData)
	return sha256.Sum256(codec.Byteset(vData).Flatten())
}

func (this *DataStore) KVs() ([]string, []interface{}) {
	return this.localCache.KVs()
}

func (this *DataStore) CachePolicy() *CachePolicy {
	return this.cachePolicy
}
