package types

import (
	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
)

type EuResult struct {
	H           string
	ID          uint32
	Transitions [][]byte
	DC          *DeferCall
	Status      uint64
	GasUsed     uint64
}

func (this *EuResult) HeaderSize() uint32 {
	return 7 * codec.UINT32_LEN
}

func (this *EuResult) Size() uint32 {
	return this.HeaderSize() +
		uint32(len(this.H)) +
		codec.UINT32_LEN +
		codec.Byteset(this.Transitions).Size() +
		this.DC.Size() +
		codec.UINT64_LEN +
		codec.UINT64_LEN
}

// Fill in the header info
func (this *EuResult) FillHeader(buffer []byte) {
	offset := uint32(0)
	codec.Uint32(6).EncodeToBuffer(buffer[codec.UINT32_LEN*0:])

	codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*1:])
	offset += codec.String(this.H).Size()

	codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*2:])
	offset += codec.Uint32(this.ID).Size()

	codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*3:])
	offset += codec.Byteset(this.Transitions).Size()

	codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*4:])
	offset += this.DC.Size()

	codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*5:])
	offset += codec.Uint64(this.Status).Size()

	codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*6:])
}

func (this *EuResult) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this *EuResult) EncodeToBuffer(buffer []byte) {
	if this == nil {
		return
	}
	this.FillHeader(buffer)

	headerLen := this.HeaderSize()
	offset := uint32(0)

	codec.String(this.H).EncodeToBuffer(buffer[headerLen+offset:])
	offset += codec.String(this.H).Size()

	codec.Uint32(this.ID).EncodeToBuffer(buffer[headerLen+offset:])
	offset += codec.Uint32(this.ID).Size()

	codec.Byteset(this.Transitions).EncodeToBuffer(buffer[headerLen+offset:])
	offset += codec.Byteset(this.Transitions).Size()

	this.DC.EncodeToBuffer(buffer[headerLen+offset:])
	offset += this.DC.Size()

	codec.Uint64(this.Status).EncodeToBuffer(buffer[headerLen+offset:])
	offset += codec.Uint64(this.Status).Size()

	codec.Uint64(this.GasUsed).EncodeToBuffer(buffer[headerLen+offset:])
}

func (this *EuResult) Decode(buffer []byte) *EuResult {
	fields := [][]byte(codec.Byteset{}.Decode(buffer).(codec.Byteset))

	this.H = string(fields[0])
	this.ID = uint32(codec.Uint32(0).Decode(fields[1]).(codec.Uint32))

	this.Transitions = [][]byte(codec.Byteset{}.Decode(fields[2]).(codec.Byteset))
	if len(fields[3]) > 0 {
		this.DC = (&DeferCall{}).Decode(fields[3])
	}
	this.Status = uint64(codec.Uint64(0).Decode(fields[4]).(codec.Uint64))
	this.GasUsed = uint64(codec.Uint64(0).Decode(fields[5]).(codec.Uint64))
	return this
}

func (this *EuResult) GobEncode() ([]byte, error) {
	return this.Encode(), nil
}

func (this *EuResult) GobDecode(buffer []byte) error {
	this.Decode(buffer)
	return nil
}

func (tar *TxAccessRecords) GobEncode() ([]byte, error) {
	return tar.Encode(), nil
}

func (tar *TxAccessRecords) GobDecode(buffer []byte) error {
	tar.Decode(buffer)
	return nil
}

type Euresults []*EuResult

func (this *Euresults) HeaderSize() uint32 {
	return uint32((len(*this) + 1) * codec.UINT32_LEN) // Header length
}

func (this *Euresults) Size() uint32 {
	total := this.HeaderSize()
	for i := 0; i < len(*this); i++ {
		total += (*this)[i].Size()
	}
	return total
}

// Fill in the header info
func (this *Euresults) FillHeader(buffer []byte) {
	codec.Uint32(len(*this)).EncodeToBuffer(buffer)

	offset := uint32(0)
	for i := 0; i < len(*this); i++ {
		codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*(i+1):])
		offset += (*this)[i].Size()
	}
}

func (this Euresults) GobEncode() ([]byte, error) {
	buffer := make([]byte, this.Size())
	this.FillHeader(buffer)

	offsets := make([]uint32, len(this)+1)
	offsets[0] = 0
	for i := 0; i < len(this); i++ {
		offsets[i+1] = offsets[i] + this[i].Size()
	}

	headerLen := this.HeaderSize()
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			this[i].EncodeToBuffer(buffer[headerLen+offsets[i]:])
		}
	}
	common.ParallelWorker(len(this), 4, worker)
	return buffer, nil
}

func (this *Euresults) GobDecode(buffer []byte) error {
	bytesset := [][]byte(codec.Byteset{}.Decode(buffer).(codec.Byteset))
	euresults := make([]*EuResult, len(bytesset))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			euresults[i] = (&EuResult{}).Decode(bytesset[i])
		}
	}
	common.ParallelWorker(len(bytesset), 4, worker)
	*this = euresults
	return nil
}
