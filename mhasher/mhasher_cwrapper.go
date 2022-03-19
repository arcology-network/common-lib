// +build !nometri

package mhasher

/*
#cgo CFLAGS: -Iinclude

#cgo LDFLAGS: -L/usr/local/lib/ -lmhasher

#include "mhasher.external.h"

#include "stdlib.h"

*/
import "C" //must flow above
import (
	"bytes"
	"errors"
	"unsafe"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/encoding"
)

const (
	HashType_160 = 20
	HashType_256 = 32
)

func UniqueKeysByMap(keys [][]byte) [][]byte {
	mp := make(map[string]int, len(keys))
	for i := range keys {
		mp[string(keys[i])] = i
	}
	outBys := make([][]byte, len(mp))
	i := 0
	for k, _ := range mp {
		outBys[i] = []byte(k)
		i = i + 1
	}
	return outBys
}

func UniqueKeys(keys [][]byte) ([][]byte, error) {
	if len(keys) == 0 {
		return [][]byte{}, nil
	}

	chars := encoding.Byteset(keys).Flatten()
	c_char := (*C.char)(unsafe.Pointer(&chars[0]))

	count := len(keys)
	c_count := C.uint64_t(count)

	step := len(keys[0])

	aIndex := make([]uint64, 1)
	a_index := (*C.uint64_t)(unsafe.Pointer(&aIndex[0]))

	ahash := make([]byte, len(chars))
	a_char := (*C.char)(unsafe.Pointer(&ahash[0]))

	_, err := C.UniqueHash256(c_char, c_count, a_char, a_index)

	outBytes := make([][]byte, int(aIndex[0]))
	if aIndex[0] == 0 {
		return [][]byte{}, errors.New("return data is empty")
	}

	for i := 0; i < int(aIndex[0]); i++ {
		outBytes[i] = ahash[step*i : step*(i+1)]
	}

	return outBytes, err
}

func SortByHash(hashes []ethCommon.Hash) ([]uint64, error) {
	if len(hashes) == 0 {
		return []uint64{}, nil
	}
	chars := ethCommon.Hashes(hashes).Flatten()
	c_char := (*C.char)(unsafe.Pointer(&chars[0]))

	count := len(hashes)
	c_count := C.uint64_t(count)

	rIndex := make([]uint64, count)
	r_index := (*C.uint64_t)(unsafe.Pointer(&rIndex[0]))

	var err error
	_, err = C.SortHash256(c_char, c_count, r_index)

	return rIndex, err
}

//is used in monaco/core/types/tx.go
func Roothash(ls [][]byte, HashType int) ([]byte, error) {
	ahash := make([]byte, HashType)
	if len(ls) == 0 {
		return ahash, nil
	}
	srcLenth := len(ls)
	var buffer bytes.Buffer
	for i := 0; i < srcLenth; i++ {

		buffer.Write(ls[i])
	}
	srcStr := buffer.Bytes()
	num := C.uint64_t(srcLenth)
	c_char := (*C.char)(unsafe.Pointer(&srcStr[0]))

	a_char := (*C.char)(unsafe.Pointer(&ahash[0]))

	var err error
	if HashType == HashType_160 {
		_, err = C.BinaryMhasherRIPEMD160(c_char, num, a_char)
	} else {
		_, err = C.BinaryMhasherKeccak256(c_char, num, a_char)
	}
	return ahash, err
}

func Keccak256(byteSet [][]byte) [][]byte {
	length := len(byteSet)
	bytes := codec.Byteset(byteSet).Flatten()
	dataLens := make([]uint64, length)
	for i := range byteSet {
		dataLens[i] = uint64(len(byteSet[i]))
	}

	buffer := make([]byte, length*32)
	C.keccak256(
		(*C.char)(unsafe.Pointer(&bytes[0])),
		(*C.uint64_t)(unsafe.Pointer(&dataLens[0])),
		(C.uint64_t)(length),
		(*C.char)(unsafe.Pointer(&buffer[0])),
	)
	return Reshape(buffer, 32)
}

func Ripemd160(byteSet [][]byte) [][]byte {
	length := len(byteSet)
	bytes := codec.Byteset(byteSet).Flatten()
	dataLens := make([]uint64, length)
	for i := range byteSet {
		dataLens[i] = uint64(len(byteSet[i]))
	}

	buffer := make([]byte, length*20)
	C.keccak256(
		(*C.char)(unsafe.Pointer(&bytes[0])),
		(*C.uint64_t)(unsafe.Pointer(&dataLens[0])),
		(C.uint64_t)(length),
		(*C.char)(unsafe.Pointer(&buffer[0])),
	)
	return Reshape(buffer, 20)
}

func Sha3256(byteSet [][]byte) [][]byte {
	length := len(byteSet)
	bytes := codec.Byteset(byteSet).Flatten()
	dataLens := make([]uint64, length)
	for i := range byteSet {
		dataLens[i] = uint64(len(byteSet[i]))
	}

	buffer := make([]byte, length*32)
	C.sha3256(
		(*C.char)(unsafe.Pointer(&bytes[0])),
		(*C.uint64_t)(unsafe.Pointer(&dataLens[0])),
		(C.uint64_t)(length),
		(*C.char)(unsafe.Pointer(&buffer[0])),
	)
	return Reshape(buffer, 32)
}

func Reshape(bytes []byte, hashSize int) [][]byte {
	hashes := make([][]byte, len(bytes)/hashSize)
	for i := range hashes {
		hashes[i] = bytes[i*hashSize : (i+1)*hashSize]
	}
	return hashes
}
