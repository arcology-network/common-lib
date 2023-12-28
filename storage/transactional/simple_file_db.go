package transactional

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/natefinch/atomic"
)

type SimpleFileDB struct {
	root string
}

func NewSimpleFileDB(root string) *SimpleFileDB {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		os.Mkdir(root, 0755)
	}

	return &SimpleFileDB{
		root: root,
	}
}

func (db *SimpleFileDB) Set(key string, value []byte) error {
	return atomic.WriteFile(db.root+key, bytes.NewReader(value))
}

func (db *SimpleFileDB) Get(key string) ([]byte, error) {
	//fmt.Printf("[SimpleFileDB.Get] root = %s, key = %s\n", db.root, key)
	return ioutil.ReadFile(db.root + key)
}

func (db *SimpleFileDB) Delete(key string) error {
	return os.Remove(db.root + key)
}
