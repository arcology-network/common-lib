package types

import (
	"encoding/gob"
	"math/big"

	evmCommon "github.com/arcology-network/evm/common"
)

var arEncoder *arbReqEncoder
var arDecoder *arbReqDecoder
var bytesPool chan []byte

func init() {
	gob.Register(&NodeRole{})
	gob.Register(&InclusiveList{})
	gob.Register(&ApplyListItem{})
	gob.Register(&ParentInfo{})
	gob.Register(&PartialHeader{})
	gob.Register(&ReapingList{})
	gob.Register(&ReceiptHashList{})
	gob.Register(&[]*EuResult{})
	//gob.Register(&Reads{})
	gob.Register(&ArbitratorRequest{})
	gob.Register(&ArbitratorResponse{})
	gob.Register(&ExecutorRequest{})
	gob.Register(&ExecutorResponses{})
	gob.Register(&StandardMessage{})
	gob.Register(StandardMessages{})
	gob.Register(Txs{})
	gob.Register(&StandardTransaction{})
	gob.Register([]*StandardMessage{})
	gob.Register([]*StandardTransaction{})
	gob.Register(&StatisticalInformation{})

	gob.Register(RequestBalance{})
	gob.Register(RequestContainer{})
	gob.Register(&RequestBlock{})
	gob.Register(&RequestReceipt{})
	gob.Register(Block{})
	gob.Register(Log{})
	gob.Register([]*QueryReceipt{})
	gob.Register([][]byte{})
	gob.Register([]byte{})
	gob.Register(&MetaBlock{})
	gob.Register(&MonacoBlock{})
	gob.Register(&big.Int{})
	gob.Register(SendingStandardMessages{})
	gob.Register(ExecutingLogs{})
	// gob.Register([]*SpawnedRelation{})
	gob.Register(map[evmCommon.Hash]evmCommon.Hash{})

	gob.Register(&[]*TxAccessRecords{})
	gob.Register(&TxAccessRecordSet{})
	gob.Register(&Euresults{})
	// gob.Register(&DeferredCall{})
	gob.Register(&RequestParameters{})
	gob.Register(&RequestBlockEth{})
	gob.Register(&RequestStorage{})

	gob.Register(&SyncStatus{})
	gob.Register(&SyncPoint{})
	gob.Register(&SyncDataRequest{})
	gob.Register(&SyncDataResponse{})

	gob.Register(&IncomingTxs{})
	gob.Register(&IncomingMsgs{})

	arEncoder = newArbReqEncoder()
	arDecoder = newArbReqDecoder()

	bytesPool = make(chan []byte, 100)
	for i := 0; i < 100; i++ {
		bytesPool <- make([]byte, 0, 2*1024*1024)
	}
}
