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
	stgintf "github.com/arcology-network/common-lib/storage/interface"
	"github.com/dgraph-io/badger"
)

var _ stgintf.ReadWriteStore[string, []byte] = (*BadgerDB)(nil)

type BadgerDB struct {
	impl    *badger.DB
	decoder func(string, any, any) (any, error)
}

func NewBadgerDB(path string, decoder ...func(string, any, any) (any, error)) *BadgerDB {
	bdg, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		panic(err)
	}

	var decode func(string, any, any) (any, error)
	if len(decoder) > 0 {
		decode = decoder[0]
	}

	return &BadgerDB{
		impl:    bdg,
		decoder: decode,
	}
}

func (db *BadgerDB) Get(key string) (value any, err error) {
	return db.GetAs(key, nil)
}

func (db *BadgerDB) GetAs(key string, typeHint any) (any, error) {
	if db == nil {
		return nil, stgintf.ErrNotFound
	}

	var value []byte
	err := db.impl.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == nil {
			value, err = item.ValueCopy(nil)
		}
		return err
	})

	if err == badger.ErrKeyNotFound {
		return nil, stgintf.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if db.decoder != nil {
		return db.decoder(key, value, typeHint)
	}
	return value, nil
}

func (db *BadgerDB) Has(key string) bool {
	_, err := db.Get(key)
	return err == nil
}

func (db *BadgerDB) Set(key string, value []byte) error {
	return db.impl.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func (db *BadgerDB) Delete(key string) error {
	return db.impl.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (db *BadgerDB) DeleteBatch(keys []string) []error {
	errs := make([]error, len(keys))
	for i, key := range keys {
		if len(key) == 0 {
			errs[i] = stgintf.ErrNotFound
			continue
		}
		if err := db.Delete(key); err != nil && err != badger.ErrKeyNotFound {
			errs[i] = err
		}
	}
	if allNil(errs) {
		return nil
	}
	return errs
}

func (db *BadgerDB) GetBatch(keys []string) (values []any, errs []error) {
	values = make([]any, len(keys))
	errs = make([]error, len(keys))
	for i := range keys {
		if len(keys[i]) == 0 {
			errs[i] = stgintf.ErrNotFound
			continue
		}
		v, err := db.Get(keys[i])
		if err != nil {
			errs[i] = err
			continue
		}
		values[i] = v
	}
	if allNil(errs) {
		return values, nil
	}
	return values, errs
}

func (db *BadgerDB) SetBatch(keys []string, values [][]byte) []error {
	errs := make([]error, len(keys))
	for i := range keys {
		if len(keys[i]) == 0 {
			errs[i] = stgintf.ErrNotFound
			continue
		}
		if err := db.Set(keys[i], values[i]); err != nil {
			errs[i] = err
		}
	}
	if allNil(errs) {
		return nil
	}
	return errs
}

func allNil(errs []error) bool {
	for _, err := range errs {
		if err != nil {
			return false
		}
	}
	return true
}

func (db *BadgerDB) Query(prefix string, checker func(string, []byte) bool) (keys []string, values [][]byte, errs []error) {
	err := db.impl.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   100,
			Prefix:         []byte(prefix),
		})
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			val, copyErr := item.ValueCopy(nil)
			if copyErr != nil {
				return copyErr
			}
			if checker != nil && !checker(key, val) {
				continue
			}
			keys = append(keys, key)
			values = append(values, val)
		}
		return nil
	})
	if err != nil {
		return nil, nil, []error{err}
	}
	return keys, values, nil
}

func (db *BadgerDB) Close() error {
	return db.impl.Close()
}
