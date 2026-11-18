// Copyright © 2023 Attestant Limited.
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

package pq

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	jsonMarshalerType   = reflect.TypeFor[json.Marshaler]()
	jsonUnmarshalerType = reflect.TypeFor[json.Unmarshaler]()
	byteType            = reflect.TypeFor[byte]()
)

func jsonFieldName(field *reflect.StructField) string {
	return strings.Split(field.Tag.Get("json"), ",")[0]
}

func marshalStruct(value any) ([]byte, error) {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return []byte("null"), nil
		}

		rv = rv.Elem()
	}

	var buf bytes.Buffer

	buf.WriteByte('{')

	for i := range rv.NumField() {
		field := rv.Type().Field(i)

		name := jsonFieldName(&field)
		if name == "" {
			return nil, fmt.Errorf("%s.%s has no json name", rv.Type(), field.Name)
		}

		encoded, err := encodeValue(rv.Field(i))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		if i > 0 {
			buf.WriteByte(',')
		}

		buf.WriteString(strconv.Quote(name))
		buf.WriteByte(':')
		buf.Write(encoded)
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func encodeValue(rv reflect.Value) ([]byte, error) {
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return []byte("null"), nil
		}

		if rv.Type().Implements(jsonMarshalerType) {
			return rv.Interface().(json.Marshaler).MarshalJSON() //nolint:forcetypeassert // checked above
		}

		return encodeValue(rv.Elem())
	}

	if reflect.PointerTo(rv.Type()).Implements(jsonMarshalerType) {
		ptr := reflect.New(rv.Type())
		ptr.Elem().Set(rv)

		return ptr.Interface().(json.Marshaler).MarshalJSON() //nolint:forcetypeassert // checked above
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Type().Elem() == byteType {
			data := make([]byte, rv.Len())
			reflect.Copy(reflect.ValueOf(data), rv)

			return []byte(strconv.Quote("0x" + hex.EncodeToString(data))), nil
		}

		var buf bytes.Buffer

		buf.WriteByte('[')

		for i := range rv.Len() {
			encoded, err := encodeValue(rv.Index(i))
			if err != nil {
				return nil, err
			}

			if i > 0 {
				buf.WriteByte(',')
			}

			buf.Write(encoded)
		}

		buf.WriteByte(']')

		return buf.Bytes(), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.Quote(strconv.FormatUint(rv.Uint(), 10))), nil
	case reflect.Bool:
		return []byte(strconv.FormatBool(rv.Bool())), nil
	default:
		return nil, fmt.Errorf("unsupported json type %s", rv.Type())
	}
}

func unmarshalStruct(input []byte, value any) error {
	rv := reflect.ValueOf(value).Elem()

	fields := make(map[string]json.RawMessage, rv.NumField())
	if err := json.Unmarshal(input, &fields); err != nil {
		return fmt.Errorf("invalid JSON for %s: %w", rv.Type(), err)
	}

	for i := range rv.NumField() {
		field := rv.Type().Field(i)
		name := jsonFieldName(&field)

		raw, ok := fields[name]
		if !ok {
			return fmt.Errorf("%s: %s missing", rv.Type(), name)
		}

		if err := decodeValue(raw, rv.Field(i)); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}

	return nil
}

func decodeValue(raw json.RawMessage, rv reflect.Value) error {
	if rv.Kind() == reflect.Pointer {
		if string(bytes.TrimSpace(raw)) == "null" {
			rv.SetZero()

			return nil
		}

		rv.Set(reflect.New(rv.Type().Elem()))

		if rv.Type().Implements(jsonUnmarshalerType) {
			return rv.Interface().(json.Unmarshaler).UnmarshalJSON(raw) //nolint:forcetypeassert // checked above
		}

		return decodeValue(raw, rv.Elem())
	}

	if reflect.PointerTo(rv.Type()).Implements(jsonUnmarshalerType) {
		return rv.Addr().Interface().(json.Unmarshaler).UnmarshalJSON(raw) //nolint:forcetypeassert // checked above
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Type().Elem() == byteType {
			return decodeBytes(raw, rv)
		}

		var items []json.RawMessage
		if err := json.Unmarshal(raw, &items); err != nil {
			return err
		}

		if rv.Kind() == reflect.Slice {
			rv.Set(reflect.MakeSlice(rv.Type(), len(items), len(items)))
		} else if len(items) != rv.Len() {
			return fmt.Errorf("expected %d items, got %d", rv.Len(), len(items))
		}

		for i, item := range items {
			if err := decodeValue(item, rv.Index(i)); err != nil {
				return fmt.Errorf("item %d: %w", i, err)
			}
		}

		return nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var text string
		if err := json.Unmarshal(raw, &text); err != nil {
			return err
		}

		number, err := strconv.ParseUint(text, 10, rv.Type().Bits())
		if err != nil {
			return err
		}

		rv.SetUint(number)

		return nil
	case reflect.Bool:
		var flag bool
		if err := json.Unmarshal(raw, &flag); err != nil {
			return err
		}

		rv.SetBool(flag)

		return nil
	default:
		return fmt.Errorf("unsupported json type %s", rv.Type())
	}
}

func decodeBytes(raw json.RawMessage, rv reflect.Value) error {
	var text string
	if err := json.Unmarshal(raw, &text); err != nil {
		return err
	}

	data, err := hex.DecodeString(strings.TrimPrefix(text, "0x"))
	if err != nil {
		return err
	}

	if rv.Kind() == reflect.Slice {
		rv.Set(reflect.MakeSlice(rv.Type(), len(data), len(data)))
	} else if len(data) != rv.Len() {
		return fmt.Errorf("expected %d bytes, got %d", rv.Len(), len(data))
	}

	reflect.Copy(rv, reflect.ValueOf(data))

	return nil
}
