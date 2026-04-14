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

package cachedkvstore

import (
	"fmt"
	"testing"
	"time"

	crdtcommon "github.com/arcology-network/common-lib/crdt/common"
	noncommutative "github.com/arcology-network/common-lib/crdt/noncommutative"
)

type testKVStore[K comparable, T any] struct {
	values           map[K]T
	currentLayerOnly bool
}

func newTestKVStore[K comparable, T any]() *testKVStore[K, T] {
	return &testKVStore[K, T]{values: map[K]T{}}
}

func (this *testKVStore[K, T]) Get(key K) (T, bool) {
	value, ok := this.values[key]
	return value, ok
}

func (this *testKVStore[K, T]) Has(key K) bool {
	_, ok := this.values[key]
	return ok
}

func (this *testKVStore[K, T]) GetBatch(keys []K) []T {
	values := make([]T, len(keys))
	for i, key := range keys {
		if value, ok := this.values[key]; ok {
			values[i] = value
		}
	}
	return values
}

func (this *testKVStore[K, T]) Len() uint64 {
	return uint64(len(this.values))
}

func (this *testKVStore[K, T]) Size() uint64 {
	return uint64(len(this.values))
}

func (this *testKVStore[K, T]) Set(key K, value T) {
	this.values[key] = value
}

func (this *testKVStore[K, T]) Delete(key K) {
	delete(this.values, key)
}

func (this *testKVStore[K, T]) SetBatch(keys []K, values []T) {
	for i := 0; i < len(keys) && i < len(values); i++ {
		this.values[keys[i]] = values[i]
	}
}

func (this *testKVStore[K, T]) DeleteBatch(keys []K) {
	for _, key := range keys {
		delete(this.values, key)
	}
}

func (this *testKVStore[K, T]) Precommit() error {
	return nil
}

func (this *testKVStore[K, T]) Commit(bool, uint64) error {
	return nil
}

func (this *testKVStore[K, T]) SetLocalOnly(yes bool) {
	this.currentLayerOnly = yes
}

func (this *testKVStore[K, T]) LocalOnly() bool {
	return this.currentLayerOnly
}

func newStringValue(v string) crdtcommon.CRDT {
	return noncommutative.NewString(v)
}

func newProfiledString(v string) *Entry[crdtcommon.CRDT] {
	value := newStringValue(v)
	return &Entry[crdtcommon.CRDT]{
		Value: value,
		Stat: Stat{
			sizeInMem: value.MemSize(),
		},
	}
}

func newBenchmarkBackend(entryCount int) (*testKVStore[string, crdtcommon.CRDT], []string, []crdtcommon.CRDT) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	keys := make([]string, entryCount)
	values := make([]crdtcommon.CRDT, entryCount)
	for i := 0; i < entryCount; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := newStringValue("value")
		keys[i] = key
		values[i] = value
		backend.Set(key, value)
	}
	return backend, keys, values
}

func TestStoreCachesReads(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	backend.Set("alpha", newStringValue("one"))
	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 1024, nil)

	first, ok := store.Get("alpha")
	if !ok || first == nil {
		t.Fatalf("expected value from backend on first read")
	}

	cached, ok := store.ConcurrentMap.Get("alpha")
	if !ok || cached == nil {
		t.Fatalf("expected first read to promote value into first layer")
	}

	second, ok := store.Get("alpha")
	if !ok || second == nil {
		t.Fatalf("expected cached value on second read")
	}
	if second != cached {
		t.Fatalf("expected second read to return the cached value")
	}
}

func TestStoreSkipsCachingOversizedEntry(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	backend.Set("alpha", newStringValue("one"))

	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 1024, func(v crdtcommon.CRDT) uint64 { return 2048 })

	first, ok := store.Get("alpha")
	if !ok || first == nil {
		t.Fatalf("expected oversized value from backend on first read")
	}
	if cached, ok := store.ConcurrentMap.Get("alpha"); ok && cached != nil {
		t.Fatalf("expected oversized entry to stay out of the first layer")
	}

	second, ok := store.Get("alpha")
	if !ok || second == nil {
		t.Fatalf("expected oversized value from backend on second read")
	}
	if cached, ok := store.ConcurrentMap.Get("alpha"); ok && cached != nil {
		t.Fatalf("expected oversized entry not to be cached after repeated reads")
	}
}

