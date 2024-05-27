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

package types

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
)

type Addresses []ethCommon.Address

// Len()
func (as Addresses) Len() int {
	return len(as)
}

// Less():
func (as Addresses) Less(i, j int) bool {
	ibys := as[i].Bytes()
	jbys := as[j].Bytes()
	for k, ib := range ibys {
		jb := jbys[k]
		if ib < jb {
			return true
		} else if ib > jb {
			return false
		}
	}
	return true
}

// Swap()
func (as Addresses) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

func (addresses Addresses) Encode() []byte {
	return Addresses(addresses).Flatten()
}

func (addresses Addresses) Decode(data []byte) []ethCommon.Address {
	addresses = make([]ethCommon.Address, len(data)/AddressLength)
	for i := 0; i < len(addresses); i++ {
		copy(addresses[i][:], data[i*AddressLength:(i+1)*AddressLength])
	}
	return addresses
}
func (addresses Addresses) Flatten() []byte {
	buffer := make([]byte, len(addresses)*AddressLength)
	for i := 0; i < len(addresses); i++ {
		copy(buffer[i*AddressLength:(i+1)*AddressLength], addresses[i][:])
	}
	return buffer
}
