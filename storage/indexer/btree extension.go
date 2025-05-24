package indexer

import "github.com/google/btree"

// FindFirstIf finds the first element in a slice that satisfies a given condition and returns its index and a pointer to the element
func FindFirstIf[T any](tree *btree.BTreeG[T], condition func(_ int, v T) bool) (int, T) {
	i := 0
	var val T
	getter := func(v T) bool {
		i++
		val = v
		return !condition(i-1, v) // Should be !condition
	}
	tree.Ascend(getter)
	return i, val
}

func Export[T any](tree *btree.BTreeG[T]) []T {
	vals := make([]T, 0, tree.Len())
	getter := func(v T) bool {
		vals = append(vals, v)
		return true
	}

	tree.Ascend(getter)
	return vals
}
