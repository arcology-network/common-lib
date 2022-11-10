package types

import (
	"fmt"
	"strings"
)

const (
	TxSourceLocal = iota
	TxSourceConsensus
	TxSourceMonacoP2p
)

var (
	SourceTypeStr = map[int]string{
		TxSourceLocal:     "loc",
		TxSourceConsensus: "con",
		TxSourceMonacoP2p: "p2p",
	}
)

type TxSource string

func NewTxSource(typ int, id string) TxSource {
	return TxSource(fmt.Sprintf("%s:%s", SourceTypeStr[typ], id))
}

func (src TxSource) BypassRepeatCheck() bool {
	return strings.HasPrefix(string(src), SourceTypeStr[TxSourceMonacoP2p])
}

func (src TxSource) IsForWaitingList() bool {
	return strings.HasPrefix(string(src), SourceTypeStr[TxSourceMonacoP2p])
}

type IncomingTxs struct {
	Txs [][]byte
	Src TxSource
}
