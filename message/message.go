// Package message provides functionalities for packing and unpacking
// ISO 8583 messages.
package message

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tomasdemarco/iso8583/bitmap"
	"github.com/tomasdemarco/iso8583/field"
	"github.com/tomasdemarco/iso8583/header"
	"github.com/tomasdemarco/iso8583/packager"
	"github.com/tomasdemarco/iso8583/utils"
	"log"
	"strconv"
)

// Message represents an ISO 8583 message, containing its structure,
// fields, and the associated packager.
type Message struct {
	Packager *packager.Packager
	Length   int
	Header   header.Header
	Trailer  interface{}
	Bitmap   *utils.BitSet
	fields   map[int]field.Field
}

// NewMessage creates and returns a new Message instance
// initialized with the provided packager.
func NewMessage(p *packager.Packager) *Message {
	return &Message{
		Packager: p,
		Bitmap:   utils.NewBitSet(64, 128),
		fields:   make(map[int]field.Field),
	}
}

// SetField sets the value of a specific field in the message.
func (m *Message) SetField(id int, v any) error {
	f, ok := m.fields[id]
	if !ok {
		f = m.createField(id)
		m.fields[id] = f
	}

	if err := f.Set(v); err != nil {
		return fmt.Errorf("failed to set field %d: %w", id, err)
	}
	m.Bitmap.Set(id)
	return nil
}

func (m *Message) SetFieldString(id int, value string) error {
	return m.SetField(id, value)
}

func (m *Message) SetFieldInt(id int, value int) error {
	return m.SetField(id, value)
}

func (m *Message) SetFieldBytes(id int, value []byte) error {
	return m.SetField(id, value)
}

func (m *Message) GetField(id int) (v any, err error) {
	f, ok := m.fields[id]
	if !ok {
		return nil, fmt.Errorf("field %d not found", id)
	}

	return f.Get(), nil
}

func (m *Message) GetFieldString(id int) (string, error) {
	f, ok := m.fields[id]
	if !ok {
		return "", fmt.Errorf("field %d not found", id)
	}

	if val, ok := f.(*field.String); ok {
		return val.Get().(string), nil
	}

	return "", fmt.Errorf("field %d is not a string field, it is %T", id, f)
}

func (m *Message) GetFieldInt(id int) (int, error) {
	f, ok := m.fields[id]
	if !ok {
		return 0, fmt.Errorf("field %d not found", id)
	}

	if val, ok := f.(*field.Int); ok {
		return val.Get().(int), nil
	}

	return 0, fmt.Errorf("field %d is not an int field, it is %T", id, f)
}

func (m *Message) GetFieldBytes(id int) ([]byte, error) {
	f, ok := m.fields[id]
	if !ok {
		return nil, fmt.Errorf("field %d not found", id)
	}

	if val, ok := f.(*field.Bytes); ok {
		return val.Get().([]byte), nil
	}

	return nil, fmt.Errorf("field %d is not a byte field, it is %T", id, f)
}

// createField is an internal function to instantiate the correct Field type
// based on the packager configuration.
func (m *Message) createField(id int) field.Field {
	if fieldSpec, ok := m.Packager.Fields[id]; ok {
		return fieldSpec.NewDataField()
	}
	// Default to a string field if no spec is found.
	// This could be made stricter if desired.
	return &field.String{}
}

// Unpack unpacks a byte slice of an ISO 8583 message
// into the Message structure, populating its fields.
func (m *Message) Unpack(messageRaw []byte) (err error) {
	lengthMti, err := m.unpackMti(messageRaw)
	if err != nil {
		return err
	}

	lengthBitmap, err := m.unpackBitmap(messageRaw, lengthMti)
	if err != nil {
		return err
	}

	err = m.unpackFields(messageRaw, lengthMti+lengthBitmap)
	return err
}

// Pack packs the message fields into an ISO 8583 byte slice.
func (m *Message) Pack() ([]byte, error) {
	msgPacked := new(bytes.Buffer)

	encodeField, err := m.packMti()
	if err != nil {
		return nil, err
	}
	msgPacked.Write(encodeField)

	encodeField, err = m.packBitmap()
	if err != nil {
		return nil, err
	}
	msgPacked.Write(encodeField)

	encodeField, err = m.packFields()
	if err != nil {
		return nil, err
	}
	msgPacked.Write(encodeField)

	return msgPacked.Bytes(), nil
}

func (m *Message) LogMsg() string {
	fieldsToLog := make(map[string]interface{})

	if m.Bitmap != nil {
		for _, id := range m.Bitmap.GetSliceString() {
			if f, ok := m.fields[id]; ok {
				var err error
				fieldsToLog[strconv.Itoa(id)], err = f.Log()
				if err != nil {
					log.Printf("error logging field %d: %v", id, err)
				}
			}
		}
	}

	jsonBytes, err := json.Marshal(fieldsToLog)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"failed to convert log to JSON: %v\"}", err)
	}

	return string(jsonBytes)
}

