package encoding

import (
	"errors"
)

func Unpack(encoding Encoding, messageRaw string, field string, position int, length int) (value string, doubleLength int, err error) {
	switch encoding {
	case Ascii:
		if len(messageRaw) < position+(length*2) {
			err = errors.New("index out of range while trying to unpack field " + field)
			return "", 0, err
		}
		result, err := AsciiDecode(messageRaw[position : position+(length*2)])
		return result, 2, err
	case Ebcdic:
		if len(messageRaw) < position+(length*2) {
			err = errors.New("index out of range while trying to unpack field " + field)
			return "", 0, err
		}
		result, _ := EbcdicDecode(messageRaw[position : position+(length*2)])
		return result, 2, nil
	case Ans:
		if len(messageRaw) < position+(length*2) {
			err = errors.New("index out of range while trying to unpack field " + field)
			return "", 0, err
		}
		return messageRaw[position : position+(length*2)], 2, nil
	default:
		if len(messageRaw) < position+length {
			err = errors.New("index out of range while trying to unpack field " + field)
			return "", 0, err
		}
		return messageRaw[position : position+length], 1, nil
	}
}

func Pack(encoding Encoding, value string) string {
	switch encoding {
	case Ascii:
		result := AsciiEncode(value)
		return result
	case Ebcdic:
		result, _ := EbcdicEncode(value)
		return result
	default:
		return value
	}
}

func PackSubField(encoding Encoding, value string) string {
	switch encoding {
	case Ascii:
		result := AsciiEncode(value)
		return result
	case Ebcdic:
		result, _ := EbcdicEncode(value)
		return result
	default:
		return value
	}
}
