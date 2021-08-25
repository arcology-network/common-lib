package types

import (
	"bytes"
	"math/big"
	"math/rand"
	"sort"
	"time"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
	ethTypes "github.com/HPISTechnologies/3rd-party/eth/types"
	"github.com/HPISTechnologies/common-lib/common"
	"github.com/HPISTechnologies/common-lib/encoding"
	"github.com/HPISTechnologies/evm/crypto"
)

const (
	concurrency = 4
)

type StandardMessage struct {
	TxHash    ethCommon.Hash
	Native    *ethTypes.Message
	TxRawData []byte
	Source    uint8
}

func MakeMessageWithDefCall(def *DeferCall, hash ethCommon.Hash) *StandardMessage {
	signature := def.Signature
	contractAddress := def.ContractAddress
	data := crypto.Keccak256([]byte(signature))[:4]
	data = append(data, ethCommon.AlignToEvmForInt(ethCommon.EvmWordSize)...)
	idLen := ethCommon.AlignToEvmForInt(len(def.DeferID))
	id := ethCommon.AlignToEvmForString(def.DeferID)
	data = append(data, idLen...)
	data = append(data, id...)
	contractAddr := ethCommon.BytesToAddress([]byte(contractAddress))
	nonce := uint64(time.Now().UnixNano())
	message := ethTypes.NewMessage(contractAddr, &contractAddr, nonce, new(big.Int).SetInt64(0), 1e9, new(big.Int).SetInt64(0), data, false)
	standardMessager := StandardMessage{
		Native: &message,
		TxHash: hash,
	}
	return &standardMessager
}

func (stdMsg *StandardMessage) Hash() ethCommon.Hash {
	return stdMsg.TxHash
}

func (stdMsg *StandardMessage) Key() string {
	return stdMsg.TxHash.String()
}

func (stdMsg *StandardMessage) Equal(other *StandardMessage) bool {
	return stdMsg.TxHash.String() == other.TxHash.String()
}

func (lft *StandardMessage) CompareHash(rgt *StandardMessage) bool {
	return bytes.Compare(lft.TxHash[:], rgt.TxHash[:]) < 0
}

func (lft *StandardMessage) CompareGas(rgt *StandardMessage) bool {
	lftFrom, rgtFrom := lft.Native.From(), rgt.Native.From()
	if bytes.Compare(lftFrom[:], rgtFrom[:]) == 0 { // by nonce if from the same address
		return lft.Native.Nonce() < rgt.Native.Nonce()
	}

	if v := lft.Native.GasPrice().Cmp(rgt.Native.GasPrice()); v == 0 { // by address if fees are the same
		return bytes.Compare(lft.TxHash[:], rgt.TxHash[:]) < 0
	} else {
		return v > 0 // by fee otherwise in descending order
	}
}

func (lft *StandardMessage) CompareFee(rgt *StandardMessage) bool {
	lftFrom, rgtFrom := lft.Native.From(), rgt.Native.From()
	if bytes.Compare(lftFrom[:], rgtFrom[:]) == 0 { // by nonce if from the same address
		return lft.Native.Nonce() < rgt.Native.Nonce()
	}

	if v := lft.Native.Fee().Cmp(rgt.Native.Fee()); v == 0 { // by address if fees are the same
		return bytes.Compare(lft.TxHash[:], rgt.TxHash[:]) < 0
	} else {
		return v > 0 // by fee otherwise in descending order
	}
}

type byFee []*StandardMessage

func (stdMsgs byFee) Len() int      { return len(stdMsgs) }
func (stdMsgs byFee) Swap(i, j int) { stdMsgs[i], stdMsgs[j] = stdMsgs[j], stdMsgs[i] }
func (stdMsgs byFee) Less(i, j int) bool {
	return stdMsgs[i].CompareFee(stdMsgs[j])
}

type byGas []*StandardMessage

func (stdMsgs byGas) Len() int      { return len(stdMsgs) }
func (stdMsgs byGas) Swap(i, j int) { stdMsgs[i], stdMsgs[j] = stdMsgs[j], stdMsgs[i] }
func (stdMsgs byGas) Less(i, j int) bool {
	return stdMsgs[i].CompareGas(stdMsgs[j])
}

type byHash []*StandardMessage

func (stdMsgs byHash) Len() int      { return len(stdMsgs) }
func (stdMsgs byHash) Swap(i, j int) { stdMsgs[i], stdMsgs[j] = stdMsgs[j], stdMsgs[i] }
func (stdMsgs byHash) Less(i, j int) bool {
	return stdMsgs[i].CompareHash(stdMsgs[j])
}

type SendingStandardMessages struct {
	Data [][]byte
}

func (stdMsgs SendingStandardMessages) Encode() ([]byte, error) {
	return encoding.Byteset(stdMsgs.Data).Encode(), nil
}
func (stdMsgs *SendingStandardMessages) Decode(data []byte) error {
	stdMsgs.Data = encoding.Byteset{}.Decode(data)
	return nil
}

