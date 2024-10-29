package message

import (
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/packager"
	"regexp"
	"sort"
)

type Message struct {
	Packager *packager.Packager
	Length   int
	Header   map[string]string
	Bitmap   []string
	Fields   map[string]Field
	TagsEmv  map[string]string
}

type Field struct {
	Value     string
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
	position := length

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

	lengthBitmap, sliceBitmap, err := m.UnpackBitmap(position, messageRaw)
	if err != nil {
		err = errors.New("could not get bitmap, " + err.Error())
		return err
	}
	fmt.Println("LEN BITMAP", lengthBitmap)
	m.Bitmap = sliceBitmap

	m.SetField("001", messageRaw[position:position+lengthBitmap])

	position += lengthBitmap

	match, _ = regexp.MatchString(m.Packager.Fields["001"].Pattern, m.Fields["001"].Value)
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

	bitmap, err := m.PackBitmap()
	if err != nil {
		return message, err
	}

	m.SetField("001", bitmap)

	keys := make([]string, 0, len(m.Fields))
	for k := range m.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fieldsAux := m.Fields
	m.Fields = nil

	for _, k := range keys {
		if fieldsAux[k].SubFields != nil {
			m.PackSubfields(fieldsAux, k)
		} else {
			m.SetField(k, fieldsAux[k].Value)
		}
	}

	for i := 0; i < len(keys); i++ {
		fieldEncode := m.PackField(keys[i])
		message += fieldEncode
	}

	return message, err
}
