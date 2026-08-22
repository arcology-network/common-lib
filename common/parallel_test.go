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

package common

import (
	"sync/atomic"
	"testing"
)

func TestParallelExecute(t *testing.T) {
	var counter atomic.Int64
	n := 10
	tasks := make([]func(), n)
	for i := range tasks {
		tasks[i] = func() { counter.Add(1) }
	}
	ParallelExecute(tasks...)
	if got := counter.Load(); got != int64(n) {
		t.Errorf("expected %d tasks executed, got %d", n, got)
	}
}

func TestParallelWorker_allElementsProcessedOnce(t *testing.T) {
	const total = 100
	hits := make([]atomic.Int32, total)
	ParallelWorker(total, 8, func(start, end, idx int, args ...any) {
		for i := start; i < end; i++ {
			hits[i].Add(1)
		}
	})
	for i := range hits {
		if hits[i].Load() != 1 {
			t.Errorf("index %d processed %d times, want 1", i, hits[i].Load())
		}
	}
}

// total not evenly divisible by nThds — the tricky stride case
func TestParallelWorker_unevenDivision(t *testing.T) {
	const total = 4
	hits := make([]atomic.Int32, total)
	ParallelWorker(total, 3, func(start, end, idx int, args ...any) {
		for i := start; i < end; i++ {
			hits[i].Add(1)
		}
	})
	for i := range hits {
		if hits[i].Load() != 1 {
			t.Errorf("index %d processed %d times, want 1", i, hits[i].Load())
		}
	}
}

func TestParallelWorker_nThdsExceedsTotal(t *testing.T) {
	const total = 3
	hits := make([]atomic.Int32, total)
	ParallelWorker(total, 100, func(start, end, idx int, args ...any) {
		for i := start; i < end; i++ {
			hits[i].Add(1)
		}
	})
	for i := range hits {
		if hits[i].Load() != 1 {
			t.Errorf("index %d processed %d times, want 1", i, hits[i].Load())
		}
	}
}

func TestParallelWorker_zeroTotal(t *testing.T) {
	called := false
	ParallelWorker(0, 4, func(start, end, idx int, args ...any) { called = true })
	if called {
		t.Error("worker should not be called when total is 0")
	}
}

func TestParallelWorker_zeroThreads(t *testing.T) {
	called := false
	ParallelWorker(10, 0, func(start, end, idx int, args ...any) { called = true })
	if called {
		t.Error("worker should not be called when nThds is 0")
	}
}

func TestParallelFor_allIndicesVisited(t *testing.T) {
	const v0, v1 = 5, 25
	hits := make([]atomic.Int32, v1)
	err := ParallelFor(v0, v1, 4, func(i int) { hits[i].Add(1) })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := v0; i < v1; i++ {
		if hits[i].Load() != 1 {
			t.Errorf("index %d visited %d times, want 1", i, hits[i].Load())
		}
	}
	// indices outside [v0, v1) must not be touched
	for i := 0; i < v0; i++ {
		if hits[i].Load() != 0 {
			t.Errorf("index %d outside range was visited", i)
		}
	}
}

func TestParallelFor_invalidRange(t *testing.T) {
	if err := ParallelFor(10, 5, 4, func(int) {}); err == nil {
		t.Error("expected error for v1 <= v0")
	}
	if err := ParallelFor(5, 5, 4, func(int) {}); err == nil {
		t.Error("expected error for equal v0 and v1")
	}
}

func TestParallelFor_zeroThreads(t *testing.T) {
	if err := ParallelFor(0, 10, 0, func(int) {}); err == nil {
		t.Error("expected error for nThds <= 0")
	}
}
