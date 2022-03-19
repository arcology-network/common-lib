package encoding

const (
	BOOL_LEN = 1
)

type Bools []bool

func (v Bools) Size() uint32 {
	return uint32(len(v))
}

func (v Bools) Encode() []byte {
	data := make([]byte, len(v))
	for i := range v {
		if v[i] {
			data[i] = 1
		} else {
			data[i] = 0
		}
	}
	return data
}

func (Bools) Decode(data []byte) []bool {
	bools := make([]bool, len(data))
	for i := range data {
		if data[i] == 1 {
			bools[i] = true
		} else {
			bools[i] = false
		}
	}
	return bools
}
