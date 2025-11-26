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

package badgerdb

import (
	"bytes"
	"os"
	"testing"

	common "github.com/arcology-network/common-lib/common"
)

func TestParaBadgerDBFunctions(t *testing.T) {
	os.RemoveAll(TEST_ROOT_PATH)

	db := NewParaBadgerDB("./badger-test/", common.Remainder)

	data := [][]byte{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
		{13, 14, 15},
		{16, 17, 18},
	}
	keys := []string{
		"a01",
		"a02",
		"a03",
		"b01",
		"c03",
		"d01",
	}

	db.BatchSet(keys, data)

	values, err := db.BatchGet(keys)
	if err != nil {
		t.Error(err)
	}

	values, _ = db.BatchGet([]string{
		"a01",
		"b01",
		"c03",
	})
	if len(values) != 3 ||
		!bytes.Equal(values[0], []byte{1, 2, 3}) ||
		!bytes.Equal(values[1], []byte{10, 11, 12}) ||
		!bytes.Equal(values[2], []byte{13, 14, 15}) {
		t.Error("BatchGet Failed")
	}

	value, _ := db.Get("d01")
	if !bytes.Equal(value, []byte{16, 17, 18}) {
		t.Error("Get Failed")
	}

	keys, values, _ = db.Query("a", nil)
	t.Log(keys)
	t.Log(values)
	os.RemoveAll("./badger-test/")
}
