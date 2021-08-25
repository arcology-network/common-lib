package types

import (
	"encoding/binary"
	"math/big"

	"github.com/HPISTechnologies/common-lib/common"
	"github.com/HPISTechnologies/common-lib/encoding"
)

const (
	// Address in form of HexString.
	HexAddrLen = 42
	// Address in form of [20]byte.
	AddrLen = 20
)

// ArrayAccess presents an array read/write operation on array 'ID', position 'Index'.
// Index = -1 means append.
type ArrayAccess struct {
	ID    string
	Index int
}

type ArrayRead struct {
	ArrayAccess
	// Version Version
}

type ArrayWrite struct {
	ArrayAccess
	Version Version
	Value   []byte
}

type HashMapAccess struct {
	ID  string
	Key string
}

type HashMapRead struct {
	HashMapAccess
	// Version Version
}

type HashMapWrite struct {
	HashMapAccess
	Version Version
	Value   []byte
}

type DeferCallWrite struct {
	DeferID   string
	Signature string
}

const (
	// QueuePosSize used in ConcurrentQueue.size().
	QueuePosSize byte = 0
	// QueuePosHead used in ConcurrentQueue.pop().
	QueuePosHead byte = 1
	// QueuePosTail used in ConcurrentQueue.push().
	QueuePosTail byte = 2

	// QueueOpCreate used in ConcurrentQueue.create().
	QueueOpCreate byte = 0
	// QueueOpPush used in ConcurrentQueue.push().
	QueueOpPush byte = 1
	// QueueOpPop used in ConcurrentQueue.pop().
	QueueOpPop byte = 2
)

// QueueRead used in ConcurrentQueue.size().
type QueueRead struct {
	ID  string
	Pos byte
}

// QueueWrite used in ConcurrentQueue.create()/push()/pop().
// create: Op = QueueOpCreate, Value = []byte{ byte(elemType) };
// push: Pos = QueuePosTail, Op = QueueOpPush, Value = ToBytes(elem);
// pop: Pos = QueuePosHead, Op = QueueOpPop.
type QueueWrite struct {
	ID  string
	Pos byte
	Op  byte
	// Value used in ConcurrentQueue.create()/push().
	Value []byte
}

type ReadSet struct {
	ArrayReads   []ArrayRead
	HashMapReads []HashMapRead
	QueueReads   []QueueRead
}

type Reads struct {
	ClibReads       map[Address]*ReadSet
	BalanceReads    map[Address]*big.Int
	EthStorageReads []Address
}

func (r Reads) MarshalBinary() ([]byte, error) {
	offset := 0
	buf := <-bytesPool
	buf = buf[:0]
	defer func() {
		bytesPool <- buf
	}()
	// Length of ClibReads, BalanceReads, EthStorageReads.
	offset = expandUint32(buf, offset, uint32(len(r.ClibReads)))
	offset = expandUint32(buf, offset, uint32(len(r.BalanceReads)))
	offset = expandUint32(buf, offset, uint32(len(r.EthStorageReads)))
	// ClibReads.
	for addr, readSet := range r.ClibReads {
		// Address.
		offset = expandHexAddress(buf, offset, addr)
		// ReadSet.
		// Length of ArrayReads, HashMapReads, QueueReads.
		offset = expandUint32(buf, offset, uint32(len(readSet.ArrayReads)))
		offset = expandUint32(buf, offset, uint32(len(readSet.HashMapReads)))
		offset = expandUint32(buf, offset, uint32(len(readSet.QueueReads)))
		// ArrayReads.
		for _, arrRead := range readSet.ArrayReads {
			// Length of ID.
			offset = expandByte(buf, offset, byte(len(arrRead.ID)))
			// ID.
			offset = expandString(buf, offset, arrRead.ID)
			// Index.
			offset = expandUint32(buf, offset, uint32(arrRead.Index))
		}
		// HashMapReads.
		for _, hmRead := range readSet.HashMapReads {
			// Length of ID, Key.
			offset = expandByte(buf, offset, byte(len(hmRead.ID)))
			offset = expandByte(buf, offset, byte(len(hmRead.Key)))
			// ID.
			offset = expandString(buf, offset, hmRead.ID)
			// Key.
			offset = expandString(buf, offset, hmRead.Key)
		}
		// QueueReads.
		for _, qRead := range readSet.QueueReads {
			// Length of ID.
			offset = expandByte(buf, offset, byte(len(qRead.ID)))
			// ID.
			offset = expandString(buf, offset, qRead.ID)
			// Pos.
			offset = expandByte(buf, offset, qRead.Pos)
		}
	}
	// BalanceReads.
	offset = encodeAddrBigIntMap(buf, offset, r.BalanceReads)
	// EthStorageReads.
	for _, addr := range r.EthStorageReads {
		offset = expandAddress(buf, offset, addr)
	}
	returnBuf := make([]byte, offset)
	copy(returnBuf[0:], buf[:offset])
	return returnBuf, nil
}

