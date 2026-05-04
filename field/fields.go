// Package field provides the core interfaces and concrete types for representing
// the data within a message field.
package field

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// Field defines the contract for all data-holding field types.
type Field interface {
	Get() any
	Set(any) error
	Log() (interface{}, error)
	Parse() (string, error)
	Unparse(string) error
}

// String is a data field that holds a string value.
type String struct {
	val string
}

func (f *String) Get() any { return f.val }

func (f *String) Set(v any) error {
	if val, ok := v.(string); ok {
		f.val = val
		return nil
	}
	return fmt.Errorf("invalid type for String: received %T, expected string", v)
}

func (f *String) Log() (interface{}, error) { return f.val, nil }

func (f *String) Parse() (string, error) {
	return f.val, nil
}

func (f *String) Unparse(data string) error {
	f.val = data
	return nil
}

// Int is a data field that holds an integer value.
type Int struct {
	val int
}

func (f *Int) Get() any { return f.val }

func (f *Int) Set(v any) error {
	if val, ok := v.(int); ok {
		f.val = val
		return nil
	}
	return fmt.Errorf("invalid type for Int: received %T, expected int", v)
}

func (f *Int) Log() (interface{}, error) { return fmt.Sprintf("%d", f.val), nil }

func (f *Int) Parse() (string, error) {
	return fmt.Sprintf("%d", f.val), nil
}

func (f *Int) Unparse(data string) error {
	val, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}
	f.val = int(val)
	return nil
}

// Bytes is a data field that holds a byte slice.
type Bytes struct {
	val []byte
}

func (f *Bytes) Get() any { return f.val }

func (f *Bytes) Set(v any) error {
	if val, ok := v.([]byte); ok {
		f.val = val
		return nil
	}
	return fmt.Errorf("invalid type for Bytes: received %T, expected []byte", v)
}

func (f *Bytes) Parse() (string, error) {
	return fmt.Sprintf("%X", f.val), nil
}

func (f *Bytes) Unparse(data string) error {
	val, err := hex.DecodeString(data)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}
	f.val = val
	return nil
}

func (f *Bytes) Log() (interface{}, error) { return fmt.Sprintf("%X", f.val), nil }

//// Struct is a generic data field that can hold any custom struct.
//type Struct[T any] struct {
//	val T
//}
//
//func (f *Struct[T]) Get() any {
//	return f.val
//}
//
//func (f *Struct[T]) Set(v any) error {
//	if val, ok := v.(T); ok {
//		f.val = val
//		return nil
//	}
//
//	return fmt.Errorf("invalid type for Struct: received %T, expected %T", v, f.val)
//}
//
//func (f *Struct[T]) LogMsg() (interface{}, error) {
//	if p, ok := any(&f.val).(CustomPacker); ok {
//		return p.LogMsg()
//	}
//
//	return f.String()
//}
//
//func (f *Struct[T]) packInterno() (string, error) {
//	if p, ok := any(&f.val).(CustomPacker); ok {
//		return p.Parse()
//	}
//	bytes, err := json.Marshal(f.val)
//	if err != nil {
//		return "", err
//	}
//	return string(bytes), nil
//}
//
//func (f *Struct[T]) unpackInterno(data string) error {
//	if p, ok := any(&f.val).(CustomPacker); ok {
//		return p.Unparse(data)
//	}
//	return json.Unmarshal([]byte(data), &f.val)
//}
//
//func (f *Struct[T]) String() (string, error) {
//	str, err := f.packInterno()
//	if err != nil {
//		return "", fmt.Errorf("serialization failed: %w", err)
//	}
//	return str, nil
//}
//
//func (f *Struct[T]) SetBytes(b []byte) error {
//	err := f.unpackInterno(string(b))
//	if err != nil {
//		return fmt.Errorf("serialization failed: %w", err)
//	}
//	return nil
//}
