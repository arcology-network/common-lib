package filedb

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

var (
	TEST_ROOT_PATH   = path.Join(os.TempDir(), "/filedb/")
	TEST_BACKUP_PATH = path.Join(os.TempDir(), "/filedb-back/")
)

func TestFileDB(t *testing.T) {
	fileDB, err := NewFileDB(TEST_ROOT_PATH, 8, 2)
	if err != nil {
		t.Error(err)
	}

	keys := []string{"123", "456"}
	values := make([][]byte, 2)
	values[0] = []byte{1, 2, 3}
	values[1] = []byte{4, 5, 6}

	if err := fileDB.Set(keys[0], values[0]); err != nil {
		t.Error(err)
	}

	if err := fileDB.Set(keys[1], values[1]); err != nil {
		t.Error(err)
	}

	if v, _ := fileDB.Get(keys[0]); !bytes.Equal(v, values[0]) {
		t.Error("Error")
	}

	if v, _ := fileDB.Get(keys[1]); !bytes.Equal(v, values[1]) {
		t.Error("Error")
	}

	// Delete the entry
	if err := fileDB.Set(keys[0], nil); err != nil {
		t.Error(err)
	}

	if v, _ := fileDB.Get(keys[0]); v != nil {
		t.Error("Error: Should have been deleted already !")
	}

	if err := fileDB.Set(keys[1], nil); err != nil {
		t.Error(err)
	}

	if v, _ := fileDB.Get(keys[1]); v != nil {
		t.Error("Error: Should have been deleted already !")
	}

	if files, err := fileDB.ListFiles(); err != nil {
		t.Error(err)
	} else {
		if len(files) != 0 {
			t.Error("Error: All deleted")
		}
	}
	os.RemoveAll(fileDB.rootpath)
}

func TestFileDBBatch(t *testing.T) {
	fileDB, err := NewFileDB(TEST_ROOT_PATH, 8, 2)
	if err != nil {
		t.Error(err)
	}

	keys := []string{"123", "456"}
	values := make([][]byte, 2)
	values[0] = []byte{1, 2, 3}
	values[1] = []byte{4, 5, 6}

	if err := fileDB.BatchSet(keys, values); err != nil {
		t.Error(err)
	}

	if v, _ := fileDB.Get(keys[0]); !bytes.Equal(v, values[0]) {
		t.Error("Error")
	}

	if v, _ := fileDB.Get(keys[1]); !bytes.Equal(v, values[1]) {
		t.Error("Error")
	}

	if v, _ := fileDB.BatchGet(keys); len(v) != 2 || !bytes.Equal(v[0], values[0]) || !bytes.Equal(v[1], values[1]) {
		t.Error("Error")
	}
	os.RemoveAll(fileDB.rootpath)
}

func TestFileDbBatch(t *testing.T) {
	fileDB, err := NewFileDB(TEST_ROOT_PATH, 16, 2)

	if err != nil {
		t.Error(err)
	}

	keys := make([]string, 10000)
	values := make([][]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		buffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(buffer, uint32(i))
		k := sha256.Sum256(buffer)
		values[i] = buffer
		keys[i] = string(k[:])
	}

	t0 := time.Now()
	if err := fileDB.BatchSet(keys, values); err != nil {
		t.Error(err)
	}

	if retrived, err := fileDB.BatchGet(keys); err == nil {
		for i := 0; i < len(keys); i++ {
			if !bytes.Equal(retrived[i], values[i]) {
				t.Error("Error: Mismatch !!!")
			}
		}
	} else {
		t.Error(err)
	}

	for i := 0; i < len(keys); i++ {
		if retrived, err := fileDB.Get(keys[i]); err == nil {
			if !bytes.Equal(retrived, values[i]) {
				t.Error("Error: Mismatch !!!")
			}
		}
	}

	fmt.Println("Batch Set ", len(keys), " Entries from files:", time.Since(t0))
	os.RemoveAll(fileDB.rootpath)
}

