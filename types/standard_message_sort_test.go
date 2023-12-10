package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	evmCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestStandardMessageEncodingAndDeconing(t *testing.T) {
	to := evmCommon.BytesToAddress(crypto.Keccak256([]byte("1"))[:])

	ethMsg_serial_0 := core.NewMessage(evmCommon.Address{}, &to, 1, big.NewInt(int64(1)), 100, big.NewInt(int64(8)), []byte{1, 2, 3}, nil, false)
	ethMsg_serial_1 := core.NewMessage(evmCommon.Address{}, &to, 3, big.NewInt(int64(100)), 200, big.NewInt(int64(9)), []byte{4, 5, 6}, nil, false)
	hash1 := RlpHash(ethMsg_serial_0)
	hash2 := RlpHash(ethMsg_serial_1)
	stdMsgs := []*StandardMessage{
		{Source: 0, Native: &ethMsg_serial_0, TxHash: hash1},
		{Source: 1, Native: &ethMsg_serial_1, TxHash: hash2},
	}
	standardMessages := StandardMessages(stdMsgs)
	data, err := standardMessages.Encode()
	if err != nil {
		fmt.Printf("StandardMessages encode err=%v\n", err)
		return
	}
	fmt.Printf("StandardMessages encode result=%v\n", data)

	standardMessages2 := new(StandardMessages)

	standardMessagesResult, err := standardMessages2.Decode(data)

	if err != nil {
		fmt.Printf("StandardMessages dncode err=%v\n", err)
		return
	}
	for _, v := range standardMessagesResult {
		fmt.Printf("StandardMessages dncode result=%v,Native=%v\n", v, v.Native)
	}

}

func TestStandardMessageSortingByFee(t *testing.T) {
	to := evmCommon.BytesToAddress(crypto.Keccak256([]byte("1"))[:])

	ethMsg_serial_0 := core.NewMessage(evmCommon.Address{}, &to, 1, big.NewInt(int64(1)), 100, big.NewInt(int64(8)), []byte{}, nil, false)
	ethMsg_serial_1 := core.NewMessage(evmCommon.Address{}, &to, 3, big.NewInt(int64(100)), 200, big.NewInt(int64(9)), []byte{}, nil, false)
	ethMsg_serial_2 := core.NewMessage(evmCommon.Address{}, &to, 2, big.NewInt(int64(500)), 100, big.NewInt(int64(1)), []byte{}, nil, false)

	from_3 := evmCommon.BytesToAddress([]byte("1"))
	ethMsg_serial_3 := core.NewMessage(from_3, &to, 1, big.NewInt(int64(200)), 500, big.NewInt(int64(8)), []byte{}, nil, false)

	from_4 := evmCommon.BytesToAddress([]byte("2"))
	ethMsg_serial_4 := core.NewMessage(from_4, &to, 1, big.NewInt(int64(10)), 2, big.NewInt(int64(9)), []byte{}, nil, false)

	stdMsgs := []*StandardMessage{
		{Source: 0, Native: &ethMsg_serial_0},
		{Source: 1, Native: &ethMsg_serial_1},
		{Source: 2, Native: &ethMsg_serial_2},
		{Source: 3, Native: &ethMsg_serial_3},
		{Source: 4, Native: &ethMsg_serial_4},
	}

	StandardMessages(stdMsgs).SortByFee()

	if (*stdMsgs[0]).Source != 3 || (*stdMsgs[1]).Source != 0 || (*stdMsgs[2]).Source != 2 || (*stdMsgs[3]).Source != 1 || (*stdMsgs[4]).Source != 4 {
		t.Error("Wrong order")
	}

}

func TestStandardMessageSortingByGas(t *testing.T) {
	to := evmCommon.BytesToAddress(crypto.Keccak256([]byte("1"))[:])

	ethMsg_serial_0 := core.NewMessage(evmCommon.Address{}, &to, 1, big.NewInt(int64(1)), 100, big.NewInt(int64(8)), []byte{}, nil, false)
	ethMsg_serial_1 := core.NewMessage(evmCommon.Address{}, &to, 3, big.NewInt(int64(100)), 200, big.NewInt(int64(9)), []byte{}, nil, false)
	ethMsg_serial_2 := core.NewMessage(evmCommon.Address{}, &to, 2, big.NewInt(int64(500)), 100, big.NewInt(int64(1)), []byte{}, nil, false)

	from_3 := evmCommon.BytesToAddress([]byte("1"))
	ethMsg_serial_3 := core.NewMessage(from_3, &to, 1, big.NewInt(int64(200)), 500, big.NewInt(int64(8)), []byte{}, nil, false)

	from_4 := evmCommon.BytesToAddress([]byte("2"))
	ethMsg_serial_4 := core.NewMessage(from_4, &to, 1, big.NewInt(int64(10)), 2, big.NewInt(int64(9)), []byte{}, nil, false)

	stdMsgs := []*StandardMessage{
		{Source: 0, Native: &ethMsg_serial_0},
		{Source: 1, Native: &ethMsg_serial_1},
		{Source: 2, Native: &ethMsg_serial_2},
		{Source: 3, Native: &ethMsg_serial_3},
		{Source: 4, Native: &ethMsg_serial_4},
	}

	StandardMessages(stdMsgs).SortByGas()

	if (*stdMsgs[0]).Source != 0 || (*stdMsgs[1]).Source != 2 || (*stdMsgs[2]).Source != 1 || (*stdMsgs[3]).Source != 4 || (*stdMsgs[4]).Source != 3 {
		t.Error("Wrong order")
	}
}

func PrepareData(max int) []*StandardMessage {
	stdMsgs := make([]*StandardMessage, max)

	for i := 0; i < len(stdMsgs); i++ {
		to := evmCommon.BytesToAddress([]byte{11, 8, 9, 10})
		bytes := sha256.Sum256([]byte{byte(i)})
		ethMsg := core.NewMessage(
			evmCommon.BytesToAddress(bytes[:]),
			&to,
			uint64(10),
			big.NewInt(12000000),
			uint64(22),
			big.NewInt(34),
			make([]byte, 128),
			nil,
			false,
		)

		stdMsgs[i] = &StandardMessage{
			Source: 1,
			TxHash: bytes,
			Native: &ethMsg,
		}

	}
	return stdMsgs
}

func TestQuickSort(t *testing.T) {
	stdmsg0 := PrepareData(4)
	stdmsg1 := PrepareData(4)

	worker := func(lft *StandardMessage, rgt *StandardMessage) bool {
		return bytes.Compare(lft.TxHash[:], rgt.TxHash[:]) < 0
	}

	StandardMessages(stdmsg0).QuickSort(worker)
	StandardMessages(stdmsg1).SortByHash()
	if !reflect.DeepEqual(stdmsg0, stdmsg1) {
		t.Error("mismatch")
	}
}

func TestQuickSortPerformance(t *testing.T) {
	stdmsgs := PrepareData(500000)
	t0 := time.Now()
	worker := func(lft *StandardMessage, rgt *StandardMessage) bool {
		return bytes.Compare(lft.TxHash[:], rgt.TxHash[:]) < 0
	}
	StandardMessages(stdmsgs).QuickSort(worker)

	fmt.Println("append:", time.Now().Sub(t0))
}
