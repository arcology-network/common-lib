package types

import (
	"bytes"
	"math/big"
	"math/rand"
	"sort"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
	ethTypes "github.com/arcology-network/3rd-party/eth/types"
	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/encoding"
	"github.com/arcology-network/evm/crypto"
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

func MakeMessageWithDefCall(def *DeferredCall, hash ethCommon.Hash, nonce uint64) *StandardMessage {
	signature := def.Signature
	contractAddress := def.ContractAddress
	data := crypto.Keccak256([]byte(signature))[:4]
	data = append(data, ethCommon.AlignToEvmForInt(ethCommon.EvmWordSize)...)
	idLen := ethCommon.AlignToEvmForInt(len(def.DeferID))
	id := ethCommon.AlignToEvmForString(def.DeferID)
	data = append(data, idLen...)
	data = append(data, id...)
	contractAddr := ethCommon.BytesToAddress([]byte(contractAddress))
	//nonce := uint64(time.Now().UnixNano())
	message := ethTypes.NewMessage(contractAddr, &contractAddr, nonce, new(big.Int).SetInt64(0), 1e9, new(big.Int).SetInt64(0), data, false)
	standardMessager := StandardMessage{
		Native: &message,
		TxHash: hash,
	}
	return &standardMessager
}

func (this *StandardMessage) Hash() ethCommon.Hash {
	return this.TxHash
}

func (this *StandardMessage) Key() string {
	return this.TxHash.String()
}

func (this *StandardMessage) Equal(other *StandardMessage) bool {
	return this.TxHash.String() == other.TxHash.String()
}

func (this *StandardMessage) CompareHash(rgt *StandardMessage) bool {
	return bytes.Compare(this.TxHash[:], rgt.TxHash[:]) < 0
}

func (this *StandardMessage) CompareGas(rgt *StandardMessage) bool {
	lftFrom, rgtFrom := this.Native.From(), rgt.Native.From()
	if bytes.Compare(lftFrom[:], rgtFrom[:]) == 0 { // by nonce if from the same address
		return this.Native.Nonce() < rgt.Native.Nonce()
	}

	if v := this.Native.GasPrice().Cmp(rgt.Native.GasPrice()); v == 0 { // by address if fees are the same
		return bytes.Compare(this.TxHash[:], rgt.TxHash[:]) < 0
	} else {
		return v > 0 // by fee otherwise in descending order
	}
}

func (this *StandardMessage) CompareFee(rgt *StandardMessage) bool {
	lftFrom, rgtFrom := this.Native.From(), rgt.Native.From()
	if bytes.Compare(lftFrom[:], rgtFrom[:]) == 0 { // by nonce if from the same address
		return this.Native.Nonce() < rgt.Native.Nonce()
	}

	if v := this.Native.Fee().Cmp(rgt.Native.Fee()); v == 0 { // by address if fees are the same
		return bytes.Compare(this.TxHash[:], rgt.TxHash[:]) < 0
	} else {
		return v > 0 // by fee otherwise in descending order
	}
}

type byFee []*StandardMessage

func (this byFee) Len() int      { return len(this) }
func (this byFee) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this byFee) Less(i, j int) bool {
	return this[i].CompareFee(this[j])
}

type byGas []*StandardMessage

func (this byGas) Len() int      { return len(this) }
func (this byGas) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this byGas) Less(i, j int) bool {
	return this[i].CompareGas(this[j])
}

type byHash []*StandardMessage

func (this byHash) Len() int      { return len(this) }
func (this byHash) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this byHash) Less(i, j int) bool {
	return this[i].CompareHash(this[j])
}

type SendingStandardMessages struct {
	Data [][]byte
}

func (this SendingStandardMessages) Encode() ([]byte, error) {
	return encoding.Byteset(this.Data).Encode(), nil
}
func (this *SendingStandardMessages) Decode(data []byte) error {
	this.Data = encoding.Byteset{}.Decode(data)
	return nil
}

func (this *SendingStandardMessages) ToMessages() []*StandardMessage {
	fields := this.Data
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

func (this StandardMessages) Hashes() []ethCommon.Hash {
	hashes := make([]ethCommon.Hash, len(this))
	for i := range this {
		hashes[i] = this[i].TxHash
	}
	return hashes
}

func (this StandardMessages) SortByFee() {
	sort.Sort(byFee(this))
}

func (this StandardMessages) SortByGas() {
	sort.Sort(byGas(this))
}

func (this StandardMessages) SortByHash() {
	sort.Sort(byHash(this))
}

func (this StandardMessages) Count(value *StandardMessage) int {
	counter := 0
	for i := range this {
		if bytes.Equal(this[i].TxHash[:], value.TxHash[:]) {
			counter++
		}
	}
	return counter
}

func (this StandardMessages) QuickSort(less func(this *StandardMessage, rgt *StandardMessage) bool) {
	if len(this) < 2 {
		return
	}
	left, right := 0, len(this)-1
	pivotIndex := rand.Int() % len(this)

	this[pivotIndex], this[right] = this[right], this[pivotIndex]
	for i := range this {
		if less(this[i], this[right]) {
			this[i], this[left] = this[left], this[i]
			left++
		}
	}
	this[left], this[right] = this[right], this[left]

	StandardMessages(this[:left]).QuickSort(less)
	StandardMessages(this[left+1:]).QuickSort(less)
}

func (this StandardMessages) EncodeToBytes() [][]byte {
	if this == nil {
		return [][]byte{}
	}
	data := make([][]byte, len(this))
	worker := func(start, end, idx int, args ...interface{}) {
		this := args[0].([]interface{})[0].(StandardMessages)
		data := args[0].([]interface{})[1].([][]byte)

		for i := start; i < end; i++ {
			if encoded, err := this[i].Native.GobEncode(); err == nil {
				tmpData := [][]byte{
					this[i].TxHash.Bytes(),
					[]byte{this[i].Source},
					encoded,
					//this[i].TxRawData,
					[]byte{}, //remove TxRawData
				}
				data[i] = encoding.Byteset(tmpData).Encode()
			}
		}
	}
	common.ParallelWorker(len(this), concurrency, worker, this, data)
	return data
}

func (this StandardMessages) Encode() ([]byte, error) {
	if this == nil {
		return []byte{}, nil
	}
	data := make([][]byte, len(this))
	worker := func(start, end, idx int, args ...interface{}) {
		this := args[0].([]interface{})[0].(StandardMessages)
		data := args[0].([]interface{})[1].([][]byte)

		for i := start; i < end; i++ {
			if encoded, err := this[i].Native.GobEncode(); err == nil {
				//data[i] = encoding.Byteset([][]byte{this[i].TxHash.Bytes()[:], {this[i].Source}, encoded}).Flatten()
				tmpData := [][]byte{
					this[i].TxHash.Bytes(),
					[]byte{this[i].Source},
					encoded,
					this[i].TxRawData,
				}
				data[i] = encoding.Byteset(tmpData).Encode()
			}
		}
	}
	common.ParallelWorker(len(this), concurrency, worker, this, data)
	return encoding.Byteset(data).Encode(), nil
}

func (this *StandardMessages) Decode(data []byte) ([]*StandardMessage, error) {
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
