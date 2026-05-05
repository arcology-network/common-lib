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
	"fmt"
	"path/filepath"
	"testing"
	"time"

	crdtcommon "github.com/arcology-network/common-lib/crdt/common"
	noncommutative "github.com/arcology-network/common-lib/crdt/noncommutative"
	badgerdb "github.com/arcology-network/common-lib/storage/badger"
	stgcodec "github.com/arcology-network/common-lib/storage/codec"
	filedb "github.com/arcology-network/common-lib/storage/filedb"
	stgintf "github.com/arcology-network/common-lib/storage/interface"
	memdb "github.com/arcology-network/common-lib/storage/memdb"
)

type demoObj interface {
	TypeName() string
	Payload() string
}

type demoObjA struct{ v string }

func (o demoObjA) TypeName() string { return "A" }
func (o demoObjA) Payload() string  { return o.v }

type demoObjB struct{ v string }

func (o demoObjB) TypeName() string { return "B" }
func (o demoObjB) Payload() string  { return o.v }

func equalDemoObj(lhs, rhs demoObj) bool {
	if lhs == nil || rhs == nil {
		return lhs == rhs
	}
	return lhs.TypeName() == rhs.TypeName() && lhs.Payload() == rhs.Payload()
}

type testKVStore[K comparable, T any] struct {
	values           map[K]T
	currentLayerOnly bool
}

type byteBackendAdapter struct {
	get         func(string) (any, error)
	getBatch    func([]string) ([]any, []error)
	set         func(string, []byte) error
	setBatch    func([]string, [][]byte) []error
	delete      func(string) error
	deleteBatch func([]string) []error
	has         func(string) bool
}

func (this *byteBackendAdapter) Get(key string) (any, error) {
	return this.get(key)
}

func (this *byteBackendAdapter) GetBatch(keys []string) ([]any, []error) {
	return this.getBatch(keys)
}

func (this *byteBackendAdapter) Set(key string, value []byte) error {
	return this.set(key, value)
}

func (this *byteBackendAdapter) SetBatch(keys []string, values [][]byte) []error {
	return this.setBatch(keys, values)
}

func (this *byteBackendAdapter) Delete(key string) error {
	return this.delete(key)
}

func (this *byteBackendAdapter) DeleteBatch(keys []string) []error {
	return this.deleteBatch(keys)
}

func (this *byteBackendAdapter) Has(key string) bool {
	return this.has(key)
}

func castAnyBytes(values [][]byte) []any {
	converted := make([]any, len(values))
	for i, value := range values {
		if value != nil {
			converted[i] = value
		}
	}
	return converted
}

func wrapMemoryByteBackend(db *memdb.MemoryDB) stgintf.BackendStore[string, []byte] {
	return &byteBackendAdapter{
		get: func(key string) (any, error) {
			value, err := db.Get(key)
			if err != nil {
				return nil, err
			}
			if value == nil {
				return []byte(nil), nil
			}
			return value.([]byte), nil
		},
		getBatch:    db.GetBatch,
		set:         db.Set,
		setBatch:    db.SetBatch,
		delete:      db.Delete,
		deleteBatch: db.DeleteBatch,
		has: func(key string) bool {
			return db.Has(key)
		},
	}
}

func wrapFileByteBackend(db *filedb.FileDB) stgintf.BackendStore[string, []byte] {
	return &byteBackendAdapter{
		get: func(key string) (any, error) {
			value, err := db.Get(key)
			if err != nil {
				return nil, err
			}
			if value == nil {
				return []byte(nil), nil
			}
			return value.([]byte), nil
		},
		getBatch:    db.GetBatch,
		set:         db.Set,
		setBatch:    db.SetBatch,
		delete:      db.Delete,
		deleteBatch: db.DeleteBatch,
		has: func(key string) bool {
			return db.Has(key)
		},
	}
}

func wrapBadgerByteBackend(db *badgerdb.BadgerDB) stgintf.BackendStore[string, []byte] {
	return &byteBackendAdapter{
		get:         db.Get,
		getBatch:    db.GetBatch,
		set:         db.Set,
		setBatch:    db.SetBatch,
		delete:      db.Delete,
		deleteBatch: db.DeleteBatch,
		has:         db.Has,
	}
}

func wrapParaBadgerByteBackend(db *badgerdb.ParaBadgerDB) stgintf.BackendStore[string, []byte] {
	return &byteBackendAdapter{
		get:         db.Get,
		getBatch:    db.GetBatch,
		set:         db.Set,
		setBatch:    db.SetBatch,
		delete:      db.Delete,
		deleteBatch: db.DeleteBatch,
		has:         db.Has,
	}
}

