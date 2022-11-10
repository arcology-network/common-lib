package types

import (
	"crypto/sha256"
	"math/big"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/encoding"
)

type ExecutingSequence struct {
	Msgs       []*StandardMessage
	Parallel   bool
	SequenceId ethCommon.Hash
	Txids      []uint32
}

func NewExecutingSequence(msgs []*StandardMessage, parallel bool) *ExecutingSequence {
	buffers := make([][]byte, len(msgs))
	for i, msg := range msgs {
		buffers[i] = msg.TxHash.Bytes()
	}

	hash := sha256.Sum256(encoding.Byteset(buffers).Encode())
	return &ExecutingSequence{
		Msgs:       msgs,
		Parallel:   parallel,
		SequenceId: ethCommon.BytesToHash(hash[:]),
		Txids:      make([]uint32, len(msgs)),
	}
}

type ExecutingSequences []*ExecutingSequence

func (this ExecutingSequences) Encode() ([]byte, error) {
	if this == nil {
		return []byte{}, nil
	}

	data := make([][]byte, len(this))
	worker := func(start, end, idx int, args ...interface{}) {
		executingSequences := args[0].([]interface{})[0].(ExecutingSequences)
		data := args[0].([]interface{})[1].([][]byte)
		for i := start; i < end; i++ {
			standardMessages := StandardMessages(executingSequences[i].Msgs)
			standardMessagesData, err := standardMessages.Encode()
			if err != nil {
				standardMessagesData = []byte{}
			}

			tmpData := [][]byte{
				standardMessagesData,
				codec.Bools([]bool{executingSequences[i].Parallel}).Encode(),
				executingSequences[i].SequenceId[:],
				codec.Uint32s(executingSequences[i].Txids).Encode(),
			}
			data[i] = codec.Byteset(tmpData).Encode()
		}
	}
	common.ParallelWorker(len(this), concurrency, worker, this, data)
	return codec.Byteset(data).Encode(), nil
}

func (this *ExecutingSequences) Decode(data []byte) ([]*ExecutingSequence, error) {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	v := ExecutingSequences(make([]*ExecutingSequence, len(fields)))
	this = &v

	worker := func(start, end, idx int, args ...interface{}) {
		datas := args[0].([]interface{})[0].(codec.Byteset)
		executingSequences := args[0].([]interface{})[1].(ExecutingSequences)

		for i := start; i < end; i++ {
			executingSequence := new(ExecutingSequence)

			datafields := codec.Byteset{}.Decode(datas[i]).(codec.Byteset)
			msgResults, err := new(StandardMessages).Decode(datafields[0])
			if err != nil {
				msgResults = StandardMessages{}
			}
			executingSequence.Msgs = msgResults
			parallels := new(encoding.Bools).Decode(datafields[1])
			if len(parallels) > 0 {
				executingSequence.Parallel = parallels[0]
			}
			executingSequence.SequenceId = ethCommon.BytesToHash(datafields[2])
			executingSequence.Txids = new(encoding.Uint32s).Decode(datafields[3])
			executingSequences[i] = executingSequence

		}
	}
	common.ParallelWorker(len(fields), concurrency, worker, fields, *this)
	return ([]*ExecutingSequence)(*this), nil
}

type ExecutorRequest struct {
	Sequences     []*ExecutingSequence
	Precedings    [][]*ethCommon.Hash
	PrecedingHash []ethCommon.Hash
	Timestamp     *big.Int
	Parallelism   uint64
	Debug         bool
}

func (this *ExecutorRequest) GobEncode() ([]byte, error) {
	executingSequences := ExecutingSequences(this.Sequences)
	executingSequencesData, err := executingSequences.Encode()
	if err != nil {
		return []byte{}, err
	}

	precedingsBytes := make([][]byte, len(this.Precedings))
	for i := range this.Precedings {
		precedings := Ptr2Arr(this.Precedings[i])
		precedingsBytes[i] = ethCommon.Hashes(precedings).Encode()
	}

	timeStampData := []byte{}
	if this.Timestamp != nil {
		timeStampData = this.Timestamp.Bytes()
	}

	data := [][]byte{
		executingSequencesData,
		encoding.Byteset(precedingsBytes).Encode(),
		ethCommon.Hashes(this.PrecedingHash).Encode(),
		timeStampData,
		common.Uint64ToBytes(this.Parallelism),
		codec.Bool(this.Debug).Encode(),
	}
	return encoding.Byteset(data).Encode(), nil
}

func (this *ExecutorRequest) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	msgResults, err := new(ExecutingSequences).Decode(fields[0])
	if err != nil {
		return err
	}
	this.Sequences = msgResults

	precedingsBytes := encoding.Byteset{}.Decode(fields[1])
	this.Precedings = make([][]*ethCommon.Hash, len(precedingsBytes))
	for i := range precedingsBytes {
		this.Precedings[i] = Arr2Ptr(ethCommon.Hashes([]ethCommon.Hash{}).Decode(precedingsBytes[i]))
	}

	this.PrecedingHash = ethCommon.Hashes([]ethCommon.Hash{}).Decode(fields[2])
	//if len(fields[3]) > 0 {
	this.Timestamp = new(big.Int).SetBytes(fields[3])
	//}
	//if len(fields[4]) > 0 {
	this.Parallelism = common.BytesToUint64(fields[4])
	//}
	this.Debug = bool(codec.Bool(this.Debug).Decode(fields[5]).(codec.Bool))
	return nil
}
