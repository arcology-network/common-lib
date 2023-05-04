package codec

import common "github.com/arcology-network/common-lib/common"

const (
	BOOL_LEN = 1
)

type Bool bool

func (this *Bool) Clone() interface{} {
	return common.New(*this)
}

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

func (this Bool) EncodeToBuffer(buffer []byte) int {
	buffer[0] = uint8(common.IfThen(bool(this), 1, 0))
	return BOOL_LEN
}

func (this Bool) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}

	this = Bool(buffer[0] > 0)
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

func (this Bools) EncodeToBuffer(buffer []byte) int {
	for i := range this {
		if this[i] {
			buffer[i] = 1
		} else {
			buffer[i] = 0
		}
	}
	return len(this) * BOOL_LEN
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
