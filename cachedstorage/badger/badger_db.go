package badgerdb

import (
	"github.com/dgraph-io/badger"
)

type BadgerDB struct {
	impl *badger.DB
}

func NewBadgerDB(path string) *BadgerDB {
	bdg, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		panic(err)
	}
	return &BadgerDB{
		impl: bdg,
	}
}

func (db *BadgerDB) Get(key string) (value []byte, err error) {
	err = db.impl.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == nil {
			value, err = item.ValueCopy(nil)
		}
		return err
	})
	return
}

func (db *BadgerDB) Set(key string, value []byte) error {
	panic("not implemented")
}

func (db *BadgerDB) BatchGet(keys []string) (values [][]byte, err error) {
	db.impl.View(func(txn *badger.Txn) error {
		for i := range keys {
			if len(keys[i]) == 0 {
				continue
			}
			item, err := txn.Get([]byte(keys[i]))
			if err != nil {
				values = append(values, nil)
			} else {
				val, err := item.ValueCopy(nil)
				if err != nil {
					values = append(values, nil)
				} else {
					values = append(values, val)
				}
			}
		}
		return nil
	})
	return
}

func (db *BadgerDB) BatchSet(keys []string, values [][]byte) error {
	index := 0
	for index < len(keys) {
		db.impl.Update(func(txn *badger.Txn) error {
			for i := index; i < len(keys); i++ {
				if len(keys[i]) == 0 {
					continue
				}

				if err := txn.Set([]byte(keys[i]), values[i]); err != nil {
					return nil
				} else {
					index++
				}
			}
			return nil
		})
	}
	return nil
}

func (db *BadgerDB) Query(prefix string, checker func(string, string) bool) (keys []string, values [][]byte, err error) {
	db.impl.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   100,
			Prefix:         []byte(prefix),
		})
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			keys = append(keys, string(item.Key()))
			val, _ := item.ValueCopy(nil)
			values = append(values, val)
		}
		return nil
	})
	return
}

func (db *BadgerDB) Close() error {
	return db.impl.Close()
}
