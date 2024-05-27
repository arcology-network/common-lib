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

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer([]byte{})
	},
}

func StrCat(ss ...string) string {
	if len(ss) <= 1 {
		panic("misuse of StrCat, len(ss) must be greater or equal to 2.")
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	for _, s := range ss {
		buf.WriteString(s)
	}

	str := buf.String()
	bufPool.Put(buf)
	return str
}
