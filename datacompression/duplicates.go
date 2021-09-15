package datacompression

func (this *CompressionLut) compressDuplicates(originals []string) []string {
	idx := 0
	minLen := len(originals[idx])
	for i := 0; i < len(originals); i++ {
		if len(originals[i]) > 0 && minLen > len(originals[i]) && originals[i][len(originals[i])-1] != ']' {
			idx = i
			minLen = len(originals[idx])
		}
	}

	this.KeyLut = append(this.KeyLut, originals[0])
	for i := 1; i < len(originals); i++ {
		originals[i] = this.compressFixedLeading(originals[i], originals[:1], int(this.offset))

	}
	return originals
}
