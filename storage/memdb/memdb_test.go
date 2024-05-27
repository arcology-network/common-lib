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
	"bytes"
	"testing"
)

func TestMemDB(t *testing.T) {
	memDB := NewMemoryDB()
	keys := []string{"123", "456"}
	values := make([][]byte, 2)
	values[0] = []byte{1, 2, 3}
	values[1] = []byte{4, 5, 6}
	memDB.BatchSet(keys, values)

	if v, _ := memDB.Get(keys[0]); !bytes.Equal(v, values[0]) {
		t.Error("Error")
	}

	if v, _ := memDB.Get(keys[1]); !bytes.Equal(v, values[1]) {
		t.Error("Error")
	}

	retrived, _ := memDB.BatchGet(append(keys, ""))
	if len(values) != 2 || !bytes.Equal(values[0], retrived[0]) || !bytes.Equal(values[1], retrived[1]) {
		t.Error("Error")
	}
}
