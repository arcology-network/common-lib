package codec

import (
	"math/big"
)

type Bigint big.Int

func (this *Bigint) Size() uint32 {
	return BOOL_LEN + uint32((*big.Int)(this).BitLen())
}

func (this *Bigint) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this *Bigint) EncodeToBuffer(buffer []byte) {
	flag := (*big.Int)(this).Cmp(big.NewInt(0)) >= 0

	Bool(flag).EncodeToBuffer(buffer)
	val := (big.Int)(*this)
	val.FillBytes(buffer[1 : val.BitLen()+1])
}

func (this *Bigint) Decode(buffer []byte) interface{} {
	if len(buffer) > 0 {
		v := new(big.Int)
		*this = *(*Bigint)(v.SetBytes(buffer[1:]))
		if !Bool(true).Decode(buffer[:1]).(Bool) { // negative value
			return (*Bigint)((&big.Int{}).Neg(v))
		}
	}
	return (*Bigint)(this)
}
