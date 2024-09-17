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

package strum_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terminalstream/strum"
)

func TestUnmarshal_indexes(t *testing.T) { //nolint:funlen
	t.Run("assigns string correctly with startIdx and endIdx", func(t *testing.T) {
		const line = "lkjasldthis is a testoi09asdfhj"

		test := &struct {
			Val string `strum:"7,21"`
		}{}

		err := strum.Unmarshal(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("assigns suffix correctly with just startIdx", func(t *testing.T) {
		const line = "lkjasldoi09asdfhjthis is a test"

		test := &struct {
			Val string `strum:"17"`
		}{}

		err := strum.Unmarshal(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("assigns suffix correctly with startIdx and endIdx", func(t *testing.T) {
		const line = "lkjasldoi09asdfhjthis is a test"

		test := &struct {
			Val string `strum:"17,31"`
		}{}

		err := strum.Unmarshal(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("assigns prefix correctly with just endIdx", func(t *testing.T) {
		const line = "this is a testlkjasldoi09asdfhj"

		test := &struct {
			Val string `strum:",14"`
		}{}

		err := strum.Unmarshal(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("assigns prefix correctly with endIdx and endIdx", func(t *testing.T) {
		const line = "this is a testlkjasldoi09asdfhj"

		test := &struct {
			Val string `strum:"0,14"`
		}{}

		err := strum.Unmarshal(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("error when given a non-pointer", func(t *testing.T) {
		err := strum.Unmarshal("", struct{}{})
		require.ErrorContains(t, err, "not a pointer")
	})

	t.Run("error when given a nil pointer", func(t *testing.T) {
		var p *string
		err := strum.Unmarshal("", p)
		require.ErrorContains(t, err, "nil pointer")
	})

	t.Run("error when pointer points to something other than a struct", func(t *testing.T) {
		p := &[]byte{}
		err := strum.Unmarshal("", p)
		require.ErrorContains(t, err, "not a struct")
	})

	t.Run("error when tag value has extra elements", func(t *testing.T) {
		test := &struct {
			Val string `strum:"0,14,"`
		}{}

		err := strum.Unmarshal("", test)
		require.ErrorContains(t, err, "invalid strum format")
	})

	t.Run("error when startIdx has invalid value", func(t *testing.T) {
		test := &struct {
			Val string `strum:"invalid"`
		}{}

		err := strum.Unmarshal("", test)
		require.ErrorContains(t, err, "invalid start index")
	})

	t.Run("error when endIdx has invalid value", func(t *testing.T) {
		test := &struct {
			Val string `strum:",invalid"`
		}{}

		err := strum.Unmarshal("", test)
		require.ErrorContains(t, err, "invalid end index")
	})

	t.Run("error when startIdx > line length", func(t *testing.T) {
		test := &struct {
			Val string `strum:"1000"`
		}{}

		err := strum.Unmarshal("", test)
		require.ErrorContains(t, err, "out of bounds")
	})

	t.Run("error when startIdx > endIdx", func(t *testing.T) {
		test := &struct {
			Val string `strum:"2,1"`
		}{}

		err := strum.Unmarshal("      ", test)
		require.ErrorContains(t, err, "end index must be greater or equal to start index")
	})

	t.Run("error if startIdx > line length", func(t *testing.T) {
		test := &struct {
			Val string `strum:"1"`
		}{}

		err := strum.Unmarshal("", test)
		require.ErrorContains(t, err, "start index out of bounds")
	})

	t.Run("error if endIdx > line length", func(t *testing.T) {
		test := &struct {
			Val string `strum:"0,100"`
		}{}

		err := strum.Unmarshal("", test)
		require.ErrorContains(t, err, "end index out of bounds")
	})

	t.Run("ignores string fields without struct tag", func(t *testing.T) {
		test := &struct {
			Val string `json:"val"`
		}{}

		err := strum.Unmarshal("a", test)
		require.NoError(t, err)
		require.Empty(t, test.Val)
	})

	t.Run("error if field is not exported", func(t *testing.T) {
		test := &struct {
			val string `json:"val"`
		}{}

		err := strum.Unmarshal("a", test)
		require.ErrorContains(t, err, "cannot assign any value to field")
		require.Empty(t, test.val)
	})
}

func TestUnmarshal_builtin(t *testing.T) { //nolint:funlen,maintidx
	one := 1

	t.Run("int", func(t *testing.T) {
		test := &struct {
			Val int `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, 1, test.Val)
	})

	t.Run("*int", func(t *testing.T) {
		test := &struct {
			Val *int `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, &one, test.Val)
	})

	t.Run("int8", func(t *testing.T) {
		test := &struct {
			Val int8 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, int8(1), test.Val)
	})

	t.Run("*int8", func(t *testing.T) {
		test := &struct {
			Val *int8 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, int8(one), *test.Val)
	})

	t.Run("int16", func(t *testing.T) {
		test := &struct {
			Val int16 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, int16(1), test.Val)
	})

	t.Run("*int16", func(t *testing.T) {
		test := &struct {
			Val *int16 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, int16(1), *test.Val)
	})

	t.Run("int32", func(t *testing.T) {
		test := &struct {
			Val int32 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, int32(1), test.Val)
	})

	t.Run("*int32", func(t *testing.T) {
		test := &struct {
			Val *int32 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, int32(1), *test.Val)
	})

	t.Run("int64", func(t *testing.T) {
		test := &struct {
			Val int64 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, int64(1), test.Val)
	})

	t.Run("*int64", func(t *testing.T) {
		test := &struct {
			Val *int64 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, int64(1), *test.Val)
	})

	t.Run("uint", func(t *testing.T) {
		test := &struct {
			Val uint `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint(1), test.Val)
	})

	t.Run("*uint", func(t *testing.T) {
		test := &struct {
			Val *uint `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint(1), *test.Val)
	})

	t.Run("uint8", func(t *testing.T) {
		test := &struct {
			Val uint8 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint8(1), test.Val)
	})

	t.Run("*uint8", func(t *testing.T) {
		test := &struct {
			Val *uint8 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint8(1), *test.Val)
	})

	t.Run("uint16", func(t *testing.T) {
		test := &struct {
			Val uint16 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint16(1), test.Val)
	})

	t.Run("*uint16", func(t *testing.T) {
		test := &struct {
			Val *uint16 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint16(1), *test.Val)
	})

	t.Run("uint32", func(t *testing.T) {
		test := &struct {
			Val uint32 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint32(1), test.Val)
	})

	t.Run("*uint32", func(t *testing.T) {
		test := &struct {
			Val *uint32 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint32(1), *test.Val)
	})

	t.Run("uint64", func(t *testing.T) {
		test := &struct {
			Val uint64 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint64(1), test.Val)
	})

	t.Run("*uint64", func(t *testing.T) {
		test := &struct {
			Val *uint64 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, uint64(1), *test.Val)
	})

	t.Run("float32", func(t *testing.T) {
		test := &struct {
			Val float32 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, float32(1), test.Val)
	})

	t.Run("*float32", func(t *testing.T) {
		test := &struct {
			Val *float32 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, float32(1), *test.Val)
	})

	t.Run("float64", func(t *testing.T) {
		test := &struct {
			Val float64 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, float64(1), test.Val)
	})

	t.Run("*float64", func(t *testing.T) {
		test := &struct {
			Val *float64 `strum:"0"`
		}{}

		err := strum.Unmarshal("1", test)
		require.NoError(t, err)
		require.Equal(t, float64(1), *test.Val)
	})

	t.Run("bool", func(t *testing.T) {
		test := &struct {
			Val bool `strum:"0"`
		}{}

		err := strum.Unmarshal("true", test)
		require.NoError(t, err)
		require.True(t, test.Val)
	})

	t.Run("*bool", func(t *testing.T) {
		test := &struct {
			Val *bool `strum:"0"`
		}{}

		err := strum.Unmarshal("true", test)
		require.NoError(t, err)
		require.True(t, *test.Val)
	})

	t.Run("string", func(t *testing.T) {
		test := &struct {
			Val string `strum:"0"`
		}{}

		err := strum.Unmarshal("abc", test)
		require.NoError(t, err)
		require.Equal(t, "abc", test.Val)
	})

	t.Run("*string", func(t *testing.T) {
		test := &struct {
			Val *string `strum:"0"`
		}{}

		err := strum.Unmarshal("abc", test)
		require.NoError(t, err)
		require.Equal(t, "abc", *test.Val)
	})

	t.Run("[]byte", func(t *testing.T) {
		test := &struct {
			Val []byte `strum:"0"`
		}{}

		err := strum.Unmarshal("abc", test)
		require.NoError(t, err)
		require.Equal(t, []byte("abc"), test.Val)
	})
}

func TestUnmarshal_formatter(t *testing.T) {
	t.Run("formats the string prior to parsing and assigning", func(t *testing.T) {
		test := &struct {
			Val string `strform:"test" strum:"0"`
		}{}

		err := strum.Unmarshal("abcdefg", test,
			strum.WithFormatter("test", func(s string) (string, error) {
				return strings.ToUpper(s), nil
			}),
		)
		require.NoError(t, err)
		require.Equal(t, "ABCDEFG", test.Val)
	})

	t.Run("returns an error if the formatter returns an error", func(t *testing.T) {
		expected := errors.New("test")

		test := &struct {
			Val string `strform:"test" strum:"0"`
		}{}

		err := strum.Unmarshal("abcdefg", test,
			strum.WithFormatter("test", func(s string) (string, error) {
				return "", expected
			}),
		)
		require.ErrorIs(t, err, expected)
	})
}

func TestUnmarshal_delimiter(t *testing.T) {
	test := &struct {
		Val string `strum:"1-3"`
	}{}

	err := strum.Unmarshal("abcde", test, strum.WithDelimiter("-"))
	require.NoError(t, err)
	require.Equal(t, "bc", test.Val)
}
