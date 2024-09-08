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

package cache

import (
	redblacktree "github.com/emirpasic/gods/v2/trees/redblacktree"
	"golang.org/x/exp/constraints"
)

type CacheIndexer[I constraints.Unsigned, K comparable, V any] struct{ *redblacktree.Tree[I, *[]*K] }

func NewCacheIndex[I constraints.Unsigned, K comparable, V any](sizer func(V) uint64) *CacheIndexer[I, K, V] {
	return &CacheIndexer[I, K, V]{Tree: redblacktree.New[I, *[]*K]()}
}

func (this *CacheIndexer[I, K, V]) Insert(idx I, keys []K, values []V) {
	for _, key := range keys {
		keyVec := this.Tree.GetNode(I(idx)).Value
		if keyVec == nil {
			this.Tree.Put(idx, keyVec)
		}
		*keyVec = append(*keyVec, &key)
	}
}

func (this *CacheIndexer[I, K, V]) Delete(idx I, keys []K) {
	for _, key := range keys {
		keyVec := this.Tree.GetNode(I(idx)).Value
		if keyVec == nil {
			this.Tree.Put(idx, keyVec)
		}
		*keyVec = append(*keyVec, &key)
	}
}
