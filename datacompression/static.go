package datacompression

import (
	"strconv"

	common "github.com/arcology-network/common-lib/common"
	performance "github.com/arcology-network/common-lib/mhasher"
)

func (this *CompressionLut) GetNewAccounts(originals []string) []string {
	acctLen := 40
	prefixLen := len("blcc://eth1.0/account/")

	keys, _ := performance.UniqueStrings(GetBetween(originals, prefixLen, prefixLen+acctLen))
	return this.filterExistingKeys(keys, this.dict) // Get new keys
}

func (this *CompressionLut) CompressStaticKey(original string) string {
	acctLen := 40
	prefixLen := len("blcc://eth1.0/account/")

	if len(original) < prefixLen {
		return original
	}

	var prefixid int
	k := original[:prefixLen-1]
	if v, ok := this.dict.Get(k); ok {
		prefixid = int(v.(uint32))
	} else {
		return original
	}

	if len(original) < prefixLen+acctLen {
		original = "[" + strconv.Itoa(int(prefixid)) + "]" + original[prefixLen:]
		return original
	}

	key := original[prefixLen : prefixLen+acctLen]
	if id, ok := this.dict.Get(key); ok {
		original = "[" + strconv.Itoa(int(prefixid)) + "]" + "/[" + strconv.Itoa(int(id.(uint32)+this.offset)) + "]" + original[prefixLen+acctLen:]
	} else {
		if id, ok := this.tempLut.dict.Get(key); ok {
			original = "[" + strconv.Itoa(int(prefixid)) + "]" + "/[" + strconv.Itoa(int(id.(uint32)+this.length)) + "]" + original[prefixLen+acctLen:]
		}
	}
	return original
}

func (this *CompressionLut) CompressStaticKeys(originals []string) []string {
	replacer := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			originals[i] = this.CompressStaticKey(originals[i])
		}
	}
	common.ParallelWorker(len(originals), 4, replacer)
	return originals
}
