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

package noncommutative

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/arcology-network/common-lib/exp/slice"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestNewBigint(t *testing.T) {
	v := NewBigint(100).(*Bigint)

	out, _, _ := v.Get()
	outV := (out.(big.Int))
	v2 := (out.(big.Int))
	if outV.Cmp(&v2) != 0 {
		t.Error("Mismatch")
	}
}

func TestBigintCodecs(t *testing.T) {
	v := NewBigint(100).(*Bigint)

	out, _, _ := v.Get()
	outV := out.(big.Int)
	v2 := (out.(big.Int))
	if outV.Cmp(&v2) != 0 {
		t.Error("Mismatch")
	}

	encoded, err := rlp.EncodeToBytes(out)
	if err != nil {
		t.Error(err)
	}

	var decoded big.Int
	err = rlp.DecodeBytes(encoded, &decoded)
	if err != nil {
		t.Error(err)
	}

	if decoded.Uint64() != 100 {
		t.Error("Mismatch expecting ", 100)
	}
}

func TestBigintRlpCodecs(t *testing.T) {
	in := NewInt64(111)
	buffer := in.StorageEncode("")
	out := new(Int64).StorageDecode("", buffer)

	if *out.(*Int64) != 111 {
		t.Error("Mismatch expecting ", 100)
	}
}

func TestInt64RlpCodec(t *testing.T) {
	v := NewInt64(12345)
	buffer := v.StorageEncode("")
	output := new(Int64).StorageDecode("", buffer)

	if *v != *output.(*Int64) {
		fmt.Println("Error: Missmatched")
	}
}

func TestUint64Codec(t *testing.T) {
	v := NewUint64(56789)
	buffer := v.StorageEncode("")
	output := new(Uint64).StorageDecode("", buffer)

	if *v != *output.(*Uint64) {
		fmt.Println("Error: Missmatched")
	}
}


func TestByteRlp(t *testing.T) {
	v2 := slice.New[byte](32, 11)
	encoded, _ := rlp.EncodeToBytes(v2)

	buf := []byte{}
	err := rlp.DecodeBytes(encoded, &buf)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(v2, buf) {
		fmt.Println("Error: Missmatched")
	}
}
