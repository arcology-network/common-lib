package types

// import (
// 	"bytes"
// 	"encoding/gob"
// 	"fmt"
// 	"math/big"
// 	"sync"
// 	"testing"

// 	ethCommon "github.com/arcology/3rd-party/eth/common"
// 	"github.com/arcology/common-lib/common"
// )

// func Test_Defcalls_encodingAndDeconing(t *testing.T) {
// 	defs := []*DeferCall{
// 		{
// 			DeferID:         "123123",
// 			ContractAddress: Address("defcalll123122121212"),
// 			Signature:       "defcall call()",
// 		},
// 		{
// 			DeferID:         "345678",
// 			ContractAddress: Address("defcalll123122121452"),
// 			Signature:       "defcall call(s int)",
// 		},
// 	}

// 	defcalls := make([]*DeferCall, 3)
// 	defcalls[0] = defs[0]

// 	defcalls[2] = defs[1]

// 	deferCalls := DeferCalls(defcalls)
// 	data := deferCalls.Encode()
// 	fmt.Printf("DeferCalls encode result=%v\n", data)

// 	deferCalls2 := new(DeferCalls)

// 	defss := deferCalls2.Decode(data)

// 	for _, v := range defss {
// 		fmt.Printf("DeferCalls dncode result=%v\n", v)

// 	}
// }

// // func Test_Defcall_encodingAndDeconing(t *testing.T) {
// var euResult *EuResult
// var reads *Reads
// var writes *Writes

// func setup() {
// 	addr1 := ethCommon.BytesToAddress([]byte{1, 2, 3})
// 	addr2 := ethCommon.BytesToAddress([]byte{4, 5, 6})
// 	addr3 := ethCommon.BytesToAddress([]byte{7, 8, 9})
// 	reads = &Reads{
// 		ClibReads: map[Address]*ReadSet{
// 			Address(string(addr2.Bytes())): {
// 				HashMapReads: []HashMapRead{
// 					{
// 						HashMapAccess: HashMapAccess{
// 							ID:  "hashmap",
// 							Key: "key",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		BalanceReads: map[Address]*big.Int{
// 			Address(string(addr1.Bytes())): new(big.Int).SetInt64(1),
// 			Address(string(addr2.Bytes())): new(big.Int).SetInt64(1),
// 			Address(string(addr3.Bytes())): new(big.Int).SetInt64(1),
// 		},
// 		EthStorageReads: nil,
// 	}
// 	writes = &Writes{
// 		ClibWrites: map[Address]*WriteSet{
// 			Address(string(addr2.Bytes())): {
// 				DeferCallWrites: []DeferCallWrite{
// 					{
// 						DeferID:   "deferid",
// 						Signature: "defer(string)",
// 					},
// 				},
// 			},
// 		},
// 		NewAccounts: nil,
// 		BalanceWrites: map[Address]*big.Int{
// 			Address(string(addr1.Bytes())): new(big.Int).SetInt64(-1),
// 			Address(string(addr3.Bytes())): new(big.Int).SetInt64(1),
// 		},
// 		BalanceOrigin: map[Address]*big.Int{
// 			Address(string(addr1.Bytes())): new(big.Int).SetInt64(100),
// 			Address(string(addr3.Bytes())): new(big.Int).SetInt64(100),
// 		},
// 		NonceWrites: map[Address]uint64{
// 			Address(string(addr1.Bytes())): 1,
// 		},
// 		CodeWrites:       nil,
// 		EthStorageWrites: nil,
// 	}
// 	euResult = &EuResult{
// 		H: string(ethCommon.Hash{}.Bytes()),
// 		R: reads,
// 		W: writes,
// 		DC: &DeferCall{
// 			DeferID:         "deferid",
// 			ContractAddress: Address(string(addr2.Bytes())),
// 			Signature:       "defer(string)",
// 		},
// 		Status:  1,
// 		GasUsed: 100,
// 		// RevertedTxs: nil,
// 	}
// }

