package filedb

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func BenchmarkFileDBBatchWrite(b *testing.B) {
	db, _ = NewFileDB(TEST_ROOT_PATH, 64, 2)

	keys, values := setup()
	timer("setup", func() {
		db.BatchSet(keys, values)
	})

	n := 10
	var sum time.Duration
	for i := 0; i < n; i++ {
		keys, values := newBlock()
		sum += timer("commit", func() {
			db.BatchSet(keys, values)
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
	db, _ := NewFileDB(TEST_ROOT_PATH, 128, 2)

	total := 0
	for i := 0; i < 256; i++ {
		timer(fmt.Sprintf("iteration %d", i), func() {
			keys, _, _ := db.Query(string([]byte{byte(i)}), func(pattern string, target string) bool {
				return strings.HasPrefix(target, pattern)
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
