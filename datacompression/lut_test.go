package datacompression

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	codec "github.com/arcology-network/common-lib/codec"
)

func TestCompressUncompressBasicSingleAccount(t *testing.T) {
	paths := []string{
		"blcc://eth1.0/account/" + "0x112345" + "/",
		"blcc://eth1.0/account/" + "0x112345" + "/code",
		"blcc://eth1.0/account/" + "0x112345" + "/nonce",
		"blcc://eth1.0/account/" + "0x112345" + "/balance",
		"blcc://eth1.0/account/" + "0x112345" + "/defer/",
		"blcc://eth1.0/account/" + "0x112345" + "/storage/",
		"blcc://eth1.0/account/" + "0x112345" + "/storage/containers/",
		"blcc://eth1.0/account/" + "0x112345" + "/storage/native/",
		"blcc://eth1.0/account/" + "0x112345" + "/storage/containers/!/",
	}

	lut := NewCompressionLut()
	t0 := time.Now()
	compressed := lut.BatchCompress(paths)
	fmt.Println("BatchCompress "+fmt.Sprint(len(paths)), " in ", time.Since(t0))

	t0 = time.Now()
	uncompressed := lut.BatchUncompress(compressed)
	fmt.Println("Uncompressed "+fmt.Sprint(len(paths)), " in ", time.Since(t0))

	if !reflect.DeepEqual(paths, uncompressed) {
		t.Error("Error: Failed to uncompress")
	}
}

func TestCompressUncompressBasicMultipleAccounts(t *testing.T) {
	paths := []string{}
	accounts := []string{"alice", "bob", "carol", "dave"}
	for _, acct := range accounts {
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

	lut := NewCompressionLut()
	compressed := lut.BatchCompress(paths)
	uncompressed := lut.BatchUncompress(compressed)

	ratio := float32(len(codec.Strings(compressed).Encode())) / float32(len(codec.Strings(paths).Encode()))
	fmt.Println("---- Compression ratio " + fmt.Sprint(ratio))

	if !reflect.DeepEqual(paths, uncompressed) {
		t.Error("Error: Failed to uncompress")
	}
}

func BenchmarkUncompressSameAccount(b *testing.B) {
	paths := []string{}
	for i := 0; i < 100000; i++ {
		acct := fmt.Sprint(rand.Float64())
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

	sections := make([][]string, 100000)
	for i := 0; i < 100000; i++ {
		sections[i] = paths[i*9 : (i+1)*9]
	}

	t0 := time.Now()
	lut := NewCompressionLut()
	for i := 0; i < len(sections); i++ {
		lut.BatchUncompress(sections[i])
	}
	fmt.Println("Compressed then Uncompressed "+fmt.Sprint(len(paths)), " in ", time.Since(t0))
}

func BenchmarkUncompressDifferentAccount100k(b *testing.B) {
	paths := []string{}
	//accounts := []string{"alice", "bob", "carol", "dave"}
	for j := 0; j < 4; j++ {
		for i := 0; i < 25000; i++ {
			acct := fmt.Sprint(rand.Float64())
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

	lut := NewCompressionLut()
	t0 := time.Now()
	compressed := lut.BatchCompress(paths)
	fmt.Println("BatchCompress "+fmt.Sprint(len(paths)), " in ", time.Since(t0))

	t0 = time.Now()
	uncompressed := lut.BatchUncompress(compressed)
	fmt.Println("Uncompressed "+fmt.Sprint(len(paths)), " in ", time.Since(t0))

	ratio := float32(len(codec.Strings(compressed).Encode())) / float32(len(codec.Strings(paths).Encode()))
	fmt.Println("---- Compression ratio " + fmt.Sprint(ratio))

	if !reflect.DeepEqual(paths, uncompressed) {
		b.Error("Error: Failed to uncompress")
	}
}
