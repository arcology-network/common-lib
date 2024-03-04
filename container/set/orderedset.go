// An ordered set is a collection of unique elements where the order of insertion is preserved. In other words,
// it combines the properties of a set (no duplicate elements) and a list (preserves the order of insertion).

package orderedset

import (
	"math"

	slice "github.com/arcology-network/common-lib/exp/slice"
	"github.com/elliotchance/orderedmap"
)

// OrderedSet represents an ordered set data structure.
type OrderedSet struct {
	dict    *orderedmap.OrderedMap // committed keys + added - removed
	keys    []string               // slice of keys in the order of insertion
	touched bool                   // indicates if the set has been modified
}

// NewOrderedSet creates a new instance of OrderedSet with the given initial keys.
// Duplicate keys are not allowed and will be ignored.
func NewOrderedSet(keys []string) *OrderedSet {
	this := &OrderedSet{
		dict:    orderedmap.NewOrderedMap(),
		touched: true,
		keys:    []string{},
	}

	for i := 0; i < len(keys); i++ {
		if this.Insert(keys[i]) {
			keys = append(keys, keys[i])
		}
	}
	return this
}

// Equal checks if the current OrderedSet is equal to another OrderedSet.
// Two OrderedSets are considered equal if they have the same keys in the same order and the same touched and synced status.
func (this *OrderedSet) Equal(other *OrderedSet) bool {
	if this == other && this == nil {
		return true
	}

	if this == nil || other == nil {
		return false
	}

	return slice.Equal(this.keys, other.keys) &&
		this.touched == other.touched &&
		this.isSynced() == other.isSynced()
}

// Length returns the number of elements in the OrderedSet.
func (this *OrderedSet) Length() int {
	return len(this.keys)
}

// isSynced checks if the OrderedSet is in sync with the underlying ordered map.
func (this *OrderedSet) isSynced() bool {
	return (this.dict.Len()) == len(this.keys)
}

// Dict returns the underlying ordered map of the OrderedSet.
// If the ordered map is not in sync with the keys, it will be synced before returning.
func (this *OrderedSet) Dict() *orderedmap.OrderedMap {
	if !this.isSynced() {
		panic("OrderedSet is not in sync with the underlying ordered map. This should never happen.")
		if this.dict.Len() > 0 {
			this.dict = orderedmap.NewOrderedMap() // This should never happen
		}

		for i := 0; i < len(this.keys); i++ {
			this.dict.Set(this.keys[i], uint64(i))
		}
	}
	return this.dict
}

// Touched returns the touched status of the OrderedSet.
func (this *OrderedSet) Touched() bool { return this.touched }

// Len returns the number of elements in the OrderedSet as a uint64.
func (this *OrderedSet) Len() uint64 { return uint64(len(this.keys)) }

// Keys returns a slice of keys in the OrderedSet.
func (this *OrderedSet) Keys() []string {
	return this.keys
}

// Clone creates a deep copy of the OrderedSet.
func (this *OrderedSet) Clone() interface{} {
	if this == nil {
		return this
	}
	set := NewOrderedSet(this.keys)
	set.touched = this.touched
	return set
}

// Exists checks if a key exists in the OrderedSet.
func (this *OrderedSet) Exists(key string) bool {
	_, ok := this.Dict().Get(key)
	return ok
}

// Delete removes a key from the OrderedSet.
// It returns true if the key was successfully deleted, false otherwise.
func (this *OrderedSet) Delete(key string) bool { return this.DeleteByKey(key) }

// IdxOf returns the index of a key in the OrderedSet.
// If the key does not exist, it returns math.MaxUint64.
func (this *OrderedSet) IdxOf(key string) uint64 {
	v, ok := this.Dict().Get(key)
	if !ok {
		return math.MaxUint64
	}
	return v.(uint64)
}

// KeyAt returns the key at the given index in the OrderedSet.
// If the index is out of range, it returns an empty string.
func (this *OrderedSet) KeyAt(idx uint64) string {
	if idx < uint64(len(this.keys)) {
		return this.keys[idx]
	}
	return ""
}

// Insert adds a new key to the OrderedSet.
// If the key already exists, it is not added again.
func (this *OrderedSet) Insert(key string) bool {
	if _, ok := this.Dict().Get(key); ok {
		return false // Already exists
	}

	if this.touched = this.Dict().Set(key, uint64(this.Dict().Len())); this.touched { // A new key will be added
		this.keys = append(this.keys, key)
	}
	return true
}

// DeleteByKey removes a key from the OrderedSet.
// It returns true if the key was successfully deleted, false otherwise.
func (this *OrderedSet) DeleteByKey(key string) bool {
	idx, ok := this.Dict().Get(key)
	if !ok {
		return false
	}

	this.touched = this.Dict().Delete(key) // A key was deleted
	this.keys = append(this.keys[:idx.(uint64)], this.keys[idx.(uint64)+1:]...)

	if idx.(uint64) == uint64(len(this.keys)) { // Pop back only
		return true
	}

	current := this.Dict().GetElement(this.keys[idx.(uint64)])
	for current != nil {
		current.Value = current.Value.(uint64) - 1
		current = current.Next()
	}
	return true
}

// DeleteByIdx removes a key at the given index from the OrderedSet.
// It returns true if the key was successfully deleted, false otherwise.
func (this *OrderedSet) DeleteByIdx(idx uint64) bool {
	if idx < uint64(len(this.keys)) {
		return this.DeleteByKey(this.keys[idx])
	}
	return false
}

// Union performs a union operation with another OrderedSet.
// It adds all the keys from the other OrderedSet to the current OrderedSet.
// The current OrderedSet is modified and returned.
func (this *OrderedSet) Union(otherSet *OrderedSet) *OrderedSet {
	other := otherSet.Dict()
	for iter := other.Front(); iter != nil; iter = iter.Next() {
		this.Insert(iter.Key.(string))
	}
	return this
}
