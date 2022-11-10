package types

type IncomingMsgs struct {
	Msgs []*StandardMessage
	Src  TxSource
}
