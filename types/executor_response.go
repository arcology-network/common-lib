package types

import (
	codec "github.com/arcology-network/common-lib/codec"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

type ExecuteResponse struct {
	// DfCall  *DeferredCall
	Hash    ethCommon.Hash
	Status  uint64
	GasUsed uint64
}

type ExecutorResponses struct {
	// DfCalls     []*DeferredCall
	HashList    []ethCommon.Hash
	StatusList  []uint64
	GasUsedList []uint64

	// SpawnedKeys       []ethCommon.Hash
	// SpawnedTxs        []ethCommon.Hash
	// RelationKeys      []ethCommon.Hash
	// RelationSizes     []uint64
	// RelationValues    []ethCommon.Hash
	ContractAddresses []ethCommon.Address
	// TxidsHash         []ethCommon.Hash
	// TxidsId           []uint32
	// TxidsAddress      []ethCommon.Address
	CallResults [][]byte
}

func (er *ExecutorResponses) GobEncode() ([]byte, error) {
	data := [][]byte{
		Hashes(er.HashList).Encode(),
		codec.Uint64s(er.StatusList).Encode(),
		codec.Uint64s(er.GasUsedList).Encode(),
		// DeferCalls(er.DfCalls).Encode(),
		// Hashes(er.SpawnedKeys).Encode(),
		// Hashes(er.SpawnedTxs).Encode(),
		// Hashes(er.RelationKeys).Encode(),
		// codec.Uint64s(er.RelationSizes).Encode(),
		// Hashes(er.RelationValues).Encode(),
		Addresses(er.ContractAddresses).Encode(),
		// Hashes(er.TxidsHash).Encode(),
		// codec.Uint32s(er.TxidsId).Encode(),
		// Addresses(er.TxidsAddress).Encode(),
		codec.Byteset(er.CallResults).Encode(),
	}
	return codec.Byteset(data).Encode(), nil
}
func (er *ExecutorResponses) GobDecode(data []byte) error {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	er.HashList = Hashes(er.HashList).Decode(fields[0])
	er.StatusList = []uint64(codec.Uint64s(er.StatusList).Decode(fields[1]).(codec.Uint64s))
	er.GasUsedList = []uint64(codec.Uint64s(er.GasUsedList).Decode(fields[2]).(codec.Uint64s))
	// er.DfCalls = new(DeferCalls).Decode(fields[3])
	// er.SpawnedKeys = Hashes(er.SpawnedKeys).Decode(fields[4])
	// er.SpawnedTxs = Hashes(er.SpawnedTxs).Decode(fields[5])
	// er.RelationKeys = Hashes(er.RelationKeys).Decode(fields[6])
	// er.RelationSizes = codec.Uint64s(er.RelationSizes).Decode(fields[7])
	// er.RelationValues = Hashes(er.RelationValues).Decode(fields[8])
	er.ContractAddresses = Addresses(er.ContractAddresses).Decode(fields[3])
	// er.TxidsHash = Hashes(er.TxidsHash).Decode(fields[10])
	// er.TxidsId = codec.Uint32s(er.TxidsId).Decode(fields[11])
	// er.TxidsAddress = Addresses(er.TxidsAddress).Decode(fields[12])
	er.CallResults = [][]byte(codec.Byteset{}.Decode(fields[4]).(codec.Byteset))
	return nil
}
