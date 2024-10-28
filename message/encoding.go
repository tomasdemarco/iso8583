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
		result := encoding.AsciiEncode(padLeft + m.Fields[field].Value + padRight)
		return result
	case "EBCDIC":
		result, _ := encoding.EbcdicEncode(padLeft + m.Fields[field].Value + padRight)
		return result
	default:
		return padLeft + m.Fields[field].Value + padRight
	}
}

func (m *Message) PackSubFieldEncoding(fieldsAux map[string]Field, field string, subfield string) string {
	switch m.Packager.Fields[field].SubFields[subfield].Encoding {
	case "ASCII":
		result := encoding.AsciiEncode(fieldsAux[field].SubFields[subfield])
		return result
	case "EBCDIC":
		result, _ := encoding.EbcdicEncode(fieldsAux[field].SubFields[subfield])
		return result
	default:
		return fieldsAux[field].SubFields[subfield]
	}
}
