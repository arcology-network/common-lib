package cachedstorage

import (
	"crypto/sha256"
	"errors"
	"path/filepath"
	"sync"

	codec "github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	ccmap "github.com/arcology-network/common-lib/container/map"
	datacompression "github.com/arcology-network/common-lib/datacompression"
)

type DataStore struct {
	db     PersistentStorageInterface
	dblock sync.RWMutex

	cachePolicy    *CachePolicy
	compressionLut *datacompression.CompressionLut
	localCache     *ccmap.ConcurrentMap

	encoder func(interface{}) []byte
	decoder func([]byte) interface{}

	partitionIDs []uint8

	keyBuffer     []string
	valueBuffer   []interface{} //this should be binary
	encodedBuffer [][]byte      //this should be binary

	//dbfilter     DbFilter
	commitLock sync.RWMutex

	globalCache map[string]interface{}
	cacheGuard  sync.RWMutex
}

func NewDataStore(args ...interface{}) *DataStore {
	dataStore := DataStore{
		partitionIDs: make([]uint8, 0, 65536),
		localCache:   ccmap.NewConcurrentMap(),
		globalCache:  make(map[string]interface{}),
	}

	if len(args) > 0 && args[0] != nil {
		dataStore.compressionLut = args[0].(*datacompression.CompressionLut)
	}

	if len(args) > 1 && args[1] != nil {
		dataStore.cachePolicy = args[1].(*CachePolicy)
	}

	if len(args) > 2 && args[2] != nil {
		dataStore.db = args[2].(PersistentStorageInterface)
		dataStore.cachePolicy.Customize(dataStore.db)
	}

	if len(args) > 3 && args[3] != nil {
		dataStore.encoder = args[3].(func(interface{}) []byte)
	}

	if len(args) > 4 && args[4] != nil {
		dataStore.decoder = args[4].(func([]byte) interface{})
	}

	// if len(args) > 5 && args[5] != nil {
	// 	dataStore.dbfilter = DbFilter(args[5].(func(PersistentStorageInterface) bool))
	// }

	if dataStore.db != nil && (dataStore.encoder == nil || dataStore.decoder == nil) {
		panic("Error: DB Encoder or Decoder haven't been specified !")
	}

	return &dataStore
}

func (this *DataStore) Encoder() func(interface{}) []byte {
	return this.encoder
}

func (this *DataStore) Decoder() func([]byte) interface{} {
	return this.decoder
}

func (this *DataStore) Size() uint32 {
	return this.localCache.Size()
}

func (this *DataStore) Cache() *ccmap.ConcurrentMap {
	return this.localCache
}

func (this *DataStore) Checksum() [32]byte {
	return this.localCache.Checksum()
}

func (this *DataStore) Query(pattern string, condition func(string, string) bool) ([]string, [][]byte, error) {
	this.dblock.RLock()
	defer this.dblock.RUnlock()

	return this.db.Query(pattern, condition)

}

// Inject directly to the local cache.
func (this *DataStore) Inject(key string, v interface{}) error {
	if this.compressionLut != nil {
		key = this.compressionLut.CompressOnTemp([]string{key})[0]
		this.compressionLut.Commit()
	}

	err := this.localCache.Set(key, v)
	if err == nil {
		return this.batchWritePersistentStorage([]string{key}, [][]byte{this.encoder(v)})
	}
	return err
}

// Inject directly to the local cache.
func (this *DataStore) BatchInject(keys []string, values []interface{}) error {
	if this.compressionLut != nil {
		this.compressionLut.CompressOnTemp(keys)
		this.compressionLut.Commit()
	}

	this.localCache.BatchSet(keys, values) //need to write to the storage as well
	return this.batchWritePersistentStorage(keys, common.ArrayCastTo[interface{}, []byte](values, func(v interface{}) []byte { return this.encoder(v) }))
}

func (this *DataStore) prefetch(key string) (uint32, uint32, error) {
	if this.db == nil {
		return 0, 0, errors.New("Error: DB not found !")
	}

	pattern := filepath.Dir(key)

	this.dblock.RLock()
	prefetchedKeys, valBytes, err := this.db.Query(pattern, Under)
	this.dblock.RUnlock()

	prefetchedValues := make([]interface{}, len(valBytes))
	for i := 0; i < len(valBytes); i++ {
		prefetchedValues[i] = this.decoder(valBytes[i])
	}

	flags, count := this.cachePolicy.BatchCheckCapacity(prefetchedKeys, prefetchedValues) // need to check the cache status first
	if count > 0 {
		this.localCache.BatchSet(prefetchedKeys, prefetchedValues, flags) // Save to the local cache
	}
	return uint32(len(prefetchedKeys)), count, err
}

func (this *DataStore) fetchPersistentStorage(key string) (interface{}, error) {
	if this.db == nil {
		return nil, errors.New("Error: DB not found")
	}

	var value interface{}
	this.dblock.RLock()
	bytes, err := this.db.Get(key)
	this.dblock.RUnlock()

	if bytes != nil && err == nil { // Get from the cache
		value = this.decoder(bytes)
	}
	return value, err
}