func newTestKVStore[K comparable, T any]() *testKVStore[K, T] {
	return &testKVStore[K, T]{values: map[K]T{}}
}

func (this *testKVStore[K, T]) Get(key K) (any, error) {
	if value, ok := this.values[key]; ok {
		return value, nil
	}
	return nil, stgintf.ErrNotFound
}

func (this *testKVStore[K, T]) Has(key K) bool {
	_, ok := this.values[key]
	return ok
}

func (this *testKVStore[K, T]) GetBatch(keys []K) ([]any, []error) {
	values := make([]any, len(keys))
	errs := make([]error, len(keys))
	for i, key := range keys {
		if value, ok := this.values[key]; ok {
			values[i] = value
			errs[i] = nil
			continue
		}
		errs[i] = stgintf.ErrNotFound
	}
	return values, errs
}
func (this *testKVStore[K, T]) Set(key K, value T) error {
	this.values[key] = value
	return nil
}

func (this *testKVStore[K, T]) Delete(key K) error {
	delete(this.values, key)
	return nil
}

func (this *testKVStore[K, T]) SetBatch(keys []K, values []T) []error {
	errs := make([]error, len(keys))
	for i := 0; i < len(keys) && i < len(values); i++ {
		this.values[keys[i]] = values[i]
		errs[i] = nil
	}
	return errs
}

func (this *testKVStore[K, T]) DeleteBatch(keys []K) []error {
	errs := make([]error, len(keys))
	for _, key := range keys {
		delete(this.values, key)
	}
	return errs
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

func newIdentityCodec[T any]() *stgcodec.StorageCodec[string, T, string, T] {
	return stgcodec.NewStorageCodec[string, T, string, T](
		func(key string, value T) (string, T, error) {
			return key, value, nil
		},
		func(key string, value T) (string, T, error) {
			return key, value, nil
		},
	)
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
		if err := backend.Set(key, value); err != nil {
			panic(err)
		}
	}
	return backend, keys, values
}

func TestStoreCachesReads(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	if err := backend.Set("alpha", newStringValue("one")); err != nil {
		t.Fatal(err)
	}
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		1024,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)

	first, err := store.Get("alpha")
	if err != nil || first == nil {
		t.Fatalf("expected value from backend on first read")
	}

	cached, err := store.cache.Get("alpha")
	if err != nil || cached == nil {
		t.Fatalf("expected backend read to promote value into first layer")
	}

	second, err := store.Get("alpha")
	if err != nil || second == nil {
		t.Fatalf("expected cached value on second read")
	}
	if second != cached {
		t.Fatalf("expected second read to return the cached value")
	}
}

func TestStoreSkipsCachingOversizedEntry(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	if err := backend.Set("alpha", newStringValue("one")); err != nil {
		t.Fatal(err)
	}

	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		1024,
		func(crdtcommon.CRDT) uint64 { return 2048 },
	)

	first, err := store.Get("alpha")
	if err != nil || first == nil {
		t.Fatalf("expected oversized value from backend on first read")
	}
	if cached, err := store.cache.Get("alpha"); err == nil && cached != nil {
		t.Fatalf("expected oversized entry to stay out of the first layer")
	}

	second, err := store.Get("alpha")
	if err != nil || second == nil {
		t.Fatalf("expected oversized value from backend on second read")
	}
	if cached, err := store.cache.Get("alpha"); err == nil && cached != nil {
		t.Fatalf("expected oversized entry not to be cached after repeated reads")
	}
}

func TestStoreLayeredReadLeavesBackendUntouched(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	if err := backend.Set("alpha", newStringValue("one")); err != nil {
		t.Fatal(err)
	}

	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		1024,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)

	first, err := store.Get("alpha")
	if err != nil || first == nil {
		t.Fatalf("expected first layer to read through backend")
	}

	second, err := store.Get("alpha")
	if err != nil || second == nil {
		t.Fatalf("expected first layer to return cached backend value")
	}
	if second != first {
		t.Fatalf("expected repeated read to return the cached entry")
	}

	local := newStringValue("two")
	store.Set("beta", local)
	if has := store.Has("beta"); !has {
		t.Fatalf("expected local write to stay in first layer")
	}
	if has := backend.Has("beta"); !has {
		t.Fatalf("expected backend to be updated by write-through set")
	}
	if value, err := store.Get("beta"); err != nil || value != local {
		t.Fatalf("expected local value to remain in cache")
	}
}

