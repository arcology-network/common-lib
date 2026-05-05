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

type CRDT interface { // value type
	TypeID() uint8
	Equal(any) bool
	Clone() any

	IsNumeric() bool
	IsCommutative() bool // If the type is commutative, the order of the operands does not matter.

	Value() any
	Delta() (any, bool)

	Limits() (any, any) // Get the limits of the type, if applicable.
	IsDeltaApplied() bool

	// Delta replication methods
	New(any, any, any, any, any) any
	CloneDelta() (any, bool)
	SetDelta(any, bool)
	SetValue(v any)
	GetCascadeSub(string, any) []string // Get the sub paths for cascade delete, if applicable.

	Get() (any, uint32, uint32) // Value, reads and writes, no deltawrites.
	Set(any, any) (any, uint32, uint32, uint32, error)
	CopyTo(any) (any, uint32, uint32, uint32) // Only a function to generate the right access counts, when assigning the value.
	ApplyDelta([]CRDT) (CRDT, int, error)
	IsDeletable(any, any) bool

	MemSize() uint64 // Size in memory
	Size() uint64    // Encoded size

	Encode() []byte
	EncodeTo([]byte) int
	Decode([]byte) any

	Preload(string, any)

	// Auxiliary methods
	Hash() [32]byte
	ShortHash() (uint64, bool) // For fast comparison only.
	Print()
}

// type Writer[T any] interface {
// 	Import([]T)
// 	Precommit(bool) error //should return a error
// 	Commit(uint64) error  //should return a error

// 	IsSync() bool // If the writer is synchronous, it will block until the commit is done.
// 	Name() string
// }

type Hasher func(CRDT) []byte

func SizeOf(T any) uint64 {
	switch v := T.(type) {
	case CRDT:
		return v.MemSize()
	case []byte:
		return uint64(len(v))
	}
	panic("Unsupported type for SizeOf")
	return 0
}
