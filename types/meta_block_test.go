package types

import (
	"reflect"
	"testing"

	evmCommon "github.com/arcology-network/evm/common"
)

func TestMetaBlock(t *testing.T) {
	ethHash := evmCommon.BytesToHash([]byte{9, 9, 9, 9, 9, 9})
	in := &MetaBlock{
		Txs:      [][]byte{{1, 2}, {3, 4}, {5, 6}, {9, 8}, {7, 6}, {5, 4}, {4, 6}, {2, 7}, {8, 0}, {1, 9}},
		Hashlist: []*evmCommon.Hash{&ethHash, &ethHash, &ethHash, &ethHash, &ethHash, &ethHash, &ethHash, &ethHash, &ethHash, &ethHash},
	}
	buffer := in.Encode()

	out := new(MetaBlock)
	out.GobDecode(buffer)

	if !reflect.DeepEqual(in, out) {
		t.Error("Error")
	}
}
