package cache

import (
	"testing"

	stgintf "github.com/arcology-network/common-lib/storage/interface"
)

// For entrySize test: type with MemSize method
type sized struct{}

func (s sized) MemSize() uint64 { return 99 }

func testStringHash(v string) uint64 {
	if len(v) == 0 {
		return 0
	}
	return uint64(v[0])
}

func TestCacheBasic(t *testing.T) {
	// Policy: count size as value itself
	policy := NewCachePolicy(100, func(v int) uint64 { return uint64(v) })
	cache := NewCache(4, testStringHash, policy)

	// Hash
	if cache.Hash("abc") != testStringHash("abc")%4 {
		t.Error("Hash should match testStringHash")
	}

	// Set/Get basic

	cache.Set("alpha", 11)
	value, err := cache.Get("alpha")
	if err != nil || value.(int) != 11 {
		t.Error("Get should return the cached value when present")
	}

	// Overwrite
	cache.Set("alpha", 22)
	value, err = cache.Get("alpha")
	if err != nil || value.(int) != 22 {
		t.Error("Get should return updated value")
	}

	// SetBatch/GetBatch
	cache.SetBatch([]string{"beta", "gamma"}, []int{22, 33})
	vals, errs := cache.GetBatch([]string{"beta", "gamma", "delta"})
	if vals[0].(int) != 22 || errs[0] != nil {
		t.Error("GetBatch should return correct value for beta")
	}
	if vals[1].(int) != 33 || errs[1] != nil {
		t.Error("GetBatch should return correct value for gamma")
	}
	if errs[2] == nil {
		t.Error("GetBatch should return error for missing key delta")
	}

	// Delete should remove key
	cache.Set("toRemove", 123)
	cache.Delete("toRemove")
	_, err = cache.Get("toRemove")
	if err == nil {
		t.Error("Deleted key should not be found")
	}

	// Delete/DeleteBatch
	cache.Delete("alpha")
	_, err = cache.Get("alpha")
	if err == nil {
		t.Error("Deleted key should not be found")
	}
	cache.DeleteBatch([]string{"beta", "gamma"})
	_, err = cache.Get("beta")
	if err == nil {
		t.Error("Deleted key beta should not be found")
	}
	_, err = cache.Get("gamma")
	if err == nil {
		t.Error("Deleted key gamma should not be found")
	}

	// Cap/Len/Clear
	cache.Set("a", 1)
	cache.Set("b", 2)
	if cache.Length() != 2 {
		t.Error("Len should be 2 after two inserts")
	}
	if cache.Cap() == 0 {
		t.Error("Cap should reflect policy size")
	}
	cache.Clear()
	if cache.Length() != 0 {
		t.Error("Clear should remove all entries")
	}

	// Policy: test eviction
	cache.Set("x", 90)
	cache.Set("y", 20) // triggers eviction if maxSize=100
	if cache.Policy().NeedEviction() {
		cache.Evict()
		if cache.Policy().NeedEviction() {
			t.Error("Evict should clear enough entries")
		}
	}

	// Policy: Remove (simulate by Update)
	cache.Set("z", 10)
	cache.Policy().Update(10, 0)
	if cache.Policy().Size() > 100 {
		t.Error("Update should decrease occupied size")
	}

	// wrap/value sizing
	entry := cache.wrap(123)
	if entry.value != 123 || entry.visits != 1 {
		t.Error("wrap should set value and visits")
	}
	if cache.Policy().ValueSize(entry.value) != 123 {
		t.Error("ValueSize should use the cache policy sizer")
	}

	// entry.Size with MemSize
	sizedCache := NewCache[string, sized](4, testStringHash, nil)
	entry2 := sizedCache.wrap(sized{})
	if entry2.Size() != 99 {
		t.Error("entry Size should use MemSize if available")
	}
}

