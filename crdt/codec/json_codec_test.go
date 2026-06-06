/*
 *   Copyright (c) 2026 Arcology Network

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package codec

import (
	json "encoding/json"
	"testing"

	commutative "github.com/arcology-network/common-lib/crdt/commutative"
	noncommutative "github.com/arcology-network/common-lib/crdt/noncommutative"
)

func TestJSONCodecEncodeDecodeString(t *testing.T) {
	codec := JSONCodec{}
	value := noncommutative.NewString("hello")
	decoded := codec.Decode("", codec.Encode("", value), nil)

	stringer, ok := decoded.(*noncommutative.String)
	if !ok || string(*stringer) != "hello" {
		t.Error("Error: JSONCodec should round-trip noncommutative strings")
	}
}

func TestJSONCodecDecodeUint64WithExplicitID(t *testing.T) {
	value := commutative.NewUnboundedUint64().(*commutative.Uint64)
	value.SetValue(uint64(7))

	buffer, err := json.Marshal(jsonEnvelope{Payload: value.Encode()})
	if err != nil {
		t.Fatalf("Error: failed to marshal the JSON codec envelope: %v", err)
	}

	decoded := (JSONCodec{ID: commutative.UINT64}).Decode("", buffer, nil)
	uint64v, ok := decoded.(*commutative.Uint64)
	if !ok || uint64v.Value().(uint64) != 7 {
		t.Error("Error: JSONCodec should decode buffers with an explicit type ID")
	}
}

func TestJSONCodecNilAndUnknownID(t *testing.T) {
	codec := JSONCodec{}
	if encoded := codec.Encode("", nil); string(encoded) != "null" {
		t.Error("Error: encoding nil should return JSON null")
	}
	if decoded := codec.Decode("", nil, nil); decoded != nil {
		t.Error("Error: decoding an empty buffer should return nil")
	}
	if decoded := codec.Decode("", []byte("null"), nil); decoded != nil {
		t.Error("Error: decoding JSON null should return nil")
	}
	if decoded := codec.Decode("", []byte(`{"typeId":255,"payload":"AQ=="}`), nil); decoded != nil {
		t.Error("Error: unknown type IDs should return nil")
	}
	if decoded := codec.Decode("", []byte(`{"payload":"AQ=="}`), nil); decoded != nil {
		t.Error("Error: missing type IDs without an explicit codec ID should return nil")
	}
	if decoded := codec.Decode("", []byte(`{"typeId":1,"payload":"!!!"}`), nil); decoded != nil {
		t.Error("Error: invalid JSON payloads should return nil")
	}
	if decoded := codec.Decode("", []byte(`{"typeId":1`), nil); decoded != nil {
		t.Error("Error: invalid JSON should return nil")
	}
}
