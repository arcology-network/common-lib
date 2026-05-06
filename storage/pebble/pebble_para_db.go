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
)

type ParaPebbleDB struct {
	impls      [16]*PebbleDB
	shardLocks [16]sync.RWMutex
	shardFunc  func(int, []byte) int
}

func NewParaPebbleDB(root string, shardFunc func(numOfShard int, key []byte) int) (*ParaPebbleDB, error) {
	var paraPebbleDB ParaPebbleDB
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err := os.MkdirAll(root, fs.ModePerm); err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(paraPebbleDB.impls); i++ {
		path := filepath.Join(root, fmt.Sprint(i))
		db, err := NewPebbleDB(path)
		if err != nil {
			return nil, err
		}
		paraPebbleDB.impls[i] = db
	}

	if shardFunc != nil {
		paraPebbleDB.shardFunc = shardFunc
	} else {
		paraPebbleDB.shardFunc = func(n int, key []byte) int {
			total := 0
			for _, b := range key {
				total += int(b)
			}
			return total % n
		}
	}
	return &paraPebbleDB, nil
}

func (this *ParaPebbleDB) Get(key []byte) ([]byte, error) {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.Get(key)
}

func (this *ParaPebbleDB) Has(key []byte) (bool, error) {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.Has(key)
}

func (this *ParaPebbleDB) Set(key []byte, value []byte) error {
	idx, db := this.getShard(key)
	this.shardLocks[idx].Lock()
	defer this.shardLocks[idx].Unlock()
	return db.Set(key, value)
}

func (this *ParaPebbleDB) Delete(key []byte) error {
	idx, db := this.getShard(key)
	this.shardLocks[idx].Lock()
	defer this.shardLocks[idx].Unlock()
	return db.Delete(key)
}

func (this *ParaPebbleDB) DeleteBatch(keys [][]byte) []error {
	categorized := make([][][]byte, len(this.impls))
	origIdx := make([][]int, len(this.impls))
	for i := range categorized {
		categorized[i] = make([][]byte, 0, len(keys)/len(this.impls)+1)
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

func (this *ParaPebbleDB) GetBatch(keys [][]byte) ([][]byte, []error) {
	categorizedKeys := make([][][]byte, len(this.impls))
	categorizedIdx := make([][]int, len(this.impls))
	for i := range categorizedKeys {
		categorizedKeys[i] = make([][]byte, 0, len(keys)/len(this.impls)+1)
		categorizedIdx[i] = make([]int, 0, len(keys)/len(this.impls)+1)
	}

	for i, key := range keys {
		idx, _ := this.getShard(key)
		categorizedKeys[idx] = append(categorizedKeys[idx], key)
		categorizedIdx[idx] = append(categorizedIdx[idx], i)
	}

	shardResults := make([][][]byte, len(this.impls))
	shardErrSlices := make([][]error, len(this.impls))
	finder := func(start, end, _ int, _ ...interface{}) {
		for i := start; i < end; i++ {
			this.shardLocks[i].RLock()
			shardResults[i], shardErrSlices[i] = this.impls[i].GetBatch(categorizedKeys[i])
			this.shardLocks[i].RUnlock()
		}
	}
	common.ParallelWorker(len(this.impls), len(this.impls), finder)

	values := make([][]byte, len(keys))
	errs := make([]error, len(keys))
	for i := range categorizedIdx {
		for j, orig := range categorizedIdx[i] {
			values[orig] = shardResults[i][j]
			errs[orig] = shardErrSlices[i][j]
		}
	}
	return values, errs
}

func (this *ParaPebbleDB) SetBatch(keys [][]byte, values [][]byte) []error {
	categorizedKeys := make([][][]byte, len(this.impls))
	categorizedVals := make([][][]byte, len(this.impls))
	origIdx := make([][]int, len(this.impls))
	for i := range categorizedKeys {
		categorizedKeys[i] = make([][]byte, 0, len(keys)/len(this.impls)+1)
		categorizedVals[i] = make([][]byte, 0, len(keys)/len(this.impls)+1)
		origIdx[i] = make([]int, 0, len(keys)/len(this.impls)+1)
	}

	for i, key := range keys {
		idx, _ := this.getShard(key)
		categorizedKeys[idx] = append(categorizedKeys[idx], key)
		categorizedVals[idx] = append(categorizedVals[idx], values[i])
		origIdx[idx] = append(origIdx[idx], i)
	}

	perShardErrs := slice.ParallelTransform(categorizedKeys, len(categorizedKeys), func(i int, _ [][]byte) []error {
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

func (this *ParaPebbleDB) Query(prefix []byte, checker func([]byte, []byte) bool) ([][]byte, [][]byte, error) {
	shardIdx, db := this.getShard(prefix)
	this.shardLocks[shardIdx].RLock()
	defer this.shardLocks[shardIdx].RUnlock()
	return db.Query(prefix, checker)
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

func (this *ParaPebbleDB) getShard(key []byte) (int, *PebbleDB) {
	shardIdx := this.shardFunc(len(this.impls), key)
	return shardIdx, this.impls[shardIdx]
}
