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

type ParaPebbleDB[K Key, T any] struct {
	impls      [16]*PebbleDB[K, T]
	shardLocks [16]sync.RWMutex
	shardFunc  func(int, K) int
}

func NewParaPebbleDB[K Key, T any](root string, shardFunc func(numOfShard int, key K) int) *ParaPebbleDB[K, T] {
	return NewParaPebbleDBWithCodec[K, T](root, shardFunc, nil, nil)
}

func NewParaPebbleDBWithCodec[K Key, T any](
	root string,
	shardFunc func(numOfShard int, key K) int,
	encoder func(K, T) ([]byte, error),
	decoder func(K, []byte, T) any,
) *ParaPebbleDB[K, T] {
	var paraPebbleDB ParaPebbleDB[K, T]
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err := os.MkdirAll(root, fs.ModePerm); err != nil {
			panic(err)
		}
	}

	for i := 0; i < len(paraPebbleDB.impls); i++ {
		path := filepath.Join(root, fmt.Sprint(i))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.MkdirAll(path, fs.ModePerm); err != nil {
				panic(err)
			}
		}
		paraPebbleDB.impls[i] = NewPebbleDBWithCodec[K, T](path, encoder, decoder)
	}

	if shardFunc != nil {
		paraPebbleDB.shardFunc = shardFunc
	} else {
		paraPebbleDB.shardFunc = paraPebbleDB.hash32
	}
	return &paraPebbleDB
}

func (this *ParaPebbleDB[K, T]) Get(key K) (value T, err error) {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.Get(key)
}

func (this *ParaPebbleDB[K, T]) Has(key K) bool {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.Has(key)
}

func (this *ParaPebbleDB[K, T]) GetAs(key K, target T) (any, error) {
	idx, db := this.getShard(key)
	this.shardLocks[idx].RLock()
	defer this.shardLocks[idx].RUnlock()
	return db.GetAs(key, target)
}

func (this *ParaPebbleDB[K, T]) Set(key K, value T) error {
	idx, db := this.getShard(key)
	this.shardLocks[idx].Lock()
	defer this.shardLocks[idx].Unlock()
	return db.Set(key, value)
}

func (this *ParaPebbleDB[K, T]) Delete(key K) error {
	idx, db := this.getShard(key)
	this.shardLocks[idx].Lock()
	defer this.shardLocks[idx].Unlock()
	return db.Delete(key)
}

func (this *ParaPebbleDB[K, T]) DeleteBatch(keys []K) error {
	categorized := make([][]K, len(this.impls))
	for i := 0; i < len(categorized); i++ {
		categorized[i] = make([]K, 0, len(keys)/len(this.impls)+100)
	}

	for i := 0; i < len(keys); i++ {
		idx, _ := this.getShard(keys[i])
		categorized[idx] = append(categorized[idx], keys[i])
	}

	for i := range categorized {
		if len(categorized[i]) == 0 {
			continue
		}
		this.shardLocks[i].Lock()
		err := this.impls[i].DeleteBatch(categorized[i])
		this.shardLocks[i].Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *ParaPebbleDB[K, T]) GetBatch(keys []K) (values []T, err error) {
	categorizedKeys := make([][]K, len(this.impls))
	categorizedIdx := make([][]int, len(this.impls))
	for i := 0; i < len(categorizedKeys); i++ {
		categorizedKeys[i] = make([]K, 0, len(keys)/len(this.impls)+100)
		categorizedIdx[i] = make([]int, 0, len(keys)/len(this.impls)+100)
	}

	for i := 0; i < len(keys); i++ {
		idx, _ := this.getShard(keys[i])
		categorizedKeys[idx] = append(categorizedKeys[idx], keys[i])
		categorizedIdx[idx] = append(categorizedIdx[idx], i)
	}

	errors := make([]error, len(categorizedKeys))
	valueSet := make([][]T, len(categorizedKeys))
	finder := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			this.shardLocks[i].RLock()
			valueSet[i], errors[i] = this.impls[i].GetBatch(categorizedKeys[i])
			this.shardLocks[i].RUnlock()
		}
	}
	common.ParallelWorker(len(categorizedKeys), len(categorizedKeys), finder)

	results := make([]T, len(keys))
	for i := range categorizedIdx {
		for j, original := range categorizedIdx[i] {
			results[original] = valueSet[i][j]
		}
	}
	return results, errors[0]
}

func (this *ParaPebbleDB[K, T]) SetBatch(keys []K, values []T) error {
	categorizedKeys := make([][]K, len(this.impls))
	categorizedVals := make([][]T, len(this.impls))
	for i := 0; i < len(categorizedKeys); i++ {
		categorizedKeys[i] = make([]K, 0, len(keys)/len(this.impls)+100)
		categorizedVals[i] = make([]T, 0, len(keys)/len(this.impls)+100)
	}

	for i := 0; i < len(keys); i++ {
		idx, _ := this.getShard(keys[i])
		categorizedKeys[idx] = append(categorizedKeys[idx], keys[i])
		categorizedVals[idx] = append(categorizedVals[idx], values[i])
	}

	errors := slice.ParallelTransform(categorizedKeys, len(categorizedKeys), func(i int, _ []K) error {
		this.shardLocks[i].Lock()
		defer this.shardLocks[i].Unlock()
		return this.impls[i].SetBatch(categorizedKeys[i], categorizedVals[i])
	})
	return errors[0]
}

func (this *ParaPebbleDB[K, T]) Query(prefix K, checker func(K, T) bool) (keys []K, values []T, err error) {
	shardIdx, db := this.getShard(prefix)
	this.shardLocks[shardIdx].RLock()
	defer this.shardLocks[shardIdx].RUnlock()
	return db.Query(prefix, checker)
}

func (this *ParaPebbleDB[K, T]) Close() error {
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

func (this *ParaPebbleDB[K, T]) getShard(key K) (int, *PebbleDB[K, T]) {
	shardIdx := this.shardFunc(len(this.impls), key)
	return shardIdx, this.impls[shardIdx]
}

func (this *ParaPebbleDB[K, T]) hash32(numOfShard int, key K) int {
	encoded, err := defaultKeyEncoder[K](key)
	if err != nil || len(encoded) == 0 {
		return 0
	}

	total := 0
	for i := 0; i < len(encoded); i++ {
		total += int(encoded[i])
	}
	return total % numOfShard
}
