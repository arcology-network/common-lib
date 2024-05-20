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
	"sync"
	"testing"
	"time"
)

func TestSlice(t *testing.T) {
	ccslice := NewSlice[string]()
	strs := ccslice.Append("a").Append("b").Append("c").MoveToSlice()

	if !reflect.DeepEqual(strs, []string{"a", "b", "c"}) || len(ccslice.values) != 0 {
		t.Error("Failed", strs)
	}

	t0 := time.Now()
	var lock sync.Mutex
	for i := 0; i < 3000; i++ {
		lock.Lock()
		lock.Unlock()
	}
	fmt.Println("Lock/Unlock: ", time.Now().Sub(t0))
}

func TestPipelineAwaitConcurrency(t *testing.T) {
	// dict := map[string]bool{}
	pipe := NewPipeline(
		"test1",
		1, // Buffer size
		5, // Sleep time 5ms
		func(k string, buffer *Slice[string]) ([]string, bool, bool) {
			if len(k) == 0 {
				v := buffer.Append(k).MoveToSlice()
				return v, true, true
			}

			buffer.Append(k + "-1")
			return buffer.ToSlice(), false, false
		},
		func(k string, buffer *Slice[string]) ([]string, bool, bool) {
			time.Sleep(1 * time.Second)

			if len(k) == 0 {
				return buffer.MoveToSlice(), true, true
			}
			buffer.Append(k + "-2")
			return nil, false, false
		},
	) // Time out after 5 seconds

	pipe.Start()
	pipe.Push("key1", "key2", "key3", "")
	output := pipe.Await()
	fmt.Println(output)

	if !reflect.DeepEqual(output, []string{"key1-1-2", "key2-1-2", "key3-1-2"}) {
		t.Error("Couple Failed", output)
	}

	pipe.Push("key21", "key22", "key23", "")
	output = pipe.Await()
	fmt.Println(output)

	if !reflect.DeepEqual(output, []string{"key21-1-2", "key22-1-2", "key23-1-2"}) {
		t.Error("Couple Failed", output)
	}
}

// func TestPipelineAwait(t *testing.T) {
// 	dict := map[string]bool{}
// 	i := 0
// 	pipe := NewPipeline(
// 		"test1",
// 		10, // Buffer size
// 		5,  // Sleep time 5ms
// 		func(k string, buffer *Slice[string]) ([]string, bool, bool) {
// 			dict["key1"] = true
// 			fmt.Println(k)
// 			if i < 2 {
// 				i++
// 				buffer.Append(k + "-12")
// 				return []string{}, false, false
// 			}

// 			if len(k) != 0 {
// 				// *buffer = append(*buffer, k+"-1")
// 				v := buffer.Append(k + "-1").MoveToSlice()
// 				// v := slice.Move(buffer)
// 				return v, true, false
// 			} else {
// 				return nil, false, false
// 			}
// 		},
// 		func(k string, buffer *Slice[string]) ([]string, bool, bool) {
// 			time.Sleep(1 * time.Second)
// 			// dict["key1"] = true

// 			return buffer.Append(k + "-2").MoveToSlice(), true, k == ""
// 		},
// 	) // Time out after 5 seconds

// 	pipe.Start()
// 	// var output []string
// 	common.ParallelExecute(
// 		func() {
// 			pipe.Push("key1", "key2", "key3")
// 		},
// 		func() {
// 			pipe.Push("key11", "key12", "key13")
// 		},
// 		// func() {
// 		// 	pipe.Push("key4", "key5", "key6")
// 		// },
// 	)
// 	pipe.Push("")
// 	pipe.Await()
// 	// common.ParallelExecute(
// 	// 	func() {
// 	// 		pipe.Push("key122", "key222", "key322")
// 	// 	},
// 	// )
// 	// output := pipe.Await()

// 	// if pipe.CountVacant() != 0 {
// 	// 	t.Error("CountVacant Failed", pipe.CountVacant())
// 	// }

// 	// pipe.Push("")
// 	// fmt.Println(output)
// }

// func(k string, buffer *Slice[string]) ([]string, bool, bool) {
// 	if len(k) == 0 {
// 		v := buffer.Append(k).MoveToSlice()
// 		return v, true, true
// 	}

// 	buffer.Append(k + "-1")
// 	return buffer.ToSlice(), false, false
// },

// func(k string, buffer *Slice[string]) ([]string, bool, bool) {
// 	time.Sleep(1 * time.Second)

// 	if len(k) == 0 {
// 		return buffer.MoveToSlice(), true, true
// 	}
// 	buffer.Append(k + "-2")
// 	return nil, false, false
// },

func TestPipelineClose(t *testing.T) {
	dict := map[string]bool{}

	pipe := NewPipeline(
		"test1",
		10, // Buffer size
		5,  // Sleep time 5ms
		func(k string, buffer *Slice[string]) ([]string, bool, bool) {
			dict["key1"] = true
			fmt.Println("Func 1:", k)
			if len(k) == 0 {
				return buffer.Append("").MoveToSlice(), true, true
			}

			buffer.Append(k + "-1")
			return nil, false, false
		},
		func(k string, buffer *Slice[string]) ([]string, bool, bool) {
			fmt.Println("Func 2:", k)
			dict["key1"] = true
			if len(k) == 0 {
				return buffer.Append("").MoveToSlice(), true, true
			}
			buffer.Append(k + "-2")
			return nil, false, false
		},
	) // Time out after 5 seconds

	pipe.Start()
	pipe.Push("key1", "key2", "key3", "")
	output := pipe.Await()
	time.Sleep(2 * time.Second)

	if !reflect.DeepEqual(output, []string{"key1-1-2", "key2-1-2", "key3-1-2", ""}) {
		t.Error("Couple Failed", output)
	}

	pipe.Close()

	pipe.Push("key11", "key12", "key13", "")
	output = pipe.Await()

	if !reflect.DeepEqual(output, []string{"key11-1-2", "key12-1-2", "key13-1-2", ""}) {
		t.Error("Couple Failed", output)
	}
	time.Sleep(2 * time.Second)
}

func TestPipelineRedirect(t *testing.T) {
	pipe := NewPipeline(
		"test1",
		10, // Buffer size
		5,  // Sleep time 5ms
		func(k string, buffer *Slice[string]) ([]string, bool, bool) {
			if len(k) == 0 {
				return buffer.Append("").MoveToSlice(), true, true
			}
			buffer.Append(k + "-1")
			return nil, false, false
		},
		func(k string, buffer *Slice[string]) ([]string, bool, bool) {
			if len(k) == 0 {
				return buffer.MoveToSlice(), true, true
			}

			buffer.Append(k + "-2")
			return nil, false, false
		},
	) // Time out after 5 seconds
	pipe.Start()
	pipe.Push("key1", "key2", "key3", "")

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
