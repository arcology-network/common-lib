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
	"sort"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	evmcore "github.com/ethereum/go-ethereum/core"
)

type StandardMessage struct {
	ID         uint64
	TxHash     [32]byte
	Native     *evmcore.Message
	Source     uint8
	PrepaidGas uint64 // Gas required for the deferred TX if it has one. Set by the scheduler only. Multiprocessor won't touch it.
	IsDeferred bool   // If the message is deferred execution.
}

// AddrAndSignature return the address and the first 4 bytes of the signature of the message.
func (this *StandardMessage) AddrAndSignature() string {
	length := min(4, len(this.Native.Data))
	return this.Native.To.String() + string(this.Native.Data[:length])
}

type StandardMessages []*StandardMessage

func (this StandardMessages) SortByFee() {
	sort.SliceStable(
		this,
		func(i, j int) bool {
			return this[i].Native.GasLimit < this[j].Native.GasLimit
		},
	)
}

func (this StandardMessages) SortByHash() {
	sort.Slice(this, func(i, j int) bool { return string(this[i].TxHash[:]) < string(this[j].TxHash[:]) })
}

func (this StandardMessages) Count(value *StandardMessage) int {
	counter := 0
	for i := range this {
		if bytes.Equal(this[i].TxHash[:], value.TxHash[:]) {
			counter++
		}
	}
	return counter
}

func (this StandardMessages) Encode() ([]byte, error) {
	if this == nil {
		return []byte{}, nil
	}
	data := make([][]byte, len(this))
	worker := func(start, end, idx int, args ...any) {
		this := args[0].([]any)[0].(StandardMessages)
		data := args[0].([]any)[1].([][]byte)

		for i := start; i < end; i++ {
			encodedMsg := []byte{}
			if encoded, err := MsgEncode(this[i].Native); err == nil {
				encodedMsg = encoded
			}

			tmpData := [][]byte{
				codec.Uint64(this[i].ID).Encode(),
				this[i].TxHash[:],
				encodedMsg,
				[]byte{this[i].Source},
			}
			data[i] = codec.Byteset(tmpData).Encode()
		}
	}
	common.ParallelWorker(len(this), Concurrency, worker, this, data)
	return codec.Byteset(data).Encode(), nil
}

func (this *StandardMessages) Decode(data []byte) ([]*StandardMessage, error) {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	msgs := make([]*StandardMessage, len(fields))

	worker := func(start, end, idx int, args ...any) {
		data := args[0].([]any)[0].([][]byte)
		messages := args[0].([]any)[1].([]*StandardMessage)

		for i := start; i < end; i++ {
			standredMessage := new(StandardMessage)

			fields := codec.Byteset{}.Decode(data[i]).(codec.Byteset)

			standredMessage.ID = uint64(codec.Uint64(0).Decode(fields[0]).(codec.Uint64))
			standredMessage.TxHash = [32]byte(fields[1])

			if len(fields[2]) > 0 {
				msg, err := MsgDecode(fields[2])
				if err != nil {
					return
				}
				standredMessage.Native = msg
			}
			standredMessage.Source = uint8(fields[3][0])
			messages[i] = standredMessage
		}
	}
	common.ParallelWorker(len(fields), Concurrency, worker, fields, msgs)

	return msgs, nil
}
