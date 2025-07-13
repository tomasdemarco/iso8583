// Package message provides functionalities for packing and unpacking
// ISO 8583 messages.
package message

import (
	"bytes"
	"encoding/json"
	"errors"
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
	//Bitmap   []int
	Bitmap2 *utils.BitSet
	Fields  map[int]string
	TagsEmv map[string]string
}

// NewMessage creates and returns a new Message instance
// initialized with the provided packager.
func NewMessage(packager *packager.Packager) *Message {
	return &Message{
		Packager: packager,
		Bitmap2:  utils.NewBitSet(64, 128),
	}
}

// Unpack unpacks a byte slice of an ISO 8583 message
// into the Message structure, populating its fields.
// It returns an error if unpacking fails.
func (m *Message) Unpack(messageRaw []byte) (err error) {
	length, err := m.unpackMti(messageRaw)
	if err != nil {
		return err
	}

	position := length

	length, err = m.unpackBitmap(messageRaw, position)
	if err != nil {
		return err
	}

	position += length

	for _, fieldId := range m.Bitmap2.GetSliceString() {
		if fieldId != 0 && fieldId != 1 {
			if _, ok := m.Packager.Fields[fieldId]; !ok {
				return fmt.Errorf("packager does not contain field %d", fieldId)
			}

			if fldPkg, ok := m.Packager.Fields[fieldId]; ok {
				value, length, err := fldPkg.Unpack(messageRaw, position)
				if err != nil {
					return fmt.Errorf("unpack field %d: %w", fieldId, err)
				}

				m.SetField(fieldId, value)
				position += length
			} else {
				return fmt.Errorf("unpack field %d: field packager not found", fieldId)
			}
		}
	}

	return nil
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

	for _, k := range m.Bitmap2.GetSliceString() {
		if _, ok := m.Packager.Fields[k]; !ok {
			return nil, fmt.Errorf("packager does not contain field %d", k)
		}

		value := m.Fields[k]

		encodeField, plainField, errPack := m.Packager.Fields[k].Pack(value)
		if errPack != nil {
			return nil, fmt.Errorf("pack field %d: %w", k, errPack)
		}
		m.SetField(k, plainField)

		msgPacked.Write(encodeField)
	}

	return msgPacked.Bytes(), nil
}

// SetField sets the value of a specific field in the message.
// If the fields map not initialized, it creates it.
func (m *Message) SetField(fieldId int, value string) {
	if m.Fields == nil {
		var fields = make(map[int]string)
		m.Fields = fields
	}

	m.Fields[fieldId] = value

	if !m.Bitmap2.Get(fieldId) {
		m.Bitmap2.Set(fieldId)
	}
}

// GetField retrieves the value of a specific field from the message.
// It returns the fields value as a string, and an error if the field does not exist.
func (m *Message) GetField(fieldId int) (string, error) {
	if fld, ok := m.Fields[fieldId]; ok {
		return fld, nil
	}

	return "", fmt.Errorf(`the message does not contain the field with the id "%d"`, fieldId)
}

func (m *Message) Log() string {

	jsonStr := "{"
	v, err := m.GetField(0)
	if err != nil {
		return ""
	}

	marshaledValue, _ := json.Marshal(v)
	jsonStr += fmt.Sprintf(`"%d":%s`, 0, string(marshaledValue))
	jsonStr += ","

	marshaledValue, _ = json.Marshal(v)
	jsonStr += fmt.Sprintf(`"%d":%s`, 1, m.Bitmap2.ToString())
	jsonStr += ","

	for i, k := range m.Bitmap2.GetSliceString() {
		v, err = m.GetField(k)
		if err != nil {
			return ""
		}

		marshaledValue, _ = json.Marshal(v)
		jsonStr += fmt.Sprintf(`"%d":%s`, k, string(marshaledValue))
		if i < len(m.Bitmap2.GetSliceString())-1 {
			jsonStr += ","
		}
	}
	jsonStr += "}"

	return jsonStr
}

func (m *Message) packMti() ([]byte, error) {
	if _, ok := m.Packager.Fields[0]; !ok {
		return nil, fmt.Errorf("packager does not contain field %d", 0)
	}

	encodeField, _, errPack := m.Packager.Fields[0].Pack(m.Fields[0])
	if errPack != nil {
		return nil, fmt.Errorf("pack mti: %w", errPack)
	}

	return encodeField, nil
}

func (m *Message) packBitmap() ([]byte, error) {
	if _, ok := m.Packager.Fields[1]; !ok {
		return nil, fmt.Errorf("packager does not contain field %d", 1)
	}

	if len(m.Bitmap2.ToBytes()) > 8 {
		m.SetField(1, fmt.Sprintf("%X", m.Bitmap2.ToBytes()))
	}

	encodeField, _, errPack := m.Packager.Fields[1].Pack(m.Bitmap2.ToString())
	if errPack != nil {
		return nil, fmt.Errorf("pack bitmap: %w", errPack)
	}

	return encodeField, nil
}

// unpackMti unpacks the Message Type Indicator (MTI) from the message.
// This is an internal helper method.
func (m *Message) unpackMti(messageRaw []byte) (int, error) {
	if _, ok := m.Packager.Fields[0]; !ok {
		return 0, errors.New("packager does not contain field 0")
	}

	value, length, err := m.Packager.Fields[0].Unpack(messageRaw, 0)
	if err != nil {
		return 0, fmt.Errorf("pack field 0: %w", err)
	}

	m.SetField(0, value)

	return length, nil
}

// unpackBitmap unpacks the bitmap from the message.
// This is an internal helper method.
func (m *Message) unpackBitmap(messageRaw []byte, offset int) (int, error) {
	if _, ok := m.Packager.Fields[1]; !ok {
		return 0, errors.New("packager does not contain field 1")
	}

	bMap, length, err := bitmap.Unpack(m.Packager.Fields[1], messageRaw, offset)
	if err != nil {
		return 0, fmt.Errorf("could not get bitmap: %w", err)
	}

	m.Bitmap2 = bMap
	m.SetField(1, bMap.ToString())

	return length, nil
}