func TestStoreTracksVisitsGenerically(t *testing.T) {
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		nil,
		codec,
		4096,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)

	store.Set("alpha", newStringValue("one"))
	alpha, err := store.cache.Get("alpha")
	if err != nil || alpha == nil {
		t.Fatalf("expected set to populate cache")
	}

	value, err := store.Get("alpha")
	if err != nil || value == nil {
		t.Fatalf("expected get to return cached value")
	}

	replacementValue := newStringValue("two")
	store.Set("alpha", replacementValue)
	replacement, err := store.cache.Get("alpha")
	if err != nil || replacement == nil {
		t.Fatalf("expected replacement to stay cached")
	}

	values, errs := store.GetBatch([]string{"alpha"})
	if len(values) != 1 || errs[0] != nil || values[0] == nil {
		t.Fatalf("expected batch get to return cached value")
	}

	store.SetBatch([]string{"beta"}, []crdtcommon.CRDT{newStringValue("three")})
	beta, err := store.cache.Get("beta")
	if err != nil || beta == nil {
		t.Fatalf("expected batch set to update cache entry")
	}

	backend := newTestKVStore[string, crdtcommon.CRDT]()
	if err := backend.Set("backend", newStringValue("backend")); err != nil {
		t.Fatal(err)
	}

	layered := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		4096,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)
	fetched, err := layered.Get("backend")
	if err != nil || fetched == nil {
		t.Fatalf("expected backend get to succeed")
	}
	fetchedEntry, err := layered.cache.Get("backend")
	if err != nil || fetchedEntry == nil {
		t.Fatalf("expected backend get to populate cache")
	}

	again, err := layered.Get("backend")
	if err != nil || again != fetched {
		t.Fatalf("expected backend value to be cached after first read")
	}
}

func TestStoreEvictsWhenFull(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		1024,
		func(crdtcommon.CRDT) uint64 { return 700 },
	)

	alpha := newStringValue("one")
	beta := newStringValue("two")

	store.Set("alpha", alpha)
	store.Set("beta", beta)

	if store.cache.Length() == 0 {
		t.Fatalf("expected at least one entry to stay in the first layer")
	}

	store.cache.Evict()

	if len(backend.values) != 2 {
		t.Fatalf("expected backend to have write-through entries, got %d", len(backend.values))
	}
	if store.cache.Cap() > 1024 {
		t.Fatalf("expected eviction to keep cache within cap, got %d", store.cache.Cap())
	}
	if store.cache.Length() >= 2 {
		t.Fatalf("expected eviction to remove at least one cached entry")
	}
}

func TestStoreBatchAndDelete(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	if err := backend.Set("alpha", newStringValue("one")); err != nil {
		t.Fatal(err)
	}
	if err := backend.Set("beta", newStringValue("two")); err != nil {
		t.Fatal(err)
	}
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		1024,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)

	values, errs := store.GetBatch([]string{"alpha", "beta", "missing"})
	if len(values) != 3 || errs[0] != nil || errs[1] != nil || errs[2] == nil {
		t.Fatalf("unexpected batch result: %#v, errs: %#v", values, errs)
	}
	if cached, err := store.cache.Get("alpha"); err != nil || cached == nil {
		t.Fatalf("expected batch backend read to populate cache for alpha")
	}
	if cached, err := store.cache.Get("beta"); err != nil || cached == nil {
		t.Fatalf("expected batch backend read to populate cache for beta")
	}

	gamma := newStringValue("three")
	store.Set("gamma", gamma)
	if has := backend.Has("gamma"); !has {
		t.Fatalf("expected set to write through to backend")
	}
	if has := store.Has("gamma"); !has {
		t.Fatalf("expected set to write through to backend")
	}
	store.Delete("alpha")
	store.DeleteBatch([]string{"beta", "gamma"})
	if has := store.Has("gamma"); has {
		t.Fatalf("expected key to stay absent after backend delete")
	}

	if has := store.Has("alpha"); has {
		t.Fatalf("expected deleted backend-backed keys to remain absent")
	}
	if has := store.Has("beta"); has {
		t.Fatalf("expected deleted backend-backed keys to remain absent")
	}
	if _, err := store.cache.Get("alpha"); err == nil {
		t.Fatalf("expected delete to evict alpha from first layer")
	}
	if _, err := store.cache.Get("beta"); err == nil {
		t.Fatalf("expected delete to evict beta from first layer")
	}
	if _, err := store.cache.Get("gamma"); err == nil {
		t.Fatalf("expected delete to evict gamma from first layer")
	}
	if len(backend.values) != 0 {
		t.Fatalf("expected backend deletes to remove all keys, got %d items", len(backend.values))
	}

	if value, err := store.Get("alpha"); err == nil || value != nil {
		t.Fatalf("expected deleted alpha to be absent")
	}
	if value, err := store.Get("beta"); err == nil || value != nil {
		t.Fatalf("expected deleted beta to be absent")
	}
	if value, err := store.Get("gamma"); err == nil && value != nil {
		t.Fatalf("expected deleted gamma to stay absent")
	}
}

