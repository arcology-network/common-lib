/*
*   Copyright (c) 2026 Arcology Network

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
package cachedkvstore

type Readable[K comparable, V any] interface {
	Get(K) (V, bool)
	Has(K) bool
	GetBatch([]K) []V
	Len() uint64
	Size() uint64
}

type Writeable[K comparable, V any] interface {
	Set(K, V)
	Delete(K)
	SetBatch([]K, []V)
	DeleteBatch([]K)
}

type ReadWriteStore[K comparable, V any] interface {
	Readable[K, V]
	Writeable[K, V]
}

type Backend[K comparable, V any] interface {
	Readable[K, V]
	Precommit() error
	Commit(bool, uint64) error
}

type KVStore[K comparable, V any] interface {
	ReadWriteStore[K, V]
	Backend[K, V]
	SetLocalOnly(yes bool)
	LocalOnly() bool
}
