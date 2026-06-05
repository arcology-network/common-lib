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
	json "encoding/json"
	unsafe "unsafe"

	stgcodec "github.com/arcology-network/common-lib/crdt/codec"
	crdtcommon "github.com/arcology-network/common-lib/crdt/common"
	"github.com/cespare/xxhash"
)

type stateCellJSON struct {
	GenerationID         uint64          `json:"generationId"`
	JobSequenceID        uint64          `json:"jobSequenceId"`
	JobID                uint64          `json:"jobId"`
	ValueType            uint8           `json:"valueType"`
	Tx                   uint64          `json:"tx"`
	Callee               uint64          `json:"callee"`
	Path                 string          `json:"path"`
	KeyHash              uint64          `json:"keyHash"`
	Reads                uint32          `json:"reads"`
	Writes               uint32          `json:"writes"`
	DeltaWrites          uint32          `json:"deltaWrites"`
	SizeInStorage        uint64          `json:"sizeInStorage"`
	GasUsed              uint64          `json:"gasUsed"`
	Message              string          `json:"message"`
	SkipConflictCheck    bool            `json:"skipConflictCheck"`
	Expanded             bool            `json:"expanded"`
	BlockBound           bool            `json:"blockBound"`
	Committed            bool            `json:"committed"`
	Deleted              bool            `json:"deleted"`
	HasCollision         bool            `json:"hasCollision"`
	Value                json.RawMessage `json:"value"`
}

func (this *StateCell) MarshalJSON() ([]byte, error) {
	path := ""
	if this.path != nil {
		path = *this.path
	}

	value := json.RawMessage("null")
	if this.value != nil {
		value = json.RawMessage(stgcodec.JSONCodec{}.Encode("", this.value))
	}

	return json.Marshal(stateCellJSON{
		GenerationID:      this.GenerationID,
		JobSequenceID:     this.JobSequenceID,
		JobID:             this.JobID,
		ValueType:         this.vType,
		Tx:                this.tx,
		Callee:            this.callee,
		Path:              path,
		KeyHash:           this.keyHash,
		Reads:             this.reads,
		Writes:            this.writes,
		DeltaWrites:       this.deltaWrites,
		SizeInStorage:     this.sizeInStorage,
		GasUsed:           this.gasUsed,
		Message:           this.msg,
		SkipConflictCheck: this.ifSkipConflictCheck,
		Expanded:          this.isExpanded,
		BlockBound:        this.isBlockBound,
		Committed:         this.isCommitted,
		Deleted:           this.isDeleted,
		HasCollision:      this.HasCollision,
		Value:             value,
	})
}

func (this *StateCell) UnmarshalJSON(buffer []byte) error {
	envelope := stateCellJSON{}
	if err := json.Unmarshal(buffer, &envelope); err != nil {
		return err
	}

	path := envelope.Path
	this.Property = Property{
		GenerationID:         envelope.GenerationID,
		JobSequenceID:        envelope.JobSequenceID,
		JobID:                envelope.JobID,
		vType:                envelope.ValueType,
		tx:                   envelope.Tx,
		callee:               envelope.Callee,
		path:                 &path,
		pathBytes:            unsafe.Slice(unsafe.StringData(path), len(path)),
		keyHash:              envelope.KeyHash,
		reads:                envelope.Reads,
		writes:               envelope.Writes,
		deltaWrites:          envelope.DeltaWrites,
		sizeInStorage:        envelope.SizeInStorage,
		gasUsed:              envelope.GasUsed,
		msg:                  envelope.Message,
		ifSkipConflictCheck:  envelope.SkipConflictCheck,
		isExpanded:           envelope.Expanded,
		isBlockBound:         envelope.BlockBound,
		isCommitted:          envelope.Committed,
		isDeleted:            envelope.Deleted,
		HasCollision:         envelope.HasCollision,
		reclaimFunc:          nil,
	}

	if this.keyHash == 0 && len(path) > 0 {
		this.keyHash = xxhash.Sum64String(path)
	}

	this.value = stgcodec.JSONCodec{ID: envelope.ValueType}.Decode("", envelope.Value, nil)
	if this.value == nil {
		this.buf = nil
		return nil
	}

	this.buf = this.value.(crdtcommon.CRDT).Encode()
	return nil
}