// func Test_WriteSet(t *testing.T) {
// 	addr1 := ethCommon.BytesToAddress([]byte{1, 2, 3})
// 	addr2 := ethCommon.BytesToAddress([]byte{4, 5, 6})
// 	addr3 := ethCommon.BytesToAddress([]byte{7, 8, 9})
// 	writes = &Writes{
// 		ClibWrites: map[Address]*WriteSet{
// 			Address(string(addr2.Bytes())): {
// 				DeferCallWrites: []DeferCallWrite{
// 					{
// 						DeferID:   "deferid",
// 						Signature: "defer(string)",
// 					},
// 				},
// 			},
// 		},
// 		NewAccounts: nil,
// 		BalanceWrites: map[Address]*big.Int{
// 			Address(string(addr1.Bytes())): new(big.Int).SetInt64(-1),
// 			Address(string(addr3.Bytes())): new(big.Int).SetInt64(1),
// 		},
// 		BalanceOrigin: map[Address]*big.Int{
// 			Address(string(addr1.Bytes())): new(big.Int).SetInt64(100),
// 			Address(string(addr3.Bytes())): new(big.Int).SetInt64(100),
// 		},
// 		NonceWrites: map[Address]uint64{
// 			Address(string(addr1.Bytes())): 1,
// 			//Address(string(addr1.Bytes())): 2,
// 		},
// 		CodeWrites:       nil,
// 		EthStorageWrites: nil,
// 	}

// 	data, err := writes.MarshalBinary()
// 	if err != nil {
// 		fmt.Printf("writes.MarshalBinary result err = %v\n", err)
// 		return
// 	}
// 	fmt.Printf("writes.MarshalBinary result  = %v\n", data)
// 	destWrites := &Writes{}
// 	err = destWrites.UnmarshalBinary(data)
// 	if err != nil {
// 		fmt.Printf("writes.UnmarshalBinary result err = %v\n", err)
// 		return
// 	}
// 	fmt.Printf("writes.UnmarshalBinary result  = %v\n", destWrites)
// }

// func Test_Gob(t *testing.T) {
// 	aw := ArrayWrite{
// 		ArrayAccess: ArrayAccess{
// 			ID:    "id111212",
// 			Index: 2,
// 		},
// 		Value:   []byte("skdlklsdgf"),
// 		Version: 0,
// 	}
// 	fmt.Printf("aw=%v\n", aw)
// 	data, err := common.GobEncode(aw)
// 	if err != nil {
// 		fmt.Printf("err=%v\n", err)
// 		return
// 	}
// 	fmt.Printf("data=%v\n", data)
// 	var src *ArrayWrite
// 	err = common.GobDecode(data, &src)
// 	if err != nil {
// 		fmt.Printf("err=%v\n", err)
// 		return
// 	}
// 	fmt.Printf("src=%v\n", src)
// }

// func BenchmarkEuResultEncode(b *testing.B) {
// 	setup()
// 	list := make([]*EuResult, 500)
// 	for i := range list {
// 		list[i] = euResult
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// for j := 0; j < 1000; j++ {
// 		var buf bytes.Buffer
// 		encoder := gob.NewEncoder(&buf)
// 		encoder.Encode(list)
// 		// }
// 	}
// }

// func BenchmarkParallelEuResultEncode(b *testing.B) {
// 	setup()
// 	list := make([]*EuResult, 500)
// 	for i := range list {
// 		list[i] = euResult
// 	}

// 	b.ResetTimer()
// 	concurrency := 4
// 	for i := 0; i < b.N; i++ {
// 		var wg sync.WaitGroup
// 		bufArray := make([][]byte, concurrency)
// 		for j := 0; j < concurrency; j++ {
// 			wg.Add(1)
// 			go func(index int) {
// 				var buf bytes.Buffer
// 				gob.NewEncoder(&buf).Encode(list[500*index/concurrency : 500*(index+1)/concurrency])
// 				bufArray[index] = buf.Bytes()
// 				wg.Done()
// 			}(j)
// 		}
// 		wg.Wait()
// 		var buf bytes.Buffer
// 		gob.NewEncoder(&buf).Encode(bufArray)
// 	}
// }

// func BenchmarkEuResultDecode(b *testing.B) {
// 	setup()
// 	list := make([]*EuResult, 500)
// 	for i := range list {
// 		list[i] = euResult
// 	}

// 	var buf bytes.Buffer
// 	encoder := gob.NewEncoder(&buf)
// 	encoder.Encode(list)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// for j := 0; j < 1000; j++ {
// 		var list []*EuResult
// 		gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&list)
// 		// }
// 	}
// }

