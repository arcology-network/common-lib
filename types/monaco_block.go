package types

import codec "github.com/arcology-network/common-lib/codec"

// "github.com/arcology-network/common-lib/common"

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
	// bys := [][]byte{codec.Byteset(mb.Headers).Flatten(), codec.Byteset(mb.Txs).Flatten(), common.Uint64ToBytes(mb.Height)}
	// sum := sha256.Sum256(codec.Byteset(bys).Flatten())
	// return sum[:]
	return []byte{}
}

func (mb MonacoBlock) GobEncode() ([]byte, error) {
	data := [][]byte{
		//	common.Uint64ToBytes(mb.Height),
		codec.Byteset(mb.Headers).Encode(),
		codec.Byteset(mb.Txs).Encode(),
	}
	return codec.Byteset(data).Encode(), nil
}
func (mb *MonacoBlock) GobDecode(data []byte) error {
	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
	// mb.Height = common.BytesToUint64(fields[0])
	mb.Headers = codec.Byteset{}.Decode(fields[1]).(codec.Byteset)
	mb.Txs = codec.Byteset{}.Decode(fields[2]).(codec.Byteset)
	return nil
}
