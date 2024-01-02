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

type Indexer struct {
	indexes map[string]interface{}
}

func NewIndexer() *Indexer {
	return &Indexer{
		indexes: map[string]interface{}{},
	}
}

func (this *MemoryDB) AddIndex(key string, v []byte) error {
	return this.db.Set(key, v)
}

func (this *MemoryDB) Set(key string, v []byte) error {
	return this.db.Set(key, v)
}

func (this *MemoryDB) Get(key string) ([]byte, error) {
	v, _ := this.db.Get(key)
	if v == nil {
		return nil, nil
	}
	return v.([]byte), nil
}