func TestStoreLayeredReadLeavesBackendUntouched(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	backend.Set("alpha", newStringValue("one"))

	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 1024, nil)

	first, ok := store.Get("alpha")
	if !ok || first == nil {
		t.Fatalf("expected first layer to read through backend")
	}

	second, ok := store.Get("alpha")
	if !ok || second == nil {
		t.Fatalf("expected first layer to return cached backend value")
	}
	if second != first {
		t.Fatalf("expected repeated read to return the cached entry")
	}

	local := newProfiledString("two")
	store.Set("beta", local)
	if !store.Has("beta") {
		t.Fatalf("expected local write to stay in first layer")
	}
	if backend.Has("beta") {
		t.Fatalf("expected backend to remain unchanged")
	}
	if value, ok := store.Get("beta"); !ok || value != local {
		t.Fatalf("expected local value to remain in cache")
	}
}

func TestStoreTracksVisitsGenerically(t *testing.T) {
	store := NewCachedKVStore[string, crdtcommon.CRDT](nil, 4096, nil)
	alpha := newProfiledString("one")
	store.UpdateVersion(11)

	store.Set("alpha", alpha)
	if alpha.visits != 1 {
		t.Fatalf("expected set to increment visits, got %d", alpha.visits)
	}
	if alpha.firstLoaded != 11 {
		t.Fatalf("expected set to assign the current version as firstLoaded")
	}

	value, ok := store.Get("alpha")
	if !ok || value == nil {
		t.Fatalf("expected get to return cached value")
	}
	if value.visits != 2 {
		t.Fatalf("expected get to increment visits, got %d", value.visits)
	}

	replacement := newProfiledString("two")
	store.Set("alpha", replacement)
	if replacement.visits != 3 {
		t.Fatalf("expected replacement to inherit and increment visits, got %d", replacement.visits)
	}
	if replacement.firstLoaded != alpha.firstLoaded {
		t.Fatalf("expected replacement to preserve firstLoaded for an existing cache key")
	}

	batch := store.GetBatch([]string{"alpha"})
	if len(batch) != 1 || batch[0] == nil {
		t.Fatalf("expected batch get to return cached value")
	}
	if batch[0].visits != 4 {
		t.Fatalf("expected batch get to increment visits, got %d", batch[0].visits)
	}

	beta := newProfiledString("three")
	store.UpdateVersion(21)
	store.SetBatch([]string{"beta"}, []*Entry[crdtcommon.CRDT]{beta})
	if beta.visits != 1 {
		t.Fatalf("expected batch set to increment visits, got %d", beta.visits)
	}
	if beta.firstLoaded != 21 {
		t.Fatalf("expected batch set to assign the current version as firstLoaded")
	}
	if beta.firstLoaded == alpha.firstLoaded {
		t.Fatalf("expected distinct version-derived firstLoaded values for distinct entries")
	}
	if cached, ok := store.ConcurrentMap.Get("beta"); !ok || cached != beta {
		t.Fatalf("expected batch set to update cache entry")
	}

	backend := newTestKVStore[string, crdtcommon.CRDT]()
	backend.Set("backend", newStringValue("backend"))

	layered := NewCachedKVStore[string, crdtcommon.CRDT](backend, 4096, nil)
	layered.UpdateVersion(33)
	fetched, ok := layered.Get("backend")
	if !ok || fetched == nil {
		t.Fatalf("expected backend get to succeed")
	}
	if fetched.visits != 1 {
		t.Fatalf("expected backend get to initialize visits, got %d", fetched.visits)
	}
	if fetched.firstLoaded != 33 {
		t.Fatalf("expected backend get to assign version-derived firstLoaded")
	}

	again, ok := layered.Get("backend")
	if !ok || again != fetched {
		t.Fatalf("expected backend value to be cached after first read")
	}
	if again.visits != 2 {
		t.Fatalf("expected cached backend value to increment visits, got %d", again.visits)
	}
}

