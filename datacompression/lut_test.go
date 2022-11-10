package datacompression

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/arcology-network/common-lib/mhasher"
)

func TestMapDelete(t *testing.T) {
	m := make(map[string]interface{})
	m["12"] = 12
	m["34"] = 34
	m["12"] = nil

	for k, v := range m {
		fmt.Println(k, v)
	}

	delete(m, "12")
	for k, v := range m {
		fmt.Println(k, v)
	}
}

func TestFlattenStrings(t *testing.T) {
	paths := make([][]string, 100000)
	acct := RandomAccount()
	for i := 0; i < len(paths); i++ {
		for j := 0; j < 10; j++ {
			paths[i] = append(paths[i], "blcc://eth1.0/account/"+acct+"/")
		}
	}

	t0 := time.Now()
	if len(Flatten(paths)) != len(paths)*10 {
		t.Error("Error")
	}
	fmt.Println("Flatten "+fmt.Sprint(100000*10), time.Since(t0))
}

func TestCompressString(t *testing.T) {
	lut := NewCompressionLut()
	path := "1//1/1/1"
	compressed := lut.TryCompress(path)
	if compressed != "1//1/1/1" {
		t.Error("Error")
	}

	acct := RandomAccount()
	acctPath := "blcc://eth1.0/account/" + acct + "/"
	acctPathCompressed := lut.TryCompress(acctPath)

	if acctPathCompressed != "[1]/"+acct+"/" {
		t.Error("Error")
	}
}

func TestSingleAccount(t *testing.T) {
	strs := []string{"2", "1", "1"}
	newKeys, _ := mhasher.UniqueStrings(strs)
	fmt.Println(newKeys)

	acct := RandomAccount()
	if len(acct) != 40 {
		t.Error("Error: Account Address must be 40 byte long")
	}

	paths := []string{
		"blcc://eth1.0/account/" + acct + "/",
		"blcc://eth1.0/account/" + acct + "/code",
		"blcc://eth1.0/account/" + acct + "/nonce",
		"blcc://eth1.0/account/" + acct + "/balance",
		"blcc://eth1.0/account/" + acct + "/defer/",
		"blcc://eth1.0/account/" + acct + "/storage/",
		"blcc://eth1.0/account/" + acct + "/storage/containers/",
		"blcc://eth1.0/account/" + acct + "/storage/native/",
		"blcc://eth1.0/account/" + acct + "/storage/containers/!/",
		"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8111111111111111111111111111111111111",
		"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8211111111111111111111111111111111111",
		"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8311111111111111111111111111111111111",
		"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8411111111111111111111111111111111111",
	}

	source := Deepcopy(paths)
	lut := NewCompressionLut()

	compressed := lut.CompressOnTemp(paths) //LUT not the temp ?
	lut.Commit()
	fmt.Println("Compression Ratio: ", lut.GetCompressionRatio(source, compressed))
	lut.TryBatchUncompress(compressed)

	for i := 0; i < len(source); i++ {
		if len(compressed) != len(source) || source[i] != compressed[i] {
			t.Error("Error: Error happened after uncompression")
		}
	}
}

func TestShortPath(t *testing.T) {
	strs := []string{"2", "1", "1"}
	newKeys, _ := mhasher.UniqueStrings(strs)
	fmt.Println(newKeys)

	acct := RandomAccount()
	if len(acct) != 40 {
		t.Error("Error: Account Address must be 40 byte long")
	}

	path := "blcc://eth1.0/account/"
	lut := NewCompressionLut()
	lut.Commit()

	compressed := lut.TryCompress(path) //LUT not the temp ?

	if compressed != "[1]/" {
		t.Error("Error: Failed to compress the orginal string")
	}

	uncompressed := lut.TryUncompress(compressed)
	if uncompressed != "blcc://eth1.0/account/" {
		t.Error("Error: Strings don't match !")
	}
}

func TestMultiAccounts(t *testing.T) {
	paths := []string{}
	for j := 0; j < 3; j++ {
		acct := RandomAccount()
		for i := 0; i < 1; i++ {
			paths = append(paths, []string{
				"blcc://eth1.0/account/" + acct + "/",
				"blcc://eth1.0/account/" + acct + "/code",
				"blcc://eth1.0/account/" + acct + "/nonce",
				"blcc://eth1.0/account/" + acct + "/balance",
				"blcc://eth1.0/account/" + acct + "/defer/",
				"blcc://eth1.0/account/" + acct + "/storage/",
				"blcc://eth1.0/account/" + acct + "/storage/containers/",
				"blcc://eth1.0/account/" + acct + "/storage/native/",
				"blcc://eth1.0/account/" + acct + "/storage/containers/!/",
			}...)
		}
	}
	source := Deepcopy(paths)
	lut := NewCompressionLut()

	compressed := lut.CompressOnTemp(paths)
	lut.Commit()
	lut.TryBatchUncompress(compressed)

	if !reflect.DeepEqual(source, compressed) {
		t.Error("Error: Failed to uncompress")
	}

	acct := RandomAccount()
	compressed = []string{"[1]/" + acct + "/"}
	lut.TryBatchUncompress(compressed)

	if compressed[0] != "blcc://eth1.0/account/"+acct+"/" {
		t.Error("Error: Failed to uncompress")
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	accounts := make([]string, 1000000)
	for i := 0; i < len(accounts); i++ {
		accounts[i] = "abcdefghijklmnopqrestuvwxyz0123456789"
	}

	t0 := time.Now()
	byteset := make([][]byte, 1000000)
	for i := 0; i < len(accounts); i++ {
		byteset[i] = []byte(accounts[i])
	}
	fmt.Println("1000000 "+fmt.Sprint(100000*9), time.Since(t0))
}

func BenchmarkStringToBytesObjects(b *testing.B) {
	accounts := make([]string, 1000000)
	for i := 0; i < len(accounts); i++ {
		accounts[i] = "abcdefghijklmnopqrestuvwxyz0123456789"
	}

	t0 := time.Now()
	byteset := make([][]byte, 1000000)
	for i := 0; i < len(accounts); i++ {
		byteset[i] = *(*[]byte)(unsafe.Pointer(&accounts[i]))
	}
	fmt.Println("1000000 "+fmt.Sprint(100000*9), time.Since(t0))
}

func BenchmarkStringToBytesUnsafePtr(b *testing.B) {
	accounts := make([]string, 1000000)
	for i := 0; i < len(accounts); i++ {
		accounts[i] = "abcdefghijklmnopqrestuvwxyz0123456789"
	}

	t0 := time.Now()
	byteset := make([]*[]byte, 1000000)
	for i := 0; i < len(accounts); i++ {
		byteset[i] = (*[]byte)(unsafe.Pointer(&accounts[i]))
	}
	fmt.Println("1000000 "+fmt.Sprint(100000*9), time.Since(t0))
}

func BenchmarkKeccak256(b *testing.B) {
	accounts := make([]string, 1000000)
	for i := 0; i < len(accounts); i++ {
		accounts[i] = RandomAccount()
	}

	t0 := time.Now()
	byteset := make([][]byte, 1000000)
	for i := 0; i < len(accounts); i++ {
		byteset[i] = *(*[]byte)(unsafe.Pointer(&accounts[i]))
	}
	fmt.Println("1000000 "+fmt.Sprint(100000*9), time.Since(t0))

	t0 = time.Now()
	mhasher.Sha3256(byteset)
	fmt.Println("mhasher.Keccak256 :", time.Since(t0))
}
