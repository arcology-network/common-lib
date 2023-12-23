package badgerdb

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

var (
	TEST_ROOT_PATH   = path.Join(os.TempDir(), "/badgerdb/")
	TEST_BACKUP_PATH = path.Join(os.TempDir(), "/badgerdb-back/")
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

var containers = []string{
	"abcde",
	"fghijkl",
	"mnopqrstuv",
	"wxy",
	"z",
}

var accounts = make([]string, 10000000)
var contracts = make([]string, 10000)

func setup() ([]string, [][]byte) {
	var keys []string
	var values [][]byte
	timer("generate urls", func() {
		for i := range accounts {
			address, ks, vs := generateAccountUrls()
			accounts[i] = address
			keys = append(keys, ks...)
			values = append(values, vs...)
		}
		for i := range contracts {
			address, ks, vs := generateContractUrls()
			contracts[i] = address
			keys = append(keys, ks...)
			values = append(values, vs...)
		}
	})

	// hashes := make([]string, len(keys))
	// timer("calculate hashes", func() {
	// 	for i := range keys {
	// 		hashes[i] = string(sum256([]byte(keys[i])))
	// 	}
	// })

	// return hashes, values
	return keys, values
}

func sum256(bytes []byte) []byte {
	sum := sha256.Sum256(bytes)
	return sum[:]
}

func timer(step string, f func()) time.Duration {
	start := time.Now()
	f()
	d := time.Since(start)
	fmt.Printf("%s: %v\n", step, d)
	return d
}

func newBlock() ([]string, [][]byte) {
	var keys []string
	var values [][]byte
	timer("generate transitions", func() {
		for i := 0; i < 25000; i++ {
			ks, vs := generateTransferUrls()
			keys = append(keys, ks...)
			values = append(values, vs...)
		}
		for i := 0; i < 25000; i++ {
			ks, vs := generateContractCallUrls()
			keys = append(keys, ks...)
			values = append(values, vs...)
		}
	})

	// hashes := make([]string, len(keys))
	// timer("calculate hashes", func() {
	// 	for i := range keys {
	// 		hashes[i] = string(sum256([]byte(keys[i])))
	// 	}
	// })

	// return hashes, values
	return keys, values
}

func generateContractCallUrls() ([]string, [][]byte) {
	from := accounts[rand.Intn(len(accounts))]
	to := contracts[rand.Intn(len(contracts))]
	return []string{
			// "blcc://eth1.0/account/" + from + "/nonce",
			// "blcc://eth1.0/account/" + from + "/balance",
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[0] + "/",
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[0] + "/" + from,
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[1] + "/",
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[1] + "/" + from,
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[2] + "/",
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[2] + "/" + from,
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[3] + "/",
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[3] + "/" + from,
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[4] + "/",
			// "blcc://eth1.0/account/" + to + "/storage/containers/" + containers[4] + "/" + from,
			from + "/nonce",
			from + "/balance",
			to + "/storage/containers/" + containers[0] + "/",
			to + "/storage/containers/" + containers[0] + "/" + from,
			to + "/storage/containers/" + containers[1] + "/",
			to + "/storage/containers/" + containers[1] + "/" + from,
			to + "/storage/containers/" + containers[2] + "/",
			to + "/storage/containers/" + containers[2] + "/" + from,
			to + "/storage/containers/" + containers[3] + "/",
			to + "/storage/containers/" + containers[3] + "/" + from,
			to + "/storage/containers/" + containers[4] + "/",
			to + "/storage/containers/" + containers[4] + "/" + from,
		}, [][]byte{
			make([]byte, 4),
			make([]byte, 8),
			make([]byte, 256),
			make([]byte, 32),
			make([]byte, 256),
			make([]byte, 32),
			make([]byte, 256),
			make([]byte, 32),
			make([]byte, 256),
			make([]byte, 32),
			make([]byte, 256),
			make([]byte, 1024),
		}
}

func generateTransferUrls() ([]string, [][]byte) {
	from := accounts[rand.Intn(len(accounts))]
	to := accounts[rand.Intn(len(accounts))]
	return []string{
			// "blcc://eth1.0/account/" + from + "/nonce",
			// "blcc://eth1.0/account/" + from + "/balance",
			// "blcc://eth1.0/account/" + to + "/balance",
			from + "/nonce",
			from + "/balance",
			to + "/balance",
		}, [][]byte{
			make([]byte, 4),
			make([]byte, 8),
			make([]byte, 8),
		}
}

func generateContractUrls() (string, []string, [][]byte) {
	address := randomHexString(20)
	return address, []string{
			// "blcc://eth1.0/account/" + address + "/",
			// "blcc://eth1.0/account/" + address + "/code",
			// "blcc://eth1.0/account/" + address + "/nonce",
			// "blcc://eth1.0/account/" + address + "/balance",
			// "blcc://eth1.0/account/" + address + "/defer",
			// "blcc://eth1.0/account/" + address + "/storage",
			// "blcc://eth1.0/account/" + address + "/storage/containers",
			// "blcc://eth1.0/account/" + address + "/storage/native/",
			// "blcc://eth1.0/account/" + address + "/storage/containers/!/",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[0] + "/",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[0] + "/!",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[0] + "/@",
			// "blcc://eth1.0/account/" + address + "/storage/containers/!/" + containers[0],
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[1] + "/",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[1] + "/!",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[1] + "/@",
			// "blcc://eth1.0/account/" + address + "/storage/containers/!/" + containers[1],
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[2] + "/",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[2] + "/!",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[2] + "/@",
			// "blcc://eth1.0/account/" + address + "/storage/containers/!/" + containers[2],
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[3] + "/",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[3] + "/!",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[3] + "/@",
			// "blcc://eth1.0/account/" + address + "/storage/containers/!/" + containers[3],
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[4] + "/",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[4] + "/!",
			// "blcc://eth1.0/account/" + address + "/storage/containers/" + containers[4] + "/@",
			// "blcc://eth1.0/account/" + address + "/storage/containers/!/" + containers[4],
			address + "/",
			address + "/code",
			address + "/nonce",
			address + "/balance",
			address + "/defer",
			address + "/storage",
			address + "/storage/containers",
			address + "/storage/native/",
			address + "/storage/containers/!/",
			address + "/storage/containers/" + containers[0] + "/",
			address + "/storage/containers/" + containers[0] + "/!",
			address + "/storage/containers/" + containers[0] + "/@",
			address + "/storage/containers/!/" + containers[0],
			address + "/storage/containers/" + containers[1] + "/",
			address + "/storage/containers/" + containers[1] + "/!",
			address + "/storage/containers/" + containers[1] + "/@",
			address + "/storage/containers/!/" + containers[1],
			address + "/storage/containers/" + containers[2] + "/",
			address + "/storage/containers/" + containers[2] + "/!",
			address + "/storage/containers/" + containers[2] + "/@",
			address + "/storage/containers/!/" + containers[2],
			address + "/storage/containers/" + containers[3] + "/",
			address + "/storage/containers/" + containers[3] + "/!",
			address + "/storage/containers/" + containers[3] + "/@",
			address + "/storage/containers/!/" + containers[3],
			address + "/storage/containers/" + containers[4] + "/",
			address + "/storage/containers/" + containers[4] + "/!",
			address + "/storage/containers/" + containers[4] + "/@",
			address + "/storage/containers/!/" + containers[4],
		}, [][]byte{
			make([]byte, 32),
			make([]byte, 4),
			make([]byte, 4),
			make([]byte, 8),
			make([]byte, 4),
			make([]byte, 4),
			make([]byte, 32),
			make([]byte, 4),
			make([]byte, 32),
			make([]byte, 32),
			make([]byte, 4),
			make([]byte, 8),
			make([]byte, 8),
			make([]byte, 32),
			make([]byte, 4),
			make([]byte, 8),
			make([]byte, 8),
			make([]byte, 32),
			make([]byte, 4),
			make([]byte, 8),
			make([]byte, 8),
			make([]byte, 32),
			make([]byte, 4),
			make([]byte, 8),
			make([]byte, 8),
			make([]byte, 32),
			make([]byte, 4),
			make([]byte, 8),
			make([]byte, 8),
		}
}

func generateAccountUrls() (string, []string, [][]byte) {
	address := randomHexString(20)
	return address, []string{
			// "blcc://eth1.0/account/" + address + "/",
			// "blcc://eth1.0/account/" + address + "/code",
			// "blcc://eth1.0/account/" + address + "/nonce",
			// "blcc://eth1.0/account/" + address + "/balance",
			// "blcc://eth1.0/account/" + address + "/defer",
			// "blcc://eth1.0/account/" + address + "/storage",
			// "blcc://eth1.0/account/" + address + "/storage/containers",
			// "blcc://eth1.0/account/" + address + "/storage/native/",
			// "blcc://eth1.0/account/" + address + "/storage/containers/!/",
			address + "/",
			address + "/code",
			address + "/nonce",
			address + "/balance",
			address + "/defer",
			address + "/storage",
			address + "/storage/containers",
			address + "/storage/native/",
			address + "/storage/containers/!/",
		}, [][]byte{
			make([]byte, 32),
			make([]byte, 4),
			make([]byte, 4),
			make([]byte, 8),
			make([]byte, 4),
			make([]byte, 4),
			make([]byte, 4),
			make([]byte, 4),
			make([]byte, 4),
		}
}

func randomHexString(nbytes int) string {
	b := make([]byte, nbytes)
	rnd.Read(b)
	return fmt.Sprintf("%x", b)
}
