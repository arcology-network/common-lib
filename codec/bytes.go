package codec

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	BYTE_LEN = 1
)

type Bytes []byte

func (*Bytes) LessAsUint64(first, second []byte) bool {
	return *(*uint64)(unsafe.Pointer((*[8]byte)(unsafe.Pointer(&first)))) <
		*(*uint64)(unsafe.Pointer((*[8]byte)(unsafe.Pointer(&second))))
}

func (this *Bytes) Get() interface{} {
	return *this
}

func (this *Bytes) Set(v interface{}) {
	*this = v.(Bytes)
}

func (this *Bytes) Sum(offset uint64) uint64 {
	total := uint64(0)
	for j := offset; j < uint64(len(*this)); j++ {
		total += uint64((*this)[j])
	}
	return total
}

func (this *Bytes) Hex() string {
	bytes := make([]byte, 2*len(*this))
	hex.Encode(bytes[:], (*this)[:])
	return string(bytes)
}

func (this Bytes) Encode() []byte {
	return []byte(this)
}

func (this Bytes) Size() uint32 {
	return uint32(len(this))
}

func (this Bytes) Clone() interface{} {
	if this == nil {
		return this
	}

	target := make([]byte, len(this))
	copy(target, this)
	return Bytes(target)
}

func (this Bytes) EncodeToBuffer(buffer []byte) int {
	copy(buffer, this)
	return len(this)
}

func (this Bytes) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return this
	}
	return Bytes(buffer)
}

func (this Bytes) ToString() string {
	return *(*string)(unsafe.Pointer(&this))
}

type Byteset [][]byte

func (this Byteset) Clone() interface{} {
	if this == nil {
		return this
	}

	target := make([][]byte, len(this))
	for i := range this {
		target[i] = make([]byte, len(this[i]))
		copy(target[i], this[i])
	}
	return Byteset(target)
}

func (this Byteset) Size() uint32 {
	if len(this) == 0 {
		return 0
	}

	total := (len(this) + 1) * UINT32_LEN // Header size
	for i := 0; i < len(this); i++ {
		total += len(this[i])
	}
	return uint32(total)
}

func (this Byteset) Sizes() Uint32s {
	sizes := make([]uint32, len(this))
	for i := range this {
		sizes[i] = uint32(len(this[i]))
	}
	return sizes
}

func (this Byteset) Flatten() []byte {
	total := 0
	for i := range this {
		total += len(this[i])
	}
	buffer := make([]byte, total)

	offset := 0
	for i := 0; i < len(this); i++ {
		copy(buffer[offset:], this[i])
		offset += len(this[i])
	}
	return buffer
}

func (this Byteset) Checksum() ethCommon.Hash {
	return sha256.Sum256(this.Flatten())
}

func (this Byteset) Encode() []byte {
	total := this.Size()
	buffer := make([]byte, total)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this Byteset) HeaderSize() uint32 {
	if len(this) == 0 {
		return 0
	}
	return uint32(len(this)+1) * UINT32_LEN
}

func (this Byteset) FillHeader(buffer []byte) {
	if len(this) == 0 {
		return
	}

	Uint32(len(this)).EncodeToBuffer(buffer)

	offset := uint32(0)
	for i := 0; i < len(this); i++ {
		Uint32(offset).EncodeToBuffer(buffer[(i+1)*UINT32_LEN:])
		offset += uint32(len(this[i]))
	}
}

func (this Byteset) EncodeToBuffer(buffer []byte) int {
	if len(buffer) == 0 {
		return 0
	}
	this.FillHeader(buffer)

	offset := this.HeaderSize()
	for i := 0; i < len(this); i++ {
		copy(buffer[offset:offset+uint32(len(this[i]))], this[i])
		offset += uint32(len(this[i]))
	}
	return int(offset)
}

func (this Byteset) Decode(buffer []byte) interface{} {
	if len(buffer) == 0 {
		return Byteset{}
	}

	count := uint32(Uint32(0).Decode(buffer[:UINT32_LEN]).(Uint32))
	this = make([][]byte, count)

	headerLen := (count + 1) * UINT32_LEN
	prev := uint32(Uint32(0).Decode(buffer[UINT32_LEN : UINT32_LEN+UINT32_LEN]).(Uint32))
	next := uint32(0)
	for i := 0; i < int(count); i++ {
		if i == int(count)-1 {
			next = uint32(len(buffer)) - headerLen
		} else {
			next = uint32(Uint32(0).Decode(buffer[UINT32_LEN+(i+1)*UINT32_LEN : UINT32_LEN+(i+2)*UINT32_LEN]).(Uint32))
		}

		this[i] = buffer[headerLen+prev : headerLen+next]
		prev = next
	}
	return Byteset(this)
}

type Bytegroup [][][]byte

func (this Bytegroup) Clone() Bytegroup {
	target := make([][][]byte, len(this))
	for i := range this {
		target[i] = make([][]byte, len(this[i]))
		for j := range this[i] {
			target[i][j] = make([]byte, len(this[i][j]))
			copy(target[i][j], this[i][j])
		}
	}
	return Bytegroup(target)
}

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
