/*
 *   Copyright (c) 2026 Arcology Network

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
	"errors"
	"math/big"

	"github.com/arcology-network/common-lib/codec"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	evmTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

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

func EncodeReceipt(receipt *evmTypes.Receipt) ([]byte, error) {
	return rlp.EncodeToBytes(receipt)
}

func DecodeReceipt(data []byte) (*evmTypes.Receipt, error) {
	receipt := evmTypes.Receipt{}
	if err := rlp.DecodeBytes(data, &receipt); err != nil {
		return nil, err
	}
	return &receipt, nil
}

// RLP encoder doesn't work with error type, so we implement our own codec here
func EncodeExecutionResult(result *core.ExecutionResult) ([]byte, error) {
	if result == nil {
		return []byte{}, nil
	}

	errString := ""
	if result.Err != nil {
		errString = result.Err.Error()
	}

	return codec.Byteset([][]byte{
		codec.Uint64(result.UsedGas).Encode(),
		codec.Uint64(result.RefundedGas).Encode(),
		codec.String(errString).Encode(),
		codec.Bytes(result.ReturnData).Encode(),
		codec.Bytes(result.ContractAddress.Bytes()).Encode(),
	}).Encode(), nil
}

func DecodeExecutionResult(data []byte) (*core.ExecutionResult, error) {
	if len(data) == 0 {
		return nil, nil
	}

	fields := [][]byte((&codec.Byteset{}).Decode(data).(codec.Byteset))
	result := &core.ExecutionResult{}
	result.UsedGas = uint64(codec.Uint64(0).Decode(fields[0]).(codec.Uint64))
	result.RefundedGas = uint64(codec.Uint64(0).Decode(fields[1]).(codec.Uint64))

	if len(fields[2]) > 0 {
		result.Err = errors.New(string(codec.String("").Decode(fields[2]).(codec.String)))
	}
	result.ReturnData = []byte(codec.Bytes(nil).Decode(fields[3]).(codec.Bytes))
	result.ContractAddress = ethcommon.BytesToAddress(fields[4])
	return result, nil
}