func (r *Reads) UnmarshalBinary(buf []byte) error {
	offset := uint32(0)
	// Length of ClibReads, BalanceReads, EthStorageReads.
	clibReadsLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	balanceReadsLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	ethStorageReadsLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	// ClibReads.
	r.ClibReads = make(map[Address]*ReadSet)
	for i := uint32(0); i < clibReadsLen; i++ {
		// Address.
		addr := Address(string(buf[offset : offset+HexAddrLen]))
		offset += HexAddrLen
		// ReadSet.
		// Length of ArrayReads, HashMapReads, QueueReads.
		arrayReadsLen := binary.BigEndian.Uint32(buf[offset : offset+4])
		offset += 4
		hashMapReadsLen := binary.BigEndian.Uint32(buf[offset : offset+4])
		offset += 4
		queueReadsLen := binary.BigEndian.Uint32(buf[offset : offset+4])
		offset += 4
		readSet := &ReadSet{
			ArrayReads:   make([]ArrayRead, arrayReadsLen),
			HashMapReads: make([]HashMapRead, hashMapReadsLen),
			QueueReads:   make([]QueueRead, queueReadsLen),
		}
		// ArrayReads.
		for j := uint32(0); j < arrayReadsLen; j++ {
			// Length of ID.
			idLen := uint32(buf[offset])
			offset++
			readSet.ArrayReads[j] = ArrayRead{
				ArrayAccess: ArrayAccess{
					ID:    string(buf[offset : offset+idLen]),
					Index: int(binary.BigEndian.Uint32(buf[offset+idLen : offset+idLen+4])),
				},
			}
			offset += idLen + 4
		}
		// HashMapReads.
		for j := uint32(0); j < hashMapReadsLen; j++ {
			// Length of ID, Key.
			idLen := uint32(buf[offset])
			offset++
			keyLen := uint32(buf[offset])
			offset++
			readSet.HashMapReads[j] = HashMapRead{
				HashMapAccess: HashMapAccess{
					ID:  string(buf[offset : offset+idLen]),
					Key: string(buf[offset+idLen : offset+idLen+keyLen]),
				},
			}
			offset += idLen + keyLen
		}
		// QueueReads.
		for j := uint32(0); j < queueReadsLen; j++ {
			// Length of ID.
			idLen := uint32(buf[offset])
			offset++
			readSet.QueueReads[j] = QueueRead{
				ID:  string(buf[offset : offset+idLen]),
				Pos: buf[offset+idLen],
			}
			offset += idLen + 1
		}
		r.ClibReads[addr] = readSet
	}
	// BalanceReads.
	r.BalanceReads = make(map[Address]*big.Int)
	for i := uint32(0); i < balanceReadsLen; i++ {
		offset = decodeAddrBigIntMap(buf, offset, r.BalanceReads)
	}
	// EthStorageReads.
	r.EthStorageReads = make([]Address, ethStorageReadsLen)
	for i := uint32(0); i < ethStorageReadsLen; i++ {
		r.EthStorageReads[i] = Address(string(buf[offset : offset+AddrLen]))
		offset += AddrLen
	}

	return nil
}

