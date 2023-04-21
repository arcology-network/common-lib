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
	out := in.Decode(data).(Uint8)

	if (in) != (out) {
		t.Error("Uint8 Mismatched !")
	}
}

func TestUint8s(t *testing.T) {
	in := []uint8{1, 2, 3, 4, 5}
	data := Uint8s(in).Encode()
	out := Uint8s(in).Decode(data).(Uint8s)

	if !reflect.DeepEqual(in, []uint8(out)) {
		t.Error("Uint8s Mismatched !")
	}
}

func TestUint32(t *testing.T) {
	in := uint32(99)
	data := Uint32(in).Encode()
	out := Uint32(in).Decode(data).(Uint32)
	if Uint32(in) != out {
		t.Error("Uint32 Mismatched !")
	}
}

func TestUint32s(t *testing.T) {
	in := []uint32{1, 2, 3, 4, 5}
	data := Uint32s(in).Encode()
	out := Uint32s(in).Decode(data).(Uint32s)

	if !reflect.DeepEqual(Uint32s(in), (out)) {
		t.Error("Uint32 Mismatched !")
	}
}

func TestUint64s(t *testing.T) {
	in := []uint64{11, 22, 33, 44, 555555}
	data := Uint64s(in).Encode()
	out := Uint64s(in).Decode(data).(Uint64s)

	if !reflect.DeepEqual(Uint64s(in), out) {
		t.Error("Uint64s Mismatched !")
	}
}

func TestBools(t *testing.T) {
	in := []bool{false, false, true, true}
	data := Bools(in).Encode()
	out := Bools(in).Decode(data).(Bools)

	if !reflect.DeepEqual(Bools(in), out) {
		t.Error("Mismatch !")
		fmt.Println(in)
		fmt.Println()
		fmt.Println(out)
	}
}

func TestString(t *testing.T) {
	in := "0x1234567890abcdef"
	bytes := String(in).Encode()
	out := String("").Decode(bytes).(String)
	if !reflect.DeepEqual(String(in), out) {
		t.Error("strings mismatch !")
	}
}

func TestStrings(t *testing.T) {
	in := []string{"", "111", "2222", "33333", ""}
	in2 := []string{"999", "111", "2222", "33333", ""}
	bytes := Strings(in).Encode()
	out := Strings([]string{}).Decode(bytes).(Strings)
	if !reflect.DeepEqual(Strings(in), (out)) {
		t.Error("strings mismatch !")
	}

	buffer := Byteset([][]byte{
		Strings(in).Encode(),
		Strings(in2).Encode(),
	}).Encode()

	fields := Byteset{}.Decode(buffer).(Byteset)

	str1 := (Strings{}).Decode(fields[0]).(Strings)[0]
	str2 := (Strings{}).Decode(fields[1]).(Strings)[2]
	if str1 != in[0] ||
		str2 != in2[2] {
		t.Error("strings mismatch !")
	}
}

func TestBytes(t *testing.T) {
	inner := [][]byte{
		Uint32s([]uint32{1, 2, 3, 4, 5}).Encode(),
		Bools([]bool{false, false, true, true}).Encode(),
	}

	in := [][]byte{
		{},
		Uint32s([]uint32{1, 2, 3, 4, 5}).Encode(),
		Uint64s([]uint64{11, 22, 33, 44, 555555}).Encode(),
		Bools([]bool{false, false, true, true}).Encode(),
		{},
		{},
		Byteset(inner).Encode(),
	}

	buffer := Byteset(in).Encode()
	out := Byteset(in).Decode(buffer).(Byteset)

	if !reflect.DeepEqual(in, [][]byte(out)) {
		t.Error("Mismatch !")
		fmt.Println(in)
		fmt.Println()
		fmt.Println(out)
	}

}

