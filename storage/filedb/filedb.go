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

package filedb

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	slice "github.com/arcology-network/common-lib/exp/slice"
	intf "github.com/arcology-network/common-lib/storage/interface"
)

const (
	MAX_DEPTH    = 4
	MAX_SHARD    = 256
	VERSION_FILE = "version.txt"
	EXTENSIION   = ".dat"
)

type FileDB struct {
	rootpath string
	dirs     []string
	files    []string
	shards   uint32
	depth    uint8
}

func LoadFileDB(rootpath string, shards uint32, depth uint8) (*FileDB, error) {
	fileDB := &FileDB{
		rootpath: path.Join(rootpath, "/") + "/",
		shards:   shards,
		depth:    depth,
	}

	if files, err := fileDB.ListFiles(); err == nil {
		fileDB.files = files
		fileDB.dirs = fileDB.Directories(fileDB.rootpath, depth)
	} else {
		return nil, err
	}
	return fileDB, nil
}

func NewFileDB(rootPath string, shards uint32, depth uint8) (*FileDB, error) {
	var err error
	if rootPath, err = filepath.Abs(rootPath); err != nil {
		return nil, err
	}

	if depth >= MAX_DEPTH {
		return nil, errors.New("Error: Excessed max depth ")
	}

	if shards >= MAX_SHARD {
		return nil, errors.New("Error: Excessed max depth ")
	}

	fileDB := &FileDB{
		rootpath: path.Join(rootPath, "/") + "/",
		shards:   shards,
		depth:    2,
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
	fileDB.files = make([]string, 0, 1024)
	return fileDB, err
}

func (this *FileDB) Type() uint8 {
	return intf.PERSISTENT_DB
}

func (this *FileDB) Root() string {
	return this.rootpath
}

func (this *FileDB) SetVer(newVer uint64) (uint64, error) {
	version, err := this.GetVer()
	if err == nil {
		if version <= newVer {
			return version, errors.New("Error: The new version has to be greater than the current !")
		}

		if version+1 != newVer {
			return version, errors.New("Error: The version numbers have to be consecutive !")
		}

		err = os.WriteFile(path.Join(this.rootpath, VERSION_FILE), codec.Uint64(newVer).Encode(), os.ModePerm)
		if err == nil {
			return this.GetVer()
		}
	}
	return math.MaxUint64, err
}

func (this *FileDB) GetVer() (uint64, error) {
	if verBytes, err := os.ReadFile(path.Join(this.rootpath, VERSION_FILE)); err == nil {
		version := codec.Uint64(0).Decode(verBytes).(uint64)
		return version, nil
	} else {
		return math.MaxUint64, err
	}
}

func (this *FileDB) Directories(folder string, depth uint8) []string {
	dirs := []string{}
	if depth < this.depth {
		for i := uint32(0); i < this.shards; i++ {
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

func (this *FileDB) locateFile(key string) string {
	path := this.findPath(key)
	path += fmt.Sprint(byte(key[this.depth]%byte(this.shards))) + EXTENSIION
	return path
}

func (this *FileDB) findPath(key string) string {
	path := this.rootpath
	for i := uint8(0); i < this.depth; i++ {
		path += fmt.Sprint(byte(key[0]%byte(this.shards))) + "/"
	}
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
	file := this.locateFile(key)
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
	file := this.locateFile(nkey)
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
			if nVal == nil { // Remove the value
				return append(keys[:i], keys[i+1:]...), append(values[:i], values[i+1:]...)
			} else {
				values[i] = nVal // Update the value
				return keys, values
			}
		}
	}
	return append(keys, nKey), append(values, nVal) // Append the new value
}

func (this *FileDB) Set(key string, v []byte) error {
	return this.writeFile(key, v)
}

func (this *FileDB) Get(key string) ([]byte, error) {
	return this.readFile(key)
}

func (this *FileDB) BatchGet(nkeys []string) ([][]byte, error) {
	files := slice.ParallelTransform(nkeys, 8, func(i int, _ string) string {
		return this.locateFile(nkeys[i]) //Must use the compressed ky to compute the shard
	})

	// Read files
	errs := make([]error, len(nkeys))
	data := make([][]byte, len(nkeys))
	t0 := time.Now()
	uniqueFiles, indices := this.CategorizeFiles(files)
	fmt.Println("niqueFiles, indices := this.CategorizeFiles(files):", time.Since(t0))

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
	slice.RemoveIf(&errs, func(_ int, v error) bool { return v == nil })

	if len(errs) > 0 {
		return data, errs[0]
	}

	return data, nil
}

func (this *FileDB) BatchSet(nkeys []string, byteset [][]byte) error {
	files := make([]string, len(nkeys))
	finder := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			files[i] = this.locateFile(nkeys[i])
		}
	}
	common.ParallelWorker(len(nkeys), 8, finder)

	errs := make([]error, len(files))
	uniqueFiles, indices := this.CategorizeFiles(files)

	newFiles := make([]string, len(uniqueFiles))
	maker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			if _, err := os.Stat(uniqueFiles[i]); err != nil { // File doesn't exist, create it
				if err := os.WriteFile(uniqueFiles[i], []byte{}, os.ModePerm); err != nil && errors.Is(err, os.ErrNotExist) {
					errs[i] = err
					return
				}
				newFiles[i] = uniqueFiles[i]
			}
		}
	}
	common.ParallelWorker(len(uniqueFiles), 4, maker)

	slice.Remove(&newFiles, "")
	slice.RemoveIf(&errs, func(_ int, v error) bool { return v == nil })

	this.files = append(this.files, newFiles...)
	if len(errs) > 0 {
		return errs[0].(error)
	}

	// Write Contents
	errs = make([]error, len(uniqueFiles))
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
	slice.RemoveIf(&errs, func(_ int, v error) bool { return v == nil })

	if len(errs) > 0 {
		return errs[0]
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
	paths = slice.Unique(paths, func(s0, s1 string) bool { return s0 < s1 })
	return this.readAll(paths)
}

func (this *FileDB) ExportAll() ([][]byte, error) {
	return this.readAll(this.dirs)
}

func (this *FileDB) readAll(paths []string) ([][]byte, error) {
	sort.Strings(paths)

	errs := make([]error, len(paths))
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
	slice.RemoveIf(&errs, func(_ int, v error) bool { return v == nil })

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

func (this *FileDB) Files() ([]string, error) {
	return this.files, nil
}

func (this *FileDB) getFilesUnder(root string) ([]string, error) {
	list := []string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == EXTENSIION {
			list = append(list, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (this *FileDB) ListFiles() ([]string, error) {
	return this.getFilesUnder(this.rootpath)
}

func (this *FileDB) Import(data [][]byte) error {
	for i := 0; i < len(data); i++ {
		combo := codec.Byteset{}.Decode(data[i]).(codec.Byteset)

		file := codec.Bytes(combo[0]).ToString() // file ROOT_PATH
		// file := this.rootpath + relativePath

		reversed := codec.String(file).Reverse()
		parts := strings.Split(reversed, "/")
		if len(parts) < int(this.depth) {
			return errors.New("Error: Wrong path !!!")
		}
		reversed = strings.Join(parts[:this.depth+1], "/")
		ROOT_PATH := this.rootpath + codec.String(reversed).Reverse()
		if err := os.WriteFile(ROOT_PATH, combo[1], os.ModePerm); err != nil {
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
