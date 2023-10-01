package encoding

import (
	"crypto/sha256"

	evmCommon "github.com/arcology-network/evm/common"
)

const (
	BYTE_LEN = 1
)

type Byteset [][]byte

func (byteset Byteset) Sizes() []uint32 {
	sizes := make([]uint32, len(byteset))
	for i := range byteset {
		sizes[i] = uint32(len(byteset[i]))
	}
	return sizes
}

func (byteset Byteset) Flatten() []byte {
	lengths := byteset.Sizes()
	buffer := make([]byte, Uint32s(lengths).Sum())

	positions := append([]uint32{0}, Uint32s(lengths).Accumulate()...)
	for i := 0; i < len(positions)-1; i++ {
		copy(buffer[positions[i]:positions[i+1]], byteset[i])
	}
	return buffer
}

func (byteset Byteset) Checksum() evmCommon.Hash {
	return sha256.Sum256(byteset.Flatten())
}

func (byteset Byteset) Encode() []byte {
	lengths := byteset.Sizes()
	header := append(Uint32(len(lengths)).Encode(), Uint32s(lengths).Encode()...)
	headerLen := uint32(len(header))

	buffer := make([]byte, headerLen+Uint32s(lengths).Sum())
	copy(buffer[:len(header)], header)

	positions := append([]uint32{0}, Uint32s(lengths).Accumulate()...)
	for i := 0; i < len(positions)-1; i++ {
		copy(buffer[positions[i]+headerLen:positions[i+1]+headerLen], byteset[i])
	}
	return buffer
}

func (Byteset) Decode(bytes []byte) [][]byte {
	count := Uint32(0).Decode(bytes[:UINT32_LEN])
	lengths := Uint32s{}.Decode(bytes[UINT32_LEN : UINT32_LEN+count*UINT32_LEN])

	byteset := make([][]byte, count)
	positions := append([]uint32{0}, Uint32s(lengths).Accumulate()...)
	for i := 0; i < len(positions)-1; i++ {
		byteset[i] = bytes[(count+1)*UINT32_LEN+positions[i] : (count+1)*UINT32_LEN+positions[i+1]]
	}
	return byteset
}

type Bytegroup [][][]byte

func (bytegroup Bytegroup) Sizes() []uint32 {
	sizes := make([]uint32, len(bytegroup))
	for i := range bytegroup {
		sizes[i] = uint32(len(bytegroup[i]))
	}
	return sizes
}

func (bytegroup Bytegroup) Flatten() [][]byte {
	lengths := bytegroup.Sizes()
	buffer := make([][]byte, Uint32s(lengths).Sum())

	positions := append([]uint32{0}, Uint32s(lengths).Accumulate()...)
	for i := 0; i < len(positions)-1; i++ {
		copy(buffer[positions[i]:positions[i+1]], bytegroup[i])
	}
	return buffer
}
