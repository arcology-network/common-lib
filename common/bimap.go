/*
 *   Copyright (c) 2024 Arcology Network

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

package common

type BiMap[K comparable, V comparable] struct {
	k2v map[K]V
	v2k map[V]K
}

func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{
		k2v: make(map[K]V),
		v2k: make(map[V]K),
	}
}

func (b *BiMap[K, V]) Add(k K, v V) {
	b.k2v[k] = v
	b.v2k[v] = k
}

func (b *BiMap[K, V]) Get(k K) V {
	return b.k2v[k]
}

func (b *BiMap[K, V]) GetInverse(v V) K {
	return b.v2k[v]
}
