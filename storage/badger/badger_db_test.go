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
	"testing"
)

func TestBadgerDBFunctions(t *testing.T) {
	db := NewBadgerDB(tempBadgerPath(t))
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close db: %v", err)
		}
	})

	if errs := db.SetBatch([]string{
		"a01",
		"a02",
		"a03",
		"b01",
		"c01",
		"d01",
	}, [][]byte{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
		{13, 14, 15},
		{16, 17, 18},
	}); errs != nil {
		t.Fatal(errs)
	}

	values, errs := db.GetBatch([]string{
		"a01",
		"b01",
		"c01",
	})
	if errs != nil {
		t.Fatal(errs)
	}
	if len(values) != 3 ||
		!bytes.Equal(values[0].([]byte), []byte{1, 2, 3}) ||
		!bytes.Equal(values[1].([]byte), []byte{10, 11, 12}) ||
		!bytes.Equal(values[2].([]byte), []byte{13, 14, 15}) {
		t.Error("GetBatch Failed")
	}

	value, _ := db.Get("d01")
	if !bytes.Equal(value.([]byte), []byte{16, 17, 18}) {
		t.Error("Get Failed")
	}
	if has := db.Has("d01"); !has {
		t.Error("Has Failed")
	}

	if err := db.Delete("d01"); err != nil {
		t.Fatal(err)
	}
	if has := db.Has("d01"); has {
		t.Error("Delete Failed")
	}

	if errs := db.DeleteBatch([]string{"a01", "b01"}); errs != nil {
		t.Fatal(errs)
	}
	h1 := db.Has("a01")
	h2 := db.Has("b01")
	if h1 || h2 {
		t.Error("DeleteBatch Failed")
	}

	queryKeys, queryValues, _ := db.Query("a", nil)
	t.Log(queryKeys)
	t.Log(queryValues)
}
