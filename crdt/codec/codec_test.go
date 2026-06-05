/*
 *   Copyright (c) 2023 Arcology Network

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

package codec

import (
	"testing"

	commutative "github.com/arcology-network/common-lib/crdt/commutative"
	noncommutative "github.com/arcology-network/common-lib/crdt/noncommutative"
)

func TestCodecEncodeDecodeString(t *testing.T) {
	this := Codec{}
	value := noncommutative.NewString("hello")
	encoded := this.Encode("", value)
	decoded := this.Decode("", encoded, nil)

	stringer, ok := decoded.(*noncommutative.String)
	if !ok || string(*stringer) != "hello" {
		t.Error("Error: Codec should round-trip noncommutative strings")
	}
}

func TestCodecDecodeUint64WithExplicitID(t *testing.T) {
	value := commutative.NewUnboundedUint64().(*commutative.Uint64)
	value.SetValue(uint64(7))

	this := Codec{ID: commutative.UINT64}
	decoded := this.Decode("", value.Encode(), nil)

	uint64v, ok := decoded.(*commutative.Uint64)
	if !ok || uint64v.Value().(uint64) != 7 {
		t.Error("Error: Codec should decode buffers with an explicit type ID")
	}
}

func TestCodecNilAndUnknownID(t *testing.T) {
	this := Codec{}
	if encoded := this.Encode("", nil); len(encoded) != 0 {
		t.Error("Error: encoding nil should return an empty buffer")
	}
	if decoded := this.Decode("", nil, nil); decoded != nil {
		t.Error("Error: decoding an empty buffer should return nil")
	}
	if decoded := (Codec{ID: 255}).Decode("", []byte{1}, nil); decoded != nil {
		t.Error("Error: unknown type IDs should return nil")
	}
}