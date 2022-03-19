package cachedstorage

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
)

type FileDB struct {
	rootpath string
	dirs     []string
	shards   uint8
	depth    uint8
	ext      string
}

func NewFileDB(filePath string, shards uint8, depth uint8) (*FileDB, error) {
	var err error
	if filePath, err = filepath.Abs(filePath); err != nil {
		return nil, err
	}
	filePath = filePath + "/"

	fileDB := &FileDB{
		rootpath: filePath,
		shards:   shards,
		depth:    2,
		ext:      ".dat",
	}

	if fileDB.depth > 2 {
		fileDB.depth = 2
	}

	if err := os.RemoveAll(fileDB.rootpath); err != nil && !os.IsExist(err) {
		return fileDB, err
	}

	if err := os.MkdirAll(fileDB.rootpath, os.ModePerm); err != nil {
		return fileDB, err
	}

	dirs, err := fileDB.makeDirectories(fileDB.rootpath, 0)
	fileDB.dirs = dirs
	return fileDB, err
}

func (this *FileDB) Root() string {
	return this.rootpath
}

func (this *FileDB) Directories(folder string, depth uint8) []string {
	dirs := []string{}
	if depth < this.depth {
		for i := uint8(0); i < this.shards; i++ {
			newDir := path.Join(folder, fmt.Sprint(i)) + "/"
			if depth+1 == this.depth {
				dirs = append(dirs, newDir)
			}
			ndirs := this.Directories(newDir, depth+1)
			dirs = append(dirs, ndirs...)
		}
	}
	return dirs
}

func (this *FileDB) makeDirectories(folder string, depth uint8) ([]string, error) {
	dirs := this.Directories(folder, depth)
	for i := 0; i < len(dirs); i++ {
		if err := os.MkdirAll(dirs[i], os.ModePerm); err != nil {
			return dirs, err
		}
	}
	return dirs, nil
}

func (this *FileDB) GetFileName(key string) string {
	path := this.rootpath
	for i := uint8(0); i < this.depth; i++ {
		path += fmt.Sprint(byte(key[0]%byte(this.shards))) + "/"
	}

	path += fmt.Sprint(byte(key[this.depth]%byte(this.shards))) + this.ext
	return path
}

func (this *FileDB) loadFile(file string) ([]string, [][]byte, error) {
	buffer, err := os.ReadFile(file)
	if len(buffer) == 0 || err != nil {
		return nil, nil, err
	}
	k, v := this.deserialize(buffer)
	return k, v, nil
}

func (this *FileDB) deserialize(buffer []byte) ([]string, [][]byte) {
	rawBytes := [][]byte(codec.Byteset([][]byte{}).Decode(buffer).(codec.Byteset))
	keys := []string(codec.Strings([]string{}).Decode(rawBytes[0]).(codec.Strings))
	values := [][]byte(codec.Byteset([][]byte{}).Decode(rawBytes[1]).(codec.Byteset))
	return keys, values
}

func (this *FileDB) readFile(key string) ([]byte, error) {
	file := this.GetFileName(key)
	keys, values, err := this.loadFile(file)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(keys); i++ {
		if keys[i] == key {
			return values[i], nil
		}
	}
	return nil, nil
}

func (this *FileDB) writeFile(nkey string, nvalue []byte) error {
	file := this.GetFileName(nkey)
	keys, values, err := this.loadFile(file)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		kBuffer := codec.Byteset([][]byte{codec.String(nkey).Encode()}).Encode()
		vBuffer := codec.Byteset([][]byte{codec.Bytes(nvalue).Encode()}).Encode()
		content := codec.Byteset([][]byte{kBuffer, vBuffer}).Encode()
		return os.WriteFile(file, content, os.ModePerm)
	}

	keys, values = this.updateContent(nkey, keys, nvalue, values)
	if len(values) == 0 {
		return os.Remove(file)
	}

	content := codec.Byteset([][]byte{codec.Strings(keys).Encode(), codec.Byteset(values).Encode()}).Encode()
	return os.WriteFile(file, content, os.ModePerm)
}

