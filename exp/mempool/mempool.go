// Mempool is responsible for managing a pool of objects of the same type. It is thread-safe.
package mempool

import (
	"sync"

	"github.com/arcology-network/common-lib/container/array"
)

// Mempool represents a pool of objects of the same type.
type Mempool[T any] struct {
	new      func() T
	parent   interface{}   // Parent Mempool
	children []interface{} // Child Mempools
	objects  *array.PagedArray[T]
	counter  int
	lock     sync.Mutex
}

// NewMempool creates a new Mempool instance with the given ID and object creation function.
func NewMempool[T any](perPage, pages int, new func() T) *Mempool[T] {
	mempool := &Mempool[T]{
		new:      new,
		parent:   nil,
		children: []interface{}{},
		objects:  array.NewPagedArray[T](perPage, pages, perPage*pages),
		counter:  0,
	}

	mempool.objects.Foreach(func(i int, v *T) {
		mempool.objects.Set(i, new())
	})
	return mempool
}

// New returns an object from the Mempool.
func (this *Mempool[T]) New() T {
	this.lock.Lock()
	defer this.lock.Unlock()

	var v T
	if this.counter >= this.objects.Size() {
		v = this.new()
		this.objects.PushBack(v)
	} else {
		v = this.objects.Get(this.counter)
	}

	this.counter++
	return v
}

// Reclaim resets the Mempool, allowing the objects to be reused.
func (this *Mempool[T]) Reclaim() {
	this.counter = 0
	this.objects.Resize(this.objects.MinSize())
}

// ReclaimRecursive resets the Mempool and all its thread-local Mempools recursively.
func (m *Mempool[T]) ReclaimRecursive() {
	for _, v := range m.children {
		v.(interface{ ReclaimRecursive() }).ReclaimRecursive()
	}
	m.Reclaim()
}

// ForEachAllocated iterates over all allocated objects in the Mempool and executes the given function on each object.
// func (m *Mempool[T]) ForEachAllocated(f func(obj T)) {
// 	for _, v := range m.children {
// 		v.(interface{ ReclaimRecursive() })..(interface{ ReclaimRecursive() })(f)
// 	}

// 	for i := 0; i < m.children; i++ {
// 		f(m.objects[i])
// 	}
// }

func (this *Mempool[T]) AddToChild(child interface{}) { this.children = append(this.children, child) }
func (this *Mempool[T]) NewChildren() []interface{}   { return this.children }
func (this *Mempool[T]) SetParent(parent interface{}) { this.parent = parent }
func (this *Mempool[T]) NewParent() interface{}       { return this.parent }
