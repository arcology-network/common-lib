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
)

// "github.com/HPISTechnologies/common-lib/common"

func TestPipeline(t *testing.T) {
	i := 0
	pipe := NewPipeline(
		10, // Buffer size
		5,  // Sleep time 5ms
		func(k ...string) (string, bool) {
			if i < 2 {
				i++
				return k[0] + "-12", false
			}
			return k[0] + "-1", true
		},
		func(k ...string) (string, bool) {
			time.Sleep(1 * time.Second)
			return k[0] + "-2", true
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
		10, // Buffer size
		5,  // Sleep time 5ms
		func(k ...string) (string, bool) {
			return k[0] + "-1", true
		},
		func(k ...string) (string, bool) {
			time.Sleep(2 * time.Second)
			return k[0] + "-2", true
		},
	) // Time out after 5 seconds

	pipe.Start()

	pipe.Push("key1", "key2", "key3")
	pipe.Close()
	output := []string{}
	for i := 0; i < len(output); i++ {
		output = append(output, <-pipe.inChans[1])
	}
}

func TestPipelineRedirect(t *testing.T) {
	pipe := NewPipeline(
		10, // Buffer size
		5,  // Sleep time 5ms
		func(k ...string) (string, bool) {
			return k[0] + "-1", true
		},
		func(k ...string) (string, bool) {
			return k[0] + "-2", true
		},
	) // Time out after 5 seconds
	pipe.Start()

	outChan := make(chan string, 3)
	pipe.RedirectTo(outChan)

	pipe.Push("key1", "key2", "key3")
	pipe.Close()

	output := ToSlice(outChan)

	if !reflect.DeepEqual(output, []string{"key1-1-2", "key2-1-2", "key3-1-2"}) {
		t.Error("Couple Failed", output)
	}
}