func (this *FileDB) updateContent(nKey string, keys []string, nVal []byte, values [][]byte) ([]string, [][]byte) {
	for i := 0; i < len(keys); i++ {
		if keys[i] == nKey {
			if nVal == nil {
				return append(keys[:i], keys[i+1:]...), append(values[:i], values[i+1:]...)
			} else {
				values[i] = nVal
				return keys, values
			}
		}
	}
	return append(keys, nKey), append(values, nVal)
}

func (this *FileDB) Set(key string, v []byte) error {
	return this.writeFile(key, v)
}

func (this *FileDB) Get(key string) ([]byte, error) {
	return this.readFile(key)
}

func (this *FileDB) BatchGet(nkeys []string) ([][]byte, error) {
	files := make([]string, len(nkeys))
	for i := 0; i < len(nkeys); i++ {
		files[i] = this.GetFileName(nkeys[i])
	}

	// Read files
	errs := make([]interface{}, len(nkeys))
	data := make([][]byte, len(nkeys))
	uniqueFiles, indices := this.CategorizeFiles(files)
	reader := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			keys, values, err := this.loadFile(uniqueFiles[i])
			if err != nil && !os.IsNotExist(err) {
				errs[i] = err
				return
			}

			for j := 0; j < len(indices[i]); j++ {
				for k := 0; k < len(keys); k++ {
					if keys[k] == nkeys[indices[i][j]] {
						data[indices[i][j]] = values[k]
					}
				}
			}
		}
	}
	common.ParallelWorker(len(uniqueFiles), 8, reader)

	common.RemoveNils(&errs)
	if len(errs) > 0 {
		return data, errs[0].(error)
	}

	return data, nil
}

func (this *FileDB) BatchSet(nkeys []string, byteset [][]byte) error {
	files := make([]string, len(nkeys))
	for i := 0; i < len(nkeys); i++ {
		files[i] = this.GetFileName(nkeys[i])
	}

	errs := make([]interface{}, len(files))
	uniqueFiles, indices := this.CategorizeFiles(files)
	maker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			if _, err := os.Stat(uniqueFiles[i]); err != nil { // File doesn't exist, create it
				if err := os.WriteFile(uniqueFiles[i], []byte{}, os.ModePerm); err != nil && errors.Is(err, os.ErrNotExist) {
					errs[i] = err
					return
				}
			}
		}
	}
	common.ParallelWorker(len(uniqueFiles), 4, maker)
	common.RemoveNils(&errs)
	if len(errs) > 0 {
		return errs[0].(error)
	}

	// Write Contents
	errs = make([]interface{}, len(uniqueFiles))
	writer := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			file := uniqueFiles[i]
			keys, values, err := this.loadFile(file)
			if err != nil && !os.IsNotExist(err) {
				errs[i] = err
				return
			}

			for j := 0; j < len(indices[i]); j++ {
				keys, values = this.updateContent(nkeys[indices[i][j]], keys, byteset[indices[i][j]], values)
			}

			content := codec.Byteset([][]byte{codec.Strings(keys).Encode(), codec.Byteset(values).Encode()}).Encode()
			if len(content) == 0 {
				errs[i] = os.Remove(file)
				return
			}

			if err := os.WriteFile(file, content, os.ModePerm); err != nil {
				errs[i] = err
				return
			}
		}
	}
	common.ParallelWorker(len(uniqueFiles), 8, writer)

	common.RemoveNils(&errs)
	if len(errs) > 0 {
		return errs[0].(error)
	}
	return nil
}

func (this *FileDB) CategorizeFiles(keys []string) ([]string, [][]uint32) {
	lookup := make(map[string][]uint32)
	for i := 0; i < len(keys); i++ {
		if v, ok := lookup[keys[i]]; !ok {
			nVal := make([]uint32, 0, len(keys)/int(this.shards))
			nVal = append(nVal, uint32(i))
			lookup[keys[i]] = nVal
		} else {
			lookup[keys[i]] = append(v, uint32(i))
		}
	}

	files := make([]string, 0, len(keys))
	indices := make([][]uint32, 0, len(lookup))
	for k, v := range lookup {
		files = append(files, k)
		indices = append(indices, v)
	}
	return files, indices
}

