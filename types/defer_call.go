package types

import (
	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	encoding "github.com/arcology-network/common-lib/encoding"
)

type DeferCall struct {
	DeferID         string
	ContractAddress Address
	Signature       string
}

func (this *DeferCall) Encode() []byte {
	buffers := [][]byte{
		codec.String(this.DeferID).ToBytes(),
		codec.String(this.ContractAddress).ToBytes(),
		codec.String(this.Signature).ToBytes(),
	}

	return codec.Byteset(buffers).Encode()
}

func (this *DeferCall) Decode(data []byte) *DeferCall {
	buffers := [][]byte(codec.Byteset{}.Decode(data).(codec.Byteset))
	this.DeferID = string(buffers[0])
	this.ContractAddress = Address(buffers[1])
	this.Signature = string(buffers[2])
	return this
}

func (this *DeferCall) HeaderSize() uint32 {
	return 4 * codec.UINT32_LEN
}

func (this *DeferCall) FillHeader(buffer []byte) {
	offset := uint32(0)
	codec.Uint32(3).EncodeToBuffer(buffer[codec.UINT32_LEN*0:])

	codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*1:])
	offset += codec.String(this.DeferID).Size()

	codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*2:])
	offset += codec.String(this.ContractAddress).Size()

	codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*3:])
}

func (this *DeferCall) Size() uint32 {
	if this == nil {
		return 0
	}

	return 4*codec.UINT32_LEN +
		uint32(len(this.DeferID)+len(this.ContractAddress)+len(this.Signature))
}

func (this *DeferCall) EncodeToBuffer(buffer []byte) {
	if this == nil {
		return
	}

	this.FillHeader(buffer)
	headerLen := this.HeaderSize()
	offset := uint32(0)

	codec.String(this.DeferID).EncodeToBuffer(buffer[headerLen+offset:])
	offset += codec.String(this.DeferID).Size()

	codec.String(this.ContractAddress).EncodeToBuffer(buffer[headerLen+offset:])
	offset += codec.String(this.ContractAddress).Size()

	codec.String(this.Signature).EncodeToBuffer(buffer[headerLen+offset:])
}

type DeferCalls []*DeferCall

func (dcs DeferCalls) Encode() []byte {
	if dcs == nil {
		return []byte{}
	}

	worker := func(start, end int, idx int, args ...interface{}) {
		defcalls := args[0].([]interface{})[0].(DeferCalls)
		dataSet := args[0].([]interface{})[1].([][]byte)

		for i := start; i < end; i++ {
			if defcall := defcalls[i]; defcall != nil {
				dataSet[i] = encoding.Byteset([][]byte{
					encoding.Uint32(len(defcall.DeferID[:])).Encode()[:],
					encoding.Uint32(len(defcall.Signature[:])).Encode()[:],
					[]byte(defcall.DeferID),
					[]byte(defcall.ContractAddress),
					[]byte(defcall.Signature),
				}).Flatten()
			}
		}
	}

	dataSet := make([][]byte, len(dcs))
	common.ParallelWorker(len(dcs), concurrency, worker, dcs, dataSet)
	return encoding.Byteset(dataSet).Encode()
}

func (dcs *DeferCalls) Decode(data []byte) []*DeferCall {
	buffers := encoding.Byteset{}.Decode(data)
	defs := make([]*DeferCall, len(buffers))

	worker := func(start, end, idx int, args ...interface{}) {
		dataSet := args[0].([]interface{})[0].([][]byte)
		defcalls := args[0].([]interface{})[1].([]*DeferCall)

		for i := start; i < end; i++ {
			if len(dataSet[i]) == 0 {
				continue
			}
			deferCall := new(DeferCall)
			DeferIDLength := 0
			DeferIDLength = int(encoding.Uint32(DeferIDLength).Decode(dataSet[i][0:encoding.UINT32_LEN]))
			SignatureLength := 0
			SignatureLength = int(encoding.Uint32(SignatureLength).Decode(dataSet[i][encoding.UINT32_LEN : encoding.UINT32_LEN*2]))

			deferCall.DeferID = string(dataSet[i][encoding.UINT32_LEN*2 : encoding.UINT32_LEN*2+DeferIDLength])
			deferCall.ContractAddress = Address(string(dataSet[i][encoding.UINT32_LEN*2+DeferIDLength : encoding.UINT32_LEN*2+DeferIDLength+AddressLength]))
			deferCall.Signature = string(dataSet[i][encoding.UINT32_LEN*2+DeferIDLength+AddressLength:])
			defcalls[i] = deferCall
		}
	}
	common.ParallelWorker(len(buffers), concurrency, worker, buffers, defs)

	return defs
}
