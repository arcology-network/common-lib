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

package async

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

// "github.com/HPISTechnologies/common-lib/common"

func TestCouple(t *testing.T) {
	couple := NewCouple(
		func(k string) string {
			return k + "-1"
		},
		func(k string) string {
			return k + "-2"
		},
		10,   // Buffer size
		5,    // Sleep time 5ms
		5000) // Time out after 5 seconds

	go couple.Start()

	couple.Push("key1", "key2", "key3")
	output, _ := couple.Await()
	fmt.Println(output)

	if !reflect.DeepEqual(output, []string{"key1-1-2", "key2-1-2", "key3-1-2"}) {
		t.Error("Couple Failed")
	}
}

func TestTriple(t *testing.T) {
	triple := NewTriple(
		func(k string) string {
			return k + "-1"
		},
		func(k string) []byte {
			return []byte(k + "-2")
		},
		func(k []byte) int {
			v := binary.LittleEndian.Uint32(k)
			// v, _ := strconv.Atoi(string(k))
			return int(v)
		},
		10,   // Buffer size
		5,    // Sleep time 5ms
		5000) // Time out after 5 seconds

	go triple.Start()

	triple.Push("key1", "key2", "key3")
	output, err := triple.Await()
	fmt.Println(output, err)

	if !reflect.DeepEqual(output, []int{830039403, 846816619, 863593835}) {
		t.Error("Couple Failed")
	}
}
