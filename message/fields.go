package message

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
)

type Field interface {
	Get() any
	String() (string, error)
	Set(any) error
	SetBytes([]byte) error
	Log() (interface{}, error)
}

type FieldAccessor struct {
	field Field
	err   error
}

func (fa *FieldAccessor) String() (string, error) {
	if fa.err != nil {
		return "", fa.err
	}

	return fa.field.String()
}

func (fa *FieldAccessor) Log() (interface{}, error) {
	if fa.err != nil {
		return "", fa.err
	}

	return fa.field.Log()
}

type StringField struct {
	val string
}

func (f *StringField) Get() any { return f.val }

func (f *StringField) Set(v any) error {
	switch val := v.(type) {
	case string:
		f.val = val
	case []byte:
		f.val = string(val)
	default:
		// Opción: convertir cualquier cosa a string
		f.val = fmt.Sprintf("%v", val)
	}
	return nil
}

func (f *StringField) String() (string, error) {
	return f.val, nil
}

func (f *StringField) SetBytes(b []byte) error {
	f.val = string(b)
	return nil
}

type IntField struct {
	val int
}

func (f *StringField) Log() (interface{}, error) { return f.String() }

func (f *IntField) Get() any { return f.val }

func (f *IntField) Set(v any) error {
	switch val := v.(type) {
	case int:
		f.val = val
	case int64:
		f.val = int(val)
	case string:
		i, err := strconv.Atoi(val)
		if err == nil {
			f.val = i
		}
	case []byte:
		i, err := strconv.Atoi(string(val))
		if err == nil {
			f.val = i
		}
	}
	return nil // Coincide con la interfaz
}

func (f *IntField) String() (string, error) {
	return fmt.Sprintf("%d", f.val), nil // Coincide con la interfaz
}

func (f *IntField) SetBytes(b []byte) error {
	i, err := strconv.Atoi(string(b))
	if err == nil {
		f.val = i
	}
	return nil
}

func (f *IntField) Log() (interface{}, error) { return f.String() }

type BytesField struct {
	val []byte
}

func (f *BytesField) Get() any { return f.val }

func (f *BytesField) Set(v any) error {
	switch val := v.(type) {
	case []byte:
		f.val = val
	case string:
		f.val = []byte(val)
	}
	return nil
}

func (f *BytesField) String() (string, error) {
	return hex.EncodeToString(f.val), nil // Coincide con la interfaz
}

func (f *BytesField) SetBytes(b []byte) error {
	f.val = b
	return nil
}

func (f *BytesField) Log() (interface{}, error) { return f.String() }

type CustomPacker interface {
	Pack() (string, error)
	Unpack(string) error
	Log() (interface{}, error)
}

type StructField[T any] struct {
	val T
}

func (f *StructField[T]) Get() any {
	return f.val
}

func (f *StructField[T]) Set(v any) error {
	if val, ok := v.(T); ok {
		f.val = val
		return nil
	}

	return fmt.Errorf("invalid type StructField: received %T, expected %T", v, f.val)
}

func (f *StructField[T]) Log() (interface{}, error) { return f.String() }

func (f *StructField[T]) packInterno() (string, error) {
	if p, ok := any(&f.val).(CustomPacker); ok {
		return p.Pack()
	}
	bytes, err := json.Marshal(f.val)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (f *StructField[T]) unpackInterno(data string) error {
	if p, ok := any(&f.val).(CustomPacker); ok {
		return p.Unpack(data)
	}
	return json.Unmarshal([]byte(data), &f.val)
}

func (f *StructField[T]) String() (string, error) {
	str, err := f.packInterno()
	if err != nil {
		return "", fmt.Errorf("serialization failed: %w", err)
	}
	return str, nil
}

func (f *StructField[T]) SetBytes(b []byte) error {
	err := f.unpackInterno(string(b))
	if err != nil {
		return fmt.Errorf("serialization failed: %w", err)
	}
	return nil
}