func (this *DataStore) batchFetchPersistentStorage(keys []string) ([][]byte, error) {
	if this.db == nil {
		return nil, errors.New("Error: DB not found")
	}

	this.dblock.RLock()
	defer this.dblock.RUnlock()
	return this.db.BatchGet(keys) // Get from the cache
}

func (this *DataStore) batchWritePersistentStorage(keys []string, encodedValues [][]byte) error {
	if this.db == nil {
		return errors.New("Error: DB not found")
	}

	this.dblock.Lock()
	defer this.dblock.Unlock()
	return this.db.BatchSet(keys, encodedValues)
}

func (this *DataStore) addToCache(key string, value interface{}) {
	if value == nil || this.cachePolicy == nil || !this.cachePolicy.IsFullCache() {
		return
	}

	if !this.cachePolicy.CheckCapacity(key, value) { // Not full yet
		this.localCache.Set(key, value)
	}
}

func (this *DataStore) batchAddToCache(keys []string, values []interface{}) {
	if this.cachePolicy == nil || !this.cachePolicy.IsFullCache() {
		return
	}

	if flags, count := this.cachePolicy.BatchCheckCapacity(keys, values); count > 0 { // need to check the cache status first
		this.localCache.BatchSet(keys, values, flags)
	}
}

func (this *DataStore) Buffers() ([]string, []interface{}, [][]byte) {
	return this.keyBuffer, this.valueBuffer, this.encodedBuffer
}

func (this *DataStore) FillCache(path string) {

}

func (this *DataStore) Retrive(key string) (interface{}, error) {
	if this.compressionLut != nil {
		key = this.compressionLut.TryCompress(key) // Convert the key
	}

	// Read the local cache first
	if v, _ := this.localCache.Get(key); v != nil {
		return v, nil
	}

	// if v == nil && this.cachePolicy != nil && !this.cachePolicy.IsFullCache() {
	v, err := this.fetchPersistentStorage(key)
	if v != nil && err == nil {
		if this.cachePolicy.CheckCapacity(key, v) { // need to check the cache status first
			if err = this.localCache.Set(key, v); err != nil { // Save to the local cache
				return nil, err
			}
			this.addToCache(key, v) //update to the local cache and add all the missing values to the cache
		}
	}
	return v, err
}

func (this *DataStore) BatchRetrive(keys []string) []interface{} {
	this.commitLock.RLock()
	defer this.commitLock.RUnlock()
	if this.compressionLut != nil {
		keys = this.compressionLut.TryBatchCompress(keys)
	}

	values := this.localCache.BatchGet(keys) // From the local cache first
	if common.Count(values, nil) == 0 {
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
				values[idx] = this.decoder(data[i])
			}
		}
		this.batchAddToCache(keys, values) //update to the local cache and add all the missing values to the cache
	}
	return values
}

func (this *DataStore) CacheRetrive(key string, valueTransformer func(interface{}) interface{}) (interface{}, error) {
	this.cacheGuard.RLock()
	if v, ok := this.globalCache[key]; ok {
		this.cacheGuard.RUnlock()
		return v, nil
	}

	if v, err := this.Retrive(key); err != nil {
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
func (this *DataStore) Precommit(keys []string, values interface{}) {
	this.commitLock.Lock()
	this.keyBuffer = common.IfThenDo1st(
		this.compressionLut != nil,
		func() []string { return this.compressionLut.CompressOnTemp(codec.Strings(keys).Clone()) },
		keys)

	this.valueBuffer = values.([]interface{})
	this.encodedBuffer = make([][]byte, len(this.valueBuffer))
	for i := 0; i < len(this.valueBuffer); i++ {
		if this.valueBuffer[i] != nil {
			this.valueBuffer[i] = this.valueBuffer[i].(AccessibleInterface).Value() // Strip access info
			this.encodedBuffer[i] = this.encoder(this.valueBuffer[i])
		}
	}

	this.partitionIDs = make([]uint8, len(this.keyBuffer))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			this.partitionIDs[i] = this.localCache.Hash8(this.keyBuffer[i]) //Must use the compressed ky to compute the shard
		}
	}
	common.ParallelWorker(len(this.keyBuffer), 4, worker)
}

func (this *DataStore) Commit() error {
	defer this.commitLock.Unlock()

	this.localCache.DirectBatchSet(this.partitionIDs, this.keyBuffer, this.valueBuffer) // update the local cache

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
	// 	objs := make([]AccessibleInterface, len(nVals))
	// 	for i := range nVals {
	// 		objs[i] = nVals[i].(AccessibleInterface)
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
	k, vs := this.Dump()
	kData := codec.Strings(k).Flatten()
	vData := make([][]byte, len(vs))
	for i, v := range vs {
		vData[i] = this.encoder(v)
	}
	vData = append(vData, kData)
	return sha256.Sum256(codec.Byteset(vData).Flatten())
}

func (this *DataStore) Dump() ([]string, []interface{}) {
	return this.localCache.Dump()
}

func (this *DataStore) KVs() ([]string, []interface{}) {
	return this.localCache.KVs()
}

func (this *DataStore) CachePolicy() *CachePolicy {
	return this.cachePolicy
}
