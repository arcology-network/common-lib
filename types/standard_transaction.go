package types

import (
	"bytes"
	"math/big"
	"math/rand"
	"sort"

	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/encoding"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	evmTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	concurrency = 4
)

const (
	TxType_Eth  = 0
	TxType_Coin = 1

	TxFrom_Remote = 1
	TxFrom_Local  = 2
	TxFrom_Block  = 3
)

const (
	TX_SOURCE_REMOTE = iota
	TX_SOURCE_LOCAL
	TX_SOURCE_BLOCK
	TX_SOURCE_DEFERRED
)

type StandardTransaction struct {
	TxHash            ethCommon.Hash
	NativeMessage     *core.Message
	NativeTransaction *evmTypes.Transaction
	TxRawData         []byte
	Source            uint8
	Signer            uint8
}

type StdTransactionPack struct {
	Txs        StandardTransactions
	Src        TxSource
	TxHashChan chan ethCommon.Hash
}

func (stdp *StdTransactionPack) GobEncode() ([]byte, error) {
	txsData, err := stdp.Txs.Encode()
	if err != nil {
		panic("StandardMessages encode failed")
	}
	data := [][]byte{
		txsData,
		[]byte(stdp.Src),
	}
	return encoding.Byteset(data).Encode(), nil
}
func (stdp *StdTransactionPack) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	stds := StandardTransactions{}
	stds, err := stds.Decode(fields[0])
	if err != nil {
		panic("StandardMessages decode failed")
	}
	stdp.Txs = stds
	stdp.Src = TxSource(fields[1])
	return nil
}

// func MakeMessageWithDefCall(def *DeferredCall, hash ethCommon.Hash, nonce uint64) *StandardMessage {
// 	signature := def.Signature
// 	contractAddress := def.ContractAddress
// 	data := crypto.Keccak256([]byte(signature))[:4]
// 	data = append(data, common.AlignToEvmForInt(common.EvmWordSize)...)
// 	idLen := common.AlignToEvmForInt(len(def.DeferID))
// 	id := common.AlignToEvmForString(def.DeferID)
// 	data = append(data, idLen...)
// 	data = append(data, id...)
// 	contractAddr := ethCommon.BytesToAddress([]byte(contractAddress))
// 	//nonce := uint64(time.Now().UnixNano())
// 	message := core.NewMessage(contractAddr, &contractAddr, nonce, new(big.Int).SetInt64(0), 1e9, new(big.Int).SetInt64(0), data, nil, false)
// 	standardMessager := StandardMessage{
// 		Native: &message,
// 		TxHash: hash,
// 	}
// 	return &standardMessager
// }

func (this *StandardTransaction) Hash() ethCommon.Hash {
	return this.TxHash
}

func (this *StandardTransaction) UnSign(signer evmTypes.Signer) error {
	otx := this.NativeTransaction
	msg, err := core.TransactionToMessage(otx, signer, nil)
	if err != nil {
		return err
	}
	// msg.SkipAccountChecks = true
	this.NativeMessage = msg
	return nil
}

func (this *StandardTransaction) Key() string {
	return this.TxHash.String()
}

func (this *StandardTransaction) Equal(other *StandardTransaction) bool {
	return this.TxHash.String() == other.TxHash.String()
}

func (this *StandardTransaction) CompareHash(rgt *StandardTransaction) bool {
	return bytes.Compare(this.TxHash[:], rgt.TxHash[:]) < 0
}

func (this *StandardTransaction) CompareGas(rgt *StandardTransaction) bool {
	lftFrom, rgtFrom := this.NativeMessage.From, rgt.NativeMessage.From
	if bytes.Compare(lftFrom[:], rgtFrom[:]) == 0 { // by nonce if from the same address
		return this.NativeMessage.Nonce < rgt.NativeMessage.Nonce
	}

	if v := this.NativeMessage.GasPrice.Cmp(rgt.NativeMessage.GasPrice); v == 0 { // by address if fees are the same
		return bytes.Compare(this.TxHash[:], rgt.TxHash[:]) < 0
	} else {
		return v > 0 // by fee otherwise in descending order
	}
}

