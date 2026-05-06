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
"path/filepath"
"testing"
)

func tempPebblePath(tb testing.TB) string {
tb.Helper()
return filepath.Join(tb.TempDir(), "pebble")
}

func tempParaPebbleRoot(tb testing.TB) string {
tb.Helper()
return filepath.Join(tb.TempDir(), "pebble-shards")
}

func TestPebbleDBFunctions(t *testing.T) {
db, err := NewPebbleDB(tempPebblePath(t))
if err != nil {
t.Fatal(err)
}
t.Cleanup(func() {
if err := db.Close(); err != nil {
t.Errorf("close db: %v", err)
}
})

keys := [][]byte{
[]byte("a01"), []byte("a02"), []byte("a03"),
[]byte("b01"), []byte("c01"), []byte("d01"),
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

got, getErrs := db.GetBatch([][]byte{[]byte("a01"), []byte("b01"), []byte("c01")})
for i, err := range getErrs {
if err != nil {
t.Fatalf("GetBatch[%d]: %v", i, err)
}
}
if len(got) != 3 || got[0][0] != 1 || got[1][0] != 10 || got[2][0] != 13 {
t.Error("GetBatch failed")
}

val, err := db.Get([]byte("d01"))
if err != nil {
t.Fatal(err)
}
if val[0] != 16 {
t.Error("Get failed")
}

ok, err := db.Has([]byte("d01"))
if err != nil || !ok {
t.Error("Has failed")
}

if err := db.Delete([]byte("d01")); err != nil {
t.Fatal(err)
}
ok, err = db.Has([]byte("d01"))
if err != nil || ok {
t.Error("Delete failed")
}

delErrs := db.DeleteBatch([][]byte{[]byte("a01"), []byte("b01")})
for i, err := range delErrs {
if err != nil {
t.Fatalf("DeleteBatch[%d]: %v", i, err)
}
}
ok1, _ := db.Has([]byte("a01"))
ok2, _ := db.Has([]byte("b01"))
if ok1 || ok2 {
t.Error("DeleteBatch failed")
}

qkeys, qvalues, err := db.Query([]byte("a"), nil)
if err != nil {
t.Fatal(err)
}
if len(qkeys) != 2 || len(qvalues) != 2 {
t.Fatalf("unexpected query result: %v %v", qkeys, qvalues)
}
t.Log(qkeys)
t.Log(qvalues)
}
