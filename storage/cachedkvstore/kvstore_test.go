package cachedkvstore

import (
	"fmt"
	"testing"
	"time"

	crdtcommon "github.com/arcology-network/common-lib/crdt/common"
	noncommutative "github.com/arcology-network/common-lib/crdt/noncommutative"
)

func newProfiledString(v string) *Entry[crdtcommon.CRDT] {
	value := noncommutative.NewString(v)
	return &Entry[crdtcommon.CRDT]{
		Value: value,
		Stat: Stat{
			sizeInMem: value.MemSize(),
		},
	}
}

func newBenchmarkBackend(entryCount int) (*CachedKVStore[string, crdtcommon.CRDT], []string, []*Entry[crdtcommon.CRDT]) {
	sampleSize := newProfiledString("value").Size()
	backend := NewCachedKVStore[string, crdtcommon.CRDT](nil, uint64(entryCount)*sampleSize, nil)
	keys := make([]string, entryCount)
	values := make([]*Entry[crdtcommon.CRDT], entryCount)
	for i := 0; i < entryCount; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := newProfiledString("value")
		keys[i] = key
		values[i] = value
		backend.Set(key, value)
	}
	_ = backend.Commit(false, 0)
	return backend, keys, values
}

func TestStoreCachesReads(t *testing.T) {
	backend := NewCachedKVStore[string, crdtcommon.CRDT](nil, 4096, nil)
	backend.Set("alpha", newProfiledString("one"))
	_ = backend.Commit(false, 0)
	store := NewCachedKVStore(backend, 1024, nil)

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
	backend := NewCachedKVStore[string, crdtcommon.CRDT](nil, 4096, nil)
	oversized := newProfiledString("one")
	backend.Set("alpha", oversized)
	_ = backend.Commit(false, 0)

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

func TestStoreLayeredReadAndCommitWriteback(t *testing.T) {
	backend := NewCachedKVStore[string, crdtcommon.CRDT](nil, 4096, nil)
	backend.Set("alpha", newProfiledString("one"))

	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 1024, nil)

	first, ok := store.Get("alpha")
	if !ok || first == nil {
		t.Fatalf("expected first layer to read through backend")
	}

	second, ok := store.Get("alpha")
	if !ok || second == nil {
		t.Fatalf("expected first layer to return cached backend value")
	}

	local := newProfiledString("two")
	store.Set("beta", local)
	if !store.Has("beta") {
		t.Fatalf("expected local write to stay in first layer")
	}
	if backend.Has("beta") {
		t.Fatalf("expected backend to remain unchanged before commit")
	}

	if err := store.Precommit(); err != nil {
		t.Fatalf("unexpected precommit error: %v", err)
	}
	if backend.Has("beta") {
		t.Fatalf("expected precommit not to flush to backend")
	}

	if err := store.Commit(true, 1); err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}
	if !backend.Has("beta") {
		t.Fatalf("expected commit to flush staged write to backend")
	}
}

func TestStoreCommitEvictsWhenFull(t *testing.T) {
	backend := NewCachedKVStore[string, crdtcommon.CRDT](nil, 4096, nil)
	store := NewCachedKVStore[string, crdtcommon.CRDT](backend, 1024, func(v crdtcommon.CRDT) uint64 { return 700 })

	alpha := newProfiledString("one")
	beta := newProfiledString("two")

	store.Set("alpha", alpha)
	store.Set("beta", beta)

	if store.ConcurrentMap.Length() != 2 {
		t.Fatalf("expected both entries to stay in the first layer before commit")
	}

	if err := store.Commit(true, 1); err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}

	if backend.Len() != 2 {
		t.Fatalf("expected both entries to be committed to backend, got %d", backend.Len())
	}
	if store.Size() > 1024 {
		t.Fatalf("expected eviction to keep cache within cap, got %d", store.Size())
	}
	if store.ConcurrentMap.Length() >= 2 {
		t.Fatalf("expected commit-time eviction to remove at least one cached entry")
	}
}

