package message

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/bitmap"
	"github.com/tomasdemarco/iso8583/field"
	"github.com/tomasdemarco/iso8583/packager"
	"regexp"
	"sort"
)

type Message struct {
	Packager *packager.Packager
	Length   int
	Header   map[string]string
	Bitmap   []string
	Fields   map[string]field.Field
	TagsEmv  map[string]string
}

func NewMessage(packager *packager.Packager) *Message {
	return &Message{
		Packager: packager,
	}
}

func (m *Message) Unpack(messageRaw []byte) (err error) {
	if _, ok := m.Packager.Fields["000"]; !ok {
		err = errors.New("packager does not contain field 000")
		return err
	}

	value, length, err := field.Unpack(m.Packager.Fields["000"], messageRaw, 0, "000")
	if err != nil {
		return err
	}

	m.SetField("000", *value)
	position := *length

	match, _ := regexp.MatchString(m.Packager.Fields["000"].Pattern, m.Fields["000"].Value)
	if !match {
		err = errors.New("invalid format in field 000")
		return err
	}

	if _, ok := m.Packager.Fields["001"]; !ok {
		err = errors.New("packager does not contain field 001")
		return err
	}

	if len(messageRaw) < position+m.Packager.Fields["001"].Length {
		err = errors.New("the message is too short to be unpacked")
		return err
	}

	lengthBitmap, sliceBitmap, err := bitmap.Unpack(m.Packager.Fields["001"], position, messageRaw)
	if err != nil {
		err = errors.New("could not get bitmap, " + err.Error())
		return err
	}

	m.Bitmap = sliceBitmap

	m.SetField("001", fmt.Sprintf("%x", messageRaw[position:position+lengthBitmap]))

	position += lengthBitmap

	match, _ = regexp.MatchString(m.Packager.Fields["001"].Pattern, m.Fields["001"].Value)
	if !match {
		err = errors.New("invalid format in field 001")
		return err
	}

	for _, fieldId := range m.Bitmap {
		if fieldId != "001" {
			value, length, err = field.Unpack(m.Packager.Fields[fieldId], messageRaw, position, fieldId)
			if err != nil {
				return err
			}

			m.SetField(fieldId, *value)
			position += *length
		}
	}

	return nil
}

func (m *Message) Pack() ([]byte, error) {

	bitmapSlice, bitmapString, err := bitmap.Pack(m.Fields)
	if err != nil {
		return nil, err
	}

	m.Bitmap = bitmapSlice

	m.SetField("001", fmt.Sprintf("%x", bitmapString))

	keys := make([]string, 0, len(m.Fields))
	for k := range m.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fieldsAux := m.Fields
	m.Fields = nil

	for _, k := range keys {
		value := fieldsAux[k].Value

		if fieldsAux[k].SubFields != nil {
			fieldAux := fieldsAux[k]
			value = fmt.Sprintf("%x", fieldAux.PackSubfields(m.Packager.Fields[k].SubFields))
		}

		m.SetField(k, value)
	}

	msgPacked := new(bytes.Buffer)

	for _, k := range keys {
		fieldEncode, errPack := field.Pack(m.Packager.Fields[k], m.Fields[k].Value)
		if errPack != nil {
			return nil, errPack
		}
		//m.SetField(k, fieldEncode) //TODO creo que no hace falta
		msgPacked.Write(fieldEncode)
	}

	return msgPacked.Bytes(), err
}

func (m *Message) SetField(fieldId string, value string) {
	if m.Fields == nil {
		var fields = make(map[string]field.Field)
		m.Fields = fields
	}
	fieldAux := m.Fields[fieldId]
	fieldAux.Value = value
	m.Fields[fieldId] = fieldAux
}

func (m *Message) GetField(fieldId string) (value string, err error) {
	if _, ok := m.Fields[fieldId]; ok {
		return m.Fields[fieldId].Value, nil
	}

	return "", errors.New(fmt.Sprintf(`the message does not contain the field with the id "%s"`, fieldId))
}