func (this *StandardTransaction) CompareFee(rgt *StandardTransaction) bool {
	lftFrom, rgtFrom := this.NativeMessage.From, rgt.NativeMessage.From
	if bytes.Compare(lftFrom[:], rgtFrom[:]) == 0 { // by nonce if from the same address
		return this.NativeMessage.Nonce < rgt.NativeMessage.Nonce
	}

	if v := MsgFee(this.NativeMessage).Cmp(MsgFee(rgt.NativeMessage)); v == 0 { // by address if fees are the same
		return bytes.Compare(this.TxHash[:], rgt.TxHash[:]) < 0
	} else {
		return v > 0 // by fee otherwise in descending order
	}
}

type byFee []*StandardTransaction

func (this byFee) Len() int      { return len(this) }
func (this byFee) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this byFee) Less(i, j int) bool {
	return this[i].CompareFee(this[j])
}

type byGas []*StandardTransaction

func (this byGas) Len() int      { return len(this) }
func (this byGas) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this byGas) Less(i, j int) bool {
	return this[i].CompareGas(this[j])
}

type byHash []*StandardTransaction

func (this byHash) Len() int      { return len(this) }
func (this byHash) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this byHash) Less(i, j int) bool {
	return this[i].CompareHash(this[j])
}

// type SendingStandardMessages struct {
// 	Data [][]byte
// }

// func (this SendingStandardMessages) Encode() ([]byte, error) {
// 	return encoding.Byteset(this.Data).Encode(), nil
// }
// func (this *SendingStandardMessages) Decode(data []byte) error {
// 	this.Data = encoding.Byteset{}.Decode(data)
// 	return nil
// }

// func (this *SendingStandardMessages) ToMessages() []*StandardMessage {
// 	fields := this.Data
// 	msgs := make([]*StandardMessage, len(fields))

// 	worker := func(start, end, idx int, args ...interface{}) {
// 		data := args[0].([]interface{})[0].([][]byte)
// 		messages := args[0].([]interface{})[1].([]*StandardMessage)

// 		for i := start; i < end; i++ {
// 			standredMessage := new(StandardMessage)

// 			fields := encoding.Byteset{}.Decode(data[i])
// 			standredMessage.TxHash = ethCommon.BytesToHash(fields[0])
// 			standredMessage.Source = uint8(fields[1][0])

// 			// msg := new(core.Message)
// 			// err := msg.GobDecode(fields[2])
// 			msg, err := MsgDecode(fields[2])
// 			if err != nil {
// 				fmt.Printf("SendingStandardMessages decode err:%v", err)
// 				return
// 			}
// 			standredMessage.NativeMessage = msg
// 			standredMessage.TxRawData = fields[3]

// 			messages[i] = standredMessage
// 		}
// 	}
// 	common.ParallelWorker(len(fields), concurrency, worker, fields, msgs)

// 	return msgs
// }

type StandardTransactions []*StandardTransaction

func (this StandardTransactions) Hashes() []ethCommon.Hash {
	hashes := make([]ethCommon.Hash, len(this))
	for i := range this {
		hashes[i] = this[i].TxHash
	}
	return hashes
}

func (this StandardTransactions) SortByFee() {
	sort.Sort(byFee(this))
}

func (this StandardTransactions) SortByGas() {
	sort.Sort(byGas(this))
}

func (this StandardTransactions) SortByHash() {
	sort.Sort(byHash(this))
}

func (this StandardTransactions) Count(value *StandardTransaction) int {
	counter := 0
	for i := range this {
		if bytes.Equal(this[i].TxHash[:], value.TxHash[:]) {
			counter++
		}
	}
	return counter
}

func (this StandardTransactions) QuickSort(less func(this *StandardTransaction, rgt *StandardTransaction) bool) {
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

	StandardTransactions(this[:left]).QuickSort(less)
	StandardTransactions(this[left+1:]).QuickSort(less)
}

