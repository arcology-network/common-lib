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
	"testing"

	"github.com/arcology-network/common-lib/codec"
)

func TestStagedRemovalSetCodec(t *testing.T) {
	removalSet := NewStagedRemovalSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	removalSet.InsertBatch([]string{"13", "15", "17"})
	removalSet.Commit(nil) // The strings are in the committed set already

	removalSet.InsertBatch([]string{"113", "115", "117"})
	removalSet.DeleteByIndex(1) // {"15"} are in the stagedRemovals set
	removalSet.DeleteByIndex(4) // {"115"} is in the stagedRemovals set
	removalSet.DeleteByIndex(5) // {"117"} is in the staged

	buff := removalSet.Encode()                                                                                                                // Encode the staged removal set
	out := NewStagedRemovalSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).Decode(buff).(*StagedRemovalSet[string]) // Decode the staged removal set

	// set2.Equal(removalSet) // Check if the decoded set is equal to the original
	if !out.Equal(removalSet) {
		t.Error("decoded set is not equal to the original")
	}
}

func TestStagedRemovalSetCodecAllDeleted(t *testing.T) {
	removalSet := NewStagedRemovalSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil)
	removalSet.InsertBatch([]string{"13", "15", "17"})
	removalSet.Commit(nil) // The strings are in the committed set already

	removalSet.InsertBatch([]string{"113", "115", "117"})
	removalSet.DeleteByIndex(1)  // {"15"} are in the stagedRemovals set
	removalSet.DeleteByIndex(4)  // {"115"} is in the stagedRemovals set
	removalSet.DeleteByIndex(5)  // {"117"} is in the staged
	removalSet.allDeleted = true // {"13", "15", "17"} are in the stagedRemovals set,

	buff := removalSet.Encode()                                                                                                                // Encode the staged removal set
	out := NewStagedRemovalSet("", 100, codec.Sizer, codec.EncodeTo, new(codec.String).DecodeTo, nil).Decode(buff).(*StagedRemovalSet[string]) // Decode the staged removal set

	if len(out.Committed().Elements()) != 0 {
		t.Error("committed set should be empty, but it is not")
	}
	if !out.Added().Equal(removalSet.Added()) {
		t.Error("added set should be equal, but it is not")
	}
	if !out.Removed().Equal(removalSet.Removed()) {
		t.Error("removed set should be equal, but it is not")
	}
}
