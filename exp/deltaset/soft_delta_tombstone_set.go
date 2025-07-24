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

package deltaset

// DeltaSet represents a mutable view over a base set, allowing staged additions and deletions.
type TombstoneSet[K comparable] struct {
	DeltaSet[K]
}

func NewTombstoneSet[K comparable](nilVal K, preAlloc int, hasher func(K) [32]byte, keys ...K) *TombstoneSet[K] {
	return &TombstoneSet[K]{
		DeltaSet: *NewDeltaSet(nilVal, preAlloc, hasher, keys...),
	}
}

func (this *TombstoneSet[K]) CloneFull() *TombstoneSet[K] {
	set := this.CloneDelta()
	set.committed = this.committed.Clone()
	return set
}

// Clone returns a new instance with the shared shared committed set.
func (this *TombstoneSet[K]) Clone() *TombstoneSet[K] {
	this.DeltaSet.Clone()
	return this
}

func (this *TombstoneSet[K]) CloneDelta() *TombstoneSet[K] {
	return &TombstoneSet[K]{
		DeltaSet: *this.DeltaSet.CloneDelta(),
	}
}

func (this *TombstoneSet[K]) Length() uint64 {
	return this.DeltaSet.NonNilCount()
}

func (this *TombstoneSet[K]) NewFrom(other *TombstoneSet[K]) *TombstoneSet[K] {
	return &TombstoneSet[K]{
		DeltaSet: *this.DeltaSet.NewFrom(&other.DeltaSet),
	}
}
