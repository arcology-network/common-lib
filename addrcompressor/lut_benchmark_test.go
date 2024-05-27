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
	"fmt"
	"testing"
	"time"
)

func BenchmarkUncompressAllTogether(b *testing.B) {
	paths := []string{}
	for i := 0; i < 70000; i++ {
		acct := RandomAccount()
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
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8111111111111111111111111111111111111",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8211111111111111111111111111111111111",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8311111111111111111111111111111111111",
			"blcc://eth1.0/account/" + acct + "/storage/containers/KittyIndexToOwner/$ad90f8411111111111111111111111111111111111",
		}...)
	}
	t0 := time.Now()
	//source := Deepcopy(paths)
	fmt.Println("Deepcopy "+fmt.Sprint(len(paths)), " in ", time.Since(t0))

	lut := NewCompressionLut()
	t0 = time.Now()
	compressed := lut.CompressOnTemp(paths)
	fmt.Println("Compress On temp Dict "+fmt.Sprint(len(paths)), " in ", time.Since(t0))

	t0 = time.Now()
	lut.Commit()
	fmt.Println("Commit "+fmt.Sprint(len(paths)), " in ", time.Since(t0))
	lut.TryBatchUncompress(compressed)

	// if !reflect.DeepEqual(source, compressed) {
	// 	b.Error("Error: Error happened after uncompression")
	// }
}
