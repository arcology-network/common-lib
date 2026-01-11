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
	"bytes"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
)

func makeStandardMessage(from ethcommon.Address, to *ethcommon.Address, selector []byte, gasPrice *big.Int) *StandardMessage {
	data := make([]byte, len(selector))
	copy(data, selector)

	return &StandardMessage{
		Native: &ethcore.Message{
			From:     from,
			To:       to,
			GasPrice: gasPrice,
			Data:     data,
		},
	}
}

func TestMessageSummaryEncodeDecodeRoundTrip(t *testing.T) {
	fromAddr := ethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	toAddr := ethcommon.HexToAddress("0x2222222222222222222222222222222222222222")
	gasPrice := big.NewInt(0).SetUint64(123456789)
	selector := []byte{0xde, 0xad, 0xbe, 0xef, 0x01}

	msg := makeStandardMessage(fromAddr, &toAddr, selector, gasPrice)

	summary := NewMessageView(msg)
	encoded, err := summary.Encode()
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	if len(encoded) != 44+len(gasPrice.Bytes()) {
		t.Fatalf("unexpected encoded length: %d", len(encoded))
	}

	decoded := (&MessageView{}).Decode(encoded).(*MessageView)

	var expectedFrom [20]byte
	copy(expectedFrom[:], fromAddr.Bytes())
	if decoded.From != expectedFrom {
		t.Fatalf("unexpected From: %x", decoded.From)
	}

	var expectedTo [20]byte
	copy(expectedTo[:], toAddr.Bytes())
	if decoded.To != expectedTo {
		t.Fatalf("unexpected To: %x", decoded.To)
	}

	expectedSelector := [4]byte{0xde, 0xad, 0xbe, 0xef}
	if decoded.Selector != expectedSelector {
		t.Fatalf("unexpected selector: %x", decoded.Selector)
	}

	if decoded.GasPrice.Cmp(gasPrice) != 0 {
		t.Fatalf("unexpected gas price: %s", decoded.GasPrice.String())
	}
}

func TestMessageSummaryEncodeZeroGasPrice(t *testing.T) {
	fromAddr := ethcommon.HexToAddress("0x3333333333333333333333333333333333333333")
	selector := []byte{}
	gasPrice := big.NewInt(0)

	msg := makeStandardMessage(fromAddr, nil, selector, gasPrice)

	summary := NewMessageView(msg)
	encoded, err := summary.Encode()
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	if len(encoded) != 44 {
		t.Fatalf("unexpected encoded length for zero gas price: %d", len(encoded))
	}

	decoded := (&MessageView{}).Decode(encoded).(*MessageView)

	var expectedFrom [20]byte
	copy(expectedFrom[:], fromAddr.Bytes())
	if decoded.From != expectedFrom {
		t.Fatalf("unexpected From: %x", decoded.From)
	}

	zeroAddr := [20]byte{}
	if decoded.To != zeroAddr {
		t.Fatalf("expected zero To address, got: %x", decoded.To)
	}

	expectedSelector := [4]byte{}
	if decoded.Selector != expectedSelector {
		t.Fatalf("unexpected selector: %x", decoded.Selector)
	}

	if decoded.GasPrice.Sign() != 0 {
		t.Fatalf("expected zero gas price, got: %s", decoded.GasPrice.String())
	}
}

func TestMessageSummaryEncodeProducesDeterministicLayout(t *testing.T) {
	from := [20]byte{1, 2, 3}
	to := [20]byte{4, 5, 6}
	selector := [4]byte{0xde, 0xad, 0xbe, 0xef}
	gasPrice := big.NewInt(0).SetBytes([]byte{0x01, 0x02, 0x03})

	summary := &MessageView{
		From:     from,
		To:       to,
		Selector: selector,
		GasPrice: gasPrice,
	}

	encoded, err := summary.Encode()
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	expected := make([]byte, 0, len(encoded))
	expected = append(expected, from[:]...)
	expected = append(expected, to[:]...)
	expected = append(expected, selector[:]...)
	expected = append(expected, gasPrice.Bytes()...)

	if !bytes.Equal(encoded, expected) {
		t.Fatalf("encoded bytes mismatch\nexpected: %x\nactual:   %x", expected, encoded)
	}
}

func TestMessageViewDecodeParsesFields(t *testing.T) {
	from := [20]byte{9, 9, 9}
	to := [20]byte{8, 8, 8}
	selector := [4]byte{0x11, 0x22, 0x33, 0x44}
	gasBytes := []byte{0xaa}

	encoded := make([]byte, 0, 44+len(gasBytes))
	encoded = append(encoded, from[:]...)
	encoded = append(encoded, to[:]...)
	encoded = append(encoded, selector[:]...)
	encoded = append(encoded, gasBytes...)

	decoded := (&MessageView{}).Decode(encoded).(*MessageView)

	if decoded.From != from {
		t.Fatalf("unexpected decoded From: %x", decoded.From)
	}

	if decoded.To != to {
		t.Fatalf("unexpected decoded To: %x", decoded.To)
	}

	if decoded.Selector != selector {
		t.Fatalf("unexpected decoded selector: %x", decoded.Selector)
	}

	if decoded.GasPrice.Cmp(big.NewInt(0).SetBytes(gasBytes)) != 0 {
		t.Fatalf("unexpected decoded gas price: %s", decoded.GasPrice.String())
	}
}