type WriteSet struct {
	ArrayWrites     []ArrayWrite
	HashMapWrites   []HashMapWrite
	QueueWrites     []QueueWrite
	DeferCallWrites []DeferCallWrite
}

type Writes struct {
	ClibWrites       map[Address]*WriteSet
	NewAccounts      []Address
	BalanceWrites    map[Address]*big.Int
	BalanceOrigin    map[Address]*big.Int
	NonceWrites      map[Address]uint64
	CodeWrites       map[Address][]byte
	EthStorageWrites map[Address]map[string]string
}

func (w Writes) MarshalBinary() ([]byte, error) {
	offset := 0
	buf := <-bytesPool
	buf = buf[:0]
	defer func() {
		bytesPool <- buf
	}()
	// Length of ClibWrites, NewAccounts, BalanceWrites, BalanceOrigin, NonceWrites, CodeWrites, EthStorageWrites.
	offset = expandUint32(buf, offset, uint32(len(w.ClibWrites)))
	offset = expandUint32(buf, offset, uint32(len(w.NewAccounts)))
	offset = expandUint32(buf, offset, uint32(len(w.BalanceWrites)))
	offset = expandUint32(buf, offset, uint32(len(w.BalanceOrigin)))
	offset = expandUint32(buf, offset, uint32(len(w.NonceWrites)))
	offset = expandUint32(buf, offset, uint32(len(w.CodeWrites)))
	offset = expandUint32(buf, offset, uint32(len(w.EthStorageWrites)))
	// ClibWrites.
	for addr, writeSet := range w.ClibWrites {
		// Address.
		offset = expandHexAddress(buf, offset, addr)
		// WriteSet.
		// Length of ArrayWrites, HashMapWrites, QueueWrites, DeferCallWrites.
		offset = expandUint32(buf, offset, uint32(len(writeSet.ArrayWrites)))
		offset = expandUint32(buf, offset, uint32(len(writeSet.HashMapWrites)))
		offset = expandUint32(buf, offset, uint32(len(writeSet.QueueWrites)))
		offset = expandUint32(buf, offset, uint32(len(writeSet.DeferCallWrites)))
		// ArrayWrites.
		for _, arrWrite := range writeSet.ArrayWrites {
			// Length of ID, Value.
			offset = expandByte(buf, offset, byte(len(arrWrite.ID)))
			offset = expandUint32(buf, offset, uint32(len(arrWrite.Value)))
			// ID.
			offset = expandString(buf, offset, arrWrite.ID)
			// Index.
			offset = expandUint32(buf, offset, uint32(arrWrite.Index))
			// Version.
			offset = expandUint64(buf, offset, uint64(arrWrite.Version))
			// Value.
			offset = expandByteArray(buf, offset, arrWrite.Value)
		}
		// HashMapWrites.
		for _, hmWrite := range writeSet.HashMapWrites {
			// Length of ID, Key, Value.
			offset = expandByte(buf, offset, byte(len(hmWrite.ID)))
			offset = expandByte(buf, offset, byte(len(hmWrite.Key)))
			offset = expandUint32(buf, offset, uint32(len(hmWrite.Value)))
			// ID.
			offset = expandString(buf, offset, hmWrite.ID)
			// Key.
			offset = expandString(buf, offset, hmWrite.Key)
			// Version.
			offset = expandUint64(buf, offset, uint64(hmWrite.Version))
			// Value.
			offset = expandByteArray(buf, offset, hmWrite.Value)
		}
		// QueueWrites.
		for _, qWrite := range writeSet.QueueWrites {
			// Length of ID, Value.
			offset = expandByte(buf, offset, byte(len(qWrite.ID)))
			offset = expandUint32(buf, offset, uint32(len(qWrite.Value)))
			// ID.
			offset = expandString(buf, offset, qWrite.ID)
			// Pos.
			offset = expandByte(buf, offset, qWrite.Pos)
			// Op.
			offset = expandByte(buf, offset, qWrite.Op)
			// Value.
			offset = expandByteArray(buf, offset, qWrite.Value)
		}
		// DeferCallWrites.
		for _, dcWrite := range writeSet.DeferCallWrites {
			// Length of DeferID, Signature.
			offset = expandByte(buf, offset, byte(len(dcWrite.DeferID)))
			offset = expandByte(buf, offset, byte(len(dcWrite.Signature)))
			// DeferID.
			offset = expandString(buf, offset, dcWrite.DeferID)
			// Signature.
			offset = expandString(buf, offset, dcWrite.Signature)
		}
	}
	// NewAccounts.
	for _, acc := range w.NewAccounts {
		offset = expandAddress(buf, offset, acc)
	}
	// BalanceWrites.
	offset = encodeAddrBigIntMap(buf, offset, w.BalanceWrites)
	// BalanceOrigin.
	offset = encodeAddrBigIntMap(buf, offset, w.BalanceOrigin)
	// NonceWrites.
	for addr := range w.NonceWrites {
		offset = expandAddress(buf, offset, addr)
	}
	// CodeWrites.
	for addr, code := range w.CodeWrites {
		// Length of Code.
		offset = expandUint32(buf, offset, uint32(len(code)))
		// Address.
		offset = expandAddress(buf, offset, addr)
		// Code.
		offset = expandByteArray(buf, offset, code)
	}
	// EthStorageWrites.
	for addr, storage := range w.EthStorageWrites {
		// Length of Storage.
		offset = expandUint32(buf, offset, uint32(len(storage)))
		// Address.
		offset = expandAddress(buf, offset, addr)
		// Storage.
		for k, v := range storage {
			offset = expandByteArray(buf, offset, []byte(k))
			offset = expandByteArray(buf, offset, []byte(v))
		}
	}
	returnBuf := make([]byte, offset)
	copy(returnBuf[0:], buf[:offset])
	return returnBuf, nil
}

