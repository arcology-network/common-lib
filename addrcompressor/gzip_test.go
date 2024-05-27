/*
 *   Copyright (c) 2024 Arcology Network

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package addrcompressor

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
