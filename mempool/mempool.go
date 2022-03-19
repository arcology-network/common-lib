package mempool

import (
	"fmt"
	"sync"
)

type Mempool struct {
	nf      func() interface{}
	id      string
	tls     map[string]*Mempool
	objects []interface{}
	next    int
	parent  *Mempool
	guard   sync.Mutex
}

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

func (m *Mempool) Get() interface{} {
	if len(m.objects) <= m.next {
		m.objects = append(m.objects, make([]interface{}, 1024)...)
	}

	if m.objects[m.next] == nil {
		m.objects[m.next] = m.nf()
	}
	m.next++
	return m.objects[m.next-1]
}

func (m *Mempool) Reclaim() {
	// fmt.Println("Mempool.Reclaim:", m.id, "next =", m.next)
	m.next = 0
}

func (m *Mempool) ReclaimRecursive() {
	// fmt.Println("Mempool.ReclaimRecursive:", m.id)
	for _, v := range m.tls {
		v.ReclaimRecursive()
	}
	m.Reclaim()
}

func (m *Mempool) ForEachAllocated(f func(obj interface{})) {
	for _, v := range m.tls {
		v.ForEachAllocated(f)
	}

	for i := 0; i < m.next; i++ {
		f(m.objects[i])
	}
}
