package types

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestEuresultEncoding(t *testing.T) {
	eu := &EuResult{
		H:           "0x1234567",
		ID:          99,
		Transitions: [][]byte{[]byte("1"), []byte("2")},
		DC:          nil,
		Status:      0,
		GasUsed:     0,
	}

	if eu.DC == nil {
		fmt.Println()
	}

	buffer := eu.Encode()
	out := (&EuResult{}).Decode(buffer)
	if !reflect.DeepEqual(eu, out) {
		t.Error("Error")
	}
}

func TestEuResultEncoding(t *testing.T) {
	dc := &DeferCall{
		DeferID:         "7777",
		ContractAddress: "45678",
		Signature:       "xxxx",
	}

	euresult := &EuResult{
		H:           "1234",
		ID:          uint32(99),
		Transitions: [][]byte{[]byte(fmt.Sprint("xxxxxx")), []byte("+++++")},
		DC:          dc,
		Status:      0,
		GasUsed:     99,
	}

	t0 := time.Now()
	buffer, _ := euresult.GobEncode()
	fmt.Println("GobDecode():", time.Now().Sub(t0))

	t0 = time.Now()
	out := new(EuResult)
	out.GobDecode(buffer)
	fmt.Println("GobDecode():", time.Now().Sub(t0))

	if !reflect.DeepEqual(*euresult, *out) {
		t.Error("Error")
	}
}

func TestEuResultsEncoding(t *testing.T) {
	dc := &DeferCall{
		DeferID:         "7777",
		ContractAddress: "45678",
		Signature:       "xxxx",
	}

	euresults := make([]*EuResult, 10)
	for i := 0; i < len(euresults); i++ {
		euresults[i] = &EuResult{
			H:           "0x1234567",
			ID:          uint32(99),
			Transitions: [][]byte{[]byte(fmt.Sprint("++++")), []byte("||||")},
			DC:          dc,
			Status:      11,
			GasUsed:     99,
		}
	}

	t0 := time.Now()
	buffer, _ := Euresults(euresults).GobEncode()
	fmt.Println("EuResults GobEncode():", time.Now().Sub(t0))

	out := new(Euresults)
	out.GobDecode(buffer)
	fmt.Println("EuResults GobDecode():", time.Now().Sub(t0))

	for i := 0; i < len(euresults); i++ {
		if !reflect.DeepEqual(euresults[i], (*out)[i]) {
			t.Error("Error")
		}
	}
}

func BenchmarkEuResultsEncoding(b *testing.B) {
	dc := &DeferCall{
		DeferID:         "7777",
		ContractAddress: "45678",
		Signature:       "xxxx",
	}

	euresults := make([]*EuResult, 1000000)
	for i := 0; i < len(euresults); i++ {
		euresults[i] = &EuResult{
			H:           "0x1234567",
			ID:          uint32(99),
			Transitions: [][]byte{[]byte(fmt.Sprint("++++")), []byte("||||")},
			DC:          dc,
			Status:      11,
			GasUsed:     99,
		}
	}

	t0 := time.Now()
	buffer, _ := Euresults(euresults).GobEncode()
	fmt.Println("EuResults GobEncode():", time.Now().Sub(t0))

	out := new(Euresults)
	out.GobDecode(buffer)
	fmt.Println("EuResults GobDecode():", time.Now().Sub(t0))
}
