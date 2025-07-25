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

package stringdeltaset

import (
	"testing"
)

func TestStagedRemovalSetCodec(t *testing.T) {
	removalSet := NewStagedRemovalSet("", 100, sizer, encodeToBuffer, decoder, nil)
	removalSet.InsertBatch([]string{"13", "15", "17"})
	removalSet.Commit(nil) // The strings are in the committed set already
	removalSet.allDeleted = true

	removalSet.InsertBatch([]string{"113", "115", "117"})
	removalSet.DeleteByIndex(1) // {"15"} are in the stagedRemovals set
	removalSet.DeleteByIndex(4) // {"115"} is in the stagedRemovals set
	removalSet.DeleteByIndex(5) // {"117"} is in the staged

	buff := removalSet.Encode()                                                                                       // Encode the staged removal set
	out := NewStagedRemovalSet("", 100, sizer, encodeToBuffer, decoder, nil).Decode(buff).(*StagedRemovalSet[string]) // Decode the staged removal set

	// set2.Equal(removalSet) // Check if the decoded set is equal to the original
	if !out.Equal(removalSet) {
		t.Error("decoded set is not equal to the original")
	}
}
