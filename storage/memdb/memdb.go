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
)

type MemoryDB struct {
	// db *ccmap.ConcurrentMap
	db *ccmap.ConcurrentMap[string, any]
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		db: ccmap.NewConcurrentMap(
			16,
			func(v any) bool { return v == nil },
			func(k string) uint64 {
				var hash uint64
				for i := 0; i < len(k); i++ {
					hash += uint64(k[i])
				}
				return hash % 16
			},
		),
	}
}

func (this *MemoryDB) Set(key string, v []byte) error {
	this.db.Set(key, v)
	return nil
}

func (this *MemoryDB) Get(key string) ([]byte, error) {
	v, _ := this.db.Get(key)
	if v == nil {
		return nil, nil
	}
	return v.([]byte), nil
}

func (this *MemoryDB) BatchGet(keys []string) ([][]byte, error) {
	values, _ := this.db.BatchGet(keys)
	byteset := make([][]byte, len(keys))
	for i, v := range values {
		if v != nil {
			byteset[i] = v.([]byte)
		}
	}
	return byteset, nil
}

func (this *MemoryDB) BatchSet(keys []string, byteset [][]byte) error {
	values := make([]interface{}, len(keys))
	for i, v := range byteset {
		if v != nil {
			values[i] = v
		}
	}

	this.db.BatchSet(keys, values)
	return nil
}

func (this *MemoryDB) Query(key string, functor func(string, string) bool) ([]string, [][]byte, error) {
	return []string{}, [][]byte{}, nil
}
