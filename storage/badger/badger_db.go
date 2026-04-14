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

func (db *BadgerDB) Has(key string) bool {
	err := db.impl.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	return err == nil
}

func (db *BadgerDB) Set(key string, value []byte) error {
	panic("not implemented")
}

func (db *BadgerDB) Delete(key string) error {
	return db.impl.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (db *BadgerDB) DeleteBatch(keys []string) error {
	return db.impl.Update(func(txn *badger.Txn) error {
		for _, key := range keys {
			if len(key) == 0 {
				continue
			}
			if err := txn.Delete([]byte(key)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *BadgerDB) GetBatch(keys []string) (values [][]byte, err error) {
	err = db.impl.View(func(txn *badger.Txn) error {
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

func (db *BadgerDB) SetBatch(keys []string, values [][]byte) error {
	index := 0
	for index < len(keys) {
		err := db.impl.Update(func(txn *badger.Txn) error {
			for i := index; i < len(keys); i++ {
				if len(keys[i]) == 0 {
					index++
					continue
				}

				if err := txn.Set([]byte(keys[i]), values[i]); err != nil {
					return err
				} else {
					index++
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *BadgerDB) Query(prefix string, checker func(string, []byte) bool) (keys []string, values [][]byte, err error) {
	err = db.impl.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   100,
			Prefix:         []byte(prefix),
		})
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			keys = append(keys, string(item.Key()))
			val, copyErr := item.ValueCopy(nil)
			if copyErr != nil {
				return copyErr
			}
			values = append(values, val)
		}
		return nil
	})
	return
}

func (db *BadgerDB) Close() error {
	return db.impl.Close()
}
