package datacompression

import (
	"strconv"
)

type Compressor struct {
	keyIdxDict map[string]uint64
	idxKeyDict map[uint64]string
	newKeys    []string
}

func NewCompressor(keys []string) *Compressor {
	dict := &Compressor{
		keyIdxDict: make(map[string]uint64),
		idxKeyDict: make(map[uint64]string),
		newKeys:    []string{},
	}
	for _, k := range keys {
		dict.Compress(k, nil)
	}
	return dict
}

func (this *Compressor) Compress(key string, getter func(string) string) string {
	if getter != nil {
		key = getter(key)
	}

	if v, ok := this.keyIdxDict[key]; ok {
		return strconv.Itoa(int(v))
	}
	this.newKeys = append(this.newKeys, key)

	this.idxKeyDict[uint64(len(this.idxKeyDict))] = key
	this.keyIdxDict[key] = uint64(len(this.keyIdxDict))
	return strconv.Itoa(len(this.keyIdxDict) - 1)
}

func (this *Compressor) Decompress(compressed string, getter func(string) string) string {
	idx, _ := strconv.Atoi(compressed)
	if v, ok := this.idxKeyDict[uint64(idx)]; ok {
		if getter != nil {
			return getter(v)
		}
		return v
	}
	return ""
}

func (this *Compressor) Length() uint64 {
	return uint64(len(this.keyIdxDict) - len(this.newKeys))
}

func (this *Compressor) NewKeys() []string {
	return this.newKeys
}