func TestStoreWithoutBackendUsesCacheOnly(t *testing.T) {
	codec := newIdentityCodec[crdtcommon.CRDT]()
	localCache := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		nil,
		codec,
		1024,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)
	localCache.Set("cached", newStringValue("value"))
	if value, err := localCache.Get("cached"); err != nil || value == nil {
		t.Fatalf("expected local store read to return cached value")
	}
	if value, err := localCache.Get("missing"); err == nil || value != nil {
		t.Fatalf("expected missing key to stay absent without a backend")
	}
	if has := localCache.Has("missing"); has {
		t.Fatalf("expected missing key to stay absent without a backend")
	}
	values, errs := localCache.cache.GetBatch([]string{"cached", "missing"})
	if len(values) != 2 || errs[0] != nil || values[0] == nil || errs[1] == nil || values[1] != nil {
		t.Fatalf("expected cache-only batch read to use local entries only")
	}
}

func TestStoreSetGet1MillionEntries(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newStringValue("value").MemSize()

	backend, keys, _ := newBenchmarkBackend(entryCount)
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		uint64(entryCount)*sampleSize,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)

	t0 := time.Now()
	for _, key := range keys {
		if _, err := store.Get(key); err != nil {
			t.Fatalf("expected warm cache get to succeed for %s", key)
		}
	}
	fmt.Printf("Set 1 Million entries: %v\n", time.Since(t0))

	t0 = time.Now()
	for i := 0; i < entryCount; i++ {
		key := keys[i%entryCount]
		value, err := store.Get(key)
		if err != nil || value == nil {
			t.Fatalf("expected layered get to succeed for %s", key)
		}
	}
	fmt.Printf("Get 1 Million entries: %v\n", time.Since(t0))
}

func TestStoreSetGet1MillionEntries1024InCache(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newStringValue("value").MemSize()

	t0 := time.Now()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		1024*sampleSize,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)
	for _, key := range keys {
		if _, err := store.Get(key); err != nil {
			t.Fatalf("expected warm cache get to succeed for %s", key)
		}
	}
	for i := 0; i < entryCount; i++ {
		key := keys[i%entryCount]
		value, err := store.Get(key)
		if err != nil || value == nil {
			t.Fatalf("expected layered get to succeed for %s", key)
		}
	}

	store.Set("committed", newStringValue("committed"))
	if has := backend.Has("committed"); !has {
		t.Fatalf("expected backend to be updated after write-through set")
	}

	fmt.Printf("Get 1 Million entries, 1024 Entry cache: %v\n", time.Since(t0))
}

func TestStoreSetGet1MillionEntries2048InCache(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newStringValue("value").MemSize()

	t0 := time.Now()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		2048*sampleSize,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)
	for _, key := range keys {
		if _, err := store.Get(key); err != nil {
			t.Fatalf("expected warm cache get to succeed for %s", key)
		}
	}
	for i := 0; i < entryCount; i++ {
		key := keys[i%entryCount]
		value, err := store.Get(key)
		if err != nil || value == nil {
			t.Fatalf("expected layered get to succeed for %s", key)
		}
	}

	store.Set("committed", newStringValue("committed"))
	if has := backend.Has("committed"); !has {
		t.Fatalf("expected backend to be updated after write-through set")
	}

	fmt.Printf("Get 1 Million entries, 2048 Entry cache: %v\n", time.Since(t0))
}

func TestStoreSetGet1MillionEntriesAllInCache(t *testing.T) {
	const entryCount = 1_000_000
	sampleSize := newStringValue("value").MemSize()

	t0 := time.Now()
	backend, keys, _ := newBenchmarkBackend(entryCount)
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		uint64(entryCount)*sampleSize,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)
	for _, key := range keys {
		if _, err := store.Get(key); err != nil {
			t.Fatalf("expected warm cache get to succeed for %s", key)
		}
	}
	for i := 0; i < entryCount; i++ {
		key := keys[i%entryCount]
		value, err := store.Get(key)
		if err != nil || value == nil {
			t.Fatalf("expected first-layer hit for %s", key)
		}
	}

	store.Set("committed", newStringValue("committed"))
	if has := backend.Has("committed"); !has {
		t.Fatalf("expected backend to be updated after write-through set")
	}

	fmt.Printf("Get 1 Million entries, all in cache: %v\n", time.Since(t0))
}

