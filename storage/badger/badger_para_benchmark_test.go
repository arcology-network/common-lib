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
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/common"
)

func BenchmarkParaBadgerBatchSet(b *testing.B) {
	os.RemoveAll(TEST_ROOT_PATH)
	fileDB := NewParaBadgerDB(TEST_ROOT_PATH, common.Remainder)

	keys := make([]string, 2000000)
	values := make([][]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		buffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(buffer, uint32(i))
		k := sha256.Sum256(buffer)
		values[i] = buffer
		keys[i] = string(k[:])
	}

	t0 := time.Now()
	if err := fileDB.BatchSet(keys, values); err != nil {
		b.Error(err)
	}
	fmt.Println("BatchSet() ", len(keys), " Entries from files:", time.Since(t0))

	t0 = time.Now()
	if _, err := fileDB.BatchGet(keys); err != nil {
		b.Error(err)
	}
	fmt.Println("BatchGet() ", len(keys), " Entries from files:", time.Since(t0))
	os.RemoveAll(TEST_ROOT_PATH)
}

func BenchmarkBadgerBatchSet2(b *testing.B) {
	os.RemoveAll(TEST_ROOT_PATH)
	fileDB := NewBadgerDB(TEST_ROOT_PATH)

	keys := make([]string, 2000000)
	values := make([][]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		buffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(buffer, uint32(i))
		k := sha256.Sum256(buffer)
		values[i] = buffer
		keys[i] = string(k[:])
	}

	t0 := time.Now()
	if err := fileDB.BatchSet(keys, values); err != nil {
		b.Error(err)
	}
	fmt.Println("BatchSet() ", len(keys), " Entries from files:", time.Since(t0))

	t0 = time.Now()
	if _, err := fileDB.BatchGet(keys); err != nil {
		b.Error(err)
	}
	fmt.Println("BatchGet() ", len(keys), " Entries from files:", time.Since(t0))
	os.RemoveAll(TEST_ROOT_PATH)
}

// func TestParaBadgerIterator(t *testing.T) {
// 	opt := badger.DefaultOptions("./badger")
// 	bdg, err := badger.Open(opt)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer bdg.Close()

// 	total := 0
// 	for i := 0; i < 256; i++ {
// 		timer(fmt.Sprintf("iteration %d", i), func() {
// 			bdg.View(func(txn *badger.Txn) error {
// 				it := txn.NewIterator(badger.IteratorOptions{
// 					PrefetchValues: true,
// 					PrefetchSize:   1000,
// 					Prefix:         []byte{byte(i)},
// 				})
// 				defer it.Close()

// 				count := 0
// 				for it.Rewind(); it.Valid(); it.Next() {
// 					if count == 0 {
// 						t.Log(it.Item().Key())
// 					}
// 					count++
// 				}
// 				t.Logf("iteration %d", count)
// 				total += count
// 				return nil
// 			})
// 		})
// 	}
// 	t.Logf("total: %d", total)
// 	os.RemoveAll("./badger")
// 	// timer("iteration", func() {
// 	// 	txn := bdg.NewTransaction(false)
// 	// 	it := txn.NewIterator(badger.IteratorOptions{
// 	// 		PrefetchValues: true,
// 	// 		PrefetchSize:   1000,
// 	// 	})

// 	// 	count := 0
// 	// 	for it.Rewind(); it.Valid(); it.Next() {
// 	// 		item := it.Item()
// 	// 		if count%10000 == 0 && count != 0 {
// 	// 			t.Log(item.Key())
// 	// 		}
// 	// 		count++
// 	// 	}
// 	// 	t.Log(count)
// 	// })
// }

// func ParaBatchSet(db *badger.DB, keys []string, values [][]byte) {
// 	index := 0
// 	for index < len(keys) {
// 		db.Update(func(txn *badger.Txn) error {
// 			for i := index; i < len(keys); i++ {
// 				if err := txn.Set([]byte(keys[i]), values[i]); err != nil {
// 					return nil
// 				} else {
// 					index++
// 				}
// 			}
// 			return nil
// 		})
// 	}
// }
