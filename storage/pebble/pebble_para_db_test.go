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

	common "github.com/arcology-network/common-lib/common"
)

func TestParaPebbleDBFunctions(t *testing.T) {
	db := NewParaPebbleDB[string, uint64](tempParaPebbleRoot(t), common.Remainder)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close db: %v", err)
		}
	})

	data := []uint64{
		1,
		4,
		7,
		10,
		13,
		16,
	}
	keys := []string{
		"a01",
		"a02",
		"a03",
		"b01",
		"c03",
		"d01",
	}

	if err := db.SetBatch(keys, data); err != nil {
		t.Fatal(err)
	}

	if _, err := db.GetBatch(keys); err != nil {
		t.Fatal(err)
	}

	values, err := db.GetBatch([]string{"a01", "b01", "c03"})
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

	queryKeys, queryValues, err := db.Query("a", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(queryKeys)
	t.Log(queryValues)
}

func TestParaPebbleDBCodecHooks(t *testing.T) {
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
	shard := func(numOfShard int, key int) int { return key % numOfShard }

	db := NewParaPebbleDBWithCodec[int, int64](tempParaPebbleRoot(t), shard, encoder, decoder)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close db: %v", err)
		}
	})

	if err := db.Set(202, -99); err != nil {
		t.Fatal(err)
	}

	loaded, err := db.GetAs(202, int64(0))
	if err != nil {
		t.Fatal(err)
	}
	decoded, ok := loaded.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", loaded)
	}
	if decoded != -99 {
		t.Fatalf("decoded value mismatch: %d", decoded)
	}
}

func TestParaPebbleDBByteKey(t *testing.T) {
	shard := func(numOfShard int, key []byte) int {
		total := 0
		for i := 0; i < len(key); i++ {
			total += int(key[i])
		}
		return total % numOfShard
	}

	db := NewParaPebbleDB[[]byte, float64](tempParaPebbleRoot(t), shard)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close db: %v", err)
		}
	})

	key := []byte{0x0a, 0x0b}
	expected := 13.25
	if err := db.Set(key, expected); err != nil {
		t.Fatal(err)
	}

	loaded, err := db.Get([]byte{0x0a, 0x0b})
	if err != nil {
		t.Fatal(err)
	}
	if loaded != expected {
		t.Fatalf("unexpected loaded value: %f", loaded)
	}
}
