package message

import (
	"errors"
	"fmt"
	"regexp"
)

func (m *Message) UnpackField(messageRaw string, position int, field string) (length int, err error) {
	switch m.Packager.Fields[field].Prefix {
	case "LL":
		if len(messageRaw) < position+2 {
			err = errors.New("index out of range while trying to unpack prefix field " + field)
			return 0, err
		}

		length, lengthPrefix := m.UnpackPrefix(field, messageRaw, position, 2)
		if len(messageRaw) < position+lengthPrefix+length {
			err = errors.New("index out of range while trying to unpack field " + field)
			return 0, err
		}

		paddingRight, paddingLeft := m.UnpackPadding(field)

		value, doubleLength, err := m.UnpackEncoding(messageRaw, field, position+lengthPrefix, length)
		if err != nil {
			return 0, err
		}

		match, _ := regexp.MatchString(m.Packager.Fields[field].Pattern, value)
		if !match {
			err = errors.New("invalid format in field " + field)
			return 0, err
		}

		m.SetField(field, value)

		if m.Packager.Fields[field].SubFields != nil {
			m.UnpackSubfields(field, value)
		}

		if length%2 != 0 {
			length += +paddingRight + paddingLeft
		}
		length = length*doubleLength + lengthPrefix

		return length, err
	case "LLL":
		if len(messageRaw) < position+4 {
			err = errors.New("index out of range while trying to unpack prefix field " + field)
			return 0, err
		}

		length, lengthPrefix := m.UnpackPrefix(field, messageRaw, position, 4)
		if len(messageRaw) < position+lengthPrefix+length {
			err = errors.New("index out of range while trying to unpack field " + field)
			return 0, err
		}

		paddingRight, paddingLeft := m.UnpackPadding(field)

		value, doubleLength, err := m.UnpackEncoding(messageRaw, field, position+lengthPrefix, length)
		if err != nil {
			return 0, err
		}

		match, _ := regexp.MatchString(m.Packager.Fields[field].Pattern, value)
		if !match {
			err = errors.New("invalid format in field " + field)
			return 0, err
		}

		m.SetField(field, value)

		if m.Packager.Fields[field].SubFields != nil {
			m.UnpackSubfields(field, value)
		}

		if length%2 != 0 {
			length += +paddingRight + paddingLeft
		}
		length = length*doubleLength + lengthPrefix

		return length, err
	default:
		paddingRight, paddingLeft := m.UnpackPadding(field)

		value, doubleLength, err := m.UnpackEncoding(messageRaw, field, position+paddingLeft, m.Packager.Fields[field].Length)
		if err != nil {
			return 0, err
		}

		match, _ := regexp.MatchString(m.Packager.Fields[field].Pattern, value)
		if !match {
			err = errors.New("invalid format in field " + field)
			return 0, err
		}

		m.SetField(field, value)

		if m.Packager.Fields[field].SubFields != nil {
			m.UnpackSubfields(field, value)
		}

		length = m.Packager.Fields[field].Length*doubleLength + paddingRight + paddingLeft

		return length, err
	}
}

func (m *Message) PackField(field string) (fieldEncode string) {
	switch m.Packager.Fields[field].Prefix {
	case "LL":
		fieldPrefix := m.PackPrefix(field, len(m.Fields[field].Value), 2)

		padRight, padLeft := m.PackPadding(field)
		if m.Packager.Fields[field].Padding.Type != "PARITY" {
			m.SetField(field, padLeft+m.Fields[field].Value+padRight)
			fieldEncode = m.PackEncoding(field, "", "")
		} else {
			fieldEncode = m.PackEncoding(field, padRight, padLeft)
		}

		return fmt.Sprintf("%02s", fieldPrefix) + fieldEncode
	case "LLL":
		fieldPrefix := m.PackPrefix(field, len(m.Fields[field].Value), 4)

		padRight, padLeft := m.PackPadding(field)
		if m.Packager.Fields[field].Padding.Type != "PARITY" {
			m.SetField(field, padLeft+m.Fields[field].Value+padRight)
			fieldEncode = m.PackEncoding(field, "", "")
		} else {
			fieldEncode = m.PackEncoding(field, padRight, padLeft)
		}

		return fmt.Sprintf("%04s", fieldPrefix) + fieldEncode
	default:
		padRight, padLeft := m.PackPadding(field)
		if m.Packager.Fields[field].Padding.Type != "PARITY" {
			m.SetField(field, padLeft+m.Fields[field].Value+padRight)
			fieldEncode = m.PackEncoding(field, "", "")
		} else {
			fieldEncode = m.PackEncoding(field, padRight, padLeft)
		}

		return fieldEncode
	}
}

func (m *Message) SetField(field string, value string) {
	if m.Fields == nil {
		var fields = make(map[string]Field)
		m.Fields = fields
	}
	fieldAux := m.Fields[field]
	fieldAux.Value = value
	m.Fields[field] = fieldAux
}

func (m *Message) GetField(field string) (value string, err error) {
	if _, ok := m.Fields[field]; ok {
		return m.Fields[field].Value, nil
	}
	err = errors.New("the message does not contain the field with the id '" + field + "'")
	return value, err
}
