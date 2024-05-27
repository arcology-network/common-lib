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
