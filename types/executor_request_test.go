package types

import (
	"fmt"
	"math/big"
	"testing"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
	ethTypes "github.com/arcology-network/3rd-party/eth/types"
)

func TestExecutorRequestEncodingAndDeconing(t *testing.T) {
	from := ethCommon.BytesToAddress([]byte{0, 0, 0, 5, 6, 7, 8, 9})
	to := ethCommon.BytesToAddress([]byte{11, 12, 13})
	from0 := ethCommon.BytesToAddress([]byte{0, 0, 0, 1, 4, 5, 6, 7, 8})
	to0 := ethCommon.BytesToAddress([]byte{11, 12, 13, 14})
	ethMsg_serial_0 := ethTypes.NewMessage(from0, &to0, 1, big.NewInt(int64(1)), 100, big.NewInt(int64(8)), []byte{1, 2, 3}, false)
	ethMsg_serial_1 := ethTypes.NewMessage(from, &to, 3, big.NewInt(int64(100)), 200, big.NewInt(int64(9)), []byte{4, 5, 6}, false)
	fmt.Printf("ethMsg_serial_0=%v\n", ethMsg_serial_0)
	fmt.Printf("ethMsg_serial_1=%v\n", ethMsg_serial_1)
	hash1 := ethCommon.RlpHash(ethMsg_serial_0)
	hash2 := ethCommon.RlpHash(ethMsg_serial_1)
	fmt.Printf("hash1=%v\n", hash1)
	fmt.Printf("hash2=%v\n", hash2)
	// msgs1 := []*StandardMessage{
	// 	{Source: 0, Native: &ethMsg_serial_0, TxHash: hash1},
	// 	{Source: 1, Native: &ethMsg_serial_1, TxHash: hash2},
	// }
	// sequence1 := NewExecutingSequence(msgs1, true)
	/*
		from1 := ethCommon.BytesToAddress([]byte{9, 9, 9})
		to1 := ethCommon.BytesToAddress([]byte{21, 22, 23})
		ethMsg_serial_2 := ethTypes.NewMessage(from1, &to1, 2, big.NewInt(int64(11)), 300, big.NewInt(int64(18)), []byte{7, 8, 9}, false)
		ethMsg_serial_3 := ethTypes.NewMessage(from1, &to1, 4, big.NewInt(int64(110)), 400, big.NewInt(int64(91)), []byte{10, 15, 16}, false)
		hash3 := ethCommon.RlpHash(ethMsg_serial_2)
		hash4 := ethCommon.RlpHash(ethMsg_serial_3)
		msgs2 := []*StandardMessage{
			{Source: 0, Native: &ethMsg_serial_2, TxHash: hash3},
			{Source: 1, Native: &ethMsg_serial_3, TxHash: hash4},
		}
		sequence2 := NewExecutingSequence(msgs2, false)

		hashes := []ethCommon.Hash{
			ethCommon.BytesToHash([]byte{1, 2, 3}),
			ethCommon.BytesToHash([]byte{4, 5, 6}),
		}
		er := ExecutorRequest{
			Sequences:     []*ExecutingSequence{sequence1, sequence2},
			Precedings:    Arr2Ptr(hashes),
			PrecedingHash: ethCommon.BytesToHash([]byte{11, 12, 13}),
			Timestamp:     big.NewInt(12),
			Parallelism:   4,
		}

		for _, sequence := range er.Sequences {
			fmt.Printf("ExecutorRequest.er sequence=%v\n", sequence)
		}

		data, err := er.GobEncode()
		if err != nil {
			fmt.Printf("ExecutorRequest encode err=%v\n", err)
			return
		}
		fmt.Printf("ExecutorRequest encode result=%v\n", data)

		erResult := ExecutorRequest{}
		err = erResult.GobDecode(data)
		if err != nil {
			fmt.Printf(" ExecutorRequest.GobDecode err=%v\n", err)
			return

		}
		fmt.Printf(" ExecutorRequest.GobDecode result=%v\n", erResult)

		for _, sequence := range erResult.Sequences {
			fmt.Printf("ExecutorRequest.GobDecode sequence=%v\n", sequence)
		}
	*/
}
