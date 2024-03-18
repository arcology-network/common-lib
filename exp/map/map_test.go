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
package mapi

import (
	"sort"
	"testing"

	"github.com/arcology-network/common-lib/common"
)

func TestMapKeys(t *testing.T) {
	_map := map[uint32]int{}
	_map[11] = 99
	_map[21] = 25

	keys := common.MapKeys(_map)
	if len(keys) != 2 || (keys[0] != 11 && keys[0] != 21) {
		t.Error("Error: Not equal")
	}
}

func TestMapValues(t *testing.T) {
	_map := map[uint32]int{}
	_map[11] = 99
	_map[21] = 25

	keys := common.MapValues(_map)
	sort.Ints(keys)
	if keys[0] != 25 || keys[1] != 99 {
		t.Error("Error: Not equal")
	}
}

func TestMapMoveIf(t *testing.T) {
	m := map[string]bool{
		"1": true,
		"2": false,
		"3": true,
		"4": false,
	}

	common.MapRemoveIf(m, func(k string, _ bool) bool { return k == "1" })
	if len(m) != 3 {
		t.Error("Error: Failed to remove nil values !")
	}

	target := map[string]bool{}
	common.MapMoveIf(m, target, func(k string, _ bool) bool { return k == "2" })
	if len(m) != 2 || len(target) != 1 {
		t.Error("Error: Failed to remove nil values !")
	}
}

func TestMapGenerics(t *testing.T) {
	m := map[string]bool{
		"1": true,
		"2": false,
		"3": true,
		"4": false,
	}

	IfNotFoundDo(m, []string{"5"}, func(k string) string { return k }, func(k string) bool { return true })
	if len(m) != 5 {
		t.Error("Error: Failed to set nil values !")
	}

	IfFoundDo(m, []string{"1", "5"}, func(k string, _ *bool) bool { return false })
	if m["1"] || m["5"] {
		t.Error("Error: Failed to set nil values !")
	}

	ParallelIfNotFoundDo(m, []string{"6"}, 2, func(k string) bool { return true })
	if len(m) != 6 {
		t.Error("Error: Failed to set nil values !")
	}

	ParalleIfFoundDo(m, []string{"6"}, 2, func(k string) bool { return false })
	if m["6"] {
		t.Error("Error: Failed to set nil values !")
	}

	m1 := map[string]bool{
		"1": true,
		"2": false,
		"3": true,
		"4": false,
	}

	m2 := map[string]bool{
		"3": true,
		"4": false,
		"1": true,
		"2": false,
	}

	if !EqualIf(m1, m2, func(v0 bool, v1 bool) bool { return v0 == v1 }) {
		t.Error("Error: Failed to compare maps !")
	}

	m1 = map[string]bool{
		"1": true,
		"2": false,
	}

	m2 = map[string]bool{
		"3": true,
		"4": false,
		"1": true,
		"2": false,
	}

	target := map[string]bool{
		"3": true,
		"4": false,
	}

	Sub(m2, m1)
	if !EqualIf(m2, target, func(v0 bool, v1 bool) bool { return v0 == v1 }) {
		t.Error("Error: Failed to compare maps !")
	}

	m3 := map[string]int{
		"3": 8,
		"4": 12,
		"1": 89,
		"2": 90,
	}

	if k, v := FindKey(m3, func(k0, k1 string) bool { return k0 > k1 }); k != "4" || v != 12 {
		t.Error("Error: Failed to get the max !", k, v)
	}

	if k, v := FindKey(m3, func(k0, k1 string) bool { return k0 < k1 }); k != "1" || v != 89 {
		t.Error("Error: Failed to get the max !", k, v)
	}

	if k, v := FindValue(m3, func(k0, k1 int) bool { return k0 < k1 }); k != "3" || v != 8 {
		t.Error("Error: Failed to get the max !", k, v)
	}

	if k, v := FindValue(m3, func(k0, k1 int) bool { return k0 > k1 }); k != "2" || v != 90 {
		t.Error("Error: Failed to get the max !", k, v)
	}

}

func TestMapMaxMinGenerics(t *testing.T) {
	m3 := map[string]int{
		"3": 8,
		"4": 12,
		"1": 89,
		"2": 90,
	}

	if k, v := FindKey(m3, func(k0, k1 string) bool { return k0 < k1 }); k != "1" || v != 89 {
		t.Error("Error: Failed to get the max !", k, v)
	}

	if k, v := FindKey(m3, func(k0, k1 string) bool { return k0 > k1 }); k != "4" || v != 12 {
		t.Error("Error: Failed to get the max !", k, v)
	}

	if k, v := FindValue(m3, func(v0, v1 int) bool { return v0 < v1 }); k != "3" || v != 8 {
		t.Error("Error: Failed to get the max !", k, v)
	}

	if k, v := FindValue(m3, func(v0, v1 int) bool { return v0 > v1 }); k != "2" || v != 90 {
		t.Error("Error: Failed to get the max !", k, v)
	}
}
