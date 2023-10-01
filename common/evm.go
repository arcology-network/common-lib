package common

const (
	EvmWordSize = 32
)

func AlignToEvmForString(str string) []byte {
	strLength := len(str)
	EvmWordBytes := strLength / EvmWordSize
	if strLength%EvmWordSize != 0 {
		EvmWordBytes = EvmWordBytes + 1
	}
	finalLengths := make([]byte, EvmWordBytes*EvmWordSize)
	for i := 0; i < len(finalLengths); i++ {
		if i < strLength {
			finalLengths[i] = str[i]
		} else {
			finalLengths[i] = byte(0)
		}
	}
	return finalLengths
}

func AlignToEvmForInt(length int) []byte {
	lens := []byte{}
	for {
		by := length % 256
		lens = append(lens, byte(by))
		length = length >> 8
		if length == 0 {
			break
		}
	}
	bysLength := len(lens)
	revertLengths := make([]byte, bysLength)
	idx := 0
	for i := bysLength - 1; i >= 0; i-- {
		revertLengths[idx] = lens[i]
		idx++
	}

	EvmWordBytes := bysLength / EvmWordSize
	if bysLength%EvmWordSize != 0 {
		EvmWordBytes = EvmWordBytes + 1
	}

	finalLengths := make([]byte, EvmWordBytes*EvmWordSize)
	idx = 0
	for ; idx < len(finalLengths)-bysLength; idx++ {
		finalLengths[idx] = byte(0)
	}
	for j := 0; j < bysLength; j++ {
		finalLengths[idx] = revertLengths[j]
		idx++
	}

	return finalLengths
}
