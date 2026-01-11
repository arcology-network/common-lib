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

// MessageView provides a concise representation of a standard message,
// capturing essential details such as sender, recipient, function selector, and gas price.
type MessageView struct {
	From     [20]byte
	To       [20]byte
	Selector [4]byte
	GasPrice *big.Int
}

func NewMessageView(msg *StandardMessage) *MessageView {
	to, selector := msg.GetAddressAndSelector()
	return &MessageView{
		From:     msg.Native.From,
		To:       to,
		Selector: selector,
		GasPrice: msg.Native.GasPrice,
	}
}

func (this *MessageView) Size() int {
	return 20 + 20 + 4 + this.GasPrice.BitLen()/8 + 1
}

func (this *MessageView) Encode() ([]byte, error) {
	buffer := make([]byte, 20+20+4+this.GasPrice.BitLen()/8+1)

	offset := codec.Bytes20(this.From).EncodeTo(buffer)
	offset += codec.Bytes20(this.To).EncodeTo(buffer[offset:])
	offset += codec.Bytes4(this.Selector).EncodeTo(buffer[offset:])
	gasPriceBytes := this.GasPrice.Bytes()
	offset += codec.Bytes(gasPriceBytes).EncodeTo(buffer[offset:])
	return buffer[:offset], nil
}

func (this *MessageView) Decode(data []byte) any {
	copy(this.From[:], data[0:20])
	copy(this.To[:], data[20:40])
	copy(this.Selector[:], data[40:44])
	this.GasPrice = new(big.Int).SetBytes(data[44:])
	return this
}
