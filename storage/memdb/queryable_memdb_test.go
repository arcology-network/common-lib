/*
 *   Copyright (c) 2023 Arcology Network

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
	"fmt"
	"testing"
	"time"

	slice "github.com/arcology-network/common-lib/exp/slice"
	// queryablecache "github.com/arcology-network/common-lib/storage/cache"
)

type CachedTx struct {
	Hash        string      // Transaction hash.
	Tx          []byte      // The encoded transaction.
	Height      uint64      // Block height.
	Index       uint64      // The index of the transaction in the block.
	ReceiptHash string      // The hash of the transaction receipt.
	ExecRecipt  interface{} // The execution receipt of the transaction.
}

type CachedBlock struct {
	Hash   string
	Height uint64
}

func TestQueryCacheTx(t *testing.T) {
	txTable := NewTable("tx",
		NewIndex("Hash", true, true, new(string)),
		NewIndex("Index", false, false, new(uint64)),
		NewIndex("Height", false, false, new(uint64)),
	)

	blockTable := NewTable("block",
		NewIndex("Height", true, true, new(uint64)),
		NewIndex("Hash", false, true, new(string)),
	)

	cache, err := NewQueryableCache(nil, txTable, blockTable)
	if err != nil {
		panic(err)
	}

	t.Run("block query", func(t *testing.T) {
		blocks := []*CachedBlock{
			{
				Height: 30,
				Hash:   "0x111111",
			},
			{
				Height: 32,
				Hash:   "0x222222",
			},
		}

		if err := cache.Add("block", slice.Transform(blocks, func(_ int, v *CachedBlock) interface{} { return v })...); err != nil {
			panic(err)
		}

		// CachedBlock table
		if raw, err := cache.FindFirst("block", "id", uint64(30)); err != nil || raw == nil {
			panic(err)
		}

		if raw, err := cache.FindFirst("block", "Hash", "0x222222"); err != nil || raw == nil {
			panic(err)
		}

		if raw, err := cache.FindFirst("block", "id", uint64(32)); err != nil || raw == nil {
			panic(err)
		}

		if raw, err := cache.FindLessThan("block", "id", uint64(40)); err != nil || len(raw) != 2 {
			panic(err)
		}

		if raw, err := cache.FindLessThan("block", "id", uint64(29)); err != nil || len(raw) != 0 {
			t.Error("should be empty")
		}

		if raw, err := cache.FindLessThan("block", "id", uint64(30)); err != nil || len(raw) != 1 {
			t.Error("should be 1")
		}

		if raw, err := cache.FindGreaterThan("block", "id", uint64(30)); err != nil || len(raw) != 2 {
			t.Error("should be 2")
		}
	})

	t.Run("tx query", func(t *testing.T) {
		txs := []*CachedTx{
			{
				Height: 5,
				Index:  1111,
				Hash:   "0x5",
			},
			{
				Height: 5,
				Index:  2222,
				Hash:   "0x6",
			},
		}

		if err := cache.Add("tx", slice.Transform(txs, func(_ int, v *CachedTx) interface{} { return v })...); err != nil {
			panic(err)
		}

		// Transtion table
		if raw, err := cache.FindFirst("tx", "id", "0x5"); err != nil || raw == nil {
			t.Error(err)
		}

		if raw, err := cache.FindFirst("tx", "id", "0x6"); err != nil || raw == nil {
			t.Error(err)
		}

		if raw, err := cache.FindFirst("tx", "Index", uint64(2222)); err != nil || raw == nil {
			t.Error(err)
		}

		if raw, err := cache.FindFirst("tx", "Index", uint64(2)); err != nil || raw != nil {
			t.Error("should be 0")
		}

		if raw, err := cache.FindLessThan("tx", "Index", uint64(6)); err != nil || raw == nil {
			t.Error("should be 0")
		}

		if raw, err := cache.FindGreaterThan("tx", "Index", uint64(6)); err != nil || len(raw) != 2 {
			t.Error("should be 2")
		}

		raw, _ := cache.FindGreaterThan("tx", "Index", uint64(1112))
		if err != nil || len(raw) != 1 {
			t.Error("should be 1")
		}

		if raw, err := cache.FindAll("tx", "Height", uint64(5)); err != nil || len(raw) != 2 {
			t.Error("should be 2")
		}

		cache.Remove("tx", raw[0])

		raw, _ = cache.FindGreaterThan("tx", "Index", uint64(1112))
		if err != nil || len(raw) != 0 {
			t.Error("should be 0")
		}

		if raw, err := cache.FindGreaterThan("tx", "Index", uint64(1111)); err != nil || len(raw) != 1 {
			t.Error("should be 2")
		}

		if raw, err := cache.FindAll("tx", "Height", uint64(5)); err != nil || len(raw) != 1 {
			t.Error("should be 1")
		}
	})
}

func TestQueryCacheTxPerformance1M(t *testing.T) {
	txTable := NewTable("tx",
		NewIndex("Hash", true, true, new(string)),
		NewIndex("Index", false, false, new(uint64)),
	)

	blockTable := NewTable("block",
		NewIndex("Height", true, true, new(uint64)),
		NewIndex("Hash", false, true, new(string)),
	)

	cache, err := NewQueryableCache(nil, txTable, blockTable)
	if err != nil {
		panic(err)
	}

	txs := make([]*CachedTx, 10)
	for i := range txs {
		txs[i] = &CachedTx{
			Height: uint64(i),
			Index:  uint64(i % 50000),
			Hash:   "0x" + fmt.Sprint(i),
		}
	}
	txInterfaces := slice.Transform(txs, func(_ int, v *CachedTx) interface{} { return v })

	t0 := time.Now()
	if err := cache.Add("tx", txInterfaces...); err != nil {
		panic(err)
	}
	fmt.Println("add txs:", len(txs), time.Since(t0))

	t0 = time.Now()
	// for i := range txs {
	if raw, err := cache.FindFirst("tx", "Index", uint64(1)); err != nil || raw == nil {
		panic(err)
	}
	// }
	fmt.Println("add txs:", len(txs), time.Since(t0))

	t0 = time.Now()
	if err := cache.Remove("tx", txInterfaces...); err != nil {
		panic(err)
	}
	fmt.Println("remove txs:", len(txs), time.Since(t0))
}
