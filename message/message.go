// Package message provides functionalities for packing and unpacking
// ISO 8583 messages.
package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tomasdemarco/iso8583/bitmap"
	"github.com/tomasdemarco/iso8583/packager"
	"github.com/tomasdemarco/iso8583/utils"
)

// Message represents an ISO 8583 message, containing its structure,
// fields, and the associated packager.
type Message struct {
	Packager *packager.Packager
	Length   int
	Header   interface{}
	Trailer  interface{}
	Bitmap   *utils.BitSet
	Fields   map[int]string
	TagsEmv  map[string]string
}

// NewMessage creates and returns a new Message instance
// initialized with the provided packager.
func NewMessage(packager *packager.Packager) *Message {
	return &Message{
		Packager: packager,
		Bitmap:   utils.NewBitSet(64, 128),
	}
}

// SetField sets the value of a specific field in the message.
// If the fields map not initialized, it creates it.
func (m *Message) SetField(fieldId int, value string) {
	if m.Fields == nil {
		var fields = make(map[int]string)
		m.Fields = fields
	}

	m.Fields[fieldId] = value

	if !m.Bitmap.Get(fieldId) {
		m.Bitmap.Set(fieldId)
	}
}

// GetField retrieves the value of a specific field from the message.
// It returns the fields value as a string, and an error if the field does not exist.
func (m *Message) GetField(fieldId int) (string, error) {
	if fld, ok := m.Fields[fieldId]; ok {
		return fld, nil
	}

	return "", fmt.Errorf("field %d: %w", fieldId, ErrNotFoundInMessage)
}

// Unpack unpacks a byte slice of an ISO 8583 message
// into the Message structure, populating its fields.
// It returns an error if unpacking fails.
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
// It calculates the bitmap and encodes each field according to the packager's configuration.
// It returns the packed message as a byte slice, and an error if packing fails.
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

func (m *Message) Log() string {

	fieldsToLog := make(map[string]interface{})
	value, err := m.GetField(0)
	if err == nil {
		fieldsToLog["0"] = value
	}

	fieldsToLog["1"] = m.Bitmap.ToString()

	for _, fieldID := range m.Bitmap.GetSliceString() {
		if fieldID != 1 {
			value, err := m.GetField(fieldID)
			if err == nil {
				fieldsToLog[fmt.Sprintf("%d", fieldID)] = value
			}
		}
	}

	jsonBytes, _ := json.Marshal(fieldsToLog)

	return string(jsonBytes)
}

func (m *Message) packMti() ([]byte, error) {
	if _, ok := m.Packager.Fields[0]; ok {
		encodeField, _, errPack := m.Packager.Fields[0].Pack(m.Fields[0])
		if errPack != nil {
			return nil, fmt.Errorf("pack mti: %w", errPack)
		}

		return encodeField, nil
	}

	return nil, ErrMTINotFoundInPackager
}

func (m *Message) packBitmap() ([]byte, error) {
	if fldPKg, ok := m.Packager.Fields[1]; ok {

		if len(m.Bitmap.ToBytes()) > 8 {
			m.SetField(1, fmt.Sprintf("%X", m.Bitmap.ToBytes()))
		}

		encodeField, _, errPack := fldPKg.Pack(m.Bitmap.ToString())
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
		if k != 1 {
			if fldPkg, ok := m.Packager.Fields[k]; ok {
				encodeField, plainField, errPack := fldPkg.Pack(m.Fields[k])
				if errPack != nil {
					return nil, fmt.Errorf("pack field %d: %w", k, errPack)
				}

				m.SetField(k, plainField)
				fieldsPacked.Write(encodeField)
			} else {
				return nil, fmt.Errorf("field %d: %w", k, ErrNotFoundInPackager)
			}
		}
	}

	return fieldsPacked.Bytes(), nil
}

// unpackMti unpacks the Message Type Indicator (MTI) from the message.
// This is an internal helper method.
func (m *Message) unpackMti(messageRaw []byte) (int, error) {
	if fldPkg, ok := m.Packager.Fields[0]; ok {
		value, length, err := fldPkg.Unpack(messageRaw, 0)
		if err != nil {
			return 0, fmt.Errorf("unpack MTI: %w", err)
		}

		m.SetField(0, value)

		return length, nil
	}

	return 0, ErrMTINotFoundInPackager
}

// unpackBitmap unpacks the bitmap from the message.
// This is an internal helper method.
func (m *Message) unpackBitmap(messageRaw []byte, offset int) (int, error) {
	if fldPkg, ok := m.Packager.Fields[1]; ok {
		bMap, length, err := bitmap.Unpack(fldPkg, messageRaw, offset)
		if err != nil {
			return 0, fmt.Errorf("unpack bitmap: %w", err)
		}

		m.Bitmap = bMap
		m.SetField(1, bMap.ToString())

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

				m.SetField(fieldId, value)
				position += length
			} else {
				return fmt.Errorf("field %d: %w", fieldId, ErrNotFoundInPackager)
			}
		}
	}
	return nil
}
