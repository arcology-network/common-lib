// Mempool is responsible for managing a pool of objects of the same type. It is thread-safe.
package mempool

import (
	"fmt"
	"sync"
)

// Mempool represents a pool of objects of the same type.
type Mempool struct {
	initializer func() interface{}
	id          string
	pools       map[string]*Mempool
	objects     []interface{}
	next        int
	parent      *Mempool
	lock        sync.Mutex
}

// NewMempool creates a new Mempool instance with the given ID and object creation function.
func NewMempool(id interface{}, initializer func() interface{}) *Mempool {
	return &Mempool{
		id:          fmt.Sprintf("%v", id),
		initializer: initializer,
		pools:       make(map[string]*Mempool),
		objects:     nil,
		next:        0,
		parent:      nil,
	}
}

// GetPool returns the thread-local Mempool associated with the given ID.
// If the thread-local Mempool does not exist, it creates a new one.
func (m *Mempool) GetPool(id interface{}) *Mempool {
	m.lock.Lock()
	defer m.lock.Unlock()

	newId := fmt.Sprintf("%v%v", m.id, id)
	if m.parent != nil {
		return m.parent.GetPool(newId)
	}

	if m.pools[newId] == nil {
		newMempool := NewMempool(newId, m.initializer)
		newMempool.parent = m
		m.pools[newId] = newMempool
	}
	return m.pools[newId]
}

// Get returns an object from the Mempool.
func (m *Mempool) Get() interface{} {
	// if len(m.objects) <= m.next {
	// 	m.objects = append(m.objects, make([]interface{}, 1024)...)
	// }

	// if m.objects[m.next] == nil {
	// 	m.objects[m.next] = m.initializer()
	// }
	// m.next++
	// return m.objects[m.next-1]

	return m.initializer()
}

// Reclaim resets the Mempool, allowing the objects to be reused.
func (m *Mempool) Reclaim() {
	m.next = 0
}

// ReclaimRecursive resets the Mempool and all its thread-local Mempools recursively.
func (m *Mempool) ReclaimRecursive() {
	for _, v := range m.pools {
		v.ReclaimRecursive()
	}
	m.Reclaim()
}

// ForEachAllocated iterates over all allocated objects in the Mempool and executes the given function on each object.
func (m *Mempool) ForEachAllocated(f func(obj interface{})) {
	for _, v := range m.pools {
		v.ForEachAllocated(f)
	}

	for i := 0; i < m.next; i++ {
		f(m.objects[i])
	}
}