func TestCacheEvictNoOpWhenNotNeeded(t *testing.T) {
	policy := NewCachePolicy[int](100, func(v int) uint64 { return uint64(v) })
	c := NewCache[string, int](4, testStringHash, policy)

	c.Set("a", 10)
	beforeLen := c.Length()
	beforeSize := c.Policy().Size()

	c.Evict()

	if c.Length() != beforeLen {
		t.Fatalf("expected Evict to be a no-op when eviction is not needed")
	}
	if c.Policy().Size() != beforeSize {
		t.Fatalf("expected occupied size to remain unchanged when eviction is not needed")
	}
}

func TestCacheEvictDeletesAndStopsWhenRecovered(t *testing.T) {
	policy := NewCachePolicy[int](10, func(v int) uint64 { return uint64(v) })
	c := NewCache[string, int](4, testStringHash, policy)

	// Keep two real entries so Evict deletes one and then returns when pressure is relieved.
	c.Set("a", 6)
	c.Set("b", 3)

	// Force eviction pressure for branch coverage. Normal Set does not permit overflow by design.
	c.Policy().occupied = 11

	c.Evict()

	if c.Policy().NeedEviction() {
		t.Fatalf("expected Evict to reduce usage below policy max")
	}
	if c.Length() >= 2 {
		t.Fatalf("expected Evict to remove at least one entry")
	}
}

func TestCacheEvictSkipsNilEntries(t *testing.T) {
	policy := NewCachePolicy[int](1, func(v int) uint64 { return uint64(v) })
	c := NewCache[string, int](4, testStringHash, policy)

	// Inject a raw nil entry pointer directly into a shard to exercise the nil-skip path.
	shardID := c.Hash("nil")
	c.ConcurrentMap.Shards()[shardID]["nil"] = nil
	c.Policy().occupied = 2

	c.Evict()

	if _, ok := c.ConcurrentMap.Get("nil"); !ok {
		t.Fatalf("expected nil entry key to remain after Evict nil-skip path")
	}
}

func TestCacheDisabledBypassesAllOperations(t *testing.T) {
	policy := NewCachePolicy[int](100, func(v int) uint64 { return uint64(v) })
	c := NewCache[string, int](4, testStringHash, policy)

	c.Set("before", 7)
	c.SetStatus(false)

	if c.Status() {
		t.Fatalf("expected cache status to report disabled")
	}
	if c.Has("before") {
		t.Fatalf("expected disabled cache to hide existing entries")
	}
	if _, err := c.Get("before"); err == nil {
		t.Fatalf("expected disabled cache Get to miss")
	}

	values, errs := c.GetBatch([]string{"before", "missing"})
	if values != nil {
		t.Fatalf("expected disabled cache GetBatch values to be nil")
	}
	if len(errs) != 2 || errs[0] != stgintf.ErrNotFound || errs[1] != stgintf.ErrNotFound {
		t.Fatalf("expected disabled cache GetBatch to miss every key")
	}

	if err := c.Set("new", 9); err != nil {
		t.Fatalf("expected disabled cache Set to stay a no-op, got %v", err)
	}
	if c.Length() != 1 {
		t.Fatalf("expected disabled cache Set to avoid mutating entries")
	}

	if gotKeys, gotValues, errs := c.Query("before", nil); gotKeys != nil || gotValues != nil || errs != nil {
		t.Fatalf("expected disabled cache Query to return empty results")
	}

	c.Policy().occupied = 200
	c.Evict()
	if c.Length() != 1 {
		t.Fatalf("expected disabled cache Evict to stay a no-op")
	}
	if c.Policy().Size() != 200 {
		t.Fatalf("expected disabled cache Evict to leave policy usage unchanged")
	}

	if err := c.Delete("before"); err != nil {
		t.Fatalf("expected disabled cache Delete to stay a no-op, got %v", err)
	}
	if c.Length() != 1 {
		t.Fatalf("expected disabled cache Delete to avoid mutating entries")
	}
}
