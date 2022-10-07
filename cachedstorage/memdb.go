package cachedstorage

import (
	"sync"

	ccctrn "github.com/HPISTechnologies/common-lib/concurrentcontainer"
)

type MemDB struct {
	mutex sync.RWMutex
	db    *ccctrn.ConcurrentMap
}

func NewMemDB() *MemDB {
	return &MemDB{
		db: ccctrn.NewConcurrentMap(),
	}
}

func (this *MemDB) Set(key string, v []byte) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.db.Set(key, v)
}

func (this *MemDB) Get(key string) ([]byte, error) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	v, _ := this.db.Get(key)
	return v.([]byte), nil
}

func (this *MemDB) BatchGet(keys []string) ([][]byte, error) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	values := this.db.BatchGet(keys)
	byteset := make([][]byte, len(keys))
	for i, v := range values {
		if v != nil {
			byteset[i] = v.([]byte)
		}
	}
	return byteset, nil
}

func (this *MemDB) BatchSet(keys []string, byteset [][]byte) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	values := make([]interface{}, len(keys))
	for i, v := range byteset {
		if v != nil {
			values[i] = v
		}
	}

	this.db.BatchSet(keys, values)
	return nil
}
