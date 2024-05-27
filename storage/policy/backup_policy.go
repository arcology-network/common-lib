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

package cachepolicy

type BackupPolicy[K comparable, V any] struct {
	datastore interface {
		Cache(interface{}) interface{}
	}
	interval uint32
}

func NewBackupPolicy[K comparable, V any](datastore interface{ Cache(interface{}) interface{} }, interval uint32) *BackupPolicy[K, V] {
	return &BackupPolicy[K, V]{
		datastore: datastore,
		interval:  interval,
	}
}

func (this *BackupPolicy[K, V]) FullBackup() {
	// keys, values := this.datastore.Cache(nil).(*expmap.ConcurrentMap[string, any]).KVs()
	// codec.Strings(keys).Encode()

	// encoder := this.datastore.Encoder(nil)
	// byteset := make([][]byte, len(keys))
	// for i := 0; i < len(keys); i++ {
	// 	byteset[i] = encoder(keys[i], values[i])
	// }
	// codec.Strings(keys).Encode()
}
