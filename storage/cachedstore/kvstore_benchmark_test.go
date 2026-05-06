/*
*   Copyright (c) 2026 Arcology Network

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

package cachedstore

import (
	"testing"

	crdtcommon "github.com/arcology-network/common-lib/crdt/common"
)

// newBenchStore returns a CachedStore backed by a pre-populated testKVStore.
// cacheCap controls how many entries fit in the cache (each entry costs 1 unit).
func newBenchStore(b *testing.B, entryCount int, cacheCap uint64) (*CachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT], []string) {
	b.Helper()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore(
		backend,
		codec,
		cacheCap,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)
	return store, keys
}

// ---------------------------------------------------------------------------
// 10K benchmarks
// ---------------------------------------------------------------------------

func BenchmarkGet_CacheHit_10K(b *testing.B) {
	const entryCount = 10_000
	store, keys := newBenchStore(b, entryCount, entryCount)
	for _, k := range keys {
		store.Get(k)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get(keys[i%entryCount])
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkGet_BackendMiss_10K(b *testing.B) {
	const entryCount = 10_000
	store, keys := newBenchStore(b, entryCount, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get(keys[i%entryCount])
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkGet_MixedCache_10K(b *testing.B) {
	const entryCount = 10_000
	store, keys := newBenchStore(b, entryCount, entryCount/2)
	for i := 0; i < entryCount/2; i++ {
		store.Get(keys[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get(keys[i%entryCount])
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkSet_10K(b *testing.B) {
	const entryCount = 10_000
	store, keys := newBenchStore(b, entryCount, entryCount)
	val := newStringValue("bench")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Set(keys[i%entryCount], val)
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkGetBatch_CacheHit_10K(b *testing.B) {
	const entryCount = 10_000
	const batchSize = 100
	store, keys := newBenchStore(b, entryCount, entryCount)
	for _, k := range keys {
		store.Get(k)
	}
	batch := keys[:batchSize]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetBatch(batch)
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkGetBatch_BackendMiss_10K(b *testing.B) {
	const entryCount = 10_000
	const batchSize = 100
	store, keys := newBenchStore(b, entryCount, 0)
	batch := keys[:batchSize]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetBatch(batch)
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkSetBatch_10K(b *testing.B) {
	const entryCount = 10_000
	const batchSize = 100
	store, keys := newBenchStore(b, entryCount, entryCount)
	batch := keys[:batchSize]
	vals := make([]crdtcommon.CRDT, batchSize)
	val := newStringValue("bench")
	for i := range vals {
		vals[i] = val
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.SetBatch(batch, vals)
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

// ---------------------------------------------------------------------------
// 1M benchmarks
// ---------------------------------------------------------------------------

func BenchmarkGet_CacheHit_1M(b *testing.B) {
	const entryCount = 1_000_000
	store, keys := newBenchStore(b, entryCount, entryCount)
	for _, k := range keys {
		store.Get(k)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get(keys[i%entryCount])
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkGet_BackendMiss_1M(b *testing.B) {
	const entryCount = 1_000_000
	store, keys := newBenchStore(b, entryCount, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get(keys[i%entryCount])
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkGet_MixedCache_1M(b *testing.B) {
	const entryCount = 1_000_000
	const cacheCap = 1024
	store, keys := newBenchStore(b, entryCount, cacheCap)
	for i := 0; i < cacheCap; i++ {
		store.Get(keys[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get(keys[i%entryCount])
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkSet_1M(b *testing.B) {
	const entryCount = 1_000_000
	store, keys := newBenchStore(b, entryCount, entryCount)
	val := newStringValue("bench")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Set(keys[i%entryCount], val)
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkGetBatch_CacheHit_1M(b *testing.B) {
	const entryCount = 1_000_000
	const batchSize = 1000
	store, keys := newBenchStore(b, entryCount, entryCount)
	for _, k := range keys {
		store.Get(k)
	}
	batch := keys[:batchSize]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetBatch(batch)
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkGetBatch_BackendMiss_1M(b *testing.B) {
	const entryCount = 1_000_000
	const batchSize = 1000
	store, keys := newBenchStore(b, entryCount, 0)
	batch := keys[:batchSize]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetBatch(batch)
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}

func BenchmarkSetBatch_1M(b *testing.B) {
	const entryCount = 1_000_000
	const batchSize = 1000
	store, keys := newBenchStore(b, entryCount, entryCount)
	batch := keys[:batchSize]
	vals := make([]crdtcommon.CRDT, batchSize)
	val := newStringValue("bench")
	for i := range vals {
		vals[i] = val
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.SetBatch(batch, vals)
	}
	b.ReportMetric(b.Elapsed().Seconds(), "s/total")
}
