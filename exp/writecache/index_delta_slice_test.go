/*
 *   Copyright (c) 2024 Arcology Network

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

package deltaslice

import (
	"sort"
	"testing"

	mapi "github.com/arcology-network/common-lib/exp/map"
	"github.com/arcology-network/common-lib/exp/slice"
)

func TestIndexedSliceString(t *testing.T) {
	// deltaSlice := NewDeltaSlice[string](10)
	// strings := []string{}
	indexed := NewIndexDeltaSlice[string, string]()

	indexed.SetByKey("aa", "AA++")
	if str := indexed.GetByKey("aa"); str == nil || *str != "AA++" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := indexed.GetByKey("bb"); str != nil {
		t.Errorf("Expected nil, got %d", str)
	}

	indexed.SetByKey("bb", "BB-")
	if str := indexed.GetByKey("bb"); str == nil || *str != "BB-" {
		t.Errorf("Expected gg, got %d", str)
	}

	indexed.SetByKey("cc", "C")
	if str := indexed.GetByKey("cc"); str == nil || *str != "C" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := indexed.GetByKey("aa"); str == nil || *str != "AA++" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := indexed.GetByIndex(0); str == nil || *str != "AA++" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := indexed.GetByIndex(1); str == nil || *str != "BB-" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := indexed.GetByIndex(11); str != nil {
		t.Errorf("Expected gg, got %d", str)
	}

	indexed.SetByIndex(1, "1234")
	if str := indexed.GetByIndex(1); str == nil || *str != "1234" {
		t.Errorf("Expected gg, got %d", str)
	}

	if str := indexed.GetByIndex(11); str != nil {
		t.Errorf("Expected gg, got %d", str)
	}

	if indexed.SetByIndex(11, "fefe") {
		t.Errorf("Should fail")
	}

	if *indexed.IndexToKey(0) != "aa" {
		t.Errorf("Expected aa, got %s", *indexed.IndexToKey(0))
	}

	if indexed.KeyToIndex("aa") != 0 {
		t.Errorf("Expected aa, got %s", *indexed.IndexToKey(0))
	}

	if indexed.KeyToIndex("aa1") != -1 {
		t.Errorf("Expected aa, got %s", *indexed.IndexToKey(0))
	}

	if !indexed.DeleteByKey("aa") {
		t.Errorf("Expected aa, got %s", *indexed.IndexToKey(0))
	}

	if !indexed.DeleteByIndex(0) {
		t.Errorf("Expected aa, got %s", *indexed.IndexToKey(0))
	}

	indexed.Commit()
	if len(indexed.index) == len(indexed.DeltaSlice.Values()) {
		t.Errorf("Expected %d, got %d", len(indexed.index), len(indexed.DeltaSlice.Values()))
	}

	keys := mapi.Keys(indexed.index)
	sort.StringSlice(keys).Sort()
	if !slice.Equal(keys, []string{"aa", "bb", "cc"}) {
		t.Errorf("Expected %v, got %v", keys, []string{"aa", "bb", "cc"})
	}
}
