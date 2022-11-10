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
