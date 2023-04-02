package orderedset

import (
	"math"

	"github.com/elliotchance/orderedmap"
)

type OrderedSet struct {
	keyDict *orderedmap.OrderedMap // committed keys + added - removed
	lookup  []string
}

func NewOrderedSet(keys []string) *OrderedSet {
	this := &OrderedSet{
		keyDict: orderedmap.NewOrderedMap(),
		lookup:  keys,
	}

	for i := 0; i < len(keys); i++ {
		this.keyDict.Set(keys[i], uint64(i))
	}
	return this
}

func (this *OrderedSet) Size() uint64 { return uint64(this.keyDict.Len()) }

func (this *OrderedSet) IdxOf(key string) (uint64, bool) {
	v, ok := this.keyDict.Get(key)
	if !ok {
		return math.MaxUint64, false
	}
	return v.(uint64), ok
}

func (this *OrderedSet) KeyOf(idx uint64) (interface{}, bool) {
	if idx < uint64(len(this.lookup)) {
		return this.lookup[idx], true
	}
	return nil, false
}

func (this *OrderedSet) Insert(key string) {
	if _, ok := this.keyDict.Get(key); ok {
		return // Already exists
	}

	if this.keyDict.Set(key, uint64(this.keyDict.Len())) {
		this.lookup = append(this.lookup, key)
	}
}

func (this *OrderedSet) DeleteByKey(key string) bool {
	idx, ok := this.keyDict.Get(key)
	if !ok {
		return false
	}
	this.keyDict.Delete(key)
	this.lookup = append(this.lookup[:idx.(uint64)], this.lookup[idx.(uint64)+1:]...)

	if idx.(uint64) == uint64(len(this.lookup)) { // Pop back only
		return true
	}

	current := this.keyDict.GetElement(this.lookup[idx.(uint64)])
	for current != nil {
		current.Value = current.Value.(uint64) - 1
		current = current.Next()
	}
	return true
}

func (this *OrderedSet) DeleteByIdx(idx uint64) bool {
	if idx < uint64(len(this.lookup)) {
		return this.DeleteByKey(this.lookup[idx])
	}
	return false
}
