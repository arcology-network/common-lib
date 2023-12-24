// The Compressor class is a data structure that is used for compressing and decompressing data. It provides a way
// to map keys to their corresponding indices and vice versa. This can be useful in scenarios where you want to reduce
// the size of data by replacing repetitive keys with their corresponding indices

package datacompression

import (
	"strconv"
)

type Compressor struct {
	keyIdxDict map[string]uint64 // maps keys to their corresponding indices.
	idxKeyDict map[uint64]string // maps indices to their corresponding keys.
	newKeys    []string          // stores newly added keys.
}

// NewCompressor creates a new Compressor instance with the given keys.
// It initializes the keyIdxDict and idxKeyDict maps.
// It also compresses each key using the Compress method.
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

// Compress compresses the given key using the getter function if provided.
// It checks if the key is already present in the keyIdxDict.
// If it is, it returns the corresponding index as a string.
// Otherwise, it adds the key to the newKeys slice, updates the dictionaries, and returns the new index as a string.
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

// Decompress decompresses the given compressed string using the getter function if provided.
// It converts the compressed string to an index and checks if it exists in the idxKeyDict.
// If it does, it returns the corresponding key.
// Otherwise, it returns an empty string.
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

// Length returns the number of keys in the Compressor instance.
// It subtracts the number of newKeys from the total number of keys in the keyIdxDict.
func (this *Compressor) Length() uint64 {
	return uint64(len(this.keyIdxDict) - len(this.newKeys))
}

// NewKeys returns the slice of newly added keys.
func (this *Compressor) NewKeys() []string {
	return this.newKeys
}
