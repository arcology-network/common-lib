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
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
)

func BenchmarkBadgerBatchSet(b *testing.B) {
	opt := badger.DefaultOptions("./badger")
	bdg, err := badger.Open(opt)
	if err != nil {
		b.Error(err)
	}
	defer bdg.Close()

	unique := make(map[string]struct{})

	keys, values := setup()
	for i := range keys {
		unique[keys[i]] = struct{}{}
	}

	timer("setup", func() {
		batchSet(bdg, keys, values)

		lsm, vlog := bdg.Size()
		b.Log(lsm, vlog)
	})

	n := 50
	var sum time.Duration
	for i := 0; i < n; i++ {
		keys, values := newBlock()
		for i := range keys {
			unique[keys[i]] = struct{}{}
		}
		sum += timer("commit", func() {
			batchSet(bdg, keys, values)

			lsm, vlog := bdg.Size()
			b.Log(lsm, vlog)
		})
	}
	b.Log(len(unique))
	b.Logf("average batch write: %v", sum/time.Duration(n))

	uniqueArray := make([]string, 0, len(unique))
	for key := range unique {
		uniqueArray = append(uniqueArray, key)
	}
	rand.Shuffle(len(uniqueArray), func(i, j int) {
		uniqueArray[i], uniqueArray[j] = uniqueArray[j], uniqueArray[i]
	})
	timer("random read", func() {
		bdg.View(func(txn *badger.Txn) error {
			for i := 0; i < 1000000; i++ {
				_, err := txn.Get([]byte(uniqueArray[i]))
				if err != nil {
					b.Error("key not found")
					return nil
				}
			}
			return nil
		})
	})

	timer("iteration", func() {
		bdg.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.IteratorOptions{
				PrefetchValues: true,
				PrefetchSize:   1000,
			})
			defer it.Close()

			count := 0
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				if count%1000000 == 0 && count != 0 {
					b.Log(item.Key())
				}
				count++
			}
			b.Log(count)
			return nil
		})
	})
}

func TestBadgerIterator(t *testing.T) {
	opt := badger.DefaultOptions("./badger")
	bdg, err := badger.Open(opt)
	if err != nil {
		t.Error(err)
	}
	defer bdg.Close()

	total := 0
	for i := 0; i < 256; i++ {
		timer(fmt.Sprintf("iteration %d", i), func() {
			bdg.View(func(txn *badger.Txn) error {
				it := txn.NewIterator(badger.IteratorOptions{
					PrefetchValues: true,
					PrefetchSize:   1000,
					Prefix:         []byte{byte(i)},
				})
				defer it.Close()

				count := 0
				for it.Rewind(); it.Valid(); it.Next() {
					if count == 0 {
						t.Log(it.Item().Key())
					}
					count++
				}
				t.Logf("iteration %d", count)
				total += count
				return nil
			})
		})
	}
	t.Logf("total: %d", total)

	// timer("iteration", func() {
	// 	txn := bdg.NewTransaction(false)
	// 	it := txn.NewIterator(badger.IteratorOptions{
	// 		PrefetchValues: true,
	// 		PrefetchSize:   1000,
	// 	})

	// 	count := 0
	// 	for it.Rewind(); it.Valid(); it.Next() {
	// 		item := it.Item()
	// 		if count%10000 == 0 && count != 0 {
	// 			t.Log(item.Key())
	// 		}
	// 		count++
	// 	}
	// 	t.Log(count)
	// })
}

func batchSet(db *badger.DB, keys []string, values [][]byte) {
	index := 0
	for index < len(keys) {
		db.Update(func(txn *badger.Txn) error {
			for i := index; i < len(keys); i++ {
				if err := txn.Set([]byte(keys[i]), values[i]); err != nil {
					return nil
				} else {
					index++
				}
			}
			return nil
		})
	}
}
