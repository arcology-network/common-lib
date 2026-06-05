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

package statecell

import (
	bytes "bytes"
	json "encoding/json"
	"testing"

	stgcodec "github.com/arcology-network/common-lib/crdt/codec"
	crdtcommon "github.com/arcology-network/common-lib/crdt/common"
	commutative "github.com/arcology-network/common-lib/crdt/commutative"
	noncommutative "github.com/arcology-network/common-lib/crdt/noncommutative"
)

func TestStateCellJSONRoundTripDifferentValueTypes(t *testing.T) {
	tests := []struct {
		name  string
		value crdtcommon.CRDT
	}{
		{name: "uint64", value: commutative.NewBoundedUint64(0, 100)},
		{name: "string", value: noncommutative.NewString("hello")},
		{name: "bytes", value: noncommutative.NewBytes([]byte("payload"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewStateCell(11, "blcc://eth1.0/account/alice/storage/ctrn-0/"+tt.name, 3, 4, 5, tt.value, nil)
			in.GenerationID = 7
			in.JobSequenceID = 8
			in.JobID = 9
			in.SetCallee(10)
			in.sizeInStorage = 12
			in.gasUsed = 13
			in.msg = "meta"
			in.ifSkipConflictCheck = true
			in.isExpanded = true
			in.isBlockBound = true
			in.isCommitted = true
			in.isDeleted = true
			in.HasCollision = true

			buffer, err := json.Marshal(in)
			if err != nil {
				t.Fatalf("Error: failed to marshal statecell to JSON: %v", err)
			}

			shadow := stateCellJSON{}
			if err := json.Unmarshal(buffer, &shadow); err != nil {
				t.Fatalf("Error: failed to inspect statecell JSON envelope: %v", err)
			}

			if shadow.ValueType != in.vType || shadow.Tx != in.tx || shadow.Path != *in.path || shadow.KeyHash != in.keyHash {
				t.Error("Error: JSON should preserve statecell metadata")
			}

			expectedValue := stgcodec.JSONCodec{}.Encode("", in.value)
			if !bytes.Equal(shadow.Value, expectedValue) {
				t.Error("Error: JSON should preserve the typed CRDT value payload")
			}

			out := &StateCell{}
			if err := json.Unmarshal(buffer, out); err != nil {
				t.Fatalf("Error: failed to unmarshal statecell JSON: %v", err)
			}

			if !in.value.(crdtcommon.CRDT).Equal(out.value.(crdtcommon.CRDT)) {
				t.Error("Error: decoded JSON should preserve the typed CRDT value")
			}

			if in.GenerationID != out.GenerationID ||
				in.JobSequenceID != out.JobSequenceID ||
				in.JobID != out.JobID ||
				in.vType != out.vType ||
				in.tx != out.tx ||
				in.callee != out.callee ||
				*in.path != *out.path ||
				in.keyHash != out.keyHash ||
				in.reads != out.reads ||
				in.writes != out.writes ||
				in.deltaWrites != out.deltaWrites ||
				in.sizeInStorage != out.sizeInStorage ||
				in.gasUsed != out.gasUsed ||
				in.msg != out.msg ||
				in.ifSkipConflictCheck != out.ifSkipConflictCheck ||
				in.isExpanded != out.isExpanded ||
				in.isBlockBound != out.isBlockBound ||
				in.isCommitted != out.isCommitted ||
				in.isDeleted != out.isDeleted ||
				in.HasCollision != out.HasCollision {
				t.Error("Error: JSON round-trip should preserve the full statecell metadata")
			}

			if !bytes.Equal(in.value.(crdtcommon.CRDT).Encode(), out.buf) {
				t.Error("Error: JSON decode should restore the cached encoded value")
			}
		})
	}
}

func TestStateCellJSONRoundTripPathAndNilValue(t *testing.T) {
	pathValue := commutative.NewPath().(*commutative.Path)
	pathValue.SetSubPaths([]string{"e-01", "e-02"})
	pathValue.SetAdded([]string{"+01"})
	pathValue.InsertRemoved([]string{"-01"})

	withPath := NewStateCell(2, "blcc://eth1.0/account/alice/storage/path", 1, 2, 3, pathValue, nil)
	buffer, err := json.Marshal(withPath)
	if err != nil {
		t.Fatalf("Error: failed to marshal path statecell to JSON: %v", err)
	}

	out := &StateCell{}
	if err := json.Unmarshal(buffer, out); err != nil {
		t.Fatalf("Error: failed to unmarshal path statecell JSON: %v", err)
	}

	if !withPath.value.(crdtcommon.CRDT).Equal(out.value.(crdtcommon.CRDT)) {
		t.Error("Error: JSON should round-trip path CRDT values")
	}

	withoutValue := NewStateCell(3, "blcc://eth1.0/account/alice/storage/deleted", 4, 5, 6, nil, nil)
	withoutValue.isCommitted = true
	withoutValue.isDeleted = true

	nilBuffer, err := json.Marshal(withoutValue)
	if err != nil {
		t.Fatalf("Error: failed to marshal nil-value statecell to JSON: %v", err)
	}

	shadow := stateCellJSON{}
	if err := json.Unmarshal(nilBuffer, &shadow); err != nil {
		t.Fatalf("Error: failed to inspect nil-value JSON envelope: %v", err)
	}
	if string(shadow.Value) != "null" {
		t.Error("Error: nil statecell values should encode as JSON null")
	}

	nilOut := &StateCell{}
	if err := json.Unmarshal(nilBuffer, nilOut); err != nil {
		t.Fatalf("Error: failed to unmarshal nil-value statecell JSON: %v", err)
	}
	if nilOut.value != nil || len(nilOut.buf) != 0 {
		t.Error("Error: nil statecell values should decode back to nil")
	}
	if !nilOut.isDeleted || !nilOut.isCommitted {
		t.Error("Error: nil-value JSON should preserve statecell metadata")
	}
}