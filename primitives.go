package strum

import (
	"reflect"
	"strconv"
)

type primitiveValuer func(s string) (*reflect.Value, error)

// TODO support []byte

var primitiveValuers = map[reflect.Kind]primitiveValuer{
	reflect.Bool: func(s string) (*reflect.Value, error) {
		return valueOrError(strconv.ParseBool(s))
	},
	reflect.Int: func(s string) (*reflect.Value, error) {
		return valueOrError(strconv.Atoi(s))
	},
	reflect.Int8: func(s string) (*reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 8)
		return valueOrError(castInt[int8](n, err))
	},
	reflect.Int16: func(s string) (*reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 16)
		return valueOrError(castInt[int16](n, err))
	},
	reflect.Int32: func(s string) (*reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 32)
		return valueOrError(castInt[int32](n, err))
	},
	reflect.Int64: func(s string) (*reflect.Value, error) {
		n, err := strconv.ParseInt(s, 10, 64)
		return valueOrError(castInt[int64](n, err))
	},
	reflect.Uint: func(s string) (*reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 64)
		return valueOrError(castUint[uint](n, err))
	},
	reflect.Uint8: func(s string) (*reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 8)
		return valueOrError(castUint[uint8](n, err))
	},
	reflect.Uint16: func(s string) (*reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 16)
		return valueOrError(castUint[uint16](n, err))
	},
	reflect.Uint32: func(s string) (*reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 32)
		return valueOrError(castUint[uint32](n, err))
	},
	reflect.Uint64: func(s string) (*reflect.Value, error) {
		n, err := strconv.ParseUint(s, 10, 64)
		return valueOrError(castUint[uint64](n, err))
	},
	reflect.Float32: func(s string) (*reflect.Value, error) {
		f, err := strconv.ParseFloat(s, 32)
		return valueOrError(castFloat[float32](f, err))
	},
	reflect.Float64: func(s string) (*reflect.Value, error) {
		f, err := strconv.ParseFloat(s, 64)
		return valueOrError(castFloat[float64](f, err))
	},
	reflect.String: func(s string) (*reflect.Value, error) {
		v := reflect.ValueOf(s)
		return &v, nil
	},
}

func valueOrError(e any, err error) (*reflect.Value, error) {
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(e)

	return &v, nil
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