func (stdMsgs *SendingStandardMessages) ToMessages() []*StandardMessage {
	fields := stdMsgs.Data
	msgs := make([]*StandardMessage, len(fields))

	worker := func(start, end, idx int, args ...interface{}) {
		data := args[0].([]interface{})[0].([][]byte)
		messages := args[0].([]interface{})[1].([]*StandardMessage)

		for i := start; i < end; i++ {
			standredMessage := new(StandardMessage)

			fields := encoding.Byteset{}.Decode(data[i])
			standredMessage.TxHash = ethCommon.BytesToHash(fields[0])
			standredMessage.Source = uint8(fields[1][0])
			msg := new(ethTypes.Message)
			err := msg.GobDecode(fields[2])
			if err != nil {
				return
			}
			standredMessage.Native = msg
			standredMessage.TxRawData = fields[3]

			messages[i] = standredMessage
		}
	}
	common.ParallelWorker(len(fields), concurrency, worker, fields, msgs)

	return msgs
}

type StandardMessages []*StandardMessage

func (stdMsgs StandardMessages) Hashes() []ethCommon.Hash {
	hashes := make([]ethCommon.Hash, len(stdMsgs))
	for i := range stdMsgs {
		hashes[i] = stdMsgs[i].TxHash
	}
	return hashes
}

func (stdMsgs StandardMessages) SortByFee() {
	sort.Sort(byFee(stdMsgs))
}

func (stdMsgs StandardMessages) SortByGas() {
	sort.Sort(byGas(stdMsgs))
}

func (stdMsgs StandardMessages) SortByHash() {
	sort.Sort(byHash(stdMsgs))
}

func (stdMsgs StandardMessages) Count(value *StandardMessage) int {
	counter := 0
	for i := range stdMsgs {
		if bytes.Equal(stdMsgs[i].TxHash[:], value.TxHash[:]) {
			counter++
		}
	}
	return counter
}

func (stdMsgs StandardMessages) QuickSort(less func(lft *StandardMessage, rgt *StandardMessage) bool) {
	if len(stdMsgs) < 2 {
		return
	}
	left, right := 0, len(stdMsgs)-1
	pivotIndex := rand.Int() % len(stdMsgs)

	stdMsgs[pivotIndex], stdMsgs[right] = stdMsgs[right], stdMsgs[pivotIndex]
	for i := range stdMsgs {
		if less(stdMsgs[i], stdMsgs[right]) {
			stdMsgs[i], stdMsgs[left] = stdMsgs[left], stdMsgs[i]
			left++
		}
	}
	stdMsgs[left], stdMsgs[right] = stdMsgs[right], stdMsgs[left]

	StandardMessages(stdMsgs[:left]).QuickSort(less)
	StandardMessages(stdMsgs[left+1:]).QuickSort(less)
}

func (stdMsgs StandardMessages) EncodeToBytes() [][]byte {
	if stdMsgs == nil {
		return [][]byte{}
	}
	data := make([][]byte, len(stdMsgs))
	worker := func(start, end, idx int, args ...interface{}) {
		stdMsgs := args[0].([]interface{})[0].(StandardMessages)
		data := args[0].([]interface{})[1].([][]byte)

		for i := start; i < end; i++ {
			if encoded, err := stdMsgs[i].Native.GobEncode(); err == nil {
				tmpData := [][]byte{
					stdMsgs[i].TxHash.Bytes(),
					[]byte{stdMsgs[i].Source},
					encoded,
					//stdMsgs[i].TxRawData,
					[]byte{}, //remove TxRawData
				}
				data[i] = encoding.Byteset(tmpData).Encode()
			}
		}
	}
	common.ParallelWorker(len(stdMsgs), concurrency, worker, stdMsgs, data)
	return data
}

func (stdMsgs StandardMessages) Encode() ([]byte, error) {
	if stdMsgs == nil {
		return []byte{}, nil
	}
	data := make([][]byte, len(stdMsgs))
	worker := func(start, end, idx int, args ...interface{}) {
		stdMsgs := args[0].([]interface{})[0].(StandardMessages)
		data := args[0].([]interface{})[1].([][]byte)

		for i := start; i < end; i++ {
			if encoded, err := stdMsgs[i].Native.GobEncode(); err == nil {
				//data[i] = encoding.Byteset([][]byte{stdMsgs[i].TxHash.Bytes()[:], {stdMsgs[i].Source}, encoded}).Flatten()
				tmpData := [][]byte{
					stdMsgs[i].TxHash.Bytes(),
					[]byte{stdMsgs[i].Source},
					encoded,
					stdMsgs[i].TxRawData,
				}
				data[i] = encoding.Byteset(tmpData).Encode()
			}
		}
	}
	common.ParallelWorker(len(stdMsgs), concurrency, worker, stdMsgs, data)
	return encoding.Byteset(data).Encode(), nil
}

func (stdMsgs *StandardMessages) Decode(data []byte) ([]*StandardMessage, error) {
	fields := encoding.Byteset{}.Decode(data)
	msgs := make([]*StandardMessage, len(fields))

	worker := func(start, end, idx int, args ...interface{}) {
		data := args[0].([]interface{})[0].([][]byte)
		messages := args[0].([]interface{})[1].([]*StandardMessage)

		for i := start; i < end; i++ {
			standredMessage := new(StandardMessage)

			fields := encoding.Byteset{}.Decode(data[i])
			standredMessage.TxHash = ethCommon.BytesToHash(fields[0])
			standredMessage.Source = uint8(fields[1][0])
			msg := new(ethTypes.Message)
			err := msg.GobDecode(fields[2])
			if err != nil {
				return
			}
			standredMessage.Native = msg
			standredMessage.TxRawData = fields[3]

			messages[i] = standredMessage
		}
	}
	common.ParallelWorker(len(fields), concurrency, worker, fields, msgs)

	return msgs, nil
}