func (w *Writes) UnmarshalBinary(buf []byte) error {
	offset := uint32(0)
	// Length of ClibWrites, NewAccounts, BalanceWrites, BalanceOrigin, NonceWrites, CodeWrites, EthStorageWrites.
	clibWritesLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	newAccountsLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	balanceWritesLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	balanceOriginLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	nonceWritesLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	codeWritesLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	ethStorageWritesLen := binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	// ClibWrites.
	w.ClibWrites = make(map[Address]*WriteSet)
	for i := uint32(0); i < clibWritesLen; i++ {
		// Address.
		addr := Address(string(buf[offset : offset+HexAddrLen]))
		offset += HexAddrLen
		// WriteSet.
		// Length of ArrayWrites, HashMapWrites, QueueWrites, DeferCallWrites.
		arrayWritesLen := binary.BigEndian.Uint32(buf[offset : offset+4])
		offset += 4
		hashMapWritesLen := binary.BigEndian.Uint32(buf[offset : offset+4])
		offset += 4
		queueWritesLen := binary.BigEndian.Uint32(buf[offset : offset+4])
		offset += 4
		deferCallWritesLen := binary.BigEndian.Uint32(buf[offset : offset+4])
		offset += 4
		writeSet := &WriteSet{
			ArrayWrites:     make([]ArrayWrite, arrayWritesLen),
			HashMapWrites:   make([]HashMapWrite, hashMapWritesLen),
			QueueWrites:     make([]QueueWrite, queueWritesLen),
			DeferCallWrites: make([]DeferCallWrite, deferCallWritesLen),
		}
		// ArrayWrites.
		for j := uint32(0); j < arrayWritesLen; j++ {
			// Length of ID, Value.
			idLen := uint32(buf[offset])
			offset++
			valueLen := binary.BigEndian.Uint32(buf[offset : offset+4])
			offset += 4
			// ArrayWrite.
			writeSet.ArrayWrites[j] = ArrayWrite{
				ArrayAccess: ArrayAccess{
					ID:    string(buf[offset : offset+idLen]),
					Index: int(binary.BigEndian.Uint32(buf[offset+idLen : offset+idLen+4])),
				},
				Version: Version(binary.BigEndian.Uint64(buf[offset+idLen+4 : offset+idLen+4+8])),
				Value:   buf[offset+idLen+4+8 : offset+idLen+4+8+valueLen],
			}
			offset += idLen + 4 + 8 + valueLen
		}
		// HashMapWrites.
		for j := uint32(0); j < hashMapWritesLen; j++ {
			// Length of ID, Key, Value.
			idLen := uint32(buf[offset])
			offset++
			keyLen := uint32(buf[offset])
			offset++
			valueLen := binary.BigEndian.Uint32(buf[offset : offset+4])
			offset += 4
			// HashMapWrite.
			writeSet.HashMapWrites[j] = HashMapWrite{
				HashMapAccess: HashMapAccess{
					ID:  string(buf[offset : offset+idLen]),
					Key: string(buf[offset+idLen : offset+idLen+keyLen]),
				},
				Version: Version(binary.BigEndian.Uint64(buf[offset+idLen+keyLen : offset+idLen+keyLen+8])),
				Value:   buf[offset+idLen+keyLen+8 : offset+idLen+keyLen+8+valueLen],
			}
			offset += idLen + keyLen + 8 + valueLen
		}
		// QueueWrites.
		for j := uint32(0); j < queueWritesLen; j++ {
			// Length of ID, Value.
			idLen := uint32(buf[offset])
			offset++
			valueLen := binary.BigEndian.Uint32(buf[offset : offset+4])
			offset += 4
			// QueueWrite.
			writeSet.QueueWrites[j] = QueueWrite{
				ID:    string(buf[offset : offset+idLen]),
				Pos:   buf[offset+idLen],
				Op:    buf[offset+idLen+1],
				Value: buf[offset+idLen+1+1 : offset+idLen+1+1+valueLen],
			}
			offset += idLen + 1 + 1 + valueLen
		}
		// DeferCallWrites.
		for j := uint32(0); j < deferCallWritesLen; j++ {
			// Length of DeferID, Signature.
			deferIDLen := uint32(buf[offset])
			offset++
			signatureLen := uint32(buf[offset])
			offset++
			// DeferCallWrite.
			writeSet.DeferCallWrites[j] = DeferCallWrite{
				DeferID:   string(buf[offset : offset+deferIDLen]),
				Signature: string(buf[offset+deferIDLen : offset+deferIDLen+signatureLen]),
			}
			offset += deferIDLen + signatureLen
		}
		w.ClibWrites[addr] = writeSet
	}
	// NewAccounts.
	w.NewAccounts = make([]Address, newAccountsLen)
	for i := uint32(0); i < newAccountsLen; i++ {
		w.NewAccounts[i] = Address(string(buf[offset : offset+AddrLen]))
		offset += AddrLen
	}
	// BalanceWrites.
	w.BalanceWrites = make(map[Address]*big.Int)
	for i := uint32(0); i < balanceWritesLen; i++ {
		offset = decodeAddrBigIntMap(buf, offset, w.BalanceWrites)
	}
	// BalanceOrigin.
	w.BalanceOrigin = make(map[Address]*big.Int)
	for i := uint32(0); i < balanceOriginLen; i++ {
		offset = decodeAddrBigIntMap(buf, offset, w.BalanceOrigin)
	}
	// NonceWrites.
	w.NonceWrites = make(map[Address]uint64)
	for i := uint32(0); i < nonceWritesLen; i++ {
		w.NonceWrites[Address(string(buf[offset:offset+AddrLen]))] = 1
		offset += AddrLen
	}
	// CodeWrites.
	w.CodeWrites = make(map[Address][]byte)
	for i := uint32(0); i < codeWritesLen; i++ {
		// Length of Code.
		codeLen := binary.BigEndian.Uint32(buf[offset : offset+4])
		offset += 4
		w.CodeWrites[Address(string(buf[offset:offset+AddrLen]))] = buf[offset+AddrLen : offset+AddrLen+codeLen]
		offset += AddrLen + codeLen
	}
	// EthStorageWrites.
	w.EthStorageWrites = make(map[Address]map[string]string)
	for i := uint32(0); i < ethStorageWritesLen; i++ {
		// Length of Storage.
		storageLen := binary.BigEndian.Uint32(buf[offset : offset+4])
		offset += 4
		// Address.
		addr := Address(string(buf[offset : offset+AddrLen]))
		offset += AddrLen
		w.EthStorageWrites[addr] = make(map[string]string)
		for j := uint32(0); j < storageLen; j++ {
			w.EthStorageWrites[addr][string(buf[offset:offset+32])] = string(buf[offset+32 : offset+32+32])
			offset += 32 + 32
		}
	}

	return nil
}

