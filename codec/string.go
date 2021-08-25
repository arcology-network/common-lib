package codec

import "github.com/HPISTechnologies/common-lib/common"

const (
	CHAR_LEN = 1
)

type String string

func (str String) Size() int {
	return CHAR_LEN * len(str)
}

func (str String) Encode() []byte {
	return []byte(str)
}

func (String) Decode(bytes []byte) string {
	return string(bytes)
}

type Strings []string

func (strs Strings) Encode() []byte {
	byteset := make([][]byte, len(strs))
	for i := range byteset {
		byteset[i] = []byte(strs[i])
	}
	return Byteset(byteset).Encode()
}

func (Strings) Decode(bytes []byte) Strings {
	fields := Byteset{}.Decode(bytes)
	strs := make([]string, len(fields))
	for i := range fields {
		strs[i] = string(fields[i])
	}
	return strs
}

func (this Strings) Flatten() []byte {
	positions := make([]int, len(this)+1)
	positions[0] = 0
	for i := 1; i < len(positions); i++ {
		positions[i] = positions[i-1] + len(this[i-1])
	}

	buffer := make([]byte, positions[len(positions)-1])
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			copy(buffer[positions[i]:positions[i+1]], []byte(this[i]))
		}
	}
	common.ParallelWorker(len(this), 4, worker)
	return buffer
}
