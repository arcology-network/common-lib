package types

import (
	"math/big"
	"time"
)

const (
	QueryType_Receipt      = "receipt"
	QueryType_Balance      = "balance"
	QueryType_Container    = "container"
	QueryType_Block        = "block"
	QueryType_RawBlock     = "rawblock"
	QueryType_Nonce        = "nonce"
	QueryType_LatestHeight = "latestheight"

	ConcurrentLibStyle_Array = "array"
	ConcurrentLibStyle_Map   = "map"
	ConcurrentLibStyle_Queue = "queue"
)

type QueryRequest struct {
	QueryType string
	Data      interface{}
}

type QueryResult struct {
	Data interface{}
}

type RequestBalance struct {
	Height  int
	Address string
}

type RequestContainer struct {
	Height  int
	Address string
	Id      string
	Style   string
	Key     string
}

type RequestBlock struct {
	Height       int
	Transactions bool
}

type RequestReceipt struct {
	Hashes             []string
	ExecutingDebugLogs bool
}

type Block struct {
	Height       int           `json:"height"`
	Hash         string        `json:"hash"`
	Coinbase     string        `json:"coinbase"`
	Number       int           `json:"number"`
	Transactions []string      `json:"transactions"`
	GasUsed      *big.Int      `json:"gasUsed"`
	ExecTime     time.Duration `json:"elapsedTime"`
	Timestamp    int           `json:"timestamp"`
}

type Log struct {
	Address     string   `json:"address"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	BlockNumber uint64   `json:"blockNumber"`
	TxHash      string   `json:"transactionHash"`
	TxIndex     uint     `json:"transactionIndex"`
	BlockHash   string   `json:"blockHash"`
	Index       uint     `json:"logIndex"`
}

type Receipt struct {
	Status          int      `json:"status"`
	ContractAddress string   `json:"contractAddress"`
	GasUsed         *big.Int `json:"gasUsed"`
	Logs            []*Log   `json:"logs"`
	ExecutingLogs   string   `json:"executing logs"`
	SpawnedTxHash   string   `json:"spawned transactionHash"`
	Height          int      `json:"height"`
}

type ClusterConfig struct {
	Parallelism int
}

type SetReply struct {
	Status int
}

type SendTransactionArgs struct {
	Txs [][]byte
}

type SendTransactionReply struct {
	Status int
}
