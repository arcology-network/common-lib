/*
 *   Copyright (c) 2023 Arcology Network
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package common

func Reference[T any](v T) *T   { return &v }
func Dereference[T any](v *T) T { return *v }

// New creates a new instance of a given type and returns a pointer to it.
func New[T any](v T) *T {
	v0 := T(v)
	return &v0
}

// IfThen returns one of two values based on a condition.
// If the condition is true, it returns v0; otherwise, it returns v1.
func IfThen[T any](condition bool, v0 T, v1 T) T {
	if condition {
		return v0
	}
	return v1
}

// IfThenDo1st returns one of two values based on a condition.
// If the condition is true, it calls f0 and returns its result; otherwise, it returns v1.
func IfThenDo1st[T any](condition bool, f0 func() T, v1 T) T {
	if condition {
		return f0()
	}
	return v1
}

// IfThenDo2nd returns one of two values based on a condition.
// If the condition is true, it returns v1; otherwise, it calls f0 and returns its result.
func IfThenDo2nd[T any](condition bool, v1 T, f0 func() T) T {
	if condition {
		return f0()
	}
	return v1
}

// IfThenDo executes one of two functions based on a condition.
// If the condition is true, it calls f0; otherwise, it calls f1.
func IfThenDo(condition bool, f0 func(), f1 func()) {
	if condition && f0 != nil {
		f0()
		return
	}

	if f1 != nil {
		f1()
	}
}

// IfThenDoEither returns one of two values based on a condition.
// If the condition is true, it calls f0 and returns its result; otherwise, it calls f1 and returns its result.
func IfThenDoEither[T any](condition bool, f0 func() T, f1 func() T) T {
	if condition {
		return f0()
	}
	return f1()
}

// EitherOf returns the first non-nil value between two values.
// If the first value is non-nil, it returns the first value; otherwise, it returns the second value.
func EitherOf[T any](lhv interface{}, rhv T) T {
	if lhv != nil {
		return lhv.(T)
	}
	return rhv
}

// EitherEqualsTo returns the first value if it is equal to a given value; otherwise, it returns the second value.
func EitherEqualsTo[T any](lhv interface{}, rhv T, equal func(v interface{}) bool) T {
	if equal(lhv) {
		return lhv.(T)
	}
	return rhv
}

// FilterFirst returns the first element of a pair.
func FilterFirst[T0, T1 any](v0 T0, v1 T1) T0 { return v0 }

// FilterSecond returns the second element of a pair.
func FilterSecond[T0, T1 any](v0 T0, v1 T1) T1 { return v1 }

// IsType checks if the given value is of the specified type.
// It returns true if the value is of the specified type, otherwise false.
func IsType[T any](v interface{}) bool {
	_, ok := v.(T)
	return ok
}

// Equal checks if two values are equal.
// It returns true if the values are equal; otherwise, it returns false.
func Equal[T comparable](lhv, rhv *T, wildcard func(*T) bool) bool {
	return (lhv == rhv) ||
		((lhv != nil) && (rhv != nil) && (*lhv == *rhv)) ||
		((lhv == nil && wildcard(rhv)) || (rhv == nil && wildcard(lhv)))
}

// EqualIf checks if two values are equal based on a given equality function.
// It returns true if the values are equal; otherwise, it returns false.
func EqualIf[T any](lhv, rhv *T, equal func(*T, *T) bool, wildcard func(*T) bool) bool {
	return (lhv == rhv) || ((lhv != nil) && (rhv != nil) && equal(lhv, rhv)) || ((lhv == nil && wildcard(rhv)) || (rhv == nil && wildcard(lhv)))
}
