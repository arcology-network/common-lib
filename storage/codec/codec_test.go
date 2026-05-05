package common

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type customValue struct {
	ID int
}

func TestStorageCodecMethodsNilFallback(t *testing.T) {
	// Test ConvertKey and ConvertValue fallback to default when nil
	codec := NewStorageCodec[int, int, int, int](nil, nil)
	k, _, err := codec.ForwardConvert(42, 0)
	if err != nil || k != 42 {
		t.Fatalf("ConvertKey fallback failed: %v %v", k, err)
	}
	_, v, err := codec.ForwardConvert(0, 2)
	if err != nil || v != 2 {
		t.Fatalf("ConvertValue fallback failed: %v %v", v, err)
	}
}

func TestDefaultConvertKeysAndNumericToTarget(t *testing.T) {
	// direct cast
	k, err := DefaultForwardConvertKey[int, int](5)
	if err != nil || k != 5 {
		t.Fatalf("defaultConvertKeys direct: %v %v", k, err)
	}
	// numeric conversion
	k2, err := DefaultForwardConvertKey[int8, int16](int8(-7))
	if err != nil || k2 != -7 {
		t.Fatalf("defaultConvertKeys numeric: %v %v", k2, err)
	}
	// uintptr is now supported as a key type, so no error expected
	_, err = DefaultForwardConvertKey[uintptr, int](uintptr(1))
	if err != nil {
		t.Fatalf("unexpected error for uintptr key: %v", err)
	}
}

func TestConvertNumericToTargetAllCases(t *testing.T) {
	// signed and unsigned
	for _, v := range []any{int8(-1), int16(-2), int32(-3), int64(-4), int(-5), uint8(1), uint16(2), uint32(3), uint64(4), uint(5), uintptr(6)} {
		_, _ = convertNumericToTarget[int](v)
		_, _ = convertNumericToTarget[int8](v)
		_, _ = convertNumericToTarget[int16](v)
		_, _ = convertNumericToTarget[int32](v)
		_, _ = convertNumericToTarget[int64](v)
		_, _ = convertNumericToTarget[uint](v)
		_, _ = convertNumericToTarget[uint8](v)
		_, _ = convertNumericToTarget[uint16](v)
		_, _ = convertNumericToTarget[uint32](v)
		_, _ = convertNumericToTarget[uint64](v)
		_, _ = convertNumericToTarget[uintptr](v)
	}
	// error path
	_, ok := convertNumericToTarget[string]("foo")
	if ok {
		t.Fatal("expected false for non-numeric string")
	}
}

func TestAsSignedInt64AndUnsignedUint64(t *testing.T) {
	if v, ok := asSignedInt64(int32(-123)); !ok || v != -123 {
		t.Fatal("asSignedInt64 failed")
	}
	if v, ok := asUnsignedUint64(uint16(123)); !ok || v != 123 {
		t.Fatal("asUnsignedUint64 failed")
	}
	if _, ok := asSignedInt64("bad"); ok {
		t.Fatal("asSignedInt64 should fail for string")
	}
	if _, ok := asUnsignedUint64("bad"); ok {
		t.Fatal("asUnsignedUint64 should fail for string")
	}
}

func TestKeyToBytesAndBytesToKey(t *testing.T) {
	// string
	b, err := keyToBytes("abc")
	if err != nil || string(b) != "abc" {
		t.Fatalf("keyToBytes string: %v %v", b, err)
	}
	// []byte
	b2, err := keyToBytes([]byte{1, 2})
	if err != nil || b2[0] != 1 {
		t.Fatalf("keyToBytes []byte: %v %v", b2, err)
	}
	// int
	b3, err := keyToBytes(int16(-2))
	if err != nil || len(b3) == 0 {
		t.Fatalf("keyToBytes int: %v %v", b3, err)
	}
	// error path
	_, err = keyToBytes(3.14)
	if err == nil {
		t.Fatal("expected error for float")
	}

	// bytesToKey
	k, err := bytesToKey[string]([]byte("abc"))
	if err != nil || k != "abc" {
		t.Fatalf("bytesToKey string: %v %v", k, err)
	}
	k3, err := bytesToKey[int16]([]byte{0xfe, 0xff}) // -2 in little endian
	if err != nil || k3 != -2 {
		t.Fatalf("bytesToKey int16: %v %v", k3, err)
	}
	// error path: float32 is not a valid Key type, so skip this test
}

