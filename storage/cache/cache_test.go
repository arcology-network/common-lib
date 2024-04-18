package cache

import (
	"testing"

	xxhash "github.com/cespare/xxhash/v2"
	"github.com/ethereum/go-ethereum/common/math"
)

func TestCache(t *testing.T) {
	readCache := NewReadCache[string](4,
		func(v int) bool {
			return v == math.MaxInt32
		},
		func(k string) uint64 {
			return uint64(xxhash.Sum64([]byte(k)))
		},
		nil,
	)

	// readCache.Update([]string{"123", "456", "789"}, []int{1, 2, 3})
	readCache.Precommit([]string{"123", "456", "789"}, []int{1, 2, 3})

	if v, ok := readCache.Get("123"); !ok || *v != 1 {
		t.Error("Error: Values mismatched !")
	}

	if v, ok := readCache.Get("456"); !ok || *v != 2 {
		t.Error("Error: Values mismatched !")
	}

	if v, ok := readCache.Get("789"); !ok || *v != 3 {
		t.Error("Error: Values mismatched !")
	}

	if readCache.Length() != 3 {
		t.Error("Error: Values mismatched !")
	}

	// readCache.Update([]string{"444", "555", "666"}, []int{4, 5, 6})
	readCache.Precommit([]string{"444", "555", "666"}, []int{4, 5, 6})

	if v, ok := readCache.Get("444"); !ok || *v != 4 {
		t.Error("Error: Values mismatched !")
	}

	if v, ok := readCache.Get("555"); !ok || *v != 5 {
		t.Error("Error: Values mismatched !")
	}

	if v, ok := readCache.Get("666"); !ok || *v != 6 {
		t.Error("Error: Values mismatched !")
	}

	if readCache.Length() != 6 {
		t.Error("Error: Values mismatched !")
	}

	readCache.Commit([]string{"444", "456", "666"}, []int{7, 8, 9})

	if v, ok := readCache.Get("444"); !ok || *v != 7 {
		t.Error("Error: Values mismatched !", *v)
	}

	if v, ok := readCache.Get("456"); !ok || *v != 8 {
		t.Error("Error: Values mismatched !", *v)
	}

	if v, ok := readCache.Get("666"); !ok || *v != 9 {
		t.Error("Error: Values mismatched !", *v)
	}

	if readCache.Length() != 6 {
		t.Error("Error: Values mismatched !")
	}

}
