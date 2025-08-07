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
	orderedset "github.com/arcology-network/common-lib/exp/orderedset"
)

// DeltaSet represents a mutable view over a base set, allowing staged
// additions and deletions. Its key difference from a regular DeltaSet
// is that it has a flag to indicate if all elements are deleted. It is
// especially useful working with a huge set where we may want to delete
// the entire set without explicitly touching each element, which would
// be a huge overhead.
type StagedRemovalSet[K comparable] struct {
	deltaset.DeltaSet[K]
	allDeleted bool // If true, all elements are deleted, i.e. the set is empty.
}

func NewStagedRemovalSet[K comparable](nilVal K, preAlloc int,
	size func(K) int,
	encodeToBuffer func(K, []byte) int,
	decoder func([]byte) K,
	hasher func(K) [32]byte,
	keys ...K) *StagedRemovalSet[K] {
	return &StagedRemovalSet[K]{
		DeltaSet: *deltaset.NewDeltaSet(nilVal, preAlloc, size, encodeToBuffer, decoder, hasher, keys...),
	}
}

func (this *StagedRemovalSet[K]) DeleteAll(
	committed *orderedset.OrderedSet[K],
	stagedAdditions *orderedset.OrderedSet[K]) {
	this.SetCommitted(committed)
	this.SetAdded(stagedAdditions.Clone())
	this.Removed().Clear()
	this.allDeleted = true
}

// For the staged removals, the committed set is always shared. This
// is because the StagedRemovalSet is used to mark elements as deleted
// from a softdelta set without actually removing anything from it.
// The committed set in the StagedRemovalSet should always be the
// empty because of its transient nature.

// This is only one exception where we want to mark all elements
// in the softdelta set as deleted, even in this case we still
// share the committed set with the softdelta set to save memory.
func (this *StagedRemovalSet[K]) CloneFull() *StagedRemovalSet[K] {
	set := this.Clone()
	return set
}

// Clone returns a new instance with the shared shared committed set.
func (this *StagedRemovalSet[K]) Clone() *StagedRemovalSet[K] {
	set := this.CloneDelta()
	set.SetCommitted(this.Committed()) // Share the committed set.
	set.allDeleted = this.allDeleted
	return set
}

func (this *StagedRemovalSet[K]) CloneDelta() *StagedRemovalSet[K] {
	return &StagedRemovalSet[K]{
		DeltaSet: *this.DeltaSet.CloneDelta(),
	}
}

func (this *StagedRemovalSet[K]) Length() uint64 {
	return this.DeltaSet.NonNilCount()
}

func (this *StagedRemovalSet[K]) NewFrom(other *StagedRemovalSet[K]) *StagedRemovalSet[K] {
	return &StagedRemovalSet[K]{
		DeltaSet: *this.DeltaSet.NewFrom(&other.DeltaSet),
	}
}

func (this *StagedRemovalSet[K]) Clear() {
	this.ResetDelta()
	this.Committed().Clear()
	this.allDeleted = false // Reset the allDeleted flag
}

func (this *StagedRemovalSet[K]) Equal(other *StagedRemovalSet[K]) bool {
	return this.allDeleted == other.allDeleted && this.DeltaSet.Equal(&other.DeltaSet)
}
