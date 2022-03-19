package types

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestAccessRecordEncoding(t *testing.T) {
	records := &TxAccessRecords{
		Hash:     "0x1234567",
		ID:       99,
		Accesses: [][]byte{[]byte("1"), []byte("2")},
	}

	buffer := records.Encode()
	out := (&TxAccessRecords{}).Decode(buffer)
	if !reflect.DeepEqual(records, out) {
		t.Error("Error")
	}
}

func BenchmarkAccessRecordSetEncoding(b *testing.B) {
	recordVec := make([]*TxAccessRecords, 1000000)
	for i := 0; i < len(recordVec); i++ {
		recordVec[i] = &TxAccessRecords{
			Hash:     "0x1234567",
			ID:       uint32(i),
			Accesses: [][]byte{[]byte(fmt.Sprint(9)), []byte("2456")},
		}
	}

	t0 := time.Now()
	buffer, _ := TxAccessRecordSet(recordVec).GobEncode()
	fmt.Println("GobEncode():", time.Now().Sub(t0))

	out := new(TxAccessRecordSet)
	t0 = time.Now()
	out.GobDecode(buffer)
	fmt.Println("GobDecode():", time.Now().Sub(t0))

	// for i := 0; i < len(recordVec); i++ {
	// 	if !reflect.DeepEqual(recordVec[i], (*out)[i]) {
	// 	//	t.Error("Error")
	// 	}
	// }
}
