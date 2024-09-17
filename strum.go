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
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	// TagName is this lib's struct tag.
	TagName = "strum"
	// DefaultDelimiter is the default one used to separate the start and end indexes.
	DefaultDelimiter = ","
)

var defaultOptions = &options{
	delimiter: DefaultDelimiter,
}

type options struct {
	delimiter string
}

// Option allows some customization of the Unmarshal process.
type Option func(*options)

// WithDelimiter uses the given delimiter instead of DefaultDelimiter.
func WithDelimiter(delimiter string) Option {
	return func(o *options) {
		o.delimiter = delimiter
	}
}

// Unmarshal unmarshals strings.
//
// If a field is tagged with 'strum' it assigns the indicated substring, otherwise the field is
// ignored.
//
// 'strum' has the format `strum:"startIdx{delimiter}endIdx"` where both startIdx and endIdx are
// optional, but at least one must be present. {delimiter} is specified by the user
// (default is ","). {delimiter} is mandatory unless only startIdx is provided. Errors are raised
// if startIdx or endIdx exceed the string's bounds.
func Unmarshal(line string, v any, opts ...Option) error { //nolint:funlen,gocyclo
	options := *defaultOptions

	for i := range opts {
		opts[i](&options)
	}

	value := reflect.ValueOf(v)

	err := validateInput(v, value)
	if err != nil {
		return err
	}

	value = value.Elem()
	t := value.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := value.Field(i)

		//nolint:godox
		// TODO support pointers to fields

		if !fv.CanSet() {
			//nolint:godox
			// TODO should we skip fields we can't set instead of returning an error?
			return fmt.Errorf("cannot assign any value to field %q", f.Name)
		}

		tagValue, ok := f.Tag.Lookup(TagName)
		if !ok {
			continue
		}

		startIdx, endIdx, err := indexes(tagValue, options.delimiter)
		if err != nil {
			return fmt.Errorf("format error on field %q: %w", f.Name, err)
		}

		if endIdx == -1 {
			endIdx = len(line)
		}

		err = validateIndexes(line, startIdx, endIdx)
		if err != nil {
			return fmt.Errorf("invalid indexes on field %q: %w", f.Name, err)
		}

		valuer, primitive := primitiveValuers[f.Type.Kind()]
		if !primitive {
			continue
		}

		val, err := valuer(line[startIdx:endIdx])
		if err != nil {
			return fmt.Errorf(
				"cannot assign value %q to field %q: %w", line[startIdx:endIdx], f.Name, err,
			)
		}

		fv.Set(*val)
	}

	return nil
}

func indexes(tagValue, delimiter string) (int, int, error) {
	parts := strings.Split(tagValue, delimiter)

	if len(parts) == 0 || len(parts) > 2 {
		return -1, -1, fmt.Errorf("invalid strum format: %q", tagValue)
	}

	var (
		startIdx int
		endIdx   = -1
		err      error
	)

	if parts[0] != "" {
		startIdx, err = strconv.Atoi(parts[0])
		if err != nil {
			return -1, -1, fmt.Errorf("invalid start index %q: %w", parts[0], err)
		}
	}

	if len(parts) > 1 {
		endIdx, err = strconv.Atoi(parts[1])
		if err != nil {
			return -1, -1, fmt.Errorf(`invalid end index "%s": %w`, parts[1], err)
		}
	}

	return startIdx, endIdx, nil
}

func validateInput(v any, value reflect.Value) error {
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer: %s", reflect.ValueOf(v).Kind())
	}

	if value.IsNil() {
		return errors.New("nil pointer")
	}

	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("not a struct: %s", reflect.ValueOf(v).Kind())
	}

	return nil
}

func validateIndexes(line string, startIdx, endIdx int) error {
	if startIdx < 0 || startIdx > len(line) {
		return errors.New("start index out of bounds")
	}

	if endIdx < 0 || endIdx > len(line) {
		return errors.New("end index out of bounds")
	}

	if endIdx < startIdx {
		return errors.New(`end index must be greater or equal to start index`)
	}

	return nil
}
