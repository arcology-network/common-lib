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

package pebbledb

import (
	"testing"
)

func TestParaPebbleDBFunctions(t *testing.T) {
	db, err := NewParaPebbleDB(tempParaPebbleRoot(t), nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close db: %v", err)
		}
	})

	keys := []string{
		"a01", "a02", "a03",
		"b01", "c03", "d01",
	}
	values := [][]byte{
		{1}, {4}, {7}, {10}, {13}, {16},
	}

	errs := db.SetBatch(keys, values)
	for i, err := range errs {
		if err != nil {
			t.Fatalf("SetBatch[%d]: %v", i, err)
		}
	}

	got, getErrs := db.GetBatch([]string{"a01", "b01", "c03"})
	for i, err := range getErrs {
		if err != nil {
			t.Fatalf("GetBatch[%d]: %v", i, err)
		}
	}
	got0, ok0 := got[0].([]byte)
	got1, ok1 := got[1].([]byte)
	got2, ok2 := got[2].([]byte)
	if len(got) != 3 || !ok0 || !ok1 || !ok2 || got0[0] != 1 || got1[0] != 10 || got2[0] != 13 {
		t.Error("GetBatch failed")
	}

	val, err := db.Get("d01")
	if err != nil {
		t.Fatal(err)
	}
	bytesVal, ok := val.([]byte)
	if !ok || bytesVal[0] != 16 {
		t.Error("Get failed")
	}

	exists := db.Has("d01")
	if !exists {
		t.Error("Has failed")
	}

	if err := db.Delete("d01"); err != nil {
		t.Fatal(err)
	}
	exists = db.Has("d01")
	if exists {
		t.Error("Delete failed")
	}

	delErrs := db.DeleteBatch([]string{"a01", "b01"})
	for i, err := range delErrs {
		if err != nil {
			t.Fatalf("DeleteBatch[%d]: %v", i, err)
		}
	}
	exists1 := db.Has("a01")
	exists2 := db.Has("b01")
	if exists1 || exists2 {
		t.Error("DeleteBatch failed")
	}

	qkeys, qvalues, err := db.Query("a", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(qkeys)
	t.Log(qvalues)
}
