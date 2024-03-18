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
	// "slice"

	"encoding/json"
	"math/big"
	"strconv"
	"testing"
	"time"

	slice "github.com/arcology-network/common-lib/exp/slice"
	"github.com/arcology-network/common-lib/storage/memdb"
)

func TestIndexerAlone(t *testing.T) {
	// Tx mocks a real transaction.
	type Tx struct {
		hash    string
		height  uint64
		gasUsed *big.Int
		time    int64 // UnixNano
		data    []byte
		Index   uint64
		From    [32]byte
		To      [32]byte
		value   *big.Int
	}

	// Create a table with two indexes.
	table := NewIndexer[*Tx](
		NewIndex("time", func(a, b *Tx) bool { return a.time < b.time }),       // Index by time.
		NewIndex("height", func(a, b *Tx) bool { return a.height < b.height }), // Index by height.
	)

	// Create 10 transactions.
	txs := slice.Transform(make([]*Tx, 10), func(i int, tx *Tx) *Tx {
		time.Sleep(1 * time.Millisecond)
		return &Tx{
			hash:    strconv.Itoa(i),
			height:  uint64(i),
			gasUsed: big.NewInt(int64(i)),
			time:    time.Now().UnixMicro(),
			data:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			Index:   uint64(i),
			From:    [32]byte{1, 2, 3, 4, 5, 6, 7, 8},
			To:      [32]byte{1, 2, 3, 4},
			value:   big.NewInt(int64(i)),
		}
	})
	table.Update(txs)

	res := table.Column("time").GreaterThan(&Tx{time: txs[5].time}) //Returns the primary keys of entries with id > 5.
	if len(res) != 4 {
		t.Error("should be 4")
	}

	res = table.Column("height").GreaterThan(&Tx{height: txs[4].height}) //Returns the primary keys of entries with id > 5.
	if len(res) != 5 {
		t.Error("should be 4")
	}
}

// The function tests the table with a database.
func TestIndexerWithDB(t *testing.T) {
	// Tx mocks a real transaction.
	type Tx struct {
		Hash    string    `json:"hash"`
		Height  uint64    `json:"height"`
		GasUsed *big.Int  `json:"gasUsed"`
		Time    time.Time `json:"time"`
		Data    []byte    `json:"data"`
		Index   uint64    `json:"index"`
		From    [32]byte  `json:"from"`
		To      [32]byte  `json:"to"`
		Value   *big.Int  `json:"value"`
	}

	// Create 10 transactions.
	txs := slice.Transform(make([]*Tx, 10), func(i int, tx *Tx) *Tx {
		time.Sleep(1 * time.Millisecond)
		return &Tx{
			Hash:    strconv.Itoa(i),
			Height:  uint64(i),
			GasUsed: big.NewInt(int64(i)),
			Time:    time.Now(),
			Data:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			Index:   uint64(i),
			From:    [32]byte{1, 2, 3, 4, 5, 6, 7, 8},
			To:      [32]byte{1, 2, 3, 4},
			Value:   big.NewInt(int64(i)),
		}
	})

	// Create a database.
	db := memdb.NewMemoryDB() // Tx DB

	// Extract the primary keys.
	keys := slice.Transform(txs, func(_ int, tx *Tx) string {
		return string(tx.Hash[:])
	})

	// Encode the transactions.
	encoded := slice.Transform(txs, func(_ int, tx *Tx) []byte {
		data, err := json.Marshal(tx)
		if err != nil {
			panic(err)
		}
		return data
	})

	// Decode the transactions.
	queryTxs := slice.Transform(encoded, func(i int, data []byte) *Tx {
		tx := &Tx{}
		if err := json.Unmarshal(data, tx); err != nil {
			panic(err)
		}
		return tx
	})

	// Save the transactions to the database.
	db.BatchSet(keys, encoded)

	// Create a table with two indexes for the transactions.
	type txIndex struct {
		time       time.Time
		height     uint64
		primaryKey string // Key is the primary key in the database.
	}

	// Create the indexes for the transactions.
	indics := slice.Transform(txs, func(i int, tx *Tx) *txIndex {
		return &txIndex{
			time:       tx.Time,
			height:     tx.Height,
			primaryKey: tx.Hash,
		}
	})

	// Create a table with two indexes.
	table := NewIndexer[*txIndex](
		NewIndex("time", func(a, b *txIndex) bool { return a.time.Nanosecond() < b.time.Nanosecond() }), // Index by time.
		NewIndex("height", func(a, b *txIndex) bool { return a.height < b.height }),                     // Index by height.
	)

	// Update the table with the indexes.
	table.Update(indics)

	// Query the table with the time index.
	res := table.Column("time").GreaterThan(&txIndex{time: txs[5].Time}) //Returns the primary keys of entries with id > 5.
	if len(res) != 4 {
		t.Error("should be 4")
	}

	// Get the primary keys of the query results.
	primaryKeys := slice.Transform(res, func(i int, tx *txIndex) string {
		return tx.primaryKey
	})

	// Get the transactions from the database using the primary keys.
	encoded, _ = db.BatchGet(primaryKeys)

	queryTxs = slice.Transform(encoded, func(i int, data []byte) *Tx {
		tx := &Tx{}
		json.Unmarshal(data, tx)
		return tx
	})

	if queryTxs[0].Hash != (*txs[6]).Hash || queryTxs[1].Hash != (*txs[7]).Hash {
		t.Error("should be equal")
	}

	// Delete the transaction indices from the table.
	table.Remove(res)

	// Query the table with the time index again, should be empty.
	res = table.Column("time").GreaterThan(&txIndex{time: txs[5].Time})
	if len(res) != 0 {
		t.Error("should be 0")
	}

	// Count the number of indices left in the table.
	res = table.Column("time").LessEqualThan(&txIndex{time: txs[5].Time})
	if len(res) != 6 {
		t.Error("should be 6")
	}
}
