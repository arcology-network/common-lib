package codec

const (
	BOOL_LEN = 1
)

type Bool bool

func (this *Bool) Get() interface{} {
	return *this
}

func (this *Bool) Set(v interface{}) {
	*this = v.(Bool)
}

func (this Bool) Size() int {
	return BOOL_LEN
}

func (this Bool) Encode() []byte {
	data := make([]byte, BOOL_LEN)
	if this {
		data[0] = 1
	} else {
		data[0] = 0
	}
	return data
}

func (Bool) Decode(data []byte) Bool {
	return Bool(data[0] > 0)
}

type Bools []bool

func (this Bools) Size() int {
	return len(this)
}

func (this Bools) Encode() []byte {
	data := make([]byte, len(this))
	for i := range this {
		if this[i] {
			data[i] = 1
		} else {
			data[i] = 0
		}
	}
	return data
}

func (Bools) Decode(data []byte) interface{} {
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