// func BenchmarkEncodeBigInt(b *testing.B) {
// 	bigs := []*big.Int{
// 		new(big.Int).SetInt64(1),
// 		new(big.Int).SetInt64(1),
// 		new(big.Int).SetInt64(1),
// 		new(big.Int).SetInt64(1),
// 		new(big.Int).SetInt64(-1),
// 		new(big.Int).SetInt64(100),
// 		new(big.Int).SetInt64(100),
// 	}
// 	list := make([][]*big.Int, 500)
// 	for i := range list {
// 		list[i] = bigs
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// for j := 0; j < 1000; j++ {
// 		var buf bytes.Buffer
// 		encoder := gob.NewEncoder(&buf)
// 		encoder.Encode(list)
// 		// }
// 	}
// }

// func BenchmarkDecodeBigInt(b *testing.B) {
// 	bigs := []*big.Int{
// 		new(big.Int).SetInt64(1),
// 		new(big.Int).SetInt64(1),
// 		new(big.Int).SetInt64(1),
// 		new(big.Int).SetInt64(1),
// 		new(big.Int).SetInt64(-1),
// 		new(big.Int).SetInt64(100),
// 		new(big.Int).SetInt64(100),
// 	}
// 	list := make([][]*big.Int, 500)
// 	for i := range list {
// 		list[i] = bigs
// 	}

// 	var buf bytes.Buffer
// 	encoder := gob.NewEncoder(&buf)
// 	encoder.Encode(list)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// for j := 0; j < 1000; j++ {
// 		var list [][]*big.Int
// 		gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&list)
// 		// }
// 	}
// }

// func BenchmarkReadsEncode(b *testing.B) {
// 	list := make([]*Reads, 500)
// 	for i := range list {
// 		list[i] = reads
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		var buf bytes.Buffer
// 		encoder := gob.NewEncoder(&buf)
// 		encoder.Encode(list)
// 	}
// }

// func BenchmarkReadsDecode(b *testing.B) {
// 	list := make([]*Reads, 500)
// 	for i := range list {
// 		list[i] = reads
// 	}
// 	var buf bytes.Buffer
// 	encoder := gob.NewEncoder(&buf)
// 	encoder.Encode(list)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		var list []*Reads
// 		gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&list)
// 	}
// }

// func BenchmarkWritesEncode(b *testing.B) {
// 	list := make([]*Writes, 500)
// 	for i := range list {
// 		list[i] = writes
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		var buf bytes.Buffer
// 		encoder := gob.NewEncoder(&buf)
// 		encoder.Encode(list)
// 	}
// }

// func BenchmarkWritesDecode(b *testing.B) {
// 	list := make([]*Writes, 500)
// 	for i := range list {
// 		list[i] = writes
// 	}
// 	var buf bytes.Buffer
// 	encoder := gob.NewEncoder(&buf)
// 	encoder.Encode(list)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		var list []*Writes
// 		gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&list)
// 	}
// }

// func TestReadsEncodeDecode(t *testing.T) {
// 	list := make([]*Reads, 10)
// 	for i := range list {
// 		addr1 := ethCommon.BytesToAddress([]byte{byte(i)})
// 		addr2 := ethCommon.BytesToAddress([]byte{byte(i + 10)})
// 		addr3 := ethCommon.BytesToAddress([]byte{byte(i + 20)})
// 		list[i] = &Reads{
// 			ClibReads: map[Address]*ReadSet{
// 				Address(string(addr2.Bytes())): {
// 					HashMapReads: []HashMapRead{
// 						{
// 							HashMapAccess: HashMapAccess{
// 								ID:  "hashmap",
// 								Key: "key",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			BalanceReads: map[Address]*big.Int{
// 				Address(string(addr1.Bytes())): new(big.Int).SetInt64(int64(i)),
// 				Address(string(addr2.Bytes())): new(big.Int).SetInt64(int64(i)),
// 				Address(string(addr3.Bytes())): new(big.Int).SetInt64(int64(i)),
// 			},
// 			EthStorageReads: nil,
// 		}
// 	}

// 	var buf bytes.Buffer
// 	err := gob.NewEncoder(&buf).Encode(list)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var list2 []*Reads
// 	gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&list2)
// 	fmt.Println(list2)
// }
