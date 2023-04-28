package codec

type Encodeable interface {
	TypeID() uint8
	Size() uint32
	Encode() []byte
	EncodeToBuffer([]byte) int
	Decode([]byte) interface{}
}

type Encodeables []Encodeable

func (this Encodeables) Size() uint32 {
	total := uint32(0)
	for i := 0; i < len(this); i++ {
		total += this[i].Size()
	}
	return total
}
