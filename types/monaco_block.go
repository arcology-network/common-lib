package types

import (

	// "github.com/arcology-network/common-lib/common"

	"github.com/arcology-network/common-lib/encoding"
)

const (
	AppType_Eth  = 0
	AppType_Coin = 1
)

type MonacoBlock struct {
	Height  uint64
	Headers [][]byte
	Txs     [][]byte
}

func (mb MonacoBlock) Hash() []byte {
	// bys := [][]byte{encoding.Byteset(mb.Headers).Flatten(), encoding.Byteset(mb.Txs).Flatten(), common.Uint64ToBytes(mb.Height)}
	// sum := sha256.Sum256(encoding.Byteset(bys).Flatten())
	// return sum[:]
	return []byte{}
}

func (mb MonacoBlock) GobEncode() ([]byte, error) {
	data := [][]byte{
		//	common.Uint64ToBytes(mb.Height),
		encoding.Byteset(mb.Headers).Encode(),
		encoding.Byteset(mb.Txs).Encode(),
	}
	return encoding.Byteset(data).Encode(), nil
}
func (mb *MonacoBlock) GobDecode(data []byte) error {
	fields := encoding.Byteset{}.Decode(data)
	// mb.Height = common.BytesToUint64(fields[0])
	mb.Headers = encoding.Byteset{}.Decode(fields[1])
	mb.Txs = encoding.Byteset{}.Decode(fields[2])
	return nil
}
