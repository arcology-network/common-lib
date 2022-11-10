package datacompression

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	codec "github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	cccontainer "github.com/arcology-network/common-lib/concurrentcontainer/map"
)

type CompressionLut struct {
	IdxToKeyLut []string
	dict        *cccontainer.ConcurrentMap
	tempLut     *CompressionLut
	length      uint32
	offset      uint32
	lock        sync.RWMutex
	depths      [][2]int
}

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
		dict:        cccontainer.NewConcurrentMap(),
		tempLut:     nil,
		length:      0,
		offset:      0,
	}

	lut := &CompressionLut{
		IdxToKeyLut: []string{},
		length:      0,
		dict:        cccontainer.NewConcurrentMap(),
		tempLut: &CompressionLut{
			IdxToKeyLut: []string{},
			dict:        cccontainer.NewConcurrentMap(),
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

func (this *CompressionLut) CompressOnTemp(originals []string) []string {
	this.tempLut.offset = this.dict.Size()

	t0 := time.Now()
	positions := this.findPositions(originals, this.depths)
	fmt.Println("findPositions", time.Since(t0))

	t0 = time.Now()
	nKeys := Flatten(this.parseKeys(originals, positions))
	fmt.Println("Flatten", time.Since(t0))

	t0 = time.Now()
	this.tempLut.insertToDict(nKeys, this.dict) // update the dictionary
	fmt.Println("insertToDict", time.Since(t0))

	t0 = time.Now()
	originals = this.replaceSubstrings(originals, positions) // Use the dictionary to compress the entires
	fmt.Println("replaceSubstrings", time.Since(t0))

	return originals
}

func (this *CompressionLut) TryCompress(original string) string {
	positions := this.findPosition(original, this.depths)
	return this.replaceSubstring(original, positions)
}

func (this *CompressionLut) TryBatchCompress(originals []string) []string {
	positions := this.findPositions(originals, this.depths)
	return this.replaceSubstrings(originals, positions)
}

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

func (this *CompressionLut) TryBatchUncompress(compressed []string) {
	if len(compressed) < 1024 {
		this.singleThreadedUncompressor(compressed)
	} else {
		this.multiThreadedUncompressor(compressed)
	}
}

func (this *CompressionLut) GetCompressionRatio(originals []string, compressed []string) float32 {
	compressedLen := 0
	for _, v := range compressed {
		compressedLen += len(v)
	}

	originalLen := len(codec.Strings(originals).Flatten())
	return float32(compressedLen) / float32(originalLen)
}
