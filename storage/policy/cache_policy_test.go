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

package cachepolicy

import (
	"fmt"
	"testing"
	"time"

	slice "github.com/arcology-network/common-lib/exp/slice"
)

func TestCachePolicy(t *testing.T) {
	t0 := time.Now()
	fmt.Println("CachePolicy FreeMemory:", time.Since(t0))
	values := []interface{}{nil, nil, 1, 2}
	// common.RemoveIf(&values, func(v interface{}) bool { return v == nil })
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2] actual: ", values)
	}

	values = []interface{}{1, nil, nil, 2}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{1, 2, nil, nil}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{1, nil, 2, nil}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{1, 2}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{nil, nil}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 0 {
		t.Error("Error: Expected [], actual: ", values)
	}
}
