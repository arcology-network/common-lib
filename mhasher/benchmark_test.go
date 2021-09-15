package mhasher

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	orderedmap "github.com/elliotchance/orderedmap"
)

func BenchmarkSortBytes1M(B *testing.B) {
	data := make([][]byte, 100000)
	for i := range data {
		data[i] = []byte("blcc:/eth10/account/alice/" + fmt.Sprint(rand.Float64()) + fmt.Sprint(rand.Float64()) + fmt.Sprint(rand.Float64()))
	}

	t0 := time.Now()
	SortBytes(data)
	fmt.Println("Sort 1M strings with SortBytes() in", time.Since(t0))

	t0 = time.Now()
	sort.SliceStable(data, func(i, j int) bool {
		return bytes.Compare(data[i], data[j]) < 0
	})
	fmt.Println("Sort 1M strings with SliceStable() in", time.Since(t0))

	t0 = time.Now()
	dict := orderedmap.NewOrderedMap()
	for i := 0; i < len(data); i++ {
		dict.Set("blcc:/eth10/account/alice/"+fmt.Sprint(rand.Float64())+fmt.Sprint(rand.Float64())+fmt.Sprint(rand.Float64()), true)
	}
	fmt.Println("ordered 1M strings with Orderedmap in", time.Since(t0))
	fmt.Println()
}

func BenchmarkUniqueBytes1MMap(B *testing.B) {
	data := make([][]byte, 1000000)
	strings := make([]string, 1000000)
	for i := range data {
		strings[i] = "blcc:/eth10/account/alice/" + fmt.Sprint(rand.Float64()) + fmt.Sprint(rand.Float64()) + fmt.Sprint(rand.Float64())
		data[i] = []byte("blcc:/eth10/account/alice/" + fmt.Sprint(rand.Float64()) + fmt.Sprint(rand.Float64()) + fmt.Sprint(rand.Float64()))
	}

	t0 := time.Now()
	UniqueBytes(data)
	fmt.Println("Unique 1M strings in", time.Since(t0))

	t0 = time.Now()
	dictMap := make(map[string]bool)
	for i := 0; i < len(data); i++ {
		dictMap[strings[i]] = true
	}
	fmt.Println("Map 1M strings in", time.Since(t0))

	t0 = time.Now()
	dict := orderedmap.NewOrderedMap()
	for i := 0; i < len(data); i++ {
		dict.Set(strings[i], true)
	}
	fmt.Println("orderedmap 1M strings in", time.Since(t0))
	fmt.Println()
}

func BenchmarkRemove100KMap(B *testing.B) {
	data := make([][]byte, 100000)
	strings := make([]string, 100000)
	for i := range data {
		strings[i] = "blcc:/eth10/account/alice/" + fmt.Sprint(rand.Float64()) + fmt.Sprint(rand.Float64()) + fmt.Sprint(rand.Float64())
		data[i] = []byte(strings[i])
	}

	t0 := time.Now()
	ordered, _ := RemoveBytes(data, data[len(strings)/4:len(strings)/2])
	fmt.Println("RemoveBytes 1M strings in", time.Since(t0), "Len:", len(ordered))

	toRemove := strings[len(strings)/4 : len(strings)/2]
	t0 = time.Now()
	dictToRemove := make(map[string]bool)
	for i := 0; i < len(toRemove); i++ {
		dictToRemove[toRemove[i]] = true
	}

	dictMap := orderedmap.NewOrderedMap()
	for i := 0; i < len(strings); i++ {
		if _, ok := dictToRemove[strings[i]]; !ok {
			dictMap.Set(strings[i], true)
		}
	}

	orderedString := make([]string, 0, dictMap.Len())
	for iter := dictMap.Front(); iter != nil; iter = iter.Next() {
		orderedString = append(orderedString, iter.Key.(string))
	}

	fmt.Println("Remove from 1M strings with Map in", time.Since(t0), "Len:", len(orderedString))
	fmt.Println()
}
