package cachedstorage

import (
	"crypto/sha256"
	"errors"
	"sync"

	codec "github.com/HPISTechnologies/common-lib/codec"
	common "github.com/HPISTechnologies/common-lib/common"
	cccontainer "github.com/HPISTechnologies/common-lib/concurrentcontainer"
	datacompression "github.com/HPISTechnologies/common-lib/datacompression"
)

type DataStore struct {
	db             PersistentStorageInterface
	encoder        func(interface{}) []byte
	decoder        func([]byte) interface{}
	localCache     *cccontainer.ConcurrentMap
	cachePolicy    *CachePolicy
	compressionLut *datacompression.CompressionLut
	partitionIDs   []uint8
	keyBuffer      []string
	valueBuffer    []interface{}
	commitLock     sync.RWMutex
}

func NewDataStore(args ...interface{}) *DataStore {
	dataStore := DataStore{
		partitionIDs: make([]uint8, 0, 65536),
		localCache:   cccontainer.NewConcurrentMap(),
	}

	if len(args) > 0 && args[0] != nil {
		dataStore.compressionLut = args[0].(*datacompression.CompressionLut)
	}

	if len(args) > 1 && args[1] != nil {
		dataStore.cachePolicy = args[1].(*CachePolicy)
	}

	if len(args) > 2 && args[2] != nil {
		dataStore.db = args[2].(PersistentStorageInterface)
	}

	if len(args) > 3 && args[3] != nil {
		dataStore.encoder = args[3].(func(interface{}) []byte)
	}

	if len(args) > 4 && args[4] != nil {
		dataStore.decoder = args[4].(func([]byte) interface{})
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

func (this *DataStore) LocalCache() *cccontainer.ConcurrentMap {
	return this.localCache
}

func (this *DataStore) Checksum() [32]byte {
	return this.localCache.Checksum()
}

// Special interface to inject directly to the local cache.
func (this *DataStore) Inject(key string, v interface{}) {
	if this.compressionLut != nil {
		key = this.compressionLut.CompressOnTemp([]string{key})[0]
		this.compressionLut.Commit()
	}

	this.localCache.Set(key, v)
}

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

	byteset := make([][]byte, len(keys))
	encoder := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			byteset[i] = this.encoder(values[i])
		}
	}
	common.ParallelWorker(len(keys), 4, encoder)
	return this.db.BatchSet(keys, byteset)
}

func (this *DataStore) readPersistentStorage(key string) (interface{}, error) {
	if this.db == nil {
		return nil, errors.New("Error: DB not found")
	}

	if bytes, err := this.db.Get(key); err == nil { // Get from the cache
		return this.decoder(bytes), nil
	} else {
		return nil, errors.New("Error: Failed to get from the DB")
	}
}

func (this *DataStore) batchReadPersistentStorage(keys []string) ([]interface{}, error) {
	if this.db == nil {
		return nil, errors.New("Error: DB not found")
	}

	values := make([]interface{}, len(keys))
	byteset, err := this.db.BatchGet(keys) // Get from the cache
	if err == nil {
		for i := 0; i < len(byteset); i++ {
			values[i] = this.decoder(byteset[i])
		}
	}
	return values, err
}

func (this *DataStore) Retrive(key string) interface{} {
	this.commitLock.RLock()
	defer this.commitLock.RUnlock()

	if this.compressionLut != nil {
		key = this.compressionLut.TryCompress(key) // Convert the key
	}

	v, _ := this.localCache.Get(key)
	if v == nil {
		if v, err := this.readPersistentStorage(key); err == nil && v != nil {
			if this.cachePolicy.CheckCapacity(key, v) { // need to check the cache status first
				this.localCache.Set(key, v) // Save to local cache if isn't full yet,
			}
			return v
		}
	}
	return v
}

func (this *DataStore) BatchRetrive(keys []string) []interface{} {
	this.commitLock.RLock()
	defer this.commitLock.RUnlock()
	if this.compressionLut != nil {
		keys = this.compressionLut.TryBatchCompress(keys)
	}

	missing := false
	values := this.localCache.BatchGet(keys)
	for i := 0; i < len(keys); i++ {
		if values[i] != nil {
			keys[i] = ""
		} else {
			missing = true
		}
	}

	if !missing { // No missing values
		return values
	}

	values, _ = this.batchReadPersistentStorage(keys)
	this.cachePolicy.BatchCheckCapacity(keys, values)
	this.localCache.BatchSet(keys, values) // Save to the local cache, need to check the cache status first
	return values
}

func (this *DataStore) Clear() {
	this.partitionIDs = this.partitionIDs[:0]
	this.keyBuffer = this.keyBuffer[:0]
	this.valueBuffer = this.valueBuffer[:0]
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
			this.valueBuffer[i] = this.valueBuffer[i].(AccessableInterface).Value() // Strip access info
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
	var err error
	if this.compressionLut != nil {
		common.ParallelExecute(
			func() { this.localCache.DirectBatchSet(this.partitionIDs, this.keyBuffer, this.valueBuffer) },
			func() { err = this.batchWritePersistentStorage(this.keyBuffer, this.valueBuffer) }, // Write data back
			func() { this.compressionLut.Commit() })
	} else {
		common.ParallelExecute(
			func() { this.localCache.DirectBatchSet(this.partitionIDs, this.keyBuffer, this.valueBuffer) },
			func() { err = this.batchWritePersistentStorage(this.keyBuffer, this.valueBuffer) })
	}

	this.Clear()
	return err
}

func (this *DataStore) UpdateCacheStats(keys []string, nVals []interface{}) {
	if this.cachePolicy != nil {
		this.CachePolicy().AddToBuffer(keys, nVals)
	}
}

func (this *DataStore) RefreshCache() {
	if this.cachePolicy != nil {
		this.CachePolicy().Refresh(this.LocalCache())
		this.CachePolicy().FreeMemory(this.LocalCache())
	}
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