func TestStoreEvictsWhenFull(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 1024, func(v crdtcommon.CRDT) uint64 { return 700 })

	alpha := newProfiledString("one")
	beta := newProfiledString("two")

	store.Set("alpha", alpha)
	store.Set("beta", beta)

	if store.ConcurrentMap.Length() != 2 {
		t.Fatalf("expected both entries to stay in the first layer before eviction")
	}

	store.Evict()

	if backend.Len() != 0 {
		t.Fatalf("expected backend to remain unchanged, got %d", backend.Len())
	}
	if store.Size() > 1024 {
		t.Fatalf("expected eviction to keep cache within cap, got %d", store.Size())
	}
	if store.ConcurrentMap.Length() >= 2 {
		t.Fatalf("expected eviction to remove at least one cached entry")
	}
}

func TestStoreBatchAndDelete(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	backend.Set("alpha", newStringValue("one"))
	backend.Set("beta", newStringValue("two"))
	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 1024, nil)

	values := store.GetBatch([]string{"alpha", "beta", "missing"})
	if len(values) != 3 || values[0] == nil || values[1] == nil || values[2] != nil {
		t.Fatalf("unexpected batch result: %#v", values)
	}
	if cached, ok := store.ConcurrentMap.Get("alpha"); !ok || cached == nil {
		t.Fatalf("expected batch read to populate cache for alpha")
	}
	if cached, ok := store.ConcurrentMap.Get("beta"); !ok || cached == nil {
		t.Fatalf("expected batch read to populate cache for beta")
	}

	gamma := newProfiledString("three")
	store.Set("gamma", gamma)
	if backend.Has("gamma") || !store.Has("gamma") {
		t.Fatalf("expected set to stay local")
	}

	if err := store.Delete("alpha"); err != nil {
		t.Fatalf("delete alpha: %v", err)
	}

	if err := store.DeleteBatch([]string{"beta", "gamma"}); err != nil {
		t.Fatalf("delete batch: %v", err)
	}
	if store.Has("gamma") {
		t.Fatalf("expected local-only key to stay absent after delete")
	}
	if !store.Has("alpha") || !store.Has("beta") {
		t.Fatalf("expected backend-backed keys to surface again after cache delete")
	}
	if _, ok := store.GetRaw("alpha"); ok {
		t.Fatalf("expected delete to evict alpha from first layer")
	}
	if _, ok := store.GetRaw("beta"); ok {
		t.Fatalf("expected delete to evict beta from first layer")
	}
	if _, ok := store.GetRaw("gamma"); ok {
		t.Fatalf("expected delete to evict gamma from first layer")
	}
	if backend.Len() != 2 {
		t.Fatalf("expected backend to remain unchanged after delete, got %d items", backend.Len())
	}

	if value, ok := store.Get("alpha"); !ok || value == nil {
		t.Fatalf("expected backend alpha to be readable after cache delete")
	}
	if value, ok := store.Get("beta"); !ok || value == nil {
		t.Fatalf("expected backend beta to be readable after cache delete")
	}
	if value, ok := store.Get("gamma"); ok || value != nil {
		t.Fatalf("expected local-only gamma to stay absent after delete")
	}
}

