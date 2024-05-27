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

package interfaces

import "reflect"

type Accessible interface {
	Value() interface{}
	Reads() uint32
	Writes() uint32
	Size() uint32
}

const (
	MEMORY_DB     = 0
	PERSISTENT_DB = 1
)

type PersistentStorage interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	BatchGet([]string) ([][]byte, error)
	BatchSet([]string, [][]byte) error
	Query(string, func(string, string) bool) ([]string, [][]byte, error)
}

type DbFilter func(PersistentStorage) bool

func NotQueryRpc(db PersistentStorage) bool { // Do not access MemoryDB
	name := reflect.TypeOf(db).String()
	return name == "*storage.ReadonlyRpcClient"
}
