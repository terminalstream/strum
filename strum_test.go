package strum

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_StringUnmarshaller_UnmarshalString(t *testing.T) {
	t.Run("assigns string correctly with startIdx and endIdx", func(t *testing.T) {
		const line = "lkjasldthis is a testoi09asdfhj"

		test := &struct {
			Val string `linePos:"7-21"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("assigns suffix correctly with just startIdx", func(t *testing.T) {
		const line = "lkjasldoi09asdfhjthis is a test"

		test := &struct {
			Val string `linePos:"17"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("assigns suffix correctly with startIdx and endIdx", func(t *testing.T) {
		const line = "lkjasldoi09asdfhjthis is a test"

		test := &struct {
			Val string `linePos:"17-31"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("assigns prefix correctly with just endIdx", func(t *testing.T) {
		const line = "this is a testlkjasldoi09asdfhj"

		test := &struct {
			Val string `linePos:"-14"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("assigns prefix correctly with endIdx and endIdx", func(t *testing.T) {
		const line = "this is a testlkjasldoi09asdfhj"

		test := &struct {
			Val string `linePos:"0-14"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString(line, test)
		require.NoError(t, err)

		require.Equal(t, "this is a test", test.Val)
	})

	t.Run("error when given a non-pointer", func(t *testing.T) {
		err := (&StringUnmarshaller{}).UnmarshalString("", struct{}{})
		require.ErrorContains(t, err, "not a pointer")
	})

	t.Run("error when given a nil pointer", func(t *testing.T) {
		var p *string
		err := (&StringUnmarshaller{}).UnmarshalString("", p)
		require.ErrorContains(t, err, "nil pointer")
	})

	t.Run("error when pointer points to something other than a struct", func(t *testing.T) {
		p := &[]byte{}
		err := (&StringUnmarshaller{}).UnmarshalString("", p)
		require.ErrorContains(t, err, "not a struct")
	})

	t.Run("error when tag value has extra elements", func(t *testing.T) {
		test := &struct {
			Val string `linePos:"0-14-"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("", test)
		require.ErrorContains(t, err, "invalid index format")
	})

	t.Run("error when startIdx has invalid value", func(t *testing.T) {
		test := &struct {
			Val string `linePos:"invalid"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("", test)
		require.ErrorContains(t, err, "invalid start index")
	})

	t.Run("error when endIdx has invalid value", func(t *testing.T) {
		test := &struct {
			Val string `linePos:"-invalid"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("", test)
		require.ErrorContains(t, err, "invalid end index")
	})

	t.Run("error when startIdx > line length", func(t *testing.T) {
		test := &struct {
			Val string `linePos:"1000"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("", test)
		require.ErrorContains(t, err, "start index 1000 greater than line length")
	})

	t.Run("error when startIdx > endIdx", func(t *testing.T) {
		test := &struct {
			Val string `linePos:"2-1"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("      ", test)
		require.ErrorContains(t, err, "end index smaller than start index")
	})

	t.Run("error if startIdx > line length", func(t *testing.T) {
		test := &struct {
			Val string `linePos:"1"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("", test)
		require.ErrorContains(t, err, "start index 1 greater than line length 0")
	})

	t.Run("error if endIdx > line length", func(t *testing.T) {
		test := &struct {
			Val string `linePos:"0-100"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("", test)
		require.ErrorContains(t, err, "end index 100 greater than line length 0")
	})

	t.Run("error if field is not of type string", func(t *testing.T) {
		test := &struct {
			Val int `linePos:"0-1"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("a", test)
		require.ErrorContains(t, err, `field "Val" is not a string type`)
	})

	t.Run("ignores string fields without struct tag", func(t *testing.T) {
		test := &struct {
			Val string `json:"val"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("a", test)
		require.NoError(t, err)
		require.Empty(t, test.Val)
	})

	t.Run("error if field is not exported", func(t *testing.T) {
		test := &struct {
			val string `json:"val"`
		}{}

		err := (&StringUnmarshaller{}).UnmarshalString("a", test)
		require.ErrorContains(t, err, "cannot assign any value to field")
		require.Empty(t, test.val)
	})
}