// func (this StandardMessages) EncodeToBytes() [][]byte {
// 	if this == nil {
// 		return [][]byte{}
// 	}
// 	data := make([][]byte, len(this))
// 	worker := func(start, end, idx int, args ...interface{}) {
// 		this := args[0].([]interface{})[0].(StandardMessages)
// 		data := args[0].([]interface{})[1].([][]byte)

// 		for i := start; i < end; i++ {
// 			if encoded, err := MsgEncode(this[i].NativeMessage); err == nil {
// 				tmpData := [][]byte{
// 					this[i].TxHash.Bytes(),
// 					[]byte{this[i].Source},
// 					encoded,
// 					//this[i].TxRawData
// 					[]byte{}, //remove TxRawData
// 				}
// 				data[i] = encoding.Byteset(tmpData).Encode()
// 			}
// 		}
// 	}
// 	common.ParallelWorker(len(this), concurrency, worker, this, data)
// 	return data
// }

func (this StandardTransactions) Encode() ([]byte, error) {
	if this == nil {
		return []byte{}, nil
	}
	data := make([][]byte, len(this))
	worker := func(start, end, idx int, args ...interface{}) {
		this := args[0].([]interface{})[0].(StandardTransactions)
		data := args[0].([]interface{})[1].([][]byte)

		for i := start; i < end; i++ {
			encodedMsg := []byte{}
			if encoded, err := MsgEncode(this[i].NativeMessage); err == nil {
				encodedMsg = encoded
			}
			encodedTx := []byte{}
			if encoded, err := TxEncode(this[i].NativeTransaction); err == nil {
				encodedTx = encoded
			}
			tmpData := [][]byte{
				this[i].TxHash.Bytes(),
				[]byte{this[i].Source},
				encodedMsg,
				encodedTx,
				this[i].TxRawData,
				[]byte{this[i].Signer},
			}
			data[i] = encoding.Byteset(tmpData).Encode()

		}
	}
	common.ParallelWorker(len(this), concurrency, worker, this, data)
	return encoding.Byteset(data).Encode(), nil
}

func (this *StandardTransactions) Decode(data []byte) ([]*StandardTransaction, error) {
	fields := encoding.Byteset{}.Decode(data)
	msgs := make([]*StandardTransaction, len(fields))

	worker := func(start, end, idx int, args ...interface{}) {
		data := args[0].([]interface{})[0].([][]byte)
		messages := args[0].([]interface{})[1].([]*StandardTransaction)

		for i := start; i < end; i++ {
			standredMessage := new(StandardTransaction)

			fields := encoding.Byteset{}.Decode(data[i])
			standredMessage.TxHash = ethCommon.BytesToHash(fields[0])
			standredMessage.Source = uint8(fields[1][0])
			if len(fields[2]) > 0 {
				msg, err := MsgDecode(fields[2])
				if err != nil {
					return
				}
				standredMessage.NativeMessage = msg
			}
			if len(fields[3]) > 0 {
				tx, err := TxDecode(fields[3])
				if err != nil {
					return
				}
				standredMessage.NativeTransaction = tx
			}
			standredMessage.TxRawData = fields[4]
			standredMessage.Signer = uint8(fields[5][0])
			messages[i] = standredMessage
		}
	}
	common.ParallelWorker(len(fields), concurrency, worker, fields, msgs)

	return msgs, nil
}

//---------------------------------------------------

func EntrySignature(m core.Message) string {
	if len(m.Data) >= 4 {
		return string(m.Data[:4])
	}
	return ""
}
func MsgEncode(m *core.Message) ([]byte, error) {
	return rlp.EncodeToBytes(m)
}
func MsgDecode(data []byte) (*core.Message, error) {
	m := core.Message{}
	return &m, rlp.DecodeBytes(data, &m)
}

func TxEncode(tx *evmTypes.Transaction) ([]byte, error) {
	return tx.MarshalBinary()
}
func TxDecode(data []byte) (*evmTypes.Transaction, error) {
	tx := evmTypes.Transaction{}
	return &tx, tx.UnmarshalBinary(data)
}

func MsgFee(m *core.Message) *big.Int {
	return big.NewInt(0).Mul(big.NewInt(int64(m.GasLimit)), m.GasPrice)
} // Max fee possible
