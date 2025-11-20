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
	slice "github.com/arcology-network/common-lib/exp/slice"
)

type ParaBadgerDB struct {
	impls      [16]*BadgerDB
	shardLocks [16]sync.RWMutex
	shardFunc  func(int, string) int
}

func NewParaBadgerDB(root string, shardFunc func(numOfShard int, key string) int) *ParaBadgerDB {
	var paraBadgerDB ParaBadgerDB
	if _, err := os.Stat(root); os.IsNotExist(err) {
		os.MkdirAll(root, fs.ModePerm)
	}

	for i := 0; i < len(paraBadgerDB.impls); i++ {
		path := path.Join(root+fmt.Sprint(i)) + "/"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, fs.ModePerm)
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

func (this *ParaBadgerDB) Get(key string) (value []byte, err error) {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.Get(key)
}

func (this *ParaBadgerDB) Set(key string, value []byte) error {
	panic("not implemented")
}

func (this *ParaBadgerDB) BatchGet(keys []string) (values [][]byte, err error) {
	categorized := make([][]string, len(this.impls))
	for i := 0; i < len(categorized); i++ {
		categorized[i] = make([]string, 0, len(keys)/len(this.impls)+100)
	}

	for i := 0; i < len(keys); i++ {
		idx, _ := this.getShard(keys[i])
		categorized[idx] = append(categorized[idx], keys[i])
	}

	errors := make([]error, len(categorized))
	valueSet := make([][][]byte, len(categorized))
	finder := func(start, end, index int, args ...interface{}) {
		for index = start; index < end; index++ {
			this.shardLocks[index].RLock()
			defer this.shardLocks[index].RUnlock() // Using start is correct, as start + 1 == end
			valueSet[index], errors[index] = this.impls[index].BatchGet(categorized[index])
		}
	}
	common.ParallelWorker(len(categorized), len(categorized), finder)

	mp := map[string][]byte{}
	for i := range categorized {
		for k := range categorized[i] {
			mp[categorized[i][k]] = valueSet[i][k]
		}
	}
	results := make([][]byte, len(keys))
	for i := range keys {
		results[i] = mp[keys[i]]
	}

	return results, errors[0]
}

func (this *ParaBadgerDB) BatchSet(keys []string, values [][]byte) error {
	categorizedKeys := make([][]string, len(this.impls))
	categorizedVals := make([][][]byte, len(this.impls))
	for i := 0; i < len(categorizedKeys); i++ {
		categorizedKeys[i] = make([]string, 0, len(keys)/len(this.impls)+100)
		categorizedVals[i] = make([][]byte, 0, len(keys)/len(this.impls)+100)
	}

	for i := 0; i < len(keys); i++ {
		idx, _ := this.getShard(keys[i])
		categorizedKeys[idx] = append(categorizedKeys[idx], keys[i])
		categorizedVals[idx] = append(categorizedVals[idx], values[i])
	}

	errors := slice.ParallelTransform(categorizedKeys, len(categorizedKeys), func(i int, _ []string) error {
		this.shardLocks[i].Lock()
		defer this.shardLocks[i].Unlock() // Using start is correct, as start + 1 == end
		return this.impls[i].BatchSet(categorizedKeys[i], categorizedVals[i])
	})
	return errors[0]
}

func (this *ParaBadgerDB) Query(prefix string, checker func(string, string) bool) (keys []string, values [][]byte, err error) {
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
		return math.MaxUint32
	}

	var total int = 0
	for j := 0; j < len(key); j++ {
		total += int(key[j])
	}
	return total % numOfShard
}
