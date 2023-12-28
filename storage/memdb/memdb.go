package memdb

import (
	ccmap "github.com/arcology-network/common-lib/container/map"
)

type MemDB struct {
	db *ccmap.ConcurrentMap
}

func NewMemDB() *MemDB {
	return &MemDB{
		db: ccmap.NewConcurrentMap(),
	}
}

func (this *MemDB) Set(key string, v []byte) error {
	return this.db.Set(key, v)
}

func (this *MemDB) Get(key string) ([]byte, error) {
	v, _ := this.db.Get(key)
	if v == nil {
		return nil, nil
	}
	return v.([]byte), nil
}

func (this *MemDB) BatchGet(keys []string) ([][]byte, error) {
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
	values := make([]interface{}, len(keys))
	for i, v := range byteset {
		if v != nil {
			values[i] = v
		}
	}

	this.db.BatchSet(keys, values)
	return nil
}

func (this *MemDB) Query(key string, functor func(string, string) bool) ([]string, [][]byte, error) {
	return []string{}, [][]byte{}, nil
}
