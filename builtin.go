// Copyright 2024 Terminal Stream Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strum

import (
	"reflect"
	"strconv"
)

type primitiveValuer func(s string) (reflect.Value, error)

//nolint:godox
// TODO support []byte

var builtin = map[reflect.Kind]primitiveValuer{
	reflect.Bool: func(s string) (reflect.Value, error) {
		return valueOrError(strconv.ParseBool(s))
	},
	reflect.Int: func(s string) (reflect.Value, error) {
		return valueOrError(strconv.Atoi(s))
	},
	reflect.Int8: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 8)

		return valueOrError(castInt[int8](n, err))
	},
	reflect.Int16: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 16)

		return valueOrError(castInt[int16](n, err))
	},
	reflect.Int32: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 32)

		return valueOrError(castInt[int32](n, err))
	},
	reflect.Int64: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 64)

		return valueOrError(castInt[int64](n, err))
	},
	reflect.Uint: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 64)

		return valueOrError(castUint[uint](n, err))
	},
	reflect.Uint8: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 8)

		return valueOrError(castUint[uint8](n, err))
	},
	reflect.Uint16: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 16)

		return valueOrError(castUint[uint16](n, err))
	},
	reflect.Uint32: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 32)

		return valueOrError(castUint[uint32](n, err))
	},
	reflect.Uint64: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 64)

		return valueOrError(castUint[uint64](n, err))
	},
	reflect.Float32: func(s string) (reflect.Value, error) {
		f, err := strconv.ParseFloat(s, 32)

		return valueOrError(castFloat[float32](f, err))
	},
	reflect.Float64: func(s string) (reflect.Value, error) {
		f, err := strconv.ParseFloat(s, 64)

		return valueOrError(castFloat[float64](f, err))
	},
	reflect.String: func(s string) (reflect.Value, error) {
		v := reflect.ValueOf(s)

		return v, nil
	},
}

var builtinPointers = map[reflect.Kind]primitiveValuer{
	reflect.Bool: func(s string) (reflect.Value, error) {
		b, err := strconv.ParseBool(s)

		return valueOrError(&b, err)
	},
	reflect.Int: func(s string) (reflect.Value, error) {
		n, err := strconv.Atoi(s)

		return valueOrError(&n, err)
	},
	reflect.Int8: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 8)
		i := int8(n)

		return valueOrError(&i, err)
	},
	reflect.Int16: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 16)
		i := int16(n)

		return valueOrError(&i, err)
	},
	reflect.Int32: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 32)
		i := int32(n)

		return valueOrError(&i, err)
	},
	reflect.Int64: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 64)
		i := int64(n)

		return valueOrError(&i, err)
	},
	reflect.Uint: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 64)
		i := uint(n)

		return valueOrError(&i, err)
	},
	reflect.Uint8: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 8)
		i := uint8(n)

		return valueOrError(&i, err)
	},
	reflect.Uint16: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 16)
		i := uint16(n)

		return valueOrError(&i, err)
	},
	reflect.Uint32: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 32)
		i := uint32(n)

		return valueOrError(&i, err)
	},
	reflect.Uint64: func(s string) (reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 64)

		return valueOrError(&n, err)
	},
	reflect.Float32: func(s string) (reflect.Value, error) {
		f, err := strconv.ParseFloat(s, 32)
		i := float32(f)

		return valueOrError(&i, err)
	},
	reflect.Float64: func(s string) (reflect.Value, error) {
		f, err := strconv.ParseFloat(s, 64)

		return valueOrError(&f, err)
	},
	reflect.String: func(s string) (reflect.Value, error) {
		v := reflect.ValueOf(&s)

		return v, nil
	},
}

func valueOrError(e any, err error) (reflect.Value, error) {
	if err != nil {
		return reflect.Value{}, err
	}

	v := reflect.ValueOf(e)

	return v, nil
}

type numeric interface {
	int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func castInt[T numeric](n int64, err error) (T, error) {
	if err != nil {
		return *new(T), err
	}

	return T(n), nil
}

func castUint[T numeric](n uint64, err error) (T, error) {
	if err != nil {
		return *new(T), err
	}

	return T(n), nil
}

func castFloat[T numeric](n float64, err error) (T, error) {
	if err != nil {
		return *new(T), err
	}

	return T(n), nil
}
