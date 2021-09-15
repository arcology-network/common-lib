package datacompression

import (
	"reflect"
	"testing"
)

func TestCodec(t *testing.T) {
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

	inLut := NewCompressionLut()
	compressed := inLut.BatchCompress(paths)
	bytes := inLut.Encode()
	outLut := (&CompressionLut{}).Decode(bytes).(*CompressionLut)

	uncompressed := outLut.BatchUncompress(compressed)
	if !reflect.DeepEqual(paths, uncompressed) {
		t.Error("Error: Failed to uncompress")
	}
}

// func TestCodecPerform(t *testing.T) {
// 	lut := NewCompressionLut()
// 	for i := 0; i < 100000; i++ {
// 		lut.FindIdx([]string{fmt.Sprint(i)})
// 	}

// 	t0 := time.Now()
// 	bytes := lut.Encode()
// 	fmt.Println("Encode():", time.Now().Sub(t0))

// 	t0 = time.Now()
// 	(&CompressionLut{}).Decode(bytes)
// 	fmt.Println("Decode():", time.Now().Sub(t0))
// }
