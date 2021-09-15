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

func (dc *DeferCall) Encode() []byte {
	tmpData := [][]byte{
		[]byte(dc.DeferID),
		[]byte(dc.ContractAddress),
		[]byte(dc.Signature),
	}
	return encoding.Byteset(tmpData).Encode()
}

func (dc *DeferCall) Decode(data []byte) {
	fields := encoding.Byteset{}.Decode(data)
	dc.DeferID = string(fields[0])
	dc.ContractAddress = Address(fields[1])
	dc.Signature = string(fields[2])
}

type EuResult struct {
	H           string
	ID          uint32
	Transitions [][]byte
	DC          *DeferCall
	Status      uint64
	GasUsed     uint64
}

func (er *EuResult) Encode() []byte {
	dcData := []byte{}
	if er.DC != nil {
		dcData = er.DC.Encode()
	}
	tmpData := [][]byte{
		[]byte(er.H),
		encoding.Uint32(er.ID).Encode(),
		encoding.Byteset(er.Transitions).Encode(),
		dcData,
		encoding.Uint64(er.Status).Encode(),
		encoding.Uint64(er.GasUsed).Encode(),
	}
	return encoding.Byteset(tmpData).Encode()
}
func (er *EuResult) Decode(data []byte) {
	fields := encoding.Byteset{}.Decode(data)

	er.H = string(fields[0])
	er.ID = encoding.Uint32(0).Decode(fields[1])
	er.Transitions = encoding.Byteset{}.Decode(fields[2])
	if len(fields[3]) > 0 {
		dc := &DeferCall{}
		dc.Decode(fields[3])
		er.DC = dc
	}
	er.Status = encoding.Uint64(0).Decode(fields[4])
	er.GasUsed = encoding.Uint64(0).Decode(fields[5])
}

func (er *EuResult) GobEncode() ([]byte, error) {
	return er.Encode(), nil
}
func (er *EuResult) GobDecode(data []byte) error {
	er.Decode(data)
	return nil
}

type TxAccessRecords struct {
	Hash     string
	ID       uint32
	Accesses [][]byte
}

func (tar *TxAccessRecords) Encode() []byte {
	tmpData := [][]byte{
		[]byte(tar.Hash),
		encoding.Uint32(tar.ID).Encode(),
		encoding.Byteset(tar.Accesses).Encode(),
	}
	return encoding.Byteset(tmpData).Encode()
}

func (tar *TxAccessRecords) Decode(data []byte) {
	fields := encoding.Byteset{}.Decode(data)

	tar.Hash = string(fields[0])
	tar.ID = encoding.Uint32(0).Decode(fields[1])
	tar.Accesses = encoding.Byteset{}.Decode(fields[2])
}

func (tar *TxAccessRecords) GobEncode() ([]byte, error) {
	return tar.Encode(), nil
}
func (tar *TxAccessRecords) GobDecode(data []byte) error {
	tar.Decode(data)
	return nil
}

type Euresults []*EuResult

func (ers Euresults) GobEncode() ([]byte, error) {
	byteset := make([][]byte, len(ers))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			byteset[i] = ers[i].Encode()
		}
	}
	common.ParallelWorker(len(ers), 4, worker)
	return codec.Byteset(byteset).Encode(), nil
}
func (ers *Euresults) GobDecode(data []byte) error {
	bytesset := codec.Byteset{}.Decode(data)
	euresults := make([]*EuResult, len(bytesset))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			euresult := &EuResult{}
			euresult.Decode(bytesset[i])
			euresults[i] = euresult
		}
	}
	common.ParallelWorker(len(bytesset), 4, worker)
	*ers = euresults
	return nil
}

type TxAccessRecordses []*TxAccessRecords

func (tars TxAccessRecordses) GobEncode() ([]byte, error) {
	byteset := make([][]byte, len(tars))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			byteset[i] = tars[i].Encode()
		}
	}
	common.ParallelWorker(len(tars), 4, worker)
	return codec.Byteset(byteset).Encode(), nil
}
func (tars *TxAccessRecordses) GobDecode(data []byte) error {
	bytesset := codec.Byteset{}.Decode(data)
	tarses := make([]*TxAccessRecords, len(bytesset))
	worker := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			tar := &TxAccessRecords{}
			tar.Decode(bytesset[i])
			tarses[i] = tar
		}
	}
	common.ParallelWorker(len(bytesset), 4, worker)
	*tars = tarses
	return nil
}
