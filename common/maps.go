package common

import (
	"golang.org/x/exp/maps"
)

func MergeMaps[M ~map[K]V, K comparable, V any](m1, m2 M) M {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	return maps.Keys(m)
}

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	return maps.Values(m)
}

func MapKVs[M ~map[K]V, K comparable, V any](m M) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
