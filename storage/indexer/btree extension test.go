package indexer

import (
	"testing"
	"time"

	"github.com/google/btree"
)

func TestArrayMap(t *testing.T) {
	m := make(map[uint64]string)
	t0 := time.Now()
	for i := 0; i < 1000000; i++ {
		m[uint64(i)] = ""
	}
	t1 := time.Now()
	t.Log("Array map time:", t1.Sub(t0))

	btree.New(2)
	bMap := make(map[uint64]*btree.BTree)
	t0 = time.Now()
	for i := 0; i < 1000000; i++ {
		bMap[uint64(i)] = btree.New(3)
	}
	t1 = time.Now()
	t.Log("Btree time:", t1.Sub(t0))
}

func TestBTreePerm(t *testing.T) {
	tree := btree.NewOrderedG[int](2)

	t0 := time.Now()
	for i := 0; i < 1000000; i++ {
		tree.ReplaceOrInsert(i)
	}
	t.Log("Btree ReplaceOrInsert:", time.Since(t0))

	t0 = time.Now()
	arr := Export(tree)
	t.Log(len(arr), "Btree Export:", time.Since(t0))

	t0 = time.Now()
	_, idx := FindFirstIf(tree, func(_ int, v int) bool {
		return v == 100
	})

	t.Log(idx, "Btree FindFirstIf:", time.Since(t0))
}
