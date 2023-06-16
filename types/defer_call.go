package types

import (
	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	encoding "github.com/arcology-network/common-lib/encoding"
)

type DeferredCall struct {
	DeferID         string
	GroupBy         [32]byte
	ContractAddress Address
	Signature       string
	Data            []byte
}

func (this *DeferredCall) Encode() []byte {
	buffers := [][]byte{
		codec.String(this.DeferID).ToBytes(),
		codec.String(this.ContractAddress).ToBytes(),
		codec.String(this.Signature).ToBytes(),
	}

	return codec.Byteset(buffers).Encode()
}

func (this *DeferredCall) Decode(data []byte) *DeferredCall {
	buffers := [][]byte(codec.Byteset{}.Decode(data).(codec.Byteset))
	this.DeferID = string(buffers[0])
	this.ContractAddress = Address(buffers[1])
	this.Signature = string(buffers[2])
	return this
}

func (this *DeferredCall) HeaderSize() uint32 {
	return 4 * codec.UINT32_LEN
}

func (this *DeferredCall) Size() uint32 {
	if this == nil {
		return 0
	}

	return 4*codec.UINT32_LEN +
		uint32(len(this.DeferID)+len(this.ContractAddress)+len(this.Signature))
}

func (this *DeferredCall) EncodeToBuffer(buffer []byte) int {
	if this == nil {
		return 0
	}

	offset := codec.Encoder{}.FillHeader(
		buffer,
		[]uint32{
			codec.String(this.DeferID).Size(),
			codec.String(this.ContractAddress).Size(),
			codec.String(this.Signature).Size(),
		},
	)

	offset += codec.String(this.DeferID).EncodeToBuffer(buffer[offset:])
	offset += codec.String(this.ContractAddress).EncodeToBuffer(buffer[offset:])
	offset += codec.String(this.Signature).EncodeToBuffer(buffer[offset:])
	return offset
}

type DeferCalls []*DeferredCall

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

func (dcs *DeferCalls) Decode(data []byte) []*DeferredCall {
	buffers := encoding.Byteset{}.Decode(data)
	defs := make([]*DeferredCall, len(buffers))

	worker := func(start, end, idx int, args ...interface{}) {
		dataSet := args[0].([]interface{})[0].([][]byte)
		defcalls := args[0].([]interface{})[1].([]*DeferredCall)

		for i := start; i < end; i++ {
			if len(dataSet[i]) == 0 {
				continue
			}
			deferCall := new(DeferredCall)
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
