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
	"sync"
	"sync/atomic"
	"time"

	"github.com/arcology-network/common-lib/exp/slice"
)

// Couple is a struct that contains two functions and two inChans.
// The first function consumes the input from the channel and produces
// an output that is sent to the second channel, which is consumed by the second function.
type Pipeline[T any] struct {
	name        string
	sleepTime   time.Duration
	channelSize int

	inChans     []chan T
	buffers     []*Slice[T] // Buffers to store the values temporarily before pushing to the downstream channel.
	workers     []func(T, *Slice[T]) ([]T, bool, bool)
	vacant      []*atomic.Bool
	totalVacant atomic.Uint64
	quit        atomic.Bool
	locks       []sync.Mutex

	counter atomic.Uint64
}

func NewPipeline[T any](name string, channelSize int, sleepTime time.Duration, workers ...func(T, *Slice[T]) ([]T, bool, bool)) *Pipeline[T] {
	pl := &Pipeline[T]{
		name:        name,
		sleepTime:   sleepTime,
		workers:     workers,
		inChans:     make([]chan T, len(workers)+1),
		buffers:     make([]*Slice[T], len(workers)),    // output outBuf
		vacant:      make([]*atomic.Bool, len(workers)), // if the worker is busy
		channelSize: channelSize,
		quit:        atomic.Bool{},
		locks:       make([]sync.Mutex, len(workers)),
	}

	for i := 0; i < len(pl.vacant); i++ {
		pl.vacant[i] = &atomic.Bool{}
		pl.vacant[i].Store(false)
	}

	for i := 0; i < len(pl.buffers); i++ {
		pl.buffers[i] = NewSlice[T]()
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
				if this.quit.Load() {
					this.workers[i] = nil
					return
				}

				this.DoTask(this.inChans[i], this.inChans[i+1], i, this.workers[i], this.buffers[i])
			}
		}(i)
	}
	return this
}

func (this *Pipeline[T]) DoTask(inQueue, outQueue chan T, i int, worker func(T, *Slice[T]) ([]T, bool, bool), buffer *Slice[T]) {
	select {
	case inv := <-inQueue:
		outVals, ok, isVacant := worker(inv, buffer)
		if ok {
			for j := 0; j < len(outVals); j++ {
				outQueue <- outVals[j] // Send to the downstream channel
			}
		}
		if isVacant {
			this.totalVacant.Add(1)
		}
		// this.vacant[i].Store(vacant)
	default:
		time.Sleep(this.sleepTime)
	}
}

// Redict all the results to be processed and returns. No new values can be pushed into the Pipeline
// after before Await() returns.
func (this *Pipeline[T]) RedirectTo(outChan chan T) {
	go func() {
		for {
			v := <-this.inChans[len(this.inChans)-1]
			outChan <- v
		}
	}()
}

// Push pushes values to the inChans.
func (this *Pipeline[T]) Push(vals ...T) {
	for _, v := range vals {
		this.inChans[0] <- v // Push to the first channel
	}
}

// Awiat waits for all the results to be processed and returns. No new values can be pushed into the Pipeline
// after before Await() returns.
func (this *Pipeline[T]) Await() []T {
	out := make([]T, 0, 1024)
	for {
		select {
		case v := <-this.inChans[len(this.inChans)-1]:
			out = append(out, v)
		default:
			if this.allDone() {
				fmt.Println("Total Time:", this.counter.Load()/(1000*1000), "ms")
				this.totalVacant.Store(0)
				// this.resetWorkerStatus()
				return out
			}
			time.Sleep(10 * this.sleepTime)
		}
	}
}

// Close the Pipeline and terminate all goroutines.
func (this *Pipeline[T]) Close() {
	this.quit.Store(true)

	for slice.CountIf(this.workers, func(i int, f *func(T, *Slice[T]) ([]T, bool, bool)) bool {
		return *f == nil
	}) != uint64(len(this.workers)) {
		time.Sleep(this.sleepTime)
	}
}

func (this *Pipeline[T]) allDone() bool {
	t0 := time.Now()
	flag := this.totalVacant.Load() == uint64(len(this.workers))
	this.counter.Add(uint64(time.Since(t0)))
	return flag
}
