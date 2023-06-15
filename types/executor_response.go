package types

import (
	ethCommon "github.com/arcology-network/3rd-party/eth/common"
	"github.com/arcology-network/common-lib/codec"
	encoding "github.com/arcology-network/common-lib/encoding"
)

type ExecuteResponse struct {
	DfCall  *DeferredCall
	Hash    ethCommon.Hash
	Status  uint64
	GasUsed uint64
}

type ExecutorResponses struct {
	DfCalls     []*DeferredCall
	HashList    []ethCommon.Hash
	StatusList  []uint64
	GasUsedList []uint64

	SpawnedKeys       []ethCommon.Hash
	SpawnedTxs        []ethCommon.Hash
	RelationKeys      []ethCommon.Hash
	RelationSizes     []uint64
	RelationValues    []ethCommon.Hash
	ContractAddresses []ethCommon.Address
	TxidsHash         []ethCommon.Hash
	TxidsId           []uint32
	TxidsAddress      []ethCommon.Address
	CallResults       [][]byte
}

func (er *ExecutorResponses) GobEncode() ([]byte, error) {
	data := [][]byte{
		ethCommon.Hashes(er.HashList).Encode(),
		encoding.Uint64s(er.StatusList).Encode(),
		encoding.Uint64s(er.GasUsedList).Encode(),
		DeferCalls(er.DfCalls).Encode(),
		ethCommon.Hashes(er.SpawnedKeys).Encode(),
		ethCommon.Hashes(er.SpawnedTxs).Encode(),
		ethCommon.Hashes(er.RelationKeys).Encode(),
		encoding.Uint64s(er.RelationSizes).Encode(),
		ethCommon.Hashes(er.RelationValues).Encode(),
		ethCommon.Addresses(er.ContractAddresses).Encode(),
		ethCommon.Hashes(er.TxidsHash).Encode(),
		encoding.Uint32s(er.TxidsId).Encode(),
		ethCommon.Addresses(er.TxidsAddress).Encode(),
		codec.Byteset(er.CallResults).Encode(),
	}
	return encoding.Byteset(data).Encode(), nil
}
func (er *ExecutorResponses) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	er.HashList = ethCommon.Hashes(er.HashList).Decode(fields[0])
	er.StatusList = encoding.Uint64s(er.StatusList).Decode(fields[1])
	er.GasUsedList = encoding.Uint64s(er.GasUsedList).Decode(fields[2])
	er.DfCalls = new(DeferCalls).Decode(fields[3])
	er.SpawnedKeys = ethCommon.Hashes(er.SpawnedKeys).Decode(fields[4])
	er.SpawnedTxs = ethCommon.Hashes(er.SpawnedTxs).Decode(fields[5])
	er.RelationKeys = ethCommon.Hashes(er.RelationKeys).Decode(fields[6])
	er.RelationSizes = encoding.Uint64s(er.RelationSizes).Decode(fields[7])
	er.RelationValues = ethCommon.Hashes(er.RelationValues).Decode(fields[8])
	er.ContractAddresses = ethCommon.Addresses(er.ContractAddresses).Decode(fields[9])
	er.TxidsHash = ethCommon.Hashes(er.TxidsHash).Decode(fields[10])
	er.TxidsId = encoding.Uint32s(er.TxidsId).Decode(fields[11])
	er.TxidsAddress = ethCommon.Addresses(er.TxidsAddress).Decode(fields[12])
	er.CallResults = [][]byte(codec.Byteset{}.Decode(fields[13]).(codec.Byteset))
	return nil
}
