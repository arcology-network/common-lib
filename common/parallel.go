package common

import (
	"math"
	"sync"
)

func GenerateRanges(length int, numThreads int) []int {
	ranges := make([]int, 0, numThreads+1)
	step := int(math.Ceil(float64(length) / float64(numThreads)))
	for i := 0; i <= numThreads; i++ {
		ranges = append(ranges, int(math.Min(float64(step*i), float64(length))))
	}
	return ranges
}

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
