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
	"testing"
)

func TestGetSet(t *testing.T) {
	db := NewSimpleFileDB("./testdb/")
	key := "key"
	value := []byte{1, 2, 3, 4, 5}

	db.Set(key, value)
	if v, err := db.Get(key); err != nil {
		t.Error("error", err)
	} else if !bytes.Equal(v, value) {
		t.Error("v", v)
	}
}
