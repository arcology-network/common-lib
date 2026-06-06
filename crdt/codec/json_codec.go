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
	bytes "bytes"
	json "encoding/json"

	crdtcommon "github.com/arcology-network/common-lib/crdt/common"
)

type JSONCodec struct {
	ID uint8
}

type jsonEnvelope struct {
	TypeID  uint8  `json:"typeId,omitempty"`
	Payload []byte `json:"payload,omitempty"`
}

func (JSONCodec) Encode(_ string, value any) []byte {
	if value == nil {
		return []byte("null")
	}

	envelope, err := json.Marshal(jsonEnvelope{
		TypeID:  value.(crdtcommon.CRDT).TypeID(),
		Payload: value.(crdtcommon.CRDT).Encode(),
	})
	if err != nil {
		return nil
	}

	return envelope
}

func (this JSONCodec) Decode(_ string, buffer []byte, _ any) any {
	trimmed := bytes.TrimSpace(buffer)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil
	}

	envelope := jsonEnvelope{}
	if err := json.Unmarshal(trimmed, &envelope); err != nil {
		return nil
	}

	if envelope.TypeID != 0 {
		return Codec{ID: envelope.TypeID}.Decode("", envelope.Payload, nil)
	}

	if this.ID == 0 {
		return nil
	}

	return Codec{ID: this.ID}.Decode("", envelope.Payload, nil)
}