func TestStoreCurrentLayerOnlySkipsBackend(t *testing.T) {
	localCache := NewCachedKVStore[string, crdtcommon.CRDT](nil, 1024, nil)
	localCache.Set("cached", newProfiledString("value"))
	if value, ok := localCache.Get("cached"); !ok || value == nil {
		t.Fatalf("expected local store read to return cached value")
	}

	backend := newTestKVStore[string, crdtcommon.CRDT]()
	backend.Set("alpha", newStringValue("one"))
	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 1024, nil)
	store.SetLocalOnly(true)

	if value, ok := store.Get("alpha"); ok || value != nil {
		t.Fatalf("expected current-layer-only get to skip backend")
	}

	values := store.GetBatch([]string{"alpha"})
	if len(values) != 1 || values[0] != nil {
		t.Fatalf("expected batch get to skip backend, got %#v", values)
	}

	if store.Has("alpha") {
		t.Fatalf("expected has to skip backend in current-layer-only mode")
	}

	store.Set("alpha", newProfiledString("local"))
	if value, ok := store.Get("alpha"); !ok || value == nil {
		t.Fatalf("expected current-layer-only get to return cached value")
	}
}

func TestStoreSetGet1MillionEntries(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newProfiledString("value").Size()

	t0 := time.Now()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, uint64(entryCount)*sampleSize, nil)
	for _, key := range keys {
		if _, ok := store.Get(key); !ok {
			t.Fatalf("expected warm cache get to succeed for %s", key)
		}
	}
	store.SetLocalOnly(true)

	for i := 0; i < entryCount; i++ {
		key := keys[i%entryCount]
		value, ok := store.Get(key)
		if !ok || value == nil {
			t.Fatalf("expected first-layer hit for %s", key)
		}
	}
	fmt.Printf("Get 1 Million entries: %v\n", time.Since(t0))
}

func TestStoreSetGet1MillionEntries1024InCache(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newProfiledString("value").Size()

	t0 := time.Now()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 1024*sampleSize, nil)
	for _, key := range keys {
		if _, ok := store.Get(key); !ok {
			t.Fatalf("expected warm cache get to succeed for %s", key)
		}
	}

	for i := 0; i < entryCount; i++ {
		key := keys[i%entryCount]
		value, ok := store.Get(key)
		if !ok || value == nil {
			t.Fatalf("expected layered get to succeed for %s", key)
		}
	}

	store.Set("committed", newProfiledString("committed"))
	if backend.Has("committed") {
		t.Fatalf("expected backend to remain unchanged after local set")
	}

	fmt.Printf("Get 1 Million entries, 1024 Entry cache: %v\n", time.Since(t0))
}

func TestStoreSetGet1MillionEntries2048InCache(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newProfiledString("value").Size()

	t0 := time.Now()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 2048*sampleSize, nil)
	for _, key := range keys {
		if _, ok := store.Get(key); !ok {
			t.Fatalf("expected warm cache get to succeed for %s", key)
		}
	}

	for i := 0; i < entryCount; i++ {
		key := keys[i%entryCount]
		value, ok := store.Get(key)
		if !ok || value == nil {
			t.Fatalf("expected layered get to succeed for %s", key)
		}
	}

	store.Set("committed", newProfiledString("committed"))
	if backend.Has("committed") {
		t.Fatalf("expected backend to remain unchanged after local set")
	}

	fmt.Printf("Get 1 Million entries, 2048 Entry cache: %v\n", time.Since(t0))
}

func TestStoreSetGet1MillionEntriesAllInCache(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newProfiledString("value").Size()

	t0 := time.Now()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, uint64(entryCount)*sampleSize, nil)
	for _, key := range keys {
		if _, ok := store.Get(key); !ok {
			t.Fatalf("expected warm cache get to succeed for %s", key)
		}
	}
	store.SetLocalOnly(true)

	for i := 0; i < entryCount; i++ {
		key := keys[i%entryCount]
		value, ok := store.Get(key)
		if !ok || value == nil {
			t.Fatalf("expected first-layer hit for %s", key)
		}
	}

	store.Set("committed", newProfiledString("committed"))
	if backend.Has("committed") {
		t.Fatalf("expected backend to remain unchanged after local set")
	}

	fmt.Printf("Get 1 Million entries, all in cache: %v\n", time.Since(t0))
}
