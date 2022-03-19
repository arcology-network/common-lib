package cachedstorage

import (
	"bytes"
	"testing"
)

func TestMemDB(t *testing.T) {
	memDB := NewMemDB()
	keys := []string{"123", "456"}
	values := make([][]byte, 2)
	values[0] = []byte{1, 2, 3}
	values[1] = []byte{4, 5, 6}
	memDB.BatchSet(keys, values)

	if v, _ := memDB.Get(keys[0]); !bytes.Equal(v, values[0]) {
		t.Error("Error")
	}

	if v, _ := memDB.Get(keys[1]); !bytes.Equal(v, values[1]) {
		t.Error("Error")
	}

	retrived, _ := memDB.BatchGet(append(keys, ""))
	if len(values) != 2 || !bytes.Equal(values[0], retrived[0]) || !bytes.Equal(values[1], retrived[1]) {
		t.Error("Error")
	}
}