// type DeferCall struct {
// 	DeferID         string
// 	ContractAddress Address
// 	Signature       string
// }

// func (d DeferCall) MarshalBinary1() ([]byte, error) {
// 	offset := 0
// 	buf := make([]byte, 128)
// 	// Length of DeferID, Signature.
// 	buf[offset] = byte(len(d.DeferID))
// 	offset++
// 	buf[offset] = byte(len(d.Signature))
// 	offset++
// 	// DeferID.
// 	copy(buf[offset:], []byte(d.DeferID))
// 	offset += len(d.DeferID)
// 	// ContractAddress.
// 	copy(buf[offset:], []byte(d.ContractAddress))
// 	offset += AddrLen
// 	// Signature.
// 	copy(buf[offset:], []byte(d.Signature))
// 	offset += len(d.Signature)

// 	return buf[:offset], nil
// }

// func (d *DeferCall) UnmarshalBinary1(buf []byte) error {
// 	return nil
// }

// type EuResult struct {
// 	H       string
// 	R       *Reads
// 	W       *Writes
// 	DC      *DeferCall
// 	Status  uint64
// 	GasUsed uint64
// 	// RevertedTxs []string
// }

func encodeAddrBigIntMap(buf []byte, offset int, m map[Address]*big.Int) int {
	for addr, bi := range m {
		// Address.
		offset = expandAddress(buf, offset, addr)
		// Balance.
		sign := bi.Sign()
		offset = expandByte(buf, offset, byte(sign))
		if sign != 0 {
			// Length of Bits.
			offset = expandByte(buf, offset, byte(len(bi.Bits())))
			// Bits.
			for _, w := range bi.Bits() {
				offset = expandUint64(buf, offset, uint64(w))
			}
		}
	}
	return offset
}

