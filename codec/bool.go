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

func (this Bool) Size() uint32 {
	return uint32(BOOL_LEN)
}

func (this Bool) Encode() []byte {
	buffer := make([]byte, BOOL_LEN)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Bool) EncodeToBuffer(buffer []byte) {
	if this {
		buffer[0] = 1
	} else {
		buffer[0] = 0
	}
}

func (this Bool) Decode(data []byte) interface{} {
	this = Bool(data[0] > 0)
	return this
}

type Bools []bool

func (this Bools) Size() int {
	return len(this)
}

func (this Bools) Encode() []byte {
	buffer := make([]byte, len(this))
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Bools) EncodeToBuffer(buffer []byte) {
	for i := range this {
		if this[i] {
			buffer[i] = 1
		} else {
			buffer[i] = 0
		}
	}
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
	return Bools(bools)
}
