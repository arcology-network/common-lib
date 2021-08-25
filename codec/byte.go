package codec

import (
	"crypto/sha256"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
)

const (
	BYTE_LEN = 1
)

type Bytes []byte

func (this *Bytes) Get() interface{} {
	return *this
}

func (this *Bytes) Set(v interface{}) {
	*this = v.(Bytes)
}

func (bytes Bytes) Encode() []byte {
	return []byte(bytes)
}

func (Bytes) Decode(bytes []byte) interface{} {
	return Bytes(bytes)
}

type Byteset [][]byte

func (byteset Byteset) Sizes() Uint32s {
	sizes := make([]Uint32, len(byteset))
	for i := range byteset {
		sizes[i] = Uint32(len(byteset[i]))
	}
	return sizes
}

func (byteset Byteset) Flatten() []byte {
	lengths := byteset.Sizes()
	buffer := make([]byte, Uint32s(lengths).Sum())

	positions := append([]Uint32{0}, Uint32s(lengths).Accumulate()...)
	for i := 0; i < len(positions)-1; i++ {
		copy(buffer[positions[i]:positions[i+1]], byteset[i])
	}
	return buffer
}

func (byteset Byteset) Checksum() ethCommon.Hash {
	return sha256.Sum256(byteset.Flatten())
}

func (byteset Byteset) Encode() []byte {
	lengths := byteset.Sizes()
	header := append(Uint32(len(lengths)).Encode(), Uint32s(lengths).Encode()...)
	headerLen := uint32(len(header))

	buffer := make([]byte, headerLen+Uint32s(lengths).Sum())
	copy(buffer[:len(header)], header)

	positions := append([]Uint32{0}, Uint32s(lengths).Accumulate()...)
	for i := 0; i < len(positions)-1; i++ {
		copy(buffer[uint32(positions[i])+headerLen:uint32(positions[i+1])+headerLen], byteset[i])
	}
	return buffer
}

func (Byteset) Decode(bytes []byte) [][]byte {
	count := uint32(Uint32(0).Decode(bytes[:UINT32_LEN]))
	lengths := Uint32s{}.Decode(bytes[UINT32_LEN : UINT32_LEN+count*UINT32_LEN]).(Uint32s)

	byteset := make([][]byte, count)
	positions := append([]Uint32{0}, (lengths).Accumulate()...)
	for i := 0; i < len(positions)-1; i++ {
		start := (count+1)*(UINT32_LEN) + uint32(positions[i])
		end := (count+1)*(UINT32_LEN) + uint32(positions[i+1])
		byteset[i] = bytes[start:end]
	}
	return byteset
}

type Bytegroup [][][]byte

func (bytegroup Bytegroup) Sizes() []Uint32 {
	sizes := make([]Uint32, len(bytegroup))
	for i := range bytegroup {
		sizes[i] = Uint32(len(bytegroup[i]))
	}
	return sizes
}

func (bytegroup Bytegroup) Flatten() [][]byte {
	lengths := bytegroup.Sizes()
	buffer := make([][]byte, Uint32s(lengths).Sum())

	positions := append([]Uint32{0}, Uint32s(lengths).Accumulate()...)
	for i := 0; i < len(positions)-1; i++ {
		copy(buffer[positions[i]:positions[i+1]], bytegroup[i])
	}
	return buffer
}
