package merklepatriciatrie

import (
	"fmt"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/merkle"
	"github.com/stretchr/testify/require"
)

func TestMpt(t *testing.T) {
	keys := make([][]byte, 10000)
	data := make([][]byte, len(keys))
	for i := 0; i < len(data); i++ {
		keys[i] = merkle.Sha256{}.Hash([]byte(fmt.Sprint(i)))
		data[i] = merkle.Sha256{}.Hash([]byte(fmt.Sprint(i)))
	}

	t.Run("Proof test", func(t *testing.T) {
		trie := NewTrie()
		for i := 0; i < len(data); i++ {
			trie.Put(keys[i], data[i])
		}

		for i := 0; i < 50; i++ {
			if _, ok := trie.Prove(keys[i]); !ok {
				t.Error("Error: Proof not found")
				return
			}
		}
	})

	t.Run("Performance test", func(t *testing.T) {
		t0 := time.Now()
		trie := NewTrie()
		for i := 0; i < len(data); i++ {
			trie.Put(keys[i], data[i])
		}
		roothash := trie.Hash()
		fmt.Println("Serial put: "+fmt.Sprint(len(data)), time.Since(t0))

		t0 = time.Now()
		trie = NewTrie()
		paraRoothash := ParallelInserter{}.Insert(trie, keys, data)
		fmt.Println("ParallelInserter put: "+fmt.Sprint(len(data)), time.Since(t0))

		require.Equal(t, paraRoothash, roothash)
	})
}
