package codec

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestUint8(t *testing.T) {
	in := Uint8(244)
	data := in.Encode()
	out := in.Decode(data)

	if uint8(in) != out {
		t.Error("Uint8 Mismatched !")
	}
}

func TestUint8s(t *testing.T) {
	in := Uint8s{1, 2, 3, 4, 5}
	data := Uint8s(in).Encode()
	out := Uint8s(in).Decode(data)

	if !reflect.DeepEqual(in, out) {
		t.Error("Uint8s Mismatched !")
	}
}

func TestUint32s(t *testing.T) {
	in := []Uint32{1, 2, 3, 4, 5}
	data := Uint32s(in).Encode()
	out := Uint32s(in).Decode(data)

	if !reflect.DeepEqual(Uint32s(in), out) {
		t.Error("Uint32 Mismatched !")
	}
}
func TestUint64s(t *testing.T) {
	in := Uint64s([]Uint64{11, 22, 33, 44, 555555})
	data := Uint64s(in).Encode()
	out := Uint64s(in).Decode(data).(Uint64s)

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

func TestStrings(t *testing.T) {
	in := []string{"111", "2222", "33333"}
	bytes := Strings(in).Encode()
	out := Strings([]string{}).Decode(bytes)
	if !reflect.DeepEqual(in, []string(out)) {
		t.Error("strings mismatch !")
	}

}

func TestBytes(t *testing.T) {
	in := [][]byte{
		Uint32s([]Uint32{1, 2, 3, 4, 5}).Encode(),
		Uint64s([]Uint64{11, 22, 33, 44, 555555}).Encode(),
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

func TestBigint(t *testing.T) {
	in := big.NewInt(1234567)
	data := (*Bigint)(in).Encode()
	out := (&Bigint{}).Decode(data)

	if !reflect.DeepEqual(in, out) {
		t.Error("Mismatch !")
		fmt.Println(in)
		fmt.Println()
		fmt.Println(out)
	}

	in = big.NewInt(-456789)
	data = (*Bigint)(in).Encode()
	out = (&Bigint{}).Decode(data)

	if !reflect.DeepEqual(in, out) {
		t.Error("Mismatch !")
		fmt.Println(in)
		fmt.Println()
		fmt.Println(out)
	}
}

func TestBytesPerformance(t *testing.T) {
	byteset := make([][]byte, 0, 500000)
	for i := 0; i < 500000; i++ {
		byteset = append(byteset, [][]byte{
			Uint32s([]Uint32{1, 2, 3, 4, 5}).Encode(),
			Uint64s([]Uint64{11, 22, 33, 44, 555555}).Encode(),
			Bools([]bool{false, false, true, true}).Encode(),
		}...)
	}

	t0 := time.Now()
	Byteset(byteset).Encode()
	fmt.Println("Bytes encoding : ", time.Now().Sub(t0))
}

func TestConcatenateStrings(t *testing.T) {
	paths := []string{}
	var concated string
	for i := 0; i < 50000; i++ {
		str := fmt.Sprint(rand.Float64())
		paths = append(paths, str)
		concated += fmt.Sprint(str)
	}

	t0 := time.Now()
	length := uint(len(paths))
	var buf bytes.Buffer
	pathLen := make([]uint32, length)
	for i, p := range paths {
		pathLen[i] = uint32(len(p))
		fmt.Fprintf(&buf, "%s", p)
	}
	fmt.Println("Concatenate Strings : ", buf.Len(), time.Now().Sub(t0))

	t0 = time.Now()
	buffer := Strings(paths).Flatten()
	fmt.Println("Flatten Strings : ", len(buffer), time.Now().Sub(t0))

	if !bytes.Equal(buf.Bytes(), buffer) {
		t.Error("Mismatch !")
	}
}
