/*
 *   Copyright (c) 2026 Arcology Network

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

type Stat struct {
	sizeInMem   uint64
	firstLoaded uint64
	visits      uint64
}

func (this *Stat) SetLoaded(version uint64) {
	this.firstLoaded = version
}

type entry[T any] struct {
	value T
	Stat
}

func (this *entry[T]) Size() uint64 {
	if this == nil {
		return 0
	}

	if this.sizeInMem != 0 {
		return this.sizeInMem
	}

	sized, ok := any(this.value).(interface{ MemSize() uint64 })
	if !ok {
		return 0
	}

	this.sizeInMem = sized.MemSize()
	return this.sizeInMem
}

func (this *entry[T]) Replace(NewValue T) (uint64, uint64) {
	oldSize := this.Size()
	this.value = NewValue
	this.sizeInMem = this.Size()
	this.visits++
	return oldSize, this.sizeInMem
}
