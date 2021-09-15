package types

import (
	"math"

	ethCommon "github.com/arcology-network/3rd-party/eth/common"
	"github.com/arcology-network/common-lib/common"
	encoding "github.com/arcology-network/common-lib/encoding"
)

type ArbitratorRequest struct {
	TxsListGroup [][]*TxElement
}

type TxElement struct {
	TxHash  *ethCommon.Hash
	Batchid uint64
}

func (tx TxElement) Size() uint32 {
	return ethCommon.Hash{}.Size() + encoding.Uint64(0).Size()
}

type TxElements []*TxElement

func (elems TxElements) Encode() []byte {
	length := int(TxElement{}.Size())
	bytes := make([]byte, len(elems)*length)
	for i := range elems {
		copy(bytes[i*length:], elems[i].TxHash[:])
		copy(bytes[i*length+len(elems[i].TxHash):], encoding.Uint64(elems[i].Batchid).Encode())
	}
	return bytes
}

func (_ TxElements) Decode(bytes []byte) TxElements {
	length := int(TxElement{}.Size())
	elems := make([]*TxElement, len(bytes)/int(length))
	for i := 0; i < len(elems); i++ {
		tx := TxElement{&ethCommon.Hash{}, 0}
		copy(tx.TxHash[:], bytes[i*length:int(math.Min(float64((i+1)*length), float64(len(bytes))-1))])
		tx.Batchid = encoding.Uint64(0).Decode(bytes[i*length+ethCommon.HashLength : int(math.Min(float64((i+1)*length+ethCommon.HashLength), float64(len(bytes))-1))])
		elems[i] = &tx
	}
	return elems
}
func (request *ArbitratorRequest) GobEncode() ([]byte, error) {
	return request.Encode(), nil
}
func (request *ArbitratorRequest) GobDecode(data []byte) error {
	req := request.Decode(data)
	request.TxsListGroup = req.TxsListGroup
	return nil
}
func (request *ArbitratorRequest) Encode() []byte {
	bytes := make([][]byte, len(request.TxsListGroup))
	worker := func(start int, end int, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			bytes[i] = TxElements(request.TxsListGroup[i]).Encode()
		}
	}
	common.ParallelWorker(len(bytes), 2, worker)
	return encoding.Byteset(bytes).Encode()
}

func (_ ArbitratorRequest) Decode(bytes []byte) *ArbitratorRequest {
	byteset := encoding.Byteset{}.Decode(bytes)
	elems := make([][]*TxElement, len(byteset))

	worker := func(start int, end int, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			elems[i] = TxElements{}.Decode(byteset[i])
		}
	}
	common.ParallelWorker(len(elems), 2, worker)
	return &ArbitratorRequest{elems}
}

/*
func (a ArbitratorRequest) MarshalBinary() ([]byte, error) {
	ar := arEncoder.Encode(&a)
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(ar)
	return buf.Bytes(), err
}

func (a *ArbitratorRequest) UnmarshalBinary(data []byte) error {
	var ar arbReq
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&ar)
	if err != nil {
		return err
	}
	a.TxsListGroup = arDecoder.Decode(&ar).TxsListGroup
	return nil
}
*/
type arbReq struct {
	Indices []uint32
	Hashes  []byte
	Batches []uint32
}

type arbReqEncoder struct {
	indexBuf []uint32
	hashBuf  []byte
	batchBuf []uint32
}

func newArbReqEncoder() *arbReqEncoder {
	maxSize := 500000
	return &arbReqEncoder{
		indexBuf: make([]uint32, maxSize*2),
		hashBuf:  make([]byte, maxSize*2*32),
		batchBuf: make([]uint32, maxSize*2),
	}
}

func (e *arbReqEncoder) Encode(r *ArbitratorRequest) *arbReq {
	if len(r.TxsListGroup) == 0 {
		return &arbReq{}
	}

	indexOffset := uint32(0)
	dataOffset := 0
	batchOffset := uint32(0)

	prevGroupSize := len(r.TxsListGroup[0])
	count := 1
	for _, elem := range r.TxsListGroup[0] {
		dataOffset += copy(e.hashBuf[dataOffset:], elem.TxHash.Bytes())
		e.batchBuf[batchOffset] = uint32(elem.Batchid)
		batchOffset++
	}

	for i := 1; i < len(r.TxsListGroup); i++ {
		if len(r.TxsListGroup[i]) != prevGroupSize {
			e.indexBuf[indexOffset] = uint32(prevGroupSize)
			e.indexBuf[indexOffset+1] = uint32(count)
			indexOffset += 2
			prevGroupSize = len(r.TxsListGroup[i])
			count = 1
		} else {
			count++
		}

		for _, elem := range r.TxsListGroup[i] {
			dataOffset += copy(e.hashBuf[dataOffset:], elem.TxHash.Bytes())
			e.batchBuf[batchOffset] = uint32(elem.Batchid)
			batchOffset++
		}
	}

	e.indexBuf[indexOffset] = uint32(prevGroupSize)
	e.indexBuf[indexOffset+1] = uint32(count)
	indexOffset += 2

	return &arbReq{
		Indices: e.indexBuf[:indexOffset],
		Hashes:  e.hashBuf[:dataOffset],
		Batches: e.batchBuf[:batchOffset],
	}
}

type arbReqDecoder struct {
	list [][]*TxElement
}

func newArbReqDecoder() *arbReqDecoder {
	list := make([][]*TxElement, 500000)
	for i := range list {
		list[i] = make([]*TxElement, 0, 8)
	}
	return &arbReqDecoder{
		list: list,
	}
}

func (d *arbReqDecoder) Decode(r *arbReq) *ArbitratorRequest {
	offset := 0
	hashOffset := 0
	batchOffset := 0
	for i := 0; i < len(r.Indices); i += 2 {
		subListSize := r.Indices[i]
		count := r.Indices[i+1]
		for j := uint32(0); j < count; j++ {
			d.list[offset] = d.list[offset][:0]
			for k := uint32(0); k < subListSize; k++ {
				hash := ethCommon.BytesToHash(r.Hashes[hashOffset : hashOffset+32])
				d.list[offset] = append(d.list[offset], &TxElement{
					TxHash:  &hash,
					Batchid: uint64(r.Batches[batchOffset]),
				})
				hashOffset += 32
				batchOffset++
			}
			offset++
		}
	}
	return &ArbitratorRequest{
		TxsListGroup: d.list[:offset],
	}
}
