package cachedstorage

import "reflect"

type AccessibleInterface interface {
	Value() interface{}
	Reads() uint32
	Writes() uint32
	Size() uint32
}

const (
	MEMORY_DB     = 0
	PERSISTENT_DB = 1
)

type PersistentStorageInterface interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	BatchGet([]string) ([][]byte, error)
	BatchSet([]string, [][]byte) error
	Query(string, func(string, string) bool) ([]string, [][]byte, error)
}

type DbFilter func(PersistentStorageInterface) bool

func NotQueryRpc(db PersistentStorageInterface) bool { // Do not access MemDB
	name := reflect.TypeOf(db).String()
	return name == "*storage.ReadonlyRpcClient"
}
