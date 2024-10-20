package message

import (
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/packager"
	"regexp"
	"sort"
)

type Message struct {
	Packager          *packager.Packager
	Length            int
	Header            map[string]string
	Bitmap            []string
	FieldAndSubFields map[string]Fields
	TagsEmv           map[string]string
}

type Fields struct {
	Field     string
	SubFields map[string]string
}

func NewMessage(packager *packager.Packager) *Message {
	return &Message{
		Packager: packager,
	}
}

func (m *Message) Unpack(messageRaw string) (err error) {
	if _, ok := m.Packager.Fields["000"]; !ok {
		err = errors.New("packager does not contain field 000")
		return err
	}

	length, err := m.UnpackField(messageRaw, 0, "000")
	if err != nil {
		return err
	}
	positionInitial := length

	match, _ := regexp.MatchString(m.Packager.Fields["000"].Pattern, m.FieldAndSubFields["000"].Field)
	if !match {
		err = errors.New("invalid format in field 000")
		return err
	}

	if _, ok := m.Packager.Fields["001"]; !ok {
		err = errors.New("packager does not contain field 001")
		return err
	}

	if len(messageRaw) < positionInitial+m.Packager.Fields["001"].Length {
		err = errors.New("the message is too short to be unpacked")
		return err
	}

	numberBitmaps, sliceBitmap, err := m.UnpackBitmap(positionInitial, messageRaw)
	if err != nil {
		err = errors.New("could not get bitmap")
		return err
	}
	m.Bitmap = sliceBitmap

	position := positionInitial + m.Packager.Fields["001"].Length*numberBitmaps
	m.SetField("001", messageRaw[positionInitial:positionInitial+m.Packager.Fields["001"].Length*numberBitmaps])
	match, _ = regexp.MatchString(m.Packager.Fields["001"].Pattern, m.FieldAndSubFields["001"].Field)
	if !match {
		err = errors.New("invalid format in field 001")
		return err
	}

	for _, field := range m.Bitmap {
		if field != "001" {
			length, err := m.UnpackField(messageRaw, position, field)
			if err != nil {
				return err
			}
			position += length
		}
	}
	return nil
}

func (m *Message) Pack() (message string, err error) {
	sliceBitmap := make([]string, 0)
	for k := range m.FieldAndSubFields {
		str := fmt.Sprintf("%03s", k)
		sliceBitmap = append(sliceBitmap, str)
	}
	sort.Strings(sliceBitmap)
	m.Bitmap = sliceBitmap
	m.SetField("001", m.PackBitmap())

	keys := make([]string, 0, len(m.FieldAndSubFields))
	for k := range m.FieldAndSubFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fieldsAux := m.FieldAndSubFields
	m.FieldAndSubFields = nil
	for _, k := range keys {
		if fieldsAux[k].SubFields != nil {
			m.PackSubfields(fieldsAux, k)
		} else {
			m.SetField(k, fieldsAux[k].Field)
		}
	}

	for i := 0; i < len(keys); i++ {
		fieldEncode := m.PackField(keys[i])
		message += fieldEncode
	}

	return message, err
}
