/*
*   Copyright (c) 2024 Arcology Network

*   This program is free software: you can redistribute it and/or modify
*   it under the terms of the GNU General Public License as published by
*   the Free Software Foundation, either version 3 of the License, or
*   (at your option) any later version.

*   This program is distributed in the hope that it will be useful,
*   but WITHOUT ANY WARRANTY; without even the implied warranty of
*   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*   GNU General Public License for more details.

*   You should have received a copy of the GNU General Public License
*   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package badgerdb

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"sync"

	common "github.com/arcology-network/common-lib/common"
	stgintf "github.com/arcology-network/common-lib/storage/interface"
)

var _ stgintf.ReadWriteStore[string, []byte] = (*ParaBadgerDB)(nil)

type ParaBadgerDB struct {
	impls      [16]*BadgerDB
	shardLocks [16]sync.RWMutex
	shardFunc  func(int, string) int
}

func NewParaBadgerDB(root string, shardFunc func(numOfShard int, key string) int) *ParaBadgerDB {
	var paraBadgerDB ParaBadgerDB
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err := os.MkdirAll(root, fs.ModePerm); err != nil {
			panic(err)
		}
	}

	for i := 0; i < len(paraBadgerDB.impls); i++ {
		path := path.Join(root+fmt.Sprint(i)) + "/"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, fs.ModePerm); err != nil {
				panic(err)
			}
		}
		paraBadgerDB.impls[i] = NewBadgerDB(path)
	}

	if shardFunc != nil {
		paraBadgerDB.shardFunc = shardFunc
	} else {
		paraBadgerDB.shardFunc = paraBadgerDB.hash32
	}
	return &paraBadgerDB
}

func (this *ParaBadgerDB) Get(key string) (value any, err error) {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.Get(key)
}

func (this *ParaBadgerDB) Has(key string) bool {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.Has(key)
}

func (this *ParaBadgerDB) Set(key string, value []byte) error {
	idx, db := this.getShard(key)
	this.shardLocks[idx].Lock()
	defer this.shardLocks[idx].Unlock()
	return db.Set(key, value)
}

func (this *ParaBadgerDB) Delete(key string) error {
	idx, db := this.getShard(key)
	this.shardLocks[idx].Lock()
	defer this.shardLocks[idx].Unlock()
	return db.Delete(key)
}

func (this *ParaBadgerDB) DeleteBatch(keys []string) []error {
	errs := make([]error, len(keys))
	for i := 0; i < len(keys); i++ {
		idx, db := this.getShard(keys[i])
		this.shardLocks[idx].Lock()
		err := db.Delete(keys[i])
		this.shardLocks[idx].Unlock()
		if err != nil {
			errs[i] = err
		}
	}
	if allNil(errs) {
		return nil
	}
	return errs
}

func (this *ParaBadgerDB) GetBatch(keys []string) ([]any, []error) {
	results := make([]any, len(keys))
	errs := make([]error, len(keys))
	finder := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			idx, db := this.getShard(keys[i])
			this.shardLocks[idx].RLock()
			v, err := db.Get(keys[i])
			this.shardLocks[idx].RUnlock()
			if err != nil {
				errs[i] = err
				continue
			}
			results[i] = v
		}
	}
	common.ParallelWorker(len(keys), len(this.impls), finder)
	if allNil(errs) {
		return results, nil
	}
	return results, errs
}

func (this *ParaBadgerDB) SetBatch(keys []string, values [][]byte) []error {
	errs := make([]error, len(keys))
	finder := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			idx, db := this.getShard(keys[i])
			this.shardLocks[idx].Lock()
			err := db.Set(keys[i], values[i])
			this.shardLocks[idx].Unlock()
			if err != nil {
				errs[i] = err
			}
		}
	}
	common.ParallelWorker(len(keys), len(this.impls), finder)
	if allNil(errs) {
		return nil
	}
	return errs
}

func (this *ParaBadgerDB) Query(prefix string, checker func(string, []byte) bool) (keys []string, values [][]byte, errs []error) {
	shardIdx, db := this.getShard(prefix)

	this.shardLocks[shardIdx].RLock()
	defer this.shardLocks[shardIdx].RUnlock()
	return db.Query(prefix, checker)
}

func (this *ParaBadgerDB) Close() error {
	for _, db := range this.impls {
		db.Close()
	}
	return nil
}

func (this *ParaBadgerDB) getShard(key string) (int, *BadgerDB) {
	shardIdx := this.shardFunc(len(this.impls), key)
	return shardIdx, this.impls[shardIdx]
}

func (this *ParaBadgerDB) hash32(numOfShard int, key string) int {
	if len(key) == 0 {
		return int(math.MaxUint32 % uint32(numOfShard))
	}

	var total int = 0
	for j := 0; j < len(key); j++ {
		total += int(key[j])
	}
	return total % numOfShard
}
