package codec

type Encodeable interface {
	Size() uint32
	Encode() []byte
	EncodeToBuffer([]byte) int
	Decode([]byte) interface{}
}

type Encodeables []Encodeable

func (this Encodeables) Size() uint32 {
	length := uint32(0)
	for i := 0; i < len(this); i++ {
		if this[i] != nil {
			length += this[i].Size()
		}
	}
	return UINT32_LEN*uint32(len(this)+1) + uint32(length)
}

func (this Encodeables) Sizes() []uint32 {
	lengths := make([]uint32, len(this))
	for i := 0; i < len(lengths); i++ {
		if this[i] != nil {
			lengths[i] += this[i].Size()
		}
	}
	return lengths
}

func (this Encodeables) FillHeader(buffer []byte) int {
	lengths := this.Sizes()
	Uint32(len(lengths)).EncodeToBuffer(buffer[UINT32_LEN*0:])
	offset := uint32(0)
	for i := 0; i < len(lengths); i++ {
		Uint32(offset).EncodeToBuffer(buffer[UINT32_LEN*(i+1):])
		offset += uint32(lengths[i])
	}
	return (len(lengths) + 1) * UINT32_LEN
}

func (this Encodeables) Encode() []byte {
	total := this.Size()
	buffer := make([]byte, total)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Encodeables) EncodeToBuffer(buffer []byte) int {
	offset := this.FillHeader(buffer)
	for i := 0; i < len(this); i++ {
		// if selectors[i] {
		offset += this[i].EncodeToBuffer(buffer[offset:])
		// }
	}
	return offset
}

func (this Encodeables) Decode(buffer []byte, decoders ...func([]byte) interface{}) []interface{} {
	fields := Byteset{}.Decode(buffer).(Byteset)
	values := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		values[i] = decoders[i](fields[i])
	}
	return values
}
