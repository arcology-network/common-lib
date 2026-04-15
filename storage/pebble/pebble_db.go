/*
 *   Copyright (c) 2026 Arcology Network
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package pebbledb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"reflect"
	"unsafe"

	"github.com/cockroachdb/pebble"
)

type Key interface {
	string | []byte | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

type PebbleDB[K Key, T any] struct {
	impl       *pebble.DB
	keyEncoder func(K) ([]byte, error)
	keyDecoder func([]byte) (K, error)
	encoder    func(K, T) ([]byte, error)
	decoder    func(K, []byte, T) any
}

func NewPebbleDB[K Key, T any](path string) *PebbleDB[K, T] {
	return NewPebbleDBWithCodec[K, T](path, nil, nil)
}

func NewPebbleDBWithCodec[K Key, T any](
	path string,
	encoder func(K, T) ([]byte, error),
	decoder func(K, []byte, T) any,
) *PebbleDB[K, T] {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		panic(err)
	}

	if encoder == nil {
		encoder = defaultEncoder[K, T]
	}
	if decoder == nil {
		decoder = defaultDecoder[K, T]
	}

	return &PebbleDB[K, T]{
		impl:       db,
		keyEncoder: defaultKeyEncoder[K],
		keyDecoder: defaultKeyDecoder[K],
		encoder:    encoder,
		decoder:    decoder,
	}
}

func (db *PebbleDB[K, T]) Get(key K) (value T, err error) {
	encodedKey, err := db.keyEncoder(key)
	if err != nil {
		return zero[T](), err
	}

	stored, closer, err := db.impl.Get(encodedKey)
	if err != nil {
		return zero[T](), err
	}
	defer func() {
		closeErr := closer.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	return db.decodeTypedValue(key, bytes.Clone(stored), zero[T]())
}

func (db *PebbleDB[K, T]) GetAs(key K, target T) (any, error) {
	encodedKey, err := db.keyEncoder(key)
	if err != nil {
		return nil, err
	}

	stored, closer, err := db.impl.Get(encodedKey)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer func() {
		closeErr := closer.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	return db.decoder(key, bytes.Clone(stored), target), nil
}

func (db *PebbleDB[K, T]) Has(key K) bool {
	encodedKey, err := db.keyEncoder(key)
	if err != nil {
		return false
	}

	_, closer, err := db.impl.Get(encodedKey)
	if err != nil {
		return false
	}
	return closer.Close() == nil
}

func (db *PebbleDB[K, T]) Set(key K, value T) error {
	encodedKey, err := db.keyEncoder(key)
	if err != nil {
		return err
	}

	encoded, err := db.encoder(key, value)
	if err != nil {
		return err
	}
	return db.impl.Set(encodedKey, encoded, pebble.NoSync)
}

func (db *PebbleDB[K, T]) Delete(key K) error {
	encodedKey, err := db.keyEncoder(key)
	if err != nil {
		return err
	}
	return db.impl.Delete(encodedKey, pebble.NoSync)
}

func (db *PebbleDB[K, T]) DeleteBatch(keys []K) error {
	batch := db.impl.NewBatch()
	defer batch.Close()

	for _, key := range keys {
		encodedKey, err := db.keyEncoder(key)
		if err != nil {
			return err
		}
		if err := batch.Delete(encodedKey, pebble.NoSync); err != nil {
			return err
		}
	}
	return batch.Commit(pebble.NoSync)
}

func (db *PebbleDB[K, T]) GetBatch(keys []K) (values []T, err error) {
	values = make([]T, len(keys))
	for i, key := range keys {
		encodedKey, keyErr := db.keyEncoder(key)
		if keyErr != nil {
			return values, keyErr
		}

		stored, closer, getErr := db.impl.Get(encodedKey)
		if getErr != nil {
			if getErr == pebble.ErrNotFound {
				continue
			}
			return values, getErr
		}

		decoded, decodeErr := db.decodeTypedValue(key, bytes.Clone(stored), zero[T]())
		closeErr := closer.Close()
		if decodeErr != nil {
			return values, decodeErr
		}
		if closeErr != nil {
			return values, closeErr
		}
		values[i] = decoded
	}
	return values, nil
}

func (db *PebbleDB[K, T]) SetBatch(keys []K, values []T) error {
	batch := db.impl.NewBatch()
	defer batch.Close()

	for i := range keys {
		encodedKey, keyErr := db.keyEncoder(keys[i])
		if keyErr != nil {
			return keyErr
		}

		encoded, err := db.encoder(keys[i], values[i])
		if err != nil {
			return err
		}
		if err := batch.Set(encodedKey, encoded, pebble.NoSync); err != nil {
			return err
		}
	}
	return batch.Commit(pebble.NoSync)
}

func (db *PebbleDB[K, T]) Query(prefix K, checker func(K, T) bool) (keys []K, values []T, err error) {
	encodedPrefix, err := db.keyEncoder(prefix)
	if err != nil {
		return nil, nil, err
	}

	iter, err := db.impl.NewIter(&pebble.IterOptions{})
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		closeErr := iter.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	for iter.First(); iter.Valid(); iter.Next() {
		if len(encodedPrefix) > 0 && !bytes.HasPrefix(iter.Key(), encodedPrefix) {
			continue
		}

		key, keyErr := db.keyDecoder(bytes.Clone(iter.Key()))
		if keyErr != nil {
			return keys, values, keyErr
		}

		value, decodeErr := db.decodeTypedValue(key, bytes.Clone(iter.Value()), zero[T]())
		if decodeErr != nil {
			return keys, values, decodeErr
		}
		if checker != nil && !checker(key, value) {
			continue
		}

		keys = append(keys, key)
		values = append(values, value)
	}
	if iterErr := iter.Error(); iterErr != nil {
		return keys, values, iterErr
	}
	return keys, values, nil
}

func (db *PebbleDB[K, T]) Close() error {
	return db.impl.Close()
}

func (db *PebbleDB[K, T]) decodeTypedValue(key K, raw []byte, target T) (T, error) {
	decoded := db.decoder(key, raw, target)
	if decoded == nil {
		return zero[T](), nil
	}

	value, ok := decoded.(T)
	if !ok {
		return zero[T](), fmt.Errorf("pebbledb decoder returned %T, expected %T", decoded, zero[T]())
	}
	return value, nil
}

func defaultKeyEncoder[K Key](key K) ([]byte, error) {
	switch v := any(key).(type) {
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
		return nil, fmt.Errorf("unsupported key type %T", key)
	}
}

func defaultKeyDecoder[K Key](raw []byte) (K, error) {
	switch any(*new(K)).(type) {
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
		return zero[K](), fmt.Errorf("unsupported key type %T", *new(K))
	}
}

func defaultEncoder[K Key, T any](_ K, value T) ([]byte, error) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeSigned(rv.Int(), rv.Type().Bits()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return encodeUnsigned(rv.Uint(), rv.Type().Bits()), nil
	case reflect.Float32:
		return encodeUnsigned(uint64(math.Float32bits(float32(rv.Float()))), 32), nil
	case reflect.Float64:
		return encodeUnsigned(math.Float64bits(rv.Float()), 64), nil
	default:
		return nil, fmt.Errorf("default pebble encoder only supports numeric values, got %T", value)
	}
}

func defaultDecoder[K Key, T any](_ K, value []byte, target T) any {
	if value == nil {
		return zero[T]()
	}

	rv := reflect.New(reflect.TypeOf(target)).Elem()
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		decoded, err := decodeSigned(value, rv.Type().Bits())
		if err != nil {
			return nil
		}
		rv.SetInt(decoded)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		decoded, err := decodeUnsigned(value, rv.Type().Bits())
		if err != nil {
			return nil
		}
		rv.SetUint(decoded)
	case reflect.Float32:
		decoded, err := decodeUnsigned(value, 32)
		if err != nil {
			return nil
		}
		rv.SetFloat(float64(math.Float32frombits(uint32(decoded))))
	case reflect.Float64:
		decoded, err := decodeUnsigned(value, 64)
		if err != nil {
			return nil
		}
		rv.SetFloat(math.Float64frombits(decoded))
	default:
		return nil
	}

	decoded, ok := rv.Interface().(T)
	if !ok {
		return nil
	}
	return decoded
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
