package codec

import (
	"math/big"
)

type Bigint big.Int

func (this *Bigint) Encode() []byte {
	flag := (*big.Int)(this).Cmp(big.NewInt(0)) >= 0
	return Byteset{
		Bool(flag).Encode(),
		(*big.Int)(this).Bytes(),
	}.Encode()
}

func (*Bigint) Decode(data []byte) *big.Int {
	fields := Byteset{}.Decode(data)
	v := new(big.Int)
	v = v.SetBytes(fields[1])
	if !bool(Bool(true).Decode(fields[0])) { // negative value
		return (&big.Int{}).Neg(v)
	}
	return v
}
