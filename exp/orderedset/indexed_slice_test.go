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

package orderedset

import (
	"reflect"
	"testing"
)

func TestIndexedSlice(t *testing.T) {
	set := NewOrderedSet[string]("", 10, "1", "2", "5")
	set.Insert("11")

	if ok, _ := set.Exists("11"); !ok {
		t.Error("Error: Key is not equal !")
	}

	if *set.At(0) != "1" {
		t.Error("Error: Key is not equal !")
	}

	if !reflect.DeepEqual(set.Elements(), []string{"1", "2", "5", "11"}) {
		t.Error("Error: Key is not equal !")
	}

	set.DeleteByIndex(0)
	if !reflect.DeepEqual(set.Elements(), []string{"2", "5", "11"}) {
		t.Error("Error: Key is not equal !")
	}

	set.DeleteByIndex(2)
	if !reflect.DeepEqual(set.Elements(), []string{"2", "5"}) {
		t.Error("Error: Key is not equal !")
	}

	set.replace(1, "11")
	if !reflect.DeepEqual(set.Elements(), []string{"2", "11"}) {
		t.Error("Error: Key is not equal !")
	}

	set.replace(0, "111")
	if !reflect.DeepEqual(set.Elements(), []string{"111", "11"}) {
		t.Error("Error: Key is not equal !")
	}

	set.Insert("111")
	if !reflect.DeepEqual(set.Elements(), []string{"111", "11"}) {
		t.Error("Error: Key is not equal !")
	}

	set.Insert("222")
	if !reflect.DeepEqual(set.Elements(), []string{"111", "11", "222"}) {
		t.Error("Error: Key is not equal !")
	}

	set.Delete("11")
	if !reflect.DeepEqual(set.Elements(), []string{"111", "222"}) {
		t.Error("Error: Key is not equal !")
	}

	set.Merge(NewOrderedSet[string]("", 10, "1", "2", "5").Elements())
	if !reflect.DeepEqual(set.Elements(), []string{"111", "222", "1", "2", "5"}) {
		t.Error("Error: Key is not equal !", set.Elements())
	}

	set.Merge(NewOrderedSet[string]("", 10, "111", "222", "1", "2", "6").Elements())
	if !reflect.DeepEqual(set.Elements(), []string{"111", "222", "1", "2", "5", "6"}) {
		t.Error("Error: Key is not equal !", set.Elements())
	}

	set.Clear()
	if set.Length() != 0 {
		t.Error("Error: Key is not equal !")
	}

	set.Merge(NewOrderedSet[string]("", 10, "1", "2", "5").Elements())
	if !reflect.DeepEqual(set.Elements(), []string{"1", "2", "5"}) {
		t.Error("Error: Key is not equal !", set.Elements())
	}
}
