package types

import (
	"github.com/arcology-network/common-lib/codec"
	encoding "github.com/arcology-network/common-lib/encoding"
	evmCommon "github.com/arcology-network/evm/common"
)

type ExecuteResponse struct {
	DfCall  *DeferCall
	Hash    evmCommon.Hash
	Status  uint64
	GasUsed uint64
}

type ExecutorResponses struct {
	DfCalls     []*DeferCall
	HashList    []evmCommon.Hash
	StatusList  []uint64
	GasUsedList []uint64

	SpawnedKeys       []evmCommon.Hash
	SpawnedTxs        []evmCommon.Hash
	RelationKeys      []evmCommon.Hash
	RelationSizes     []uint64
	RelationValues    []evmCommon.Hash
	ContractAddresses []evmCommon.Address
	TxidsHash         []evmCommon.Hash
	TxidsId           []uint32
	TxidsAddress      []evmCommon.Address
	CallResults       [][]byte
}

func (er *ExecutorResponses) GobEncode() ([]byte, error) {
	data := [][]byte{
		Hashes(er.HashList).Encode(),
		encoding.Uint64s(er.StatusList).Encode(),
		encoding.Uint64s(er.GasUsedList).Encode(),
		DeferCalls(er.DfCalls).Encode(),
		Hashes(er.SpawnedKeys).Encode(),
		Hashes(er.SpawnedTxs).Encode(),
		Hashes(er.RelationKeys).Encode(),
		encoding.Uint64s(er.RelationSizes).Encode(),
		Hashes(er.RelationValues).Encode(),
		Addresses(er.ContractAddresses).Encode(),
		Hashes(er.TxidsHash).Encode(),
		encoding.Uint32s(er.TxidsId).Encode(),
		Addresses(er.TxidsAddress).Encode(),
		codec.Byteset(er.CallResults).Encode(),
	}
	return encoding.Byteset(data).Encode(), nil
}
func (er *ExecutorResponses) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	er.HashList = Hashes(er.HashList).Decode(fields[0])
	er.StatusList = encoding.Uint64s(er.StatusList).Decode(fields[1])
	er.GasUsedList = encoding.Uint64s(er.GasUsedList).Decode(fields[2])
	er.DfCalls = new(DeferCalls).Decode(fields[3])
	er.SpawnedKeys = Hashes(er.SpawnedKeys).Decode(fields[4])
	er.SpawnedTxs = Hashes(er.SpawnedTxs).Decode(fields[5])
	er.RelationKeys = Hashes(er.RelationKeys).Decode(fields[6])
	er.RelationSizes = encoding.Uint64s(er.RelationSizes).Decode(fields[7])
	er.RelationValues = Hashes(er.RelationValues).Decode(fields[8])
	er.ContractAddresses = Addresses(er.ContractAddresses).Decode(fields[9])
	er.TxidsHash = Hashes(er.TxidsHash).Decode(fields[10])
	er.TxidsId = encoding.Uint32s(er.TxidsId).Decode(fields[11])
	er.TxidsAddress = Addresses(er.TxidsAddress).Decode(fields[12])
	er.CallResults = [][]byte(codec.Byteset{}.Decode(fields[13]).(codec.Byteset))
	return nil
}
