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
	"bytes"

	"github.com/cockroachdb/pebble"
)

type PebbleDB struct {
	impl *pebble.DB
}

func NewPebbleDB(path string) (*PebbleDB, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &PebbleDB{impl: db}, nil
}

func (this *PebbleDB) Get(key []byte) ([]byte, error) {
	stored, closer, err := this.impl.Get(key)
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	return bytes.Clone(stored), nil
}

func (this *PebbleDB) GetBatch(keys [][]byte) ([][]byte, []error) {
	values := make([][]byte, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		stored, closer, err := this.impl.Get(key)
		if err != nil {
			errs[i] = err
			continue
		}
		values[i] = bytes.Clone(stored)
		errs[i] = closer.Close()
	}
	return values, errs
}

func (this *PebbleDB) Set(key []byte, value []byte) error {
	return this.impl.Set(key, value, pebble.NoSync)
}

func (this *PebbleDB) SetBatch(keys [][]byte, values [][]byte) []error {
	batch := this.impl.NewBatch()
	defer batch.Close()

	errs := make([]error, len(keys))
	for i := range keys {
		if err := batch.Set(keys[i], values[i], pebble.NoSync); err != nil {
			errs[i] = err
		}
	}
	if err := batch.Commit(pebble.NoSync); err != nil {
		for i := range errs {
			if errs[i] == nil {
				errs[i] = err
			}
		}
	}
	return errs
}

func (this *PebbleDB) Delete(key []byte) error {
	return this.impl.Delete(key, pebble.NoSync)
}

func (this *PebbleDB) DeleteBatch(keys [][]byte) []error {
	batch := this.impl.NewBatch()
	defer batch.Close()

	errs := make([]error, len(keys))
	for i, key := range keys {
		if err := batch.Delete(key, pebble.NoSync); err != nil {
			errs[i] = err
		}
	}
	if err := batch.Commit(pebble.NoSync); err != nil {
		for i := range errs {
			if errs[i] == nil {
				errs[i] = err
			}
		}
	}
	return errs
}

func (this *PebbleDB) Has(key []byte) (bool, error) {
	_, closer, err := this.impl.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, closer.Close()
}

func (this *PebbleDB) Query(prefix []byte, checker func([]byte, []byte) bool) ([][]byte, [][]byte, error) {
	iter, err := this.impl.NewIter(&pebble.IterOptions{})
	if err != nil {
		return nil, nil, err
	}
	defer iter.Close()

	var keys [][]byte
	var values [][]byte
	for iter.First(); iter.Valid(); iter.Next() {
		if len(prefix) > 0 && !bytes.HasPrefix(iter.Key(), prefix) {
			continue
		}
		k := bytes.Clone(iter.Key())
		v := bytes.Clone(iter.Value())
		if checker != nil && !checker(k, v) {
			continue
		}
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values, iter.Error()
}

func (this *PebbleDB) Close() error {
	return this.impl.Close()
}
