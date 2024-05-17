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

	// outChan chan T
	inChans      []chan T
	buffers      []*Slice[T] // Buffers to store the values temporarily before pushing to the downstream channel.
	workers      []func(T, *Slice[T]) ([]T, bool)
	locks        []sync.Mutex
	isWorkerBusy []*atomic.Bool
	quit         atomic.Bool
}

func NewPipeline[T any](name string, channelSize int, sleepTime time.Duration, workers ...func(T, *Slice[T]) ([]T, bool)) *Pipeline[T] {
	pl := &Pipeline[T]{
		name:         name,
		sleepTime:    sleepTime,
		workers:      workers,
		inChans:      make([]chan T, len(workers)+1),
		buffers:      make([]*Slice[T], len(workers)),    // output outBuf
		isWorkerBusy: make([]*atomic.Bool, len(workers)), // if the worker is busy
		channelSize:  channelSize,
		quit:         atomic.Bool{},
		locks:        make([]sync.Mutex, len(workers)),
	}

	for i := 0; i < len(pl.isWorkerBusy); i++ {
		pl.isWorkerBusy[i] = &atomic.Bool{}
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
				if !this.isWorkerBusy[i].Load() && this.quit.Load() {
					return
				}

				// No job in the channel, this may affect performance
				if len(this.inChans[i]) == 0 {
					time.Sleep(10 * time.Millisecond)
					continue
				}

				this.DoTask(this.inChans[i], this.inChans[i+1], i, this.workers[i], this.buffers[i], this.isWorkerBusy[i])
			}
		}(i)
	}
	return this
}

func (this *Pipeline[T]) DoTask(inQueue, outQueue chan T, i int, worker func(T, *Slice[T]) ([]T, bool), buffer *Slice[T], flag *atomic.Bool) bool {
	this.locks[i].Lock()
	defer this.locks[i].Unlock()

	flag.Store(true)
	if inv, ok := <-inQueue; ok {
		if outVals, ok := worker(inv, buffer); ok {
			for j := 0; j < len(outVals); j++ {
				outQueue <- outVals[j] // Send to the downstream channel
			}
		}
	}
	flag.Store(false)
	return !this.isWorkerBusy[i].Load() && this.quit.Load()
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
	arr := this.windUp()
	this.inChans[0] = make(chan T, this.channelSize) // Reopen the entrance channel
	return arr
}

// Close the Pipeline and terminate all goroutines.
func (this *Pipeline[T]) Close() {
	this.windUp()
	this.quit.Store(true)

	// The entrance channel is closed already by the windUp function.
	for i := 1; i < len(this.inChans); i++ {
		close(this.inChans[i])
	}
}

// windUp Closes the entrance channel and
// waits for all the results to be processed.
func (this *Pipeline[T]) windUp() []T {
	close(this.inChans[0]) // No more values to push

	out := make([]T, 0, 1024)
	for {
		select {
		case v := <-this.inChans[len(this.inChans)-1]:
			out = append(out, v)
		default:
			if this.allDone() {
				return out
			}
		}
	}
}

func (this *Pipeline[T]) allDone() bool {
	activeWorkers := slice.CountIf[*atomic.Bool, int](this.isWorkerBusy, func(_ int, b **atomic.Bool) bool {
		return (*b).Load()
	})

	activeChans := slice.CountIf[chan T, int](this.inChans[:len(this.inChans)-1], func(_ int, b *chan T) bool {
		return len(*b) > 0
	})

	activeBuffs := slice.CountIf[*Slice[T], int](this.buffers, func(_ int, buffer **Slice[T]) bool {
		return (*buffer).Length() > 0
	})
	return activeWorkers+activeChans+activeBuffs == 0
}
