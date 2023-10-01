package cachedstorage

import (
	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
)

func (this *FileDB) Query(pattern string, condition func(string, string) bool) ([]string, [][]byte, error) {
	parentPath := this.findPath(pattern) // match file parent path first
	if files, err := this.getFilesUnder(parentPath); err == nil {
		keyset := make([][]string, len(files))
		valSet := make([][][]byte, len(files))

		for i := 0; i < len(files); i++ {
			keys, valBytes, err := this.loadFile(files[i])
			if err != nil {
				return []string{}, [][]byte{}, err
			}

			for j := 0; j < len(keys); j++ {
				if !condition(pattern, keys[j]) {
					keys[j] = ""
					valBytes[j] = valBytes[j][:0]
				}
			}

			common.Remove(&keys, "")
			common.RemoveIf(&valBytes, func(v []byte) bool { return len(v) == 0 })

			keyset[i] = keys
			valSet[i] = valBytes
		}
		return codec.Stringset(keyset).Flatten(), codec.Bytegroup(valSet).Flatten(), nil
	} else {
		return []string{}, [][]byte{}, err
	}
}
