package field

import (
	"errors"
	"fmt"
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/encoding"
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/packager"
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/padding"
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/prefix"
	"regexp"
)

type Field struct {
	Value     string
	SubFields map[string]string
}

func Unpack(fieldPackager packager.Field, messageRaw string, position int, field string) (*string, *int, error) {

	if fieldPackager.Prefix.Type != prefix.Fixed {
		if len(messageRaw) < position+(fieldPackager.Prefix.Type.EnumIndex()*1) {
			err := errors.New("index out of range while trying to unpack prefix field " + field)
			return nil, nil, err
		}

		length, lengthPrefix, err := prefix.Unpack(fieldPackager.Prefix, messageRaw[position:])
		if err != nil {
			return nil, nil, err
		}

		if len(messageRaw) < position+lengthPrefix+length {
			err = errors.New("index out of range while trying to unpack field " + field)
			return nil, nil, err
		}

		paddingRight, paddingLeft := padding.Unpack(fieldPackager.Padding)

		value, doubleLength, err := encoding.Unpack(fieldPackager.Encoding, messageRaw, field, position+lengthPrefix, length)
		if err != nil {
			return nil, nil, err
		}

		match, _ := regexp.MatchString(fieldPackager.Pattern, value)
		if !match {
			err = errors.New("invalid format in field " + field)
			return nil, nil, err
		}

		//if fieldPackager.SubFields != nil {
		//	m.UnpackSubfields(field, value)
		//}

		if length%2 != 0 {
			length += +paddingRight + paddingLeft
		}
		length = length*doubleLength + lengthPrefix

		return &value, &length, err
	} else {
		paddingRight, paddingLeft := padding.Unpack(fieldPackager.Padding)

		value, doubleLength, err := encoding.Unpack(fieldPackager.Encoding, messageRaw, field, position+paddingLeft, fieldPackager.Length)
		if err != nil {
			return nil, nil, err
		}

		match, _ := regexp.MatchString(fieldPackager.Pattern, value)
		if !match {
			err = errors.New("invalid format in field " + field)
			return nil, nil, err
		}

		//if fieldPackager.SubFields != nil {
		//	m.UnpackSubfields(field, value)
		//}

		length := fieldPackager.Length*doubleLength + paddingRight + paddingLeft

		return &value, &length, err
	}
}

func Pack(fieldPackager packager.Field, value string) string {
	fieldPrefix, _ := prefix.Pack(fieldPackager.Prefix, len(value))

	padRight, padLeft := padding.Pack(fieldPackager.Padding, len(value), fieldPackager.Length)

	fieldEncode := encoding.Pack(fieldPackager.Encoding, padLeft+value+padRight)

	switch fieldPackager.Prefix.Type {
	case prefix.LL:
		return fmt.Sprintf("%02s", fieldPrefix) + fieldEncode
	case prefix.LLL:
		return fmt.Sprintf("%04s", fieldPrefix) + fieldEncode
	default:
		return fieldEncode
	}
}