func decodeAddrBigIntMap(buf []byte, offset uint32, m map[Address]*big.Int) uint32 {
	addr := Address(string(buf[offset : offset+AddrLen]))
	offset += AddrLen
	sign := buf[offset]
	offset++
	bi := new(big.Int)
	if sign != 0 {
		bitsLen := uint32(buf[offset])
		offset++
		bits := make([]big.Word, bitsLen)
		for j := uint32(0); j < bitsLen; j++ {
			bits[j] = big.Word(binary.BigEndian.Uint64(buf[offset : offset+8]))
			offset += 8
		}
		bi.SetBits(bits)
		if sign != 1 {
			bi = new(big.Int).Neg(bi)
		}
	}
	m[addr] = bi
	return offset
}

func expandByte(buf []byte, offset int, b byte) int {
	buf = buf[0 : offset+1]
	buf[offset] = b
	return offset + 1
}

func expandUint32(buf []byte, offset int, elem uint32) int {
	buf = buf[0 : offset+4]
	binary.BigEndian.PutUint32(buf[offset:], elem)
	return offset + 4
}

func expandUint64(buf []byte, offset int, elem uint64) int {
	buf = buf[0 : offset+8]
	binary.BigEndian.PutUint64(buf[offset:], elem)
	return offset + 8
}

func expandByteArray(buf []byte, offset int, bytes []byte) int {
	buf = buf[0 : offset+len(bytes)]
	copy(buf[offset:], bytes)
	return offset + len(bytes)
}

