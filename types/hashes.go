package types

import (
	"sync"

	ethCommon "github.com/ethereum/go-ethereum/common"

	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"

	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func Size(hash ethCommon.Hash) uint32 {
	return uint32(ethCommon.HashLength)
}

func Encode(hash ethCommon.Hash) []byte {
	return hash[:]
}

func Decode(data []byte) ethCommon.Hash {
	hash := ethCommon.Hash{}
	copy(hash[:], data)
	return hash
}

func Checksum(hash ethCommon.Hash) []byte {
	return Encode(hash)
}

func ToUint32(hash ethCommon.Hash) uint32 {
	return binary.BigEndian.Uint32(hash[0:4])
}

type Hashes []ethCommon.Hash

func (hashes Hashes) Intersected(lft []ethCommon.Hash, rgt []ethCommon.Hash) bool {
	for i := range lft {
		for j := range rgt {
			if bytes.Equal(lft[i][:], rgt[j][:]) {
				return true
			}
		}
	}
	return false
}

func (hashes Hashes) Checksum() ethCommon.Hash {
	combined := make([]ethCommon.Hash, 64)
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

func (hashes Hashes) Decode(data []byte) []ethCommon.Hash {
	hashes = make([]ethCommon.Hash, len(data)/ethCommon.HashLength)
	for i := 0; i < len(hashes); i++ {
		copy(hashes[i][:], data[i*ethCommon.HashLength:(i+1)*ethCommon.HashLength])
	}
	return hashes
}

func (hashes Hashes) Size() uint32 {
	return uint32(len(hashes) * ethCommon.HashLength)
}

func (hashes Hashes) Flatten() []byte {
	buffer := make([]byte, len(hashes)*ethCommon.HashLength)
	for i := 0; i < len(hashes); i++ {
		copy(buffer[i*ethCommon.HashLength:(i+1)*ethCommon.HashLength], hashes[i][:])
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

type Hashset [][]ethCommon.Hash

func (hashes Hashset) Flatten() []ethCommon.Hash {
	buffer := make([]ethCommon.Hash, len(hashes))
	for i := 0; i < len(hashes); i++ {
		buffer = append(buffer, hashes[i]...)
	}
	return buffer
}

type Hashgroup [][][]ethCommon.Hash

func (hashgroup Hashgroup) Flatten() [][]ethCommon.Hash {
	buffer := make([][]ethCommon.Hash, len(hashgroup))
	for i := 0; i < len(hashgroup); i++ {
		buffer = append(buffer, hashgroup[i]...)
	}
	return buffer
}

func RlpHash(x interface{}) (h ethCommon.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func TxHash(data []byte) (h ethCommon.Hash) {
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
