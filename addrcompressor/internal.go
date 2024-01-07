package addrcompressor

import (
	"bytes"
	"strconv"

	common "github.com/arcology-network/common-lib/common"
	ccmap "github.com/arcology-network/common-lib/container/map"
)

func (this *CompressionLut) singleThreadedUncompressor(compressed []string) {
	for i := 0; i < len(compressed); i++ {
		compressed[i] = this.TryUncompress(compressed[i])
	}
}

func (this *CompressionLut) multiThreadedUncompressor(compressed []string) {
	// worker := func(start, end, idx int, args ...interface{}) {
	// 	for i := start; i < end; i++ {
	// 		compressed[i] = this.TryUncompress(compressed[i])
	// 	}
	// }
	// common.ParallelWorker(len(compressed), 4, worker)
	common.ParallelForeach(compressed, 6, func(i int, _ *string) {
		compressed[i] = this.TryUncompress(compressed[i])
	})
}

func (this *CompressionLut) findPositions(originals []string, depths [][2]int) [][][2]int {
	positions := make([][][2]int, len(originals))

	common.ParallelForeach(originals, 6, func(i int, _ *string) {
		positions[i] = make([][2]int, len(depths))
		for j := 0; j < len(depths); j++ {
			positions[i][j][0] = IndexN(originals[i], "/", depths[j][0])
			positions[i][j][1] = IndexN(originals[i], "/", depths[j][1])
		}
	})
	return positions
}

func (this *CompressionLut) findPosition(original string, depths [][2]int) [][2]int {
	positions := make([][2]int, len(depths))
	for i := 0; i < len(depths); i++ {
		positions[i][0] = IndexN(original, "/", depths[i][0])
		positions[i][1] = IndexN(original, "/", depths[i][1])
	}
	return positions
}

func (this *CompressionLut) parseKeys(originals []string, positions [][][2]int) [][]string {
	keySet := make([][]string, len(originals))
	worker := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			for j := 0; j < len(positions[i]); j++ {
				if positions[i][j][0] < 0 || positions[i][j][1] < 0 {
					continue
				}

				if positions[i][j][0] == positions[i][j][1] {
					continue
				}

				if positions[i][j][0] == 0 {
					keySet[i] = append(keySet[i], originals[i][positions[i][j][0]:positions[i][j][1]])
				} else {
					p0 := positions[i][j][0] + 1
					p1 := positions[i][j][1]
					keySet[i] = append(keySet[i], originals[i][p0:p1])
				}

			}
		}
	}
	common.ParallelWorker(len(originals), 4, worker)
	return keySet
}

func (this *CompressionLut) filterExistingKeys(keys []string, dict *ccmap.ConcurrentMap) []string {
	values := dict.BatchGet(keys)
	nKeys := make([]string, 0, len(values))
	for i := range values {
		if values[i] == nil && len(keys[i]) > 0 {
			nKeys = append(nKeys, keys[i])
		}
	}
	return nKeys
}

func (this *CompressionLut) replaceSubstring(original string, pos [][2]int) string {
	var buffer bytes.Buffer
	prefix := original[:pos[0][0]]
	buffer.WriteString(prefix)
	for i := 0; i < len(pos); i++ {
		if pos[i][0] < 0 {
			return original
		}

		if pos[i][1] < 0 {
			connection := original[pos[i][0]+1:]
			this.searchInDict(connection, &buffer)
			return buffer.String()
		}

		if pos[i][0] == pos[i][1] {
			continue
		}

		var key string
		if pos[i][0] == 0 {
			key = original[pos[i][0]:pos[i][1]]
		} else {
			key = original[pos[i][0]+1 : pos[i][1]]
		}
		this.searchInDict(key, &buffer)

		if pos[i][1]+1 >= len(original) {
			key = original[pos[i][1]:]
			buffer.WriteString(key)
			break
		}

		if i+1 >= len(pos) {
			connection := original[pos[i][1]:]
			buffer.WriteString(connection)
			break
		}
		connection := original[pos[i][1] : pos[i+1][0]+1]
		buffer.WriteString(connection)
	}
	return buffer.String()
}

func (this *CompressionLut) searchInDict(key string, buffer *bytes.Buffer) {
	if v, ok := this.dict.Get(key); ok {
		buffer.WriteString("[" + strconv.Itoa(int(v.(uint32)+this.offset)) + "]")
	} else {
		if v, ok := this.tempLut.dict.Get(key); ok {
			buffer.WriteString("[" + strconv.Itoa(int(v.(uint32)+this.tempLut.offset)) + "]")
		} else {
			buffer.WriteString(key)
		}
	}
}

func (this *CompressionLut) insertToDict(newKeys []string, dict *ccmap.ConcurrentMap) {
	this.lock.Lock()
	defer this.lock.Unlock()

	uniqueKeys := common.Unique(newKeys, func(str0, str1 string) bool { return str0 < str1 })
	newKeys = this.filterExistingKeys(uniqueKeys, dict)
	if len(newKeys) == 0 {
		return
	}

	offset := this.dict.Size()
	indices := make([]interface{}, len(newKeys))
	for i := uint32(0); i < uint32(len(newKeys)); i++ {
		indices[i] = i + offset
	}
	this.dict.BatchSet(newKeys, indices)
}

func (this *CompressionLut) merge(nKeys []string, nValues []interface{}) {
	length := this.dict.Size()
	for i := range nValues {
		nValues[i] = nValues[i].(uint32) + length
	}
	this.dict.BatchSet(nKeys, nValues)
}

func (this *CompressionLut) Commit() {
	nKeys := this.tempLut.dict.Keys() //SHOULD NOT INCLUDE ALL THE KEYS
	nValues := this.tempLut.dict.BatchGet(nKeys)

	length := len(this.IdxToKeyLut)
	this.IdxToKeyLut = append(this.IdxToKeyLut, make([]string, this.tempLut.dict.Size())...)
	for i := 0; i < len(nValues); i++ {
		this.IdxToKeyLut[length+int(nValues[i].(uint32))] = nKeys[i]
	}

	this.merge(nKeys, nValues)
	this.length = this.dict.Size()
	this.reset()
}

func (this *CompressionLut) reset() {
	this.tempLut.IdxToKeyLut = this.tempLut.IdxToKeyLut[:0]
	this.tempLut.dict = ccmap.NewConcurrentMap()
	this.tempLut.offset = this.dict.Size()
}
