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

package mempool

import (
	"fmt"
	"testing"
	"time"
)

func TestPagedIntArray(t *testing.T) {
	type CustomType struct {
		a int
		b [20]byte
		e string
	}

	i := 0
	pool := NewMempool[*int](1, 2, func() *int {
		i++
		return &i
	})

	if *pool.New() != 1 {
		t.Error("Error: Wrong value")
	}

	if *pool.New() != 2 {
		t.Error("Error: Wrong value")
	}

	if *pool.New() != 3 {
		t.Error("Error: Wrong value")
	}
}

func TestPagedArrayCustomTypes(t *testing.T) {
	type CustomType struct {
		a int
		b [20]byte
		e string
	}

	i := 0
	pool := NewMempool[*CustomType](1, 2, func() *CustomType {
		i++
		return &CustomType{
			a: i,
			b: [20]byte{},
			e: "hello" + fmt.Sprint(i),
		}
	})

	if pool.New().a != 1 {
		t.Error("Error: Wrong value")
	}

	if pool.New().a != 2 {
		t.Error("Error: Wrong value")
	}

	if pool.New().a != 3 {
		t.Error("Error: Wrong value")
	}
	pool.Reclaim()

	// Reset the init function
	i = 99
	pool.new = func() *CustomType {
		return &CustomType{
			a: i,
			b: [20]byte{},
			e: "hello" + fmt.Sprint(i),
		}
	}
	if v := pool.New().a; v != 1 {
		t.Error("Error: Wrong value", v)
	}

	if pool.New().a != 2 {
		t.Error("Error: Wrong value")
	}

	if pool.New().a != 99 {
		t.Error("Error: Wrong value")
	}

	v := pool.New()
	v.a = 10

	v = pool.New()
	v.a = 11

	v = pool.New()
	v.a = 12

	v = pool.New()
	v.a = 13

	v = pool.New()
	v.a = 14
}

func BenchmarkTestPagedArrayCustomTypes(t *testing.B) {
	type CustomType struct {
		a int
		b [20]byte
		e string
	}

	i := 0
	pool := NewMempool[*CustomType](4096, 156, func() *CustomType {
		i++
		return &CustomType{
				a: i,
				b: [20]byte{},
				e: "hello" + fmt.Sprint(i),
			},
			
	},func(_ *CustomType) {}})


	vs := make([]*CustomType, 1000000)
	t0 := time.Now()
	for i := 0; i < 1000000; i++ {
		vs[i] = &CustomType{
			a: i,
			b: [20]byte{},
			e: "hello" + fmt.Sprint(i),
		}
	}
	fmt.Println("New 1 ", "1000000", time.Since(t0))

	t0 = time.Now()
	for i := 0; i < 1000000; i++ {
		pool.New()
	}
	pool.Reset()
	fmt.Println("pool.New() 1 ", "1000000", time.Since(t0))

	t0 = time.Now()
	for i := 0; i < 1000000; i++ {
		pool.New()
	}
	pool.Reset()
	fmt.Println("pool.New() 2 ", "1000000", time.Since(t0))
}
