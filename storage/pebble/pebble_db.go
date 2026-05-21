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
	"unsafe"

	stgintf "github.com/arcology-network/common-lib/storage/interface"
	"github.com/cockroachdb/pebble"
)

type PebbleDB struct {
	impl    *pebble.DB
	decoder func(string, any, any) (any, error)
}

func NewPebbleDB(path string, decoder ...func(string, any, any) (any, error)) (*PebbleDB, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}

	var decode func(string, any, any) (any, error)
	if len(decoder) > 0 {
		decode = decoder[0]
	}
	return &PebbleDB{impl: db, decoder: decode}, nil
}

func (this *PebbleDB) Get(key string) (any, error) {
	return this.GetAs(key, nil)
}

func (this *PebbleDB) GetAs(key string, typeHint any) (any, error) {
	if this == nil {
		return nil, stgintf.ErrNotFound
	}

	stored, closer, err := this.impl.Get(unsafe.Slice(unsafe.StringData(key), len(key)))
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, stgintf.ErrNotFound
		}
		return nil, err
	}
	defer closer.Close()
	value := bytes.Clone(stored)
	if this.decoder != nil {
		return this.decoder(key, value, typeHint)
	}
	return value, nil
}

func (this *PebbleDB) GetBatch(keys []string) ([]any, []error) {
	values := make([]any, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		stored, closer, err := this.impl.Get(unsafe.Slice(unsafe.StringData(key), len(key)))
		if err != nil {
			errs[i] = err
			continue
		}
		values[i] = bytes.Clone(stored)
		errs[i] = closer.Close()
	}
	return values, errs
}

func (this *PebbleDB) Set(key string, value []byte) error {
	return this.impl.Set(unsafe.Slice(unsafe.StringData(key), len(key)), value, pebble.NoSync)
}

func (this *PebbleDB) SetBatch(keys []string, values [][]byte) []error {
	batch := this.impl.NewBatch()
	defer batch.Close()

	errs := make([]error, len(keys))
	for i := range keys {
		if err := batch.Set(unsafe.Slice(unsafe.StringData(keys[i]), len(keys[i])), values[i], pebble.NoSync); err != nil {
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

func (this *PebbleDB) Delete(key string) error {
	return this.impl.Delete(unsafe.Slice(unsafe.StringData(key), len(key)), pebble.NoSync)
}

func (this *PebbleDB) DeleteBatch(keys []string) []error {
	batch := this.impl.NewBatch()
	defer batch.Close()

	errs := make([]error, len(keys))
	for i, key := range keys {
		if err := batch.Delete(unsafe.Slice(unsafe.StringData(key), len(key)), pebble.NoSync); err != nil {
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

func (this *PebbleDB) Has(key string) bool {
	_, closer, err := this.impl.Get(unsafe.Slice(unsafe.StringData(key), len(key)))
	if err != nil {
		if err == pebble.ErrNotFound {
			return false
		}
		return false
	}
	defer closer.Close()
	return true
}

func (this *PebbleDB) Query(prefix string, checker func(string, []byte) bool) ([]string, [][]byte, error) {
	iter, err := this.impl.NewIter(&pebble.IterOptions{})
	if err != nil {
		return nil, nil, err
	}
	defer iter.Close()

	var keys []string
	var values [][]byte
	prefixBytes := unsafe.Slice(unsafe.StringData(prefix), len(prefix))
	for iter.First(); iter.Valid(); iter.Next() {
		if len(prefixBytes) > 0 && !bytes.HasPrefix(iter.Key(), prefixBytes) {
			continue
		}
		k := string(iter.Key())
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