func TestBigint(t *testing.T) {
	in := big.NewInt(1234567)
	data := (*Bigint)(in).Encode()
	out := (&Bigint{}).Decode(data).(*Bigint)

	if !reflect.DeepEqual(in, (*big.Int)(out)) {
		t.Error("Mismatch !")
		fmt.Println(in)
		fmt.Println()
		fmt.Println(out)
	}

	in = big.NewInt(-456789)
	data = (*Bigint)(in).Encode()
	out = (&Bigint{}).Decode(data).(*Bigint)

	if !reflect.DeepEqual((in), (*big.Int)(out)) {
		t.Error("Mismatch !")
		fmt.Println(in)
		fmt.Println()
		fmt.Println(out)
	}
}

func TestByteSetAndClone(t *testing.T) {
	byteset := make([][]byte, 0, 600000)
	for i := 0; i < 20; i++ {
		byteset = append(byteset, [][]byte{
			Uint32s([]uint32{1, 2, 3, 4, 5}).Encode(),
			Uint64s([]uint64{11, 22, 33, 44, 555555}).Encode(),
			Bools([]bool{false, false, true, true}).Encode(),
		}...)
	}

	for i := 0; i < 20; i++ {
		byteset = append(byteset, [][]byte{
			Uint32s([]uint32{31, 42, 53, 24, 15}).Encode(),
			Uint64s([]uint64{211, 622, 733, 484, 3555555}).Encode(),
			Bools([]bool{false, false, true, true}).Encode(),
		}...)
	}

	clone := Byteset(byteset).Clone()
	for i := 0; i < len(byteset); i++ {
		if !reflect.DeepEqual(clone[i], byteset[i]) {
			t.Error("Mismatch !")
		}
	}
}

func TestByteGroupClone(t *testing.T) {
	_1 := make([][]byte, 0, 2)
	for i := 0; i < 20; i++ {
		_1 = append(_1, [][]byte{
			Uint32s([]uint32{1, 2, 3, 4, 5}).Encode(),
			Uint64s([]uint64{11, 22, 33, 44, 555555}).Encode(),
			Bools([]bool{false, false, true, true}).Encode(),
		}...)
	}

	_2 := make([][]byte, 0, 2)
	for i := 0; i < 20; i++ {
		_2 = append(_2, [][]byte{
			Uint32s([]uint32{8, 8, 7, 6, 5}).Encode(),
			Uint64s([]uint64{411, 522, 363, 44, 5755555}).Encode(),
			Bools([]bool{false, false, true, true}).Encode(),
		}...)
	}

	byteGroup := [][][]byte{_1, _2}
	clone := Bytegroup(byteGroup).Clone()

	for i := 0; i < len(byteGroup); i++ {
		for j := 0; j < len(byteGroup[i]); j++ {
			if !reflect.DeepEqual(clone[i][j], byteGroup[i][j]) {
				t.Error("Mismatch !")
			}
		}
	}
}

