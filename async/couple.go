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
	"errors"
	"sync/atomic"
	"time"
)

const (
	TIME_OUT_RATE = 100
)

// Couple is a struct that contains two functions and two channels.
// The first function consumes the input from the channel and produces
// an output that is sent to the second channel, which is consumed by the second function.
type Couple[T0, T1, T2 any] struct {
	total     atomic.Int64
	timeout   time.Duration
	sleepTime time.Duration
	_1stChan  chan T0
	_2ndChan  chan T1
	_3rdChan  chan T2
	First     func(T0) T1
	Second    func(T1) T2
}

func NewCouple[T0, T1, T2 any](
	First func(T0) T1,
	Second func(T1) T2,
	bufferSize int,
	sleepTime time.Duration,
	timeout time.Duration) *Couple[T0, T1, T2] {
	return &Couple[T0, T1, T2]{
		sleepTime: sleepTime,
		timeout:   timeout,
		_1stChan:  make(chan T0, bufferSize),
		_2ndChan:  make(chan T1, bufferSize),
		_3rdChan:  make(chan T2, bufferSize),
		First:     First,
		Second:    Second,
	}
}

// Push pushes values to the channels.
func (this *Couple[T0, T1, T2]) Push(vals ...T0) {
	this.total.Add(int64(len(vals)))
	for _, val := range vals {
		this._1stChan <- val
	}
}

// Start starts the goroutines.
func (this *Couple[T0, T1, T2]) Start() *Couple[T0, T1, T2] {
	go func() {
		for {
			select {
			case val := <-this._1stChan:
				this._2ndChan <- this.First(val)
			default:
				time.Sleep(this.sleepTime * time.Millisecond)
			}
		}
	}()

	go func() {
		for {
			select {
			case val := <-this._2ndChan:
				this._3rdChan <- this.Second(val)
			default:
				time.Sleep(this.sleepTime * time.Millisecond)
			}
		}
	}()
	return this
}

func (this *Couple[T0, T1, T2]) Await() ([]T2, error) {
	close(this._1stChan)
	arr := make([]T2, 0, this.total.Load())
	for {
		select {
		case v, ok := <-this._3rdChan:
			if !ok { // Channel closed
				return arr, errors.New("Error: Channel closed!")
			}
			arr = append(arr, v)
			if len(arr) == int(this.total.Load()) {
				this._1stChan = make(chan T0, cap(this._1stChan)) // Reset the channel
				return arr, nil                                   // All results received
			}

		case <-time.After(this.timeout * time.Millisecond):
			return arr, nil // Timeout if all results are not received within sleepTime
		}
	}
}
