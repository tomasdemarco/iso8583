package message

import (
	"errors"
	"github.com/tomasdemarco/iso8583/encoding"
)

func (m *Message) UnpackEncoding(messageRaw string, field string, position int, length int) (value string, doubleLength int, err error) {
	switch m.Packager.Fields[field].Encoding {
	case "ASCII":
		if len(messageRaw) < position+(length*2) {
			err = errors.New("index out of range while trying to unpack field " + field)
			return "", 0, err
		}
		result, err := encoding.AsciiDecode(messageRaw[position : position+(length*2)])
		return result, 2, err
	case "HEX":
		if m.Packager.Fields[field].Type == "STRING" {
			if len(messageRaw) < position+length {
				err = errors.New("index out of range while trying to unpack field " + field)
				return "", 0, err
			}
			result := encoding.HexDecode(value)
			return result, 2, nil
		} else {
			if len(messageRaw) < position+length {
				err = errors.New("index out of range while trying to unpack field " + field)
				return "", 0, err
			}
			result := encoding.HexDecode(value)
			return result, 1, nil
		}
	case "EBCDIC":
		if len(messageRaw) < position+(length*2) {
			err = errors.New("index out of range while trying to unpack field " + field)
			return "", 0, err
		}
		result, _ := encoding.EbcdicDecode(messageRaw[position : position+(length*2)])
		return result, 2, nil
	default:
		if len(messageRaw) < position+length {
			err = errors.New("index out of range while trying to unpack field " + field)
			return "", 0, err
		}
		return messageRaw[position : position+length], 1, nil
	}
}

func (m *Message) PackEncoding(field string, padRight string, padLeft string) (value string) {
	switch m.Packager.Fields[field].Encoding {
	case "ASCII":
		result := encoding.AsciiEncode(padLeft + m.FieldAndSubFields[field].Field + padRight)
		return result
	case "HEX":
		result := encoding.HexEncode(padLeft + m.FieldAndSubFields[field].Field + padRight)
		return result
	case "EBCDIC":
		result, _ := encoding.EbcdicEncode(padLeft + m.FieldAndSubFields[field].Field + padRight)
		return result
	default:
		return padLeft + m.FieldAndSubFields[field].Field + padRight
	}
}

func (m *Message) PackSubFieldEncoding(fieldsAux map[string]Fields, field string, subfield string) string {
	switch m.Packager.Fields[field].SubFields[subfield].Encoding {
	case "ASCII":
		result := encoding.AsciiEncode(fieldsAux[field].SubFields[subfield])
		return result
	case "HEX":
		result := encoding.HexEncode(fieldsAux[field].SubFields[subfield])
		return result
	case "EBCDIC":
		result, _ := encoding.EbcdicEncode(fieldsAux[field].SubFields[subfield])
		return result
	default:
		return fieldsAux[field].SubFields[subfield]
	}
}
