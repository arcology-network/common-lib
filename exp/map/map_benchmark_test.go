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

package mapi

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	slice "github.com/arcology-network/common-lib/exp/slice"
)

func BenchmarkMinMax(b *testing.B) {
	ccmap := NewConcurrentMap[string, int](8, func(v int) bool { return false }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	keys := make([]string, 1000000)
	values := make([]int, len(keys))
	for i := 0; i < len(keys); i++ {
		keys[i] = fmt.Sprint(i)
		values[i] = i
	}

	ccmap.BatchSet(keys, values)

	minv := math.MaxInt
	less := func(_ string, rhs *int) {
		if minv > *rhs {
			minv = *rhs
		}
	}
	ccmap.Traverse(less)

	if minv != 0 {
		b.Error("Error: Wrong min value")
	}

	maxv := math.MinInt
	greater := func(_ string, v *int) {
		if maxv < *v {
			maxv = *v
		}
	}
	ccmap.Traverse(greater)

	if maxv != 1000000-1 {
		b.Error("Error: Wrong min value")
	}
}

func BenchmarkForeach(b *testing.B) {
	ccmap := NewConcurrentMap[string, int](8, func(v int) bool { return false }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	keys := make([]string, 1000000)
	values := make([]int, len(keys))
	for i := 0; i < len(keys); i++ {
		keys[i] = fmt.Sprint(i)
		values[i] = i
	}
	ccmap.BatchSet(keys, values)

	t0 := time.Now()
	adder := func(v int) int {
		return v + 10
	}
	ccmap.Foreach(adder)
	fmt.Println("Foreach + 10 ", time.Since(t0))
}

// func TestChecksum(t *testing.T) {
// 	ccmap := NewConcurrentMap()
// 	flags := []bool{true, true, true, true}
// 	keys := []string{"1", "2", "3", "4"}
// 	values := []interface{}{codec.String("1"), codec.String("2"), codec.Int64(3), codec.String("4")}

// 	ccmap.BatchSet(keys, values, flags)
// 	if !reflect.DeepEqual(ccmap.Checksum(), ccmap.Checksum()) {
// 		t.Error("Error: Checksums don't match")
// 	}
// }

func BenchmarkCcmapInit(b *testing.B) {
	ccmaps := make([]*ConcurrentMap[string, int], 1000)
	t0 := time.Now()
	for i := 0; i < len(ccmaps); i++ {
		ccmaps[i] = NewConcurrentMap(4, func(v int) bool { return false }, func(k string) uint64 {
			return uint64(slice.Sum[byte, int]([]byte(k)))
		})
	}
	fmt.Println("Init ConcurrentMaps", len(ccmaps), time.Since(t0))

	nativeMaps := make([]map[string]int, len(ccmaps))
	t0 = time.Now()
	for i := 0; i < len(ccmaps); i++ {
		// nativeMaps := make([]*ConcurrentMap[string, int], len(ccmaps))
		nativeMaps[i] = make(map[string]int)
	}
	fmt.Println("Init native maps", len(nativeMaps), time.Since(t0))
}

func BenchmarkCcmapBatchSet(b *testing.B) {
	genString := func() string {
		var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
		rand.Seed(time.Now().UnixNano())
		b := make([]rune, 40)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	values := make([]string, 0, 100000)
	paths := make([]string, 0, 100000)
	for i := 0; i < 100000; i++ {
		acct := genString()
		paths = append(paths, []string{
			"blcc://eth1.0/account/" + acct + "/",
			"blcc://eth1.0/account/" + acct + "/code",
			"blcc://eth1.0/account/" + acct + "/nonce",
			"blcc://eth1.0/account/" + acct + "/balance",
			"blcc://eth1.0/account/" + acct + "/defer/",
			"blcc://eth1.0/account/" + acct + "/storage/",
			"blcc://eth1.0/account/" + acct + "/storage/containers/",
			"blcc://eth1.0/account/" + acct + "/storage/native/",
			"blcc://eth1.0/account/" + acct + "/storage/containers/!/",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8111111111111111111111111111111111111",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8211111111111111111111111111111111111",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8311111111111111111111111111111111111",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8411111111111111111111111111111111111",
		}...)
	}

	for i := 0; i < len(paths); i++ {
		values = append(values, paths[i])
	}

	t0 := time.Now()
	var masterLock sync.RWMutex
	for i := 0; i < 1000000; i++ {
		masterLock.Lock()
		masterLock.Unlock()
	}
	fmt.Println("Lock() 1000000 "+fmt.Sprint(len(paths)), " in ", time.Since(t0))

	t0 = time.Now()
	ccmap := NewConcurrentMap[string, string](8, func(v string) bool { return len(v) == 0 }, func(k string) uint64 {
		return uint64(slice.Sum[byte, int]([]byte(k)))
	})

	ccmap.BatchSet(paths, values)
	fmt.Println("BatchSet "+fmt.Sprint(len(paths)), " in ", time.Since(t0))

	fmt.Println("sum:", slice.Sum[int, int](ccmap.Sizes()))

}