func expandString(buf []byte, offset int, s string) int {
	buf = buf[0 : offset+len(s)]
	copy(buf[offset:], []byte(s))
	return offset + len(s)
}

func expandAddress(buf []byte, offset int, addr Address) int {
	buf = buf[0 : offset+AddrLen]
	copy(buf[offset:], []byte(addr))
	return offset + AddrLen
}

func expandHexAddress(buf []byte, offset int, addr Address) int {
	buf = buf[0 : offset+HexAddrLen]
	copy(buf[offset:], []byte(addr))
	return offset + HexAddrLen
}

type DeferCalls []*DeferCall

func (dcs DeferCalls) Encode() []byte {
	if dcs == nil {
		return []byte{}
	}

	worker := func(start, end int, idx int, args ...interface{}) {
		defcalls := args[0].([]interface{})[0].(DeferCalls)
		dataSet := args[0].([]interface{})[1].([][]byte)

		for i := start; i < end; i++ {
			if defcall := defcalls[i]; defcall != nil {
				dataSet[i] = encoding.Byteset([][]byte{
					encoding.Uint32(len(defcall.DeferID[:])).Encode()[:],
					encoding.Uint32(len(defcall.Signature[:])).Encode()[:],
					[]byte(defcall.DeferID),
					[]byte(defcall.ContractAddress),
					[]byte(defcall.Signature),
				}).Flatten()
			}
		}
		/*
			for i := range defcalls[start:end] {
				idx := i + start
				defcall := defcalls[idx]
				if defcall == nil {
					continue
				}
				DeferIDLength := len(defcall.DeferID[:])
				SignatureLength := len(defcall.Signature[:])
				length := DeferIDLength + AddressLength + len(defcall.Signature[:]) + encoding.UINT32_LEN*2
				msgData := make([]byte, 0, length)

				msgData = append(msgData, encoding.Uint32(DeferIDLength).Encode()...)
				msgData = append(msgData, encoding.Uint32(SignatureLength).Encode()...)
				msgData = append(msgData, defcall.DeferID[:]...)
				msgData = append(msgData, string(defcall.ContractAddress)[:]...)
				msgData = append(msgData, defcall.Signature[:]...)
				datas[idx] = msgData
			}
		*/
	}

	dataSet := make([][]byte, len(dcs))
	common.ParallelWorker(len(dcs), concurrency, worker, dcs, dataSet)

	return encoding.Byteset(dataSet).Encode()
}

func (dcs *DeferCalls) Decode(data []byte) []*DeferCall {
	fields := encoding.Byteset{}.Decode(data)
	defs := make([]*DeferCall, len(fields))

	worker := func(start, end, idx int, args ...interface{}) {
		dataSet := args[0].([]interface{})[0].([][]byte)
		defcalls := args[0].([]interface{})[1].([]*DeferCall)

		for i := start; i < end; i++ {
			if len(dataSet[i]) == 0 {
				continue
			}
			deferCall := new(DeferCall)
			DeferIDLength := 0
			DeferIDLength = int(encoding.Uint32(DeferIDLength).Decode(dataSet[i][0:encoding.UINT32_LEN]))
			SignatureLength := 0
			SignatureLength = int(encoding.Uint32(SignatureLength).Decode(dataSet[i][encoding.UINT32_LEN : encoding.UINT32_LEN*2]))

			deferCall.DeferID = string(dataSet[i][encoding.UINT32_LEN*2 : encoding.UINT32_LEN*2+DeferIDLength])
			deferCall.ContractAddress = Address(string(dataSet[i][encoding.UINT32_LEN*2+DeferIDLength : encoding.UINT32_LEN*2+DeferIDLength+AddressLength]))
			deferCall.Signature = string(dataSet[i][encoding.UINT32_LEN*2+DeferIDLength+AddressLength:])
			defcalls[i] = deferCall
		}
	}
	common.ParallelWorker(len(fields), concurrency, worker, fields, defs)

	return defs
}
