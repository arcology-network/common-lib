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

package indexer

import (
	"fmt"
	"reflect"

	// "slice"
	"testing"
	"time"

	slice "github.com/arcology-network/common-lib/exp/slice"
	btree "github.com/google/btree"
)

type Int int

func TestBTree(t *testing.T) {
	tr := btree.New(4)

	for i := 0; i < 10; i++ {
		if min := tr.Min(); min != nil {
			t.Fatalf("empty min, got %+v", min)
		}
		if max := tr.Max(); max != nil {
			t.Fatalf("empty max, got %+v", max)
		}
	}
}

func TestInteger(t *testing.T) {
	type Int int

	index := NewSortedIndex("id", func(a, b Int) bool { return a < b })
	newVals := slice.ParallelTransform(make([]Int, 10), 4, func(i int, _ Int) Int { return Int(i) })

	t0 := time.Now()
	index.Add(newVals)
	fmt.Println("Insert 1000000 integers in", time.Since(t0))

	res := index.GreaterThan(5)
	if !reflect.DeepEqual(res, []Int{6, 7, 8, 9}) {
		t.Error("mismatch!!", res)
	}

	res = index.GreaterEqualThan(5)
	if !reflect.DeepEqual(res, []Int{5, 6, 7, 8, 9}) {
		t.Error("mismatch!!")
	}

	if _, ok := index.Find(Int(len(newVals))); ok {
		t.Error("Shouldn't be found!!")
	}

	res = index.LessThan(5)
	if !reflect.DeepEqual(res, []Int{4, 3, 2, 1, 0}) {
		t.Error("mismatch!!", res)
	}

	res = index.LessEqualThan(5)
	if !reflect.DeepEqual(res, []Int{5, 4, 3, 2, 1, 0}) {
		t.Error("mismatch!!", res)
	}

	res = index.Between(5, 7)
	if !reflect.DeepEqual(res, []Int{5, 6, 7}) {
		t.Error("mismatch!!", res)
	}

	index.Remove(res)

	res = index.Between(5, 7)
	if !reflect.DeepEqual(res, []Int{}) {
		t.Error("mismatch!!", res)
	}

	res = index.Export()
	if !reflect.DeepEqual(res, []Int{0, 1, 2, 3, 4, 8, 9}) {
		t.Error("mismatch!!", res)
	}
}

func TestIndex(t *testing.T) {
	type Tx struct {
		id     string
		height uint64
	}

	index := NewSortedIndex("id", func(a, b *Tx) bool { return a.id < b.id })
	txs := slice.ParallelTransform(make([]*Tx, 10), 4, func(i int, _ *Tx) *Tx { return &Tx{id: fmt.Sprint(i), height: uint64(i)} })
	index.Add(txs)

	res := index.GreaterThan(&Tx{id: "1"})
	if len(res) != 8 {
		t.Error("mismatch!!", res)
	}

	res = index.GreaterThan(&Tx{id: "2"})
	if len(res) != 7 {
		t.Error("mismatch!!", res)
	}
}

func BenchmarkInteger(t *testing.B) {
	type Int int

	index := NewSortedIndex("id", func(a, b Int) bool { return a < b })
	newVals := slice.ParallelTransform(make([]Int, 1000000), 4, func(i int, _ Int) Int { return Int(i) })

	t0 := time.Now()
	index.Add(newVals)
	fmt.Println("Add 1000000 entires in", time.Since(t0))

	t0 = time.Now()
	for i := 0; i < len(newVals); i++ {
		_, ok := index.Find(Int(i))
		if !ok {
			t.Error("not found!!")
		}
	}
	fmt.Println("Find 1000000 entires in", time.Since(t0))

	if _, ok := index.Find(Int(len(newVals))); ok {
		t.Error("Shouldn't be found!!")
	}
}
