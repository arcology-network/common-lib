package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"

	"unsafe"

	"github.com/arcology-network/common-lib/encoding"
	evmCommon "github.com/arcology-network/evm/common"
	"github.com/google/uuid"
)

const (
	TypeByteSize = 4
	ThreadNum    = 4
)

func ToNewHash(h evmCommon.Hash, height, round uint64) evmCommon.Hash {
	keys := Uint64ToBytes(height)
	keys = append(keys, Uint64ToBytes(round)...)
	keys = append(keys, h.Bytes()...)
	newhash := sha256.Sum256(keys)
	return evmCommon.BytesToHash(newhash[:])
}

func HexToString(src []byte) string {
	shex := hex.EncodeToString(src)
	return "0x" + shex
}

func BytesToUint32(array []byte) uint32 {
	return binary.BigEndian.Uint32(array)
}

func Uint32ToBytes(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}

func Uint16ToBytes(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

func BytesToUint16(array []byte) uint16 {
	return binary.BigEndian.Uint16(array)
}

func Uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

func BytesToUint64(array []byte) uint64 {
	return binary.BigEndian.Uint64(array)
}

func Int64ToUint64(src1 int64) uint64 {
	return *(*uint64)(unsafe.Pointer(&src1))
}

func Uint64ToInt64(src1 uint64) int64 {
	return *(*int64)(unsafe.Pointer(&src1))
}

func GenerateUUID() uint64 {
	uuid := uuid.New()
	return binary.BigEndian.Uint64(uuid[0:9])
}

func Transpose(slice [][]string) [][]string {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

func GobEncode(x interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(x)
	if err != nil {
		return nil, err
	}
	retData := buf.Bytes()
	return retData, nil
}

func GobDecode(data []byte, x interface{}) error {
	bufTo := bytes.NewBuffer(data)
	dec := gob.NewDecoder(bufTo)
	err := dec.Decode(x)
	if err != nil {
		return err
	}
	return nil
}

// func UniqueIf[T comparable](strs *[]T, equal func(interface{}, interface{}) bool) []T {
// 	dict := make(map[T]bool)
// 	for i := 0; i < len(*strs); i++ {
// 		dict[(*strs)[i]] = true
// 	}

//		uniques := make([]T, 0, len(dict))
//		for k := range dict {
//			uniques = append(uniques, k)
//		}
//		return uniques
//	}

func CalculateHash(hashes []*evmCommon.Hash) evmCommon.Hash {
	if len(hashes) == 0 {
		return evmCommon.Hash{}
	}
	datas := make([][]byte, len(hashes))
	for i := range hashes {
		datas[i] = hashes[i].Bytes()
	}
	hash := sha256.Sum256(encoding.Byteset(datas).Encode())
	return evmCommon.BytesToHash(hash[:])
}
