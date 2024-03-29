package types

import (
	"time"

	// "github.com/arcology-network/common-lib/common"
	encoding "github.com/arcology-network/common-lib/encoding"
)

type StatisticalInformation struct {
	Key      string
	Value    string
	TimeUsed time.Duration
}

func (si StatisticalInformation) EncodeToBytes() []byte {
	data := [][]byte{
		[]byte(si.Key),
		[]byte(si.Value),
		// common.Uint64ToBytes(common.Int64ToUint64(int64(si.TimeUsed))),
	}
	return encoding.Byteset(data).Encode()
}
func (si *StatisticalInformation) Decode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	si.Key = string(fields[0])
	si.Value = string(fields[1])
	// si.TimeUsed = time.Duration(common.Uint64ToInt64(common.BytesToUint64(fields[2])))
	return nil
}
