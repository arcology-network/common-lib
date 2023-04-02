package common

func Reverse[Type any](values *[]Type) {
	for i, j := 0, len(*values)-1; i < j; i, j = i+1, j-1 {
		(*values)[i], (*values)[j] = (*values)[j], (*values)[i]
	}
}

func Fill[Type any](values *[]Type, v Type) {
	for i := 0; i < len(*values); i++ {
		(*values)[i] = v
	}
}

func RemoveIf[Type any](values *[]Type, condition func(Type) bool) {
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

func Remove[Type comparable](values *[]Type, target Type) {
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

func Foreach[Type any](values *[]Type, predicate func(v Type)) {
	for i := 0; i < len(*values); i++ {
		predicate((*values)[i])
	}
}

func FindFirst[Type comparable](values *[]Type, v Type) int {
	for i := 0; i < len(*values); i++ {
		if (*values)[i] == v {
			return i
		}
	}
	return -1
}

// Find the leftmost index of the element meeting the criteria
func FindFirstIf[Type any](values *[]Type, condition func(v Type) bool) int {
	for i := 0; i < len(*values); i++ {
		if condition((*values)[i]) {
			return i
		}
	}
	return -1
}

func FindLast[Type comparable](values *[]Type, v Type) int {
	for i := len(*values) - 1; i >= 0; i-- {
		if (*values)[i] == v {
			return i
		}
	}
	return -1
}

// Find the rightmost index of the element meeting the criteria
func FindLastIf[Type any](values *[]Type, condition func(v Type) bool) int {
	for i := len(*values) - 1; i >= 0; i-- {
		if condition((*values)[i]) {
			return i
		}
	}
	return -1
}

func DeepCopy[T any](src []T) []T {
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
