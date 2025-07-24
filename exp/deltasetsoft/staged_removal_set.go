/*
 *   Copyright (c) 2025 Arcology Network

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

package softdeltaset

import (
	"github.com/arcology-network/common-lib/exp/deltaset"
)

// DeltaSet represents a mutable view over a base set, allowing staged additions and deletions.
type StagedRemovalset[K comparable] struct {
	deltaset.DeltaSet[K]
}

func NewStagedRemovalset[K comparable](nilVal K, preAlloc int, hasher func(K) [32]byte, keys ...K) *StagedRemovalset[K] {
	return &StagedRemovalset[K]{
		DeltaSet: *deltaset.NewDeltaSet(nilVal, preAlloc, hasher, keys...),
	}
}

func (this *StagedRemovalset[K]) CloneFull() *StagedRemovalset[K] {
	set := this.CloneDelta()
	set.SetCommitted(this.Committed().Clone())
	return set
}

// Clone returns a new instance with the shared shared committed set.
func (this *StagedRemovalset[K]) Clone() *StagedRemovalset[K] {
	this.DeltaSet.Clone()
	return this
}

func (this *StagedRemovalset[K]) CloneDelta() *StagedRemovalset[K] {
	return &StagedRemovalset[K]{
		DeltaSet: *this.DeltaSet.CloneDelta(),
	}
}

func (this *StagedRemovalset[K]) Length() uint64 {
	return this.DeltaSet.NonNilCount()
}

func (this *StagedRemovalset[K]) NewFrom(other *StagedRemovalset[K]) *StagedRemovalset[K] {
	return &StagedRemovalset[K]{
		DeltaSet: *this.DeltaSet.NewFrom(&other.DeltaSet),
	}
}

func (this *StagedRemovalset[K]) Clear() {
	this.ResetDelta()
	this.Committed().Clear()
}