func BenchmarkEncodeBytesPerformance(t *testing.B) {
	byteset := make([][]byte, 0, 600000)
	for i := 0; i < 600000; i++ {
		byteset = append(byteset, [][]byte{
			Uint32s([]uint32{1, 2, 3, 4, 5}).Encode(),
			Uint64s([]uint64{11, 22, 33, 44, 555555}).Encode(),
			Bools([]bool{false, false, true, true}).Encode(),
		}...)
	}

	t0 := time.Now()
	Byteset(byteset).Encode()
	fmt.Println("Bytes encoding : ", time.Since(t0))
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

func TestEncoderUint32(t *testing.T) {
	// n1 := Uint32(999999)
	// buffer := make([]byte, Encoder{}.Size([]interface{}{n1, n1}))
	// Encoder{}.ToBuffer(buffer, []interface{}{n1, n1})

	// fields := [][]byte(Byteset{}.Decode(buffer).(Byteset))
	// if n1 != Uint32(0).Decode(fields[0]) ||
	// 	n1 != Uint32(0).Decode(fields[1]) {
	// 	t.Error("Mismatch !")
	// }
}

func TestEncoderBigint(t *testing.T) {
	// v := big.NewInt(-999999)
	// n1 := Bigint(*v)

	// v0 := big.NewInt(11)
	// n2 := Bigint(*v0)

	// buffer := make([]byte, Encoder{}.Size([]interface{}{&n1, &n2}))
	// Encoder{}.ToBuffer(buffer, []interface{}{&n1, &n2})

	// fields := [][]byte(Byteset{}.Decode(buffer).(Byteset))
	// lft := (*big.Int)((&Bigint{}).Decode(fields[0]).(*Bigint))
	// rgt := (*big.Int)((&Bigint{}).Decode(fields[1]).(*Bigint))

	// buf := n1.Encode()
	// fmt.Print(buf)
	// if v.Cmp(lft) != 0 || v0.Cmp(rgt) != 0 {
	// 	t.Error("Mismatch !")
	// }
}

func TestEncoderBigintAndNil(t *testing.T) {
	v := big.NewInt(-999999)
	n1 := Bigint(*v)

	buffer := make([]byte, Encoder{}.Size([]interface{}{nil, nil, nil}))
	Encoder{}.ToBuffer(buffer, []interface{}{nil, nil, nil})

	fields := [][]byte(Byteset{}.Decode(buffer).(Byteset))
	_0 := (*big.Int)((&Bigint{}).Decode(fields[0]).(*Bigint))
	_1 := (*big.Int)((&Bigint{}).Decode(fields[1]).(*Bigint))
	_2 := (*big.Int)((&Bigint{}).Decode(fields[2]).(*Bigint))

	buf := n1.Encode()
	fmt.Print(buf)
	if _0.Cmp(&big.Int{}) != 0 || _1.Cmp(&big.Int{}) != 0 || _2.Cmp(&big.Int{}) != 0 {
		t.Error("Mismatch !")
	}
}

func TestStringsetFlatten(t *testing.T) {
	str0 := []string{"123456", "987654"}
	str1 := []string{"abcdef", "zqwert"}

	flattened := Stringset([][]string{str0, str1}).Flatten()
	if flattened[0] != "123456" ||
		flattened[1] != "987654" ||
		flattened[2] != "abcdef" ||
		flattened[3] != "zqwert" {
		t.Error("Mismatch !")
	}
}

// func TestStringsetCodec(t *testing.T) {
// 	str0 := []string{"123456", "987654"}
// 	str1 := []string{"abcdef", "zqwert"}

// 	buffer := Stringset([][]string{str0, str1}).Encode()
// 	out := Stringset{}.Decode(buffer).(Stringset)

// 	if out[0][0] != "123456" ||
// 		out[0][1] != "987654" ||
// 		out[1][0] != "abcdef" ||
// 		out[1][1] != "zqwert" {
// 		t.Error("Mismatch !")
// 	}
// }

func TestHash16s(t *testing.T) {
	in := [][16]byte{{1, 2, 3, 4, 5}, {5, 6, 7, 8, 9}}

	data := Hash16s(in).Encode()
	out := Hash16s(in).Decode(data).(Hash16s)

	if !reflect.DeepEqual(in, ([][16]byte)(out)) {
		t.Error("Uint8s Mismatched !")
	}
}

func TestHash32s(t *testing.T) {
	in := [][32]byte{{1, 2, 3, 4, 5}, {5, 6, 7, 8, 9}}

	data := Hash32s(in).Encode()
	out := Hash32s(in).Decode(data).(Hash32s)

	if !reflect.DeepEqual(in, ([][32]byte)(out)) {
		t.Error("Uint8s Mismatched !")
	}
}

func TestHash64s(t *testing.T) {
	in := [][64]byte{{1, 2, 3, 4, 5}, {5, 6, 7, 8, 9}}

	data := Hash64s(in).Encode()
	out := Hash64s(in).Decode(data).(Hash64s)

	if !reflect.DeepEqual(in, ([][64]byte)(out)) {
		t.Error("Uint8s Mismatched !")
	}

	clone := Hash64s(in).Clone()
	if !reflect.DeepEqual(clone[0], in[0]) {
		t.Error("Hash64s Mismatched !")
	}

	if !reflect.DeepEqual(clone[1], in[1]) {
		t.Error("Hash64s Mismatched !")
	}
}