func TestGetBatch_HalfCacheHalfBackend(t *testing.T) {
	backend := newTestKVStore[string, crdtcommon.CRDT]()
	codec := newIdentityCodec[crdtcommon.CRDT]()
	store := NewCachedStore[string, crdtcommon.CRDT, string, crdtcommon.CRDT](
		backend,
		codec,
		1024,
		func(crdtcommon.CRDT) uint64 { return 1 },
	)

	// Prepare keys and values
	keys := []string{"a", "b", "c", "d"}
	vals := []crdtcommon.CRDT{
		newStringValue("A"),
		newStringValue("B"),
		newStringValue("C"),
		newStringValue("D"),
	}

	// Put first half in cache, second half in backend only
	store.Set(keys[0], vals[0])
	store.Set(keys[1], vals[1])
	if err := backend.Set(keys[2], vals[2]); err != nil {
		t.Fatal(err)
	}
	if err := backend.Set(keys[3], vals[3]); err != nil {
		t.Fatal(err)
	}

	// Query all
	got, errs := store.GetBatch(keys)
	for i := range keys {
		if errs[i] != nil || got[i] == nil {
			t.Fatalf("missing value for key %s", keys[i])
		}
		if !got[i].(crdtcommon.CRDT).Equal(vals[i]) {
			t.Fatalf("wrong value for key %s: got %v, want %v", keys[i], got[i], vals[i])
		}
	}

	// Ensure all are now cached
	for _, k := range keys {
		entry, err := store.cache.Get(k)
		if (k == "a" || k == "b") && (err != nil || entry == nil) {
			t.Fatalf("key %s should remain cached after GetBatch", k)
		}
		if (k == "c" || k == "d") && (err != nil || entry == nil) {
			t.Fatalf("key %s should be cached after GetBatch", k)
		}
	}
}

func TestSetBatch_CoversAllBranches(t *testing.T) {
	identityCodec := newIdentityCodec[string]()

	// Cover early return when backend is nil.
	nilBackendStore := NewCachedStore[string, string, string, string](
		nil,
		identityCodec,
		1024,
		func(string) uint64 { return 1 },
	)
	nilBackendStore.SetBatch([]string{"local"}, []string{"value"})
	if v, err := nilBackendStore.cache.Get("local"); err != nil || v.(string) != "value" {
		t.Fatalf("expected nil-backend SetBatch to update local cache")
	}

	backend := newTestKVStore[string, string]()
	codec := stgcodec.NewStorageCodec[string, string, string, string](
		func(k string, v string) (string, string, error) {
			if k == "bad-key" {
				return "", "", fmt.Errorf("key conversion failed")
			}
			if v == "" {
				return "bk:" + k, "", nil
			}
			if v == "bad-val" {
				return "", "", fmt.Errorf("value conversion failed")
			}
			return "bk:" + k, "bv:" + v, nil
		},
		func(k string, v string) (string, string, error) {
			return k, v, nil
		},
	)

	store := NewCachedStore[string, string, string, string](
		backend,
		codec,
		1024,
		func(string) uint64 { return 1 },
	)

	keys := []string{"ok", "bad-key", "bad-val"}
	values := []string{"good", "skip-key", "bad-val"}
	store.SetBatch(keys, values)

	if got, err := backend.Get("bk:ok"); err != nil || got.(string) != "bv:good" {
		t.Fatalf("expected successful key/value conversion to reach backend")
	}
	if _, err := backend.Get("bk:bad-key"); err == nil {
		t.Fatalf("did not expect backend entry for key conversion failure")
	}
	if _, err := backend.Get("bk:bad-val"); err == nil {
		t.Fatalf("did not expect backend entry for value conversion failure")
	}
}

