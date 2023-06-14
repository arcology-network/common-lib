package types

import (
	evmCommon "github.com/arcology-network/evm/common"
)

const (
	Conclusion_Success  = byte(0)
	Conclusion_Failed   = byte(1)
	Conclusion_Rollback = byte(2)

	TxType_Transfer      = byte(0)
	TxType_SmartContract = byte(1)

	ListType_Arbitrator = byte(0)
	ListType_Sproc      = byte(1)

	SelectedMode_FillTx       = byte(0) //default
	SelectedMode_FillReceipts = byte(1)
	SelectedMode_All          = byte(2)
)

type ApplyListItem struct {
	Txhash     evmCommon.Hash
	Conclusion byte
	TxType     byte
}

type ApplyList struct {
	ListType     byte
	Lists        *[]*ApplyListItem
	SelectedMode byte
}

// Range calls f on each key and value present in the map.
func (al ApplyList) Size() int {
	return 0
}

// Range calls f on each key and value present in the map.
func (al ApplyList) Range(f func(hash *evmCommon.Hash, selected bool)) {
	if al.Lists != nil {
		switch al.SelectedMode {
		case SelectedMode_FillTx:
			for _, ali := range *al.Lists {
				if ali == nil {
					continue
				}
				if ali.Conclusion == Conclusion_Success {
					f(&ali.Txhash, true)
				} else if ali.Conclusion == Conclusion_Failed {
					f(&ali.Txhash, false)
				} else {
					f(nil, false)
				}
			}
		case SelectedMode_FillReceipts:
			for _, ali := range *al.Lists {
				if ali == nil {
					continue
				}
				if ali.Conclusion == Conclusion_Success {
					f(&ali.Txhash, true)
				} else {
					f(&ali.Txhash, false)
				}
			}
		case SelectedMode_All:
			for _, ali := range *al.Lists {
				if ali == nil {
					continue
				}
				f(&ali.Txhash, true)
			}
		}

	}

}
