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

package filedb

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func BenchmarkFileDBBatchWrite(b *testing.B) {
	root := testFileDBRoot(b)
	db, _ = NewFileDB(root, 64, 2)

	keys, values := setup()
	timer("setup", func() {
		if err := db.SetBatch(keys, values); err != nil {
			b.Fatal(err)
		}
	})

	n := 10
	var sum time.Duration
	for i := 0; i < n; i++ {
		keys, values := newBlock()
		sum += timer("commit", func() {
			if err := db.SetBatch(keys, values); err != nil {
				b.Fatal(err)
			}
		})
	}
	b.Logf("average batch write: %v", sum/time.Duration(n))

	// total := 0
	// for i := 0; i < 256; i++ {
	// 	timer(fmt.Sprintf("iteration %d", i), func() {
	// 		keys, _, _ := db.Query(string([]byte{byte(i)}), func(pattern string, target string) bool {
	// 			return strings.HasPrefix(target, pattern)
	// 		})
	// 		if len(keys) != 0 {
	// 			b.Log([]byte(keys[0]))
	// 		}
	// 		b.Log(len(keys))
	// 		total += len(keys)
	// 	})
	// }
	// b.Logf("total: %d", total)
}

func BenchmarkFileDBQuery(b *testing.B) {
	root := testFileDBRoot(b)
	db, _ := NewFileDB(root, 128, 2)

	total := 0
	for i := 0; i < 256; i++ {
		timer(fmt.Sprintf("iteration %d", i), func() {
			keys, _, _ := db.Query(string([]byte{byte(i)}), func(pattern string, target []byte) bool {
				return strings.HasPrefix(string(target), pattern)
			})
			if len(keys) != 0 {
				b.Log(keys[0])
			}
			b.Log(len(keys))
			total += len(keys)
		})
	}
	b.Logf("total: %d", total)
}
