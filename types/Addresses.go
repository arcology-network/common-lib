package types

import (
	evmCommon "github.com/arcology-network/evm/common"
)

type Addresses []evmCommon.Address

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

func (addresses Addresses) Decode(data []byte) []evmCommon.Address {
	addresses = make([]evmCommon.Address, len(data)/AddressLength)
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
