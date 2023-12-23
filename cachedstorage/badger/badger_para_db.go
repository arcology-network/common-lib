package badgerdb

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"sync"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
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
		this.shardLocks[start].RLock()
		defer this.shardLocks[start].RUnlock() // Using start is correct, as start + 1 == end

		valueSet[start], errors[start] = this.impls[start].BatchGet(categorized[start])
	}
	common.ParallelWorker(len(categorized), len(categorized), finder)
	return codec.Bytegroup(valueSet).Flatten(), errors[0]
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

	errors := make([]error, len(categorizedKeys))
	finder := func(start, end, index int, args ...interface{}) {
		this.shardLocks[start].Lock()
		defer this.shardLocks[start].Unlock() // Using start is correct, as start + 1 == end
		errors[start] = this.impls[start].BatchSet(categorizedKeys[start], categorizedVals[start])
	}
	common.ParallelWorker(len(categorizedKeys), len(categorizedKeys), finder)
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
