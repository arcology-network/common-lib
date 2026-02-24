/*
*   Copyright (c) 2026 Arcology Network

*   This program is free software: you can redistribute it and/or modify
*   it under the terms of the GNU General Public License as published by
*   the Free Software Foundation, either version 3 of the License, or
*   (at your option) any later version.

*   This program is distributed in the hope that it will be useful,
*   but WITHOUT ANY WARRANTY; without even the implied warranty of
*   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*   GNU General Public License for more details.

*   You should have received a copy of the GNU General Public License
*   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package bimap

type BiMap[K comparable, V comparable] struct {
	kv map[K]V
	vk map[V]K
}

func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{kv: map[K]V{}, vk: map[V]K{}}
}

func (m *BiMap[K, V]) Set(k K, v V) {
	if oldV, ok := m.kv[k]; ok {
		delete(m.vk, oldV)
	}
	if oldK, ok := m.vk[v]; ok {
		delete(m.kv, oldK)
	}
	m.kv[k] = v
	m.vk[v] = k
}

func (m *BiMap[K, V]) GetByKey(k K) (V, bool) {
	v, ok := m.kv[k]
	return v, ok
}

func (m *BiMap[K, V]) GetByValue(v V) (K, bool) {
	k, ok := m.vk[v]
	return k, ok
}

func (m *BiMap[K, V]) DeleteByKey(k K) {
	if v, ok := m.kv[k]; ok {
		delete(m.kv, k)
		delete(m.vk, v)
	}
}

func (m *BiMap[K, V]) DeleteByValue(v V) {
	if k, ok := m.vk[v]; ok {
		delete(m.vk, v)
		delete(m.kv, k)
	}
}
