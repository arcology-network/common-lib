package types

import (
	"github.com/HPISTechnologies/common-lib/codec"
	"github.com/HPISTechnologies/common-lib/common"
)

type TxAccessRecords struct {
	Hash     string
	ID       uint32
	Accesses [][]byte
}

func (this *TxAccessRecords) HeaderSize() uint32 {
	return 4 * codec.UINT32_LEN
}

func (this *TxAccessRecords) Size() uint32 {
	return this.HeaderSize() +
		codec.String(this.Hash).Size() +
		codec.UINT32_LEN +
		codec.Byteset(this.Accesses).Size()
}

func (this *TxAccessRecords) FillHeader(buffer []byte) {
	codec.Uint32(3).EncodeToBuffer(buffer)
	codec.Uint32(0).EncodeToBuffer(buffer[codec.UINT32_LEN*1:])
	codec.Uint32(codec.String(this.Hash).Size()).EncodeToBuffer(buffer[codec.UINT32_LEN*2:])
	codec.Uint32(codec.String(this.Hash).Size() + codec.Uint32(this.ID).Size()).EncodeToBuffer(buffer[codec.UINT32_LEN*3:])
}

func (this *TxAccessRecords) Encode() []byte {
	buffer := make([]byte, this.Size())
	this.FillHeader(buffer)
	this.EncodeToBuffer(buffer)
	return buffer
}

func (this *TxAccessRecords) EncodeToBuffer(buffer []byte) {
	if this == nil {
		return
	}

	headerLen := this.HeaderSize()
	offset := uint32(0)

	codec.String(this.Hash).EncodeToBuffer(buffer[headerLen+offset:])
	offset += codec.String(this.Hash).Size()

	codec.Uint32(this.ID).EncodeToBuffer(buffer[headerLen+offset:])
	offset += codec.Uint32(this.ID).Size()

	codec.Byteset(this.Accesses).EncodeToBuffer(buffer[headerLen+offset:])
}

func (this *TxAccessRecords) Decode(buffer []byte) *TxAccessRecords {
	fields := codec.Byteset{}.Decode(buffer).(codec.Byteset)
	this.Hash = codec.Bytes(fields[0]).ToString()
	this.ID = uint32(codec.Uint32(0).Decode(fields[1]).(codec.Uint32))
	this.Accesses = codec.Byteset{}.Decode(fields[2]).(codec.Byteset)
	return this
}

type TxAccessRecordSet []*TxAccessRecords

func (this *TxAccessRecordSet) HeaderSize() uint32 {
	return uint32((len(*this) + 1) * codec.UINT32_LEN)
}

func (this *TxAccessRecordSet) Size() uint32 {
	total := this.HeaderSize()        // Header length
	for i := 0; i < len(*this); i++ { // Body  length
		total += (*this)[i].Size()
	}
	return total
}

// Fill in the header info
func (this *TxAccessRecordSet) FillHeader(buffer []byte) {
	offset := uint32(0)
	codec.Uint32(len(*this)).EncodeToBuffer(buffer)
	for i := 0; i < len(*this); i++ {
		codec.Uint32(offset).EncodeToBuffer(buffer[(i+1)*codec.UINT32_LEN:])
		offset += (*this)[i].Size()
	}
}

func (this TxAccessRecordSet) GobEncode() ([]byte, error) {
	buffer := make([]byte, this.Size())
	this.FillHeader(buffer)

	headerLen := this.HeaderSize()
	offsets := make([]uint32, len(this)+1)
	offsets[0] = 0
	for i := 0; i < len(this); i++ {
		offsets[i+1] = offsets[i] + this[i].Size()
	}

	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			this[i].FillHeader(buffer[headerLen+offsets[i]:])
			this[i].EncodeToBuffer(buffer[headerLen+offsets[i]:])
		}
	}
	common.ParallelWorker(len(this), 4, worker)
	return buffer, nil
}

func (this *TxAccessRecordSet) GobDecode(data []byte) error {
	bytesset := codec.Byteset{}.Decode(data).(codec.Byteset)
	records := make([]*TxAccessRecords, len(bytesset))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			this := &TxAccessRecords{}
			this.Decode(bytesset[i])
			records[i] = this
		}
	}
	common.ParallelWorker(len(bytesset), 6, worker)
	*this = records
	return nil
}
