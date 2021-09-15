package datacompression

import (
	"strconv"
	"strings"

	"github.com/arcology-network/common-lib/common"
)

type CompressionLut struct {
	KeyLut []string
	Dict   map[string]uint32
	offset uint32 // offset for builtins
}

func NewCompressionLut() *CompressionLut {
	lut := &CompressionLut{
		KeyLut: []string{
			"blcc://eth1.0/account",
			// "code",
			// "nonce",
			// "balance",
			// "defer/",
			// "storage/",
			"storage/containers",
			"storage/native",
			"storage/containers/!",
		},
		Dict: map[string]uint32{},
	}
	lut.offset = uint32(len(lut.KeyLut))
	return lut
}

func (this *CompressionLut) Compress(original string) string {
	line := this.compressFixedLeading(original, this.KeyLut[:this.offset], 0) // Partially compressed
	if first := strings.Index(line, "]/"); first > -1 {
		first += 2
		if second := strings.Index(line[first:], "/"); second > -1 {
			second += first
			key := line[first:second]
			if _, ok := this.Dict[line[first:second]]; !ok {
				this.Dict[key] = uint32(len(this.Dict) + len(this.KeyLut))
				this.KeyLut = append(this.KeyLut, key)
			}
			line = line[:first] + "[" + strconv.Itoa(int(this.Dict[key])) + "]" + line[second:]
		}
	}

	return this.compressFixedTrailing(line, this.KeyLut[:this.offset], 0)
}

func (this *CompressionLut) Uncompress(compressed string) string {
	pos := make([]uint32, 0, 6)
	for j := range compressed {
		if compressed[j] == '[' || compressed[j] == ']' {
			pos = append(pos, uint32(j))
		}
	}

	if len(pos) > 0 {
		original := compressed[:pos[0]]
		for j := 0; j < len(pos)/2; j++ {
			idxStr := compressed[pos[j*2]+1 : pos[j*2+1]]
			idx, _ := strconv.Atoi(idxStr)
			original += this.KeyLut[idx]

			if (j+1)*2 < len(pos) {
				original += compressed[pos[j*2+1]+1 : pos[(j+1)*2]]
			}
		}

		if uint32(len(compressed)) > pos[len(pos)-1]+1 {
			original += compressed[pos[len(pos)-1]+1:]
		}
		return original
	}
	return compressed
}

func (this *CompressionLut) BatchCompress(originals []string) []string {
	worker := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			originals[i] = this.Compress(originals[i])
		}
	}
	common.ParallelWorker(len(originals), 6, worker)
	return originals
}

func (this *CompressionLut) BatchUncompress(compressed []string) []string {
	worker := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			compressed[i] = this.Uncompress(compressed[i])
			compressed[i] = this.Uncompress(compressed[i]) // Twice
		}
	}
	common.ParallelWorker(len(compressed), 6, worker)
	return compressed
}
