package datacompression

import (
	"reflect"
	"testing"
)

func TestCodec(t *testing.T) {
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

	inLut := NewCompressionLut()
	//inLut.CompressOnTemp(paths)
	compressed := inLut.CompressOnTemp(paths)
	inLut.Commit()
	bytes := inLut.Encode()
	outLut := (&CompressionLut{}).Decode(bytes).(*CompressionLut)

	outLut.TryBatchUncompress(compressed)
	if !reflect.DeepEqual(paths, compressed) {
		t.Error("Error: Failed to uncompress")
	}
}
