// Mempool is responsible for managing a pool of objects of the same type. It is thread-safe.
package mempool

import (
	"fmt"
	"sync"
)

// Mempool represents a pool of objects of the same type.
type Mempool struct {
	nf      func() interface{}
	id      string
	tls     map[string]*Mempool
	objects []interface{}
	next    int
	parent  *Mempool
	guard   sync.Mutex
}

// NewMempool creates a new Mempool instance with the given ID and object creation function.
func NewMempool(id interface{}, nf func() interface{}) *Mempool {
	return &Mempool{
		id:      fmt.Sprintf("%v", id),
		nf:      nf,
		tls:     make(map[string]*Mempool),
		objects: nil,
		next:    0,
		parent:  nil,
	}
}

// GetTlsMempool returns the thread-local Mempool associated with the given ID.
// If the thread-local Mempool does not exist, it creates a new one.
func (m *Mempool) GetTlsMempool(id interface{}) *Mempool {
	m.guard.Lock()
	defer m.guard.Unlock()

	newId := fmt.Sprintf("%v%v", m.id, id)
	if m.parent != nil {
		return m.parent.GetTlsMempool(newId)
	}

	if m.tls[newId] == nil {
		newMempool := NewMempool(newId, m.nf)
		newMempool.parent = m
		m.tls[newId] = newMempool
	}
	return m.tls[newId]
}

// Get returns an object from the Mempool.
func (m *Mempool) Get() interface{} {
	// if len(m.objects) <= m.next {
	// 	m.objects = append(m.objects, make([]interface{}, 1024)...)
	// }

	// if m.objects[m.next] == nil {
	// 	m.objects[m.next] = m.nf()
	// }
	// m.next++
	// return m.objects[m.next-1]

	return m.nf()
}

// Reclaim resets the Mempool, allowing the objects to be reused.
func (m *Mempool) Reclaim() {
	m.next = 0
}

// ReclaimRecursive resets the Mempool and all its thread-local Mempools recursively.
func (m *Mempool) ReclaimRecursive() {
	for _, v := range m.tls {
		v.ReclaimRecursive()
	}
	m.Reclaim()
}

// ForEachAllocated iterates over all allocated objects in the Mempool and executes the given function on each object.
func (m *Mempool) ForEachAllocated(f func(obj interface{})) {
	for _, v := range m.tls {
		v.ForEachAllocated(f)
	}

	for i := 0; i < m.next; i++ {
		f(m.objects[i])
	}
}
