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
	db             PersistentStorageInterface
	encoder        func(interface{}) []byte
	decoder        func([]byte) interface{}
	localCache     *ccmap.ConcurrentMap
	cachePolicy    *CachePolicy
	compressionLut *datacompression.CompressionLut
	partitionIDs   []uint8
	keyBuffer      []string
	valueBuffer    []interface{}
	dbfilter       DbFilter
	commitLock     sync.RWMutex
	dblock         sync.RWMutex

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

	if len(args) > 5 && args[5] != nil {
		dataStore.dbfilter = DbFilter(args[5].(func(PersistentStorageInterface) bool))
	}

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

func (this *DataStore) LocalCache() *ccmap.ConcurrentMap {
	return this.localCache
}

func (this *DataStore) Checksum() [32]byte {
	return this.localCache.Checksum()
}

// Inject directly to the local cache.
func (this *DataStore) Inject(key string, v interface{}) {
	if this.compressionLut != nil {
		key = this.compressionLut.CompressOnTemp([]string{key})[0]
		this.compressionLut.Commit()
	}

	this.localCache.Set(key, v)
}

func (this *DataStore) Query(pattern string, condition func(string, string) bool) ([]string, [][]byte, error) {
	this.dblock.RLock()
	defer this.dblock.RUnlock()

	return this.db.Query(pattern, condition)

}

// Inject directly to the local cache.
func (this *DataStore) BatchInject(keys []string, values []interface{}) {
	if this.compressionLut != nil {
		this.compressionLut.CompressOnTemp(keys)
		this.compressionLut.Commit()
	}

	this.localCache.BatchSet(keys, values)
}

func (this *DataStore) batchWritePersistentStorage(keys []string, values []interface{}) error {
	if this.db == nil {
		return nil
	}

	this.dblock.Lock()
	go func() {
		byteset := make([][]byte, len(keys))
		encoder := func(start, end, index int, args ...interface{}) {
			for i := start; i < end; i++ {
				byteset[i] = this.encoder(values[i])
			}
		}

		common.ParallelWorker(len(keys), 4, encoder)

		this.db.BatchSet(keys, byteset)
		this.dblock.Unlock()
	}()

	return nil
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

func (this *DataStore) batchFetchPersistentStorage(keys []string) ([]interface{}, error) {
	if this.db == nil {
		return nil, errors.New("Error: DB not found")
	}

	values := make([]interface{}, len(keys))
	this.dblock.RLock()
	byteset, err := this.db.BatchGet(keys) // Get from the cache
	this.dblock.RUnlock()

	if err == nil {
		for i := 0; i < len(byteset); i++ {
			if byteset[i] != nil {
				values[i] = this.decoder(byteset[i])
			}
		}
	}
	return values, err
}

func (this *DataStore) addToCache(keys []string, values []interface{}) {
	if this.cachePolicy == nil {
		return
	}

	flags, count := this.cachePolicy.BatchCheckCapacity(keys, values) // need to check the cache status first
	if count > 0 {
		this.localCache.BatchSet(keys, values, flags)
	}
}

func (this *DataStore) FillCache(path string) {

}

func (this *DataStore) Retrive(key string) (interface{}, error) {
	if this.compressionLut != nil {
		key = this.compressionLut.TryCompress(key) // Convert the key
	}

	var err error
	// var ok bool
	v, _ := this.localCache.Get(key) // Read the local cache first
	if v == nil && this.cachePolicy != nil && !this.cachePolicy.IsFullCache() {
		if v, err = this.fetchPersistentStorage(key); err == nil && v != nil {
			if this.cachePolicy.CheckCapacity(key, v) { // need to check the cache status first
				if err = this.localCache.Set(key, v); err != nil { // Save to the local cache
					return nil, err
				}

				// if _, _, err := this.prefetch(key); err != nil { // Fetch the related entries
				// 	return nil, err
				// }
			}
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

	/* Find the missing values */
	queryKeys := make([]string, 0, len(keys))
	queryIdxes := make([]int, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		if values[i] == nil {
			queryKeys = append(queryKeys, keys[i])
			queryIdxes = append(queryIdxes, i)
		}
	}

	if len(queryKeys) == 0 || this.cachePolicy == nil || this.cachePolicy.IsFullCache() { // No missing values
		return values
	}

	/* Filter based on the persistent storage type */
	if this.dbfilter != nil {
		this.dblock.RLock()
		ok := this.dbfilter(this.db)
		this.dblock.RUnlock()
		if ok {
			return values
		}

	}

	/* Search in the persistent storage for the missing ones */
	if queryvalues, err := this.batchFetchPersistentStorage(queryKeys); err == nil {
		for i, idx := range queryIdxes {
			values[idx] = queryvalues[i]
		}
		this.addToCache(queryKeys, queryvalues) //adding to the local cache
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
	this.keyBuffer = keys
	if this.compressionLut != nil {
		this.keyBuffer = this.compressionLut.CompressOnTemp(codec.Strings(keys).Clone())
	}

	this.valueBuffer = values.([]interface{})
	for i := 0; i < len(this.valueBuffer); i++ {
		if this.valueBuffer[i] != nil {
			this.valueBuffer[i] = this.valueBuffer[i].(AccessibleInterface).Value() // Strip access info
		} else {
			this.valueBuffer[i] = nil
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

	this.localCache.DirectBatchSet(this.partitionIDs, this.keyBuffer, this.valueBuffer)

	var err error
	if this.compressionLut != nil {
		common.ParallelExecute(
			//func() { this.localCache.DirectBatchSet(this.partitionIDs, this.keyBuffer, this.valueBuffer) },
			func() { err = this.batchWritePersistentStorage(this.keyBuffer, this.valueBuffer) }, // Write data back
			func() { this.compressionLut.Commit() })

	} else {
		// common.ParallelExecute(
		// 	func() { this.localCache.DirectBatchSet(this.partitionIDs, this.keyBuffer, this.valueBuffer) },
		// 	func() { err = this.batchWritePersistentStorage(this.keyBuffer, this.valueBuffer) },
		// )
		err = this.batchWritePersistentStorage(this.keyBuffer, this.valueBuffer)
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
	return this.CachePolicy().Refresh(this.LocalCache())
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
