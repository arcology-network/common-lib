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

package addrcompressor

import (
	"strconv"

	common "github.com/arcology-network/common-lib/common"
	slice "github.com/arcology-network/common-lib/exp/slice"
)

func (this *CompressionLut) GetNewAccounts(originals []string) []string {
	acctLen := 40
	prefixLen := len("blcc://eth1.0/account/")

	keys := slice.Transform(originals, func(_ int, v string) string { return v[prefixLen : prefixLen+acctLen] })
	return this.filterExistingKeys(slice.Unique(keys, func(str0, str1 string) bool { return str0 < str1 }), this.dict) // Get new keys
}

func (this *CompressionLut) CompressStaticKey(original string) string {
	acctLen := 40
	prefixLen := len("blcc://eth1.0/account/")

	if len(original) < prefixLen {
		return original
	}

	var prefixid int
	k := original[:prefixLen-1]
	if v, ok := this.dict.Get(k); ok {
		prefixid = int(v.(uint32))
	} else {
		return original
	}

	if len(original) < prefixLen+acctLen {
		original = "[" + strconv.Itoa(int(prefixid)) + "]" + original[prefixLen:]
		return original
	}

	key := original[prefixLen : prefixLen+acctLen]
	if id, ok := this.dict.Get(key); ok {
		original = "[" + strconv.Itoa(int(prefixid)) + "]" + "/[" + strconv.Itoa(int(id.(uint32)+this.offset)) + "]" + original[prefixLen+acctLen:]
	} else {
		if id, ok := this.tempLut.dict.Get(key); ok {
			original = "[" + strconv.Itoa(int(prefixid)) + "]" + "/[" + strconv.Itoa(int(id.(uint32)+this.length)) + "]" + original[prefixLen+acctLen:]
		}
	}
	return original
}

func (this *CompressionLut) CompressStaticKeys(originals []string) []string {
	replacer := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			originals[i] = this.CompressStaticKey(originals[i])
		}
	}
	common.ParallelWorker(len(originals), 4, replacer)
	return originals
}