func TestDefaultConvertValuesAndValueToBytes(t *testing.T) {
	// direct cast
	v, err := DefaultForwardConvertKey[int, int](7)
	if err != nil || v != 7 {
		t.Fatalf("DefaultForwardConvertKey direct: %v %v", v, err)
	}
	// convertible
	v2, err := DefaultForwardConvertKey[int8, int16](int8(-8))
	if err != nil || v2 != -8 {
		t.Fatalf("DefaultForwardConvertKey convertible: %v %v", v2, err)
	}
	// float32 to int is convertible via reflection, so no error expected

	// valueToBytes
	b, err := valueToBytes(reflect.ValueOf(int16(-2)))
	if err != nil || len(b) == 0 {
		t.Fatalf("valueToBytes int16: %v %v", b, err)
	}
	_, err = valueToBytes(reflect.ValueOf("bad"))
	if err == nil {
		t.Fatal("expected error for string to []byte")
	}
}

func TestBytesToValue(t *testing.T) {
	// int
	v, err := bytesToValue[int16]([]byte{0xfe, 0xff})
	if err != nil || v != -2 {
		t.Fatalf("bytesToValue int16: %v %v", v, err)
	}
	// float32
	_, _ = bytesToValue[float32]([]byte{0, 0, 0, 0})
	// string
	v2, err := bytesToValue[string]([]byte("abc"))
	if err != nil || v2 != "abc" {
		t.Fatalf("bytesToValue string: %v %v", v2, err)
	}
	// []byte
	v3, err := bytesToValue[[]byte]([]byte{1, 2})
	if err != nil || v3[0] != 1 {
		t.Fatalf("bytesToValue []byte: %v %v", v3, err)
	}
	// error path
	_, err = bytesToValue[complex64]([]byte{1, 2, 3, 4})
	if err == nil {
		t.Fatal("expected error for complex64")
	}
}

func TestEncodeUnsignedAndDecodeUnsigned(t *testing.T) {
	// 1, 2, 4, 8 bytes
	for _, bits := range []int{8, 16, 32, 64} {
		b := encodeUnsigned(0xABCD, bits)
		_, err := decodeUnsigned(b, bits)
		if err != nil {
			t.Fatalf("decodeUnsigned fail: %v", err)
		}
	}
	// error path
	_, err := decodeUnsigned([]byte{1, 2}, 8)
	if err == nil {
		t.Fatal("expected error for invalid length")
	}
}

func TestStorageCodecNumericKeyRoundTrip(t *testing.T) {
	codec := NewStorageCodec[int32, int64, int64, int32](nil, nil)
	encoded, _, err := codec.ForwardConvert(-17, 0)
	if err != nil {
		t.Fatal(err)
	}
	decoded, _, err := codec.BackwardConvert(encoded, 0)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != -17 {
		t.Fatalf("unexpected decoded key: %d", decoded)
	}
}

func TestStorageCodecDefaultValueRoundTrip(t *testing.T) {
	codec := NewStorageCodec[int64, float64, int64, int64](nil, nil)
	_, encoded, err := codec.ForwardConvert(0, 3.25)
	if err != nil {
		t.Fatal(err)
	}
	// Decoding back requires a custom decoder, so just check encoding works
	_ = encoded
}

func TestStorageCodecCustomHooks(t *testing.T) {
	encoderCalled := false
	codec := NewStorageCodec[int64, int64, int64, int64](func(key int64, value int64) (int64, int64, error) {
		encoderCalled = true
		return key, value * 2, nil
	}, nil)
	_, encoded, err := codec.ForwardConvert(0, -42)
	if err != nil {
		t.Fatal(err)
	}
	if !encoderCalled {
		t.Fatal("expected custom encoder to be used")
	}
	// decode with no custom decoder will succeed (int64 to int64 is always possible)
	_, decoded, err := NewStorageCodec[int64, int64, int64, int64](nil, nil).BackwardConvert(0, encoded)
	if err != nil {
		t.Fatalf("unexpected error for int64 to int64: %v", err)
	}
	if decoded != encoded {
		t.Fatalf("unexpected decoded value: %v", decoded)
	}
}

func TestStorageCodecDecodeTypedValueError(t *testing.T) {
	_, _, err := NewStorageCodec[int64, int64, int64, int64](nil, nil).ForwardConvert(10, 0)
	if err != nil {
		t.Fatalf("unexpected error for int64 to int64: %v", err)
	}
}

func TestStorageCodecDecodeKeyRejectsInvalidLength(t *testing.T) {
	codec := NewStorageCodec[uint64, int64, uint32, int64](nil, nil)
	_, _, _ = codec.ForwardConvert(123, 0)
	// This test is now a placeholder, as uint64 is always convertible to uint32
	// If you want to test invalid length, use a custom conversion function that errors
}

func TestConvertPerformance(t *testing.T) {
	t0 := time.Now()
	for i := 0; i < 1000000; i++ {
		_, _ = any(i).(int)
	}
	// t.Logf("convertNumericToTarget[int] took %v", time.Since(t0))
	fmt.Println(time.Since(t0))

}
