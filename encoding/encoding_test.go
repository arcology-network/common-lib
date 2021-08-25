package encoding

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestUint64ss(t *testing.T) {
	uins := []uint64{1, 4, 3, 2, 6, 5, 6}

	bys := Uint64s(uins).Encode()
	fmt.Printf("bys=%v\n", bys)

	nuins := Uint64s(uins).Unique()
	fmt.Printf("nuins=%v\n", nuins)

	bys = Uint64s(uins).Encode()
	fmt.Printf("bys=%v\n", bys)
}

func TestUint8(t *testing.T) {
	in := uint8(244)
	data := Uint8(in).Encode()
	out := Uint8(in).Decode(data)

	if !reflect.DeepEqual(in, out) {
		t.Error("Uint8 Mismatched !")
	}
}

func TestUint8s(t *testing.T) {
	in := []uint8{1, 2, 3, 4, 5}
	data := Uint8s(in).Encode()
	out := Uint8s(in).Decode(data)

	if !reflect.DeepEqual(in, out) {
		t.Error("Uint8s Mismatched !")
	}
}

func TestUint32s(t *testing.T) {
	in := []uint32{1, 2, 3, 4, 5}
	data := Uint32s(in).Encode()
	out := Uint32s(in).Decode(data)

	if !reflect.DeepEqual(in, out) {
		t.Error("Uint32 Mismatched !")
	}
}
func TestUint64s(t *testing.T) {
	in := []uint64{11, 22, 33, 44, 555555}
	data := Uint64s(in).Encode()
	out := Uint64s(in).Decode(data)

	if !reflect.DeepEqual(in, out) {
		t.Error("Uint64s Mismatched !")
	}
}

func TestBools(t *testing.T) {
	in := []bool{false, false, true, true}
	data := Bools(in).Encode()
	out := Bools(in).Decode(data)

	if !reflect.DeepEqual(in, out) {
		t.Error("Mismatch !")
		fmt.Println(in)
		fmt.Println()
		fmt.Println(out)
	}
}

func TestBytes(t *testing.T) {
	in := [][]byte{
		Uint32s([]uint32{1, 2, 3, 4, 5}).Encode(),
		Uint64s([]uint64{11, 22, 33, 44, 555555}).Encode(),
		Bools([]bool{false, false, true, true}).Encode(),
	}

	data := Byteset(in).Encode()
	out := Byteset(in).Decode(data)

	if !reflect.DeepEqual(in, out) {
		t.Error("Mismatch !")
		fmt.Println(in)
		fmt.Println()
		fmt.Println(out)
	}

	byteset := [][]byte{{byte(1)}, {byte(2)}, {byte(3)}}
	if !reflect.DeepEqual(Byteset(byteset).Flatten(), append([]byte{byte(1)}, append([]byte{byte(2)}, []byte{byte(3)}[:]...)[:]...)) {
		t.Error("Mismatch !")
		fmt.Println(Byteset(byteset).Flatten())
		fmt.Println()
		fmt.Println(append([]byte{byte(1)}, append([]byte{byte(2)}, []byte{byte(3)}[:]...)[:]...))
	}
}

func TestBytesPerformance(t *testing.T) {
	byteset := make([][]byte, 0, 500000)
	for i := 0; i < 500000; i++ {
		byteset = append(byteset, [][]byte{
			Uint32s([]uint32{1, 2, 3, 4, 5}).Encode(),
			Uint64s([]uint64{11, 22, 33, 44, 555555}).Encode(),
			Bools([]bool{false, false, true, true}).Encode(),
		}...)
	}

	t0 := time.Now()
	bytes := Byteset(byteset).Encode()
	fmt.Println("Encode Bytes : ", time.Now().Sub(t0))

	t0 = time.Now()
	Byteset(byteset).Decode(bytes)
	fmt.Println("Decode Bytes : ", time.Now().Sub(t0))
}
