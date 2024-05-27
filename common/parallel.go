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
	"math"
	"runtime"
	"sync"
)

// GenerateRanges generates a slice of ranges based on the length and number of threads.
// Each range represents a portion of the total length that can be processed by a single thread.
func GenerateRanges(length int, numThreads int) []int {
	numThreads = Min(Min(numThreads, length), runtime.NumCPU()) // limit the number of threads to the number of CPUs

	ranges := make([]int, 0, numThreads+1)
	step := int(math.Ceil(float64(length) / float64(numThreads)))
	for i := 0; i <= numThreads; i++ {
		ranges = append(ranges, int(math.Min(float64(step*i), float64(length))))
	}
	return ranges
}

// ParallelExecute executes the given tasks in parallel using goroutines.
// It waits for all the tasks to complete before returning.
func ParallelExecute(tasks ...interface{}) {
	var wg sync.WaitGroup
	for i := 0; i < len(tasks); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tasks[i].(func())()
		}(i)
	}
	wg.Wait()
}

// ParallelWorker divides the total work into multiple ranges and assigns each range to a worker function.
// The worker function is called in parallel for each range.
// The number of threads determines the number of ranges and worker functions.
func ParallelWorker(total, nThds int, worker func(start, end, idx int, args ...interface{}), args ...interface{}) {
	idxRanges := GenerateRanges(total, nThds)
	var wg sync.WaitGroup
	for i := 0; i < len(idxRanges)-1; i++ {
		wg.Add(1)
		go func(start int, end int, idx int) {
			defer wg.Done()
			if start != end {
				worker(start, end, idx, args)
			}
		}(idxRanges[i], idxRanges[i+1], i)
	}
	wg.Wait()
}

// ParallelForeach applies a function to each element in a slice in parallel using multiple threads.
func ParallelFor(v0, v1, nThds int, do func(int)) {
	processor := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			do(i + v0)
		}
	}
	ParallelWorker(v1-v0, nThds, processor)
}
