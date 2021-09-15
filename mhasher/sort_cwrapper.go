// +build !nometri

package mhasher

/*
#cgo CFLAGS: -Iinclude
#cgo LDFLAGS: -L/usr/local/lib/ -lmhasher
#include "bytes.external.h"
#include "stdlib.h"

*/
import "C" //must flow above
import (
	"unsafe"

	"github.com/arcology-network/common-lib/codec"
)

type StringEngine struct {
	enginePtr unsafe.Pointer
}

func (s *StringEngine) Stop() {
	C.Stop(s.enginePtr)
}
func (s *StringEngine) Clear() {
	C.Clear(s.enginePtr)
}

func Start() *StringEngine {
	s := StringEngine{}
	s.enginePtr = C.Start()
	return &s
}

func (s *StringEngine) ToBuffer(paths []string) error {
	//ToBuffer(void* engine, char* path, char* pathLen, uint32_t count);
	pathLenth := len(paths)
	pathLenthVec := make([]uint32, pathLenth)
	for i := range paths {
		pathLenthVec[i] = uint32(len(paths[i]))
	}
	_, err := C.ToBuffer(
		s.enginePtr,
		(*C.char)(unsafe.Pointer(&paths[0])),
		(*C.char)(unsafe.Pointer(&pathLenthVec[0])),
		C.uint32_t(pathLenth),
	)
	return err
}

func (s *StringEngine) FromBuffer(paths []string) ([]string, error) {
	//FromBuffer(void* engine, char* buffer, char* outputLens, uint32_t* count);
	BufferSize := 2048
	DataLenth := 100
	buffer := make([]byte, BufferSize*DataLenth)

	dataLengthVec := make([]uint32, DataLenth)
	dataNums := make([]uint32, 1)
	_, err := C.FromBuffer(
		s.enginePtr,
		(*C.char)(unsafe.Pointer(&buffer[0])),
		(*C.char)(unsafe.Pointer(&dataLengthVec[0])),
		(*C.uint32_t)(unsafe.Pointer(&dataNums[0])),
	)

	if err != nil {
		return []string{}, err
	}
	retDataSize := int(dataNums[0])
	retpaths := make([]string, retDataSize)
	pos := 0
	for i := 0; i < retDataSize; i++ {
		retpaths[i] = string(buffer[pos : pos+int(dataLengthVec[i])])
		pos = pos + int(dataLengthVec[i])
	}
	return retpaths, nil
}

func SortString(datas []string) ([]string, error) {
	lengthVec, dataLength := orderDataString(datas)
	totalBytes := codec.Strings(datas).Flatten()
	indices := make([]uint32, dataLength)
	_, err := C.Sort(
		(*C.char)(unsafe.Pointer(&totalBytes[0])),
		(*C.uint32_t)(unsafe.Pointer(&lengthVec[0])),
		C.uint32_t(dataLength),
		(*C.uint32_t)(unsafe.Pointer(&indices[0])),
	)
	if err != nil {
		return datas, err
	}
	results := make([]string, dataLength)
	for i, idx := range indices {
		results[i] = datas[int(idx)]
	}
	return results, nil
}
func UniqueString(datas []string) ([]string, error) {
	lengthVec, dataLength := orderDataString(datas)
	totalBytes := codec.Strings(datas).Flatten()
	indices := make([]uint8, dataLength)
	_, err := C.Unique(
		(*C.char)(unsafe.Pointer(&totalBytes[0])),
		(*C.uint32_t)(unsafe.Pointer(&lengthVec[0])),
		C.uint32_t(dataLength),
		(*C.uint8_t)(unsafe.Pointer(&indices[0])),
	)
	if err != nil {
		return datas, err
	}
	results := make([]string, 0, dataLength)
	for i, flag := range indices {
		if flag == uint8(255) {
			results = append(results, datas[i])
		}
	}
	return results, nil
}
func RemoveString(datas, toRemove []string) ([]string, error) {
	lengthVec, dataLength := orderDataString(datas)
	removeLengthVec, removeDataLength := orderDataString(toRemove)
	totalBytes := codec.Strings(datas).Flatten()
	totalBytesRemove := codec.Strings(toRemove).Flatten()
	indices := make([]uint8, dataLength)
	_, err := C.Remove(
		(*C.char)(unsafe.Pointer(&totalBytes[0])),
		(*C.uint32_t)(unsafe.Pointer(&lengthVec[0])),
		C.uint32_t(dataLength),
		(*C.char)(unsafe.Pointer(&totalBytesRemove[0])),
		(*C.uint32_t)(unsafe.Pointer(&removeLengthVec[0])),
		C.uint32_t(removeDataLength),
		(*C.uint8_t)(unsafe.Pointer(&indices[0])),
	)

	if err != nil {
		return datas, err
	}

	results := make([]string, 0, dataLength)
	for i, flag := range indices {
		if flag == uint8(255) {
			results = append(results, datas[i])
		}
	}
	return results, nil
}
func orderDataString(datas []string) ([]uint32, uint32) {
	dataLength := len(datas)
	if dataLength == 0 {
		return []uint32{}, 0
	}
	lengthVec := make([]uint32, dataLength)
	for i, data := range datas {
		lengthVec[i] = uint32(len(data))
	}
	return lengthVec, uint32(dataLength)
}

