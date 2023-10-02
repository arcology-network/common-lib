package datacompression

import (
	"bytes"
	"fmt"
	"testing"

	codec "github.com/arcology-network/common-lib/codec"
)

func TestCompression(t *testing.T) {
	paths := []string{
		"blcc://eth1.0/account/" + "0x123456780x123456780x123456780x12345678" + "/",
		"blcc://eth1.0/account/" + "0x123456780x123456780x123456780x12345678" + "/code",
		"blcc://eth1.0/account/" + "0x123456780x123456780x123456780x12345678" + "/nonce",
		"blcc://eth1.0/account/" + "0x123456780x123456780x123456780x12345678" + "/balance",
		"blcc://eth1.0/account/" + "0x123456780x123456780x123456780x12345678" + "/defer/",
		"blcc://eth1.0/account/" + "0x123456780x123456780x123456780x12345678" + "/storage/",
		"blcc://eth1.0/account/" + "0x123456780x123456780x123456780x12345678" + "/storage/containers/",
		"blcc://eth1.0/account/" + "0x123456780x123456780x123456780x12345678" + "/storage/native/",
		"blcc://eth1.0/account/" + "0x123456780x123456780x123456780x12345678" + "/storage/containers/!/",
	}
	str := codec.Strings(paths).Flatten()
	compressed, _ := CompressGZip(str, "test", "A test string")
	fmt.Println("Uncompressed size:", len(str), " Compressed size:", len(compressed), " Ratio:", float64(len(compressed))/float64(len(str)))

	original, name, comment, _ := DecompressGZip(compressed)
	if name != "test" || comment != "A test string" || !bytes.Equal(str, original) {
		t.Error("Mismatch")
	}
}