func TestPartialCache_ByteBackend_WithEmbeddedTypeInfo(t *testing.T) {
	backend := newTestKVStore[string, []byte]()

	codec := stgcodec.NewStorageCodec[string, demoObj, string, []byte](
		func(k string, v demoObj) (string, []byte, error) {
			if v == nil {
				return k, nil, nil
			}
			return k, append([]byte{v.TypeName()[0]}, []byte(v.Payload())...), nil
		},
		func(k string, raw []byte) (string, demoObj, error) {
			if len(raw) == 0 {
				return k, nil, fmt.Errorf("empty encoded payload")
			}
			payload := string(raw[1:])
			switch raw[0] {
			case 'A':
				return k, demoObjA{v: payload}, nil
			case 'B':
				return k, demoObjB{v: payload}, nil
			default:
				return k, nil, fmt.Errorf("unknown embedded type %q", raw[0])
			}
		},
	)

	store := NewCachedStore[string, demoObj, string, []byte](
		backend,
		codec,
		1024,
		func(demoObj) uint64 { return 1 },
	)

	// Partial cache: alpha is cached; beta exists only in backend as bytes.
	alpha := demoObjA{v: "alpha"}
	store.Set("alpha", alpha)
	if err := backend.Set("beta", append([]byte{'B'}, []byte("beta")...)); err != nil {
		t.Fatal(err)
	}

	if got, err := store.Get("alpha"); err != nil || !equalDemoObj(got.(demoObj), alpha) {
		t.Fatalf("expected cache hit for alpha")
	}

	betaExpected := demoObjB{v: "beta"}
	got, err := store.Get("beta")
	if err != nil || !equalDemoObj(got.(demoObj), betaExpected) {
		t.Fatalf("expected backend decode with embedded type info for beta")
	}

	if cached, err := store.cache.Get("beta"); err != nil || !equalDemoObj(cached.(demoObj), betaExpected) {
		t.Fatalf("expected backend read to populate cache for beta")
	}
}

func TestPartialCache_ByteBackend_WithoutTypeInfo_WithExternalHint(t *testing.T) {
	backend := newTestKVStore[string, []byte]()
	decodeHints := map[string]string{}

	codec := stgcodec.NewStorageCodec[string, demoObj, string, []byte](
		func(k string, v demoObj) (string, []byte, error) {
			if v == nil {
				return k, nil, nil
			}
			return k, []byte(v.Payload()), nil
		},
		func(k string, raw []byte) (string, demoObj, error) {
			hint, ok := decodeHints[k]
			if !ok {
				return k, nil, fmt.Errorf("missing decode hint for key %q", k)
			}
			switch hint {
			case "A":
				return k, demoObjA{v: string(raw)}, nil
			case "B":
				return k, demoObjB{v: string(raw)}, nil
			default:
				return k, nil, fmt.Errorf("unknown decode hint %q", hint)
			}
		},
	)

	store := NewCachedStore[string, demoObj, string, []byte](
		backend,
		codec,
		1024,
		func(demoObj) uint64 { return 1 },
	)

	// Partial cache: alpha is cached; beta exists only in backend as raw bytes without type tag.
	alpha := demoObjA{v: "alpha"}
	store.Set("alpha", alpha)
	if err := backend.Set("beta", []byte("beta")); err != nil {
		t.Fatal(err)
	}
	decodeHints["beta"] = "B"

	if got, err := store.Get("alpha"); err != nil || !equalDemoObj(got.(demoObj), alpha) {
		t.Fatalf("expected cache hit for alpha")
	}

	betaExpected := demoObjB{v: "beta"}
	got, err := store.Get("beta")
	if err != nil || !equalDemoObj(got.(demoObj), betaExpected) {
		t.Fatalf("expected backend decode with external type hint for beta")
	}

	if cached, err := store.cache.Get("beta"); err != nil || !equalDemoObj(cached.(demoObj), betaExpected) {
		t.Fatalf("expected backend read to populate cache for beta")
	}
}

