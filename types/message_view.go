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
	"math/big"

	"github.com/arcology-network/common-lib/codec"
)

// TransactionView provides a concise representation of a standard Tx,
// capturing essential details such as sender, recipient, function selector, and gas price.
// It is mainly used for conflict detection and resolution in concurrent transaction processing.
type TransactionView struct {
	Hash     [32]byte // Transaction Hash
	ID       uint64   // Transaction ID
	From     [20]byte
	To       [20]byte
	Selector [4]byte
	GasPrice *big.Int
}

func NewTransactionView(msg *StandardMessage) *TransactionView {
	to, selector := msg.GetAddressAndSelector()
	return &TransactionView{
		Hash:     msg.TxHash,
		ID:       msg.ID,
		From:     msg.Native.From,
		To:       to,
		Selector: selector,
		GasPrice: msg.Native.GasPrice,
	}
}

func (this *TransactionView) Size() int {
	return 32 + 8 + 20 + 20 + 4 + this.GasPrice.BitLen()/8 + 1
}

func (this *TransactionView) Encode() ([]byte, error) {
	buffer := make([]byte, 32+8+20+20+4+this.GasPrice.BitLen()/8+1)

	offset := codec.Bytes32(this.Hash).EncodeTo(buffer)
	offset += codec.Uint64(this.ID).EncodeTo(buffer[offset:])
	offset += codec.Bytes20(this.From).EncodeTo(buffer[offset:])
	offset += codec.Bytes20(this.To).EncodeTo(buffer[offset:])
	offset += codec.Bytes4(this.Selector).EncodeTo(buffer[offset:])
	// gasPriceBytes := this.GasPrice.Bytes()
	offset += codec.Bytes(this.GasPrice.Bytes()).EncodeTo(buffer[offset:])
	return buffer[:offset], nil
}

func (this *TransactionView) Decode(data []byte) any {
	copy(this.Hash[:], data[0:32])
	this.ID = uint64(codec.Uint64(0).Decode(data[32:40]).(codec.Uint64))
	copy(this.From[:], data[40:60])
	copy(this.To[:], data[60:80])
	copy(this.Selector[:], data[80:84])
	this.GasPrice = new(big.Int).SetBytes(data[84:])
	return this
}
