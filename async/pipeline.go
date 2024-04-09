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
	"sync/atomic"
	"time"

	"github.com/arcology-network/common-lib/exp/slice"
)

// Couple is a struct that contains two functions and two channels.
// The first function consumes the input from the channel and produces
// an output that is sent to the second channel, which is consumed by the second function.
type Pipeline[T any] struct {
	total      atomic.Int64
	sleepTime  time.Duration
	bufferSize int

	channels []chan T
	workers  []func(T) (T, bool)
	quit     chan struct{}
}

func NewPipeline[T any](bufferSize int, sleepTime time.Duration, workers ...func(T) (T, bool)) *Pipeline[T] {
	pl := &Pipeline[T]{
		sleepTime:  sleepTime,
		workers:    workers,
		channels:   make([]chan T, len(workers)+1),
		bufferSize: bufferSize,
		quit:       make(chan struct{}),
	}

	for i := 0; i < len(pl.channels); i++ {
		pl.channels[i] = make(chan T, bufferSize)
	}
	return pl
}

// Start starts the goroutines.
func (this *Pipeline[T]) Start() {
	for i := 0; i < len(this.workers); i++ {
		go func(i int) {
			for {
				select {
				case val, ok := <-this.channels[i]:
					if ok {
						if v, ok := this.workers[i](val); ok {
							fmt.Println("Processed 1 msg from", i, "channel, pushing to the next channel")
							this.channels[i+1] <- v
						}
					}

				case <-this.quit:
					this.workers[i] = nil
					return

				default:
					time.Sleep(this.sleepTime * time.Millisecond)
				}
			}
		}(i)
	}
}

// Push pushes values to the channels.
func (this *Pipeline[T]) Push(vals ...T) {
	this.total.Add(int64(len(vals)))
	for _, v := range vals {
		this.channels[0] <- v // Push to the first channel
	}
}

func (this *Pipeline[T]) Await() []T {
	fmt.Println("Awaiting results...")
	arr := this.windUp()
	this.channels[0] = make(chan T, this.bufferSize) // Reset the entrance channel
	return arr
}

// Forcefully close the pipeline and terminate all goroutines.
func (this *Pipeline[T]) Close() {
	this.windUp()
	for i := 0; i < len(this.workers); i++ {
		this.quit <- struct{}{}
	}

	for {
		idx, _ := slice.FindFirstIf(this.workers, func(f func(T) (T, bool)) bool {
			return f != nil
		})

		if idx == -1 {
			break
		} else {
			time.Sleep(this.sleepTime * time.Millisecond)
		}
	}

	for i := 1; i < len(this.channels); i++ {
		close(this.channels[i])
	}
}

func (this *Pipeline[T]) windUp() []T {
	if len(this.channels) == 0 {
		return []T{}
	}
	close(this.channels[0]) // No more values to push

	arr := make([]T, 0, this.total.Load())
	for {
		v := <-this.channels[len(this.channels)-1]
		if arr = append(arr, v); len(arr) == int(this.total.Load()) {
			this.total.Store(0) // Reset the total count
			break               // All results received
		}
	}
	return arr
}