func TestStoreBatchAndDelete(t *testing.T) {
	backend := NewCachedKVStore[string, crdtcommon.CRDT](nil, 4096, nil)
	backend.Set("alpha", newProfiledString("one"))
	backend.Set("beta", newProfiledString("two"))
	_ = backend.Commit(false, 0)
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
		t.Fatalf("expected set to stay local until commit")
	}

	if err := store.Precommit(); err != nil {
		t.Fatalf("unexpected precommit error: %v", err)
	}
	if err := store.Commit(true, 1); err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}
	if !backend.Has("gamma") {
		t.Fatalf("expected committed set to update backend")
	}

	store.Delete("alpha")
	store.DeleteBatch([]string{"beta", "gamma"})
	if store.Has("alpha") || store.Has("beta") || store.Has("gamma") {
		t.Fatalf("expected delete operations to evict values from first layer")
	}
	if backend.Len() != 3 {
		t.Fatalf("expected backend to remain unchanged before delete commit, got %d items", backend.Len())
	}

	if err := store.Precommit(); err != nil {
		t.Fatalf("unexpected precommit error: %v", err)
	}
	if err := store.Commit(true, 2); err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}
	if backend.Len() != 0 {
		t.Fatalf("expected backend to be empty after committed delete, got %d items", backend.Len())
	}
}

func TestStoreCurrentLayerOnlySkipsBackend(t *testing.T) {
	localCache := NewCachedKVStore[string, crdtcommon.CRDT](nil, 1024, nil)
	localCache.Set("cached", newProfiledString("value"))
	if value, ok := localCache.Get("cached"); !ok || value == nil {
		t.Fatalf("expected local store read to return cached value")
	}

	backend := NewCachedKVStore[string, crdtcommon.CRDT](nil, 4096, nil)
	backend.Set("alpha", newProfiledString("one"))
	_ = backend.Commit(false, 0)
	store := NewCachedKVStore(backend, 1024, nil)
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
	store := NewCachedKVStore(backend, uint64(entryCount)*sampleSize, nil)
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
	store := NewCachedKVStore(backend, 1024*sampleSize, nil)
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

	commitKey := "committed"
	store.Set(commitKey, newProfiledString("committed"))
	if backend.Has(commitKey) {
		t.Fatalf("expected backend to remain unchanged before commit")
	}

	if err := store.Commit(true, 1); err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}
	if !backend.Has(commitKey) {
		t.Fatalf("expected commit to flush staged write to backend")
	}

	fmt.Printf("Get 1 Million entries, 1024 Entry cache: %v\n", time.Since(t0))
}

func TestStoreSetGet1MillionEntries2048InCache(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newProfiledString("value").Size()

	t0 := time.Now()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	store := NewCachedKVStore(backend, 2048*sampleSize, nil)
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

	commitKey := "committed"
	store.Set(commitKey, newProfiledString("committed"))
	if backend.Has(commitKey) {
		t.Fatalf("expected backend to remain unchanged before commit")
	}

	if err := store.Commit(true, 1); err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}
	if !backend.Has(commitKey) {
		t.Fatalf("expected commit to flush staged write to backend")
	}

	fmt.Printf("Get 1 Million entries, 2048 Entry cache: %v\n", time.Since(t0))
}

func TestStoreSetGet1MillionEntriesAllInCache(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newProfiledString("value").Size()

	t0 := time.Now()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	store := NewCachedKVStore(backend, uint64(entryCount)*sampleSize, nil)
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

	commitKey := "committed"
	store.Set(commitKey, newProfiledString("committed"))
	if backend.Has(commitKey) {
		t.Fatalf("expected backend to remain unchanged before commit")
	}

	if err := store.Commit(true, 1); err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}
	if !backend.Has(commitKey) {
		t.Fatalf("expected commit to flush staged write to backend")
	}

	fmt.Printf("Get 1 Million entries, all in cache: %v\n", time.Since(t0))
}
