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
	"sync/atomic"
	"time"

	"github.com/arcology-network/common-lib/exp/slice"
)

// Couple is a struct that contains two functions and two inChans.
// The first function consumes the input from the channel and produces
// an output that is sent to the second channel, which is consumed by the second function.
type Pipeline[T any] struct {
	total       atomic.Int64
	sleepTime   time.Duration
	channelSize int

	inChans []chan T
	outBufs [][]T // Buffers to store the values temporarily before pushing to the next channel.
	workers []func(...T) (T, bool)
	quit    chan struct{}
}

func NewPipeline[T any](channelSize int, sleepTime time.Duration, workers ...func(...T) (T, bool)) *Pipeline[T] {
	pl := &Pipeline[T]{
		sleepTime:   sleepTime,
		workers:     workers,
		inChans:     make([]chan T, len(workers)+1),
		outBufs:     make([][]T, len(workers)), // output outBuf
		channelSize: channelSize,
		quit:        make(chan struct{}),
	}

	for i := 0; i < len(pl.outBufs); i++ {
		pl.outBufs[i] = make([]T, 0, channelSize)
	}

	for i := 0; i < len(pl.inChans); i++ {
		pl.inChans[i] = make(chan T, channelSize)
	}
	return pl
}

// Start starts the goroutines.
func (this *Pipeline[T]) Start() *Pipeline[T] {
	for i := 0; i < len(this.workers); i++ {
		go func(i int) {
			for {
				select {
				case val, ok := <-this.inChans[i]:
					if ok {
						if v, ok := this.workers[i](val); ok {
							this.outBufs[i] = append(this.outBufs[i], v)
							for j := 0; j < len(this.outBufs[i]); j++ {
								this.inChans[i+1] <- this.outBufs[i][j] // Send to the downstream channel
							}
							this.outBufs[i] = this.outBufs[i][:0] // Reset the outBuf
						} else {
							this.outBufs[i] = append(this.outBufs[i], v)
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
	return this
}

// Push pushes values to the inChans.
func (this *Pipeline[T]) Push(vals ...T) {
	this.total.Add(int64(len(vals)))
	for _, v := range vals {
		this.inChans[0] <- v // Push to the first channel
	}
}

// Awiat waits for all the results to be processed and returns. No new values can be pushed into the pipeline
// after before Await() returns.
func (this *Pipeline[T]) Await() []T {
	arr := this.windUp()
	this.inChans[0] = make(chan T, this.channelSize) // Reopen the entrance channel
	return arr
}

// Close the pipeline and terminate all goroutines.
func (this *Pipeline[T]) Close() {
	this.windUp()
	for i := 0; i < len(this.workers); i++ {
		this.quit <- struct{}{}
	}

	for {
		idx, _ := slice.FindFirstIf(this.workers, func(f func(...T) (T, bool)) bool {
			return f != nil
		})

		if idx == -1 {
			break
		} else {
			time.Sleep(this.sleepTime * time.Millisecond)
		}
	}

	for i := 1; i < len(this.inChans); i++ {
		close(this.inChans[i])
	}
}

// windUp Closes the entrance channel and
// waits for all the results to be processed.
func (this *Pipeline[T]) windUp() []T {
	if this.total.Load() == 0 {
		return []T{}
	}

	close(this.inChans[0]) // No more values to push
	arr := make([]T, 0, this.total.Load())
	for {
		v := <-this.inChans[len(this.inChans)-1]
		if arr = append(arr, v); len(arr) == int(this.total.Load()) {
			this.total.Store(0) // Reset the total count
			break               // All results received
		}
	}
	return arr
}