func TestFileDbExport(t *testing.T) {
	fileDB, err := NewFileDB(TEST_ROOT_PATH, 4, 2)
	if err != nil {
		t.Error(err)
	}

	keys := make([]string, 10)
	values := make([][]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		buffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(buffer, uint32(i))
		k := sha256.Sum256(buffer)
		values[i] = buffer
		keys[i] = string(k[:])
	}

	if err := fileDB.BatchSet(keys, values); err != nil {
		t.Error(err)
	}

	prefixes := [][]byte{{0}}
	if data, err := fileDB.Export(prefixes); err != nil || len(data) != 2 {
		t.Error(err)
	}
	os.RemoveAll(fileDB.rootpath)
}

func TestFileDbExportAll(t *testing.T) {
	fileDB, err := NewFileDB(TEST_ROOT_PATH, 4, 2)
	if err != nil {
		t.Error(err)
	}

	keys := make([]string, 10)
	values := make([][]byte, len(keys))
	inHashes := make([][32]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		buffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(buffer, uint32(i))
		k := sha256.Sum256(buffer)
		values[i] = buffer
		keys[i] = string(k[:])
		inHashes[i] = sha256.Sum256(buffer)
	}

	if err := fileDB.BatchSet(keys, values); err != nil {
		t.Error(err)
	}

	data, err := fileDB.ExportAll()
	if err != nil || len(data) != 8 {
		t.Error(err)
	}
	fs, err := fileDB.ListFiles()
	if fs == nil || err != nil {
		t.Error(err)
	}

	fileDb, err := LoadFileDB(TEST_ROOT_PATH, 4, 2)
	if err == nil {
		if err := fileDb.Import(data); err != nil {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}

	if !fileDB.Equal(fileDb) {
		t.Error("Error: Two files are different")
	}
	os.RemoveAll(fileDB.rootpath)
	os.RemoveAll(fileDb.rootpath)
}

func TestLoadFileDB(t *testing.T) {
	fileDB, err := NewFileDB(TEST_ROOT_PATH, 4, 2)
	if err != nil {
		t.Error(err)
	}

	keys := make([]string, 10)
	values := make([][]byte, len(keys))
	inHashes := make([][32]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		buffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(buffer, uint32(i))
		k := sha256.Sum256(buffer)
		values[i] = buffer
		keys[i] = string(k[:])
		inHashes[i] = sha256.Sum256(buffer)
	}

	if err := fileDB.BatchSet(keys, values); err != nil {
		t.Error(err)
	}

	data, err := fileDB.ExportAll()
	if err != nil || len(data) != 8 {
		t.Error(err)
	}

	fileDb, err := LoadFileDB(TEST_ROOT_PATH, 4, 2)
	if err == nil {
		if err := fileDb.Import(data); err != nil {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}

	if !fileDB.Equal(fileDb) {
		t.Error("Error: Two files are different")
	}

	os.RemoveAll(TEST_ROOT_PATH)
	os.RemoveAll(TEST_BACKUP_PATH)
}

func BenchmarkFileDbBatch(b *testing.B) {
	fileDB, err := NewFileDB(TEST_ROOT_PATH, 128, 2)
	if err != nil {
		b.Error(err)
	}

	keys := make([]string, 2000000)
	values := make([][]byte, len(keys))
	for i := 0; i < len(keys); i++ {
		buffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(buffer, uint32(i))
		k := sha256.Sum256(buffer)
		values[i] = buffer
		keys[i] = string(k[:])
	}

	t0 := time.Now()
	if err := fileDB.BatchSet(keys, values); err != nil {
		b.Error(err)
	}
	fmt.Println("BatchSet() ", len(keys), " Entries from files:", time.Since(t0))

	t0 = time.Now()
	if _, err := fileDB.BatchGet(keys); err != nil {
		b.Error(err)
	}
	fmt.Println("BatchGet() ", len(keys), " Entries from files:", time.Since(t0))

	t0 = time.Now()
	if err := fileDB.BatchSet(keys, values); err != nil {
		b.Error(err)
	}
	fmt.Println("BatchSet() ", len(keys), " Entries from files:", time.Since(t0))

	os.RemoveAll(fileDB.rootpath)
}
