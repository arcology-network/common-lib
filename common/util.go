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

func HexToString(src []byte) string {
	shex := hex.EncodeToString(src)
	return "0x" + shex
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

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
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

func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

func ArrayCopy(data []byte) []byte {
	datas := make([]byte, len(data))
	copy(datas, data)
	return datas
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

func AppendToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content + "\n")
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return err
	}
	file.Sync()

	return nil
}

func AddToLogFile(filename, field string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		AppendToFile(filename, "Marshal err : "+err.Error())
		return
	}
	AppendToFile(filename, field+" : "+string(data))
}

func RemoveNils(values *[]interface{}) {
	pos := int64(-1)
	for i := 0; i < len((*values)); i++ {
		if pos < 0 && (*values)[i] == nil {
			pos = int64(i)
			continue
		}

		if pos < 0 && (*values)[i] != nil {
			continue
		}

		if pos >= 0 && (*values)[i] == nil {
			(*values)[pos] = (*values)[i]
			continue
		}

		(*values)[pos] = (*values)[i]
		pos++
	}

	if pos >= 0 {
		(*values) = (*values)[:pos]
	}
}

func RemoveDuplicateStrings(strs *[]string) []string {
	dict := make(map[string]bool)
	for i := 0; i < len(*strs); i++ {
		dict[(*strs)[i]] = true
	}

	uniqueStrs := make([]string, 0, len(dict))
	for k := range dict {
		uniqueStrs = append(uniqueStrs, k)
	}
	return uniqueStrs
}

func Remove(values *[]interface{}, condition func(interface{}) bool) {
	pos := int64(-1)
	for i := 0; i < len((*values)); i++ {
		if pos < 0 && condition((*values)[i]) {
			pos = int64(i)
			continue
		}

		if pos < 0 && condition((*values)[i]) {
			continue
		}

		if pos >= 0 && condition((*values)[i]) {
			(*values)[pos] = (*values)[i]
			continue
		}

		(*values)[pos] = (*values)[i]
		pos++
	}

	if pos >= 0 {
		(*values) = (*values)[:pos]
	}
}

func ReverseString(s string) string {
	reversed := []byte(s)
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	return *(*string)(unsafe.Pointer(&reversed))
}