func SortBytes(data [][]byte) ([][]byte, error) {
	totalBytes, lengthVec, dataLength := orderData(data)

	indices := make([]uint32, dataLength)
	_, err := C.Sort(
		(*C.char)(unsafe.Pointer(&totalBytes[0])),
		(*C.uint32_t)(unsafe.Pointer(&lengthVec[0])),
		C.uint32_t(dataLength),
		(*C.uint32_t)(unsafe.Pointer(&indices[0])),
	)
	if err != nil {
		return data, err
	}
	results := make([][]byte, dataLength)
	for i, idx := range indices {
		results[i] = data[int(idx)]
	}
	return results, nil
}

func UniqueBytes(data [][]byte) ([][]byte, error) {
	totalBytes, lengthVec, dataLength := orderData(data)

	indices := make([]uint8, dataLength)
	_, err := C.Unique(
		(*C.char)(unsafe.Pointer(&totalBytes[0])),
		(*C.uint32_t)(unsafe.Pointer(&lengthVec[0])),
		C.uint32_t(dataLength),
		(*C.uint8_t)(unsafe.Pointer(&indices[0])),
	)
	if err != nil {
		return data, err
	}
	results := make([][]byte, 0, dataLength)
	for i, flag := range indices {
		if flag == uint8(255) {
			results = append(results, data[i])
		}
	}
	return results, nil
}

func RemoveBytes(data, toRemove [][]byte) ([][]byte, error) {
	totalBytes, lengthVec, dataLength := orderData(data)
	removeTotalBytes, removeLengthVec, removeDataLength := orderData(toRemove)

	indices := make([]uint8, dataLength)
	_, err := C.Remove(
		(*C.char)(unsafe.Pointer(&totalBytes[0])),
		(*C.uint32_t)(unsafe.Pointer(&lengthVec[0])),
		C.uint32_t(dataLength),
		(*C.char)(unsafe.Pointer(&removeTotalBytes[0])),
		(*C.uint32_t)(unsafe.Pointer(&removeLengthVec[0])),
		C.uint32_t(removeDataLength),
		(*C.uint8_t)(unsafe.Pointer(&indices[0])),
	)

	if err != nil {
		return data, err
	}

	results := make([][]byte, 0, dataLength)
	for i, flag := range indices {
		if flag == uint8(255) {
			results = append(results, data[i])
		}
	}
	return results, nil
}

func orderData(datas [][]byte) ([]byte, []uint32, uint32) {
	dataLength := len(datas)
	if dataLength == 0 {
		return []byte{}, []uint32{}, 0
	}
	lengthVec := make([]uint32, dataLength)
	for i, data := range datas {
		lengthVec[i] = uint32(len(data))
	}
	totalBytes := codec.Byteset(datas).Flatten()
	return totalBytes, lengthVec, uint32(dataLength)
}
