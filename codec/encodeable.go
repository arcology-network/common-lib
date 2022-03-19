package codec

type Encodeable interface {
	Size() uint32
	EncodeToBuffer([]byte)
	Decode([]byte) interface{}
}

type Encoder struct{}

func (Encoder) Size(args []interface{}) uint32 {
	length := uint32(0)
	for i := 0; i < len(args); i++ {
		if args[i] != nil {
			length += args[i].(Encodeable).Size()
		}
	}
	return UINT32_LEN*uint32(len(args)+1) + uint32(length)
}

func (this Encoder) ToBuffer(buffer []byte, args []interface{}) {
	offset := uint32(0)
	Uint32(len(args)).EncodeToBuffer(buffer)
	for i := 0; i < len(args); i++ {
		Uint32(offset).EncodeToBuffer(buffer[(i+1)*UINT32_LEN:]) // Fill header info
		if args[i] != nil {
			offset += args[i].(Encodeable).Size()
		}
	}
	headerSize := uint32((len(args) + 1) * UINT32_LEN)

	offset = uint32(0)
	for i := 0; i < len(args); i++ {
		if args[i] != nil {
			end := headerSize + offset + args[i].(Encodeable).Size()
			args[i].(Encodeable).EncodeToBuffer(buffer[headerSize+offset : end])
			offset += args[i].(Encodeable).Size()
		}
	}
}
