/*
*   Copyright (c) 2026 Arcology Network

*   This program is free software: you can redistribute it and/or modify
*   it under the terms of the GNU General Public License as published by
*   the Free Software Foundation, either version 3 of the License, or
*   (at your option) any later version.

*   This program is distributed in the hope that it will be useful,
*   but WITHOUT ANY WARRANTY; without even the implied warranty of
*   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*   GNU General Public License for more details.

*   You should have received a copy of the GNU General Public License
*   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"reflect"
	"unsafe"

	stgintf "github.com/arcology-network/common-lib/storage/interface"
)

type Convertible interface {
	stgintf.Key | []byte
}

type StorageCodec[K0 Convertible, V0 any, K1 Convertible, V1 any] struct {
	// Converts the key and value from the cached  format to the backend format.
	ForwardConvert func(K0, V0) (K1, V1, error)

	// Converts the key and value from the backend  format to the cached format.
	BackwardConvert func(K1, V1) (K0, V0, error)
}

func NewStorageCodec[K0 Convertible, V0 any, K1 Convertible, V1 any](
	forwardConvert func(K0, V0) (K1, V1, error),
	backwardConvert func(K1, V1) (K0, V0, error),
) *StorageCodec[K0, V0, K1, V1] {
	stgCodec := &StorageCodec[K0, V0, K1, V1]{
		ForwardConvert:  forwardConvert,
		BackwardConvert: backwardConvert,
	}

	if stgCodec.ForwardConvert == nil {
		stgCodec.ForwardConvert = DefaultForwardConvert[K0, V0, K1, V1]
	}

	if stgCodec.BackwardConvert == nil {
		stgCodec.BackwardConvert = DefaultBackwardConvert[K0, V0, K1, V1]
	}
	return stgCodec
}

func DefaultForwardConvert[K0 Convertible, V0 any, K1 Convertible, V1 any](key K0, value V0) (K1, V1, error) {
	convertedKey, err := DefaultForwardConvertKey[K0, K1](key)
	if err != nil {
		return zero[K1](), zero[V1](), err
	}

	convertedValue, err := DefaultForwardConvertValue[V0, V1](value)
	if err != nil {
		return zero[K1](), zero[V1](), err
	}
	return convertedKey, convertedValue, nil
}

func DefaultBackwardConvert[K0 Convertible, V0 any, K1 Convertible, V1 any](key K1, value V1) (K0, V0, error) {
	convertedKey, err := DefaultForwardConvertKey[K1, K0](key)
	if err != nil {
		return zero[K0](), zero[V0](), err
	}

	convertedValue, err := DefaultForwardConvertValue[V1, V0](value)
	if err != nil {
		return zero[K0](), zero[V0](), err
	}
	return convertedKey, convertedValue, nil
}

func DefaultForwardConvertKey[K0, K1 Convertible](key K0) (K1, error) {
	if converted, ok := any(key).(K1); ok {
		return converted, nil
	}

	source := any(key)
	if converted, ok := convertNumericToTarget[K1](source); ok {
		return converted, nil
	}

	raw, err := keyToBytes(source)
	if err != nil {
		return zero[K1](), err
	}

	return bytesToKey[K1](raw)
}

func convertNumericToTarget[K Convertible](value any) (K, bool) {
	var target K
	signed, signedOK := asSignedInt64(value)
	unsigned, unsignedOK := asUnsignedUint64(value)

	if !signedOK && !unsignedOK {
		return zero[K](), false
	}

	switch any(target).(type) {
	case int:
		if signedOK {
			return any(int(signed)).(K), true
		}
		return any(int(unsigned)).(K), true
	case int8:
		if signedOK {
			return any(int8(signed)).(K), true
		}
		return any(int8(unsigned)).(K), true
	case int16:
		if signedOK {
			return any(int16(signed)).(K), true
		}
		return any(int16(unsigned)).(K), true
	case int32:
		if signedOK {
			return any(int32(signed)).(K), true
		}
		return any(int32(unsigned)).(K), true
	case int64:
		if signedOK {
			return any(signed).(K), true
		}
		return any(int64(unsigned)).(K), true
	case uint:
		if signedOK {
			return any(uint(signed)).(K), true
		}
		return any(uint(unsigned)).(K), true
	case uint8:
		if signedOK {
			return any(uint8(signed)).(K), true
		}
		return any(uint8(unsigned)).(K), true
	case uint16:
		if signedOK {
			return any(uint16(signed)).(K), true
		}
		return any(uint16(unsigned)).(K), true
	case uint32:
		if signedOK {
			return any(uint32(signed)).(K), true
		}
		return any(uint32(unsigned)).(K), true
	case uint64:
		if signedOK {
			return any(uint64(signed)).(K), true
		}
		return any(unsigned).(K), true
	case uintptr:
		if signedOK {
			return any(uintptr(signed)).(K), true
		}
		return any(uintptr(unsigned)).(K), true
	default:
		return zero[K](), false
	}
}

func asSignedInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	default:
		return 0, false
	}
}

func asUnsignedUint64(value any) (uint64, bool) {
	switch v := value.(type) {
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case uintptr:
		return uint64(v), true
	default:
		return 0, false
	}
}

func keyToBytes(value any) ([]byte, error) {
	switch v := value.(type) {
	case string:
		if len(v) == 0 {
			return []byte{}, nil
		}
		return unsafe.Slice(unsafe.StringData(v), len(v)), nil
	case []byte:
		return v, nil
	case int:
		return encodeSigned(int64(v), bits.UintSize), nil
	case int8:
		return encodeSigned(int64(v), 8), nil
	case int16:
		return encodeSigned(int64(v), 16), nil
	case int32:
		return encodeSigned(int64(v), 32), nil
	case int64:
		return encodeSigned(v, 64), nil
	case uint:
		return encodeUnsigned(uint64(v), bits.UintSize), nil
	case uint8:
		return encodeUnsigned(uint64(v), 8), nil
	case uint16:
		return encodeUnsigned(uint64(v), 16), nil
	case uint32:
		return encodeUnsigned(uint64(v), 32), nil
	case uint64:
		return encodeUnsigned(v, 64), nil
	case uintptr:
		return encodeUnsigned(uint64(v), bits.UintSize), nil
	default:
		return nil, fmt.Errorf("unsupported key type %T", value)
	}
}

func bytesToKey[K Convertible](raw []byte) (K, error) {
	var target K
	switch any(target).(type) {
	case string:
		if len(raw) == 0 {
			return any("").(K), nil
		}
		return any(unsafe.String(unsafe.SliceData(raw), len(raw))).(K), nil
	case []byte:
		return any(bytes.Clone(raw)).(K), nil
	case int:
		decoded, err := decodeSigned(raw, bits.UintSize)
		return any(int(decoded)).(K), err
	case int8:
		decoded, err := decodeSigned(raw, 8)
		return any(int8(decoded)).(K), err
	case int16:
		decoded, err := decodeSigned(raw, 16)
		return any(int16(decoded)).(K), err
	case int32:
		decoded, err := decodeSigned(raw, 32)
		return any(int32(decoded)).(K), err
	case int64:
		decoded, err := decodeSigned(raw, 64)
		return any(decoded).(K), err
	case uint:
		decoded, err := decodeUnsigned(raw, bits.UintSize)
		return any(uint(decoded)).(K), err
	case uint8:
		decoded, err := decodeUnsigned(raw, 8)
		return any(uint8(decoded)).(K), err
	case uint16:
		decoded, err := decodeUnsigned(raw, 16)
		return any(uint16(decoded)).(K), err
	case uint32:
		decoded, err := decodeUnsigned(raw, 32)
		return any(uint32(decoded)).(K), err
	case uint64:
		decoded, err := decodeUnsigned(raw, 64)
		return any(decoded).(K), err
	case uintptr:
		decoded, err := decodeUnsigned(raw, bits.UintSize)
		return any(uintptr(decoded)).(K), err
	default:
		return zero[K](), fmt.Errorf("unsupported key conversion target %T", target)
	}
}

func DefaultForwardConvertValue[T0, T1 any](value T0) (T1, error) {
	if converted, ok := any(value).(T1); ok {
		return converted, nil
	}

	var target T1
	targetType := reflect.TypeOf(target)
	if targetType == nil {
		return zero[T1](), fmt.Errorf("unsupported value conversion from %T to %T", value, target)
	}

	sourceValue := reflect.ValueOf(value)
	if !sourceValue.IsValid() {
		return zero[T1](), fmt.Errorf("unsupported value conversion from <invalid> to %T", target)
	}

	if sourceValue.Type().ConvertibleTo(targetType) {
		converted := sourceValue.Convert(targetType).Interface().(T1)
		return converted, nil
	}

	if targetType.Kind() == reflect.Slice && targetType.Elem().Kind() == reflect.Uint8 {
		raw, err := valueToBytes(sourceValue)
		if err != nil {
			return zero[T1](), err
		}
		return any(raw).(T1), nil
	}

	if sourceBytes, ok := any(value).([]byte); ok {
		decoded, err := bytesToValue[T1](sourceBytes)
		if err != nil {
			return zero[T1](), err
		}
		return decoded, nil
	}

	return zero[T1](), fmt.Errorf("unsupported value conversion from %T to %T", value, target)
}

func valueToBytes(value reflect.Value) ([]byte, error) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeSigned(value.Int(), int(value.Type().Bits())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return encodeUnsigned(value.Uint(), int(value.Type().Bits())), nil
	case reflect.Float32:
		return encodeUnsigned(uint64(math.Float32bits(float32(value.Float()))), 32), nil
	case reflect.Float64:
		return encodeUnsigned(math.Float64bits(value.Float()), 64), nil
	default:
		return nil, fmt.Errorf("unsupported value conversion from %T to []byte", value.Interface())
	}
}

func bytesToValue[T any](raw []byte) (T, error) {
	var target T
	targetType := reflect.TypeOf(target)
	if targetType == nil {
		return zero[T](), fmt.Errorf("unsupported []byte conversion target %T", target)
	}

	value := reflect.New(targetType).Elem()
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		decoded, err := decodeSigned(raw, int(value.Type().Bits()))
		if err != nil {
			return zero[T](), err
		}
		value.SetInt(decoded)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		decoded, err := decodeUnsigned(raw, int(value.Type().Bits()))
		if err != nil {
			return zero[T](), err
		}
		value.SetUint(decoded)
	case reflect.Float32:
		decoded, err := decodeUnsigned(raw, 32)
		if err != nil {
			return zero[T](), err
		}
		value.SetFloat(float64(math.Float32frombits(uint32(decoded))))
	case reflect.Float64:
		decoded, err := decodeUnsigned(raw, 64)
		if err != nil {
			return zero[T](), err
		}
		value.SetFloat(math.Float64frombits(decoded))
	case reflect.String:
		if len(raw) == 0 {
			value.SetString("")
		} else {
			value.SetString(unsafe.String(unsafe.SliceData(raw), len(raw)))
		}
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			value.SetBytes(bytes.Clone(raw))
		} else {
			return zero[T](), fmt.Errorf("unsupported []byte conversion target %T", target)
		}
	default:
		return zero[T](), fmt.Errorf("unsupported []byte conversion target %T", target)
	}

	converted, ok := value.Interface().(T)
	if !ok {
		return zero[T](), fmt.Errorf("unsupported []byte conversion target %T", target)
	}
	return converted, nil
}

func encodeSigned(value int64, bits int) []byte {
	return encodeUnsigned(uint64(value), bits)
}

func encodeUnsigned(value uint64, bits int) []byte {
	width := bits / 8
	if width == 0 {
		width = 1
	}

	buf := make([]byte, width)
	switch width {
	case 1:
		buf[0] = byte(value)
	case 2:
		binary.LittleEndian.PutUint16(buf, uint16(value))
	case 4:
		binary.LittleEndian.PutUint32(buf, uint32(value))
	default:
		binary.LittleEndian.PutUint64(buf, value)
	}
	return buf
}

func decodeSigned(raw []byte, bits int) (int64, error) {
	decoded, err := decodeUnsigned(raw, bits)
	if err != nil {
		return 0, err
	}
	shift := 64 - bits
	return int64(decoded<<shift) >> shift, nil
}

func decodeUnsigned(raw []byte, bits int) (uint64, error) {
	width := bits / 8
	if width == 0 {
		width = 1
	}
	if len(raw) != width {
		return 0, fmt.Errorf("invalid numeric encoding length %d for %d-bit value", len(raw), bits)
	}

	switch width {
	case 1:
		return uint64(raw[0]), nil
	case 2:
		return uint64(binary.LittleEndian.Uint16(raw)), nil
	case 4:
		return uint64(binary.LittleEndian.Uint32(raw)), nil
	default:
		return binary.LittleEndian.Uint64(raw), nil
	}
}

func zero[T any]() T {
	var v T
	return v
}

// func IsNil[T any](value T) bool {
// 	if value == nil {
// 		return true
// 	}

// }
