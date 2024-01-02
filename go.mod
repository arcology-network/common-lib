module github.com/arcology-network/common-lib

go 1.21

replace github.com/arcology-network/common-lib => ../common-lib/

replace github.com/arcology-network/arcology-network/btree => ../btree/

replace github.com/ethereum/go-ethereum => ../concurrent-evm/

require (
	github.com/dgraph-io/badger v1.6.2
	github.com/elliotchance/orderedmap v1.5.0
	github.com/ethereum/go-ethereum v1.13.5
	github.com/holiman/uint256 v1.2.3
	github.com/natefinch/atomic v1.0.1
	github.com/stretchr/testify v1.8.4
	golang.org/x/crypto v0.14.0
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
)

// require github.com/arcology-network/3rd-party v1.7.1

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/bits-and-blooms/bitset v1.7.0 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/crate-crypto/go-kzg-4844 v0.7.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/dgraph-io/ristretto v0.0.2 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/ethereum-optimism/superchain-registry/superchain v0.0.0-20231030223232-e16eae11e492 // indirect
	github.com/ethereum/c-kzg-4844 v0.4.0 // indirect
	github.com/getsentry/sentry-go v0.18.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/supranational/blst v0.3.11 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)
