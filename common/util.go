package common

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"sync"
	"unsafe"

	ethCommon "github.com/HPISTechnologies/3rd-party/eth/common"
	"github.com/HPISTechnologies/common-lib/encoding"
	"github.com/google/uuid"
)

const (
	TypeByteSize = 4
	ThreadNum    = 4
)

func ToNewHash(h ethCommon.Hash, height, round uint64) ethCommon.Hash {
	keys := Uint64ToBytes(height)
	keys = append(keys, Uint64ToBytes(round)...)
	keys = append(keys, h.Bytes()...)
	newhash := sha256.Sum256(keys)
	return ethCommon.BytesToHash(newhash[:])
}

func JsonEncode(x interface{}) ([]byte, error) {
	return json.Marshal(x)
}
func JsonDecode(data []byte, x interface{}) error {
	return json.Unmarshal(data, x)
}

func ToHexStr(src []byte) string {
	shex := hex.EncodeToString(src)
	return "0x" + shex
}

// find b in a ,return idx ,if not exist return -1
func FindinArrays(a [][]byte, b []byte) int {
	for i, v := range a {
		if bytes.Equal(v, b) {
			return i
		}
	}
	return -1
}

func BytesToUint64(array []byte) uint64 {
	return binary.BigEndian.Uint64(array)
}
func BytesToUint32(array []byte) uint32 {
	return binary.BigEndian.Uint32(array)
}
func BytesToUint16(array []byte) uint16 {
	return binary.BigEndian.Uint16(array)
}
func Uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
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

type TestElement struct {
	D2Array [][]byte
}

func Int64ToUint64(src1 int64) uint64 {
	return *(*uint64)(unsafe.Pointer(&src1))
}

func Uint64ToInt64(src1 uint64) int64 {
	return *(*int64)(unsafe.Pointer(&src1))
}

func GetUniqueValue(relations *map[string]string) []string {
	unique := map[string]string{}
	for _, topic := range *relations {
		unique[topic] = topic // get unique topics
	}

	topics := []string{}
	for topic, _ := range unique {
		topics = append(topics, topic)
	}
	return topics
}

func GenerateUUID() uint64 {
	uuid := uuid.New()
	return binary.BigEndian.Uint64(uuid[0:9])
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

func FileToLines(fileName string) []string {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func JsonToCsv(lines []string) ([]string, [][]string) {
	logs := make(map[string][]string)
	var result map[string]interface{}
	for _, line := range lines {
		json.Unmarshal([]byte(line), &result)
		for k, v := range result {
			logs[k] = append(logs[k], fmt.Sprintf("%v", v))
		}
	}

	columns := make([]string, 0, len(logs))
	rows := make([][]string, 0, len(logs))
	for k, v := range logs {
		columns = append(columns, k)
		rows = append(rows, v)
	}
	return columns, Transpose(rows)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Make a deep copy from src into dst.
func Deepcopy(dst interface{}, src interface{}) error {
	if dst == nil {
		return fmt.Errorf("dst cannot be nil")
	}
	if src == nil {
		return fmt.Errorf("src cannot be nil")
	}
	bytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("Unable to marshal src: %s", err)
	}
	err = json.Unmarshal(bytes, dst)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal into dst: %s", err)
	}
	return nil
}

func GenerateRanges(length int, numThreads int) []int {
	ranges := make([]int, 0, numThreads+1)
	step := int(math.Ceil(float64(length) / float64(numThreads)))
	for i := 0; i <= numThreads; i++ {
		ranges = append(ranges, int(math.Min(float64(step*i), float64(length))))
	}
	return ranges
}

func ParallelWorker(total, nThds int, worker func(start, end, idx int, args ...interface{}), args ...interface{}) {
	idxRanges := GenerateRanges(total, nThds)
	var wg sync.WaitGroup
	for i := 0; i < len(idxRanges)-1; i++ {
		wg.Add(1)
		go func(start int, end int, idx int) {
			defer wg.Done()
			if start != end {
				worker(start, end, idx, args)
			}
		}(idxRanges[i], idxRanges[i+1], i)
	}
	wg.Wait()
}

func Intersected(lft [][32]byte, rgt [][32]byte) bool {
	for i := range lft {
		for j := range rgt {
			if bytes.Equal(lft[i][:], rgt[j][:]) {
				return true
			}
		}
	}
	return false
}
func Serialization(filename string, obj interface{}) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		fmt.Printf("openfile err=%v\n", err)
		return
	}
	defer file.Close()
	data, err := JsonEncode(obj)
	if err == nil {
		file.WriteString(fmt.Sprintf("%x", data) + "\n")
	} else {
		fmt.Printf("JsonEncode err=%v\n", err)
	}
	file.Sync()
}

func CalculateHash(hashes []*ethCommon.Hash) ethCommon.Hash {
	if len(hashes) == 0 {
		return ethCommon.Hash{}
	}
	datas := make([][]byte, len(hashes))
	for i := range hashes {
		datas[i] = hashes[i].Bytes()
	}
	hash := sha256.Sum256(encoding.Byteset(datas).Encode())
	return ethCommon.BytesToHash(hash[:])
}

func ArrayCopy(data []byte) []byte {
	datas := make([]byte, len(data))
	copy(datas, data)
	return datas
}
func Array2DCopy(src [][]byte) [][]byte {
	elements := make([][]byte, len(src))
	for i, row := range src {
		elements[i] = ArrayCopy(row)
	}
	return elements
}
func Flatten(src [][]byte) []byte {
	totalSize := 0
	for _, data := range src {
		totalSize = totalSize + len(data)
	}
	buffer := make([]byte, totalSize)
	positions := 0
	for i := range src {
		positions = positions + copy(buffer[positions:], src[i])
	}

	return buffer
}

func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}
