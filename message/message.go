// Package message provides functionalities for packing and unpacking
// ISO 8583 messages.
package message

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/bitmap"
	"github.com/tomasdemarco/iso8583/packager"
)

// Message represents an ISO 8583 message, containing its structure,
// fields, and the associated packager.
type Message struct {
	Packager *packager.Packager
	Length   int
	Header   interface{}
	Trailer  interface{}
	Bitmap   []string
	Fields   map[string]string
	TagsEmv  map[string]string
}

// NewMessage creates and returns a new Message instance
// initialized with the provided packager.
func NewMessage(packager *packager.Packager) *Message {
	return &Message{
		Packager: packager,
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

	for _, fieldId := range m.Bitmap {
		if fieldId != "001" {
			value, length, err := m.Packager.Fields[fieldId].Unpack(messageRaw, position)
			if err != nil {
				return fmt.Errorf("unpack field %s: %w", fieldId, err)
			}

			m.SetField(fieldId, value)
			position += length
		}
	}

	return nil
}

// Pack packs the message fields into an ISO 8583 byte slice.
// It calculates the bitmap and encodes each field according to the packager's configuration.
// It returns the packed message as a byte slice, and an error if packing fails.
func (m *Message) Pack() ([]byte, error) {

	bitmapSlice, bitmapString, err := bitmap.Pack(m.Fields)
	if err != nil {
		return nil, err
	}

	m.SetField("001", fmt.Sprintf("%x", bitmapString))

	m.Bitmap = append(bitmapSlice[:1], append([]string{"001"}, bitmapSlice[1:]...)...)

	msgPacked := new(bytes.Buffer)
	for _, k := range m.Bitmap {
		value := m.Fields[k]

		encodeField, plainField, errPack := m.Packager.Fields[k].Pack(value)
		if errPack != nil {
			return nil, fmt.Errorf("pack field %s: %w", k, errPack)
		}
		m.SetField(k, plainField)

		msgPacked.Write(encodeField)
	}

	return msgPacked.Bytes(), err
}

// SetField sets the value of a specific field in the message.
// If the fields map not initialized, it creates it.
func (m *Message) SetField(fieldId string, value string) {
	if m.Fields == nil {
		var fields = make(map[string]string)
		m.Fields = fields
	}
	m.Fields[fieldId] = value
}

// GetField retrieves the value of a specific field from the message.
// It returns the fields value as a string, and an error if the field does not exist.
func (m *Message) GetField(fieldId string) (string, error) {
	if fld, ok := m.Fields[fieldId]; ok {
		return fld, nil
	}

	return "", fmt.Errorf(`the message does not contain the field with the id "%s"`, fieldId)
}

// unpackMti unpacks the Message Type Indicator (MTI) from the message.
// This is an internal helper method.
func (m *Message) unpackMti(messageRaw []byte) (int, error) {
	if _, ok := m.Packager.Fields["000"]; !ok {
		return 0, errors.New("packager does not contain field 000")
	}

	value, length, err := m.Packager.Fields["000"].Unpack(messageRaw, 0)
	if err != nil {
		return 0, fmt.Errorf("pack field 000: %w", err)
	}

	if !m.Packager.Fields["000"].Pattern.MatchString(value) {
		return 0, errors.New("invalid format in field 000")
	}

	m.SetField("000", value)

	return length, nil
}

// unpackBitmap unpacks the bitmap from the message.
// This is an internal helper method.
func (m *Message) unpackBitmap(messageRaw []byte, offset int) (int, error) {

	if _, ok := m.Packager.Fields["001"]; !ok {
		return 0, errors.New("packager does not contain field 001")
	}

	if len(messageRaw) < offset+m.Packager.Fields["001"].Length {
		return 0, errors.New("the message is too short to be unpacked")
	}

	length, sliceBitmap, err := bitmap.Unpack(m.Packager.Fields["001"], messageRaw, offset)
	if err != nil {
		return 0, fmt.Errorf("could not get bitmap: %w", err)
	}

	m.Bitmap = sliceBitmap

	m.Packager.Fields["001"].Encoding.SetLength(length)
	value, err := m.Packager.Fields["001"].Encoding.Decode(messageRaw[offset:])
	if err != nil {
		return 0, fmt.Errorf("unpack field 001: %w", err)
	}

	if !m.Packager.Fields["001"].Pattern.MatchString(value) {
		return 0, fmt.Errorf("invalid format in field 001")
	}

	m.SetField("001", value)

	return length, nil
}
