package common

func Reverse[T any](values *[]T) {
	for i, j := 0, len(*values)-1; i < j; i, j = i+1, j-1 {
		(*values)[i], (*values)[j] = (*values)[j], (*values)[i]
	}
}

func Fill[T any](values *[]T, v T) {
	for i := 0; i < len(*values); i++ {
		(*values)[i] = v
	}
}

func CopyRemoveIf[T any](values []T, condition func(T, ...interface{}) bool) []T {
	array := Clone(values)
	RemoveIf(&array, condition)
	return array
}

func RemoveIf[T any](values *[]T, condition func(T, ...interface{}) bool) {
	pos := 0
	for i := 0; i < len(*values); i++ {
		if condition((*values)[i]) {
			pos = i
			break
		}
	}

	for i := pos; i < len(*values); i++ {
		if !condition((*values)[i]) {
			(*values)[pos], (*values)[i] = (*values)[i], (*values)[pos]
			pos++
		}
	}
	(*values) = (*values)[:pos]
}

func KeepIf[T any](values *[]T, condition func(T, ...interface{}) bool) {
	pos := 0
	for i := 0; i < len(*values); i++ {
		if !condition((*values)[i]) {
			pos = i
			break
		}
	}

	for i := pos; i < len(*values); i++ {
		if condition((*values)[i]) {
			(*values)[pos], (*values)[i] = (*values)[i], (*values)[pos]
			pos++
		}
	}
	(*values) = (*values)[:pos]
}

func IfThen[T any](condition bool, v0 T, v1 T) T {
	if condition {
		return v0
	}
	return v1
}

func IfThenDo1st[T any](condition bool, f0 func() T, v1 T) T {
	if condition {
		return f0()
	}
	return v1
}

func IfThenDo2nd[T any](condition bool, f0 func() T, v1 T) T {
	if condition {
		return f0()
	}
	return v1
}

func IfThenDo[T any](condition bool, f0 func() T, f1 func() T) T {
	if condition {
		return f0()
	}
	return f1()
}

func EitherOf[T any](lhv interface{}, rhv T) T {
	if lhv != nil {
		return lhv.(T)
	}
	return rhv
}

func Remove[T comparable](values *[]T, target T) {
	pos := 0
	for i := 0; i < len(*values); i++ {
		if target == (*values)[i] {
			pos = i
			break
		}
	}

	for i := pos; i < len(*values); i++ {
		if target != (*values)[i] {
			(*values)[pos], (*values)[i] = (*values)[i], (*values)[pos]
			pos++
		}
	}
	(*values) = (*values)[:pos]
}

func Foreach[T any](values *[]T, predicate func(v T)) {
	for i := 0; i < len(*values); i++ {
		predicate((*values)[i])
	}
}

func FindFirst[T comparable](values *[]T, v T) (int, *T) {
	for i := 0; i < len(*values); i++ {
		if (*values)[i] == v {
			return i, &(*values)[i]
		}
	}
	return -1, nil
}

// Find the leftmost index of the element meeting the criteria
func FindFirstIf[T any](values *[]T, condition func(v T) bool) (int, *T) {
	for i := 0; i < len(*values); i++ {
		if condition((*values)[i]) {
			return i, &(*values)[i]
		}
	}
	return -1, nil
}

func FindLast[T comparable](values *[]T, v T) (int, *T) {
	for i := len(*values) - 1; i >= 0; i-- {
		if (*values)[i] == v {
			return i, &(*values)[i]
		}
	}
	return -1, nil
}

// Find the rightmost index of the element meeting the criteria
func FindLastIf[T any](values *[]T, condition func(v T) bool) (int, *T) {
	for i := len(*values) - 1; i >= 0; i-- {
		if condition((*values)[i]) {
			return i, &(*values)[i]
		}
	}
	return -1, nil
}

func Clone[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

func Flatten[T any](src [][]T) []T {
	totalSize := 0
	for _, data := range src {
		totalSize = totalSize + len(data)
	}
	buffer := make([]T, totalSize)
	positions := 0
	for i := range src {
		positions = positions + copy(buffer[positions:], src[i])
	}
	return buffer
}

func ConcateFrom[T0, T1 any](array []T0, getter func(T0) []T1) []T1 {
	total := 0
	for i := 0; i < len(array); i++ {
		total += len(getter(array[i]))
	}
	output := make([]T1, total) // Pre-allocation for better performance

	offset := 0
	for i := 0; i < total; i++ {
		elems := getter(array[i])
		copy(output[offset:], elems)
		offset += len(elems)
	}
	return output
}

func CastTo[T0, T1 any](src []T0, predicate func(T0) T1) []T1 {
	target := make([]T1, len(src))
	for i := range src {
		target[i] = predicate(src[i])
	}
	return target
}

func EqualArray[T comparable](lhv []T, rhv []T) bool {
	if len(lhv) != len(rhv) {
		return false
	}

	for _, v0 := range lhv {
		flag := false
		for _, v1 := range rhv {
			if v0 == v1 {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}

	for _, v0 := range rhv {
		flag := false
		for _, v1 := range lhv {
			if v0 == v1 {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func MergeMaps[M ~map[K]V, K comparable, V any](from, to M) M {
	for k, v := range to {
		from[k] = v
	}
	return from
}

// func MergeMapsIf[M ~map[K]V, K comparable, V any](from, to M, func()) M {
// 	for k, v := range to {
// 		from[k] = v
// 	}
// 	return from
// }

func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

func MapKVs[M ~map[K]V, K comparable, V any](m M) ([]K, []V) {
	keys := make([]K, len(m))
	values := make([]V, len(m))
	i := 0
	for k, v := range m {
		keys[i] = k
		values[i] = v
		i++
	}
	return keys, values
}
