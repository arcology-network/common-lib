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
	"time"
)

type Triple[T0, T1, T2, T3 any] struct {
	*Couple[T0, T1, T2]
	_4thChan chan T3
	Third    func(T2) T3
}

// NewTriple creates a new Triple instance, that is a combination of three functions.
// Each function consumes the output of the previous function.
func NewTriple[T0, T1, T2, T3 any](
	first func(T0) T1,
	second func(T1) T2,
	third func(T2) T3,
	bufferSize int,
	sleepTime time.Duration,
	timeout time.Duration) *Triple[T0, T1, T2, T3] {
	return &Triple[T0, T1, T2, T3]{
		Couple:   NewCouple[T0, T1, T2](first, second, bufferSize, sleepTime, timeout),
		_4thChan: make(chan T3, bufferSize),
		Third:    third,
	}
}

// Push pushes values to the channels.
func (this *Triple[T0, T1, T2, T3]) Push(vals ...T0) {
	this.total.Add(int64(len(vals)))
	for _, val := range vals {
		this._1stChan <- val
	}
}

// Start starts the goroutines.
func (this *Triple[T0, T1, T2, T3]) Start() *Triple[T0, T1, T2, T3] {
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

	go func() {
		for {
			select {
			case val := <-this._3rdChan:
				this._4thChan <- this.Third(val)
			default:
				time.Sleep(this.sleepTime * time.Millisecond)
			}
		}
	}()
	return this
}

// Await waits for all the goroutines to finish.
func (this *Triple[T0, T1, T2, T3]) Await() ([]T3, error) {
	close(this._1stChan)
	arr := make([]T3, 0, this.total.Load())
	for {
		select {
		case v, ok := <-this._4thChan:
			if !ok { // Channel closed
				return arr, errors.New("Error: Channel closed!")
			}
			arr = append(arr, v)
			if len(arr) == int(this.total.Load()) {
				this._1stChan = make(chan T0, cap(this._1stChan)) // Reset the channel
				return arr, nil                                   // All results received
			}

		case <-time.After(this.timeout * time.Millisecond):
			return arr, errors.New("Error: Time out!") // Timeout if all results are not received within sleepTime
		}
	}
}
