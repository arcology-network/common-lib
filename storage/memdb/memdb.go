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

package memdb

import (
	ccmap "github.com/arcology-network/common-lib/exp/map"
	stgintf "github.com/arcology-network/common-lib/storage/interface"
)

var _ stgintf.ReadWriteStore[string, []byte] = (*MemoryDB)(nil)

type MemoryDB struct {
	db      *ccmap.ConcurrentMap[string, []byte]
	decoder func(string, any, any) (any, error)
}

func NewMemoryDB(decoder ...func(string, any, any) (any, error)) *MemoryDB {
	var decode func(string, any, any) (any, error)
	if len(decoder) > 0 {
		decode = decoder[0]
	}

	return &MemoryDB{
		db: ccmap.NewConcurrentMap(
			16,
			func(v []byte) bool { return v == nil },
			func(k string) uint64 {
				var hash uint64
				for i := 0; i < len(k); i++ {
					hash += uint64(k[i])
				}
				return hash % 16
			},
		),
		decoder: decode,
	}
}

func (this *MemoryDB) Set(key string, v []byte) error {
	this.db.Set(key, v)
	return nil
}

func (this *MemoryDB) Get(key string) (any, error) {
	return this.GetAs(key, nil)
}

func (this *MemoryDB) GetAs(key string, typeHint any) (any, error) {
	if this == nil {
		return nil, stgintf.ErrNotFound
	}

	v, ok := this.db.Get(key)
	if !ok || v == nil {
		return nil, stgintf.ErrNotFound
	}

	if this.decoder != nil {
		return this.decoder(key, v, typeHint)
	}
	return v, nil
}

func (this *MemoryDB) Has(key string) bool {
	v, ok := this.db.Get(key)
	return ok && v != nil
}

func (this *MemoryDB) GetBatch(keys []string) ([]any, []error) {
	values, oks := this.db.GetBatch(keys)
	byteset := make([]any, len(keys))
	errs := make([]error, len(keys))
	for i, v := range values {
		if oks[i] && v != nil {
			byteset[i] = v
			continue
		}
		errs[i] = stgintf.ErrNotFound
	}
	return byteset, errs
}

func (this *MemoryDB) SetBatch(keys []string, byteset [][]byte) []error {
	this.db.SetBatch(keys, byteset)
	return make([]error, len(keys))
}

func (this *MemoryDB) Delete(key string) error {
	this.db.Set(key, nil)
	return nil
}

func (this *MemoryDB) DeleteBatch(keys []string) []error {
	values := make([][]byte, len(keys))
	this.db.SetBatch(keys, values)
	return make([]error, len(keys))
}

func (this *MemoryDB) Query(key string, functor func(string, []byte) bool) ([]string, [][]byte, []error) {
	keys, values := this.db.KVs()
	matchedKeys := make([]string, 0, len(keys))
	matchedValues := make([][]byte, 0, len(keys))

	for i, storedKey := range keys {
		storedValue := values[i]
		if storedValue == nil {
			continue
		}
		if functor != nil {
			if !functor(storedKey, storedValue) {
				continue
			}
		} else if storedKey != key {
			continue
		}
		matchedKeys = append(matchedKeys, storedKey)
		matchedValues = append(matchedValues, storedValue)
	}
	return matchedKeys, matchedValues, nil
}
