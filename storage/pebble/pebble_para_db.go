/*
 *   Copyright (c) 2026 Arcology Network
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package pebbledb

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	common "github.com/arcology-network/common-lib/common"
	slice "github.com/arcology-network/common-lib/exp/slice"
	stgintf "github.com/arcology-network/common-lib/storage/interface"
)

type ParaPebbleDB struct {
	impls      [16]*PebbleDB
	shardLocks [16]sync.RWMutex
	shardFunc  func(int, string) int
}

func NewParaPebbleDB(root string, shardFunc func(numOfShard int, key string) int, decoder ...func(string, any, any) (any, error)) (*ParaPebbleDB, error) {
	var paraPebbleDB ParaPebbleDB
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err := os.MkdirAll(root, fs.ModePerm); err != nil {
			return nil, err
		}
	}

	var decode func(string, any, any) (any, error)
	if len(decoder) > 0 {
		decode = decoder[0]
	}

	for i := 0; i < len(paraPebbleDB.impls); i++ {
		path := filepath.Join(root, fmt.Sprint(i))
		db, err := NewPebbleDB(path, decode)
		if err != nil {
			return nil, err
		}
		paraPebbleDB.impls[i] = db
	}

	if shardFunc != nil {
		paraPebbleDB.shardFunc = shardFunc
	} else {
		paraPebbleDB.shardFunc = func(n int, key string) int {
			total := 0
			for i := 0; i < len(key); i++ {
				total += int(key[i])
			}
			return total % n
		}
	}
	return &paraPebbleDB, nil
}

func (this *ParaPebbleDB) Get(key string) (any, error) {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.Get(key)
}

func (this *ParaPebbleDB) GetAs(key string, typeHint any) (any, error) {
	if this == nil {
		return nil, stgintf.ErrNotFound
	}

	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.GetAs(key, typeHint)
}

func (this *ParaPebbleDB) Has(key string) bool {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.Has(key)
}

func (this *ParaPebbleDB) Set(key string, value []byte) error {
	idx, db := this.getShard(key)
	this.shardLocks[idx].Lock()
	defer this.shardLocks[idx].Unlock()
	return db.Set(key, value)
}

func (this *ParaPebbleDB) Delete(key string) error {
	idx, db := this.getShard(key)
	this.shardLocks[idx].Lock()
	defer this.shardLocks[idx].Unlock()
	return db.Delete(key)
}

func (this *ParaPebbleDB) DeleteBatch(keys []string) []error {
	categorized := make([][]string, len(this.impls))
	origIdx := make([][]int, len(this.impls))
	for i := range categorized {
		categorized[i] = make([]string, 0, len(keys)/len(this.impls)+1)
		origIdx[i] = make([]int, 0, len(keys)/len(this.impls)+1)
	}

	for i, key := range keys {
		idx, _ := this.getShard(key)
		categorized[idx] = append(categorized[idx], key)
		origIdx[idx] = append(origIdx[idx], i)
	}

	errs := make([]error, len(keys))
	for i := range categorized {
		if len(categorized[i]) == 0 {
			continue
		}
		this.shardLocks[i].Lock()
		shardErrs := this.impls[i].DeleteBatch(categorized[i])
		this.shardLocks[i].Unlock()
		for j, err := range shardErrs {
			errs[origIdx[i][j]] = err
		}
	}
	return errs
}

func (this *ParaPebbleDB) GetBatch(keys []string) ([]any, []error) {
	categorizedKeys := make([][]string, len(this.impls))
	categorizedIdx := make([][]int, len(this.impls))
	for i := range categorizedKeys {
		categorizedKeys[i] = make([]string, 0, len(keys)/len(this.impls)+1)
		categorizedIdx[i] = make([]int, 0, len(keys)/len(this.impls)+1)
	}

	for i, key := range keys {
		idx, _ := this.getShard(key)
		categorizedKeys[idx] = append(categorizedKeys[idx], key)
		categorizedIdx[idx] = append(categorizedIdx[idx], i)
	}

	shardResults := make([][]any, len(this.impls))
	shardErrSlices := make([][]error, len(this.impls))
	finder := func(start, end, _ int, _ ...interface{}) {
		for i := start; i < end; i++ {
			this.shardLocks[i].RLock()
			shardResults[i], shardErrSlices[i] = this.impls[i].GetBatch(categorizedKeys[i])
			this.shardLocks[i].RUnlock()
		}
	}
	common.ParallelWorker(len(this.impls), len(this.impls), finder)

	values := make([]any, len(keys))
	errs := make([]error, len(keys))
	for i := range categorizedIdx {
		for j, orig := range categorizedIdx[i] {
			values[orig] = shardResults[i][j]
			errs[orig] = shardErrSlices[i][j]
		}
	}
	return values, errs
}

func (this *ParaPebbleDB) SetBatch(keys []string, values [][]byte) []error {
	categorizedKeys := make([][]string, len(this.impls))
	categorizedVals := make([][][]byte, len(this.impls))
	origIdx := make([][]int, len(this.impls))
	for i := range categorizedKeys {
		categorizedKeys[i] = make([]string, 0, len(keys)/len(this.impls)+1)
		categorizedVals[i] = make([][]byte, 0, len(keys)/len(this.impls)+1)
		origIdx[i] = make([]int, 0, len(keys)/len(this.impls)+1)
	}

	for i, key := range keys {
		idx, _ := this.getShard(key)
		categorizedKeys[idx] = append(categorizedKeys[idx], key)
		categorizedVals[idx] = append(categorizedVals[idx], values[i])
		origIdx[idx] = append(origIdx[idx], i)
	}

	perShardErrs := slice.ParallelTransform(categorizedKeys, len(categorizedKeys), func(i int, _ []string) []error {
		this.shardLocks[i].Lock()
		defer this.shardLocks[i].Unlock()
		return this.impls[i].SetBatch(categorizedKeys[i], categorizedVals[i])
	})

	errs := make([]error, len(keys))
	for i, shardErrs := range perShardErrs {
		for j, err := range shardErrs {
			errs[origIdx[i][j]] = err
		}
	}
	return errs
}

func (this *ParaPebbleDB) Query(prefix string, checker func(string, []byte) bool) ([]string, [][]byte, error) {
	keys := make([]string, 0)
	values := make([][]byte, 0)
	for i, db := range this.impls {
		if db == nil {
			continue
		}
		this.shardLocks[i].RLock()
		shardKeys, shardValues, err := db.Query(prefix, checker)
		this.shardLocks[i].RUnlock()
		if err != nil {
			return nil, nil, err
		}
		keys = append(keys, shardKeys...)
		values = append(values, shardValues...)
	}
	return keys, values, nil
}

func (this *ParaPebbleDB) Close() error {
	for _, db := range this.impls {
		if db == nil {
			continue
		}
		if err := db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (this *ParaPebbleDB) getShard(key string) (int, *PebbleDB) {
	shardIdx := this.shardFunc(len(this.impls), key)
	return shardIdx, this.impls[shardIdx]
}
