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

package filedb

import (
	"github.com/arcology-network/common-lib/codec"
	slice "github.com/arcology-network/common-lib/exp/slice"
)

func (this *FileDB) Query(pattern string, condition func(string, string) bool) ([]string, [][]byte, error) {
	parentPath := this.findPath(pattern) // match file parent path first
	if files, err := this.getFilesUnder(parentPath); err == nil {
		keyset := make([][]string, len(files))
		valSet := make([][][]byte, len(files))

		for i := 0; i < len(files); i++ {
			keys, valBytes, err := this.loadFile(files[i])
			if err != nil {
				return []string{}, [][]byte{}, err
			}

			for j := 0; j < len(keys); j++ {
				if !condition(pattern, keys[j]) {
					keys[j] = ""
					valBytes[j] = valBytes[j][:0]
				}
			}

			slice.Remove(&keys, "")
			slice.RemoveIf(&valBytes, func(_ int, v []byte) bool { return len(v) == 0 })

			keyset[i] = keys
			valSet[i] = valBytes
		}
		return codec.Stringset(keyset).Flatten(), codec.Bytegroup(valSet).Flatten(), nil
	} else {
		return []string{}, [][]byte{}, err
	}
}
