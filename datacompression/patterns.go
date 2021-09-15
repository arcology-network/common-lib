package datacompression

import (
	"bytes"
	"strconv"
)

// func (this *CompressionLut) compressBuiltins(original string, patterns []string, idStart int) string {

// }

func (this *CompressionLut) compressPattens(originals []string) []string {
	compressed := make([]string, len(originals))
	//worker := func(start, end, idx int, args ...interface{}) {
	for i := 0; i < len(originals); i++ {
		line := this.compressFixedLeading(originals[i], this.KeyLut[:this.offset], 0) // Partially compressed
		compressed[i] = this.compressFixedTrailing(line, this.KeyLut[:this.offset], 0)
	}
	//}
	//common.ParallelWorker(len(originals), 6, worker)
	return compressed
}

func (this *CompressionLut) compressFixedLeading(original string, patterns []string, idStart int) string {
	compressedLine := original
	for i := 0; i < len(patterns); i++ {
		if len(original) >= len(patterns[i]) {
			if bytes.Equal([]byte(patterns[i]), []byte(original[:len(patterns[i])])) {
				idStr := strconv.Itoa(i + idStart)
				compressedLine = "[" + idStr + "]" + original[len(patterns[i]):]
			}
		}
	}
	return compressedLine
}

func (this *CompressionLut) compressFixedTrailing(original string, patterns []string, idStart int) string {
	compressedLine := original
	for i := 0; i < len(patterns); i++ {
		if len(original) >= len(patterns[i]) {
			lhs := original[len(original)-len(patterns[i]):]
			if bytes.Equal([]byte(patterns[i]), []byte(lhs)) {
				idStr := strconv.Itoa(i + idStart)
				compressedLine = original[:len(original)-len(lhs)] + "[" + idStr + "]"
			}
		}
	}
	return compressedLine
}
