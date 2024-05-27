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

// The package offers fast data compression using a lookup table, replacing addresses with index numbers.

package addrcompressor

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	codec "github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	ccmap "github.com/arcology-network/common-lib/container/map"
	slice "github.com/arcology-network/common-lib/exp/slice"
)

type CompressionLut struct {
	IdxToKeyLut []string
	dict        *ccmap.ConcurrentMap
	tempLut     *CompressionLut
	length      uint32
	offset      uint32
	lock        sync.RWMutex
	depths      [][2]int
}

// NewCompressionLut creates a new instance of CompressionLut.
func NewCompressionLut() *CompressionLut {
	syspath := []string{
		"blcc://eth1.0/account",
		"code",
		"nonce",
		"balance",
		"defer",
		//"storage",
		"storage/containers",
		"storage/native",
	}

	tempLut := &CompressionLut{
		IdxToKeyLut: syspath,
		dict:        ccmap.NewConcurrentMap(),
		tempLut:     nil,
		length:      0,
		offset:      0,
	}

	lut := &CompressionLut{
		IdxToKeyLut: []string{},
		length:      0,
		dict:        ccmap.NewConcurrentMap(),
		tempLut: &CompressionLut{
			IdxToKeyLut: []string{},
			dict:        ccmap.NewConcurrentMap(),
			tempLut:     tempLut,
			length:      0,
			offset:      0,
		},

		depths: [][2]int{{-1, 3}, {3, 4}, {4, 8}},
	}

	lut.tempLut.insertToDict(lut.filterExistingKeys(syspath, lut.dict), lut.dict)
	lut.Commit()
	return lut
}

// CompressOnTemp compresses the given slice of strings using the temporary lookup table.
func (this *CompressionLut) CompressOnTemp(originals []string) []string {
	this.tempLut.offset = this.dict.Size()

	t0 := time.Now()
	positions := this.findPositions(originals, this.depths)
	fmt.Println("findPositions", time.Since(t0))

	t0 = time.Now()
	nKeys := slice.Flatten(this.parseKeys(originals, positions))
	fmt.Println("Flatten", time.Since(t0))

	t0 = time.Now()
	this.tempLut.insertToDict(nKeys, this.dict) // update the dictionary
	fmt.Println("insertToDict", time.Since(t0))

	t0 = time.Now()
	originals = this.replaceSubstrings(originals, positions) // Use the dictionary to compress the entires
	fmt.Println("replaceSubstrings", time.Since(t0))

	return originals
}

// TryCompress tries to compress the given string using the lookup table.
func (this *CompressionLut) TryCompress(original string) string {
	positions := this.findPosition(original, this.depths)
	return this.replaceSubstring(original, positions)
}

// TryBatchCompress compresses the given slice of strings using the lookup table.
func (this *CompressionLut) TryBatchCompress(originals []string) []string {
	positions := this.findPositions(originals, this.depths)
	return this.replaceSubstrings(originals, positions)
}

// replaceSubstrings replaces substrings in the given slice of strings based on the positions provided.
func (this *CompressionLut) replaceSubstrings(originals []string, positions [][][2]int) []string {
	//this.insertToDict(Flatten(this.parseKeys(originals, positions)))
	compressor := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			originals[i] = this.replaceSubstring(originals[i], positions[i])
		}
	}
	common.ParallelWorker(len(originals), 4, compressor)
	return originals
}

// TryUncompress tries to uncompress the given compressed string using the lookup table.
func (this *CompressionLut) TryUncompress(compressed string) string {
	if strings.Count(compressed, "[") == 0 {
		return compressed
	}

	p0 := make([]uint32, 0, 6)
	p1 := make([]uint32, 0, 6)
	for j := range compressed {
		if compressed[j] == '[' {
			p0 = append(p0, uint32(j))
		}

		if compressed[j] == ']' {
			p1 = append(p1, uint32(j))
		}
	}

	if len(p0) == 0 || len(p1) == 0 {
		return compressed
	}

	var buffer bytes.Buffer
	prefix := compressed[:p0[0]]
	buffer.WriteString(prefix)
	for i := 0; i < len(p0); i++ {
		idxStr := compressed[p0[i]+1 : p1[i]]
		idx, _ := strconv.Atoi(idxStr)
		if idx >= len(this.IdxToKeyLut) {
			panic("Error: Wrong Uncompression LUT")
		}
		buffer.WriteString(this.IdxToKeyLut[idx])

		if i+1 < len(p0) {
			buffer.WriteString(compressed[p1[i]+1 : p0[i+1]])
		} else {
			buffer.WriteString(compressed[p1[i]+1:])
			break
		}
	}
	return buffer.String()
}

// TryBatchUncompress uncompresses the given slice of compressed strings using the lookup table.
func (this *CompressionLut) TryBatchUncompress(compressed []string) {
	if len(compressed) < 1024 {
		this.singleThreadedUncompressor(compressed)
	} else {
		this.multiThreadedUncompressor(compressed)
	}
}

// GetCompressionRatio calculates the compression ratio between the original and compressed strings.
func (this *CompressionLut) GetCompressionRatio(originals []string, compressed []string) float32 {
	compressedLen := 0
	for _, v := range compressed {
		compressedLen += len(v)
	}

	originalLen := len(codec.Strings(originals).Flatten())
	return float32(compressedLen) / float32(originalLen)
}
