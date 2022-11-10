package common

func SliceToDict[T comparable](s []T) map[T]struct{} {
	dict := make(map[T]struct{})
	for _, elem := range s {
		dict[elem] = struct{}{}
	}
	return dict
}

func ToDereferencedSlice[T any](s []*T) []T {
	res := make([]T, len(s))
	for i := range s {
		res[i] = *s[i]
	}
	return res
}

func ToReferencedSlice[T any](s []T) []*T {
	res := make([]*T, len(s))
	for i := range s {
		res[i] = &s[i]
	}
	return res
}
