package types

import (
	"bytes"
	"math/big"
	"math/rand"
	"sort"

	codec "github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	evmTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	Concurrency = 4
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
	return codec.Byteset(data).Encode(), nil
}

func (stdp *StdTransactionPack) GobDecode(data []byte) error {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	stds := StandardTransactions{}
	stds, err := stds.Decode(fields[0])
	if err != nil {
		panic("StandardMessages decode failed")
	}
	stdp.Txs = stds
	stdp.Src = TxSource(fields[1])
	return nil
}

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

	pivotIndex := rand.Intn(len(this)) //rnd.Int() % len(this)

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
			data[i] = codec.Byteset(tmpData).Encode()

		}
	}
	common.ParallelWorker(len(this), Concurrency, worker, this, data)
	return codec.Byteset(data).Encode(), nil
}

func (this *StandardTransactions) Decode(data []byte) ([]*StandardTransaction, error) {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	msgs := make([]*StandardTransaction, len(fields))

	worker := func(start, end, idx int, args ...interface{}) {
		data := args[0].([]interface{})[0].([][]byte)
		messages := args[0].([]interface{})[1].([]*StandardTransaction)

		for i := start; i < end; i++ {
			standredMessage := new(StandardTransaction)

			fields := codec.Byteset{}.Decode(data[i]).(codec.Byteset)
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
	common.ParallelWorker(len(fields), Concurrency, worker, fields, msgs)

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
