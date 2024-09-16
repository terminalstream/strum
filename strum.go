package strum

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const defaultDelimiter = ","

// StringUnmarshaller unmarshals strings.
//
// If StringUnmarshaller sees the field is tagged with 'strum' it assigns the indicated substring, otherwise
// the field is ignored.
//
// 'strum' has the format `strum:"startIdx{delimiter}endIdx"` where both startIdx and endIdx are optional, but at least
// one must be present. {delimiter} is specified by the user (default is ","). {delimiter} is mandatory unless only
// startIdx is provided. Errors are raised if startIdx or endIdx exceed the string's bounds.
type StringUnmarshaller struct {
	// TODO: support complex datatypes (types that implement func UnmarshalString(string, any) error, also fields
	//       of complex datatypes that are tagged with linePos)

	Delimiter string
}

func (l *StringUnmarshaller) UnmarshalString(line string, v any) error {
	if l.Delimiter == "" {
		l.Delimiter = defaultDelimiter
	}

	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer: %s", reflect.TypeOf(v).Kind())
	}

	if value.IsNil() {
		return errors.New("nil pointer")
	}

	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("not a struct: %s", reflect.TypeOf(v).Elem().Kind())
	}

	value = value.Elem()

	t := value.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := value.Field(i)

		if !fv.CanSet() {
			return fmt.Errorf("cannot assign any value to field %q", f.Name)
		}

		if fv.Kind() != reflect.String {
			return fmt.Errorf("field %q is not a string type", f.Name)
		}

		pos, ok := f.Tag.Lookup("linePos")
		if !ok {
			continue
		}

		parts := strings.Split(pos, l.Delimiter)

		if len(parts) > 2 {
			return fmt.Errorf("invalid index format on field %q: %q", f.Name, pos)
		}

		var (
			startIdx int
			endIdx   int
			err      error
		)

		if len(parts) > 0 {
			if parts[0] != "" {
				startIdx, err = strconv.Atoi(parts[0])
				if err != nil {
					return fmt.Errorf("invalid start index %q field %q: %w", parts[0], f.Name, err)
				}
			}

			if startIdx > len(line) {
				return fmt.Errorf("start index %d greater than line length %d on field %q", startIdx, len(line), f.Name)
			}
		}

		if len(parts) > 1 {
			endIdx, err = strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid end index %q on field %q: %w", parts[1], f.Name, err)
			}

			if endIdx < startIdx {
				return fmt.Errorf("end index smaller than start index on field %q", f.Name)
			}

			if endIdx > len(line) {
				return fmt.Errorf("end index %d greater than line length %d on field %q", endIdx, len(line), f.Name)
			}
		} else {
			endIdx = len(line)
		}

		fv.SetString(line[startIdx:endIdx])
	}

	return nil
}
