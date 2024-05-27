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
