// Mempool is responsible for managing a pool of pool of the same type. It is thread-safe.
package mempool

import (
	"sync"

	indexedslice "github.com/arcology-network/common-lib/container/slice"
)

// Mempool represents a pool of pool of the same type.
type Mempool[T any] struct {
	new      func() T
	resetter func(T)
	parent   interface{}   // Parent Mempool
	children []interface{} // Child Mempools
	pool     *indexedslice.PagedSlice[T]
	counter  int
	lock     sync.Mutex
}

// NewMempool creates a new Mempool instance with the given ID and object creation function.
func NewMempool[T any](perPage, numPages int, new func() T, resetter func(T)) *Mempool[T] {
	mempool := &Mempool[T]{
		new:      new,
		resetter: resetter,
		parent:   nil,
		children: []interface{}{},
		pool:     indexedslice.NewPagedSlice[T](perPage, numPages, perPage*numPages),
		counter:  0,
	}

	mempool.pool.Foreach(func(i int, v *T) {
		mempool.pool.Set(i, new())
	})
	return mempool
}

func (this *Mempool[T]) Size() int    { return this.pool.Size() }
func (this *Mempool[T]) MinSize() int { return this.pool.MinSize() }

// New returns an object from the Mempool.

// Note: This function is very slow because of the lock.
// It about 6.5 times slower than using the new() function, directly
// It is better to remove the lock but this may cause some problems
// So if a thead can have its own mempool, then there is no need for the lock.
// This will tremendously improve the performance.
func (this *Mempool[T]) New() T {
	return this.new()
	// this.lock.Lock()
	// defer this.lock.Unlock()

	// var v T
	// if this.counter >= this.pool.Size() {
	// 	v = this.new()
	// } else {
	// 	// v = this.pool.Get(this.counter)
	// 	v = this.pool.PopBack()
	// }

	// this.counter++
	// return v
}

func (this *Mempool[T]) Return(v T) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.resetter(v)
	this.pool.PushBack(v)
}

// New returns an array of pool from the Mempool.
func (this *Mempool[T]) NewArray(num int) []T {
	this.lock.Lock()
	defer this.lock.Unlock()

	arr := make([]T, num)
	for i := 0; i < num; i++ {
		var v T
		if this.counter >= this.pool.Size() {
			arr[i] = this.new()
			this.pool.PushBack(v)
		} else {
			arr[i] = this.pool.Get(this.counter)
		}
	}
	this.counter += num
	return arr
}

// Reclaim resets the Mempool, allowing the pool to be reused.
func (this *Mempool[T]) Reset() {
	this.counter = 0
	this.pool.Resize(this.pool.MinSize())
	// this.pool.Foreach(func(i int, v *T) {
	// 	this.resetter(*v)
	// })

	this.pool = indexedslice.NewPagedSlice[T](this.pool.PageSize(), this.pool.NumPages(), this.pool.PageSize()*this.pool.NumPages())
	this.pool.Foreach(func(i int, v *T) {
		this.pool.Set(i, this.new())
	})
}

// ReclaimRecursive resets the Mempool and all its thread-local Mempools recursively.
func (m *Mempool[T]) ReclaimRecursive() {
	for _, v := range m.children {
		v.(interface{ ReclaimRecursive() }).ReclaimRecursive()
	}
	m.Reset()
}

func (this *Mempool[T]) AddToChild(child interface{}) { this.children = append(this.children, child) }
func (this *Mempool[T]) NewChildren() []interface{}   { return this.children }
func (this *Mempool[T]) SetParent(parent interface{}) { this.parent = parent }
func (this *Mempool[T]) NewParent() interface{}       { return this.parent }
