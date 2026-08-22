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

// Package common provides common utility functions for parallel execution.

package common

import (
	"errors"
	"math"
	"sync"
)

// ParallelExecute executes the given tasks in parallel using goroutines.
// It waits for all the tasks to complete before returning.
func ParallelExecute(tasks ...func()) {
	ParallelWorker(len(tasks), len(tasks), func(start, end, _ int, _ ...any) {
		for i := start; i < end; i++ {
			tasks[i]()
		}
	})
}

func ParallelWorker(total, nThds int, worker func(start, end, idx int, args ...any), args ...any) {
	if total <= 0 || nThds <= 0 {
		return
	}

	nThds = min(total, nThds)
	stride := int(math.Ceil(float64(total) / float64(nThds)))

	var wg sync.WaitGroup
	for i := 0; i < nThds; i++ {
		end := int(min((i+1)*stride, total))
		wg.Add(1)
		go func(start int, end int, idx int) {
			defer wg.Done()
			worker(start, end, idx, args...)
		}(i*stride, end, i)
	}
	wg.Wait()
}

// ParallelForeach applies a function to each element in a slice in parallel using multiple threads.
func ParallelFor(v0, v1, nThds int, do func(int)) error {
	if nThds <= 0 {
		return errors.New("The number of threads must be greater than 0")
	}

	processor := func(start, end, index int, args ...any) {
		for i := start; i < end; i++ {
			do(i + v0)
		}
	}

	if v1 > v0 {
		ParallelWorker((v1 - v0), nThds, processor)
		return nil
	}
	return errors.New("Range must be positive")
}
