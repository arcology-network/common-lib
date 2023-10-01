package cachedstorage

import (
	"bytes"
	"os"
	"testing"
)

func TestBadgerDBFunctions(t *testing.T) {
	db := NewBadgerDB("./badger-test/")
	db.BatchSet([]string{
		"a01",
		"a02",
		"a03",
		"b01",
		"c01",
		"d01",
	}, [][]byte{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
		{13, 14, 15},
		{16, 17, 18},
	})

	values, _ := db.BatchGet([]string{
		"a01",
		"b01",
		"c01",
	})
	if len(values) != 3 ||
		!bytes.Equal(values[0], []byte{1, 2, 3}) ||
		!bytes.Equal(values[1], []byte{10, 11, 12}) ||
		!bytes.Equal(values[2], []byte{13, 14, 15}) {
		t.Error("BatchGet Failed")
	}

	value, _ := db.Get("d01")
	if !bytes.Equal(value, []byte{16, 17, 18}) {
		t.Error("Get Failed")
	}

	keys, values, _ := db.Query("a", nil)
	t.Log(keys)
	t.Log(values)

	err := os.RemoveAll("./badger-test/")
	if err != nil {
		t.Log(keys)
	}
}
