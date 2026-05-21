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
	if errs := memDB.SetBatch(keys, values); len(errs) != len(keys) {
		t.Fatal("unexpected batch error slice length")
	}
	h0 := memDB.Has(keys[0])
	h1 := memDB.Has(keys[1])
	if !h0 || !h1 {
		t.Error("Error")
	}

	v0, err := memDB.Get(keys[0])
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(v0.([]byte), values[0]) {
		t.Error("Error")
	}

	v1, err := memDB.Get(keys[1])
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(v1.([]byte), values[1]) {
		t.Error("Error")
	}

	retrieved, errs := memDB.GetBatch(append(keys, ""))
	if len(errs) != 3 || errs[2] == nil {
		t.Fatal("expected missing key error")
	}
	if len(values) != 2 || !bytes.Equal(values[0], retrieved[0].([]byte)) || !bytes.Equal(values[1], retrieved[1].([]byte)) {
		t.Error("Error")
	}

	if err := memDB.Delete(keys[0]); err != nil {
		t.Fatal(err)
	}
	h0 = memDB.Has(keys[0])
	if h0 {
		t.Error("Error")
	}

	if errs := memDB.DeleteBatch([]string{keys[1]}); len(errs) != 1 {
		t.Fatal("unexpected delete batch error slice length")
	}
	h1 = memDB.Has(keys[1])
	if h1 {
		t.Error("Error")
	}
}

func TestMemDBGetDelegatesToGetAsNilDecoder(t *testing.T) {
	memDB := NewMemoryDB()
	expected := []byte{7, 8, 9}
	if err := memDB.Set("alpha", expected); err != nil {
		t.Fatal(err)
	}

	got, err := memDB.Get("alpha")
	if err != nil {
		t.Fatal(err)
	}

	viaGetAs, err := memDB.GetAs("alpha", nil)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got.([]byte), viaGetAs.([]byte)) {
		t.Fatalf("expected Get to match GetAs with nil decoder")
	}
}
