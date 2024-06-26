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

package merkle

// import (
// 	"fmt"
// 	"testing"
// 	"time"

// 	common "github.com/arcology-network/common-lib/common"
// 	mempool "github.com/arcology-network/common-lib/mempool"
// )

// func BenchmarkMerkle10kAcct(b *testing.B) {
// 	t0 := time.Now()
// 	trees := make([]*Merkle, 100000)
// 	for i := 0; i < len(trees); i++ {
// 		trees[i] = NewMerkle(2, Sha256)
// 	}

// 	bytes := make([][]byte, 0)
// 	for j := 0; j < 10; j++ {
// 		bytes = append(bytes, []byte(fmt.Sprint(j)))
// 	}
// 	fmt.Println("append", fmt.Sprint(len(bytes)), "leaf nodes in ", time.Since(t0))

// 	t0 = time.Now()
// 	mempool := mempool.NewMempool("trees", func() interface{} {
// 		return NewNode()
// 	})

// 	// nodePool := mempool.NewMempool("nodes", func() interface{} {
// 	// 	return NewNode()
// 	// })

// 	worker := func(start, end, index int, args ...interface{}) {
// 		for i := start; i < end; i++ {
// 			trees[i].Init(bytes, mempool)
// 		}
// 	}
// 	common.ParallelWorker(len(trees), 6, worker)
// 	fmt.Println("Build NewMerkle with", fmt.Sprint(len(bytes)), "leaf nodes in ", time.Since(t0))
// }

// func BenchmarkMerkle(b *testing.B) {
// 	t0 := time.Now()
// 	bytes := make([][]byte, 1000000)
// 	for i := 0; i < len(bytes); i++ {
// 		bytes[i] = []byte(fmt.Sprint(i))
// 	}
// 	fmt.Println("append", fmt.Sprint(len(bytes)), "leaf nodes in ", time.Since(t0))

// 	t0 = time.Now()
// 	treePool := mempool.NewMempool("trees", func() interface{} {
// 		return NewNode()
// 	})

// 	NewMerkle(16, Sha256, treePool).Init(bytes)
// 	fmt.Println("Build NewMerkle with", fmt.Sprint(len(bytes)), "leaf nodes in ", time.Since(t0))
// }
