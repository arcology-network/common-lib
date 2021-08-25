package types

import (
	"crypto/sha256"
	"math/big"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
	"github.com/HPISTechnologies/common-lib/common"
	"github.com/HPISTechnologies/common-lib/encoding"
)

type ExecutingSequence struct {
	Msgs       []*StandardMessage
	Parallel   bool
	SequenceId ethCommon.Hash
	Txids      []uint32
}

func NewExecutingSequence(msgs []*StandardMessage, parallel bool) *ExecutingSequence {
	datas := make([][]byte, len(msgs))
	for i, msg := range msgs {
		datas[i] = msg.TxHash.Bytes()
	}

	hash := sha256.Sum256(encoding.Byteset(datas).Encode())
	return &ExecutingSequence{
		Msgs:       msgs,
		Parallel:   parallel,
		SequenceId: ethCommon.BytesToHash(hash[:]),
		Txids:      make([]uint32, len(msgs)),
	}
}

type ExecutingSequences []*ExecutingSequence

func (ess ExecutingSequences) Encode() ([]byte, error) {
	if ess == nil {
		return []byte{}, nil
	}
	data := make([][]byte, len(ess))
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
				encoding.Bools([]bool{executingSequences[i].Parallel}).Encode(),
				executingSequences[i].SequenceId[:],
				encoding.Uint32s(executingSequences[i].Txids).Encode(),
			}
			data[i] = encoding.Byteset(tmpData).Encode()
		}
	}

	common.ParallelWorker(len(ess), concurrency, worker, ess, data)
	return encoding.Byteset(data).Encode(), nil
}
func (ess *ExecutingSequences) Decode(data []byte) ([]*ExecutingSequence, error) {
	fields := encoding.Byteset{}.Decode(data)
	esss := make([]*ExecutingSequence, len(fields))

	worker := func(start, end, idx int, args ...interface{}) {
		datas := args[0].([]interface{})[0].([][]byte)
		executingSequences := args[0].([]interface{})[1].([]*ExecutingSequence)

		for i := start; i < end; i++ {
			executingSequence := new(ExecutingSequence)

			fields := encoding.Byteset{}.Decode(datas[i])

			msgResults, err := new(StandardMessages).Decode(fields[0])
			if err != nil {
				msgResults = StandardMessages{}
			}
			executingSequence.Msgs = msgResults
			parallels := new(encoding.Bools).Decode(fields[1])
			if len(parallels) > 0 {
				executingSequence.Parallel = parallels[0]
			}
			executingSequence.SequenceId = ethCommon.BytesToHash(fields[2])
			executingSequence.Txids = new(encoding.Uint32s).Decode(fields[3])
			executingSequences[i] = executingSequence

		}
	}
	common.ParallelWorker(len(fields), concurrency, worker, fields, esss)

	return esss, nil
}

type ExecutorRequest struct {
	Sequences     []*ExecutingSequence
	Precedings    []*ethCommon.Hash
	PrecedingHash ethCommon.Hash
	Timestamp     *big.Int
	Parallelism   uint64
}

func (er *ExecutorRequest) GobEncode() ([]byte, error) {
	precedings := Ptr2Arr(er.Precedings)
	executingSequences := ExecutingSequences(er.Sequences)
	executingSequencesData, err := executingSequences.Encode()
	if err != nil {
		return []byte{}, err
	}
	timeStampData := []byte{}
	if er.Timestamp != nil {
		timeStampData = er.Timestamp.Bytes()
	}
	data := [][]byte{
		executingSequencesData,
		ethCommon.Hashes(precedings).Encode(),
		er.PrecedingHash.Bytes(),
		timeStampData,
		common.Uint64ToBytes(er.Parallelism),
	}
	return encoding.Byteset(data).Encode(), nil
}
func (er *ExecutorRequest) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	msgResults, err := new(ExecutingSequences).Decode(fields[0])
	if err != nil {
		return err
	}
	er.Sequences = msgResults
	arrs := []ethCommon.Hash{}
	arrs = ethCommon.Hashes(arrs).Decode(fields[1])
	er.Precedings = Arr2Ptr(arrs)
	er.PrecedingHash = ethCommon.BytesToHash(fields[2])
	if len(fields[3]) > 0 {
		er.Timestamp = new(big.Int).SetBytes(fields[3])
	}
	if len(fields[4]) > 0 {
		er.Parallelism = common.BytesToUint64(fields[4])
	}

	return nil
}
