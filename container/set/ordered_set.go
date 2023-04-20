package orderedset

import (
	"math"

	"github.com/arcology-network/common-lib/common"
	"github.com/elliotchance/orderedmap"
)

type OrderedSet struct {
	dict   *orderedmap.OrderedMap // committed keys + added - removed
	lookup []string
}

func NewOrderedSet(keys []string) *OrderedSet {
	this := &OrderedSet{
		dict:   orderedmap.NewOrderedMap(),
		lookup: keys,
	}

	for i := 0; i < len(keys); i++ {
		this.dict.Set(keys[i], uint64(i))
	}
	return this
}

func (this *OrderedSet) Len() uint64         { return uint64(this.dict.Len()) }
func (this *OrderedSet) Keys() []interface{} { return this.dict.Keys() }
func (this *OrderedSet) Deepcopy() *OrderedSet {
	if this == nil {
		return this
	}
	return &OrderedSet{this.dict.Copy(), common.DeepCopy(this.lookup)}
}

func (this *OrderedSet) Exists(key string) bool {
	_, ok := this.dict.Get(key)
	return ok
}

func (this *OrderedSet) IdxOf(key string) (uint64, bool) {
	v, ok := this.dict.Get(key)
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
	if _, ok := this.dict.Get(key); ok {
		return // Already exists
	}

	if this.dict.Set(key, uint64(this.dict.Len())) {
		this.lookup = append(this.lookup, key)
	}
}

func (this *OrderedSet) DeleteByKey(key string) bool {
	idx, ok := this.dict.Get(key)
	if !ok {
		return false
	}
	this.dict.Delete(key)
	this.lookup = append(this.lookup[:idx.(uint64)], this.lookup[idx.(uint64)+1:]...)

	if idx.(uint64) == uint64(len(this.lookup)) { // Pop back only
		return true
	}

	current := this.dict.GetElement(this.lookup[idx.(uint64)])
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

func (this *OrderedSet) Union(other *orderedmap.OrderedMap) {
	for iter := other.Front(); iter != nil; iter = iter.Next() {
		this.Insert(iter.Key.(string))
	}
}

func (this *OrderedSet) Difference(other *orderedmap.OrderedMap) {
	for iter := other.Front(); iter != nil; iter = iter.Next() {
		this.DeleteByKey(iter.Key.(string)) // could have serious performance problem
	}

	if this.dict.Len()*other.Len() > 65536 {
		for iter := other.Front(); iter != nil; iter = iter.Next() {
			this.dict.Delete(iter.Key.(string)) // could have serious performance problem
		}
	}

	for iter := other.Front(); iter != nil; iter = iter.Next() {
		this.DeleteByKey(iter.Key.(string))
		return
	}

	// could better problems
	this.lookup = this.lookup[:0]
	for iter := this.dict.Front(); iter != nil; iter = iter.Next() {
		this.lookup = append(this.lookup, iter.Key.(string))
	}
}