func (m *Message) packMti() ([]byte, error) {
	fldPkg, ok := m.Packager.Fields[0]
	if !ok {
		return nil, nil
		return nil, ErrMTINotFoundInPackager
	}

	fld, err := m.GetFieldString(0)
	if err != nil {
		return nil, fmt.Errorf("pack mti: %w", err)
	}

	encodeField, _, err := fldPkg.Pack(fld)
	if err != nil {
		return nil, fmt.Errorf("pack mti: %w", err)
	}

	return encodeField, nil
}

func (m *Message) packBitmap() ([]byte, error) {
	if m.Packager.Bitmap != nil {
		// El bitmap se maneja como BytesField, su String() devuelve hex
		if len(m.Bitmap.ToBytes()) > m.Packager.Bitmap.Length() {
			// Si el bitmap es secundario, SetField(1, ...) lo actualizará
			err := m.SetField(1, m.Bitmap.ToBytes()[m.Packager.Bitmap.Length():])
			if err != nil {
				return nil, err
			}
		}

		encodeField, _, errPack := m.Packager.Bitmap.Pack(m.Bitmap.ToString())
		if errPack != nil {
			return nil, fmt.Errorf("pack bitmap: %w", errPack)
		}

		return encodeField, nil
	}
	return nil, ErrBitmapNotFoundInPackager
}

func (m *Message) packFields() ([]byte, error) {
	fieldsPacked := new(bytes.Buffer)

	for _, k := range m.Bitmap.GetSliceString() {
		if k == 0 || k == 1 {
			continue
		}

		fldPkg, ok := m.Packager.Fields[k]
		if !ok {
			return nil, fmt.Errorf("field %d: %w", k, ErrNotFoundInPackager)
		}

		fld, ok := m.fields[k]
		if !ok {
			return nil, fmt.Errorf("field %d not set in message", k)
		}

		fldStr, err := fld.Parse()
		if err != nil {
			return nil, fmt.Errorf("pack field %d: could not get string value: %w", k, err)
		}

		encodedField, _, errPack := fldPkg.Pack(fldStr)
		if errPack != nil {
			return nil, fmt.Errorf("pack field %d: %w", k, errPack)
		}
		fieldsPacked.Write(encodedField)
	}

	return fieldsPacked.Bytes(), nil
}

func (m *Message) unpackMti(messageRaw []byte) (int, error) {
	if fldPkg, ok := m.Packager.Fields[0]; ok {
		value, length, err := fldPkg.Unpack(messageRaw, 0)
		if err != nil {
			return 0, fmt.Errorf("unpack MTI: %w", err)
		}
		if err := m.SetField(0, value); err != nil {
			return 0, fmt.Errorf("unpack MTI: failed to set field: %w", err)
		}
		return length, nil
	}
	return 0, nil
	return 0, ErrMTINotFoundInPackager
}

func (m *Message) unpackBitmap(messageRaw []byte, offset int) (int, error) {
	if fldPkg, ok := m.Packager.Fields[1]; ok {
		bMap, length, err := bitmap.Unpack(fldPkg, messageRaw, offset)
		if err != nil {
			return 0, fmt.Errorf("unpack bitmap: %w", err)
		}

		m.Bitmap = bMap

		if len(m.Bitmap.ToBytes()) > fldPkg.Length() {
			err = m.SetField(1, bMap.ToBytes()[fldPkg.Length():])
			if err != nil {
				return 0, err
			}
		}

		return length, nil
	}

	return 0, ErrBitmapNotFoundInPackager
}

func (m *Message) unpackFields(messageRaw []byte, position int) error {
	for _, fieldId := range m.Bitmap.GetSliceString() {
		if fieldId != 0 && fieldId != 1 {
			if fldPkg, ok := m.Packager.Fields[fieldId]; ok {
				value, length, err := fldPkg.Unpack(messageRaw, position)
				if err != nil {
					return fmt.Errorf("unpack field %d: %w", fieldId, err)
				}

				f := m.createField(fieldId)
				if err := f.Unparse(value); err != nil {
					return fmt.Errorf("unpack field %d: failed to set bytes: %w", fieldId, err)
				}
				m.fields[fieldId] = f
				position += length
			} else {
				return fmt.Errorf("field %d: %w", fieldId, ErrNotFoundInPackager)
			}
		}
	}
	return nil
}

func (m *Message) Get() any { return m }

func (m *Message) Set(v any) error {
	val, ok := v.(*Message)
	if !ok {
		return fmt.Errorf("invalid type for Message field: received %T, expected *Message", v)
	}
	m.Packager = val.Packager
	m.Length = val.Length
	m.Header = val.Header
	m.Trailer = val.Trailer
	m.Bitmap = val.Bitmap
	m.fields = val.fields
	return nil
}

func (m *Message) Unparse(data string) error {
	rawBytes, err := hex.DecodeString(data)
	if err != nil {
		return fmt.Errorf("invalid hex data for nested message: %w", err)
	}
	return m.Unpack(rawBytes)
}

func (m *Message) Parse() (string, error) {
	packedBytes, err := m.Pack()
	if err != nil {
		return "", fmt.Errorf("failed to pack nested message: %w", err)
	}
	return hex.EncodeToString(packedBytes), nil
}

func (m *Message) Log() (interface{}, error) { return m.LogMsg(), nil }
