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
	codec "github.com/arcology-network/common-lib/codec"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

type ReceiptHash struct {
	Txhash      *ethCommon.Hash
	Receipthash *ethCommon.Hash
	GasUsed     uint64
}

type ReceiptHashList struct {
	TxHashList      []ethCommon.Hash
	ReceiptHashList []ethCommon.Hash
	GasUsedList     []uint64
}

func (rhl *ReceiptHashList) GobEncode() ([]byte, error) {
	data := [][]byte{
		Hashes(rhl.TxHashList).Encode(),
		Hashes(rhl.ReceiptHashList).Encode(),
		codec.Uint64s(rhl.GasUsedList).Encode(),
	}
	return codec.Byteset(data).Encode(), nil
}
func (rhl *ReceiptHashList) GobDecode(data []byte) error {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	rhl.TxHashList = Hashes(rhl.TxHashList).Decode(fields[0])
	rhl.ReceiptHashList = Hashes(rhl.ReceiptHashList).Decode(fields[1])
	rhl.GasUsedList = codec.Uint64s(rhl.GasUsedList).Decode(fields[2]).(codec.Uint64s)
	return nil
}