func TestStoreWithMemoryDBBackend(t *testing.T) {
	rawDB := memdb.NewMemoryDB()
	backend := rawDB

	// Encode: prefix byte = TypeName()[0], followed by Payload bytes.
	// Decode: first byte identifies the type, remaining bytes are the payload.
	codec := stgcodec.NewStorageCodec[string, demoObj, string, []byte](
		func(k string, v demoObj) (string, []byte, error) {
			if v == nil {
				return k, nil, nil
			}
			return k, append([]byte{v.TypeName()[0]}, []byte(v.Payload())...), nil
		},
		func(k string, raw []byte) (string, demoObj, error) {
			if len(raw) == 0 {
				return k, nil, fmt.Errorf("empty encoded payload")
			}
			payload := string(raw[1:])
			switch raw[0] {
			case 'A':
				return k, demoObjA{v: payload}, nil
			case 'B':
				return k, demoObjB{v: payload}, nil
			default:
				return k, nil, fmt.Errorf("unknown type byte 0x%x", raw[0])
			}
		},
	)

	store := NewCachedStore[string, demoObj, string, []byte](
		wrapMemoryByteBackend(backend),
		codec,
		1024,
		func(demoObj) uint64 { return 1 },
	)

	// Write through: Set should store in cache AND persist to MemoryDB.
	alpha := demoObjA{v: "hello"}
	store.Set("alpha", alpha)

	raw, err := rawDB.Get("alpha")
	if err != nil || len(raw.([]byte)) == 0 {
		t.Fatalf("expected Set to write through to MemoryDB backend")
	}

	// Cache hit: second Get should return the cached value.
	got, err := store.Get("alpha")
	if err != nil || !equalDemoObj(got.(demoObj), alpha) {
		t.Fatalf("expected cache hit for alpha, got %v", got)
	}

	// Backend read: insert directly into MemoryDB and confirm cache miss triggers
	// backend lookup which then populates the cache.
	betaEncoded := append([]byte{'B'}, []byte("world")...)
	if err := rawDB.Set("beta", betaEncoded); err != nil {
		t.Fatal(err)
	}

	betaExpected := demoObjB{v: "world"}
	got, err = store.Get("beta")
	if err != nil || !equalDemoObj(got.(demoObj), betaExpected) {
		t.Fatalf("expected backend read via MemoryDB for beta, got %v", got)
	}
	if cached, err := store.cache.Get("beta"); err != nil || !equalDemoObj(cached.(demoObj), betaExpected) {
		t.Fatalf("expected MemoryDB backend read to populate cache for beta")
	}

	// Batch write: SetBatch should reach MemoryDB.
	store.SetBatch([]string{"g1", "g2"}, []demoObj{demoObjA{v: "v1"}, demoObjB{v: "v2"}})
	hg1 := rawDB.Has("g1")
	hg2 := rawDB.Has("g2")
	if !hg1 || !hg2 {
		t.Fatalf("expected SetBatch to write through to MemoryDB")
	}

	// Batch read: GetBatch with mix of cache and backend entries.
	if err := rawDB.Set("g3", append([]byte{'A'}, []byte("v3")...)); err != nil {
		t.Fatal(err)
	}
	vals, errs := store.GetBatch([]string{"g1", "g3", "missing"})
	if errs[0] != nil || errs[1] != nil || errs[2] == nil {
		t.Fatalf("unexpected GetBatch errs: %v", errs)
	}
	if !equalDemoObj(vals[0].(demoObj), demoObjA{v: "v1"}) {
		t.Fatalf("expected g1 = A:v1, got %v", vals[0])
	}
	if !equalDemoObj(vals[1].(demoObj), demoObjA{v: "v3"}) {
		t.Fatalf("expected g3 = A:v3, got %v", vals[1])
	}

	// Delete: should remove from cache and MemoryDB.
	store.Delete("alpha")
	if has := store.Has("alpha"); has {
		t.Fatalf("expected Delete to remove alpha from store")
	}
	if has := rawDB.Has("alpha"); has {
		t.Fatalf("expected Delete to remove alpha from MemoryDB")
	}

	// DeleteBatch: should remove multiple entries.
	store.DeleteBatch([]string{"g1", "g2"})
	if has := store.Has("g1"); has {
		t.Fatalf("expected DeleteBatch to remove g1 and g2")
	}
	if has := store.Has("g2"); has {
		t.Fatalf("expected DeleteBatch to remove g1 and g2")
	}
	h1 := rawDB.Has("g1")
	h2 := rawDB.Has("g2")
	if h1 || h2 {
		t.Fatalf("expected DeleteBatch to remove g1 and g2 from MemoryDB")
	}
}

func newDemoObjCodec() *stgcodec.StorageCodec[string, demoObj, string, []byte] {
	return stgcodec.NewStorageCodec[string, demoObj, string, []byte](
		func(k string, v demoObj) (string, []byte, error) {
			if v == nil {
				return k, nil, nil
			}
			return k, append([]byte{v.TypeName()[0]}, []byte(v.Payload())...), nil
		},
		func(k string, raw []byte) (string, demoObj, error) {
			if len(raw) == 0 {
				return k, nil, fmt.Errorf("empty encoded payload")
			}
			payload := string(raw[1:])
			switch raw[0] {
			case 'A':
				return k, demoObjA{v: payload}, nil
			case 'B':
				return k, demoObjB{v: payload}, nil
			default:
				return k, nil, fmt.Errorf("unknown type byte 0x%x", raw[0])
			}
		},
	)
}

