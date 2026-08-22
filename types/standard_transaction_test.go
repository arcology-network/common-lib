/*
 *   Copyright (c) 2024 Arcology Network

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestStandardMessageSortingByFee(t *testing.T) {
	to := ethCommon.BytesToAddress(crypto.Keccak256([]byte("1"))[:])

	ethMsg_serial_0 := core.Message{From: ethCommon.Address{}, To: &to, Nonce: 1, Value: big.NewInt(1), GasLimit: 100, GasPrice: big.NewInt(8)}
	ethMsg_serial_1 := core.Message{From: ethCommon.Address{}, To: &to, Nonce: 3, Value: big.NewInt(100), GasLimit: 200, GasPrice: big.NewInt(9)}
	ethMsg_serial_2 := core.Message{From: ethCommon.Address{}, To: &to, Nonce: 2, Value: big.NewInt(500), GasLimit: 100, GasPrice: big.NewInt(1)}

	from_3 := ethCommon.BytesToAddress([]byte("1"))
	ethMsg_serial_3 := core.Message{From: from_3, To: &to, Nonce: 1, Value: big.NewInt(200), GasLimit: 500, GasPrice: big.NewInt(8)}

	from_4 := ethCommon.BytesToAddress([]byte("2"))
	ethMsg_serial_4 := core.Message{From: from_4, To: &to, Nonce: 1, Value: big.NewInt(10), GasLimit: 2, GasPrice: big.NewInt(9)}

	stdMsgs := []*StandardTransaction{
		{Source: 0, NativeMessage: &ethMsg_serial_0},
		{Source: 1, NativeMessage: &ethMsg_serial_1},
		{Source: 2, NativeMessage: &ethMsg_serial_2},
		{Source: 3, NativeMessage: &ethMsg_serial_3},
		{Source: 4, NativeMessage: &ethMsg_serial_4},
	}

	StandardTransactions(stdMsgs).SortByFee()

	if (*stdMsgs[0]).Source != 3 || (*stdMsgs[1]).Source != 0 || (*stdMsgs[2]).Source != 2 || (*stdMsgs[3]).Source != 1 || (*stdMsgs[4]).Source != 4 {
		t.Error("Wrong order")
	}

}

func TestStandardMessageSortingByGas(t *testing.T) {
	to := ethCommon.BytesToAddress(crypto.Keccak256([]byte("1"))[:])

	ethMsg_serial_0 := core.Message{From: ethCommon.Address{}, To: &to, Nonce: 1, Value: big.NewInt(1), GasLimit: 100, GasPrice: big.NewInt(8)}
	ethMsg_serial_1 := core.Message{From: ethCommon.Address{}, To: &to, Nonce: 3, Value: big.NewInt(100), GasLimit: 200, GasPrice: big.NewInt(9)}
	ethMsg_serial_2 := core.Message{From: ethCommon.Address{}, To: &to, Nonce: 2, Value: big.NewInt(500), GasLimit: 100, GasPrice: big.NewInt(1)}

	from_3 := ethCommon.BytesToAddress([]byte("1"))
	ethMsg_serial_3 := core.Message{From: from_3, To: &to, Nonce: 1, Value: big.NewInt(200), GasLimit: 500, GasPrice: big.NewInt(8)}

	from_4 := ethCommon.BytesToAddress([]byte("2"))
	ethMsg_serial_4 := core.Message{From: from_4, To: &to, Nonce: 1, Value: big.NewInt(10), GasLimit: 2, GasPrice: big.NewInt(9)}

	stdMsgs := []*StandardTransaction{
		{Source: 0, NativeMessage: &ethMsg_serial_0},
		{Source: 1, NativeMessage: &ethMsg_serial_1},
		{Source: 2, NativeMessage: &ethMsg_serial_2},
		{Source: 3, NativeMessage: &ethMsg_serial_3},
		{Source: 4, NativeMessage: &ethMsg_serial_4},
	}

	StandardTransactions(stdMsgs).SortByGas()

	if (*stdMsgs[0]).Source != 0 || (*stdMsgs[1]).Source != 2 || (*stdMsgs[2]).Source != 1 || (*stdMsgs[3]).Source != 4 || (*stdMsgs[4]).Source != 3 {
		t.Error("Wrong order")
	}
}

func PrepareData(max int) []*StandardTransaction {
	stdMsgs := make([]*StandardTransaction, max)

	for i := 0; i < len(stdMsgs); i++ {
		to := ethCommon.BytesToAddress([]byte{11, 8, 9, 10})
		bytes := sha256.Sum256([]byte{byte(i)})
		ethMsg := core.Message{
			From:     ethCommon.BytesToAddress(bytes[:]),
			To:       &to,
			Nonce:    10,
			Value:    big.NewInt(12000000),
			GasLimit: 22,
			GasPrice: big.NewInt(34),
			Data:     make([]byte, 128),
		}

		stdMsgs[i] = &StandardTransaction{
			Source:        1,
			TxHash:        bytes,
			NativeMessage: &ethMsg,
		}

	}
	return stdMsgs
}

func TestQuickSort(t *testing.T) {
	stdmsg0 := PrepareData(4)
	stdmsg1 := PrepareData(4)

	worker := func(lft *StandardTransaction, rgt *StandardTransaction) bool {
		return bytes.Compare(lft.TxHash[:], rgt.TxHash[:]) < 0
	}

	StandardTransactions(stdmsg0).QuickSort(worker)
	StandardTransactions(stdmsg1).SortByHash()
	if !reflect.DeepEqual(stdmsg0, stdmsg1) {
		t.Error("mismatch")
	}
}

func TestQuickSortPerformance(t *testing.T) {
	stdmsgs := PrepareData(500000)
	t0 := time.Now()
	worker := func(lft *StandardTransaction, rgt *StandardTransaction) bool {
		return bytes.Compare(lft.TxHash[:], rgt.TxHash[:]) < 0
	}
	StandardTransactions(stdmsgs).QuickSort(worker)

	fmt.Println("append:", time.Now().Sub(t0))
}
