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
	"strconv"
	"testing"
)

func TestPebbleDBFunctions(t *testing.T) {
	db := NewPebbleDB[string, uint64](tempPebblePath(t))
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close db: %v", err)
		}
	})

	if err := db.SetBatch([]string{
		"a01",
		"a02",
		"a03",
		"b01",
		"c01",
		"d01",
	}, []uint64{
		1,
		4,
		7,
		10,
		13,
		16,
	}); err != nil {
		t.Fatal(err)
	}

	values, err := db.GetBatch([]string{"a01", "b01", "c01"})
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 3 || values[0] != 1 || values[1] != 10 || values[2] != 13 {
		t.Error("GetBatch Failed")
	}

	value, err := db.Get("d01")
	if err != nil {
		t.Fatal(err)
	}
	if value != 16 {
		t.Error("Get Failed")
	}
	if !db.Has("d01") {
		t.Error("Has Failed")
	}

	if err := db.Delete("d01"); err != nil {
		t.Fatal(err)
	}
	if db.Has("d01") {
		t.Error("Delete Failed")
	}

	if err := db.DeleteBatch([]string{"a01", "b01"}); err != nil {
		t.Fatal(err)
	}
	if db.Has("a01") || db.Has("b01") {
		t.Error("DeleteBatch Failed")
	}

	keys, queried, err := db.Query("a", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 || len(queried) != 2 {
		t.Fatalf("unexpected query result: %v %v", keys, queried)
	}
	t.Log(keys)
	t.Log(queried)
}

func TestPebbleDBCodecHooks(t *testing.T) {
	encoder := func(_ int, value int64) ([]byte, error) {
		return []byte(strconv.FormatInt(value, 10)), nil
	}
	decoder := func(_ int, raw []byte, _ int64) any {
		decoded, err := strconv.ParseInt(string(raw), 10, 64)
		if err != nil {
			return nil
		}
		return decoded
	}

	db := NewPebbleDBWithCodec[int, int64](tempPebblePath(t), encoder, decoder)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close db: %v", err)
		}
	})

	if err := db.Set(101, -77); err != nil {
		t.Fatal(err)
	}

	loaded, err := db.GetAs(101, int64(0))
	if err != nil {
		t.Fatal(err)
	}
	decoded, ok := loaded.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", loaded)
	}
	if decoded != -77 {
		t.Fatalf("decoded value mismatch: %d", decoded)
	}
}

func TestPebbleDBDefaultCodecGenericKey(t *testing.T) {
	db := NewPebbleDB[int, int64](tempPebblePath(t))
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close db: %v", err)
		}
	})

	const expected int64 = 11
	if err := db.Set(-7, expected); err != nil {
		t.Fatal(err)
	}

	loaded, err := db.Get(-7)
	if err != nil {
		t.Fatal(err)
	}
	if loaded != expected {
		t.Fatalf("unexpected loaded value: %d", loaded)
	}
}

func TestPebbleDBByteKey(t *testing.T) {
	db := NewPebbleDB[[]byte, float64](tempPebblePath(t))
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close db: %v", err)
		}
	})

	key := []byte{0x01, 0x02, 0x03}
	expected := 12.5
	if err := db.Set(key, expected); err != nil {
		t.Fatal(err)
	}

	loaded, err := db.Get([]byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Fatal(err)
	}
	if loaded != expected {
		t.Fatalf("unexpected loaded value: %f", loaded)
	}
}
