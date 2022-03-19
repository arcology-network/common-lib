package cachedstorage

type AccessableInterface interface {
	Value() interface{}
	Reads() uint32
	Writes() uint32
}

type MeasurableInterface interface {
	Size() uint32
}

type PersistentStorageInterface interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	BatchGet([]string) ([][]byte, error)
	BatchSet([]string, [][]byte) error
}
