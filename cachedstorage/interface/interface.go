package interfaces

import "reflect"

type Accessible interface {
	Value() interface{}
	Reads() uint32
	Writes() uint32
	Size() uint32
}

const (
	MEMORY_DB     = 0
	PERSISTENT_DB = 1
)

type PersistentStorage interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	BatchGet([]string) ([][]byte, error)
	BatchSet([]string, [][]byte) error
	Query(string, func(string, string) bool) ([]string, [][]byte, error)
}

type DbFilter func(PersistentStorage) bool

func NotQueryRpc(db PersistentStorage) bool { // Do not access MemDB
	name := reflect.TypeOf(db).String()
	return name == "*storage.ReadonlyRpcClient"
}
