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
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/arcology-network/common-lib/exp/slice"
)

// "github.com/HPISTechnologies/common-lib/common"

func TestSlice(t *testing.T) {
	ccslice := NewSlice[string]()
	strs := ccslice.Append("a").Append("b").Append("c").MoveToSlice()

	if !reflect.DeepEqual(strs.values, []string{"a", "b", "c"}) || len(ccslice.values) != 0 {
		t.Error("Failed", strs.values)
	}
}

func TestPipeline(t *testing.T) {
	i := 0
	pipe := NewPipeline(
		"test1",
		10, // Buffer size
		5,  // Sleep time 5ms
		func(k string, buffer *[]string) ([]string, bool) {
			if i < 2 {
				i++
				*buffer = append(*buffer, k+"-12")
				return []string{}, false
			}

			if len(k) != 0 {
				*buffer = append(*buffer, k+"-1")
				v := slice.Move(buffer)
				return v, true
			} else {
				return nil, false
			}
		},
		func(k string, buffer *[]string) ([]string, bool) {
			time.Sleep(1 * time.Second)

			*buffer = append(*buffer, k+"-2")
			return slice.Move(buffer), true
		},
	) // Time out after 5 seconds

	pipe.Start()
	pipe.Push("key1", "key2", "key3")
	output := pipe.Await()
	fmt.Println(output)

	if !reflect.DeepEqual(output, []string{"key1-12-2", "key2-12-2", "key3-1-2"}) {
		t.Error("Couple Failed", output)
	}
}

func TestPipelineClose(t *testing.T) {
	pipe := NewPipeline(
		"test1",
		10, // Buffer size
		5,  // Sleep time 5ms
		func(k string, buffer *[]string) ([]string, bool) {
			if len(k) == 0 {
				return []string{}, true
			}
			return []string{k + "-1"}, true
		},
		func(k string, buffer *[]string) ([]string, bool) {
			if len(k) == 0 {
				return nil, false
			}
			return []string{k + "-2"}, true
		},
	) // Time out after 5 seconds

	pipe.Start()
	pipe.Push("key1", "key2", "key3")
	pipe.Close()
}

func TestPipelineRedirect(t *testing.T) {
	pipe := NewPipeline(
		"test1",
		10, // Buffer size
		5,  // Sleep time 5ms
		func(k string, buffer *[]string) ([]string, bool) {
			if len(k) == 0 {
				return []string{}, true
			}
			return []string{k + "-1"}, true
		},
		func(k string, buffer *[]string) ([]string, bool) {
			if len(k) == 0 {
				return nil, false
			}
			return []string{k + "-2"}, true
		},
	) // Time out after 5 seconds
	pipe.Start()
	pipe.Push("key1", "key2", "key3")

	outChan := make(chan string, 3)
	pipe.RedirectTo(outChan)

	output := []string{}
	for i := 0; i < 3; i++ {
		v := <-outChan
		fmt.Println(v)
		output = append(output, v)
	}

	if !reflect.DeepEqual(output, []string{"key1-1-2", "key2-1-2", "key3-1-2"}) {
		t.Error("Couple Failed", output)
	}
}
