package codec

type Encodable interface {
	Clone() interface{}
	Size() uint32
	Encode() []byte
	EncodeToBuffer([]byte) int
	Decode([]byte) interface{}
}

func ToEncodable() {

}

type Encodables []Encodable

func (this Encodables) Size() uint32 {
	length := uint32(0)
	for i := 0; i < len(this); i++ {
		if this[i] != nil {
			length += this[i].Size()
		}
	}
	return UINT32_LEN*uint32(len(this)+1) + uint32(length)
}

func (this Encodables) Sizes() []uint32 {
	lengths := make([]uint32, len(this))
	for i := 0; i < len(lengths); i++ {
		if this[i] != nil {
			lengths[i] += this[i].Size()
		}
	}
	return lengths
}

func (this Encodables) FillHeader(buffer []byte) int {
	lengths := this.Sizes()
	Uint32(len(lengths)).EncodeToBuffer(buffer[UINT32_LEN*0:])
	offset := uint32(0)
	for i := 0; i < len(lengths); i++ {
		Uint32(offset).EncodeToBuffer(buffer[UINT32_LEN*(i+1):])
		offset += uint32(lengths[i])
	}
	return (len(lengths) + 1) * UINT32_LEN
}

func (this Encodables) Encode() []byte {
	total := this.Size()
	buffer := make([]byte, total)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Encodables) EncodeToBuffer(buffer []byte) int {
	offset := this.FillHeader(buffer)
	for i := 0; i < len(this); i++ {
		// if selectors[i] {
		offset += this[i].EncodeToBuffer(buffer[offset:])
		// }
	}
	return offset
}

func (this Encodables) Decode(buffer []byte, decoders ...func([]byte) interface{}) []interface{} {
	fields := Byteset{}.Decode(buffer).(Byteset)
	values := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		values[i] = decoders[i](fields[i])
	}
	return values
}