func (this *FileDB) Export(prefixes [][]byte) ([][]byte, error) {
	paths := []string{}
	for i := 0; i < len(prefixes); i++ {
		for j := 0; j < len(prefixes[i]); j++ {
			for k := 0; k < len(this.dirs); k++ {
				if this.dirs[k][len(this.rootpath)+i]-'0' == prefixes[i][j]%byte(this.shards) {
					paths = append(paths, this.dirs[k])
				}
			}
		}
	}
	paths = common.RemoveDuplicateStrings(&paths)
	return this.readAll(paths)
}

func (this *FileDB) ExportAll() ([][]byte, error) {
	return this.readAll(this.dirs)
}

func (this *FileDB) readAll(paths []string) ([][]byte, error) {
	sort.Strings(paths)

	errs := make([]interface{}, len(paths))
	data := make([][][]byte, len(paths))
	reader := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			var err error
			var fileInfo fs.FileInfo
			fileInfo, err = os.Stat(paths[i])
			if err != nil {
				errs[i] = err
				break
			}

			if fileInfo.IsDir() { // A directory
				var files []fs.DirEntry
				if files, err = os.ReadDir(paths[i]); err == nil { // Get all files in the directory
					for j := 0; j < len(files); j++ {
						if buffer, err := os.ReadFile(paths[i] + files[j].Name()); err == nil {
							file := paths[i] + files[j].Name()
							file = file[len(this.rootpath):]
							buffer = codec.Byteset([][]byte{codec.String(file).ToBytes(), buffer}).Encode()
							data[i] = append(data[i], buffer)
						} else {
							errs[i] = err
							break
						}
					}
				} else {
					errs[i] = err
					break
				}
			} else { // A file
				var buffer []byte
				if buffer, err = os.ReadFile(paths[i]); err == nil {
					buffer = codec.Byteset([][]byte{codec.String(paths[i]).ToBytes(), buffer}).Encode()
					data[i] = append(data[i], buffer)
				} else {
					errs[i] = err
					break
				}
			}
		}
	}
	common.ParallelWorker(len(paths), 8, reader)

	common.RemoveNils(&errs)
	if len(errs) > 0 {
		return nil, errs[0].(error)
	}

	// Flatten
	flattened := make([][]byte, 0, len(paths))
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			flattened = append(flattened, data[i][j])
		}
	}

	if len(errs) > 0 {
		return flattened, errs[0].(error)
	}
	return flattened, nil
}

func (this *FileDB) ListFiles() ([]string, error) {
	list := []string{}
	err := filepath.Walk(this.rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == this.ext {
			list = append(list, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (this *FileDB) Import(data [][]byte) error {
	for i := 0; i < len(data); i++ {
		combo := codec.Byteset{}.Decode(data[i]).(codec.Byteset)

		file := codec.Bytes(combo[0]).ToString() // file name
		// file := this.rootpath + relativePath
		reversed := common.ReverseString(file)
		parts := strings.Split(reversed, "/")
		if len(parts) < int(this.depth) {
			return errors.New("Error: Wrong path !!!")
		}
		reversed = strings.Join(parts[:this.depth+1], "/")
		name := this.rootpath + common.ReverseString(reversed)
		if err := os.WriteFile(name, combo[1], os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (this *FileDB) Equal(other *FileDB) bool {
	lfs, lerr := this.ListFiles()
	rfs, rerr := other.ListFiles()
	if len(lfs) != len(rfs) || lerr != nil || lerr != rerr {
		return false
	}

	for i := 0; i < len(lfs); i++ {
		if lfs[i][len(this.rootpath):] != rfs[i][len(other.rootpath):] { // Check file names
			return false
		}

		lhs, lerr := os.ReadFile(lfs[i])
		rhs, rerr := os.ReadFile(rfs[i])
		if !bytes.Equal(lhs, rhs) || lerr != rerr || rerr != nil { // Check file content
			return false
		}
	}
	return true
}
