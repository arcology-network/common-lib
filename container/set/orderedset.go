package orderedset

import (
	"math"

	"github.com/arcology-network/common-lib/common"
	"github.com/elliotchance/orderedmap"
)

type OrderedSet struct {
	dict    *orderedmap.OrderedMap // committed keys + added - removed
	keys    []string
	touched bool
}

func NewOrderedSet(keys []string) *OrderedSet {
	this := &OrderedSet{
		dict:    orderedmap.NewOrderedMap(),
		keys:    keys,
		touched: false,
	}
	return this
}

func (this *OrderedSet) Equal(other *OrderedSet) bool {
	if this == other && this == nil {
		return true
	}

	if this == nil || other == nil {
		return false
	}

	return common.EqualArray(this.keys, other.keys) &&
		this.touched == other.touched &&
		this.isSynced() == other.isSynced()
}

func (this *OrderedSet) Length() int {
	return len(this.keys)
}

func (this *OrderedSet) isSynced() bool {
	return (this.dict.Len()) == len(this.keys)
}

// Sync the look up with the
func (this *OrderedSet) Dict() *orderedmap.OrderedMap {
	if !this.isSynced() {
		if this.dict.Len() > 0 {
			this.dict = orderedmap.NewOrderedMap() // This should never happen
		}

		for i := 0; i < len(this.keys); i++ {
			this.dict.Set(this.keys[i], uint64(i))
		}
	}
	return this.dict
}

func (this *OrderedSet) Touched() bool { return this.touched }
func (this *OrderedSet) Len() uint64   { return uint64(len(this.keys)) }
func (this *OrderedSet) Keys() []string {
	return this.keys
}

// func (this *OrderedSet) Dict() *orderedmap.OrderedMap { return this.Sync() }
func (this *OrderedSet) Clone() interface{} {
	if this == nil {
		return this
	}
	set := NewOrderedSet(this.keys)
	set.touched = this.touched
	return set
}

func (this *OrderedSet) Exists(key string) bool {
	_, ok := this.Dict().Get(key)
	return ok
}

func (this *OrderedSet) Delete(key string) bool { return this.DeleteByKey(key) }

func (this *OrderedSet) IdxOf(key string) uint64 {
	v, ok := this.Dict().Get(key)
	if !ok {
		return math.MaxUint64
	}
	return v.(uint64)
}

func (this *OrderedSet) KeyAt(idx uint64) string {
	if idx < uint64(len(this.keys)) {
		return this.keys[idx]
	}
	return ""
}

func (this *OrderedSet) Insert(key string) {
	if _, ok := this.Dict().Get(key); ok {
		return // Already exists
	}

	if this.touched = this.Dict().Set(key, uint64(this.Dict().Len())); this.touched { // A new key will be added
		this.keys = append(this.keys, key)
	}
}

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

func (this *OrderedSet) DeleteByIdx(idx uint64) bool {
	if idx < uint64(len(this.keys)) {
		return this.DeleteByKey(this.keys[idx])
	}
	return false
}

func (this *OrderedSet) Union(otherSet *OrderedSet) *OrderedSet {
	other := otherSet.Dict()
	for iter := other.Front(); iter != nil; iter = iter.Next() {
		this.Insert(iter.Key.(string))
	}
	return this
}

// func (this *OrderedSet) Difference(otherSet *OrderedSet) *OrderedSet {
// 	other := otherSet.Dict()
// 	for iter := other.Front(); iter != nil; iter = iter.Next() {
// 		this.DeleteByKey(iter.Key.(string)) // could have serious performance problem
// 	}

// 	if this.Dict().Len()*other.Len() > 65536 {
// 		for iter := other.Front(); iter != nil; iter = iter.Next() {
// 			this.Dict().Delete(iter.Key.(string)) // could have serious performance problem
// 		}
// 	}

// 	for iter := other.Front(); iter != nil; iter = iter.Next() {
// 		this.DeleteByKey(iter.Key.(string))
// 		return
// 	}

// 	// could better problems
// 	this.keys = this.keys[:0]
// 	dict := this.Dict()
// 	for iter := dict.Front(); iter != nil; iter = iter.Next() {
// 		this.keys = append(this.keys, iter.Key.(string))
// 	}
// 	return this
// }
