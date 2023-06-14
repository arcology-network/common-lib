package types

import (
	"sync"

	evmCommon "github.com/arcology-network/evm/common"

	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"

	"github.com/arcology-network/evm/rlp"
	"golang.org/x/crypto/sha3"
)

func Size(hash evmCommon.Hash) uint32 {
	return uint32(evmCommon.HashLength)
}

func Encode(hash evmCommon.Hash) []byte {
	return hash[:]
}

func Decode(data []byte) evmCommon.Hash {
	hash := evmCommon.Hash{}
	copy(hash[:], data)
	return hash
}

func Checksum(hash evmCommon.Hash) []byte {
	return Encode(hash)
}

func ToUint32(hash evmCommon.Hash) uint32 {
	return binary.BigEndian.Uint32(hash[0:4])
}

type Hashes []evmCommon.Hash

func (hashes Hashes) Intersected(lft []evmCommon.Hash, rgt []evmCommon.Hash) bool {
	for i := range lft {
		for j := range rgt {
			if bytes.Equal(lft[i][:], rgt[j][:]) {
				return true
			}
		}
	}
	return false
}

func (hashes Hashes) Checksum() evmCommon.Hash {
	combined := make([]evmCommon.Hash, 64)
	worker := func(start, end int, args ...interface{}) {
		stride := int(math.Ceil(float64(len(hashes)) / float64(len(combined))))
		i := int(math.Ceil(float64(start) / float64(stride)))
		combined[i] = sha256.Sum256(Hashes(hashes)[start:end].Flatten())
	}
	ParallelWorker(len(hashes), len(combined), worker)
	return sha256.Sum256(Hashes(combined).Flatten())
}

func (hashes Hashes) Encode() []byte {
	return Hashes(hashes).Flatten()
}

func (hashes Hashes) Decode(data []byte) []evmCommon.Hash {
	hashes = make([]evmCommon.Hash, len(data)/evmCommon.HashLength)
	for i := 0; i < len(hashes); i++ {
		copy(hashes[i][:], data[i*evmCommon.HashLength:(i+1)*evmCommon.HashLength])
	}
	return hashes
}

func (hashes Hashes) Size() uint32 {
	return uint32(len(hashes) * evmCommon.HashLength)
}

func (hashes Hashes) Flatten() []byte {
	buffer := make([]byte, len(hashes)*evmCommon.HashLength)
	for i := 0; i < len(hashes); i++ {
		copy(buffer[i*evmCommon.HashLength:(i+1)*evmCommon.HashLength], hashes[i][:])
	}
	return buffer
}

func (hashes Hashes) ToUint32s() []uint32 {
	keys := make([]uint32, len(hashes))
	converter := func(start, end int, args ...interface{}) {
		for i := start; i < end; i++ {
			keys[i] = ToUint32(hashes[i])
		}
	}
	ParallelWorker(len(keys), 8, converter)
	return keys
}

type Hashset [][]evmCommon.Hash

func (hashes Hashset) Flatten() []evmCommon.Hash {
	buffer := make([]evmCommon.Hash, len(hashes))
	for i := 0; i < len(hashes); i++ {
		buffer = append(buffer, hashes[i]...)
	}
	return buffer
}

type Hashgroup [][][]evmCommon.Hash

func (hashgroup Hashgroup) Flatten() [][]evmCommon.Hash {
	buffer := make([][]evmCommon.Hash, len(hashgroup))
	for i := 0; i < len(hashgroup); i++ {
		buffer = append(buffer, hashgroup[i]...)
	}
	return buffer
}

func RlpHash(x interface{}) (h evmCommon.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func TxHash(data []byte) (h evmCommon.Hash) {
	return RlpHash(data[1:])
}
func ParallelWorker(total, nThds int, worker func(start, end int, args ...interface{}), args ...interface{}) {
	idxRanges := GenerateRanges(total, nThds)
	var wg sync.WaitGroup
	for i := 0; i < len(idxRanges)-1; i++ {
		wg.Add(1)
		go func(start int, end int) {
			defer wg.Done()
			if start != end {
				worker(start, end, args)
			}
		}(idxRanges[i], idxRanges[i+1])
	}
	wg.Wait()
}
func GenerateRanges(length int, numThreads int) []int {
	ranges := make([]int, 0, numThreads+1)
	step := int(math.Ceil(float64(length) / float64(numThreads)))
	for i := 0; i <= numThreads; i++ {
		ranges = append(ranges, int(math.Min(float64(step*i), float64(length))))
	}
	return ranges
}