func runStoreWithByteBackend(t *testing.T, backend stgintf.BackendStore[string, []byte], backendName string) {
	t.Helper()

	codec := newDemoObjCodec()
	store := NewCachedStore[string, demoObj, string, []byte](
		backend,
		codec,
		1024,
		func(demoObj) uint64 { return 1 },
	)

	alpha := demoObjA{v: "hello"}
	store.Set("alpha", alpha)
	raw, err := backend.Get("alpha")
	if err != nil || len(raw.([]byte)) == 0 {
		t.Fatalf("expected Set to write through to %s backend", backendName)
	}

	got, err := store.Get("alpha")
	if err != nil || !equalDemoObj(got.(demoObj), alpha) {
		t.Fatalf("expected cache hit for alpha, got %v", got)
	}

	betaEncoded := append([]byte{'B'}, []byte("world")...)
	if err := backend.Set("beta", betaEncoded); err != nil {
		t.Fatalf("failed to seed backend value for beta: %v", err)
	}

	betaExpected := demoObjB{v: "world"}
	got, err = store.Get("beta")
	if err != nil || !equalDemoObj(got.(demoObj), betaExpected) {
		t.Fatalf("expected backend read via %s for beta, got %v", backendName, got)
	}
	if cached, err := store.cache.Get("beta"); err != nil || !equalDemoObj(cached.(demoObj), betaExpected) {
		t.Fatalf("expected %s backend read to populate cache for beta", backendName)
	}

	store.SetBatch([]string{"key-g1", "key-g2"}, []demoObj{demoObjA{v: "v1"}, demoObjB{v: "v2"}})
	hg1 := backend.Has("key-g1")
	hg2 := backend.Has("key-g2")
	if !hg1 || !hg2 {
		t.Fatalf("expected SetBatch to write through to %s", backendName)
	}

	if err := backend.Set("key-g3", append([]byte{'A'}, []byte("v3")...)); err != nil {
		t.Fatalf("failed to seed backend value for g3: %v", err)
	}
	vals, errs := store.GetBatch([]string{"key-g1", "key-g3", "missing"})
	if errs[0] != nil || errs[1] != nil || errs[2] == nil {
		t.Fatalf("unexpected GetBatch errs: %v", errs)
	}
	if !equalDemoObj(vals[0].(demoObj), demoObjA{v: "v1"}) {
		t.Fatalf("expected g1 = A:v1, got %v", vals[0])
	}
	if !equalDemoObj(vals[1].(demoObj), demoObjA{v: "v3"}) {
		t.Fatalf("expected g3 = A:v3, got %v", vals[1])
	}

	store.Delete("alpha")
	if has := store.Has("alpha"); has {
		t.Fatalf("expected Delete to remove alpha from store")
	}
	if has := backend.Has("alpha"); has {
		t.Fatalf("expected Delete to remove alpha from %s", backendName)
	}

	store.DeleteBatch([]string{"key-g1", "key-g2"})
	if has := store.Has("key-g1"); has {
		t.Fatalf("expected DeleteBatch to remove g1")
	}
	if has := store.Has("key-g2"); has {
		t.Fatalf("expected DeleteBatch to remove g2")
	}
	h1 := backend.Has("key-g1")
	h2 := backend.Has("key-g2")
	if h1 || h2 {
		t.Fatalf("expected DeleteBatch to remove g1 and g2 from %s", backendName)
	}
}

func TestStoreWithFileDBBackend(t *testing.T) {
	db, err := filedb.NewFileDB(filepath.Join(t.TempDir(), "filedb"), 8, 2)
	if err != nil {
		t.Fatalf("failed to create filedb backend: %v", err)
	}
	runStoreWithByteBackend(t, wrapFileByteBackend(db), "FileDB")
}

func TestStoreWithBadgerDBBackend(t *testing.T) {
	db := badgerdb.NewBadgerDB(filepath.Join(t.TempDir(), "badger"))
	t.Cleanup(func() {
		_ = db.Close()
	})
	runStoreWithByteBackend(t, wrapBadgerByteBackend(db), "BadgerDB")
}

func TestStoreWithParaBadgerDBBackend(t *testing.T) {
	db := badgerdb.NewParaBadgerDB(filepath.Join(t.TempDir(), "parabadger")+"/", nil)
	t.Cleanup(func() {
		_ = db.Close()
	})
	runStoreWithByteBackend(t, wrapParaBadgerByteBackend(db), "ParaBadgerDB")
}
